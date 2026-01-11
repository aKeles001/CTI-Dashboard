# Wails Development Docker Image
# Based on Go 1.24 and Ubuntu 24.04 with cross-compilation support

FROM ubuntu:24.04

# Set environment variables
ENV DEBIAN_FRONTEND=noninteractive
ENV GO_VERSION=1.24.0
ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=1
ENV PATH="/usr/local/go/bin:${PATH}"

# Install system dependencies
RUN apt-get update && apt-get install -y \
    # Basic build tools
    build-essential \
    pkg-config \
    cmake \
    ninja-build \
    git \
    wget \
    curl \
    unzip \
    software-properties-common \
    apt-transport-https \
    ca-certificates \
    gnupg \
    lsb-release \
    && rm -rf /var/lib/apt/lists/*

# Install Go 1.24
RUN wget -q https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz \
    && tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz \
    && rm go${GO_VERSION}.linux-amd64.tar.gz

# Install WebKit2GTK and GTK development libraries
RUN apt-get update && apt-get install -y \
    # WebKit2GTK and GTK
    libwebkit2gtk-4.1-dev \
    libgtk-3-dev \
    libgdk-pixbuf2.0-dev \
    libcairo2-dev \
    libpango1.0-dev \
    libatk1.0-dev \
    libgirepository1.0-dev \
    libglib2.0-dev \
    libgstreamer1.0-dev \
    libgstreamer-plugins-base1.0-dev \
    # X11 development
    libx11-dev \
    libxrandr-dev \
    libxinerama-dev \
    libxcursor-dev \
    libxcomposite-dev \
    libxdamage-dev \
    libxext-dev \
    libxfixes-dev \
    libxrender-dev \
    libxss-dev \
    libxtst-dev \
    # Additional dependencies
    libssl-dev \
    libsqlite3-dev \
    libpulse-dev \
    libasound2-dev \
    libdbus-1-dev \
    libudev-dev \
    libusb-1.0-0-dev \
    libavcodec-dev \
    libavformat-dev \
    libswscale-dev \
    libavutil-dev \
    libavfilter-dev \
    libavdevice-dev \
    libswresample-dev \
    && rm -rf /var/lib/apt/lists/*

# Install Google Chrome
RUN wget -q https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb \
    && apt-get update && apt-get install -y ./google-chrome-stable_current_amd64.deb \
    && rm google-chrome-stable_current_amd64.deb \
    && rm -rf /var/lib/apt/lists/*

# Install Node.js and npm (latest LTS version)
RUN curl -fsSL https://deb.nodesource.com/setup_lts.x | bash - \
    && apt-get install -y nodejs \
    && npm install -g npm@latest

# Install cross-compilation toolchains
RUN apt-get update && apt-get install -y \
    # Windows cross-compilation
    gcc-mingw-w64 \
    g++-mingw-w64 \
    # macOS cross-compilation (using osxcross)
    clang \
    llvm \
    libxml2-dev \
    libssl-dev \
    libbz2-dev \
    libreadline-dev \
    libsqlite3-dev \
    libncurses5-dev \
    libncursesw5-dev \
    xz-utils \
    tk-dev \
    libffi-dev \
    liblzma-dev \
    && rm -rf /var/lib/apt/lists/*

# Set up working directory
WORKDIR /app

# Create a non-root user for development (UID 1000 for Distrobox compatibility)
RUN if id 1000 >/dev/null 2>&1; then \
        # UID 1000 exists, remove it and create wailsdev \
        userdel -r $(id -un 1000) 2>/dev/null || true; \
    fi && \
    useradd -m -s /bin/bash -u 1000 wailsdev && \
    echo "Created user wailsdev with UID 1000" && \
    chown -R wailsdev:wailsdev /app

# Install Wails CLI and xgo as wailsdev user
USER wailsdev
RUN go install github.com/wailsapp/wails/v2/cmd/wails@latest
RUN go install src.techknowlogick.com/xgo@latest
USER root

# Set up Go environment for the user
ENV GOPATH=/home/wailsdev/go
ENV PATH="/home/wailsdev/go/bin:/usr/local/go/bin:${PATH}"

# Create necessary directories with proper permissions for Distrobox
RUN mkdir -p /home/wailsdev/go/bin \
    && mkdir -p /home/wailsdev/go/pkg \
    && mkdir -p /home/wailsdev/go/src \
    && chown -R wailsdev:wailsdev /home/wailsdev/go

# Create entrypoint script
COPY entrypoint.sh /entrypoint.sh
RUN chown wailsdev:wailsdev /entrypoint.sh && chmod +x /entrypoint.sh

# Switch to non-root user (use the username we determined)
USER wailsdev

ENTRYPOINT ["/entrypoint.sh"]
CMD ["bash"] 
