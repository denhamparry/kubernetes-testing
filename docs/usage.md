<!-- markdownlint-disable MD013 -->
# Usage Guide

## Installation

```bash
# Clone repository
git clone https://github.com/denhamparry/kubernetes-testing.git
cd kubernetes-testing

# Run setup
./scripts/setup.sh
```

## Prerequisites

- Go 1.21 or higher
- kubectl
- Access to a Kubernetes cluster
- Valid kubeconfig file

## Running Tests

### Conformance Tests

> **Note:** Conformance testing integration with Sonobuoy is currently in development. The command structure exists but full Sonobuoy integration is not yet implemented. See the project roadmap for planned implementation.

Run Kubernetes conformance tests using Sonobuoy (coming soon):

```bash
# Planned - full Sonobuoy integration in development
./bin/ktest conformance --kubeconfig ~/.kube/config --mode quick
```

Planned modes:

- `quick`: Fast conformance tests (~10 minutes)
- `certified-conformance`: Full conformance test suite (~2 hours)

### Operational Tests

Run operational tests for networking, storage, and workloads:

```bash
# All operational tests
./bin/ktest operational --kubeconfig ~/.kube/config

# Specific tests
./bin/ktest operational --tests networking,storage --kubeconfig ~/.kube/config

# With custom namespace
./bin/ktest operational --namespace test-ns --kubeconfig ~/.kube/config
```

Available test categories:

- `networking`: DNS, pod-to-pod, service connectivity
- `storage`: PVC creation, storage classes
- `workload`: Deployments, StatefulSets, DaemonSets

### Performance Tests

Run performance and load tests:

```bash
./bin/ktest performance \
  --endpoint http://my-app.example.com \
  --duration 10m \
  --rps 200 \
  --kubeconfig ~/.kube/config
```

Parameters:

- `--endpoint`: Target endpoint to test (required)
- `--duration`: Test duration (default: 5m)
- `--rps`: Requests per second (default: 100)

## Configuration

Tests can be configured via YAML files in `configs/tests/`.

### Networking Configuration

Edit `configs/tests/networking.yaml`:

```yaml
tests:
  - name: dns-resolution
    enabled: true
    timeout: 60s
  - name: pod-to-pod
    enabled: true
    timeout: 120s
```

### Storage Configuration

Edit `configs/tests/storage.yaml`:

```yaml
tests:
  - name: pvc-creation
    enabled: true
    storageClass: standard
    size: 1Gi
```

### Workload Configuration

Edit `configs/tests/workload.yaml`:

```yaml
tests:
  - name: deployment
    enabled: true
    replicas: 3
```

## Reports

Test reports are generated in multiple formats:

- **HTML**: `reports/test-report.html`
- **JSON**: `reports/test-report.json`
- **Console**: Real-time output

## Examples

### Quick cluster validation

```bash
# Run all operational tests
./scripts/run-tests.sh operational
```

### Full test suite

```bash
# Run conformance and operational tests
./scripts/run-tests.sh all
```

### Performance baseline

```bash
# Test application performance
./bin/ktest performance --endpoint http://my-app.example.com --duration 5m --rps 50
```

## Troubleshooting

### kubectl not found

Ensure kubectl is installed and in your PATH:

```bash
which kubectl
```

### Invalid kubeconfig

Verify your kubeconfig is valid:

```bash
kubectl cluster-info --kubeconfig ~/.kube/config
```

### Permission denied

Make scripts executable:

```bash
chmod +x scripts/*.sh
```

### Sonobuoy not found

Install Sonobuoy for conformance testing:

```bash
# macOS
brew install sonobuoy

# Linux
curl -LO https://github.com/vmware-tanzu/sonobuoy/releases/download/v0.57.0/sonobuoy_0.57.0_linux_amd64.tar.gz
tar -xzf sonobuoy_0.57.0_linux_amd64.tar.gz
sudo mv sonobuoy /usr/local/bin/
```

## CI/CD Integration

### GitHub Actions

The repository includes a CI workflow that runs tests automatically:

- Lint: golangci-lint
- Unit tests: All packages
- Build: CLI binary
- Integration tests: kind cluster

### Custom Pipelines

Integrate into your pipeline:

```yaml
- name: Run Kubernetes tests
  run: |
    ./scripts/setup.sh
    ./bin/ktest operational --kubeconfig $KUBECONFIG
```
