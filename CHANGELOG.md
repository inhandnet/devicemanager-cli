# v0.1.1 (2026-04-29)

## Improvements

### Self-Update
- Add S3 mirror (cn-north-1) as automatic fallback when GitHub is unreachable, improving update reliability for users in China

### Device Traffic
- Add `devicemanager device traffic hourly` command for hour-level traffic queries with `--after`/`--before` date range support (max 6 days)
- Simplify `traffic monthly` command: replace `--device` flag with positional `<device-id>` argument, consistent with `daily` and `hourly`

---

# v0.1.0 (2026-04-29)

Initial public release of the Device Manager CLI.

## Features

### Authentication
- Browser-based OAuth 2.0 login with automatic token exchange
- Multi-environment context management (China / Global / custom domains)
- Token auto-refresh on 401 responses
- Dynamic OAuth client discovery from platform

### Device Management
- List, get, create, kick, and reboot devices
- Filter by name, serial number, model, online status
- Signal quality history queries
- Monthly and daily traffic statistics
- Device client listing and batch queries
- Device alert listing with status filters
- Device configuration get/set

### Device Groups
- Full CRUD for device groups and subgroups
- Add/remove/list devices within groups
- Recursive device listing across subgroups
- Available devices query for group assignment

### Remote Tunnels
- Create, update, delete tunnels
- Connect and disconnect tunnel sessions

### DRC Configuration Templates
- Template CRUD with model filtering
- Assign/remove/restart devices on templates

### Edge Computing
- Edge engine upload, management, and device deployment
- Edge application CRUD
- Application version upload, deploy, update, delete
- Application configuration management and deployment
- Remote control (start/stop/restart) for edge apps

### Firmware Management
- Firmware upload and record creation
- Single device and batch upgrade
- Upgrade task device management by group

### API
- Generic `devicemanager api` command for arbitrary API calls
- Support for GET/POST/PUT/DELETE with query params, body fields, stdin input
- Custom headers and file download support

### Output
- JSON (colorized pretty / compact), table, YAML output formats
- jq expression filtering
- Verbose field control (1-100)
- Cursor-based pagination

### Self-Update
- `devicemanager update` command to check and install new versions from GitHub Releases
- SHA256 checksum verification
- Support for specific version targeting and JSON output

### Developer Experience
- golangci-lint v2 configuration with 12 linters + gofmt/goimports
- GitHub Actions CI (lint + test) and automated release workflow
- Cross-platform builds (linux/darwin/windows, amd64/arm64)
