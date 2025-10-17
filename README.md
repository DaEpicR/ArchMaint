# ArchMaint

A comprehensive CLI maintenance tool for Arch Linux systems, providing automated package management, system health monitoring, and preventive maintenance capabilities with safety-first design principles.

[![Go Version](https://img.shields.io/badge/go-1.21+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Version](https://img.shields.io/badge/version-1.1.0-blue.svg)](CHANGELOG.md)

## Overview

ArchMaint simplifies routine Arch Linux system maintenance by consolidating essential tasks into an intuitive CLI interface. Designed for both automated workflows and interactive use, it emphasizes safety through dry-run capabilities, confirmation prompts, and automatic backups.

## Core Features

### System Maintenance
- **Package Management**: Update, clean cache, remove orphans
- **System Cleanup**: Logs, temporary files, user cache
- **Automated Backups**: Pre-update snapshots with restore functionality
- **Health Monitoring**: 6-point system diagnostics with scoring

### Safety & Control
- **Dry-run Mode**: Preview all changes before execution
- **Safe Mode**: Enhanced confirmations for destructive operations
- **Dependency Checking**: Validation before package operations
- **Reboot Detection**: Kernel update notifications

### Advanced Features
- **Package Search**: Interactive repository browsing
- **Btrfs Snapshots**: System rollback capability
- **Configuration Management**: Customizable retention policies
- **Progress Indicators**: Real-time operation feedback

## Installation

### Requirements
- Arch Linux
- Go 1.21+
- `sudo` privileges
- `pacman-contrib` (for paccache)

### Build from Source

```bash
# Clone repository
git clone https://github.com/yourusername/archmaint
cd archmaint

# Install dependencies
go mod download
go mod tidy

# Build
go build -o archmaint main.go

# Install system-wide (optional)
sudo cp archmaint /usr/local/bin/
```

### Using Make

```bash
make build      # Build binary
make install    # Build and install
make run        # Build and run
make clean      # Remove artifacts
```

## Quick Start

### Interactive Mode
```bash
archmaint                    # Launch menu
archmaint --dry-run         # Preview mode
archmaint --safe            # Extra confirmations
```

### Common Commands
```bash
archmaint status            # System overview
archmaint update            # Update packages (with backup)
archmaint clean             # Clean cache and logs
archmaint orphans           # Remove unused packages
archmaint health            # System diagnostics
archmaint maintenance       # Full maintenance routine
```

### Advanced Usage
```bash
archmaint search <package>  # Search packages
archmaint backup            # Create manual backup
archmaint restore           # Restore from backup
archmaint snapshot          # Create btrfs snapshot
archmaint config            # Manage settings
```

## Command Reference

| Command | Alias | Description |
|---------|-------|-------------|
| `status` | `s` | Display system information and status |
| `update` | `u` | Update packages with optional backup |
| `clean` | `c` | Clean cache, logs, and temporary files |
| `orphans` | `o` | Identify and remove unused packages |
| `services` | `sv` | Monitor systemd service health |
| `logs` | `l` | View recent system logs |
| `health` | `h` | Run comprehensive health check |
| `maintenance` | `m` | Execute full maintenance routine |
| `search` | `se` | Search package repositories |
| `backup` | `b` | Create system package backup |
| `restore` | `r` | Restore from previous backup |
| `snapshot` | `sn` | Create btrfs filesystem snapshot |
| `config` | `cfg` | Configure tool settings |
| `help` | `-h` | Display help information |
| `version` | `-v` | Show version information |

## Configuration

### Default Behavior
- Dry-run: Disabled
- Safe mode: Disabled
- Backups: Enabled
- Cache retention: 30 days
- Log retention: 7 days

### Configuration File
```
~/.config/archmaint/config.conf
```

### Settings
```bash
archmaint config            # Interactive configuration
```

Available options:
- Toggle dry-run mode
- Toggle safe mode
- Enable/disable automatic backups
- Set cache retention period
- Set log retention period
- Toggle verbose logging

### Backup Locations
```
~/.archmaint/backups/       # Backup storage
```

## Health Check Metrics

The `health` command evaluates:
1. **Disk Space** - Root partition usage < 90%
2. **Memory Usage** - Available memory > 10%
3. **Failed Services** - No systemd service failures
4. **Package Database** - Database integrity verified
5. **System Errors** - Minimal errors in recent logs
6. **Security Updates** - No critical package updates pending

Output: Health score (0-100%) with detailed results

## Safety Mechanisms

### Confirmation System
- Non-destructive operations: Yes/No confirmation
- Dangerous operations: Yes/No confirmation with warnings
- Safe mode: Requires "yes" phrase for destructive ops

### Backup System
- Automatic pre-update backups
- Stores: explicit packages, all packages, AUR packages
- Manual backup creation with timestamp
- Restore capability with version selection

### Dry-run Preview
```bash
archmaint --dry-run <command>   # Preview without changes
```

Shows:
- Commands that would execute
- Files that would be modified
- Changes without applying them

## Usage Patterns

### Daily Maintenance
```bash
archmaint status            # Check system status
archmaint health            # Run diagnostics
```

### Weekly Maintenance
```bash
archmaint --dry-run update  # Preview updates
archmaint maintenance       # Full maintenance routine
```

### Before Major Changes
```bash
archmaint backup            # Create backup
archmaint snapshot          # Create snapshot (btrfs)
archmaint --dry-run update  # Preview changes
```

### Troubleshooting
```bash
archmaint health            # Identify issues
archmaint logs              # Review system logs
archmaint services          # Check service status
```

### Development Setup
```bash
git clone <repository>
cd archmaint
make deps
make build
./archmaint
```

## Troubleshooting

### Common Issues

**Build Error: "command not found: go"**
```bash
sudo pacman -S go
```

**Module Not Found**
```bash
go mod tidy
go mod download
go build -o archmaint main.go
```

**Permission Denied**
```bash
# Most operations require sudo
sudo archmaint update
sudo archmaint maintenance
```

**Missing Dependencies**
```bash
sudo pacman -S pacman-contrib lm-sensors
sudo sensors-detect
```

## Performance Considerations

- Typical operation time: 5-30 seconds depending on system
- Backup creation: 1-5 seconds
- Health check: 10-15 seconds
- Full maintenance: 2-5 minutes

## Security Notes

- Requires `sudo` for system modifications
- Confirmation prompts for destructive operations
- Backups stored in user home directory
- No remote connections or network access
- All operations logged locally

## Support

- Issues: GitHub Issues
- Documentation: This README
- Help: `archmaint help`

## Changelog

### v1.1.0
- Dry-run mode
- Safe mode
- Backup/restore system
- Package search
- Btrfs snapshots
- Configuration manager
- Progress indicators
- Enhanced health checks

### v1.0.0
- Initial release
- Basic maintenance tasks
- System status display
- Service monitoring

## Roadmap

**Current**: v1.1.0 - Stable with enhanced safety
**Next**: v1.2.0 - AUR support (maybe), notifications
**Future**: v2.0.0 - cooking

---

**Maintained by**: [DaEpicR]
**Status**: Active Development
