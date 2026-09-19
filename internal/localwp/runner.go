package localwp

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type CommandSpec struct {
	Binary string
	Args   []string
	Dir    string
	Env    []string
}

func BuildWPCommand(site Site, env Environment, args []string) CommandSpec {
	commandArgs := []string{"-c", env.PHPIni, env.WPCLIPhar, "--path=" + site.WebRoot}
	commandArgs = append(commandArgs, args...)

	processEnv := mergeEnv(os.Environ(), map[string]string{
		"PHPRC":                            filepath.Dir(env.PHPIni),
		"WP_ENV":                           "local-by-flywheel",
		"WP_CLI_DISABLE_AUTO_CHECK_UPDATE": "1",
	})

	if env.MySQLConfDir != "" && dirExists(env.MySQLConfDir) {
		processEnv = mergeEnv(processEnv, map[string]string{"MYSQL_HOME": env.MySQLConfDir})
	}
	if env.WPCLIConfig != "" {
		processEnv = mergeEnv(processEnv, map[string]string{"WP_CLI_CONFIG_PATH": env.WPCLIConfig})
	}

	var prepend []string
	if env.MySQLBinDir != "" {
		prepend = append(prepend, env.MySQLBinDir)
	}
	if env.PHPBinary != "" {
		prepend = append(prepend, filepath.Dir(env.PHPBinary))
	}
	if len(prepend) > 0 {
		current := envValue(processEnv, "PATH")
		newPath := strings.Join(append(prepend, current), string(os.PathListSeparator))
		processEnv = mergeEnv(processEnv, map[string]string{"PATH": newPath})
	}

	return CommandSpec{Binary: env.PHPBinary, Args: commandArgs, Dir: site.WebRoot, Env: processEnv}
}

func ExecuteWP(ctx context.Context, spec CommandSpec, stdin io.Reader, stdout, stderr io.Writer) int {
	cmd := exec.CommandContext(ctx, spec.Binary, spec.Args...)
	cmd.Dir = spec.Dir
	cmd.Env = spec.Env
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return exitErr.ExitCode()
		}
		fmt.Fprintf(stderr, "localwp: execute WP-CLI: %v\n", err)
		return 1
	}
	return 0
}

func OpenShell(site Site, dataDir string, stdin io.Reader, stdout, stderr io.Writer) int {
	entry, err := ShellEntryForSite(site, dataDir)
	if err != nil {
		fmt.Fprintf(stderr, "localwp: %v\n", err)
		return 1
	}

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd.exe", "/K", entry)
	} else {
		shell := os.Getenv("SHELL")
		if shell == "" {
			shell = "/bin/bash"
		}
		cmd = exec.Command(shell, entry)
	}
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return exitErr.ExitCode()
		}
		fmt.Fprintf(stderr, "localwp: open Local Site Shell: %v\n", err)
		return 1
	}
	return 0
}

func mergeEnv(base []string, updates map[string]string) []string {
	values := make(map[string]string, len(base)+len(updates))
	order := make([]string, 0, len(base)+len(updates))
	for _, item := range base {
		key, value, ok := strings.Cut(item, "=")
		if !ok {
			continue
		}
		normalized := key
		if runtime.GOOS == "windows" {
			normalized = strings.ToUpper(key)
		}
		if _, exists := values[normalized]; !exists {
			order = append(order, normalized)
		}
		values[normalized] = value
	}
	for key, value := range updates {
		normalized := key
		if runtime.GOOS == "windows" {
			normalized = strings.ToUpper(key)
		}
		if _, exists := values[normalized]; !exists {
			order = append(order, normalized)
		}
		values[normalized] = value
	}
	result := make([]string, 0, len(order))
	for _, key := range order {
		result = append(result, key+"="+values[key])
	}
	return result
}

func envValue(env []string, key string) string {
	for _, item := range env {
		k, value, ok := strings.Cut(item, "=")
		if ok && strings.EqualFold(k, key) {
			return value
		}
	}
	return ""
}
