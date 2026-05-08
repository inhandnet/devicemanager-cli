# Device Manager CLI Installation Guide

This document is designed for AI assistants (e.g. Claude Code) to read and execute.

## Steps

### 1. Detect Platform

Determine the OS and architecture of the current machine:

| OS      | Architecture | Binary Name                       |
|---------|-------------|-----------------------------------|
| Linux   | x86_64      | `devicemanager-linux-amd64`       |
| Linux   | arm64       | `devicemanager-linux-arm64`       |
| macOS   | x86_64      | `devicemanager-darwin-amd64`      |
| macOS   | arm64       | `devicemanager-darwin-arm64`      |
| Windows | x86_64      | `devicemanager-windows-amd64.exe` |

Run `uname -s` and `uname -m` to detect. On macOS, if `uname -m` returns `x86_64`, check `sysctl -n sysctl.proc_translated` — if it returns `1`, the shell is running under Rosetta 2 and the native architecture is `arm64`.

### 2. Download Binary and Checksums

Try **GitHub Releases** first. If it times out or fails, fall back to the **S3 mirror**.

**GitHub Releases (primary):**
```
https://github.com/inhandnet/devicemanager-cli/releases/latest/download/{BINARY_NAME}
https://github.com/inhandnet/devicemanager-cli/releases/latest/download/checksums.txt
```

**S3 mirror (fallback, China):**

First, fetch the latest version tag from the manifest:
```
https://elms-dm5-iot.s3.cn-north-1.amazonaws.com.cn/devicemanager-cli/manifest.yaml
```

The manifest is a YAML file. Read the `latest` field to get the version tag (e.g. `v0.1.2`), then download:
```
https://elms-dm5-iot.s3.cn-north-1.amazonaws.com.cn/devicemanager-cli/{TAG}/{BINARY_NAME}
https://elms-dm5-iot.s3.cn-north-1.amazonaws.com.cn/devicemanager-cli/{TAG}/checksums.txt
```

### 3. Verify Checksum

The `checksums.txt` file contains SHA256 checksums in the format:
```
<hash>  <filename>
```

Verify the downloaded binary:
- macOS: `shasum -a 256 <binary>`
- Linux: `sha256sum <binary>`
- Windows (PowerShell): `(Get-FileHash <binary> -Algorithm SHA256).Hash.ToLower()`

Compare the output hash with the corresponding entry in `checksums.txt`. **Do not proceed if the checksum does not match.**

### 4. Install

#### macOS / Linux

Make the binary executable and move it to the install path:

1. Try `/usr/local/bin/devicemanager` — if permission denied, use `sudo` (ask the user first)
2. If the user prefers no sudo, install to `~/.local/bin/devicemanager` instead (create the directory if needed, and remind the user to add `~/.local/bin` to their PATH if it's not already there)

```bash
chmod +x <binary>
mv <binary> /usr/local/bin/devicemanager
```

#### Windows

Rename the binary and move it to a directory in PATH:

```powershell
# Create install directory
New-Item -ItemType Directory -Force -Path "$env:LOCALAPPDATA\devicemanager"

# Move and rename
Move-Item <binary> "$env:LOCALAPPDATA\devicemanager\devicemanager.exe"

# Add to user PATH (persistent, takes effect in new terminal sessions)
$currentPath = [Environment]::GetEnvironmentVariable('Path', 'User')
if ($currentPath -notlike "*$env:LOCALAPPDATA\devicemanager*") {
    [Environment]::SetEnvironmentVariable('Path', "$currentPath;$env:LOCALAPPDATA\devicemanager", 'User')
}
```

After modifying PATH, refresh the current session or open a new terminal:

```powershell
# Refresh PATH in current session (other open terminals still need to be reopened)
$env:Path = [Environment]::GetEnvironmentVariable('Path', 'Machine') + ';' + [Environment]::GetEnvironmentVariable('Path', 'User')
```

### 5. Verify

Run `devicemanager version` to confirm the installation succeeded.

### 6. Login

```bash
devicemanager auth login
```

This opens a browser for OAuth authorization. It defaults to the Global region and creates a `default` context automatically.

Two production regions are available:

| Region | Short name | Domain                 | Command                                  |
|--------|-----------|------------------------|------------------------------------------|
| Global | `global`  | iot.inhandnetworks.com | `devicemanager auth login` (default)     |
| China  | `cn`      | iot.inhand.com.cn      | `devicemanager auth login --host cn`     |

Ask the user which region they need. After login, verify with `devicemanager auth status`.
