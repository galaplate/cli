#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# GitHub repository details
REPO_OWNER="galaplate"
REPO_NAME="galaplate"
BINARY_NAME="galaplate"

# Default installation directory
INSTALL_DIR="/usr/local/bin"

# Function to print colored output
print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Function to detect OS and architecture
detect_platform() {
    local os
    local arch
    
    # Detect OS
    case "$(uname -s)" in
        Linux*)     os="linux";;
        Darwin*)    os="darwin";;
        CYGWIN*|MINGW*|MSYS*) os="windows";;
        *)          os="unknown";;
    esac
    
    # Detect architecture
    case "$(uname -m)" in
        x86_64|amd64) arch="amd64";;
        arm64|aarch64) arch="arm64";;
        armv7l) arch="arm";;
        i386|i686) arch="386";;
        *) arch="unknown";;
    esac
    
    if [ "$os" = "unknown" ] || [ "$arch" = "unknown" ]; then
        print_error "Unsupported platform: $(uname -s) $(uname -m)"
        exit 1
    fi
    
    echo "${os}_${arch}"
}

# Function to get the latest release version
get_latest_version() {
    local api_url="https://api.github.com/repos/${REPO_OWNER}/${REPO_NAME}/releases/latest"
    
    # Try to get the latest release version
    if command -v curl >/dev/null 2>&1; then
        curl -s "$api_url" | grep '"tag_name":' | sed 's/.*"tag_name": *"\([^"]*\)".*/\1/' | head -1
    elif command -v wget >/dev/null 2>&1; then
        wget -qO- "$api_url" | grep '"tag_name":' | sed 's/.*"tag_name": *"\([^"]*\)".*/\1/' | head -1
    else
        print_error "Neither curl nor wget is available. Please install one of them."
        exit 1
    fi
}

# Function to download and install the binary
install_galaplate() {
    local version="$1"
    local platform="$2"
    local install_dir="$3"
    
    # Remove 'v' prefix if present
    version=$(echo "$version" | sed 's/^v//')
    
    # Construct download URL
    local binary_name="${BINARY_NAME}"
    if [[ "$platform" == *"windows"* ]]; then
        binary_name="${BINARY_NAME}.exe"
    fi
    
    local download_url="https://github.com/${REPO_OWNER}/${REPO_NAME}/releases/download/v${version}/${BINARY_NAME}_${platform}"
    
    print_info "Downloading Galaplate CLI v${version} for ${platform}..."
    print_info "Download URL: ${download_url}"
    
    # Create temporary file
    local temp_file=$(mktemp)
    
    # Download the binary
    if command -v curl >/dev/null 2>&1; then
        if ! curl -L -o "$temp_file" "$download_url"; then
            print_error "Failed to download from $download_url"
            rm -f "$temp_file"
            exit 1
        fi
    elif command -v wget >/dev/null 2>&1; then
        if ! wget -O "$temp_file" "$download_url"; then
            print_error "Failed to download from $download_url"
            rm -f "$temp_file"
            exit 1
        fi
    else
        print_error "Neither curl nor wget is available"
        rm -f "$temp_file"
        exit 1
    fi
    
    # Check if download was successful
    if [ ! -s "$temp_file" ]; then
        print_error "Downloaded file is empty"
        rm -f "$temp_file"
        exit 1
    fi
    
    # Create install directory if it doesn't exist
    if [ ! -d "$install_dir" ]; then
        print_info "Creating directory: $install_dir"
        if ! mkdir -p "$install_dir"; then
            print_error "Failed to create directory: $install_dir"
            print_info "You may need to run this script with sudo"
            rm -f "$temp_file"
            exit 1
        fi
    fi
    
    # Install the binary
    local target_path="$install_dir/$binary_name"
    
    print_info "Installing to: $target_path"
    
    if ! mv "$temp_file" "$target_path"; then
        print_error "Failed to install binary to $target_path"
        print_info "You may need to run this script with sudo"
        rm -f "$temp_file"
        exit 1
    fi
    
    # Make executable
    chmod +x "$target_path"
    
    print_success "Galaplate CLI v${version} installed successfully!"
    
    # Verify installation
    if command -v "$BINARY_NAME" >/dev/null 2>&1; then
        print_success "Installation verified: $(which $BINARY_NAME)"
        print_info "Run 'galaplate help' to get started"
    else
        print_warning "Binary installed but not found in PATH"
        print_info "You may need to add $install_dir to your PATH"
        print_info "Or run the binary directly: $target_path"
    fi
}

# Function to show help
show_help() {
    echo "Galaplate CLI Installer"
    echo ""
    echo "Usage: $0 [options]"
    echo ""
    echo "Options:"
    echo "  -d, --dir <dir>     Installation directory (default: /usr/local/bin)"
    echo "  -v, --version <ver> Install specific version (default: latest)"
    echo "  -h, --help          Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0                           # Install latest version to /usr/local/bin"
    echo "  $0 -d ~/.local/bin           # Install to custom directory"
    echo "  $0 -v 1.0.0                 # Install specific version"
    echo ""
}

# Main installation function
main() {
    local version=""
    local install_dir="$INSTALL_DIR"
    
    # Parse command line arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            -d|--dir)
                install_dir="$2"
                shift 2
                ;;
            -v|--version)
                version="$2"
                shift 2
                ;;
            -h|--help)
                show_help
                exit 0
                ;;
            *)
                print_error "Unknown option: $1"
                show_help
                exit 1
                ;;
        esac
    done
    
    print_info "🚀 Galaplate CLI Installer"
    echo ""
    
    # Detect platform
    local platform
    platform=$(detect_platform)
    print_info "Detected platform: $platform"
    
    # Get version
    if [ -z "$version" ]; then
        print_info "Getting latest release version..."
        version=$(get_latest_version)
        if [ -z "$version" ]; then
            print_error "Failed to get latest version"
            exit 1
        fi
    fi
    
    print_info "Version to install: $version"
    print_info "Installation directory: $install_dir"
    echo ""
    
    # Install
    install_galaplate "$version" "$platform" "$install_dir"
    
    echo ""
    print_success "🎉 Installation complete!"
    print_info "Create your first project with: galaplate new my-project"
}

# Run main function if script is executed directly
if [ "${BASH_SOURCE[0]}" == "${0}" ]; then
    main "$@"
fi