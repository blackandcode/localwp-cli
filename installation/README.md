# Install localwp

This guide installs a ready-to-use binary from [GitHub Releases](https://github.com/blackandcode/localwp-cli/releases/latest). You do not need Go, Git, or administrator access.

First, install [Local](https://localwp.com/), create or import your WordPress site, and start it in Local.

## Guided installation

Use the instructions for your computer. Both helpers download the latest published release and check its SHA-256 checksum before installing it. Running the helper again updates the installed binary to the latest release.

### Windows

Supported release: Windows x64.

1. Download [install-release.ps1](https://raw.githubusercontent.com/blackandcode/localwp-cli/main/installation/install-release.ps1) into your **Downloads** folder. If the browser displays text, use **Save as** and keep the `.ps1` extension.
2. You can read the [script source](install-release.ps1) or open the downloaded file in Notepad before running it.
3. Open your Downloads folder in File Explorer. Right-click an empty area and select **Open in Terminal**. Use a **PowerShell** tab. Alternatively, open PowerShell and type `cd "$HOME\Downloads"` if that is where you saved the file.
4. Paste this command and press Enter:

   ```powershell
   powershell.exe -NoProfile -ExecutionPolicy Bypass -File .\install-release.ps1
   ```

The execution-policy option applies only to this process; it does not change your saved PowerShell policy. Organization-managed policies may still prevent scripts from running; use the manual steps below in that case.

The helper installs `localwp.exe` into `%LOCALAPPDATA%\localwp-cli\bin` and adds that folder to your **User PATH**. PATH tells Windows where to find commands, so you can type `localwp` in any project folder.

Close all windows of your terminal application, then reopen it. Restart your editor too if you use its integrated terminal. Run:

```powershell
localwp --version
localwp --sites
```

If Windows still cannot find the command, sign out and back in, or follow the PATH troubleshooting below.

### macOS and Linux

Supported releases: Intel/x86-64 and Apple Silicon/ARM64. The helper detects the platform and CPU. It needs `curl`, `tar`, and either `shasum` or `sha256sum`.

1. Download [install-release.sh](https://raw.githubusercontent.com/blackandcode/localwp-cli/main/installation/install-release.sh) into your **Downloads** folder. If the browser displays text, save it with the `.sh` extension.
2. You can read the [script source](install-release.sh) or open the downloaded file in a text editor before running it.
3. Open **Terminal**, then run these commands (adjust the folder if you saved elsewhere):

   ```sh
   cd "$HOME/Downloads"
   sh ./install-release.sh
   ```

The helper installs `localwp` into `~/.local/bin`. It adds that folder to PATH in `~/.zshrc` for zsh (respecting `ZDOTDIR`), or in `~/.bashrc` and your active login profile for Bash. Existing contents are preserved. If you use another shell, it prints the folder to add; see the manual instructions below.

Open a new terminal and run:

```sh
localwp --version
localwp --sites
```

## What the helpers change

- Download an archive and `SHA256SUMS.txt` from this repository's GitHub Releases over HTTPS. Windows queries GitHub's API to select the latest release; macOS/Linux follows the latest-release link.
- Compare the archive's checksum with the published checksum before extracting the binary. This detects a damaged or mismatched download; it is not an independent publisher signature.
- Install or replace only `localwp` in the user-level folder described above and configure PATH for the supported shells.
- Remove temporary downloads afterward. They do not install or change Local, PHP, WordPress, databases, or WP-CLI.

The download links above track this repository's `main` branch. Review the source before execution, or use the manual method to choose a specific release without running an installer script.

## Manual installation (advanced users)

### Download and verify

Open [Releases](https://github.com/blackandcode/localwp-cli/releases/latest) and download the archive for your computer, plus `SHA256SUMS.txt` from the same release. Do not choose GitHub's **Source code** downloads.

| Computer | Archive name ends with |
| --- | --- |
| Windows x64 | `_windows_amd64.zip` |
| Intel Mac | `_darwin_amd64.tar.gz` |
| Apple Silicon Mac | `_darwin_arm64.tar.gz` |
| Linux x86-64 | `_linux_amd64.tar.gz` |
| Linux ARM64 | `_linux_arm64.tar.gz` |

Compute the downloaded archive's checksum, replacing `ARCHIVE` with its actual filename:

```powershell
Get-FileHash -Algorithm SHA256 .\ARCHIVE.zip
```

On macOS use `shasum -a 256 ARCHIVE.tar.gz`; on Linux use `sha256sum ARCHIVE.tar.gz`. Compare the full hash with the line for that exact archive in `SHA256SUMS.txt`. Stop if it differs.

### Windows: install and set PATH

1. Right-click the downloaded ZIP and choose **Extract All**. Find `localwp.exe` inside the extracted folder.
2. Paste `%LOCALAPPDATA%` into File Explorer's address bar. Create a `localwp-cli` folder and a `bin` folder inside it. Copy `localwp.exe` into `bin`.
3. Search the Start menu for **Edit environment variables for your account** and open it.
4. Under **User variables**, select **Path**, then **Edit**, then **New**. Add `%LOCALAPPDATA%\localwp-cli\bin`. If Path does not exist, create it with that value. Preserve all existing entries. Add the folder, not the `.exe` filename.
5. Save with **OK** and close/reopen your terminal application and editor.
6. Run `where.exe localwp`, then `localwp --version`.

### macOS/Linux: install and set PATH

Extract the archive using your file manager or `tar -xzf ARCHIVE.tar.gz`. From the extracted directory containing `localwp`, run:

```sh
mkdir -p "$HOME/.local/bin"
cp ./localwp "$HOME/.local/bin/localwp"
chmod +x "$HOME/.local/bin/localwp"
```

For zsh, add the following line to `~/.zshrc` (or `$ZDOTDIR/.zshrc` if configured). For Bash, add it to `~/.bashrc` and the first existing login file among `~/.bash_profile`, `~/.bash_login`, and `~/.profile` (create `~/.profile` if none exists):

```sh
export PATH="$HOME/.local/bin:$PATH"
```

For fish, run `fish_add_path "$HOME/.local/bin"` in fish. For another shell, use its PATH configuration syntax.

Open a new terminal, then run `command -v localwp` and `localwp --version`.

## First WordPress command

Start your site in Local. From a directory inside that site's project, run:

```sh
localwp --doctor
localwp plugin list
```

Or select the intended site explicitly:

```sh
localwp --local-site "My Site" plugin list
```

## Troubleshooting and uninstalling

- **Command not found:** restart your terminal application/editor; check `where.exe localwp` on Windows or `command -v localwp` on macOS/Linux. The install folder must be on PATH. If multiple copies are found, check which one is first.
- **Download or checksum error:** the installer stops. Check your internet connection and try again, or download and verify the release manually.
- **Local site/runtime error:** see the [Local troubleshooting guide](../docs/INSTALLATION.md#troubleshooting).
- **Uninstall:** delete `localwp.exe` from `%LOCALAPPDATA%\localwp-cli\bin` on Windows, or `localwp` from `~/.local/bin` on macOS/Linux. Remove the Windows PATH entry if no longer needed. On POSIX, remove the helper's `# localwp-cli` line and following PATH line from the startup files only if you no longer want that PATH setup; other tools may also use `~/.local/bin`.

To build or install a development checkout, see [Contributing](../CONTRIBUTING.md#build-and-install-from-source).
