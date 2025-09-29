#!/bin/sh

# tscgit installer for Linux/macOS
# Usage: curl -fsSL https://raw.githubusercontent.com/rohit746/tscgit/main/install.sh | sh

set -e

# Configuration
REPO="rohit746/tscgit"
INSTALL_DIR="${INSTALL_DIR:-$HOME/.local/bin}"
BINARY_NAME="tscgit"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

log_info() {
    printf "${GREEN}%s${NC}\n" "$1"
}

log_warn() {
    printf "${YELLOW}%s${NC}\n" "$1"
}

log_error() {
    printf "${RED}%s${NC}\n" "$1" >&2
}

log_cyan() {
    printf "${CYAN}%s${NC}\n" "$1"
}

# Detect OS and architecture
detect_platform() {
    OS="$(uname -s)"
    ARCH="$(uname -m)"
    
    case "$OS" in
        Linux*)
            OS="linux"
            ;;
        Darwin*)
            OS="darwin"
            ;;
        *)
            log_error "Unsupported operating system: $OS"
            exit 1
            ;;
    esac
    
    case "$ARCH" in
        x86_64)
            ARCH="amd64"
            ;;
        aarch64|arm64)
            ARCH="arm64"
            ;;
        armv7l)
            ARCH="arm"
            ;;
        *)
            log_error "Unsupported architecture: $ARCH"
            exit 1
            ;;
    esac
}

# Check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Download file
download() {
    url="$1"
    output="$2"
    
    if command_exists curl; then
        curl -fsSL "$url" -o "$output"
    elif command_exists wget; then
        wget -q "$url" -O "$output"
    else
        log_error "Neither curl nor wget found. Please install one of them."
        exit 1
    fi
}

# Main installation function
install_tscgit() {
    log_info "Installing tscgit..."
    
    detect_platform
    log_cyan "Detected platform: $OS-$ARCH"
    
    # Create install directory
    mkdir -p "$INSTALL_DIR"
    
    # Get latest release info
    log_warn "Fetching latest release info..."
    TEMP_DIR="$(mktemp -d)"
    RELEASE_INFO="$TEMP_DIR/release.json"
    
    download "https://api.github.com/repos/$REPO/releases/latest" "$RELEASE_INFO"
    
    # Extract version and download URL
    VERSION="$(grep '"tag_name"' "$RELEASE_INFO" | sed -E 's/.*"tag_name": "([^"]+)".*/\1/')"
    log_cyan "Latest version: $VERSION"
    
    # Find the appropriate asset
    ASSET_NAME="${BINARY_NAME}_${VERSION#v}_${OS}_${ARCH}.tar.gz"
    DOWNLOAD_URL="https://github.com/$REPO/releases/download/$VERSION/$ASSET_NAME"
    
    log_warn "Downloading $ASSET_NAME..."
    ARCHIVE_PATH="$TEMP_DIR/$ASSET_NAME"
    
    if ! download "$DOWNLOAD_URL" "$ARCHIVE_PATH"; then
        log_error "Failed to download $ASSET_NAME"
        log_error "You can download manually from: https://github.com/$REPO/releases"
        exit 1
    fi
    
    # Extract and install
    log_warn "Extracting..."
    cd "$TEMP_DIR"
    
    if command_exists tar; then
        tar -xzf "$ASSET_NAME"
    else
        log_error "tar command not found. Please install tar."
        exit 1
    fi
    
    # Move binary to install directory
    if [ -f "$BINARY_NAME" ]; then
        chmod +x "$BINARY_NAME"
        mv "$BINARY_NAME" "$INSTALL_DIR/"
        BINARY_PATH="$INSTALL_DIR/$BINARY_NAME"
    else
        log_error "Binary not found in archive"
        exit 1
    fi
    
    # Clean up
    rm -rf "$TEMP_DIR"
    
    log_info "✓ tscgit installed to $BINARY_PATH"
    
    # Check if install directory is in PATH
    case ":$PATH:" in
        *":$INSTALL_DIR:"*)
            ;;
        *)
            log_warn "⚠ $INSTALL_DIR is not in your PATH"
            log_cyan "Add the following line to your shell profile (.bashrc, .zshrc, etc.):"
            printf "  export PATH=\"\$PATH:%s\"\n" "$INSTALL_DIR"
            log_cyan "Or run tscgit with the full path: $BINARY_PATH"
            ;;
    esac
    
    # Test installation
    log_warn "Testing installation..."
    if "$BINARY_PATH" version >/dev/null 2>&1; then
        log_info "✓ Installation successful!"
        "$BINARY_PATH" version
        echo
        log_cyan "To get started, run:"
        printf "  %s lessons\n" "$BINARY_NAME"
    else
        log_error "Installation completed but binary test failed"
        exit 1
    fi
}

# Handle errors
trap 'log_error "Installation failed"; exit 1' ERR

# Run installation
install_tscgit