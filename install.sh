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

# Project configuration
REPOOWNER="Praatibh"
REPONAME="xang"
BINNAME="${BINNAME:-xang}"
BINDIR="${BINDIR:-/usr/local/bin}"

output "🚀 Installing xang - AI Terminal Assistant with Gemini API" "info"
echo

# Get latest release tag
output "Fetching latest release information..." "info"
RELEASETAG=$(curl -s "https://api.github.com/repos/$REPOOWNER/$REPONAME/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$RELEASETAG" ]; then
    output "Failed to fetch latest release tag. Installing from source..." "warning"
    
    # Fallback: Clone and build from source
    output "Cloning repository..." "info"
    git clone "https://github.com/$REPOOWNER/$REPONAME.git" /tmp/xang-build
    cd /tmp/xang-build
    
    output "Building xang..." "info"
    go build -o $BINNAME .
    
    if [ $? -eq 0 ]; then
        chmod +x $BINNAME
        sudo mv $BINNAME $BINDIR/$BINNAME
        cd /tmp && rm -rf /tmp/xang-build
        output "Installation from source complete!" "success"
    else
        output "Build failed. Please check Go installation and dependencies." "error"
        exit 1
    fi
else
    # Detect OS
    KERNEL=$(uname -s 2>/dev/null || /usr/bin/uname -s)
    case ${KERNEL} in
        "Linux"|"linux")
            KERNEL="linux"
            ;;
        "Darwin"|"darwin")
            KERNEL="darwin"
            ;;
        *)
            output "OS '${KERNEL}' not supported" "error"
            exit 1
            ;;
    esac

    # Detect architecture
    MACHINE=$(uname -m 2>/dev/null || /usr/bin/uname -m)
    case ${MACHINE} in
        arm|armv7*)
            MACHINE="arm"
            ;;
        aarch64*|armv8*|arm64)
            MACHINE="arm64"
            ;;
        i[36]86)
            MACHINE="386"
            if [ "darwin" = "${KERNEL}" ]; then
                output "Your architecture (${MACHINE}) is not supported on macOS" "error"
                exit 1
            fi
            ;;
        x86_64)
            MACHINE="amd64"
            ;;
        *)
            output "Your architecture (${MACHINE}) is not currently supported" "error"
            exit 1
            ;;
    esac

    # Download and install binary
    FILENAME="xang_${RELEASETAG}_${KERNEL}_${MACHINE}.tar.gz"
    URL="https://github.com/$REPOOWNER/$REPONAME/releases/download/${RELEASETAG}/${FILENAME}"
    
    output "Downloading xang version $RELEASETAG for $KERNEL/$MACHINE..." "info"
    output "URL: $URL" "info"
    echo
    
    curl -q --fail --location --progress-bar --output "$FILENAME" "$URL"
    
    if [ $? -ne 0 ]; then
        output "Download failed. Trying to install from source..." "warning"
        
        # Fallback to source installation
        output "Cloning repository..." "info"
        git clone "https://github.com/$REPOOWNER/$REPONAME.git" /tmp/xang-build
        cd /tmp/xang-build
        
        output "Building xang..." "info"
        go build -o $BINNAME .
        
        if [ $? -eq 0 ]; then
            chmod +x $BINNAME
            sudo mv $BINNAME $BINDIR/$BINNAME
            cd /tmp && rm -rf /tmp/xang-build
            output "Installation from source complete!" "success"
        else
            output "Build failed. Please check Go installation and dependencies." "error"
            exit 1
        fi
    else
        # Extract and install binary
        output "Extracting archive..." "info"
        tar xzf "$FILENAME"
        
        if [ -f "$BINNAME" ]; then
            chmod +x $BINNAME
            sudo mv $BINNAME $BINDIR/$BINNAME
            rm "$FILENAME"
            output "Installation of xang version $RELEASETAG complete!" "success"
        else
            output "Binary not found in archive. Installation failed." "error"
            exit 1
        fi
    fi
fi

echo
output "🎉 xang has been installed successfully!" "success"
output "Location: $BINDIR/$BINNAME" "info"
echo
output "To get started:" "info"
echo "  1. Run 'xang' to start the REPL mode"
echo "  2. You'll be prompted to enter your Gemini API key on first run"
echo "  3. Get your free Gemini API key at: https://aistudio.google.com/app/apikey"
echo
output "Usage examples:" "info"
echo "  xang                    # Start interactive REPL mode"
echo "  xang 'list files'       # Execute single command"
echo "  echo 'data' | xang      # Process piped input"
echo
output "Need help? Run 'xang --help' or press Ctrl+H in REPL mode" "info"
echo
output "Happy coding! 🚀" "success"
