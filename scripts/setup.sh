#!/bin/bash
set -e

echo "Setting up Kubernetes testing environment..."

# Check prerequisites
echo "Checking prerequisites..."
command -v kubectl >/dev/null 2>&1 || { echo "ERROR: kubectl is required but not installed. Aborting." >&2; exit 1; }
command -v go >/dev/null 2>&1 || { echo "ERROR: Go is required but not installed. Aborting." >&2; exit 1; }

echo "  ✓ kubectl found"
echo "  ✓ go found"

# Check Go version
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
echo "  Go version: $GO_VERSION"

# Install dependencies
echo "Installing Go dependencies..."
go mod download
go mod tidy

# Create bin directory if it doesn't exist
mkdir -p bin

# Build CLI
echo "Building ktest CLI..."
go build -o bin/ktest cmd/ktest/main.go

echo "  ✓ ktest binary created at bin/ktest"

# Check if Sonobuoy is installed
if ! command -v sonobuoy >/dev/null 2>&1; then
    echo "WARNING: Sonobuoy is not installed. Conformance tests will not work without it."
    echo "Install from: https://sonobuoy.io/docs/latest/"
else
    echo "  ✓ sonobuoy found"
fi

echo ""
echo "========================================="
echo "Setup complete!"
echo "========================================="
echo ""
echo "Run tests with:"
echo "  ./bin/ktest --help"
echo "  ./bin/ktest conformance --kubeconfig ~/.kube/config"
echo "  ./bin/ktest operational --kubeconfig ~/.kube/config"
echo "  ./bin/ktest performance --endpoint http://example.com"
echo ""
