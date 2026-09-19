package localwp

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestBuildWPCommandPreservesArguments(t *testing.T) {
	root := t.TempDir()
	site := Site{WebRoot: filepath.Join(root, "app", "public")}
	env := Environment{
		PHPBinary:    filepath.Join(root, "php"),
		PHPIni:       filepath.Join(root, "php.ini"),
		WPCLIPhar:    filepath.Join(root, "wp-cli.phar"),
		WPCLIConfig:  filepath.Join(root, "config.yaml"),
		MySQLBinDir:  filepath.Join(root, "mysql"),
		MySQLConfDir: filepath.Join(root, "mysql-conf"),
	}
	if err := os.MkdirAll(env.MySQLConfDir, 0o755); err != nil {
		t.Fatal(err)
	}

	input := []string{"search-replace", "http://old.local", "https://new.local", "--dry-run", "--format=json"}
	spec := BuildWPCommand(site, env, input)

	expected := []string{"-c", env.PHPIni, env.WPCLIPhar, "--path=" + site.WebRoot}
	expected = append(expected, input...)
	if !reflect.DeepEqual(spec.Args, expected) {
		t.Fatalf("args changed\nwant: %#v\n got: %#v", expected, spec.Args)
	}
	if spec.Binary != env.PHPBinary || spec.Dir != site.WebRoot {
		t.Fatalf("unexpected command spec: %#v", spec)
	}
	if !strings.Contains(envValue(spec.Env, "PATH"), env.MySQLBinDir) {
		t.Fatalf("expected DB binaries in PATH: %s", envValue(spec.Env, "PATH"))
	}
}

func TestParseWrapperArgsLeavesWPCLIArgsUntouched(t *testing.T) {
	selector, remaining, err := parseWrapperArgs([]string{"plugin", "list", "--local-site", "Demo", "--format=json"})
	if err != nil {
		t.Fatal(err)
	}
	if selector != "Demo" {
		t.Fatalf("selector: %q", selector)
	}
	expected := []string{"plugin", "list", "--format=json"}
	if !reflect.DeepEqual(remaining, expected) {
		t.Fatalf("want %#v, got %#v", expected, remaining)
	}
}

func TestExecuteWPPropagatesOutputAndExitCode(t *testing.T) {
	if os.Getenv("LOCALWP_TEST_HELPER") == "1" {
		for i, arg := range os.Args {
			if arg == "--" {
				_, _ = os.Stdout.WriteString(strings.Join(os.Args[i+1:], "|") + "\n")
				os.Exit(7)
			}
		}
		os.Exit(8)
	}

	var out strings.Builder
	var errOut strings.Builder
	spec := CommandSpec{
		Binary: os.Args[0],
		Args:   []string{"-test.run=TestExecuteWPPropagatesOutputAndExitCode", "--", "plugin", "list", "--format=json"},
		Dir:    t.TempDir(),
		Env:    append(os.Environ(), "LOCALWP_TEST_HELPER=1"),
	}
	code := ExecuteWP(context.Background(), spec, nil, &out, &errOut)
	if code != 7 {
		t.Fatalf("expected exit 7, got %d; stderr=%s", code, errOut.String())
	}
	if !strings.Contains(out.String(), "plugin|list|--format=json") {
		t.Fatalf("unexpected helper output: %q", out.String())
	}
}
