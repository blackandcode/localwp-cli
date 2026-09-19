package localwp

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveEnvironmentOverrides(t *testing.T) {
	root := t.TempDir()
	siteRoot := filepath.Join(root, "site")
	webRoot := filepath.Join(siteRoot, "app", "public")
	php := filepath.Join(root, "fake-php")
	phar := filepath.Join(root, "wp-cli.phar")
	mysqlDir := filepath.Join(root, "mysql-bin")
	phpIni := filepath.Join(root, "run", "site1", "conf", "php", "php.ini")

	for _, dir := range []string{webRoot, mysqlDir, filepath.Dir(phpIni)} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	for _, file := range []string{php, phar, phpIni} {
		if err := os.WriteFile(file, []byte("test"), 0o755); err != nil {
			t.Fatal(err)
		}
	}

	t.Setenv("LOCALWP_PHP_BINARY", php)
	t.Setenv("LOCALWP_WP_CLI_PHAR", phar)
	t.Setenv("LOCALWP_MYSQL_BIN_DIR", mysqlDir)

	site := Site{ID: "site1", Path: siteRoot, WebRoot: webRoot, PHPVersion: "8.3.14"}
	env := ResolveEnvironment(site, root)

	if env.PHPBinary != php || env.WPCLIPhar != phar || env.MySQLBinDir != mysqlDir {
		t.Fatalf("unexpected environment: %#v", env)
	}
	if env.PHPIni != phpIni {
		t.Fatalf("expected php.ini %q, got %q", phpIni, env.PHPIni)
	}
	if err := env.Validate(site); err != nil {
		t.Fatalf("expected valid environment: %v", err)
	}
}
