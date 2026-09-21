#!/bin/sh
set -eu

# Install a published release only; no Go, sudo, or Local runtime downloads.
case "$(uname -s)" in
    Darwin) os=darwin ;;
    Linux) os=linux ;;
    *) echo 'This installer supports macOS and Linux.' >&2; exit 1 ;;
esac
case "$(uname -m)" in
    x86_64|amd64) arch=amd64 ;;
    arm64|aarch64) arch=arm64 ;;
    *) echo 'Supported CPUs: x86-64 and ARM64.' >&2; exit 1 ;;
esac
for tool in curl tar mktemp; do
    command -v "$tool" >/dev/null 2>&1 || { echo "Required command missing: $tool" >&2; exit 1; }
done
if command -v sha256sum >/dev/null 2>&1; then
    hash_tool=sha256sum
elif command -v shasum >/dev/null 2>&1; then
    hash_tool=shasum
else
    echo 'Install sha256sum or shasum before continuing.' >&2; exit 1
fi
repo=https://github.com/blackandcode/localwp-cli
url=$(curl --fail --silent --show-error --location --proto '=https' --proto-redir '=https' -o /dev/null -w '%{url_effective}' "$repo/releases/latest")
tag=${url##*/}
case "$tag" in v[0-9]*) ;; *) echo 'Unexpected release tag.' >&2; exit 1 ;; esac
case "$tag" in *[!a-zA-Z0-9._-]*) echo 'Invalid release tag.' >&2; exit 1 ;; esac
name="localwp_${tag#v}_${os}_${arch}"
archive="$name.tar.gz"
temp_dir=$(mktemp -d)
trap 'rm -rf "$temp_dir"' EXIT HUP INT TERM
base="$repo/releases/download/$tag"
echo "Downloading localwp $tag from GitHub Releases..."
curl --fail --silent --show-error --location --proto '=https' --proto-redir '=https' "$base/$archive" -o "$temp_dir/$archive"
curl --fail --silent --show-error --location --proto '=https' --proto-redir '=https' "$base/SHA256SUMS.txt" -o "$temp_dir/SHA256SUMS.txt"
expected=$(awk -v file="$archive" '$2 == file {print $1}' "$temp_dir/SHA256SUMS.txt")
if [ "$hash_tool" = sha256sum ]; then
    actual=$(sha256sum "$temp_dir/$archive")
else
    actual=$(shasum -a 256 "$temp_dir/$archive")
fi
actual=${actual%% *}
if [ "${#expected}" -ne 64 ] || [ "$actual" != "$expected" ]; then
    echo 'Missing checksum or checksum mismatch; installation stopped.' >&2; exit 1
fi
tar -xzf "$temp_dir/$archive" -C "$temp_dir" "$name/localwp"
install_dir="$HOME/.local/bin"
mkdir -p "$install_dir"
cp "$temp_dir/$name/localwp" "$install_dir/localwp"
chmod +x "$install_dir/localwp"
path_line='export PATH="$HOME/.local/bin:$PATH"'
add_path() {
    if ! grep -Fqx "$path_line" "$1" 2>/dev/null; then
        printf '\n# localwp-cli\n%s\n' "$path_line" >> "$1"
    fi
    echo "PATH configured in $1"
}
case "${SHELL##*/}" in
    zsh) add_path "${ZDOTDIR:-$HOME}/.zshrc" ;;
    bash)
        add_path "$HOME/.bashrc"
        # Bash login shells read only the first existing login profile.
        if [ -f "$HOME/.bash_profile" ]; then add_path "$HOME/.bash_profile"
        elif [ -f "$HOME/.bash_login" ]; then add_path "$HOME/.bash_login"
        else add_path "$HOME/.profile"
        fi ;;
    *) echo "Add $install_dir to PATH using your shell's configuration. See installation/README.md." ;;
esac
echo "Installed $tag to $install_dir/localwp"
echo 'Open a new terminal, then run: localwp --version'
