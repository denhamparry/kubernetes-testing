#!/bin/bash
set -e

KUBECONFIG="${KUBECONFIG:-$HOME/.kube/config}"
TEST_SUITE="${1:-all}"

echo "========================================="
echo "Kubernetes Testing Framework"
echo "========================================="
echo "Kubeconfig: $KUBECONFIG"
echo "Test suite: $TEST_SUITE"
echo ""

# Check if ktest binary exists
if [ ! -f "./bin/ktest" ]; then
    echo "ERROR: ktest binary not found. Run ./scripts/setup.sh first."
    exit 1
fi

case $TEST_SUITE in
  conformance)
    echo "Running conformance tests..."
    ./bin/ktest conformance --kubeconfig "$KUBECONFIG"
    ;;
  operational)
    echo "Running operational tests..."
    ./bin/ktest operational --kubeconfig "$KUBECONFIG"
    ;;
  performance)
    if [ -z "$2" ]; then
        echo "ERROR: Performance tests require an endpoint."
        echo "Usage: $0 performance <endpoint>"
        exit 1
    fi
    echo "Running performance tests..."
    ./bin/ktest performance --endpoint "$2" --kubeconfig "$KUBECONFIG"
    ;;
  all)
    echo "Running all test suites..."
    echo ""
    echo "1/2: Conformance tests..."
    ./bin/ktest conformance --kubeconfig "$KUBECONFIG" || true
    echo ""
    echo "2/2: Operational tests..."
    ./bin/ktest operational --kubeconfig "$KUBECONFIG" || true
    echo ""
    echo "Note: Performance tests require an endpoint and must be run separately:"
    echo "  $0 performance <endpoint>"
    ;;
  *)
    echo "ERROR: Unknown test suite: $TEST_SUITE"
    echo "Usage: $0 {conformance|operational|performance|all} [endpoint]"
    exit 1
    ;;
esac

echo ""
echo "========================================="
echo "Tests complete!"
echo "========================================="
