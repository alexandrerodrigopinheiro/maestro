#!/bin/bash

# This script builds and installs the Maestro CLI on a Debian system.

# Update package lists and install Go if not already installed
if ! command -v go &> /dev/null; then
  echo "Go not found. Installing Go..."
  sudo apt update
  sudo apt install -y golang
else
  echo "Go is already installed."
fi

# Set environment variables for Go
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin

# Clone the Maestro repository if it does not exist
if [ ! -d "$HOME/maestro" ]; then
  echo "Cloning Maestro repository..."
  git clone https://github.com/alexandrerodrigopinheiro/maestro.git $HOME/maestro
fi

# Navigate to Maestro directory
cd $HOME/maestro || exit

# Build the Maestro binary
echo "Building Maestro..."
go build -o build/maestro cmd/main.go

# Make the binary accessible from ~/.local/share
mkdir -p ~/.local/share/maestro
mv build/maestro ~/.local/share/maestro/

# Add Maestro to PATH and set environment variable if not already added
if ! grep -q 'export PATH="$PATH:$HOME/.local/share/maestro"' ~/.bashrc; then
  echo 'export PATH="$PATH:$HOME/.local/share/maestro"' >> ~/.bashrc
  export PATH="$PATH:$HOME/.local/share/maestro"
fi

if ! grep -q 'export MAESTRO_HOME="$HOME/.local/share/maestro"' ~/.bashrc; then
  echo 'export MAESTRO_HOME="$HOME/.local/share/maestro"' >> ~/.bashrc
  export MAESTRO_HOME="$HOME/.local/share/maestro"
fi

# Confirm installation
echo "Maestro installed successfully."
echo "You can now use Maestro by running: maestro <command>"
