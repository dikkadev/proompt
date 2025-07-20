#!/bin/bash
# =============================================================================
# Proompt Integration Test Suite
# =============================================================================
# This script runs all integration tests for the Proompt server.
# It's designed to run inside the test-runner container.
# =============================================================================

set -euo pipefail

# Configuration
API_URL="${PROOMPT_API_URL:-http://proompt-server:8080}"
TEST_TIMEOUT="${TEST_TIMEOUT:-60}"
VERBOSE="${VERBOSE:-0}"
RESULTS_DIR="/results"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging functions
log() {
    echo -e "${BLUE}[$(date +'%Y-%m-%d %H:%M:%S')]${NC} $*"
}

success() {
    echo -e "${GREEN}✓${NC} $*"
}

error() {
    echo -e "${RED}✗${NC} $*" >&2
}

warning() {
    echo -e "${YELLOW}⚠${NC} $*"
}

# Test result tracking
TESTS_RUN=0
TESTS_PASSED=0
TESTS_FAILED=0
FAILED_TESTS=()

# Function to run a test and track results
run_test() {
    local test_name="$1"
    local test_script="$2"
    
    log "Running test: $test_name"
    TESTS_RUN=$((TESTS_RUN + 1))
    
    if [[ $VERBOSE -eq 1 ]]; then
        if bash "$test_script"; then
            success "$test_name passed"
            TESTS_PASSED=$((TESTS_PASSED + 1))
        else
            error "$test_name failed"
            TESTS_FAILED=$((TESTS_FAILED + 1))
            FAILED_TESTS+=("$test_name")
        fi
    else
        if bash "$test_script" >/dev/null 2>&1; then
            success "$test_name passed"
            TESTS_PASSED=$((TESTS_PASSED + 1))
        else
            error "$test_name failed"
            TESTS_FAILED=$((TESTS_FAILED + 1))
            FAILED_TESTS+=("$test_name")
        fi
    fi
}

# Wait for server to be ready
wait_for_server() {
    log "Waiting for server to be ready at $API_URL..."
    local max_attempts=30
    local attempt=1
    
    while [[ $attempt -le $max_attempts ]]; do
        if curl -s -f "$API_URL/api/health" >/dev/null 2>&1; then
            success "Server is ready!"
            return 0
        fi
        
        if [[ $attempt -eq $max_attempts ]]; then
            error "Server failed to start within $max_attempts attempts"
            return 1
        fi
        
        log "Attempt $attempt/$max_attempts - waiting 2 seconds..."
        sleep 2
        attempt=$((attempt + 1))
    done
}

# Generate test report
generate_report() {
    local report_file="$RESULTS_DIR/test_report.txt"
    
    {
        echo "Proompt Integration Test Report"
        echo "=============================="
        echo "Date: $(date)"
        echo "API URL: $API_URL"
        echo ""
        echo "Summary:"
        echo "  Tests Run: $TESTS_RUN"
        echo "  Passed: $TESTS_PASSED"
        echo "  Failed: $TESTS_FAILED"
        echo ""
        
        if [[ $TESTS_FAILED -gt 0 ]]; then
            echo "Failed Tests:"
            for test in "${FAILED_TESTS[@]}"; do
                echo "  - $test"
            done
            echo ""
        fi
        
        if [[ $TESTS_FAILED -eq 0 ]]; then
            echo "Result: ALL TESTS PASSED ✓"
        else
            echo "Result: SOME TESTS FAILED ✗"
        fi
    } > "$report_file"
    
    log "Test report saved to $report_file"
}

# Main test execution
main() {
    log "Starting Proompt Integration Test Suite"
    log "API URL: $API_URL"
    log "Test Timeout: ${TEST_TIMEOUT}s"
    log "Verbose: $([[ $VERBOSE -eq 1 ]] && echo "enabled" || echo "disabled")"
    echo ""
    
    # Wait for server
    if ! wait_for_server; then
        error "Server is not ready, aborting tests"
        exit 1
    fi
    
    echo ""
    log "Running test suite..."
    echo ""
    
    # Run tests in order
    run_test "Smoke Tests" "/tests/smoke_tests.sh"
    run_test "API Tests" "/tests/api_tests.sh"
    run_test "Persistence Tests" "/tests/persistence_tests.sh"
    
    echo ""
    log "Test suite completed"
    echo ""
    
    # Print summary
    echo "Test Summary:"
    echo "============="
    echo "Tests Run: $TESTS_RUN"
    success "Passed: $TESTS_PASSED"
    if [[ $TESTS_FAILED -gt 0 ]]; then
        error "Failed: $TESTS_FAILED"
    else
        echo "Failed: $TESTS_FAILED"
    fi
    echo ""
    
    # Generate report
    generate_report
    
    # Exit with appropriate code
    if [[ $TESTS_FAILED -eq 0 ]]; then
        success "All tests passed! 🎉"
        exit 0
    else
        error "Some tests failed! 😞"
        exit 1
    fi
}

# Run main function
main "$@"