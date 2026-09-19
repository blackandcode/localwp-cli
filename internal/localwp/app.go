package localwp

import (
	"context"
	"fmt"
	"io"
	"runtime"
	"sort"
	"strings"
	"text/tabwriter"
)

func Run(args []string, version string, stdin io.Reader, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		printHelp(stdout, version)
		return 0
	}

	if len(args) == 1 {
		switch args[0] {
		case "--help", "-h":
			printHelp(stdout, version)
			return 0
		case "--version":
			fmt.Fprintln(stdout, version)
			return 0
		}
	}

	selector, remaining, err := parseWrapperArgs(args)
	if err != nil {
		fmt.Fprintf(stderr, "localwp: %v\n", err)
		return 2
	}

	dataDir, err := DefaultDataDir()
	if err != nil {
		fmt.Fprintf(stderr, "localwp: resolve Local data directory: %v\n", err)
		return 1
	}
	sites, err := LoadSites(dataDir)
	if err != nil {
		fmt.Fprintf(stderr, "localwp: %v\n", err)
		return 1
	}

	if len(remaining) == 1 && remaining[0] == "--sites" {
		printSites(stdout, sites)
		return 0
	}

	site, err := SelectSite(sites, selector, "")
	if err != nil {
		fmt.Fprintf(stderr, "localwp: %v\n", err)
		return 1
	}

	if len(remaining) == 1 && remaining[0] == "--shell" {
		return OpenShell(site, dataDir, stdin, stdout, stderr)
	}

	env := ResolveEnvironment(site, dataDir)
	if len(remaining) == 1 && remaining[0] == "--doctor" {
		printDoctor(stdout, site, env)
		return 0
	}

	if len(remaining) == 0 {
		printHelp(stdout, version)
		return 0
	}

	if err := env.Validate(site); err != nil {
		fmt.Fprintf(stderr, "localwp: %v\n", err)
		return 1
	}

	if !SiteRunning(site) {
		fmt.Fprintf(stderr, "localwp: warning: Local site %q does not appear to be running; commands that need the database may fail\n", site.Name)
	}

	spec := BuildWPCommand(site, env, remaining)
	return ExecuteWP(context.Background(), spec, stdin, stdout, stderr)
}

func parseWrapperArgs(args []string) (string, []string, error) {
	var selector string
	remaining := make([]string, 0, len(args))
	for i := 0; i < len(args); i++ {
		if args[i] == "--local-site" {
			if i+1 >= len(args) {
				return "", nil, fmt.Errorf("--local-site requires a site name, domain, or ID")
			}
			selector = args[i+1]
			i++
			continue
		}
		remaining = append(remaining, args[i])
	}
	return selector, remaining, nil
}

func printHelp(w io.Writer, version string) {
	fmt.Fprintf(w, `localwp %s

Run WP-CLI against Local WordPress sites from your normal terminal.

USAGE
  localwp <wp-cli command...>
  localwp --local-site "<site>" <wp-cli command...>

LOCALWP OPTIONS
  --local-site <name|domain|id>  Explicitly choose a Local site
  --sites                        List Local sites
  --doctor                       Show resolved Local runtime details
  --shell                        Open Local's interactive Site Shell
  --version                      Show localwp version
  --help                         Show this help

EXAMPLES
  localwp plugin list
  localwp core version
  localwp cache flush
  localwp option get siteurl
  localwp search-replace "http://old.local" "https://new.local" --dry-run
  localwp --local-site "My Site" plugin list

SITE SELECTION
  1. --local-site, when supplied
  2. the Local site containing the current directory
  3. the only Local site that appears to be running

SUPPORTED HOSTS
  Windows, macOS and Linux versions supported by Local.
`, version)
}

func printSites(w io.Writer, sites []Site) {
	sort.Slice(sites, func(i, j int) bool { return strings.ToLower(sites[i].Name) < strings.ToLower(sites[j].Name) })
	tw := tabwriter.NewWriter(w, 0, 4, 2, ' ', 0)
	fmt.Fprintln(tw, "RUNNING\tNAME\tDOMAIN\tPHP\tDATABASE\tID\tPATH")
	for _, site := range sites {
		running := "no"
		if SiteRunning(site) {
			running = "yes"
		}
		db := strings.TrimSpace(site.DBName + " " + site.DBVersion)
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n", running, site.Name, site.Domain, site.PHPVersion, db, site.ID, site.Path)
	}
	_ = tw.Flush()
}

func printDoctor(w io.Writer, site Site, env Environment) {
	running := "no / not detected"
	if SiteRunning(site) {
		running = "yes"
	}
	fmt.Fprintf(w, `Local site
  Name:       %s
  ID:         %s
  Domain:     %s
  Path:       %s
  WordPress:  %s
  Running:    %s
  PHP:        %s
  Database:   %s %s
  DB port:    %d

Resolved runtime
  OS/arch:        %s/%s
  Local data:     %s
  Resources:      %s
  PHP binary:     %s
  php.ini:        %s
  WP-CLI:         %s
  WP-CLI config:  %s
  DB binaries:    %s
  DB config:      %s
  Shell entry:    %s
`, site.Name, site.ID, site.Domain, site.Path, site.WebRoot, running, site.PHPVersion, site.DBName, site.DBVersion, site.DBPort,
		runtime.GOOS, runtime.GOARCH, env.DataDir, env.ResourcesRoot, env.PHPBinary, env.PHPIni, env.WPCLIPhar, env.WPCLIConfig, env.MySQLBinDir, env.MySQLConfDir, env.ShellEntry)
}
