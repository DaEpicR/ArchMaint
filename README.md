## 🚀 **Usage Examples**

### **Interactive Mode:**

```bash
./archmaint
```

Shows beautiful menu with numbered options

### **Direct Commands:**

```bash
./archmaint status    # Quick system overview
./archmaint update    # Update packages
./archmaint clean     # Clean system
./archmaint maintenance # Full maintenance routine
```

## 🔧 **Installation Steps**

1. **Install Go dependencies:**

```bash
go mod init archmaint
go get github.com/fatih/color@v1.16.0
go get github.com/olekukonko/tablewriter@v0.0.5
```

2. **Build the tool:**

```bash
go build -o archmaint .
# or
make build
```

3. **Install system-wide (optional):**

```bash
make install
```

## 🎯 **What Makes It Special**

### **Comprehensive Coverage**

- Package management (updates, orphans, cache)
- System monitoring (services, logs, health)
- Cleanup operations (temp files, logs, cache)
- System information display

### **Production Ready**

- Proper error handling
- Safe defaults with confirmations
- Modular, maintainable code structure
- Cross-platform Go implementation

### **User-Friendly Design**

- Clear visual feedback
- Intuitive command structure
- Helpful descriptions and warnings
- Both beginner and power-user friendly

The tool covers all the essential daily and weekly maintenance tasks that Arch Linux users need:

- **Package Management**: Updates, orphan removal, cache cleaning
- **System Health**: Service monitoring, log analysis, disk/memory checks
- **Maintenance**: Automated routines, comprehensive cleanup
- **Information**: Real-time system status and statistics
