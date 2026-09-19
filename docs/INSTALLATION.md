# Installation

`localwp-cli` requires [Local](https://localwp.com/) and supports Windows, macOS, and Linux hosts supported by Local.

For current Local operating-system requirements, use the official [Local installation guide](https://localwp.com/help-docs/getting-started/installing-local/).

## Install Local first

Download and install Local from:

https://localwp.com/

Create or import at least one WordPress site and make sure Local can start it normally.

## Install a prebuilt localwp binary

GitHub Releases are the preferred end-user installation method because they do not require Go.

Choose the archive matching your platform and CPU:

- `windows_amd64`
- `darwin_amd64` for Intel Macs
- `darwin_arm64` for Apple Silicon Macs
- `linux_amd64`
- `linux_arm64`

Extract the archive and put `localwp` or `localwp.exe` in a directory on `PATH`.

Each release contains a `SHA256SUMS.txt` file for integrity verification.

## Install from source: Windows

Requirements:

- Local
- Git
- Go 1.24+

Run:

```cmd
.\install-localwp.cmd
```

Use the explicit `.\` prefix. The unique launcher name prevents Windows from resolving an unrelated generic `install.cmd` from `PATH` (for example, an NVM for Windows installer).

The installer builds and installs:

```text
%LOCALAPPDATA%\localwp-cli\bin\localwp.exe
```

and adds the directory to your User PATH.

Open a **new** terminal after installation:

```cmd
where localwp
localwp --version
localwp --sites
```

## Install from source: macOS / Linux

Requirements:

- Local
- Git
- Go 1.24+

Run:

```sh
./install-localwp.sh
```

Default destination:

```text
~/.local/bin/localwp
```

You can choose another destination:

```sh
LOCALWP_INSTALL_DIR="$HOME/bin" ./install-localwp.sh
```

If the destination is not on `PATH`, add it to your shell profile. Example:

```sh
export PATH="$HOME/.local/bin:$PATH"
```

Then open a new terminal and verify:

```sh
command -v localwp
localwp --version
localwp --sites
```

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

Windows:

```cmd
.\uninstall-localwp.cmd
```

macOS / Linux:

```sh
./uninstall-localwp.sh
```
