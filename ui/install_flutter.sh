#!/bin/bash

# Script to install Flutter on Linux with desktop support

set -e  # Exit immediately if a command exits with a non-zero status

echo "Installing Flutter on Linux with desktop support..."

# Create a directory for Flutter in the user's home
FLUTTER_HOME="$HOME/flutter"
echo "Installing Flutter to: $FLUTTER_HOME"

# Check if Flutter is already installed
if [ -d "$FLUTTER_HOME" ]; then
    echo "Flutter is already installed at $FLUTTER_HOME"
    echo "Do you want to reinstall? (y/N): "
    read -r response
    if [[ ! "$response" =~ ^([yY][eE][sS]|[yY])$ ]]; then
        echo "Exiting without reinstalling."
        exit 0
    fi
    rm -rf "$FLUTTER_HOME"
fi
set -x
# Download the latest Flutter SDK
echo "Downloading Flutter SDK..."
cd /tmp
# LATEST_FLUTTER_URL=$(curl -s https://storage.googleapis.com/flutter_infra_release/releases/releases_linux.json | grep -o 'https://storage.googleapis.com/flutter_infra_release/releases/stable/linux/flutter-linux-.*-stable.tar.xz' | head -n 1)

# if [ -z "$LATEST_FLUTTER_URL" ]; then
#     echo "Could not fetch the latest Flutter release URL. Please check https://flutter.dev/docs/get-started/install/linux for the latest download link."
#     exit 1
# fi

# echo "Latest Flutter URL: $LATEST_FLUTTER_URL"
# if [[ ! -f flutter.tar.xz ]];then
#     wget "$LATEST_FLUTTER_URL" -O flutter.tar.xz
# fi
# # Extract Flutter to the home directory
# echo "Extracting Flutter..."
# tar xf flutter.tar.xz -C "$HOME"
# rm flutter.tar.xz

# Add Flutter to PATH in .bashrc or .zshrc
echo "Adding Flutter to PATH..."

# Detect the shell being used
SHELL_NAME=$(basename "$SHELL")
CONFIG_FILE=""

if [ "$SHELL_NAME" = "bash" ]; then
    CONFIG_FILE="$HOME/.bashrc"
elif [ "$SHELL_NAME" = "zsh" ]; then
    CONFIG_FILE="$HOME/.zshrc"
else
    # Default to .bashrc
    CONFIG_FILE="$HOME/.bashrc"
fi

# Add Flutter to PATH if not already added
if ! grep -q "export PATH.*\$PATH" "$CONFIG_FILE" 2>/dev/null || ! grep -q "$FLUTTER_HOME/bin" "$CONFIG_FILE" 2>/dev/null; then
    echo "" >> "$CONFIG_FILE"
    echo "# Flutter" >> "$CONFIG_FILE"
    echo "export PATH=\"\$PATH:$FLUTTER_HOME/bin\"" >> "$CONFIG_FILE"
    echo "Added Flutter to PATH in $CONFIG_FILE"
else
    echo "Flutter already in PATH"
fi

# Source the updated configuration for the current session
export PATH="$PATH:$FLUTTER_HOME/bin"

echo "Flutter installed successfully!"
echo "Current Flutter version:"
flutter --version

# Install dependencies for Flutter desktop development on Linux
echo "Installing dependencies for Flutter desktop development..."

# Check if we're on a Debian-based system (Ubuntu, etc.)
if [ -f /etc/debian_version ] || grep -qi "ubuntu" /etc/os-release; then
    echo "Detected Debian-based system, installing dependencies..."
    sudo apt update
    sudo apt install -y clang cmake ninja-build pkg-config libgtk-3-dev liblzma-dev
elif [ -f /etc/redhat-release ] || grep -qi "fedora\|centos\|rhel" /etc/os-release; then
    echo "Detected Red Hat-based system, installing dependencies..."
    sudo dnf install -y clang cmake ninja-build pkg-config gtk3-devel libX11-devel libXext-devel libXrender-devel libXrandr-devel libXcursor-devel libXinerama-devel libXi-devel libXss-devel libwayland-devel libxkbcommon-devel
elif [ -f /etc/arch-release ] || grep -qi "arch" /etc/os-release; then
    echo "Detected Arch-based system, installing dependencies..."
    sudo pacman -S clang cmake ninja pkg-config gtk3 libx11 libxext libxrender libxrandr libxinerama libxcursor libxi libxss wayland libxkbcommon
else
    echo "Please install the necessary build dependencies for your Linux distribution."
    echo "For Debian-based systems: clang cmake ninja-build pkg-config libgtk-3-dev"
    echo "For Red Hat-based systems: clang cmake ninja-build pkg-config gtk3-devel"
    echo "For Arch-based systems: clang cmake ninja pkg-config gtk3"
fi

echo "Flutter installation complete!"
echo "Run 'source $CONFIG_FILE' or restart your terminal to update your PATH."
echo "Then run 'flutter doctor' to verify the installation."
echo "To enable desktop support, run 'flutter config --enable-linux-desktop'"