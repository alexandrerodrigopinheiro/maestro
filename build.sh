#!/bin/bash

# Nome do binário
BINARY_NAME="maestro"

# Diretório de build
BUILD_DIR="build"

# Criação do diretório de build, caso não exista
if [ ! -d "$BUILD_DIR" ]; then
  mkdir "$BUILD_DIR"
fi

# Compilando o projeto
echo "Building the project..."
go build -o "$BUILD_DIR/$BINARY_NAME" ./cmd/main.go

# Verificando se a compilação foi bem-sucedida
if [ $? -eq 0 ]; then
  echo "Build successful! The binary is located at $BUILD_DIR/$BINARY_NAME"
else
  echo "Build failed!"
fi
