#!/bin/bash

echo "Building the project..."

# Nome do binário
BINARY_NAME="maestro"

# Diretório de build
BUILD_DIR="build"

# Criação do diretório de build, caso não exista
if [ ! -d "$BUILD_DIR" ]; then
  mkdir "$BUILD_DIR"
fi

# Step 1: Install required dependencies
echo "Installing Go dependencies..."
go mod tidy

# Compilando o projeto
echo "Building the project..."
go build -o "$BUILD_DIR/$BINARY_NAME" ./cmd/main.go

# Verificando se a compilação foi bem-sucedida
if [ $? -eq 0 ]; then
  echo "Build successful! The binary is located at $BUILD_DIR/$BINARY_NAME"
else
  echo "Build failed!"
fi

# Gabiarra
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
