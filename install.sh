#!/bin/bash

# Function to get latest release version
get_latest_version() {
    curl -sS https://api.github.com/repos/bjarneo/peachy/releases/latest | grep "tag_name" | cut -d '"' -f 4
}

# Function to detect OS and architecture
detect_system() {
    local os
    local arch

    # Detect OS
    case "$(uname -s)" in
        Darwin*)  os="darwin" ;;
        Linux*)   os="linux" ;;
        *)        echo "Unsupported operating system" && exit 1 ;;
    esac

    # Detect architecture
    case "$(uname -m)" in
        x86_64)  arch="amd64" ;;
        arm64)   arch="arm64" ;;
        aarch64) arch="arm64" ;;
        *)       echo "Unsupported architecture" && exit 1 ;;
    esac

    echo "${os}-${arch}"
}

# Main installation process
main() {
    echo "Detecting system..."
    local system=$(detect_system)
    local version=$(get_latest_version)
    
    echo "Latest version: ${version}"
    echo "System detected: ${system}"
    
    # Check if peachy already exists
    if command -v peachy >/dev/null 2>&1; then
        echo "peachy is already installed at $(which peachy)"
        # Only prompt if running interactively (not piped)
        if [ -t 0 ]; then
            read -p "Do you want to override the existing installation? (y/N) " response
            case "$response" in
                [yY][eE][sS]|[yY])
                    echo "Proceeding with installation..."
                    ;;
                *)
                    echo "Installation cancelled"
                    exit 0
                    ;;
            esac
        else
            echo "Upgrading..."
        fi
    fi
    
    local binary_name="peachy-${system}"
    local download_url="https://github.com/bjarneo/peachy/releases/download/${version}/${binary_name}"
    
    echo "Downloading peachy..."
    if ! curl -sSL -o peachy "${download_url}"; then
        echo "Failed to download peachy"
        exit 1
    fi
    
    echo "Making binary executable..."
    chmod +x peachy
    
    echo "Moving to /usr/local/bin..."
    if ! sudo mv peachy /usr/local/bin/peachy; then
        echo "Failed to move peachy to /usr/local/bin"
        echo "Please run with sudo or check permissions"
        exit 1
    fi
    
    echo "peachy has been successfully installed!"
    echo "You can now use it by running: peachy --help"
}

main
