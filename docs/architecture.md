<!-- markdownlint-disable MD013 -->
# Architecture

## Overview

The Kubernetes testing framework is designed as a modular, extensible tool for testing Kubernetes clusters across multiple dimensions: conformance, operational readiness, and performance.

## Component Diagram

```text
┌─────────────────────────────────────────────┐
│         CLI Layer (cmd/ktest)               │
│  ┌──────────────────────────────────────┐   │
│  │  conformance | operational | perf    │   │
│  └──────────────────────────────────────┘   │
└───────────────────┬─────────────────────────┘
                    │
┌───────────────────▼─────────────────────────┐
│       Test Execution Layer (pkg/)           │
│  ┌──────────────────────────────────────┐   │
│  │  • Conformance (Sonobuoy)            │   │
│  │  • Networking Tests                  │   │
│  │  • Storage Tests                     │   │
│  │  • Workload Tests                    │   │
│  │  • Performance Tests                 │   │
│  │  • Report Generation                 │   │
│  └──────────────────────────────────────┘   │
└───────────────────┬─────────────────────────┘
                    │
┌───────────────────▼─────────────────────────┐
│   Kubernetes Client Layer (client-go)       │
│  ┌──────────────────────────────────────┐   │
│  │  kubeconfig → clientset → API calls  │   │
│  └──────────────────────────────────────┘   │
└───────────────────┬─────────────────────────┘
                    │
┌───────────────────▼─────────────────────────┐
│         Kubernetes Cluster                  │
│  ┌──────────────────────────────────────┐   │
│  │  Pods, Services, Deployments, etc.   │   │
│  └──────────────────────────────────────┘   │
└─────────────────────────────────────────────┘
```

## Layers

### 1. CLI Layer

**Location**: `cmd/ktest/`, `pkg/cmd/`

**Responsibilities**:

- Command-line interface using Cobra
- Flag parsing and validation
- User interaction and output formatting
- Command routing

**Commands**:

- `conformance`: Run Sonobuoy-based conformance tests
- `operational`: Run networking, storage, and workload tests
- `performance`: Run load and performance tests

### 2. Test Execution Layer

**Location**: `pkg/*/`

**Responsibilities**:

- Implement test logic
- Interact with Kubernetes API
- Collect and aggregate results
- Generate reports

**Packages**:

- `conformance/`: Sonobuoy integration for conformance testing _(in development)_
- `networking/`: DNS, connectivity, service tests
- `storage/`: PVC, storage class tests
- `workload/`: Deployment, StatefulSet, DaemonSet tests
- `performance/`: Load testing and metrics
- `report/`: Report generation (HTML, JSON, console)

### 3. Kubernetes Client Layer

**Location**: `pkg/kubeconfig/`

**Responsibilities**:

- Load kubeconfig files
- Create Kubernetes clientset
- Manage API client connections

**Key Types**:

```go
type Client struct {
    Clientset *kubernetes.Clientset
    Config    *rest.Config
}
```

## Test Flow

### Conformance Test Flow _(Planned - In Development)_

> **Note:** Full Sonobuoy integration is not yet implemented. The flow below describes the planned architecture.

1. User executes `ktest conformance`
2. CLI parses flags (mode, kubeconfig)
3. Kubeconfig is loaded, clientset created
4. Sonobuoy client initialized (planned)
5. Conformance tests deployed to cluster (planned)
6. Results retrieved and formatted (planned)
7. Report generated (planned)

### Operational Test Flow

1. User executes `ktest operational`
2. CLI parses test selection flags
3. Kubeconfig loaded, clientset created
4. For each selected test category:
   - Create test resources (pods, services, etc.)
   - Verify expected behavior
   - Collect results
   - Clean up resources
5. Aggregate results
6. Generate report

### Performance Test Flow

1. User executes `ktest performance`
2. CLI parses endpoint, duration, RPS
3. Load test initialized
4. HTTP requests sent at specified rate
5. Latency and throughput metrics collected
6. Statistical analysis performed
7. Performance report generated

## Data Flow

```text
User Input
    ↓
CLI Flags
    ↓
Kubeconfig → K8s Client
    ↓
Test Execution
    ↓
Results Collection
    ↓
Report Generation
    ↓
Output (Console/File)
```

## Testing Strategy

### Unit Tests

**Location**: `tests/unit/`

- Mock Kubernetes clientsets
- Test individual functions
- Fast execution
- No cluster required

### Integration Tests

**Location**: `tests/integration/`

- Real Kubernetes cluster (kind)
- End-to-end test scenarios
- Slower execution
- Cluster required

## Configuration

### Test Configuration

**Location**: `configs/tests/`

YAML files define test parameters:

- Timeouts
- Replica counts
- Storage classes
- Enable/disable flags

### Sonobuoy Configuration

**Location**: `configs/sonobuoy/`

Sonobuoy-specific settings:

- Plugin configuration
- Result formats
- Mode settings

## Extensibility

### Adding New Tests

1. Create test package in `pkg/`
2. Implement test functions
3. Add CLI command in `pkg/cmd/`
4. Create configuration file
5. Add to operational runner

Example:

```go
// pkg/mytests/custom.go
func TestCustomFeature(ctx context.Context, clientset *kubernetes.Clientset) error {
    // Implementation
}

// pkg/cmd/operational.go
if runCustom {
    mytests.TestCustomFeature(ctx, client.Clientset)
}
```

### Adding Report Formats

Extend `pkg/report/report.go`:

```go
func (r *TestReport) GenerateMarkdown() (string, error) {
    // Implementation
}
```

## Dependencies

### Core

- `k8s.io/client-go`: Kubernetes client library
- `github.com/spf13/cobra`: CLI framework
- `github.com/spf13/viper`: Configuration management

### Testing

- `github.com/stretchr/testify`: Assertion library
- `k8s.io/client-go/kubernetes/fake`: Mock Kubernetes client

### Optional

- `github.com/vmware-tanzu/sonobuoy`: Conformance testing (external binary)

## Security Considerations

- Kubeconfig files contain sensitive credentials
- Never commit kubeconfig to version control
- Test resources are cleaned up on completion
- RBAC permissions required for cluster operations

## Performance Considerations

- Tests run sequentially to avoid resource conflicts
- Configurable timeouts prevent hanging
- Resource cleanup is deferred for reliability
- Load tests use rate limiting

## Future Enhancements

1. Parallel test execution
2. Custom plugin framework
3. Web UI for results
4. Historical result tracking
5. Multi-cluster support
6. Enhanced reporting (PDF, email)
