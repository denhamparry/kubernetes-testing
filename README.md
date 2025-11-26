# Kubernetes Testing Framework

A comprehensive testing framework for Kubernetes clusters, supporting conformance testing, operational validation, and performance testing.

## Features

- Kubernetes conformance testing (via Sonobuoy)
- Operational testing (networking, storage, workloads)
- Performance testing and benchmarking
- Works with any Kubernetes cluster (kubeconfig-based)
- Detailed reporting (HTML, JSON, console)
- CI/CD integration ready

## Quick Start

```bash
# Setup
./scripts/setup.sh

# Run all tests
./scripts/run-tests.sh all

# Run specific test suite
./bin/ktest conformance --kubeconfig ~/.kube/config
```

## Installation

### Prerequisites

- Go 1.21+
- kubectl
- Access to a Kubernetes cluster
- Valid kubeconfig file

### Build from Source

```bash
git clone https://github.com/denhamparry/kubernetes-testing.git
cd kubernetes-testing
./scripts/setup.sh
```

## Usage

### Conformance Tests

Validate cluster compliance with Kubernetes standards:

```bash
./bin/ktest conformance --kubeconfig ~/.kube/config --mode quick
```

### Operational Tests

Test networking, storage, and workloads:

```bash
# All operational tests
./bin/ktest operational --kubeconfig ~/.kube/config

# Specific test categories
./bin/ktest operational --tests networking,storage
```

### Performance Tests

Run load tests and collect performance metrics:

```bash
./bin/ktest performance \
  --endpoint http://my-app.example.com \
  --duration 5m \
  --rps 100
```

## Test Categories

### Conformance

- Sonobuoy integration
- CNCF conformance validation
- Quick and certified modes

### Networking

- DNS resolution
- Pod-to-pod connectivity
- Service endpoint access

### Storage

- PVC creation and binding
- Storage class validation
- Dynamic provisioning

### Workloads

- Deployment creation and scaling
- StatefulSet ordered deployment
- DaemonSet node coverage

### Performance

- HTTP load testing
- Latency metrics (avg, p95, p99)
- Throughput analysis
- Error rate tracking

## Configuration

Tests can be customized via YAML files in `configs/tests/`:

- `networking.yaml` - Network test parameters
- `storage.yaml` - Storage test configuration
- `workload.yaml` - Workload test settings
- `sonobuoy/quick-config.yaml` - Conformance test config

## Documentation

- [Usage Guide](docs/usage.md) - Detailed usage instructions
- [Architecture](docs/architecture.md) - System design and components
- [Contributing](CONTRIBUTING.md) - How to contribute

## CI/CD Integration

The project includes GitHub Actions workflows:

- Lint (golangci-lint)
- Unit tests with coverage
- Build verification
- Integration tests with kind

### Example Pipeline Integration

```yaml
- name: Run Kubernetes tests
  run: |
    ./scripts/setup.sh
    ./bin/ktest operational --kubeconfig $KUBECONFIG
```

## Development

### Running Tests

```bash
# Unit tests
go test -v ./...

# Integration tests (requires cluster)
go test -v -tags=integration ./tests/integration/...

# With coverage
go test -v -race -coverprofile=coverage.txt ./...
```

### Building

```bash
# Build CLI
go build -o bin/ktest cmd/ktest/main.go

# Run locally
./bin/ktest --help
```

## Project Structure

```text
kubernetes-testing/
├── cmd/ktest/          # CLI entry point
├── pkg/                # Core packages
│   ├── cmd/           # CLI commands
│   ├── conformance/   # Sonobuoy integration
│   ├── networking/    # Network tests
│   ├── storage/       # Storage tests
│   ├── workload/      # Workload tests
│   ├── performance/   # Performance tests
│   ├── kubeconfig/    # Kubeconfig handling
│   └── report/        # Report generation
├── tests/             # Test files
│   ├── unit/         # Unit tests
│   └── integration/  # Integration tests
├── configs/           # Configuration files
│   ├── tests/        # Test configs
│   └── sonobuoy/     # Sonobuoy configs
├── scripts/           # Helper scripts
└── docs/             # Documentation
```

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

[Add your license here]

## Support

For issues and questions:

- Open an issue: <https://github.com/denhamparry/kubernetes-testing/issues>
- Documentation: See `docs/` directory
- Claude Code: <https://docs.claude.com/en/docs/claude-code>

## Acknowledgments

- Built with [Kubernetes client-go](https://github.com/kubernetes/client-go)
- Conformance testing via [Sonobuoy](https://sonobuoy.io/)
- CLI framework: [Cobra](https://github.com/spf13/cobra)

---

**Version:** 1.0.0
**Status:** Active Development
