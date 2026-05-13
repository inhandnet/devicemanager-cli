# devicemanager CLI

Command-line tool for the InHand Device Manager (DM) platform. Supports authentication, multi-environment context switching, device management, and multiple output formats.

## Installation

### Build from source

```bash
# Requires Go 1.25+
make build    # Output to bin/devicemanager
make install  # Install to $GOPATH/bin
```

> On macOS, `CGO_ENABLED=0` is required (already set in Makefile) to avoid dyld LC_UUID errors.

### Cross-platform build

CI automatically builds binaries for the following platforms:

- `linux/amd64`, `linux/arm64`
- `darwin/amd64`, `darwin/arm64`
- `windows/amd64`

## Quick start

### 1. Login

```bash
devicemanager auth login                          # Login to global region (iot.inhandnetworks.com)
devicemanager auth login --host cn                # Login to China region (iot.inhand.com.cn)
devicemanager auth login --host iot.example.com   # Custom domain
devicemanager auth login --context prod           # Create/update a named context
```

Login uses the OAuth 2.0 Authorization Code flow — it opens a browser for authorization. The CLI reuses the platform's SPA OAuth client. A local callback server (default `http://localhost:18920/callback`) receives the authorization code and exchanges it for a token automatically.

### 2. Verify

```bash
devicemanager auth status
devicemanager device list
```

## Command reference

### Authentication

```bash
devicemanager auth login                    # Browser-based OAuth login
devicemanager auth status                   # View current auth status
devicemanager auth logout                   # Log out
devicemanager auth impersonate --org <oid>             # Impersonate org admin
devicemanager auth impersonate --org <oid> --user <uid> # Impersonate specific user
devicemanager auth impersonate --stop                  # Stop impersonation
devicemanager auth orgs                    # List your organizations
devicemanager auth switch-org <org-id>     # Switch to another organization
```

### Context management

Contexts are created/updated at login via `--context`. Other subcommands are for switching, viewing, and deleting:

```bash
devicemanager config use-context <name>
devicemanager config current-context
devicemanager config list-contexts
devicemanager config delete-context <name>
```

### API calls

```bash
devicemanager api /api/users/this                                 # GET request
devicemanager api /api/devices -q page=0 -q limit=10              # With query params
devicemanager api /api/devices -X POST -f name=test               # POST with body fields
echo '{}' | devicemanager api /api/devices -X POST --input -      # Read JSON body from stdin
devicemanager api /api/users/this -H "Sudo: user@example.com"     # Custom header
```

### Device management

```bash
devicemanager device list                                          # List devices (default limit 20)
devicemanager device list --online 1 --model IR615                 # Filter by status/model
devicemanager device list --name router-01 --serial-number GL5022  # Filter by name/SN
devicemanager device list --cursor 20 --limit 50                   # Pagination: skip 20, take 50
devicemanager device list -o json                                  # JSON output (verbose=100 by default)

devicemanager device get <device-id>                               # Device details
devicemanager device create --name <name> --serial-number <sn>     # Add a device
devicemanager device import devices.xlsx                           # Batch import devices from Excel
devicemanager device import devices.xlsx --group <group-id>        # Import and assign to group
devicemanager device import devices.xlsx --no-overwrite            # Import without overwriting existing
devicemanager device models                                        # List supported device models
devicemanager device stats                                         # Device overview (online/total counts)
devicemanager device signal <device-id> --after <ISO>                 # Signal quality (from start to now)
devicemanager device signal <device-id> --after <ISO> --before <ISO>  # Signal quality (time range)
devicemanager device kick <device-id>                              # Force disconnect
devicemanager device reboot <device-id> --timeout 15000            # Reboot (milliseconds)

# Device traffic
devicemanager device traffic monthly 202604 <device-id>            # Monthly traffic
devicemanager device traffic daily 202604 <device-id>              # Daily traffic
devicemanager device traffic hourly <device-id>                    # Hourly traffic (default: last 2 days)
devicemanager device traffic hourly <device-id> --after 2026-04-25 --before 2026-04-27  # Custom range (max 6 days)
devicemanager device traffic top --date 2026-04 --limit 10         # Top 10 devices by traffic
devicemanager device traffic stats --after 2026-05-01 --before 2026-05-10  # Traffic stats per device

# Device count trends
devicemanager device count online --start-time 1714492800 --end-time 1717084800  # Online count (unix timestamp)
devicemanager device count total --start-time 2026-04-01 --end-time 2026-05-01   # Total count (YYYY-MM-DD)

# Device clients
devicemanager device clients list <device-id>                      # List connected clients
devicemanager device clients batch <device-id>...                  # Batch query clients

# Remote web management
devicemanager device web <device-id>                               # Start remote web management
devicemanager device web <device-id> --port 443 --proto https      # Custom port/protocol

# Device update & delete
devicemanager device update <device-id> --name "new-name"          # Rename device
devicemanager device update <device-id> --description "office"     # Update description
devicemanager device update <device-id> --mobile-number "1234567"  # Update mobile number
devicemanager device delete <device-id>                            # Delete a device

# Device alerts
devicemanager device alert                                         # List alerts
devicemanager device alert --device-name router --state unconfirmed # Filter by condition
devicemanager device alert-ack <alert-id>                          # Acknowledge an alert

# Alert rules
devicemanager device alert-rule list                               # List alert rules
devicemanager device alert-rule list --device-name router          # Filter by device
devicemanager device alert-rule get <rule-id>                      # Rule details
devicemanager device alert-rule create \
  --name "offline-alert" \
  --alert-type offline                                             # Create an alert rule
devicemanager device alert-rule create \
  --name "traffic-alert" \
  --alert-type daily_traffic_excess \
  --for-device-type DEVICE \
  --for-device-value id1,id2 \
  --notify-users uid1,uid2 \
  --notify-types email,webhook \
  --webhook-url https://example.com/hook                           # Create with notifications
devicemanager device alert-rule update <rule-id> --name "new-name" # Update rule
devicemanager device alert-rule enable <rule-id>                   # Enable rule
devicemanager device alert-rule disable <rule-id>                  # Disable rule
devicemanager device alert-rule delete <rule-id>                   # Delete rule

# Online statistics (per device, with pagination)
devicemanager device online-stats \
  --start-time 2026-05-01 \
  --end-time 2026-05-09                                            # All devices
devicemanager device online-stats \
  --start-time 2026-05-01 \
  --end-time 2026-05-09 --name router --online 1                   # Filter by name + online

# Device event logs (for troubleshooting)
devicemanager device online-events <device-id> \
  --start-time 2026-04-29 --end-time 2026-04-30                   # Online/offline event timeline
devicemanager device register-events <serial-number>               # Registration event history

# Device configuration
devicemanager device config get <device-id>                        # Get running config
devicemanager device config set <device-id> --content "..."        # Push configuration
devicemanager device config export <device-id>                     # Export to current directory
devicemanager device config export <device-id> --file ./config.dat # Export to specific path
```

### Device groups

```bash
devicemanager devicegroup list                                     # List groups
devicemanager devicegroup list --parent <parent-id>                # Filter by parent group
devicemanager devicegroup get <group-id>                           # Group details
devicemanager devicegroup create --name "Factory A"                # Create a group
devicemanager devicegroup create --name "Line 1" --parent <id>     # Create a subgroup
devicemanager devicegroup update <group-id> --name "New Name"      # Rename a group
devicemanager devicegroup delete <group-id>                        # Delete a group

# Devices within a group
devicemanager devicegroup devices <group-id> list                  # List devices in group
devicemanager devicegroup devices <group-id> list --recursive      # Include subgroup devices
devicemanager devicegroup devices <group-id> add <device-id>...    # Add devices to group
devicemanager devicegroup devices <group-id> remove <device-id>... # Remove devices from group
devicemanager devicegroup devices <group-id> available             # Devices available to add
```

### Remote tunnels

```bash
devicemanager tunnel list                                          # List tunnels
devicemanager tunnel list --device-id <id>                         # Filter by device
devicemanager tunnel create \
  --name ssh-tunnel \
  --device-id <id> \
  --proto tcp \
  --local-address 127.0.0.1 \
  --local-port 22               # Create a tunnel
devicemanager tunnel update <tunnel-id> --name "new-name"          # Update tunnel
devicemanager tunnel delete <tunnel-id>                            # Delete tunnel
devicemanager tunnel connect <tunnel-id>                           # Connect tunnel
devicemanager tunnel disconnect <tunnel-id>                        # Disconnect tunnel
```

### DRC configuration templates

```bash
devicemanager drc list                                             # List templates
devicemanager drc list --model IR615                               # Filter by device model
devicemanager drc get <template-id>                                # Template details
devicemanager drc create \
  --name "IR615-default" \
  --model IR615 \
  --content "..."               # Create a template
devicemanager drc delete <template-id>                             # Delete a template

# Template device management
devicemanager drc devices <template-id> list                       # List assigned devices
devicemanager drc devices <template-id> list --status running      # Filter by status
devicemanager drc devices <template-id> add <device-id>...         # Assign devices
devicemanager drc devices <template-id> add --group <group-id>     # Assign by device group
devicemanager drc devices <template-id> remove <device-id>         # Remove a device
devicemanager drc devices <template-id> restart <device-id>        # Restart device task
```

### Edge computing

#### Edge engines

```bash
devicemanager edge agent list                                # List engines
devicemanager edge agent list --version v1.0                 # Filter by version
devicemanager edge agent get <agent-id>                       # Engine details
devicemanager edge agent upload <file-path> --description "IR615 engine"  # Upload engine
devicemanager edge agent update <agent-id> --description "new desc"       # Update engine
devicemanager edge agent delete <agent-id>                    # Delete engine
devicemanager edge agent devices <agent-id>                   # List deployed devices (default: pending)
devicemanager edge agent devices <agent-id> --status ready    # Filter by status (pending/installing/downloading/ready/failed)
devicemanager edge agent deploy <agent-id> <device-id>...     # Deploy agent to devices
devicemanager edge agent deploy <agent-id> --group <group-id> # Deploy agent to device group
devicemanager edge agent undeploy <agent-id> <device-id>...   # Remove agent from devices
```

#### Edge applications

```bash
devicemanager edge app list                                   # List applications
devicemanager edge app get <app-id>                           # Application details
devicemanager edge app create --name "my-app" --description "..."          # Create application
devicemanager edge app update <app-id> --description "new desc"           # Update application
devicemanager edge app delete <app-id>                        # Delete application
devicemanager edge app logs <device-id> <app-name>            # View app runtime logs on device
devicemanager edge app deploy <app-id> --version <ver> <device-id>...  # Deploy app to devices
devicemanager edge app deploy <app-id> --version <ver> --group <id>    # Deploy app to group
devicemanager edge app undeploy <app-id> <device-id>...       # Cancel app deployment from devices
```

#### Application versions

```bash
devicemanager edge version list <app-id>                      # List versions
devicemanager edge version upload <file-path> --app <app-id>  # Upload version
devicemanager edge version update <app-id> <version> --notes "Release notes"  # Update notes
devicemanager edge version delete <app-id> <version>          # Delete version
devicemanager edge version deploy <app-id> <version> --device <id> --group <id>  # Deploy version
```

#### Application configuration

```bash
devicemanager edge config list <app-id>                       # List configs
devicemanager edge config list <app-id> --version v1.0        # Filter by version
devicemanager edge config get <app-id> <config-id>            # Config details
devicemanager edge config create <app-id> --version v1.0 --content "..."    # Create config
devicemanager edge config update <app-id> <config-id> --description "..."   # Update config
devicemanager edge config delete <app-id> <config-id>         # Delete config
devicemanager edge config deploy <app-id> <version> --device <id> --group <id>  # Deploy config
devicemanager edge config undeploy <app-id> <device-id>...    # Cancel config deployment from devices
```

#### Remote control

```bash
devicemanager edge control start <device-id> <app-id>         # Start application
devicemanager edge control stop <device-id> <app-id>          # Stop application
devicemanager edge control restart <device-id> <app-id>       # Restart application
devicemanager edge control remove <device-id> <app-id>        # Uninstall application
```

#### Device edge status

```bash
devicemanager edge device <device-id>                         # Get edge status (agent & apps) of a device
```

### Task management

```bash
devicemanager task list                                            # List all tasks
devicemanager task list --status running                           # Filter by status (running/waiting/failed/completed)
devicemanager task list --type config-apply                        # Filter by type (config-apply/interactive-command/fetch-config/import-firmware/vpn-channel/vpn-link-order/token-cleanup/traffic-stats/idle-notice/remote-web)
devicemanager task list --object-id <device-id>                    # Filter by device ID
devicemanager task cancel <task-id>                                # Cancel a task
devicemanager task restart <task-id>                               # Restart a task
```

### System management

#### Users

```bash
devicemanager system user list                                     # List users in organization
devicemanager system user list --oid <org-id>                      # List users in a specific organization
devicemanager system user get <user-id>                            # User details
devicemanager system user create \
  --name "test" \
  --email "test@example.com"               # Create a user (invitation email sent)
devicemanager system user create \
  --email "ext@example.com" --external     # Create an external user
devicemanager system user update <user-id> --name "new-name"       # Update user
devicemanager system user update <user-id> --role-id <role-id>     # Change role
devicemanager system user delete <user-id>                         # Delete user
```

#### Roles

```bash
devicemanager system role list                                     # List roles in organization
```

#### Device permissions

```bash
devicemanager system permission list                               # List permission groups
devicemanager system permission get <group-id>                     # Group details
devicemanager system permission create --name "office-devices"     # Create permission group
devicemanager system permission update <group-id> --name "new"     # Update group name
devicemanager system permission update <group-id> --description "desc"  # Update description
devicemanager system permission delete <group-id>                  # Delete group
devicemanager system permission users list <group-id>              # List users in group
devicemanager system permission users add <group-id> <uid>...      # Add users to group
devicemanager system permission users remove <group-id> <uid>...   # Remove users from group
devicemanager system permission devices list <group-id>            # List devices in group
devicemanager system permission devices add <group-id> <did>...    # Add devices to group
devicemanager system permission devices remove <group-id> <did>... # Remove devices from group
devicemanager system permission devicegroups list <group-id>       # List device groups in group
devicemanager system permission devicegroups add <group-id> <dgid>...    # Add device groups
devicemanager system permission devicegroups remove <group-id> <dgid>... # Remove device groups
devicemanager system permission unassigned-users                   # List users without a permission group
```

#### Organization

```bash
devicemanager system org list                                      # List organizations
devicemanager system org list --name "InHand"                      # Filter by name
devicemanager system org list --email "info@example.com"           # Filter by email
devicemanager system org get                                       # View current org info
devicemanager system org update <org-id> --name "New Org Name"     # Update org name
devicemanager system org update <org-id> --email "org@example.com" # Update org email
devicemanager system org update <org-id> --country US              # Update country (ISO 3166-1 alpha-2)
```

#### Audit logs

```bash
devicemanager system log list                                      # List audit logs (defaults to last 7 days)
devicemanager system log list --start-time 2026-04-24 --end-time 2026-04-30  # Filter by date
devicemanager system log list --level warning                      # Filter by level
```

### Firmware management

```bash
devicemanager firmware list                                      # List firmware
devicemanager firmware list --model IR615                        # Filter by model
devicemanager firmware get <firmware-id>                         # Firmware details
devicemanager firmware upload <file-path>                        # Upload firmware file
devicemanager firmware create \
  --fid <file-id> \
  --name "IR615-v2.0" \
  --version 2.0.0 \
  --model IR615             # Create firmware record
devicemanager firmware delete <firmware-id>                      # Delete a firmware
devicemanager firmware upgrade <device-id> --firmware-id <id>    # Upgrade a single device (timeout default 600s)

# Batch upgrade management
devicemanager firmware devices <firmware-id> list                # List devices in upgrade task
devicemanager firmware devices <firmware-id> add <device-id>...  # Add devices for batch upgrade
devicemanager firmware devices <firmware-id> add --group <group-id>  # Upgrade by device group
devicemanager firmware devices <firmware-id> remove <device-id>  # Remove device from upgrade
devicemanager firmware job-stats <firmware-id>                   # Upgrade job statistics
devicemanager firmware cancel <firmware-id> <device-id>          # Cancel upgrade for a device
devicemanager firmware retry <firmware-id> <device-id>           # Retry failed upgrade for a device
```

### Device model documentation

```bash
devicemanager docs list                                         # List all available device models
devicemanager docs list --model ER805                           # List documents for a specific model
devicemanager docs get <path>                                   # Get a document by path
devicemanager docs get --model ER805 <path>                     # Get a document under a specific model
devicemanager docs search <keyword>                             # Search documents across all models
devicemanager docs search --model ER805 <keyword>               # Search within a specific model
```

Documents are fetched from the [inhandnet/model-reference](https://github.com/inhandnet/model-reference) GitHub repository. Override the source with `DEVICEMANAGER_DOCS_REPO`.

### Debugging

```bash
devicemanager device list --debug                              # Print config/auth/HTTP debug info to stderr
DEVICEMANAGER_DEBUG=1 devicemanager device list                     # Enable via environment variable
devicemanager device list --debug -o json 2>/tmp/debug.log     # Write debug to file, keep stdout clean
```

### Global flags

```bash
devicemanager --context prod auth status            # Temporarily switch context
devicemanager --debug device list                   # Enable debug output
devicemanager --verbose 50 device list              # Reduce response detail level
devicemanager --jq '.[].name' device list           # Filter JSON with jq expression
devicemanager version                                # Show version info
```

## Output formats

Use `-o` to specify the output format. Default is `table` in a terminal and `json` when piped.

| Format | TTY behavior | Pipe behavior |
|--------|-------------|---------------|
| `table` | Aligned table | TSV |
| `json` | Colorized pretty JSON | Compact JSON |
| `yaml` | YAML | YAML |

```bash
devicemanager device list -o table                   # Table output
devicemanager device list -o yaml                    # YAML output
devicemanager device list --jq '.[] | .name'         # Filter with jq expression
```

The `{"result": ...}` envelope from the server is automatically unwrapped in yaml/json/jq modes — only the contents of `result` are shown.

### Field verbosity (`--verbose`)

The DM API uses a `verbose` parameter to control how many fields are returned (1-100, higher = more detailed). The CLI provides a global `--verbose` flag (default 100) that is automatically applied to all GET (query) requests. POST/PUT/DELETE requests are not affected.

```bash
devicemanager device list                            # verbose=100 by default (all fields)
devicemanager --verbose 15 device list               # Minimal fields
```

### Pagination

```bash
devicemanager device list --cursor 0 --limit 20       # First page (default)
devicemanager device list --cursor 20 --limit 50      # Skip 20, take 50
```

The DM platform uses `cursor` (skip offset) + `limit` pagination. The `list` commands also accept `--page-size`/`--per-page` as hidden aliases for `--limit`.

## Environment variables

| Variable | Purpose |
|----------|---------|
| `DEVICEMANAGER_CONTEXT` | Override the current context |
| `DEVICEMANAGER_HOST` | Override the host in the current context |
| `DEVICEMANAGER_TOKEN` | Override the token in the current context |
| `DEVICEMANAGER_DEBUG` | Set to any non-empty value to enable debug output |

## Configuration file

Path: `~/.config/devicemanager/config.yaml` (permissions `0600`)

The configuration file stores all context information (host, token, etc.) and is managed via the `devicemanager config` subcommands.

## Development

### Prerequisites

- Go 1.25+
- [golangci-lint](https://golangci-lint.run/)

### Build & test

```bash
make build       # Build to bin/devicemanager
make build-all   # Cross-platform build
make install     # Install to GOPATH
make test        # Run tests
make fmt         # Format code (gofmt + goimports)
make lint        # Run golangci-lint
make clean       # Clean build artifacts
```

### Project structure

```
cmd/devicemanager/       # CLI entry point
internal/
  api/              # OAuth, token transport & auto-refresh, REST client, callback server
  build/            # Injected Version/Commit/Date
  cmd/              # Subcommand implementations
    auth/           # Login, logout, auth status, impersonate, switch-org, orgs
    config/         # Context management
    device/         # Device management
    devicegroup/    # Device group management
    tunnel/         # Remote tunnel management
    drc/            # DRC configuration template management
    edge/           # Edge computing (engine/app/version/config/control)
    firmware/       # Firmware management & upgrades
    task/           # Task management (unified DRC/firmware task view)
    docs/           # Device model reference documentation (from GitHub)
    system/         # System management (users, roles, permissions, org, audit logs)
    version/        # Version info
  cmdutil/          # Shared list flags (cursor/limit), query builder
  config/           # Config file I/O, context model
  debug/            # Debug output (--debug / DEVICEMANAGER_DEBUG)
  factory/          # Dependency injection factory
  iostreams/        # Terminal output, formatters (JSON/Table/YAML/jq)
```
