# Installation

`localwp-cli` requires [Local](https://localwp.com/) and supports Windows, macOS, and Linux hosts supported by Local.

For current Local operating-system requirements, use the official [Local installation guide](https://localwp.com/help-docs/getting-started/installing-local/).

## Install Local first

Download and install Local from:

https://localwp.com/

Create or import at least one WordPress site and make sure Local can start it normally.

## Install localwp

Follow the [guided release installation](../installation/README.md#guided-installation) for Windows or macOS/Linux. The helpers download a published binary, verify its checksum, and configure PATH for supported shells. No Go or Git is required.

The same guide documents [what the helpers change](../installation/README.md#what-the-helpers-change), [manual installation and PATH setup](../installation/README.md#manual-installation-advanced-users), updates, and uninstalling. Developers building from source should use [Contributing](../CONTRIBUTING.md#build-and-install-from-source).

## First real command

Start a site in Local and `cd` anywhere under that site's directory.

Then run:

```sh
localwp --doctor
localwp plugin list
```

`--doctor` should show the expected site path, Local data directory, PHP binary, `php.ini`, WP-CLI path and database tooling.

## Explicit site selection

If the current directory is not under a Local project or more than one site is running:

```sh
localwp --local-site "Site Name" plugin list
```

List available sites:

```sh
localwp --sites
```

## Local application data locations

The CLI detects Local using these application data roots:

```text
Windows   %APPDATA%\Local
macOS     ~/Library/Application Support/Local
Linux     ~/.config/Local
```

On Linux, `XDG_CONFIG_HOME` is respected when set.

## Troubleshooting

### `sites.json` cannot be found

Confirm Local is installed and has been opened at least once. Run Local, create/import a site, then retry.

For advanced debugging, override the Local data location:

```sh
LOCALWP_DATA_DIR="/custom/path/to/Local" localwp --sites
```

On Windows PowerShell:

```powershell
$env:LOCALWP_DATA_DIR = 'C:\custom\Local'
localwp --sites
```

### Local runtime cannot be resolved

Start the target site in Local and run:

```sh
localwp --doctor
```

The site-specific `php.ini` is generated under Local's runtime directory while the site is provisioned/running.

### `--shell` says the shell entry is missing

Normal one-shot commands do not use the interactive shell entry. `localwp --shell` does.

If the shell entry has not been generated, use **Open Site Shell** once in Local for that site and retry `localwp --shell`.

### Database commands fail

Verify the site is started in Local. Commands such as `plugin list`, `db export`, `search-replace`, and many option commands need access to the running database service.

Use:

```sh
localwp --doctor
```

to confirm the detected DB version, DB port and binary directory.

## Uninstall

See [uninstall instructions](../installation/README.md#troubleshooting-and-uninstalling). If you have a development checkout, the root `uninstall-localwp.cmd` (Windows) and `uninstall-localwp.sh` (macOS/Linux) scripts are also available.
