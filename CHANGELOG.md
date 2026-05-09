# v0.5.4 (2026-05-09)

## New Features

- Add `device traffic stats` command — query traffic statistics per device with pagination, filtering by name/model/online status, and auto-merging device info with traffic data
- `device online-stats` redesigned — auto-fetches device list instead of requiring `--device-id`; supports `--name`, `--model`, `--online` filters and pagination; formats duration columns in table mode
- `device online-stats --end-time` now sets end of day to 23:59:59 (if today, uses current time)

## Improvements

- `device config export` now downloads the config file directly (previously only returned metadata). Supports `--file` for custom output path, defaults to server filename in current directory

---

# v0.5.3 (2026-05-09)

## New Features

### Config Get Refresh
- `device config get` now sends a "GET RUNNING CONFIG" task to the device before fetching config, ensuring the latest configuration is returned
- Use `--skip-refresh` to skip the task and return cached config directly

### Traffic Table Formatting
- Traffic columns (`send`, `receive`, `total`, `max`) now display human-readable sizes in table mode (e.g. `4.0 MiB` instead of `4194304`)
- JSON/YAML output retains raw byte values

### Empty Result Hint
- Display `No results.` when a query returns empty data, instead of showing nothing

## Improvements

- `device config set --content-file` now implemented (was previously TODO)
- `device traffic top --date` accepts `YYYY-MM` format (e.g. `2026-05`), defaults to current month if omitted
- All time range flags now document left-closed right-open semantics (`inclusive`/`exclusive`)
- `device traffic hourly` default range changed from 24 hours to 2 days (yesterday to tomorrow)
- `system log list --end-time` default changed to tomorrow (right-open)

---

# v0.5.2 (2026-05-08)

## Bug Fixes

### API Path Fixes
- **firmware list**: `/api/firmwares` → `/api/firmware`
- **firmware devices list**: `/api/jobs/{id}/devices` → `/api/job/{id}/devices`
- **firmware devices add**: `/api/firmwares/{id}/devices` → `/api/firmware/{id}/devices`
- **firmware devices remove**: `/api/jobs/{id}/devices/{did}` → `/api/job/{id}/devices/{did}`
- **firmware upgrade**: `/api/devices/{id}/upgrade` → `/api/device/{id}/upgrade`
- **edge control start/stop/restart**: GET → POST, path `.../app/{appId}/start` → `.../apps/start` + body `{apps:[appId]}`

### Parameter Fixes
- **task list**: Parameter names `status`→`states`, `type`→`types`, `device_name`→`object_id`; flag `--device-name` → `--object-id`
- **device count total**: Fix `--start-time/--end-time` format from unix timestamp to YYYY-MM-DD
- **device traffic top**: Fix `--month` → `--date`, format from YYYYMM to YYYY-MM-DD; `--limit` defaults to 10
- **device alert**: Fix `--start-time/--end-time` to convert YYYY-MM-DD to unix timestamp (supports both formats)
- **alert-rule update**: Fix GET response `result` unwrap; use raw JSON merge to preserve field types; send all required mutable fields on PUT

### Missing Field / Default Fixes
- **firmware upload**: Add missing `filename` form field
- **system log list**: Remove extra `language=1`; default to last 7 days
- **system user create**: Mark `--role-id` as required
- **device signal**: Make `--before` optional (defaults to now)
- **device traffic hourly**: Default to last 24 hours when `--after/--before` not specified
- **devicegroup list**: Remove invalid `--cursor/--limit`; add `--max-depth`
- **tunnel list**: Remove invalid `--cursor/--limit`

## Improvements

### Help Documentation
- Add Long descriptions and Examples to destructive commands (device delete/kick/reboot, permission add/remove)
- Add Long description and Example to `firmware upload` explaining the fid workflow
- Improve flag help text with format examples, enum values, and default value explanations across device, firmware, system, tunnel, edge, and alert-rule commands

- Fix gofmt formatting issues in kick.go and traffic_top.go

---

# v0.5.1 (2026-05-08)

## Bug Fixes

- **devicegroup list**: Remove invalid `--cursor`/`--limit` flags (API returns full tree, not paginated); add `--max-depth` flag
- **tunnel list**: Remove invalid `--cursor`/`--limit` flags (API does not support pagination)
- **system log list**: Default to last 7 days when `--start-time` is not specified (previously returned API error)

---

# v0.5.0 (2026-05-07)

## New Features

### Permission Group Member Management
- Restructure `permission users/devices/devicegroups` as subcommand groups with `list/add/remove`
- Add `system permission users add/remove` to manage users in a permission group
- Add `system permission devices add/remove` to manage devices in a permission group
- Add `system permission devicegroups list/add/remove` to manage device groups in a permission group

### Improvements
- `system org update --country` now accepts ISO 3166-1 alpha-2 code (e.g. CN, US) and sends `{code, name}` object matching platform frontend format
- Remove `--role` flag from user create/update (only `--role-id`)

## Bug Fixes

- **API command in Git Bash**: Fix MSYS path conversion issue where `/api/xxx` was converted to `E:/Git/api/xxx` on Windows Git Bash

---

# v0.4.0 (2026-05-07)

## New Features

### New Query Commands
- Add `device models` to list supported device models
- Add `device stats` for device overview (online/total counts)
- Add `device count online/total` for device count trends with date range
- Add `device traffic top` for monthly traffic ranking
- Add `device config export` to export device configuration
- Add `firmware get` for firmware details
- Add `firmware job-stats` for upgrade job statistics
- Add `system role list` to list roles in organization
- Add `system org list` to list organizations
- Add `system permission unassigned-users` to list users without a permission group

---

# v0.3.0 (2026-05-07)

## Breaking Changes

### Default Login Region
- Default `--host` changed from `cn` (iot.inhand.com.cn) to `global` (iot.inhandnetworks.com)

### Alert Rule Field Model
- `device alert-rule create` now uses `--alert-type`, `--for-device-type`, `--for-device-value`, `--notify-users`, `--notify-types`, `--notify-delay`, `--webhook-url`, `--webhook-secret`, `--locale` instead of the old `--metric`, `--condition`, `--threshold`, `--duration` flags
- `device alert-rule update` now supports `--for-device-type`, `--for-device-value`, `--notify-users`, `--notify-types`, `--notify-delay`, `--webhook-url`, `--webhook-secret`; removed `--threshold`

### User Create (Invitation Email Flow)
- `system user create` no longer requires `--password`; users receive an invitation email to set their password
- `--name` is now optional (not required for external users)
- Added `--role-id`, `--external`, `--lang` flags

## New Features

### Auth Impersonate
- `--org` is now required; `--user` is optional and must be combined with `--org`
- Automatically resolves org admin when only `--org` is given
- Pre-checks token expiry before calling impersonate API
- Improved error messages for API error responses

### Device Update
- Add `--mobile-number` flag

### DRC Create
- Add `--group-ids` flag for permission group association

### System User Update
- Add `--role-id` flag

### System Permission Update
- Add `--description` flag

### System Org Update
- Add `--email`, `--country`, `--biz-category` flags

## Bug Fixes

- **Edge version upload**: `--app` parameter was silently ignored; now correctly passed as `?app=` query parameter
- **Edge agent upload**: `--description` parameter was silently ignored; now correctly passed as `?description=` query parameter
- **Device config set**: Fix API path from `/config/set2` to `/config/set` to match platform frontend
- **Firmware create**: Fix API path from `/api/firmwares` to `/api/firmware` to match platform frontend

---

# v0.2.2 (2026-05-06)

## New Features

### Device Documentation
- Add `docs` command group to browse device model reference documentation from [model-reference](https://github.com/inhandnet/model-reference) repo
- `docs list` — list available models or model-specific document index
- `docs get <path>` — fetch a specific document, supports `--model` prefix
- `docs search <keyword>` — search across root and model indexes with structured output
- Override default repo via `DEVICEMANAGER_DOCS_REPO` environment variable

---

# v0.2.1 (2026-05-06)

## New Features

### Impersonate & Organization Switching
- Add `auth impersonate` command to impersonate another user (requires ROOT privilege)
- Supports `--user` (auto-resolves internal org), `--org` (auto-resolves org admin), or both
- `--stop` to restore admin identity; `auth status` shows impersonation state
- Add `auth orgs` to list organizations you belong to
- Add `auth switch-org <org-id>` to switch to a different organization

---

# v0.2.0 (2026-04-30)

## New Features

### Device Update & Delete
- Add `device update` command to edit device name and description
- Add `device delete` command to remove devices (with confirmation prompt)

### Alert Rules
- Add `device alert-rule` command group with full CRUD: `list`, `get`, `create`, `update`, `delete`
- Add `device alert-rule enable/disable` to toggle rules on/off
- Add `device alert-ack` to acknowledge (confirm) alerts

### Online Statistics
- Add `device online-stats` command to query device online rate, max online/offline duration, login count
- Accepts `--device-id` (repeatable) with `--start-time`/`--end-time` date range

### Task Management
- Add `task list` command for unified view of all tasks (DRC, firmware, etc.) with status/type/device filtering
- Add `task cancel` and `task restart` commands

### Impersonate & Organization Switching
- Add `auth impersonate` command to impersonate another user (requires ROOT privilege)
- Supports `--user` (auto-resolves internal org), `--org` (auto-resolves org admin), or both
- `--stop` to restore admin identity; `auth status` shows impersonation state
- Add `auth orgs` to list organizations you belong to
- Add `auth switch-org <org-id>` to switch to a different organization

### Device Diagnostics
- Add `device online-events` command to view device online/offline event timeline for troubleshooting disconnections
- Add `device register-events` command to view device registration history by serial number
- Add `edge app logs` command to view edge application runtime logs on a device

### System Management
- Add `system` command group with `user`, `permission`, `org`, and `log` subcommands
- `system user list/get/create/update/delete` — organization user management
- `system permission list/get/create/update/delete/users/devices` — device permission group management
- `system org get/update` — organization info
- `system log list` — audit log query with date range and level filtering

### Documentation
- Add `INSTALL.md` — AI-executable installation guide for automated CLI setup (platform detection, download with S3 China mirror fallback, checksum verification, install, login)
- Add `cmd/docgen` for auto-generating command reference markdown docs

---

# v0.1.2 (2026-04-29)

## New Features

### Confirmation Prompts
- Add interactive confirmation for all destructive operations (delete, remove, kick, reboot)
- Use `--yes/-y` to skip confirmation in scripts/CI
- Non-TTY environments (pipes) automatically skip confirmation

### Output
- TTY-adaptive default output format: `table` in terminal, `json` when piped
- Replace hardcoded ANSI colors with termenv-based Colorizer that auto-disables in non-TTY
- Add column formatters: `FormatBytes`, `FormatDuration`, `FormatRelativeTime`, `FormatPercent`, `TruncateRunes`
- Add `WithTransform` and `WithFormatters` FormatOption for table data transformation
- Add `normalizePage` to convert 0-based page numbers to 1-based

### API
- Add `ResultIDName` helper to extract `_id` and `name` from standard API responses

---

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
