package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

// ArchMaintenance represents the main application
type ArchMaintenance struct {
	version string
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
}

// Colors for beautiful output
var (
	headerColor  = color.New(color.FgCyan, color.Bold)
	successColor = color.New(color.FgGreen, color.Bold)
	warningColor = color.New(color.FgYellow, color.Bold)
	errorColor   = color.New(color.FgRed, color.Bold)
	infoColor    = color.New(color.FgBlue)
	dangerColor  = color.New(color.FgRed, color.Bold, color.BgYellow)
)

func main() {
	app := &ArchMaintenance{version: "1.0.0"}

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
		case "help", "--help", "-h":
			app.showHelp()
		case "version", "--version", "-v":
			app.showVersion()
		default:
			app.showHelp()
		}
	} else {
		app.showMainMenu()
	}
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
	infoColor.Printf("                        ArchLinux Maintenance Tool (DaEpicR) v%s\n", a.version)
	fmt.Println()
}

func (a *ArchMaintenance) showMainMenu() {
	a.showBanner()

	options := [][]string{
		{"1", "System Status", "Show system information and status"},
		{"2", "System Update", "Update system packages"},
		{"3", "System Clean", "Clean package cache and temporary files"},
		{"4", "Remove Orphans", "Remove orphaned packages"},
		{"5", "System Services", "View system services status"},
		{"6", "System Logs", "View recent system logs"},
		{"7", "Health Check", "Comprehensive system health check"},
		{"8", "Full Maintenance", "Run complete maintenance routine"},
		{"9", "Help", "Show help information"},
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
		a.showHelp()
	case "0":
		successColor.Println("Goodbye! Keep your Arch system running smoothly! 🐧")
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

	// Package information
	fmt.Println()
	a.showPackageInfo()

	a.waitForContinue()
}

func (a *ArchMaintenance) getSystemInfo() SystemInfo {
	info := SystemInfo{}

	// Kernel version
	if output, err := exec.Command("uname", "-r").Output(); err == nil {
		info.Kernel = strings.TrimSpace(string(output))
	}

	// Uptime
	if output, err := exec.Command("uptime", "-p").Output(); err == nil {
		info.Uptime = strings.TrimSpace(string(output))
	}

	// Load average
	if output, err := exec.Command("cat", "/proc/loadavg").Output(); err == nil {
		fields := strings.Fields(string(output))
		if len(fields) >= 3 {
			info.LoadAvg = fmt.Sprintf("%s %s %s", fields[0], fields[1], fields[2])
		}
	}

	// Memory usage
	if output, err := exec.Command("free", "-h").Output(); err == nil {
		lines := strings.Split(string(output), "\n")
		if len(lines) > 1 {
			fields := strings.Fields(lines[1])
			if len(fields) >= 3 {
				info.MemoryUsage = fmt.Sprintf("%s / %s", fields[2], fields[1])
			}
		}
	}

	// Disk usage
	if output, err := exec.Command("df", "-h", "/").Output(); err == nil {
		lines := strings.Split(string(output), "\n")
		if len(lines) > 1 {
			fields := strings.Fields(lines[1])
			if len(fields) >= 5 {
				info.DiskUsage = fmt.Sprintf("%s / %s (%s)", fields[2], fields[1], fields[4])
			}
		}
	}

	return info
}

func (a *ArchMaintenance) showPackageInfo() {
	infoColor.Println("Package Information:")

	// Installed packages
	if output, err := exec.Command("pacman", "-Q").Output(); err == nil {
		count := len(strings.Split(strings.TrimSpace(string(output)), "\n"))
		fmt.Printf("  Installed packages: %d\n", count)
	}

	// Explicitly installed packages
	if output, err := exec.Command("pacman", "-Qe").Output(); err == nil {
		count := len(strings.Split(strings.TrimSpace(string(output)), "\n"))
		fmt.Printf("  Explicitly installed: %d\n", count)
	}

	// Orphaned packages
	if output, err := exec.Command("pacman", "-Qtdq").Output(); err == nil {
		orphans := strings.TrimSpace(string(output))
		if orphans != "" {
			count := len(strings.Split(orphans, "\n"))
			warningColor.Printf("  Orphaned packages: %d\n", count)
		} else {
			fmt.Printf("  Orphaned packages: 0\n")
		}
	}
}

func (a *ArchMaintenance) systemUpdate() {
	headerColor.Println("\n=== SYSTEM UPDATE ===")

	if !a.confirmAction("This will update your system. Continue?", false) {
		return
	}

	infoColor.Println("Syncing package databases...")
	a.runCommand("sudo", "pacman", "-Sy")

	infoColor.Println("\nChecking for updates...")
	cmd := exec.Command("pacman", "-Qu")
	output, err := cmd.Output()

	if err != nil || strings.TrimSpace(string(output)) == "" {
		successColor.Println("System is up to date!")
		a.waitForContinue()
		return
	}

	fmt.Println("\nAvailable updates:")
	fmt.Println(string(output))

	if a.confirmAction("Proceed with system update?", false) {
		infoColor.Println("Updating system...")
		a.runCommand("sudo", "pacman", "-Su")
		successColor.Println("System update completed!")
	}

	a.waitForContinue()
}

func (a *ArchMaintenance) systemClean() {
	headerColor.Println("\n=== SYSTEM CLEAN ===")

	tasks := []Task{
		{
			Name:        "Package Cache",
			Description: "Clean pacman cache (keep last 3 versions)",
			Command:     []string{"sudo", "paccache", "-r"},
			Dangerous:   false,
			Frequency:   "Weekly",
		},
		{
			Name:        "Orphaned Packages",
			Description: "Remove packages no longer needed",
			Command:     []string{"sudo", "pacman", "-Rns", "$(pacman -Qtdq)"},
			Dangerous:   true,
			Frequency:   "Weekly",
		},
		{
			Name:        "System Logs",
			Description: "Clean old journal logs (keep 1 week)",
			Command:     []string{"sudo", "journalctl", "--vacuum-time=1week"},
			Dangerous:   false,
			Frequency:   "Weekly",
		},
		{
			Name:        "Temporary Files",
			Description: "Clean /tmp and /var/tmp",
			Command:     []string{"sudo", "find", "/tmp", "/var/tmp", "-type", "f", "-atime", "+7", "-delete"},
			Dangerous:   false,
			Frequency:   "Daily",
		},
	}

	for _, task := range tasks {
		fmt.Printf("\n%s (%s)\n", task.Name, task.Frequency)
		fmt.Printf("Description: %s\n", task.Description)

		if task.Dangerous {
			dangerColor.Printf("[CAUTION] This action can be dangerous!\n")
		}

		if task.Name == "Orphaned Packages" {
			// Special handling for orphaned packages
			a.removeOrphans()
			continue
		}

		if a.confirmAction(fmt.Sprintf("Run %s cleanup?", task.Name), task.Dangerous) {
			a.runCommandSlice(task.Command)
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
	fmt.Println("Orphaned packages:")
	fmt.Println(orphans)

	if a.confirmAction("Remove these orphaned packages?", true) {
		orphanList := strings.Fields(orphans)
		args := append([]string{"-Rns"}, orphanList...)
		a.runCommand("sudo", append([]string{"pacman"}, args...)...)
		successColor.Println("Orphaned packages removed!")
	}
}

func (a *ArchMaintenance) showServices() {
	headerColor.Println("\n=== SYSTEM SERVICES ===")

	infoColor.Println("Failed services:")
	a.runCommand("systemctl", "--failed")

	fmt.Println()
	infoColor.Println("Most recent service status:")
	a.runCommand("systemctl", "status", "--no-pager", "-l")

	a.waitForContinue()
}

func (a *ArchMaintenance) showLogs() {
	headerColor.Println("\n=== SYSTEM LOGS ===")

	infoColor.Println("Recent critical and error logs:")
	a.runCommand("journalctl", "-p", "3", "-x", "--no-pager", "--since", "today")

	fmt.Println()
	infoColor.Println("Boot messages:")
	a.runCommand("journalctl", "-b", "--no-pager", "-n", "20")

	a.waitForContinue()
}

func (a *ArchMaintenance) systemHealthCheck() {
	headerColor.Println("\n=== SYSTEM HEALTH CHECK ===")

	checks := []struct {
		name string
		cmd  []string
		desc string
	}{
		{"Disk Health", []string{"df", "-h"}, "Checking disk usage"},
		{"Memory Usage", []string{"free", "-h"}, "Checking memory usage"},
		{"Failed Services", []string{"systemctl", "--failed", "--no-legend"}, "Checking failed services"},
		{"Package Database", []string{"pacman", "-Dk"}, "Checking package database integrity"},
		{"System Errors", []string{"journalctl", "-p", "3", "-x", "--since", "today", "--no-pager"}, "Checking recent errors"},
	}

	for _, check := range checks {
		fmt.Printf("\n%s: %s\n", check.name, check.desc)
		fmt.Println(strings.Repeat("-", 50))
		a.runCommandSlice(check.cmd)
	}

	successColor.Println("\nHealth check completed!")
	a.waitForContinue()
}

func (a *ArchMaintenance) fullMaintenance() {
	headerColor.Println("\n=== FULL SYSTEM MAINTENANCE ===")

	dangerColor.Println("⚠️  WARNING: This will perform comprehensive system maintenance!")
	fmt.Println("This includes:")
	fmt.Println("- System update")
	fmt.Println("- Package cache cleanup")
	fmt.Println("- Orphaned package removal")
	fmt.Println("- Log cleanup")
	fmt.Println("- Temporary file cleanup")

	if !a.confirmAction("Are you sure you want to continue?", true) {
		return
	}

	steps := []struct {
		name string
		fn   func()
	}{
		{"System Update", a.systemUpdate},
		{"System Clean", a.systemClean},
		{"Health Check", a.systemHealthCheck},
	}

	for i, step := range steps {
		headerColor.Printf("\n[%d/%d] %s\n", i+1, len(steps), step.name)
		step.fn()
	}

	successColor.Println("\n🎉 Full maintenance completed successfully!")
	infoColor.Println("Recommendation: Reboot your system to ensure all changes take effect.")

	if a.confirmAction("Reboot now?", true) {
		a.runCommand("sudo", "reboot")
	}

	a.waitForContinue()
}

func (a *ArchMaintenance) showHelp() {
	a.showBanner()
	fmt.Println("USAGE:")
	fmt.Println("  archmaint [COMMAND]")
	fmt.Println()
	fmt.Println("COMMANDS:")

	commands := [][]string{
		{"status, s", "Show system status and information"},
		{"update, u", "Update system packages"},
		{"clean, c", "Clean system (cache, logs, temp files)"},
		{"orphans, o", "Remove orphaned packages"},
		{"services, sv", "Show system services status"},
		{"logs, l", "Show recent system logs"},
		{"health, h", "Run system health check"},
		{"maintenance, m", "Run full maintenance routine"},
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
	fmt.Println("  archmaint status    # Show system status")
	fmt.Println("  archmaint update    # Update system")
	fmt.Println("  archmaint clean     # Clean system")
	fmt.Println("  archmaint           # Interactive mode")

	a.waitForContinue()
}

func (a *ArchMaintenance) showVersion() {
	a.showBanner()
	fmt.Printf("Version: %s\n", a.version)
	fmt.Println("Built for Arch Linux")
	fmt.Println("https://github.com/yourusername/archmaint")
}

func (a *ArchMaintenance) runCommand(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		errorColor.Printf("Error running command: %v\n", err)
	}
}

func (a *ArchMaintenance) runCommandSlice(cmdSlice []string) {
	if len(cmdSlice) == 0 {
		return
	}

	name := cmdSlice[0]
	args := cmdSlice[1:]
	a.runCommand(name, args...)
}

func (a *ArchMaintenance) confirmAction(message string, dangerous bool) bool {
	if dangerous {
		dangerColor.Printf("⚠️  %s [y/N]: ", message)
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

	if len(os.Args) == 1 { // Interactive mode
		a.showMainMenu()
	}
}
