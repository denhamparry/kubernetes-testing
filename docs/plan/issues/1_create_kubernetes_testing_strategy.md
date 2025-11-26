# GitHub Issue #1: Create kubernetes testing strategy

**Issue:** [#1](https://github.com/denhamparry/kubernetes-testing/issues/1)
**Status:** Complete
**Date:** 2025-11-26
**Labels:** enhancement, testing, kubernetes

## Problem Statement

There is a need to test Kubernetes clusters for conformance, operational readiness, and performance. This includes day-to-day operations (networking, storage, workloads) and cluster performance testing. The testing should work with any Kubernetes cluster accessible via a kubeconfig file.

### Current Behavior

- No Kubernetes testing framework in place
- No conformance testing capability
- No operational testing for networking, storage, workloads
- No performance testing infrastructure

### Expected Behavior

- Ability to run conformance tests against any Kubernetes cluster
- Test day-to-day operations (networking, storage, workloads)
- Performance testing capabilities
- Use kubeconfig file for cluster access
- Automated and repeatable testing framework

## Current State Analysis

### Relevant Code/Config

This is a new project created from the GitHub template for Claude Code projects. The current repository structure includes:

- Template files: README.md, CLAUDE.md, CONTRIBUTING.md
- GitHub workflows: `.github/workflows/ci.yml` (placeholder)
- Pre-commit configuration: `.pre-commit-config.yaml`
- Custom slash commands in `.claude/commands/`
- No Kubernetes-specific code yet

### Related Context

- Repository: `denhamparry/kubernetes-testing`
- Branch: `denhamparry.co.uk/feat/gh-issue-001`
- This is a template repository being customized for Kubernetes testing
- Need to align with TDD principles emphasized in the template

## Solution Design

### Approach

The solution will implement a comprehensive Kubernetes testing framework using industry-standard tools:

1. **Sonobuoy** - For Kubernetes conformance testing
   - Official CNCF conformance testing tool
   - Supports custom plugins
   - Well-documented and widely adopted

2. **Custom Test Suites** - For operational testing
   - Networking tests (connectivity, DNS, service mesh)
   - Storage tests (PV/PVC, CSI drivers)
   - Workload tests (deployments, statefulsets, daemonsets)

3. **Performance Testing** - Using k6 or similar
   - Load testing for applications
   - Cluster resource utilization
   - Scalability testing

4. **Test Orchestration** - Using Go or Python
   - CLI tool for running test suites
   - Test result aggregation
   - Report generation

### Implementation

The implementation will follow this structure:

```text
kubernetes-testing/
├── cmd/
│   └── ktest/              # CLI tool for running tests
│       └── main.go
├── pkg/
│   ├── conformance/        # Sonobuoy integration
│   ├── networking/         # Network tests
│   ├── storage/            # Storage tests
│   ├── workload/           # Workload tests
│   ├── performance/        # Performance tests
│   └── kubeconfig/         # Kubeconfig handling
├── tests/
│   ├── unit/              # Unit tests
│   └── integration/       # Integration tests
├── configs/
│   ├── sonobuoy/          # Sonobuoy configurations
│   └── tests/             # Test configurations
├── scripts/
│   ├── setup.sh           # Setup script
│   └── run-tests.sh       # Test runner
└── docs/
    ├── usage.md           # Usage guide
    └── architecture.md    # Architecture docs
```

### Technology Stack Decision

**Option 1: Go** (Recommended)

- ✅ Native Kubernetes ecosystem language
- ✅ Excellent Kubernetes client libraries (client-go)
- ✅ Strong typing and compilation
- ✅ Great CLI frameworks (cobra, viper)
- ✅ Easy containerization
- ❌ More verbose than Python

**Option 2: Python**

- ✅ Easier scripting and testing
- ✅ Kubernetes Python client available
- ✅ Rich testing frameworks (pytest)
- ❌ Less common in Kubernetes tooling
- ❌ Dependency management complexity

**Decision:** Go - aligns with Kubernetes ecosystem standards and provides better performance for CLI tooling.

### Benefits

- Comprehensive testing coverage for Kubernetes clusters
- Automated conformance validation
- Reusable test suites for operational validation
- Performance baseline establishment
- Works with any Kubernetes cluster (on-prem, cloud, local)
- CI/CD integration ready

## Implementation Plan

### Step 1: Project Setup and Go Module Initialization

**File:** `go.mod`, `go.sum`

**Changes:**

- Initialize Go module: `go mod init github.com/denhamparry/kubernetes-testing`
- Add dependencies:
  - `k8s.io/client-go` - Kubernetes client
  - `github.com/spf13/cobra` - CLI framework
  - `github.com/spf13/viper` - Configuration management
  - Testing libraries: `github.com/stretchr/testify`

**Testing:**

```bash
go mod tidy
go mod verify
```

### Step 2: Create CLI Structure

**File:** `cmd/ktest/main.go`

**Changes:**

```go
package main

import (
    "github.com/denhamparry/kubernetes-testing/pkg/cmd"
)

func main() {
    cmd.Execute()
}
```

**File:** `pkg/cmd/root.go`

**Changes:**

```go
package cmd

import (
    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "ktest",
    Short: "Kubernetes cluster testing tool",
    Long:  `A comprehensive testing tool for Kubernetes clusters`,
}

func Execute() error {
    return rootCmd.Execute()
}

func init() {
    rootCmd.PersistentFlags().String("kubeconfig", "", "path to kubeconfig file")
}
```

**Testing:**

```bash
go run cmd/ktest/main.go --help
```

### Step 3: Implement Kubeconfig Handling

**File:** `pkg/kubeconfig/kubeconfig.go`

**Changes:**

```go
package kubeconfig

import (
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/clientcmd"
)

type Client struct {
    Clientset *kubernetes.Clientset
    Config    *rest.Config
}

func NewClient(kubeconfigPath string) (*Client, error) {
    // Load kubeconfig
    // Create kubernetes clientset
    // Return client
}
```

**Testing:**

```go
// tests/unit/kubeconfig_test.go
func TestNewClient(t *testing.T) {
    // Test kubeconfig loading
    // Test client creation
}
```

### Step 4: Implement Conformance Testing (Sonobuoy Integration)

**File:** `pkg/conformance/sonobuoy.go`

**Changes:**

```go
package conformance

import (
    "github.com/vmware-tanzu/sonobuoy/pkg/client"
)

type ConformanceTest struct {
    client *client.SonobuoyClient
}

func NewConformanceTest(kubeconfig string) (*ConformanceTest, error) {
    // Initialize Sonobuoy client
}

func (c *ConformanceTest) Run(mode string) error {
    // Run conformance tests
    // mode: "quick" or "certified-conformance"
}

func (c *ConformanceTest) GetResults() (*Results, error) {
    // Retrieve test results
}
```

**File:** `configs/sonobuoy/quick-config.yaml`

**Changes:**

```yaml
sonobuoy-config:
  driver: Job
  plugin-name: e2e
  result-format: junit
```

**Testing:**

```bash
# Integration test (requires cluster)
go test -tags=integration ./pkg/conformance/...
```

### Step 5: Implement Networking Tests

**File:** `pkg/networking/dns.go`

**Changes:**

```go
package networking

func TestDNS(clientset *kubernetes.Clientset) error {
    // Create test pod
    // Test DNS resolution (kubernetes.default, external DNS)
    // Clean up
}
```

**File:** `pkg/networking/connectivity.go`

**Changes:**

```go
package networking

func TestPodToPod(clientset *kubernetes.Clientset) error {
    // Deploy two pods
    // Test connectivity between pods
    // Clean up
}

func TestServiceConnectivity(clientset *kubernetes.Clientset) error {
    // Deploy service and pods
    // Test service endpoint access
    // Clean up
}
```

**Testing:**

```go
// tests/unit/networking_test.go
func TestNetworkingFunctions(t *testing.T) {
    // Unit tests with mock clientset
}
```

### Step 6: Implement Storage Tests

**File:** `pkg/storage/pvc.go`

**Changes:**

```go
package storage

func TestPVCCreation(clientset *kubernetes.Clientset, storageClass string) error {
    // Create PVC
    // Wait for bound status
    // Create pod using PVC
    // Write data to volume
    // Verify data persistence
    // Clean up
}

func TestStorageClass(clientset *kubernetes.Clientset) error {
    // List available storage classes
    // Test default storage class
    // Test dynamic provisioning
}
```

**Testing:**

```go
// tests/integration/storage_test.go
func TestStorageIntegration(t *testing.T) {
    // Integration tests with real cluster
}
```

### Step 7: Implement Workload Tests

**File:** `pkg/workload/deployment.go`

**Changes:**

```go
package workload

func TestDeployment(clientset *kubernetes.Clientset) error {
    // Create deployment
    // Verify replicas ready
    // Test rolling update
    // Test scaling
    // Clean up
}
```

**File:** `pkg/workload/statefulset.go`

**Changes:**

```go
package workload

func TestStatefulSet(clientset *kubernetes.Clientset) error {
    // Create statefulset
    // Verify ordered creation
    // Test persistent volumes
    // Test scaling
    // Clean up
}
```

**Testing:**

```bash
go test ./pkg/workload/...
```

### Step 8: Implement Performance Testing Framework

**File:** `pkg/performance/load.go`

**Changes:**

```go
package performance

type LoadTest struct {
    Duration  time.Duration
    RPS       int
    Endpoint  string
}

func (l *LoadTest) Run() (*Metrics, error) {
    // Execute load test
    // Collect metrics (latency, throughput, errors)
    // Return results
}
```

**File:** `pkg/performance/metrics.go`

**Changes:**

```go
package performance

type Metrics struct {
    AverageLatency time.Duration
    P95Latency     time.Duration
    P99Latency     time.Duration
    Throughput     float64
    ErrorRate      float64
}

func (m *Metrics) Report() string {
    // Generate performance report
}
```

**Testing:**

```go
// tests/unit/performance_test.go
func TestLoadTestMetrics(t *testing.T) {
    // Test metrics calculation
}
```

### Step 9: Create CLI Commands

**File:** `pkg/cmd/conformance.go`

**Changes:**

```go
package cmd

var conformanceCmd = &cobra.Command{
    Use:   "conformance",
    Short: "Run Kubernetes conformance tests",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Execute conformance tests
    },
}

func init() {
    rootCmd.AddCommand(conformanceCmd)
    conformanceCmd.Flags().String("mode", "quick", "Test mode: quick or certified-conformance")
}
```

**File:** `pkg/cmd/operational.go`

**Changes:**

```go
package cmd

var operationalCmd = &cobra.Command{
    Use:   "operational",
    Short: "Run operational tests",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Execute networking, storage, workload tests
    },
}

func init() {
    rootCmd.AddCommand(operationalCmd)
    operationalCmd.Flags().StringSlice("tests", []string{"all"}, "Tests to run: networking, storage, workload, all")
}
```

**File:** `pkg/cmd/performance.go`

**Changes:**

```go
package cmd

var performanceCmd = &cobra.Command{
    Use:   "performance",
    Short: "Run performance tests",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Execute performance tests
    },
}

func init() {
    rootCmd.AddCommand(performanceCmd)
    performanceCmd.Flags().Duration("duration", 5*time.Minute, "Test duration")
    performanceCmd.Flags().Int("rps", 100, "Requests per second")
}
```

**Testing:**

```bash
go run cmd/ktest/main.go conformance --help
go run cmd/ktest/main.go operational --help
go run cmd/ktest/main.go performance --help
```

### Step 10: Create Test Configuration Files

**File:** `configs/tests/networking.yaml`

**Changes:**

```yaml
tests:
  - name: dns-resolution
    enabled: true
    timeout: 60s
  - name: pod-to-pod
    enabled: true
    timeout: 120s
  - name: service-connectivity
    enabled: true
    timeout: 120s
```

**File:** `configs/tests/storage.yaml`

**Changes:**

```yaml
tests:
  - name: pvc-creation
    enabled: true
    storageClass: standard
    size: 1Gi
  - name: dynamic-provisioning
    enabled: true
```

**File:** `configs/tests/workload.yaml`

**Changes:**

```yaml
tests:
  - name: deployment
    enabled: true
    replicas: 3
  - name: statefulset
    enabled: true
    replicas: 3
  - name: daemonset
    enabled: true
```

**Testing:**

```bash
# Validate YAML syntax
yamllint configs/tests/*.yaml
```

### Step 11: Implement Reporting

**File:** `pkg/report/report.go`

**Changes:**

```go
package report

type TestReport struct {
    TestSuite   string
    StartTime   time.Time
    EndTime     time.Time
    TotalTests  int
    Passed      int
    Failed      int
    Skipped     int
    Results     []TestResult
}

type TestResult struct {
    Name     string
    Status   string
    Duration time.Duration
    Message  string
}

func (r *TestReport) GenerateHTML() (string, error) {
    // Generate HTML report
}

func (r *TestReport) GenerateJSON() (string, error) {
    // Generate JSON report
}

func (r *TestReport) Print() {
    // Print summary to console
}
```

**Testing:**

```go
// tests/unit/report_test.go
func TestReportGeneration(t *testing.T) {
    // Test report formats
}
```

### Step 12: Create Setup and Run Scripts

**File:** `scripts/setup.sh`

**Changes:**

```bash
#!/bin/bash
set -e

echo "Setting up Kubernetes testing environment..."

# Check prerequisites
command -v kubectl >/dev/null 2>&1 || { echo "kubectl is required"; exit 1; }
command -v go >/dev/null 2>&1 || { echo "Go is required"; exit 1; }

# Install dependencies
go mod download

# Build CLI
go build -o bin/ktest cmd/ktest/main.go

# Install Sonobuoy (if not present)
if ! command -v sonobuoy >/dev/null 2>&1; then
    echo "Installing Sonobuoy..."
    # Installation commands
fi

echo "Setup complete!"
```

**File:** `scripts/run-tests.sh`

**Changes:**

```bash
#!/bin/bash
set -e

KUBECONFIG="${KUBECONFIG:-$HOME/.kube/config}"
TEST_SUITE="${1:-all}"

echo "Running Kubernetes tests..."
echo "Kubeconfig: $KUBECONFIG"
echo "Test suite: $TEST_SUITE"

case $TEST_SUITE in
  conformance)
    ./bin/ktest conformance --kubeconfig "$KUBECONFIG"
    ;;
  operational)
    ./bin/ktest operational --kubeconfig "$KUBECONFIG"
    ;;
  performance)
    ./bin/ktest performance --kubeconfig "$KUBECONFIG"
    ;;
  all)
    ./bin/ktest conformance --kubeconfig "$KUBECONFIG"
    ./bin/ktest operational --kubeconfig "$KUBECONFIG"
    ./bin/ktest performance --kubeconfig "$KUBECONFIG"
    ;;
  *)
    echo "Unknown test suite: $TEST_SUITE"
    exit 1
    ;;
esac

echo "Tests complete!"
```

**Testing:**

```bash
chmod +x scripts/*.sh
./scripts/setup.sh
```

### Step 13: Update GitHub CI Workflow

**File:** `.github/workflows/ci.yml`

**Changes:**

```yaml
name: CI

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  lint:
    name: Lint
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.21'
      - name: golangci-lint
        uses: golangci/golangci-lint-action@v3

  test:
    name: Test
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.21'
      - name: Run unit tests
        run: go test -v -race -coverprofile=coverage.txt ./...
      - name: Upload coverage
        uses: codecov/codecov-action@v3

  build:
    name: Build
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.21'
      - name: Build CLI
        run: go build -v -o bin/ktest cmd/ktest/main.go
      - name: Upload artifact
        uses: actions/upload-artifact@v3
        with:
          name: ktest
          path: bin/ktest

  integration:
    name: Integration Tests
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.21'
      - name: Setup kind cluster
        uses: helm/kind-action@v1
      - name: Run integration tests
        run: go test -v -tags=integration ./tests/integration/...
```

**Testing:**

```bash
# Validate workflow syntax
actionlint .github/workflows/ci.yml
```

### Step 14: Create Documentation

**File:** `docs/usage.md`

**Changes:**

```markdown
# Usage Guide

## Installation

```bash
# Clone repository
git clone https://github.com/denhamparry/kubernetes-testing.git
cd kubernetes-testing

# Run setup
./scripts/setup.sh
```

## Running Tests

### Conformance Tests

```bash
./bin/ktest conformance --kubeconfig ~/.kube/config --mode quick
```

### Operational Tests

```bash
# All operational tests
./bin/ktest operational --kubeconfig ~/.kube/config

# Specific tests
./bin/ktest operational --tests networking,storage
```

### Performance Tests

```bash
./bin/ktest performance --kubeconfig ~/.kube/config --duration 10m --rps 200
```

## Configuration

Tests can be configured via YAML files in `configs/tests/`.

## Reports

Test reports are generated in multiple formats:

- HTML: `reports/test-report.html`
- JSON: `reports/test-report.json`
- Console output

**File:** `docs/architecture.md`

**Changes:**

```markdown
# Architecture

## Overview

The Kubernetes testing framework consists of:

1. **CLI Layer** - Cobra-based command interface
2. **Test Execution Layer** - Test orchestration and execution
3. **Kubernetes Client Layer** - Cluster interaction via client-go
4. **Reporting Layer** - Result aggregation and formatting

## Component Diagram

```text
┌─────────────────────────────────────┐
│         CLI (cmd/ktest)             │
├─────────────────────────────────────┤
│  conformance | operational | perf   │
└─────────────────┬───────────────────┘
                  │
┌─────────────────▼───────────────────┐
│       Test Execution Layer          │
├─────────────────────────────────────┤
│ • Conformance (Sonobuoy)            │
│ • Networking Tests                  │
│ • Storage Tests                     │
│ • Workload Tests                    │
│ • Performance Tests                 │
└─────────────────┬───────────────────┘
                  │
┌─────────────────▼───────────────────┐
│    Kubernetes Client (client-go)    │
└─────────────────┬───────────────────┘
                  │
┌─────────────────▼───────────────────┐
│      Kubernetes Cluster             │
└─────────────────────────────────────┘
```

## Test Flow

1. User executes CLI command
2. CLI parses flags and loads configuration
3. Kubeconfig is loaded and client is created
4. Test suite is executed
5. Results are collected
6. Report is generated

**Testing:**

```bash
# Verify docs render correctly
mdl docs/*.md
```

### Step 15: Update CLAUDE.md and README.md

**File:** `CLAUDE.md`

**Changes:**

Update the Quick Commands section:

```markdown
## Quick Commands

```bash
# Setup
./scripts/setup.sh

# Build
go build -o bin/ktest cmd/ktest/main.go

# Test
go test -v ./...

# Lint
golangci-lint run

# Run conformance tests
./bin/ktest conformance --kubeconfig ~/.kube/config

# Run operational tests
./bin/ktest operational --kubeconfig ~/.kube/config

# Run performance tests
./bin/ktest performance --kubeconfig ~/.kube/config
```

**File:** `README.md`

**Changes:**

Update to reflect Kubernetes testing focus:

```markdown
# Kubernetes Testing Framework

A comprehensive testing framework for Kubernetes clusters, supporting conformance testing, operational validation, and performance testing.

## Features

- ✅ Kubernetes conformance testing (via Sonobuoy)
- ✅ Operational testing (networking, storage, workloads)
- ✅ Performance testing and benchmarking
- ✅ Works with any Kubernetes cluster (kubeconfig-based)
- ✅ Detailed reporting (HTML, JSON, console)
- ✅ CI/CD integration ready

## Quick Start

```bash
# Setup
./scripts/setup.sh

# Run all tests
./scripts/run-tests.sh all

# Run specific test suite
./bin/ktest conformance --kubeconfig ~/.kube/config
```

## Documentation

- [Usage Guide](docs/usage.md)
- [Architecture](docs/architecture.md)
- [Contributing](CONTRIBUTING.md)

## Requirements

- Go 1.21+
- kubectl
- Access to a Kubernetes cluster

**Testing:**

```bash
# Validate markdown
mdl README.md CLAUDE.md
```

## Testing Strategy

### Unit Testing

**Approach:**

- Test all core functionality with mock Kubernetes clientsets
- Use `k8s.io/client-go/kubernetes/fake` for mocking
- Aim for >80% code coverage

**Example Tests:**

```go
// tests/unit/kubeconfig_test.go
func TestNewClient(t *testing.T) {
    // Test valid kubeconfig
    // Test invalid kubeconfig
    // Test missing kubeconfig
}

// tests/unit/networking_test.go
func TestDNSValidation(t *testing.T) {
    // Test DNS resolution logic
}
```

**Run:**

```bash
go test -v -race -coverprofile=coverage.txt ./...
```

### Integration Testing

**Test Case 1: Conformance Testing**

1. Setup: Create kind cluster
2. Execute: Run `ktest conformance --mode quick`
3. Verify: Check test results are generated
4. Verify: Ensure no errors in execution
5. Cleanup: Delete kind cluster

**Test Case 2: Operational Testing - Networking**

1. Setup: Create kind cluster
2. Execute: Run `ktest operational --tests networking`
3. Verify: DNS test passes
4. Verify: Pod-to-pod connectivity test passes
5. Verify: Service connectivity test passes
6. Cleanup: Delete kind cluster and test resources

**Test Case 3: Operational Testing - Storage**

1. Setup: Create kind cluster with storage class
2. Execute: Run `ktest operational --tests storage`
3. Verify: PVC creation test passes
4. Verify: Dynamic provisioning test passes
5. Verify: Test resources are cleaned up
6. Cleanup: Delete kind cluster

**Test Case 4: Workload Testing**

1. Setup: Create kind cluster
2. Execute: Run `ktest operational --tests workload`
3. Verify: Deployment test passes (creation, scaling, updates)
4. Verify: StatefulSet test passes
5. Verify: All test workloads are removed
6. Cleanup: Delete kind cluster

**Test Case 5: Performance Testing**

1. Setup: Create kind cluster with sample application
2. Execute: Run `ktest performance --duration 1m --rps 10`
3. Verify: Performance metrics are collected
4. Verify: Report is generated
5. Cleanup: Delete kind cluster

**Run:**

```bash
# Setup kind cluster
kind create cluster --name test-cluster

# Run integration tests
go test -v -tags=integration ./tests/integration/...

# Cleanup
kind delete cluster --name test-cluster
```

### Regression Testing

**Existing Functionality to Verify:**

- Template structure (CLAUDE.md, README.md, docs/)
- Pre-commit hooks still functional
- GitHub workflows validate
- Custom slash commands work

**Test:**

```bash
# Verify pre-commit hooks
pre-commit run --all-files

# Validate workflows
actionlint .github/workflows/*.yml

# Check template files
ls -la CLAUDE.md README.md CONTRIBUTING.md
```

### End-to-End Testing

**Scenario:** Complete test suite execution

1. Clone repository
2. Run `./scripts/setup.sh`
3. Create test cluster: `kind create cluster`
4. Run: `./scripts/run-tests.sh all`
5. Verify: All test suites complete successfully
6. Verify: Reports are generated
7. Cleanup: `kind delete cluster`

## Success Criteria

- [ ] Go module initialized with required dependencies
- [ ] CLI tool (`ktest`) builds successfully
- [ ] Kubeconfig loading and client creation works
- [ ] Conformance testing via Sonobuoy integration implemented
- [ ] Networking tests implemented (DNS, pod-to-pod, service)
- [ ] Storage tests implemented (PVC, storage classes)
- [ ] Workload tests implemented (deployment, statefulset, daemonset)
- [ ] Performance testing framework implemented
- [ ] CLI commands for all test suites working
- [ ] Test configuration files created
- [ ] Reporting functionality implemented (HTML, JSON, console)
- [ ] Setup and run scripts created and tested
- [ ] GitHub CI workflow updated and working
- [ ] Documentation complete (usage.md, architecture.md)
- [ ] CLAUDE.md and README.md updated
- [ ] Unit tests written with >80% coverage
- [ ] Integration tests passing on kind cluster
- [ ] Pre-commit hooks configured for Go
- [ ] All tests pass in CI/CD pipeline

## Files Modified

1. `go.mod` - Go module initialization
2. `go.sum` - Dependency lock file
3. `cmd/ktest/main.go` - CLI entry point
4. `pkg/cmd/root.go` - Root command
5. `pkg/cmd/conformance.go` - Conformance command
6. `pkg/cmd/operational.go` - Operational command
7. `pkg/cmd/performance.go` - Performance command
8. `pkg/kubeconfig/kubeconfig.go` - Kubeconfig handling
9. `pkg/conformance/sonobuoy.go` - Sonobuoy integration
10. `pkg/networking/dns.go` - DNS tests
11. `pkg/networking/connectivity.go` - Connectivity tests
12. `pkg/storage/pvc.go` - Storage tests
13. `pkg/workload/deployment.go` - Deployment tests
14. `pkg/workload/statefulset.go` - StatefulSet tests
15. `pkg/performance/load.go` - Load testing
16. `pkg/performance/metrics.go` - Performance metrics
17. `pkg/report/report.go` - Reporting
18. `configs/sonobuoy/quick-config.yaml` - Sonobuoy config
19. `configs/tests/networking.yaml` - Network test config
20. `configs/tests/storage.yaml` - Storage test config
21. `configs/tests/workload.yaml` - Workload test config
22. `scripts/setup.sh` - Setup script
23. `scripts/run-tests.sh` - Test runner script
24. `.github/workflows/ci.yml` - CI workflow
25. `docs/usage.md` - Usage documentation
26. `docs/architecture.md` - Architecture documentation
27. `CLAUDE.md` - Updated with Go commands
28. `README.md` - Updated with project description
29. `.pre-commit-config.yaml` - Add Go hooks
30. `tests/unit/*_test.go` - Unit tests
31. `tests/integration/*_test.go` - Integration tests

## Related Issues and Tasks

### Depends On

- None (this is the initial implementation)

### Blocks

- Future enhancements (custom test plugins, advanced reporting)
- Integration with CD pipelines
- Multi-cluster testing support

### Related

- GitHub Issue #1: Create kubernetes testing strategy

### Enables

- Automated Kubernetes cluster validation
- CI/CD integration for cluster testing
- Conformance certification testing
- Performance benchmarking
- Operational readiness validation

## References

- [GitHub Issue #1](https://github.com/denhamparry/kubernetes-testing/issues/1)
- [Sonobuoy Documentation](https://sonobuoy.io/)
- [Kubernetes client-go](https://github.com/kubernetes/client-go)
- [Cobra CLI Framework](https://github.com/spf13/cobra)
- [CNCF Conformance Testing](https://github.com/cncf/k8s-conformance)
- [kind - Kubernetes in Docker](https://kind.sigs.k8s.io/)

## Notes

### Key Insights

1. **Sonobuoy Integration** - Official CNCF tool provides standardized conformance testing
2. **Go Ecosystem** - Using Go aligns with Kubernetes ecosystem and provides excellent tooling
3. **Modular Design** - Separate packages for each test category enables extensibility
4. **kubeconfig-based** - Works with any Kubernetes cluster accessible via kubeconfig
5. **TDD Approach** - Following template's TDD principles with comprehensive test coverage
6. **CI/CD Ready** - GitHub Actions workflow with kind cluster for automated testing

### Alternative Approaches Considered

1. **Python-based Framework** ❌
   - Easier scripting but less aligned with Kubernetes ecosystem
   - Less performant for CLI tools
   - Dependency management complexity

2. **Shell Scripts Only** ❌
   - Simple but not maintainable or extensible
   - Difficult to test
   - Poor error handling

3. **Go with Custom E2E Tests** ✅ (Chosen)
   - Native Kubernetes language
   - Excellent client libraries
   - Strong typing and testing
   - Sonobuoy integration for conformance

4. **Helm Tests + Custom Scripts** ❌
   - Limited to Helm chart testing
   - Not comprehensive enough
   - Doesn't cover conformance testing

### Best Practices

1. **Error Handling** - Comprehensive error handling with meaningful messages
2. **Resource Cleanup** - Always clean up test resources (defer cleanup in tests)
3. **Configuration** - Use YAML for test configuration (version controlled)
4. **Reporting** - Multiple report formats for different use cases
5. **Logging** - Structured logging for debugging
6. **Timeouts** - Configurable timeouts for all operations
7. **Idempotency** - Tests can be run multiple times safely
8. **Documentation** - Comprehensive docs for usage and architecture
9. **CI Integration** - Automated testing on every PR
10. **Versioning** - Semantic versioning for releases

### Monitoring and Observability

- Test execution metrics (duration, pass/fail rates)
- Resource usage during tests
- Performance test results trending
- Integration with monitoring systems (future enhancement)

### Security Considerations

- Kubeconfig security (never commit kubeconfigs)
- RBAC requirements for tests documented
- Secret handling in performance tests
- Test isolation (namespaces)

### Future Enhancements

1. Multi-cluster testing support
2. Custom plugin framework for extensibility
3. Web UI for test execution and results
4. Historical test result database
5. Slack/email notifications
6. Integration with ArgoCD/Flux for GitOps testing
7. Cost analysis for cloud resources during tests
8. Chaos engineering integration
