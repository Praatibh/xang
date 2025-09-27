#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Output function with colors
output() {
    local message=$1
    local type=${2:-"info"}
    
    case $type in
        "error")
            echo -e "${RED}[ERROR]${NC} $message"
            ;;
        "success")
            echo -e "${GREEN}[SUCCESS]${NC} $message"
            ;;
        "warning")
            echo -e "${YELLOW}[WARNING]${NC} $message"
            ;;
        "info")
            echo -e "${BLUE}[INFO]${NC} $message"
            ;;
        *)
            echo "$message"
            ;;
    esac
}

# Configuration
BINNAME="${BINNAME:-xang}"
BINDIR="${BINDIR:-/usr/local/bin}"
CONFIG_FILE="${HOME}/.config/${BINNAME}.json"

output "Uninstalling xang - AI Terminal Assistant..." "info"
echo

# Check if binary exists
if [ -f "$BINDIR/$BINNAME" ]; then
    output "Removing binary from $BINDIR/$BINNAME" "info"
    sudo rm "$BINDIR/$BINNAME"
    
    if [ $? -eq 0 ]; then
        output "Binary removed successfully" "success"
    else
        output "Failed to remove binary. You may need to run with sudo." "error"
        exit 1
    fi
else
    output "Binary not found at $BINDIR/$BINNAME" "warning"
fi

# Check if config file exists
if [ -f "$CONFIG_FILE" ]; then
    output "Removing configuration file from $CONFIG_FILE" "info"
    rm "$CONFIG_FILE"
    
    if [ $? -eq 0 ]; then
        output "Configuration file removed successfully" "success"
    else
        output "Failed to remove configuration file" "error"
    fi
else
    output "Configuration file not found at $CONFIG_FILE" "warning"
fi

# Check for any remaining config directory
CONFIG_DIR="${HOME}/.config"
if [ -d "$CONFIG_DIR" ] && [ -z "$(ls -A $CONFIG_DIR 2>/dev/null | grep $BINNAME)" ]; then
    output "No additional xang configuration files found" "info"
fi

# Check if command is still accessible
if command -v $BINNAME >/dev/null 2>&1; then
    output "Warning: $BINNAME command is still accessible. It may be installed elsewhere." "warning"
    output "Run 'which $BINNAME' to locate other installations" "info"
else
    output "Command $BINNAME is no longer accessible" "success"
fi

echo
output "Uninstallation of xang complete!" "success"
echo
output "Thank you for using xang!" "info"
output "If you encountered any issues, please report them at:" "info"
output "https://github.com/Praatibh/xang/issues" "info"
