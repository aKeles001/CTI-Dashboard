#!/bin/bash

# Wails Development Docker Image Build Script

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
IMAGE_NAME="wails-dev"
TAG="latest"
FULL_IMAGE_NAME="${IMAGE_NAME}:${TAG}"

# Functions
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Build the Docker image
build_image() {
    print_info "Building Docker image: $FULL_IMAGE_NAME"
    
    # Check if Docker is running
    if ! docker info > /dev/null 2>&1; then
        print_error "Docker is not running. Please start Docker and try again."
        exit 1
    fi
    
    # Build the image
    docker build -t "$FULL_IMAGE_NAME" .
    
    if [ $? -eq 0 ]; then
        print_success "Docker image built successfully: $FULL_IMAGE_NAME"
    else
        print_error "Failed to build Docker image"
        exit 1
    fi
}

# Test the Docker image
test_image() {
    print_info "Testing Docker image: $FULL_IMAGE_NAME"
    
    # Test if the image exists
    if ! docker image inspect "$FULL_IMAGE_NAME" > /dev/null 2>&1; then
        print_error "Image $FULL_IMAGE_NAME not found. Please build it first."
        exit 1
    fi
    
    # Test Go installation
    print_info "Testing Go installation..."
    GO_VERSION=$(docker run --rm --entrypoint="" "$FULL_IMAGE_NAME" go version)
    print_success "Go version: $GO_VERSION"
    
    # Test Wails CLI
    print_info "Testing Wails CLI..."
    WAILS_VERSION=$(docker run --rm --entrypoint="" "$FULL_IMAGE_NAME" wails version)
    print_success "Wails version: $WAILS_VERSION"
    
    # Test WebKit2GTK
    print_info "Testing WebKit2GTK installation..."
    WEBKIT_VERSION=$(docker run --rm --entrypoint="" "$FULL_IMAGE_NAME" pkg-config --modversion webkit2gtk-4.1 2>/dev/null || echo "Not found")
    if [ "$WEBKIT_VERSION" != "Not found" ]; then
        print_success "WebKit2GTK version: $WEBKIT_VERSION"
    else
        print_warning "WebKit2GTK not found or not properly installed"
    fi
    
    # Test xgo
    print_info "Testing xgo installation..."
    XGO_PATH=$(docker run --rm --entrypoint="" "$FULL_IMAGE_NAME" which xgo 2>/dev/null)
    if [ -n "$XGO_PATH" ]; then
        print_success "xgo is available at $XGO_PATH"
    else
        print_warning "xgo not found or not properly installed"
    fi
    
    # Test Node.js
    print_info "Testing Node.js installation..."
    NODE_VERSION=$(docker run --rm --entrypoint="" "$FULL_IMAGE_NAME" node --version 2>/dev/null)
    if [ -n "$NODE_VERSION" ]; then
        print_success "Node.js version: $NODE_VERSION"
    else
        print_warning "Node.js not found or not properly installed"
    fi
    
    # Test npm
    print_info "Testing npm installation..."
    NPM_VERSION=$(docker run --rm --entrypoint="" "$FULL_IMAGE_NAME" npm --version 2>/dev/null)
    if [ -n "$NPM_VERSION" ]; then
        print_success "npm version: $NPM_VERSION"
    else
        print_warning "npm not found or not properly installed"
    fi
    
    print_success "Image test completed successfully"
}

# Run the container interactively
run_container() {
    print_info "Running container interactively..."
    
    # Check if image exists
    if ! docker image inspect "$FULL_IMAGE_NAME" > /dev/null 2>&1; then
        print_error "Image $FULL_IMAGE_NAME not found. Please build it first."
        exit 1
    fi
    
    # Run the container
    docker run -it --rm \
        -v "$(pwd):/app" \
        -v wails-go-cache:/home/wailsdev/go \
        -e DISPLAY="${DISPLAY:-:0}" \
        -v /tmp/.X11-unix:/tmp/.X11-unix:rw \
        --network host \
        "$FULL_IMAGE_NAME"
}

# Setup Distrobox
setup_distrobox() {
    print_info "Setting up Distrobox container..."
    
    # Check if image exists
    if ! docker image inspect "$FULL_IMAGE_NAME" > /dev/null 2>&1; then
        print_error "Image $FULL_IMAGE_NAME not found. Please build it first."
        exit 1
    fi
    
    # Check if distrobox is installed
    if ! command -v distrobox > /dev/null 2>&1; then
        print_error "Distrobox is not installed. Please install it first."
        print_info "Installation: https://github.com/89luca89/distrobox#installation"
        exit 1
    fi
    
    # Create Distrobox container
    print_info "Creating Distrobox container..."
    distrobox create wails-dev --image "$FULL_IMAGE_NAME"
    
    if [ $? -eq 0 ]; then
        print_success "Distrobox container created successfully"
        print_info "You can now enter the container with: distrobox enter wails-dev"
    else
        print_error "Failed to create Distrobox container"
        exit 1
    fi
}

# Clean up
cleanup() {
    print_info "Cleaning up Docker resources..."
    
    # Remove the image
    if docker image inspect "$FULL_IMAGE_NAME" > /dev/null 2>&1; then
        docker rmi "$FULL_IMAGE_NAME"
        print_success "Removed image: $FULL_IMAGE_NAME"
    fi
    
    # Remove unused volumes
    docker volume prune -f
    print_success "Cleaned up unused volumes"
}

# Show usage
show_usage() {
    echo "Usage: $0 [COMMAND]"
    echo ""
    echo "Commands:"
    echo "  build     Build the Docker image"
    echo "  test      Test the Docker image"
    echo "  run       Run the container interactively"
    echo "  distrobox Create and setup Distrobox container"
    echo "  clean     Clean up Docker resources"
    echo "  all       Build, test, and run the container"
    echo "  help      Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 build"
    echo "  $0 test"
    echo "  $0 run"
    echo "  $0 distrobox"
    echo "  $0 all"
}

# Main script logic
case "${1:-help}" in
    build)
        build_image
        ;;
    test)
        test_image
        ;;
    run)
        run_container
        ;;
    distrobox)
        setup_distrobox
        ;;
    clean)
        cleanup
        ;;
    all)
        build_image
        test_image
        run_container
        ;;
    help|--help|-h)
        show_usage
        ;;
    *)
        print_error "Unknown command: $1"
        show_usage
        exit 1
        ;;
esac 