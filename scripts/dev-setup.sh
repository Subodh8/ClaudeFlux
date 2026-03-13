#!/usr/bin/env bash
set -euo pipefail

echo "→ Setting up ClaudeFlux development environment"

# Check Go
if ! command -v go &> /dev/null; then
    echo "✗ Go is not installed. Install Go 1.23+ from https://go.dev/dl/"
    exit 1
fi
echo "✓ Go $(go version | awk '{print $3}')"

# Check Node.js
if ! command -v node &> /dev/null; then
    echo "✗ Node.js is not installed. Install Node.js 20+ from https://nodejs.org/"
    exit 1
fi
echo "✓ Node.js $(node --version)"

# Check Git
if ! command -v git &> /dev/null; then
    echo "✗ Git is not installed."
    exit 1
fi
echo "✓ Git $(git --version | awk '{print $3}')"

# Install Go dependencies
echo "→ Downloading Go dependencies..."
go mod download
echo "✓ Go dependencies installed"

# Install dashboard dependencies
echo "→ Installing dashboard dependencies..."
cd dashboard
npm install
cd ..
echo "✓ Dashboard dependencies installed"

# Build
echo "→ Building ClaudeFlux..."
make build
echo "✓ Build complete"

# Run tests
echo "→ Running tests..."
make test
echo "✓ Tests passed"

echo ""
echo "✓ Development environment ready!"
echo ""
echo "  Quick commands:"
echo "    make build        Build the binary"
echo "    make test         Run tests"
echo "    make dev          Run with hot reload"
echo "    make example      Run example workflow"
echo "    make docker       Start with Docker Compose"
