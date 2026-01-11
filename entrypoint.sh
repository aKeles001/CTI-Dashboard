#!/bin/bash

# Entrypoint script for Wails development container

# Set up environment for Ubuntu 24.04 WebKit compatibility
export WEBKIT_TAGS="-tags webkit2_41"

# Ensure PATH includes Go and user binaries
export PATH="/home/wailsdev/go/bin:/usr/local/go/bin:$PATH"

# Ensure Node.js and npm are available
export PATH="/usr/bin:$PATH"

# Handle Distrobox environment
if [ -n "$DISTROBOX_HOST_HOME" ]; then
    echo "Running in Distrobox environment"
    # Ensure proper permissions for Distrobox user
    if [ "$(id -u)" = "1000" ]; then
        echo "Running as user 1000 (Distrobox user)"
        # Ensure Go directories exist and are accessible
        mkdir -p "$HOME/go/bin" "$HOME/go/pkg" "$HOME/go/src"
        # Also ensure wailsdev symlink exists
        if [ ! -L "/home/wailsdev" ]; then
            ln -sf "$HOME" /home/wailsdev
        fi
    fi
fi

# Function to run wails commands with proper tags
wails_with_tags() {
    if [[ "$1" == "build" ]] || [[ "$1" == "dev" ]]; then
        # Add webkit2_41 tags for Ubuntu 24.04 compatibility
        wails "$@" $WEBKIT_TAGS
    else
        wails "$@"
    fi
}

# Function to run xgo with proper configuration
xgo_build() {
    local targets=${1:-"linux/amd64,linux/arm64,windows/amd64,windows/386,darwin/amd64,darwin/arm64"}
    local package=${2:-"./"}
    
    echo "Building for targets: $targets"
    xgo -v -targets="$targets" -out="wails-app" $WEBKIT_TAGS "$package"
}

# Export functions for use in container
export -f wails_with_tags
export -f xgo_build

# Print welcome message
echo "=========================================="
echo "Wails Development Environment"
echo "=========================================="
echo "Go version: $(go version)"
echo "Wails version: $(wails version)"
echo "Node.js version: $(node --version)"
echo "npm version: $(npm --version)"
echo "WebKit tags: $WEBKIT_TAGS"
echo ""
echo "Available commands:"
echo "  wails_with_tags build    - Build with webkit2_41 tags"
echo "  wails_with_tags dev      - Run dev server with webkit2_41 tags"
echo "  xgo_build [targets] [pkg] - Cross-compile with xgo"
echo "  wails [command]          - Standard wails commands"
echo ""
echo "Example usage:"
echo "  wails_with_tags build"
echo "  xgo_build 'linux/amd64,windows/amd64'"
echo "=========================================="

# Execute the command passed to the container
exec "$@" 