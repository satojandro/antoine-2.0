const InstallScriptContent = `#!/bin/bash

# Antoine CLI Installation Script
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
GOLD='\033[1;33m'
NC='\033[0m' # No Color

# Antoine ASCII Art
print_logo() {
    echo -e "${GOLD}"
    cat << "EOF"
███████╗ ███╗  ██╗ ████████╗  ██████╗  ██╗ ███╗  ██╗ ███████╗
██╔══██║ ████╗ ██║ ╚══██╔══╝ ██╔═══██╗ ██║ ████╗ ██║ ██╔════╝
███████║ ██╔██╗██║    ██║    ██║   ██║ ██║ ██╔██╗██║ █████╗
██╔══██║ ██║╚████║    ██║    ██║   ██║ ██║ ██║╚████║ ██╔══╝
██║  ██║ ██║ ╚███║    ██║    ╚██████╔╝ ██║ ██║ ╚███║ ███████╗
╚═╝  ╚═╝ ╚═╝  ╚══╝    ╚═╝     ╚═════╝  ╚═╝ ╚═╝  ╚══╝ ╚══════╝

           Your Ultimate Hackathon Mentor 🤖✨
EOF
    echo -e "${NC}"
}

# Utility functions
log() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
    exit 1
}

# Detect OS and architecture
detect_platform() {
    local os
    local arch

    case "$OSTYPE" in
        linux*)   os="linux" ;;
        darwin*)  os="darwin" ;;
        msys*|cygwin*) os="windows" ;;
        *)        error "Unsupported operating system: $OSTYPE" ;;
    esac

    case "$(uname -m)" in
        x86_64|amd64) arch="amd64" ;;
        arm64|aarch64) arch="arm64" ;;
        *)            arch="amd64" ;;
    esac

    PLATFORM="${os}-${arch}"
}

# Check prerequisites
check_prerequisites() {
    log "Checking prerequisites..."

    # Check if curl or wget is available
    if ! command -v curl >/dev/null 2>&1 && ! command -v wget >/dev/null 2>&1; then
        error "curl or wget is required but not installed."
    fi

    # Check if we have write permissions
    if [ ! -w "/usr/local/bin" ] && [ ! -w "$HOME/.local/bin" ]; then
        warning "No write permissions to /usr/local/bin or ~/.local/bin"
        warning "You may need to run with sudo or choose a different install location"
    fi
}

# Download binary
download_binary() {
    local version="${1:-latest}"
    local base_url="https://github.com/username/antoine-cli/releases"
    local download_url

    if [ "$version" = "latest" ]; then
        download_url="${base_url}/latest/download/antoine-${PLATFORM}"
    else
        download_url="${base_url}/download/${version}/antoine-${PLATFORM}"
    fi

    if [ "$PLATFORM" = "windows-amd64" ]; then
        download_url="${download_url}.exe"
    fi

    log "Downloading Antoine CLI from: $download_url"

    local temp_file="/tmp/antoine-${PLATFORM}"

    if command -v curl >/dev/null 2>&1; then
        curl -fsSL "$download_url" -o "$temp_file"
    elif command -v wget >/dev/null 2>&1; then
        wget -q "$download_url" -O "$temp_file"
    else
        error "Neither curl nor wget found"
    fi

    if [ ! -f "$temp_file" ]; then
        error "Failed to download Antoine CLI"
    fi

    chmod +x "$temp_file"
    BINARY_PATH="$temp_file"
}

# Install binary
install_binary() {
    local install_dir

    # Determine install directory
    if [ -w "/usr/local/bin" ]; then
        install_dir="/usr/local/bin"
    elif [ -w "$HOME/.local/bin" ]; then
        install_dir="$HOME/.local/bin"
        mkdir -p "$install_dir"
    else
        install_dir="$HOME/bin"
        mkdir -p "$install_dir"
        warning "Installing to $install_dir - make sure it's in your PATH"
    fi

    local binary_name="antoine"
    if [ "$PLATFORM" = "windows-amd64" ]; then
        binary_name="antoine.exe"
    fi

    log "Installing Antoine CLI to ${install_dir}/${binary_name}"

    cp "$BINARY_PATH" "${install_dir}/${binary_name}"

    # Verify installation
    if command -v antoine >/dev/null 2>&1; then
        success "Antoine CLI installed successfully!"
        antoine version
    else
        warning "Antoine CLI installed but not found in PATH"
        warning "You may need to add ${install_dir} to your PATH"
        echo "Add this line to your shell profile (.bashrc, .zshrc, etc.):"
        echo "export PATH=\"${install_dir}:\$PATH\""
    fi
}

# Setup initial configuration
setup_config() {
    log "Setting up initial configuration..."

    # Create config directory
    local config_dir="$HOME/.config/antoine"
    mkdir -p "$config_dir"

    # Create default config if it doesn't exist
    if [ ! -f "$HOME/.antoine.yaml" ]; then
        cat > "$HOME/.antoine.yaml" << EOF
api:
  base_url: "https://api.antoine.ai"
  timeout: "30s"

mcp:
  servers:
    exa: "mcp://localhost:8001"
    github: "mcp://localhost:8002"
    deepwiki: "mcp://localhost:8003"

ui:
  theme: "dark"
  animations: true
  ascii_art: true

cache:
  enabled: true
  ttl: "30m"

analytics:
  enabled: true
  anonymous: true
EOF
        success "Created default configuration at ~/.antoine.yaml"
    fi
}

# Main installation function
main() {
    print_logo

    log "Starting Antoine CLI installation..."

    detect_platform
    log "Detected platform: $PLATFORM"

    check_prerequisites
    download_binary "$1"
    install_binary
    setup_config

    success "Installation complete! 🎉"
    echo
    echo -e "${GOLD}Quick Start:${NC}"
    echo "  antoine              # Show dashboard"
    echo "  antoine search hackathons --tech AI"
    echo "  antoine mentor start"
    echo
    echo -e "${BLUE}Documentation:${NC} https://docs.antoine.ai"
    echo -e "${BLUE}Community:${NC} https://discord.gg/antoine"
    echo
    echo "Happy hacking! 🚀"
}

# Run main function with all arguments
main "$@"`
