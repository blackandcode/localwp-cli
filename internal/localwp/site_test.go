package localwp

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestLoadSitesAndSelectByDirectory(t *testing.T) {
	root := t.TempDir()
	siteRoot := filepath.Join(root, "Local Sites", "demo")
	if err := os.MkdirAll(filepath.Join(siteRoot, "app", "public", "wp-content", "plugins", "demo"), 0o755); err != nil {
		t.Fatal(err)
	}

	json := `{
      "site123": {
        "name": "Demo Site",
        "domain": "demo.local",
        "path": ` + quoteJSON(siteRoot) + `,
        "status": "running",
        "services": {
          "php": {"version": "8.3.14"},
          "mysql": {"name": "mysql", "version": "8.0.35", "ports": {"MYSQL": [10018]}}
        }
      }
    }`
	if err := os.WriteFile(filepath.Join(root, "sites.json"), []byte(json), 0o644); err != nil {
		t.Fatal(err)
	}

	sites, err := LoadSites(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(sites) != 1 {
		t.Fatalf("expected one site, got %d", len(sites))
	}
	site := sites[0]
	if site.ID != "site123" || site.Name != "Demo Site" || site.PHPVersion != "8.3.14" || site.DBPort != 10018 {
		t.Fatalf("unexpected site: %#v", site)
	}

	cwd := filepath.Join(siteRoot, "app", "public", "wp-content", "plugins", "demo")
	selected, err := SelectSite(sites, "", cwd)
	if err != nil {
		t.Fatal(err)
	}
	if selected.ID != "site123" {
		t.Fatalf("expected site123, got %s", selected.ID)
	}
}

func TestSelectSiteByNameDomainAndID(t *testing.T) {
	sites := []Site{
		{ID: "abc", Name: "Alpha", Domain: "alpha.local", Path: filepath.Join(t.TempDir(), "alpha")},
		{ID: "def", Name: "Beta", Domain: "beta.local", Path: filepath.Join(t.TempDir(), "beta")},
	}

	for _, selector := range []string{"abc", "Alpha", "alpha.local"} {
		site, err := SelectSite(sites, selector, "")
		if err != nil {
			t.Fatalf("selector %q: %v", selector, err)
		}
		if site.ID != "abc" {
			t.Fatalf("selector %q selected %s", selector, site.ID)
		}
	}
}

func quoteJSON(value string) string {
	b := []byte{'"'}
	for _, r := range value {
		switch r {
		case '\\':
			b = append(b, '\\', '\\')
		case '"':
			b = append(b, '\\', '"')
		default:
			b = append(b, string(r)...)
		}
	}
	b = append(b, '"')
	return string(b)
}

func TestDefaultDataDirMatchesHostPlatform(t *testing.T) {
	t.Setenv("LOCALWP_DATA_DIR", "")

	switch runtime.GOOS {
	case "windows":
		base := filepath.Join(t.TempDir(), "Roaming")
		t.Setenv("APPDATA", base)
		got, err := DefaultDataDir()
		if err != nil {
			t.Fatal(err)
		}
		want := filepath.Join(base, "Local")
		if got != want {
			t.Fatalf("want %q, got %q", want, got)
		}
	case "darwin":
		home := t.TempDir()
		t.Setenv("HOME", home)
		got, err := DefaultDataDir()
		if err != nil {
			t.Fatal(err)
		}
		want := filepath.Join(home, "Library", "Application Support", "Local")
		if got != want {
			t.Fatalf("want %q, got %q", want, got)
		}
	case "linux":
		xdg := filepath.Join(t.TempDir(), "config")
		t.Setenv("XDG_CONFIG_HOME", xdg)
		got, err := DefaultDataDir()
		if err != nil {
			t.Fatal(err)
		}
		want := filepath.Join(xdg, "Local")
		if got != want {
			t.Fatalf("want %q, got %q", want, got)
		}
	default:
		t.Skip("unsupported host for Local")
	}
}
