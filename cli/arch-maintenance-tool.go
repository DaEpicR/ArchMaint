package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/schollz/progressbar/v3"
)

// ArchMaintenance represents the main application
type ArchMaintenance struct {
	version string
	config  *Config
}

// Config holds application configuration
type Config struct {
	DryRun               bool
	AutoConfirm          bool
	BackupEnabled        bool
	BackupPath           string
	CacheRetentionDays   int
	LogRetentionDays     int
	NotificationsEnabled bool
	VerboseMode          bool
	SafeMode             bool
	CustomCommands       map[string]CustomCommand
}

// CustomCommand represents a user-defined command
type CustomCommand struct {
	Name        string
	Description string
	Command     []string
	Dangerous   bool
}

// Task represents a maintenance task
type Task struct {
	Name        string
	Description string
	Command     []string
	Dangerous   bool
	Frequency   string
}

// SystemInfo holds system information
type SystemInfo struct {
	Kernel      string
	Uptime      string
	LoadAvg     string
	MemoryUsage string
	DiskUsage   string
	CPUTemp     string
}

// Colors for beautiful output
var (
	headerColor   = color.New(color.FgCyan, color.Bold)
	successColor  = color.New(color.FgGreen, color.Bold)
	warningColor  = color.New(color.FgYellow, color.Bold)
	errorColor    = color.New(color.FgRed, color.Bold)
	infoColor     = color.New(color.FgBlue)
	dangerColor   = color.New(color.FgRed, color.Bold, color.BgYellow)
	progressColor = color.New(color.FgMagenta)
)

func main() {
	app := &ArchMaintenance{
		version: "1.1.0",
		config:  loadDefaultConfig(),
	}

	// Load user config if exists
	if err := app.loadConfig(); err == nil {
		infoColor.Println("Loaded custom configuration")
	}

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "status", "s":
			app.showSystemStatus()
		case "update", "u":
			app.systemUpdate()
		case "clean", "c":
			app.systemClean()
		case "orphans", "o":
			app.removeOrphans()
		case "services", "sv":
			app.showServices()
		case "logs", "l":
			app.showLogs()
		case "health", "h":
			app.systemHealthCheck()
		case "maintenance", "m":
			app.fullMaintenance()
		case "backup", "b":
			app.createBackup()
		case "restore", "r":
			app.restoreBackup()
		case "snapshot", "sn":
			app.createSnapshot()
		case "config", "cfg":
			app.configManager()
		case "search", "se":
			if len(os.Args) > 2 {
				app.searchPackages(os.Args[2])
			} else {
				errorColor.Println("Please provide a search term")
			}
		case "help", "--help", "-h":
			app.showHelp()
		case "version", "--version", "-v":
			app.showVersion()
		case "--dry-run":
			app.config.DryRun = true
			infoColor.Println("DRY RUN MODE: No changes will be made")
			if len(os.Args) > 2 {
				os.Args = append(os.Args[:1], os.Args[2:]...)
				main()
			}
		case "--safe":
			app.config.SafeMode = true
			successColor.Println("SAFE MODE: Extra confirmations enabled")
			if len(os.Args) > 2 {
				os.Args = append(os.Args[:1], os.Args[2:]...)
				main()
			}
		default:
			app.showHelp()
		}
	} else {
		app.showMainMenu()
	}
}

func loadDefaultConfig() *Config {
	homeDir, _ := os.UserHomeDir()
	return &Config{
		DryRun:               false,
		AutoConfirm:          false,
		BackupEnabled:        true,
		BackupPath:           filepath.Join(homeDir, ".archmaint/backups"),
		CacheRetentionDays:   30,
		LogRetentionDays:     7,
		NotificationsEnabled: true,
		VerboseMode:          false,
		SafeMode:             false,
		CustomCommands:       make(map[string]CustomCommand),
	}
}

func (a *ArchMaintenance) loadConfig() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(homeDir, ".config/archmaint/config.conf")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return err
	}

	return nil
}

func (a *ArchMaintenance) showBanner() {
	banner := `
    ██████╗  ██████╗██╗  ██╗    ███╗   ███╗ █████╗ ██╗███╗   ██╗████████╗
   ██╔══██╗██╔════╝██║  ██║    ████╗ ████║██╔══██╗██║████╗  ██║╚══██╔══╝
   ███████║██║     ███████║    ██╔████╔██║███████║██║██╔██╗ ██║   ██║
   ██╔══██║██║     ██╔══██║    ██║╚██╔╝██║██╔══██║██║██║╚██╗██║   ██║
   ██║  ██║╚██████╗██║  ██║    ██║ ╚═╝ ██║██║  ██║██║██║ ╚████║   ██║
   ╚═╝  ╚═╝ ╚═════╝╚═╝  ╚═╝    ╚═╝     ╚═╝╚═╝  ╚═╝╚═╝╚═╝  ╚═══╝   ╚═╝
`
	headerColor.Println(banner)
	infoColor.Printf("                    Arch Linux Maintenance Tool v%s\n", a.version)
	if a.config.DryRun {
		warningColor.Println("                          [DRY RUN MODE]")
	}
	if a.config.SafeMode {
		successColor.Println("                           [SAFE MODE]")
	}
	fmt.Println()
}

func (a *ArchMaintenance) showMainMenu() {
	a.showBanner()

	options := [][]string{
		{"1", "System Status", "Show system information and status"},
		{"2", "System Update", "Update system packages (with backup)"},
		{"3", "System Clean", "Clean package cache and temporary files"},
		{"4", "Remove Orphans", "Remove orphaned packages"},
		{"5", "System Services", "View system services status"},
		{"6", "System Logs", "View recent system logs"},
		{"7", "Health Check", "Comprehensive system health check"},
		{"8", "Full Maintenance", "Run complete maintenance routine"},
		{"9", "Search Packages", "Search for packages"},
		{"10", "Create Backup", "Backup package list and important files"},
		{"11", "Create Snapshot", "Create system snapshot (btrfs)"},
		{"12", "Configuration", "Manage settings and preferences"},
		{"h", "Help", "Show help information"},
		{"0", "Exit", "Exit the application"},
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Option", "Command", "Description"})
	table.SetBorder(false)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetTablePadding("\t")
	table.SetNoWhiteSpace(true)

	for _, option := range options {
		table.Append(option)
	}
	table.Render()

	fmt.Print("\nEnter your choice: ")
	reader := bufio.NewReader(os.Stdin)
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	switch choice {
	case "1":
		a.showSystemStatus()
	case "2":
		a.systemUpdate()
	case "3":
		a.systemClean()
	case "4":
		a.removeOrphans()
	case "5":
		a.showServices()
	case "6":
		a.showLogs()
	case "7":
		a.systemHealthCheck()
	case "8":
		a.fullMaintenance()
	case "9":
		fmt.Print("Enter search term: ")
		term, _ := reader.ReadString('\n')
		a.searchPackages(strings.TrimSpace(term))
	case "10":
		a.createBackup()
	case "11":
		a.createSnapshot()
	case "12":
		a.configManager()
	case "h", "H":
		a.showHelp()
	case "0":
		successColor.Println("Goodbye! Keep your Arch system running smoothly!")
		os.Exit(0)
	default:
		errorColor.Println("Invalid choice. Please try again.")
		time.Sleep(2 * time.Second)
		a.showMainMenu()
	}
}

func (a *ArchMaintenance) showSystemStatus() {
	headerColor.Println("\n=== SYSTEM STATUS ===")

	info := a.getSystemInfo()

	data := [][]string{
		{"Kernel", info.Kernel},
		{"Uptime", info.Uptime},
		{"Load Average", info.LoadAvg},
		{"Memory Usage", info.MemoryUsage},
		{"Root Disk Usage", info.DiskUsage},
		{"CPU Temperature", info.CPUTemp},
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetBorder(false)
	table.SetCenterSeparator("")
	table.SetColumnSeparator(" : ")
	table.SetRowSeparator("")
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	for _, row := range data {
		table.Append(row)
	}
	table.Render()

	fmt.Println()
	a.showPackageInfo()

	fmt.Println()
	a.showDiskHealth()

	a.waitForContinue()
}

func (a *ArchMaintenance) getSystemInfo() SystemInfo {
	info := SystemInfo{}

	if output, err := exec.Command("uname", "-r").Output(); err == nil {
		info.Kernel = strings.TrimSpace(string(output))
	}

	if output, err := exec.Command("uptime", "-p").Output(); err == nil {
		info.Uptime = strings.TrimSpace(string(output))
	}

	if output, err := exec.Command("cat", "/proc/loadavg").Output(); err == nil {
		fields := strings.Fields(string(output))
		if len(fields) >= 3 {
			info.LoadAvg = fmt.Sprintf("%s %s %s", fields[0], fields[1], fields[2])
		}
	}

	if output, err := exec.Command("free", "-h").Output(); err == nil {
		lines := strings.Split(string(output), "\n")
		if len(lines) > 1 {
			fields := strings.Fields(lines[1])
			if len(fields) >= 3 {
				info.MemoryUsage = fmt.Sprintf("%s / %s", fields[2], fields[1])
			}
		}
	}

	if output, err := exec.Command("df", "-h", "/").Output(); err == nil {
		lines := strings.Split(string(output), "\n")
		if len(lines) > 1 {
			fields := strings.Fields(lines[1])
			if len(fields) >= 5 {
				info.DiskUsage = fmt.Sprintf("%s / %s (%s)", fields[2], fields[1], fields[4])
			}
		}
	}

	if output, err := exec.Command("bash", "-c", "sensors 2>/dev/null | grep -i 'Package id 0' | awk '{print $4}' || echo 'N/A'").Output(); err == nil {
		temp := strings.TrimSpace(string(output))
		if temp == "" {
			temp = "N/A"
		}
		info.CPUTemp = temp
	} else {
		info.CPUTemp = "N/A"
	}

	return info
}

func (a *ArchMaintenance) showPackageInfo() {
	infoColor.Println("Package Information:")

	if output, err := exec.Command("pacman", "-Q").Output(); err == nil {
		count := len(strings.Split(strings.TrimSpace(string(output)), "\n"))
		fmt.Printf("  Installed packages: %d\n", count)
	}

	if output, err := exec.Command("pacman", "-Qe").Output(); err == nil {
		count := len(strings.Split(strings.TrimSpace(string(output)), "\n"))
		fmt.Printf("  Explicitly installed: %d\n", count)
	}

	if output, err := exec.Command("pacman", "-Qtdq").Output(); err == nil {
		orphans := strings.TrimSpace(string(output))
		if orphans != "" {
			count := len(strings.Split(orphans, "\n"))
			warningColor.Printf("  Orphaned packages: %d\n", count)
		} else {
			fmt.Printf("  Orphaned packages: 0\n")
		}
	}

	if output, err := exec.Command("pacman", "-Qu").Output(); err == nil {
		updates := strings.TrimSpace(string(output))
		if updates != "" {
			count := len(strings.Split(updates, "\n"))
			warningColor.Printf("  Packages to update: %d\n", count)
		} else {
			successColor.Printf("  Packages to update: 0\n")
		}
	}
}

func (a *ArchMaintenance) showDiskHealth() {
	infoColor.Println("Disk Health:")

	if output, err := exec.Command("df", "-h", "-x", "tmpfs", "-x", "devtmpfs").Output(); err == nil {
		lines := strings.Split(string(output), "\n")
		for i, line := range lines {
			if i == 0 || strings.TrimSpace(line) == "" {
				continue
			}
			fields := strings.Fields(line)
			if len(fields) >= 5 {
				usage := strings.TrimSuffix(fields[4], "%")
				if val := parseInt(usage); val > 90 {
					errorColor.Printf("  WARNING %s: %s used (Critical!)\n", fields[5], fields[4])
				} else if val > 80 {
					warningColor.Printf("  WARNING %s: %s used\n", fields[5], fields[4])
				} else {
					fmt.Printf("  OK %s: %s used\n", fields[5], fields[4])
				}
			}
		}
	}
}

func parseInt(s string) int {
	var result int
	fmt.Sscanf(s, "%d", &result)
	return result
}

func (a *ArchMaintenance) systemUpdate() {
	headerColor.Println("\n=== SYSTEM UPDATE ===")

	if a.config.DryRun {
		warningColor.Println("DRY RUN: Showing what would be updated")
	}

	if a.config.BackupEnabled && !a.config.DryRun {
		if a.confirmAction("Create backup before updating?", false) {
			a.createBackup()
		}
	}

	if !a.confirmAction("This will update your system. Continue?", false) {
		return
	}

	infoColor.Println("Syncing package databases...")
	if !a.config.DryRun {
		a.runCommandWithProgress("sudo", "pacman", "-Sy")
	} else {
		fmt.Println("  Would run: sudo pacman -Sy")
	}

	infoColor.Println("\nChecking for updates...")
	cmd := exec.Command("pacman", "-Qu")
	output, err := cmd.Output()

	if err != nil || strings.TrimSpace(string(output)) == "" {
		successColor.Println("System is up to date!")
		a.waitForContinue()
		return
	}

	updates := strings.Split(strings.TrimSpace(string(output)), "\n")
	fmt.Printf("\nAvailable updates (%d packages):\n", len(updates))

	displayCount := 20
	for i, update := range updates {
		if i >= displayCount {
			infoColor.Printf("... and %d more packages\n", len(updates)-displayCount)
			break
		}
		fmt.Printf("  • %s\n", update)
	}

	if a.confirmAction(fmt.Sprintf("Proceed with updating %d packages?", len(updates)), false) {
		infoColor.Println("Updating system...")
		if !a.config.DryRun {
			a.runCommandWithProgress("sudo", "pacman", "-Su", "--noconfirm")
			successColor.Println("System update completed!")

			if a.needsReboot() {
				warningColor.Println("\nSystem reboot recommended to apply updates")
			}
		} else {
			fmt.Println("  Would run: sudo pacman -Su")
		}
	}

	a.waitForContinue()
}

func (a *ArchMaintenance) needsReboot() bool {
	cmd := exec.Command("bash", "-c", "uname -r | cut -d'-' -f1")
	currentKernel, _ := cmd.Output()

	cmd = exec.Command("bash", "-c", "pacman -Q linux 2>/dev/null | awk '{print $2}' | cut -d'-' -f1")
	installedKernel, _ := cmd.Output()

	return strings.TrimSpace(string(currentKernel)) != strings.TrimSpace(string(installedKernel))
}

func (a *ArchMaintenance) systemClean() {
	headerColor.Println("\n=== SYSTEM CLEAN ===")

	if a.config.DryRun {
		warningColor.Println("DRY RUN: Showing what would be cleaned")
	}

	tasks := []Task{
		{
			Name:        "Package Cache",
			Description: fmt.Sprintf("Clean pacman cache (keep %d days)", a.config.CacheRetentionDays),
			Command:     []string{"sudo", "paccache", "-rk3"},
			Dangerous:   false,
			Frequency:   "Weekly",
		},
		{
			Name:        "Uninstalled Package Cache",
			Description: "Remove cache for uninstalled packages",
			Command:     []string{"sudo", "paccache", "-ruk0"},
			Dangerous:   false,
			Frequency:   "Weekly",
		},
		{
			Name:        "System Logs",
			Description: fmt.Sprintf("Clean old journal logs (keep %d days)", a.config.LogRetentionDays),
			Command:     []string{"sudo", "journalctl", fmt.Sprintf("--vacuum-time=%dd", a.config.LogRetentionDays)},
			Dangerous:   false,
			Frequency:   "Weekly",
		},
		{
			Name:        "Temporary Files",
			Description: "Clean /tmp and /var/tmp (older than 7 days)",
			Command:     []string{"sudo", "find", "/tmp", "/var/tmp", "-type", "f", "-atime", "+7", "-delete"},
			Dangerous:   false,
			Frequency:   "Daily",
		},
		{
			Name:        "User Cache",
			Description: "Clean user cache directories",
			Command:     []string{"bash", "-c", "find ~/.cache -type f -atime +30 -delete 2>/dev/null || true"},
			Dangerous:   false,
			Frequency:   "Monthly",
		},
	}

	for _, task := range tasks {
		fmt.Printf("\n%s (%s)\n", task.Name, task.Frequency)
		fmt.Printf("Description: %s\n", task.Description)

		if task.Dangerous {
			dangerColor.Printf("[CAUTION] This action can be dangerous!\n")
		}

		if a.confirmAction(fmt.Sprintf("Run %s cleanup?", task.Name), task.Dangerous) {
			if !a.config.DryRun {
				a.runCommandWithProgress(task.Command[0], task.Command[1:]...)
			} else {
				fmt.Printf("  Would run: %s\n", strings.Join(task.Command, " "))
			}
		}
	}

	successColor.Println("\nSystem cleanup completed!")
	a.waitForContinue()
}

func (a *ArchMaintenance) removeOrphans() {
	headerColor.Println("\n=== REMOVE ORPHANED PACKAGES ===")

	cmd := exec.Command("pacman", "-Qtdq")
	output, err := cmd.Output()

	if err != nil || strings.TrimSpace(string(output)) == "" {
		successColor.Println("No orphaned packages found!")
		return
	}

	orphans := strings.TrimSpace(string(output))
	orphanList := strings.Split(orphans, "\n")

	fmt.Printf("Found %d orphaned packages:\n", len(orphanList))
	for i, pkg := range orphanList {
		if i >= 20 {
			infoColor.Printf("... and %d more packages\n", len(orphanList)-20)
			break
		}
		fmt.Printf("  • %s\n", pkg)
	}

	if a.confirmAction(fmt.Sprintf("Remove these %d orphaned packages?", len(orphanList)), true) {
		if !a.config.DryRun {
			args := append([]string{"pacman", "-Rns", "--noconfirm"}, orphanList...)
			a.runCommandWithProgress("sudo", args...)
			successColor.Println("Orphaned packages removed!")
		} else {
			fmt.Println("  Would run: sudo pacman -Rns " + strings.Join(orphanList, " "))
		}
	}
}

func (a *ArchMaintenance) showServices() {
	headerColor.Println("\n=== SYSTEM SERVICES ===")

	infoColor.Println("Failed services:")
	a.runCommand("systemctl", "--failed")

	fmt.Println()
	infoColor.Println("Service status summary:")

	cmd := exec.Command("systemctl", "list-units", "--type=service", "--all", "--no-pager")
	output, _ := cmd.Output()

	active := 0
	failed := 0
	inactive := 0

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "active") && strings.Contains(line, "running") {
			active++
		} else if strings.Contains(line, "failed") {
			failed++
		} else if strings.Contains(line, "inactive") {
			inactive++
		}
	}

	fmt.Printf("  Active: %d\n", active)
	if failed > 0 {
		errorColor.Printf("  Failed: %d\n", failed)
	} else {
		successColor.Printf("  Failed: 0\n")
	}
	fmt.Printf("  Inactive: %d\n", inactive)

	a.waitForContinue()
}

func (a *ArchMaintenance) showLogs() {
	headerColor.Println("\n=== SYSTEM LOGS ===")

	infoColor.Println("Recent critical and error logs:")
	a.runCommand("journalctl", "-p", "3", "-x", "--no-pager", "--since", "today", "-n", "50")

	fmt.Println()
	infoColor.Println("Boot messages:")
	a.runCommand("journalctl", "-b", "--no-pager", "-n", "20")

	a.waitForContinue()
}

func (a *ArchMaintenance) systemHealthCheck() {
	headerColor.Println("\n=== SYSTEM HEALTH CHECK ===")

	checks := []struct {
		name string
		fn   func() bool
		desc string
	}{
		{"Disk Space", a.checkDiskSpace, "Checking available disk space"},
		{"Memory Usage", a.checkMemory, "Checking memory usage"},
		{"Failed Services", a.checkServices, "Checking for failed services"},
		{"Package Database", a.checkPackageDB, "Verifying package database integrity"},
		{"System Errors", a.checkSystemErrors, "Checking for recent system errors"},
		{"Security Updates", a.checkSecurityUpdates, "Checking for security updates"},
	}

	passedChecks := 0
	totalChecks := len(checks)

	for i, check := range checks {
		fmt.Printf("\n[%d/%d] %s\n", i+1, totalChecks, check.name)
		infoColor.Printf("     %s...\n", check.desc)

		if check.fn() {
			successColor.Println("     PASSED")
			passedChecks++
		} else {
			errorColor.Println("     FAILED")
		}
	}

	fmt.Println()
	fmt.Println(strings.Repeat("=", 50))

	percentage := (passedChecks * 100) / totalChecks
	if percentage == 100 {
		successColor.Printf("Health Score: %d%% (%d/%d checks passed)\n", percentage, passedChecks, totalChecks)
	} else if percentage >= 80 {
		warningColor.Printf("Health Score: %d%% (%d/%d checks passed)\n", percentage, passedChecks, totalChecks)
	} else {
		errorColor.Printf("Health Score: %d%% (%d/%d checks passed)\n", percentage, passedChecks, totalChecks)
	}

	a.waitForContinue()
}

func (a *ArchMaintenance) checkDiskSpace() bool {
	cmd := exec.Command("df", "-h", "/")
	output, _ := cmd.Output()
	lines := strings.Split(string(output), "\n")
	if len(lines) > 1 {
		fields := strings.Fields(lines[1])
		if len(fields) >= 5 {
			usage := strings.TrimSuffix(fields[4], "%")
			if parseInt(usage) < 90 {
				return true
			}
		}
	}
	return false
}

func (a *ArchMaintenance) checkMemory() bool {
	cmd := exec.Command("free")
	output, _ := cmd.Output()
	lines := strings.Split(string(output), "\n")
	if len(lines) > 1 {
		fields := strings.Fields(lines[1])
		if len(fields) >= 3 {
			total := parseInt(fields[1])
			used := parseInt(fields[2])
			if total > 0 && (used*100/total) < 90 {
				return true
			}
		}
	}
	return false
}

func (a *ArchMaintenance) checkServices() bool {
	cmd := exec.Command("systemctl", "--failed", "--no-legend")
	output, _ := cmd.Output()
	return strings.TrimSpace(string(output)) == ""
}

func (a *ArchMaintenance) checkPackageDB() bool {
	cmd := exec.Command("pacman", "-Dk")
	err := cmd.Run()
	return err == nil
}

func (a *ArchMaintenance) checkSystemErrors() bool {
	cmd := exec.Command("journalctl", "-p", "3", "--since", "today", "--no-pager")
	output, _ := cmd.Output()
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	return len(lines) < 5 || (len(lines) == 1 && lines[0] == "")
}

func (a *ArchMaintenance) checkSecurityUpdates() bool {
	cmd := exec.Command("pacman", "-Qu")
	output, _ := cmd.Output()
	updates := strings.TrimSpace(string(output))

	if updates == "" {
		return true
	}

	criticalPackages := []string{"linux", "systemd", "glibc", "openssl"}
	for _, pkg := range criticalPackages {
		if strings.Contains(updates, pkg) {
			return false
		}
	}

	return true
}

func (a *ArchMaintenance) fullMaintenance() {
	headerColor.Println("\n=== FULL SYSTEM MAINTENANCE ===")

	dangerColor.Println("WARNING: This will perform comprehensive system maintenance!")
	fmt.Println("\nThis routine includes:")
	fmt.Println("  1. Create system backup")
	fmt.Println("  2. System update")
	fmt.Println("  3. Package cache cleanup")
	fmt.Println("  4. Orphaned package removal")
	fmt.Println("  5. Log cleanup")
	fmt.Println("  6. Temporary file cleanup")
	fmt.Println("  7. System health check")

	if !a.confirmAction("Are you sure you want to continue with full maintenance?", true) {
		return
	}

	steps := []struct {
		name string
		fn   func()
	}{
		{"Creating Backup", a.createBackup},
		{"Updating System", a.systemUpdate},
		{"Cleaning System", a.systemClean},
		{"Removing Orphans", a.removeOrphans},
		{"Running Health Check", a.systemHealthCheck},
	}

	for i, step := range steps {
		headerColor.Printf("\n[Step %d/%d] %s\n", i+1, len(steps), step.name)
		step.fn()
	}

	successColor.Println("\nFull maintenance completed successfully!")
	infoColor.Println("Recommendation: Consider rebooting if kernel was updated.")

	if a.needsReboot() {
		warningColor.Println("Kernel update detected - reboot recommended!")
		if a.confirmAction("Reboot now?", true) {
			if !a.config.DryRun {
				a.runCommand("sudo", "reboot")
			} else {
				fmt.Println("  Would run: sudo reboot")
			}
		}
	}

	a.waitForContinue()
}

func (a *ArchMaintenance) createBackup() {
	headerColor.Println("\n=== CREATE BACKUP ===")

	if a.config.DryRun {
		warningColor.Println("DRY RUN: Showing what would be backed up")
	}

	if err := os.MkdirAll(a.config.BackupPath, 0755); err != nil {
		errorColor.Printf("Failed to create backup directory: %v\n", err)
		return
	}

	timestamp := time.Now().Format("2006-01-02_15-04-05")
	backupDir := filepath.Join(a.config.BackupPath, timestamp)

	if !a.config.DryRun {
		if err := os.MkdirAll(backupDir, 0755); err != nil {
			errorColor.Printf("Failed to create backup directory: %v\n", err)
			return
		}
	}

	infoColor.Println("Creating backup...")

	backupItems := []struct {
		name string
		cmd  []string
		file string
	}{
		{
			"Package list (explicitly installed)",
			[]string{"pacman", "-Qqe"},
			"packages_explicit.txt",
		},
		{
			"Package list (all installed)",
			[]string{"pacman", "-Qq"},
			"packages_all.txt",
		},
		{
			"Package list (foreign/AUR)",
			[]string{"pacman", "-Qqm"},
			"packages_foreign.txt",
		},
	}

	bar := progressbar.NewOptions(len(backupItems),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetWidth(50),
		progressbar.OptionSetDescription("Backing up..."),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))

	for _, item := range backupItems {
		if a.config.DryRun {
			fmt.Printf("  Would backup: %s\n", item.name)
		} else {
			cmd := exec.Command(item.cmd[0], item.cmd[1:]...)
			output, err := cmd.Output()
			if err == nil {
				outputFile := filepath.Join(backupDir, item.file)
				if err := os.WriteFile(outputFile, output, 0644); err == nil {
					if a.config.VerboseMode {
						successColor.Printf("  Backed up: %s\n", item.name)
					}
				}
			}
		}
		bar.Add(1)
	}

	fmt.Println()

	if !a.config.DryRun {
		successColor.Printf("Backup created: %s\n", backupDir)

		cmd := exec.Command("du", "-sh", backupDir)
		if output, err := cmd.Output(); err == nil {
			size := strings.Fields(string(output))[0]
			infoColor.Printf("  Backup size: %s\n", size)
		}

		a.listBackups()
	}
}

func (a *ArchMaintenance) listBackups() {
	files, err := os.ReadDir(a.config.BackupPath)
	if err != nil {
		return
	}

	if len(files) > 0 {
		fmt.Println("\nRecent backups:")
		count := 0
		for i := len(files) - 1; i >= 0 && count < 5; i-- {
			if files[i].IsDir() {
				info, _ := files[i].Info()
				fmt.Printf("  - %s (%s)\n", files[i].Name(),
					info.ModTime().Format("2006-01-02 15:04:05"))
				count++
			}
		}
	}
}

func (a *ArchMaintenance) restoreBackup() {
	headerColor.Println("\n=== RESTORE BACKUP ===")

	files, err := os.ReadDir(a.config.BackupPath)
	if err != nil || len(files) == 0 {
		errorColor.Println("No backups found!")
		return
	}

	fmt.Println("Available backups:")
	backups := []os.DirEntry{}
	for i := len(files) - 1; i >= 0; i-- {
		if files[i].IsDir() {
			backups = append(backups, files[i])
			fmt.Printf("  %d. %s\n", len(backups), files[i].Name())
		}
	}

	fmt.Print("\nSelect backup to restore (0 to cancel): ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	choice := parseInt(strings.TrimSpace(input))

	if choice <= 0 || choice > len(backups) {
		infoColor.Println("Restore cancelled.")
		return
	}

	selectedBackup := backups[choice-1]
	backupPath := filepath.Join(a.config.BackupPath, selectedBackup.Name())

	dangerColor.Println("\nWARNING: This will install packages from the backup!")
	if !a.confirmAction("Continue with restore?", true) {
		return
	}

	pkgFile := filepath.Join(backupPath, "packages_explicit.txt")
	if data, err := os.ReadFile(pkgFile); err == nil {
		packages := strings.Split(strings.TrimSpace(string(data)), "\n")
		infoColor.Printf("Restoring %d packages...\n", len(packages))

		if !a.config.DryRun {
			args := append([]string{"pacman", "-S", "--needed", "--noconfirm"}, packages...)
			a.runCommandWithProgress("sudo", args...)
			successColor.Println("Packages restored!")
		} else {
			fmt.Println("  Would install:", len(packages), "packages")
		}
	}

	a.waitForContinue()
}

func (a *ArchMaintenance) createSnapshot() {
	headerColor.Println("\n=== CREATE SYSTEM SNAPSHOT ===")

	cmd := exec.Command("findmnt", "-n", "-o", "FSTYPE", "/")
	output, err := cmd.Output()

	if err != nil || !strings.Contains(string(output), "btrfs") {
		warningColor.Println("Root filesystem is not btrfs")
		infoColor.Println("Snapshots are only supported on btrfs filesystems.")
		infoColor.Println("Consider using Timeshift or Snapper for advanced snapshot management.")
		a.waitForContinue()
		return
	}

	timestamp := time.Now().Format("2006-01-02_15-04-05")
	snapshotName := fmt.Sprintf("archmaint_%s", timestamp)
	snapshotPath := fmt.Sprintf("/.snapshots/%s", snapshotName)

	if a.config.DryRun {
		fmt.Printf("  Would create snapshot: %s\n", snapshotPath)
		return
	}

	infoColor.Println("Creating btrfs snapshot...")

	exec.Command("sudo", "mkdir", "-p", "/.snapshots").Run()

	cmd = exec.Command("sudo", "btrfs", "subvolume", "snapshot", "/", snapshotPath)
	if err := cmd.Run(); err == nil {
		successColor.Printf("Snapshot created: %s\n", snapshotPath)

		cmd = exec.Command("sudo", "btrfs", "subvolume", "list", "/")
		if output, err := cmd.Output(); err == nil {
			lines := strings.Split(string(output), "\n")
			count := 0
			fmt.Println("\nRecent snapshots:")
			for i := len(lines) - 1; i >= 0 && count < 5; i-- {
				if strings.Contains(lines[i], "archmaint_") {
					fmt.Printf("  - %s\n", lines[i])
					count++
				}
			}
		}
	} else {
		errorColor.Printf("Failed to create snapshot: %v\n", err)
	}

	a.waitForContinue()
}

func (a *ArchMaintenance) searchPackages(query string) {
	headerColor.Printf("\n=== SEARCH PACKAGES: %s ===\n", query)

	infoColor.Println("Searching in official repositories...")
	cmd := exec.Command("pacman", "-Ss", query)
	output, _ := cmd.Output()

	if strings.TrimSpace(string(output)) != "" {
		lines := strings.Split(string(output), "\n")
		count := 0
		for i := 0; i < len(lines) && count < 20; i++ {
			line := lines[i]
			if strings.TrimSpace(line) == "" {
				continue
			}

			if strings.HasPrefix(line, " ") {
				fmt.Println(line)
			} else {
				if strings.Contains(line, "[installed]") {
					successColor.Print(line[:strings.Index(line, "[installed]")])
					infoColor.Println(" [installed]")
				} else {
					fmt.Println(line)
				}
				count++
			}
		}

		if len(lines) > 40 {
			infoColor.Printf("\n... and more results (showing first 20)\n")
		}
	} else {
		warningColor.Println("No packages found in official repositories.")
	}

	cmd = exec.Command("pacman", "-Qi", query)
	if err := cmd.Run(); err == nil {
		fmt.Println()
		successColor.Printf("Package '%s' is installed\n", query)

		cmd = exec.Command("pacman", "-Qi", query)
		output, _ := cmd.Output()
		fmt.Println(string(output))
	}

	a.waitForContinue()
}

func (a *ArchMaintenance) configManager() {
	headerColor.Println("\n=== CONFIGURATION MANAGER ===")

	fmt.Println("Current Configuration:")
	fmt.Printf("  Dry Run Mode: %v\n", a.config.DryRun)
	fmt.Printf("  Safe Mode: %v\n", a.config.SafeMode)
	fmt.Printf("  Backup Enabled: %v\n", a.config.BackupEnabled)
	fmt.Printf("  Backup Path: %s\n", a.config.BackupPath)
	fmt.Printf("  Cache Retention: %d days\n", a.config.CacheRetentionDays)
	fmt.Printf("  Log Retention: %d days\n", a.config.LogRetentionDays)
	fmt.Printf("  Verbose Mode: %v\n", a.config.VerboseMode)

	fmt.Println("\nConfiguration Options:")
	fmt.Println("  1. Toggle Dry Run Mode")
	fmt.Println("  2. Toggle Safe Mode")
	fmt.Println("  3. Toggle Backup")
	fmt.Println("  4. Set Cache Retention")
	fmt.Println("  5. Set Log Retention")
	fmt.Println("  6. Toggle Verbose Mode")
	fmt.Println("  7. Export Configuration")
	fmt.Println("  0. Back")

	fmt.Print("\nSelect option: ")
	reader := bufio.NewReader(os.Stdin)
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	switch choice {
	case "1":
		a.config.DryRun = !a.config.DryRun
		successColor.Printf("Dry Run Mode: %v\n", a.config.DryRun)
	case "2":
		a.config.SafeMode = !a.config.SafeMode
		successColor.Printf("Safe Mode: %v\n", a.config.SafeMode)
	case "3":
		a.config.BackupEnabled = !a.config.BackupEnabled
		successColor.Printf("Backup Enabled: %v\n", a.config.BackupEnabled)
	case "4":
		fmt.Print("Enter cache retention days (default 30): ")
		input, _ := reader.ReadString('\n')
		if days := parseInt(strings.TrimSpace(input)); days > 0 {
			a.config.CacheRetentionDays = days
			successColor.Printf("Cache retention set to %d days\n", days)
		}
	case "5":
		fmt.Print("Enter log retention days (default 7): ")
		input, _ := reader.ReadString('\n')
		if days := parseInt(strings.TrimSpace(input)); days > 0 {
			a.config.LogRetentionDays = days
			successColor.Printf("Log retention set to %d days\n", days)
		}
	case "6":
		a.config.VerboseMode = !a.config.VerboseMode
		successColor.Printf("Verbose Mode: %v\n", a.config.VerboseMode)
	case "7":
		a.exportConfig()
	case "0":
		return
	}

	time.Sleep(2 * time.Second)
	a.configManager()
}

func (a *ArchMaintenance) exportConfig() {
	homeDir, _ := os.UserHomeDir()
	configDir := filepath.Join(homeDir, ".config/archmaint")
	os.MkdirAll(configDir, 0755)

	configFile := filepath.Join(configDir, "config.conf")
	content := fmt.Sprintf(`# ArchMaint Configuration
# Generated: %s

DRY_RUN=%v
SAFE_MODE=%v
BACKUP_ENABLED=%v
BACKUP_PATH=%s
CACHE_RETENTION_DAYS=%d
LOG_RETENTION_DAYS=%d
VERBOSE_MODE=%v
`,
		time.Now().Format("2006-01-02 15:04:05"),
		a.config.DryRun,
		a.config.SafeMode,
		a.config.BackupEnabled,
		a.config.BackupPath,
		a.config.CacheRetentionDays,
		a.config.LogRetentionDays,
		a.config.VerboseMode,
	)

	if err := os.WriteFile(configFile, []byte(content), 0644); err == nil {
		successColor.Printf("Configuration exported to: %s\n", configFile)
	} else {
		errorColor.Printf("Failed to export configuration: %v\n", err)
	}

	time.Sleep(2 * time.Second)
}

func (a *ArchMaintenance) showHelp() {
	a.showBanner()
	fmt.Println("USAGE:")
	fmt.Println("  archmaint [OPTIONS] [COMMAND]")
	fmt.Println()
	fmt.Println("OPTIONS:")
	fmt.Println("  --dry-run          Show what would be done without making changes")
	fmt.Println("  --safe             Enable safe mode with extra confirmations")
	fmt.Println()
	fmt.Println("COMMANDS:")

	commands := [][]string{
		{"status, s", "Show system status and information"},
		{"update, u", "Update system packages (with backup)"},
		{"clean, c", "Clean system (cache, logs, temp files)"},
		{"orphans, o", "Remove orphaned packages"},
		{"services, sv", "Show system services status"},
		{"logs, l", "Show recent system logs"},
		{"health, h", "Run comprehensive health check"},
		{"maintenance, m", "Run full maintenance routine"},
		{"search, se", "Search for packages"},
		{"backup, b", "Create system backup"},
		{"restore, r", "Restore from backup"},
		{"snapshot, sn", "Create btrfs snapshot"},
		{"config, cfg", "Manage configuration"},
		{"help, --help, -h", "Show this help message"},
		{"version, --version, -v", "Show version information"},
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetBorder(false)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("   ")
	table.SetRowSeparator("")
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	for _, cmd := range commands {
		table.Append(cmd)
	}
	table.Render()

	fmt.Println("\nEXAMPLES:")
	fmt.Println("  archmaint status              # Show system status")
	fmt.Println("  archmaint --dry-run update    # Preview system updates")
	fmt.Println("  archmaint --safe clean        # Clean with extra safety")
	fmt.Println("  archmaint search firefox      # Search for firefox package")
	fmt.Println("  archmaint backup              # Create system backup")
	fmt.Println("  archmaint                     # Interactive mode")

	fmt.Println("\nFEATURES (v1.1):")
	fmt.Println("  - Dry-run mode to preview changes")
	fmt.Println("  - Safe mode with extra confirmations")
	fmt.Println("  - Automatic backups before updates")
	fmt.Println("  - Package search functionality")
	fmt.Println("  - Btrfs snapshot support")
	fmt.Println("  - Configuration management")
	fmt.Println("  - Progress bars for long operations")
	fmt.Println("  - Enhanced health checks")
	fmt.Println("  - Backup and restore system")

	a.waitForContinue()
}

func (a *ArchMaintenance) showVersion() {
	a.showBanner()
	fmt.Printf("Version: %s\n", a.version)
	fmt.Println("Built for Arch Linux")
	fmt.Println()
	fmt.Println("New in v1.1:")
	fmt.Println("  - Dry-run mode")
	fmt.Println("  - Safe mode")
	fmt.Println("  - Backup/Restore system")
	fmt.Println("  - Package search")
	fmt.Println("  - Btrfs snapshots")
	fmt.Println("  - Configuration manager")
	fmt.Println("  - Progress indicators")
	fmt.Println("  - Enhanced health checks")
	fmt.Println()
	fmt.Println("https://github.com/yourusername/archmaint")
}

func (a *ArchMaintenance) runCommand(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		if !a.config.DryRun {
			errorColor.Printf("Error running command: %v\n", err)
		}
	}
}

func (a *ArchMaintenance) runCommandWithProgress(name string, args ...string) {
	if a.config.VerboseMode {
		infoColor.Printf("Running: %s %s\n", name, strings.Join(args, " "))
	}

	cmd := exec.Command(name, args...)

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	cmd.Stdin = os.Stdin

	if err := cmd.Start(); err != nil {
		errorColor.Printf("Error starting command: %v\n", err)
		return
	}

	done := make(chan bool)
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			if a.config.VerboseMode {
				fmt.Println(scanner.Text())
			}
		}
		done <- true
	}()

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			if a.config.VerboseMode {
				errorColor.Println(scanner.Text())
			}
		}
	}()

	if !a.config.VerboseMode {
		go func() {
			for {
				select {
				case <-done:
					return
				case <-time.After(500 * time.Millisecond):
					progressColor.Print(".")
				}
			}
		}()
	}

	<-done
	if err := cmd.Wait(); err != nil {
		errorColor.Printf("\nCommand failed: %v\n", err)
	} else if !a.config.VerboseMode {
		fmt.Println()
	}
}

func (a *ArchMaintenance) confirmAction(message string, dangerous bool) bool {
	if a.config.AutoConfirm && !dangerous {
		return true
	}

	if a.config.SafeMode && dangerous {
		dangerColor.Println("DANGEROUS OPERATION - Extra confirmation required!")
		fmt.Print("Type 'yes' to confirm: ")
		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		return strings.TrimSpace(strings.ToLower(response)) == "yes"
	}

	if dangerous {
		dangerColor.Printf("WARNING %s [y/N]: ", message)
	} else {
		warningColor.Printf("? %s [y/N]: ", message)
	}

	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')
	response = strings.ToLower(strings.TrimSpace(response))

	return response == "y" || response == "yes"
}

func (a *ArchMaintenance) waitForContinue() {
	fmt.Print("\nPress Enter to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	if len(os.Args) == 1 || (len(os.Args) > 1 && strings.HasPrefix(os.Args[1], "--")) {
		a.showMainMenu()
	}
}

func copyFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}
