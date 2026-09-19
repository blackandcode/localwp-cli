package localwp

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Site struct {
	ID         string
	Name       string
	Domain     string
	Path       string
	WebRoot    string
	Status     string
	PHPVersion string
	DBName     string
	DBVersion  string
	DBPort     int
}

func DefaultDataDir() (string, error) {
	if override := strings.TrimSpace(os.Getenv("LOCALWP_DATA_DIR")); override != "" {
		return filepath.Clean(expandHome(override)), nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	switch runtime.GOOS {
	case "windows":
		if appData := os.Getenv("APPDATA"); appData != "" {
			return filepath.Join(appData, "Local"), nil
		}
		return filepath.Join(home, "AppData", "Roaming", "Local"), nil
	case "darwin":
		return filepath.Join(home, "Library", "Application Support", "Local"), nil
	case "linux":
		if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
			return filepath.Join(xdg, "Local"), nil
		}
		return filepath.Join(home, ".config", "Local"), nil
	default:
		return "", fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

func LoadSites(dataDir string) ([]Site, error) {
	content, err := os.ReadFile(filepath.Join(dataDir, "sites.json"))
	if err != nil {
		return nil, fmt.Errorf("read Local sites.json: %w", err)
	}

	var root any
	if err := json.Unmarshal(content, &root); err != nil {
		return nil, fmt.Errorf("parse Local sites.json: %w", err)
	}

	container := root
	if m, ok := root.(map[string]any); ok {
		if sites, exists := m["sites"]; exists {
			container = sites
		}
	}

	var sites []Site
	switch value := container.(type) {
	case map[string]any:
		keys := make([]string, 0, len(value))
		for key := range value {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			raw, ok := value[key].(map[string]any)
			if !ok {
				continue
			}
			site := siteFromMap(key, raw)
			if site.Path != "" && site.Name != "" {
				sites = append(sites, site)
			}
		}
	case []any:
		for i, item := range value {
			raw, ok := item.(map[string]any)
			if !ok {
				continue
			}
			id := stringValue(raw["id"])
			if id == "" {
				id = strconv.Itoa(i)
			}
			site := siteFromMap(id, raw)
			if site.Path != "" && site.Name != "" {
				sites = append(sites, site)
			}
		}
	default:
		return nil, errors.New("Local sites.json has an unsupported structure")
	}

	if len(sites) == 0 {
		return nil, errors.New("no Local sites found")
	}
	return sites, nil
}

func siteFromMap(fallbackID string, raw map[string]any) Site {
	id := stringValue(raw["id"])
	if id == "" {
		id = fallbackID
	}
	path := expandHome(stringValue(raw["path"]))
	cleanPath := ""
	webRoot := ""
	if path != "" {
		cleanPath = filepath.Clean(path)
		webRoot = filepath.Join(cleanPath, "app", "public")
	}

	site := Site{
		ID:      id,
		Name:    stringValue(raw["name"]),
		Domain:  stringValue(raw["domain"]),
		Path:    cleanPath,
		Status:  stringValue(raw["status"]),
		WebRoot: webRoot,
	}

	services, _ := raw["services"].(map[string]any)
	if php, ok := services["php"].(map[string]any); ok {
		site.PHPVersion = stringValue(php["version"])
	}
	if site.PHPVersion == "" {
		site.PHPVersion = stringValue(raw["phpVersion"])
	}

	for _, dbKey := range []string{"mysql", "mariadb"} {
		if db, ok := services[dbKey].(map[string]any); ok {
			site.DBName = dbKey
			if name := stringValue(db["name"]); name != "" {
				site.DBName = name
			}
			site.DBVersion = stringValue(db["version"])
			if ports, ok := db["ports"].(map[string]any); ok {
				for _, portKey := range []string{"MYSQL", "mysql"} {
					if port := firstPort(ports[portKey]); port > 0 {
						site.DBPort = port
						break
					}
				}
			}
			break
		}
	}

	if site.DBVersion == "" {
		database := stringValue(raw["database"])
		if parts := strings.SplitN(database, "-", 2); len(parts) == 2 {
			site.DBName, site.DBVersion = parts[0], parts[1]
		}
	}

	return site
}

func SelectSite(sites []Site, selector, cwd string) (Site, error) {
	if selector != "" {
		return selectByName(sites, selector)
	}

	if cwd == "" {
		var err error
		cwd, err = os.Getwd()
		if err != nil {
			return Site{}, err
		}
	}
	cwd, _ = filepath.Abs(cwd)

	type match struct {
		site Site
		n    int
	}
	var matches []match
	for _, site := range sites {
		sitePath, err := filepath.Abs(site.Path)
		if err != nil {
			continue
		}
		if pathContains(sitePath, cwd) {
			matches = append(matches, match{site: site, n: len(sitePath)})
		}
	}
	if len(matches) > 0 {
		sort.Slice(matches, func(i, j int) bool { return matches[i].n > matches[j].n })
		return matches[0].site, nil
	}

	var running []Site
	for _, site := range sites {
		if SiteRunning(site) {
			running = append(running, site)
		}
	}
	if len(running) == 1 {
		return running[0], nil
	}
	if len(running) > 1 {
		names := make([]string, 0, len(running))
		for _, site := range running {
			names = append(names, site.Name)
		}
		return Site{}, fmt.Errorf("current directory is not inside a Local site and multiple sites appear to be running: %s; use --local-site", strings.Join(names, ", "))
	}

	return Site{}, errors.New("could not determine the Local site; run from inside a Local site directory, start exactly one site, or use --local-site")
}

func selectByName(sites []Site, selector string) (Site, error) {
	needle := strings.ToLower(strings.TrimSpace(selector))
	var exact []Site
	for _, site := range sites {
		values := []string{site.ID, site.Name, site.Domain, filepath.Base(site.Path)}
		for _, value := range values {
			if strings.EqualFold(value, needle) {
				exact = append(exact, site)
				break
			}
		}
	}
	if len(exact) == 1 {
		return exact[0], nil
	}
	if len(exact) > 1 {
		return Site{}, fmt.Errorf("multiple Local sites exactly match %q; use the site ID", selector)
	}

	var partial []Site
	for _, site := range sites {
		if strings.Contains(strings.ToLower(site.Name), needle) || strings.Contains(strings.ToLower(site.Domain), needle) {
			partial = append(partial, site)
		}
	}
	if len(partial) == 1 {
		return partial[0], nil
	}
	if len(partial) > 1 {
		return Site{}, fmt.Errorf("Local site selector %q is ambiguous", selector)
	}
	return Site{}, fmt.Errorf("Local site %q was not found", selector)
}

func SiteRunning(site Site) bool {
	if os.Getenv("LOCALWP_SKIP_RUNNING_CHECK") == "1" {
		return strings.EqualFold(site.Status, "running") || site.DBPort > 0
	}
	if site.DBPort > 0 {
		conn, err := net.DialTimeout("tcp", net.JoinHostPort("127.0.0.1", strconv.Itoa(site.DBPort)), 200*time.Millisecond)
		if err != nil {
			return false
		}
		_ = conn.Close()
		return true
	}
	return strings.EqualFold(site.Status, "running")
}

func pathContains(parent, child string) bool {
	rel, err := filepath.Rel(parent, child)
	if err != nil {
		return false
	}
	return rel == "." || (rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)))
}

func stringValue(v any) string {
	switch value := v.(type) {
	case string:
		return value
	case json.Number:
		return value.String()
	case float64:
		return strconv.FormatFloat(value, 'f', -1, 64)
	default:
		return ""
	}
}

func firstPort(v any) int {
	switch value := v.(type) {
	case []any:
		for _, item := range value {
			if port := firstPort(item); port > 0 {
				return port
			}
		}
	case float64:
		return int(value)
	case string:
		p, _ := strconv.Atoi(value)
		return p
	}
	return 0
}

func expandHome(path string) string {
	if path == "" || path[0] != '~' {
		return path
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return path
	}
	if path == "~" {
		return home
	}
	if len(path) > 1 && (path[1] == '/' || path[1] == '\\') {
		return filepath.Join(home, path[2:])
	}
	return path
}
