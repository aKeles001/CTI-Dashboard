# Wails Development Docker Image

A modern Docker image for Wails development based on Go 1.24 and Ubuntu 24.04, with full cross-platform compilation support and WebKit2GTK-4.1 compatibility.

## Features

- **Go 1.24**: Latest stable Go version
- **Ubuntu 24.04**: Modern Linux distribution
- **WebKit2GTK-4.1**: Latest WebKit engine with `-tags webkit2_41` support
- **Node.js LTS**: Latest LTS version for frontend development
- **npm**: Latest npm version for package management
- **Cross-platform compilation**: Support for Linux, Windows, and macOS targets
- **X11 development**: Full X11 development libraries
- **GTK development**: Complete GTK3 development environment
- **Security**: Non-root user for development

## Supported Platforms

### Linux
- `linux/amd64` (x86_64)
- `linux/arm64` (ARM64)
- `linux/arm-7` (ARMv7)

### Windows
- `windows/amd64` (x86_64)
- `windows/386` (x86)

### macOS
- `darwin/amd64` (Intel Mac)
- `darwin/arm64` (Apple Silicon)

## Quick Start

### Using Distrobox (Recommended for Development)

1. **Build the image:**
   ```bash
   ./build.sh build
   ```

2. **Create Distrobox container:**
   ```bash
   distrobox create wails-dev --image wails-dev:latest
   ```

3. **Enter the container:**
   ```bash
   distrobox enter wails-dev
   ```

4. **Create a new Wails project:**
   ```bash
   wails init -n my-app
   cd my-app
   ```

5. **Build with WebKit2_41 tags:**
   ```bash
   wails_with_tags build
   ```

### Using Docker Compose

1. **Build the image:**
   ```bash
   docker-compose build
   ```

2. **Start development environment:**
   ```bash
   docker-compose up -d
   ```

3. **Access the container:**
   ```bash
   docker-compose exec wails-dev bash
   ```

4. **Create a new Wails project:**
   ```bash
   wails init -n my-app
   cd my-app
   ```

5. **Build with WebKit2_41 tags:**
   ```bash
   wails_with_tags build
   ```

### Using Docker directly

1. **Build the image:**
   ```bash
   docker build -t wails-dev:latest .
   ```

2. **Run the container:**
   ```bash
   docker run -it --rm \
     -v $(pwd):/app \
     -v go-cache:/home/wailsdev/go \
     -e DISPLAY=$DISPLAY \
     -v /tmp/.X11-unix:/tmp/.X11-unix:rw \
     wails-dev:latest
   ```

## Usage Examples

### Development

```bash
# Start development server with WebKit2_41 tags
wails_with_tags dev

# Build for current platform with WebKit2_41 tags
wails_with_tags build

# Build for specific platform
wails build -platform linux/amd64 -tags webkit2_41
```

### Cross-platform Compilation

```bash
# Build for multiple platforms using xgo
xgo_build "linux/amd64,windows/amd64,darwin/amd64"

# Build for specific platform
xgo_build "linux/arm64"

# Build with custom package path
xgo_build "linux/amd64" "./cmd/myapp"
```

### Standard Wails Commands

```bash
# Initialize new project
wails init -n my-app

# Generate bindings
wails generate module

# Package application
wails build -package

# Build with specific tags
wails build -tags webkit2_41
```

## Environment Variables

- `WEBKIT_TAGS`: Automatically set to `-tags webkit2_41` for Ubuntu 24.04 compatibility
- `CGO_ENABLED`: Set to `1` for CGO support
- `GOPATH`: Set to `/home/wailsdev/go`
- `DISPLAY`: For X11 GUI applications
- `DISTROBOX_HOST_HOME`: Set to host home directory for Distrobox compatibility

## Distrobox Compatibility

This Docker image is specifically designed for use with [Distrobox](https://github.com/89luca89/distrobox). Key features:

- **User ID 1000**: Container user is created with UID 1000 to match typical host user
- **Proper permissions**: All directories are accessible to the Distrobox user
- **Environment detection**: Automatically detects Distrobox environment
- **Go workspace setup**: Creates and configures Go workspace directories

### Distrobox Usage

```bash
# Create container with the Docker image
distrobox create wails-dev --image wails-dev:latest

# Enter the container
distrobox enter wails-dev

# Your host user (UID 1000) will have full access to all directories
ls -la /home/wailsdev/go/
```

### Advanced Distrobox Options

You can also specify additional options when creating the container:

```bash
# Create with specific user mapping
distrobox create wails-dev --image wails-dev:latest --user 1000:1000

# Create with additional environment variables
distrobox create wails-dev --image wails-dev:latest --additional-flags "--env WEBKIT_TAGS=-tags webkit2_41"

# Create with X11 support
distrobox create wails-dev --image wails-dev:latest --additional-flags "-e DISPLAY=\$DISPLAY -v /tmp/.X11-unix:/tmp/.X11-unix:rw"
```

### Managing Distrobox Container

```bash
# List all containers
distrobox list

# Stop the container
distrobox stop wails-dev

# Start the container
distrobox start wails-dev

# Remove the container
distrobox rm wails-dev

# Export container to image
distrobox export wails-dev wails-dev-exported
```

## WebKit2GTK-4.1 Compatibility

This image includes WebKit2GTK-4.1 and automatically applies the `-tags webkit2_41` flag for Ubuntu 24.04 compatibility. The `wails_with_tags` function automatically adds these tags to build and dev commands.

## Cross-compilation with xgo

The image includes [xgo](https://github.com/crazy-max/xgo) for cross-compilation. Use the `xgo_build` function for easy cross-platform builds:

```bash
# Build for all major platforms
xgo_build "linux/amd64,linux/arm64,windows/amd64,windows/386,darwin/amd64,darwin/arm64"

# Build for specific platforms
xgo_build "linux/amd64,windows/amd64"
```

## Dependencies Included

### System Libraries
- WebKit2GTK-4.1 development libraries
- GTK3 development libraries
- X11 development libraries
- GStreamer development libraries
- Audio/video codec libraries

### Build Tools
- Go 1.24
- GCC/G++ compilers
- MinGW-w64 (Windows cross-compilation)
- Clang/LLVM (macOS cross-compilation)
- CMake and Ninja build systems

### Development Tools
- Git
- Wails CLI
- xgo cross-compiler
- pkg-config
- Node.js LTS
- npm

## Troubleshooting

### X11 Display Issues

If you encounter X11 display issues:

1. **Allow X11 connections:**
   ```bash
   xhost +local:docker
   ```

2. **Check DISPLAY variable:**
   ```bash
   echo $DISPLAY
   ```

3. **Use X11 forwarding:**
   ```bash
   docker run -it --rm \
     -e DISPLAY=$DISPLAY \
     -v /tmp/.X11-unix:/tmp/.X11-unix:rw \
     wails-dev:latest
   ```

### WebKit Issues

If you encounter WebKit-related issues:

1. **Verify WebKit installation:**
   ```bash
   pkg-config --modversion webkit2gtk-4.1
   ```

2. **Check WebKit tags:**
   ```bash
   echo $WEBKIT_TAGS
   ```

3. **Build with explicit tags:**
   ```bash
   wails build -tags webkit2_41
   ```

### Cross-compilation Issues

For cross-compilation problems:

1. **Check available targets:**
   ```bash
   xgo -h
   ```

2. **Verify toolchain installation:**
   ```bash
   which x86_64-w64-mingw32-gcc
   which clang
   ```

3. **Build with verbose output:**
   ```bash
   xgo -v -targets=linux/amd64 .
   ```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test the Docker image
5. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- [Wails](https://wails.io/) - The amazing framework this image supports
- [xgo](https://github.com/crazy-max/xgo) - Cross-compilation tool
- [WebKit2GTK](https://webkitgtk.org/) - Web engine for GTK applications 