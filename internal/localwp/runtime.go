package localwp

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
)

type Environment struct {
	DataDir       string
	ResourcesRoot string
	PHPBinary     string
	PHPIni        string
	WPCLIPhar     string
	WPCLIConfig   string
	MySQLBinDir   string
	MySQLConfDir  string
	ShellEntry    string
}

func ResolveEnvironment(site Site, dataDir string) Environment {
	roots := resourceRoots(dataDir)

	env := Environment{
		DataDir:      dataDir,
		PHPIni:       filepath.Join(dataDir, "run", site.ID, "conf", "php", "php.ini"),
		MySQLConfDir: filepath.Join(dataDir, "run", site.ID, "conf", "mysql"),
	}

	if runtime.GOOS == "windows" {
		env.ShellEntry = filepath.Join(dataDir, "ssh-entry", site.ID+".bat")
	} else {
		env.ShellEntry = filepath.Join(dataDir, "ssh-entry", site.ID+".sh")
	}

	if override := os.Getenv("LOCALWP_PHP_BINARY"); override != "" {
		env.PHPBinary = filepath.Clean(expandHome(override))
	} else {
		env.PHPBinary = findServiceBinary(dataDir, roots, "php", site.PHPVersion, phpRelativeCandidates())
	}

	if override := os.Getenv("LOCALWP_WP_CLI_PHAR"); override != "" {
		env.WPCLIPhar = filepath.Clean(expandHome(override))
	} else {
		for _, root := range roots {
			candidate := filepath.Join(root, "bin", "wp-cli", "wp-cli.phar")
			if fileExists(candidate) {
				env.ResourcesRoot = root
				env.WPCLIPhar = candidate
				break
			}
		}
	}

	if env.WPCLIPhar != "" {
		config := filepath.Join(filepath.Dir(env.WPCLIPhar), "config.yaml")
		if fileExists(config) {
			env.WPCLIConfig = config
		}
		if env.ResourcesRoot == "" {
			env.ResourcesRoot = filepath.Clean(filepath.Join(filepath.Dir(env.WPCLIPhar), "..", ".."))
		}
	}

	if override := os.Getenv("LOCALWP_MYSQL_BIN_DIR"); override != "" {
		env.MySQLBinDir = filepath.Clean(expandHome(override))
	} else if site.DBVersion != "" {
		serviceNames := []string{site.DBName, "mysql", "mariadb"}
		seen := map[string]bool{}
		for _, service := range serviceNames {
			service = strings.TrimSpace(service)
			if service == "" || seen[service] {
				continue
			}
			seen[service] = true
			dir := findServiceDir(dataDir, roots, service, site.DBVersion, dbRelativeCandidates())
			if dir != "" {
				env.MySQLBinDir = dir
				break
			}
		}
	}

	return env
}

func (env Environment) Validate(site Site) error {
	var missing []string
	if !dirExists(site.WebRoot) {
		missing = append(missing, "WordPress root: "+site.WebRoot)
	}
	if !fileExists(env.PHPBinary) {
		missing = append(missing, "Local PHP binary")
	}
	if !fileExists(env.PHPIni) {
		missing = append(missing, "site php.ini: "+env.PHPIni)
	}
	if !fileExists(env.WPCLIPhar) {
		missing = append(missing, "Local wp-cli.phar")
	}
	if len(missing) > 0 {
		return fmt.Errorf("could not resolve the Local runtime (%s); start the site in Local and run localwp --doctor", strings.Join(missing, "; "))
	}
	return nil
}

func resourceRoots(dataDir string) []string {
	var roots []string
	if override := os.Getenv("LOCALWP_RESOURCES_DIR"); override != "" {
		for _, item := range filepath.SplitList(override) {
			roots = appendUnique(roots, filepath.Clean(expandHome(item)))
		}
	}

	home, _ := os.UserHomeDir()
	switch runtime.GOOS {
	case "windows":
		if local := os.Getenv("LOCALAPPDATA"); local != "" {
			roots = appendUnique(roots, filepath.Join(local, "Programs", "Local", "resources", "extraResources"))
		}
		if pf := os.Getenv("ProgramFiles"); pf != "" {
			roots = appendUnique(roots, filepath.Join(pf, "Local", "resources", "extraResources"))
		}
		if pf86 := os.Getenv("ProgramFiles(x86)"); pf86 != "" {
			roots = appendUnique(roots, filepath.Join(pf86, "Local", "resources", "extraResources"))
		}
	case "darwin":
		roots = appendUnique(roots, "/Applications/Local.app/Contents/Resources/extraResources")
		roots = appendUnique(roots, filepath.Join(home, "Applications", "Local.app", "Contents", "Resources", "extraResources"))
	case "linux":
		roots = appendUnique(roots, "/opt/Local/resources/extraResources")
		roots = appendUnique(roots, "/usr/lib/Local/resources/extraResources")
		roots = appendUnique(roots, "/usr/share/Local/resources/extraResources")
	}

	return roots
}

func findServiceBinary(dataDir string, resourceRoots []string, service, version string, relatives []string) string {
	if service == "" || version == "" {
		return ""
	}

	serviceRoots := []string{filepath.Join(dataDir, "lightning-services")}
	for _, root := range resourceRoots {
		serviceRoots = appendUnique(serviceRoots, filepath.Join(root, "lightning-services"))
	}

	for _, root := range serviceRoots {
		dirs, _ := filepath.Glob(filepath.Join(root, service+"-"+version+"*"))
		sort.Sort(sort.Reverse(sort.StringSlice(dirs)))
		for _, dir := range dirs {
			for _, rel := range relatives {
				candidate := filepath.Join(dir, filepath.FromSlash(rel))
				if fileExists(candidate) {
					return candidate
				}
			}
		}
	}
	return ""
}

func findServiceDir(dataDir string, resourceRoots []string, service, version string, relatives []string) string {
	if service == "" || version == "" {
		return ""
	}

	serviceRoots := []string{filepath.Join(dataDir, "lightning-services")}
	for _, root := range resourceRoots {
		serviceRoots = appendUnique(serviceRoots, filepath.Join(root, "lightning-services"))
	}

	for _, root := range serviceRoots {
		dirs, _ := filepath.Glob(filepath.Join(root, service+"-"+version+"*"))
		sort.Sort(sort.Reverse(sort.StringSlice(dirs)))
		for _, dir := range dirs {
			for _, rel := range relatives {
				candidate := filepath.Join(dir, filepath.FromSlash(rel))
				if dirExists(candidate) {
					return candidate
				}
			}
		}
	}
	return ""
}

func phpRelativeCandidates() []string {
	switch runtime.GOOS {
	case "windows":
		return []string{"bin/win64/php.exe", "bin/win32/php.exe"}
	case "darwin":
		if runtime.GOARCH == "arm64" {
			return []string{"bin/darwin-arm64/bin/php", "bin/darwin/bin/php"}
		}
		return []string{"bin/darwin/bin/php", "bin/darwin-arm64/bin/php"}
	case "linux":
		return []string{"bin/linux/bin/php"}
	default:
		return nil
	}
}

func dbRelativeCandidates() []string {
	switch runtime.GOOS {
	case "windows":
		return []string{"bin/win64/bin", "bin/win64", "bin/win32/bin", "bin/win32"}
	case "darwin":
		if runtime.GOARCH == "arm64" {
			return []string{"bin/darwin-arm64/bin", "bin/darwin/bin"}
		}
		return []string{"bin/darwin/bin", "bin/darwin-arm64/bin"}
	case "linux":
		return []string{"bin/linux/bin"}
	default:
		return nil
	}
}

func ShellEntryForSite(site Site, dataDir string) (string, error) {
	ext := ".sh"
	if runtime.GOOS == "windows" {
		ext = ".bat"
	}
	entry := filepath.Join(dataDir, "ssh-entry", site.ID+ext)
	if !fileExists(entry) {
		return "", errors.New("Local Site Shell entry was not found; use Open Site Shell once for this site, then retry")
	}
	return entry, nil
}

func appendUnique(values []string, value string) []string {
	if value == "" {
		return values
	}
	for _, existing := range values {
		if existing == value {
			return values
		}
	}
	return append(values, value)
}

func fileExists(path string) bool {
	if path == "" {
		return false
	}
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func dirExists(path string) bool {
	if path == "" {
		return false
	}
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}
