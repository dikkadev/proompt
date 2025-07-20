#!/bin/bash
# =============================================================================
# Test Environment Setup Script
# =============================================================================
# This script sets up the test environment and runs the integration tests.
# =============================================================================

set -euo pipefail

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
TESTS_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

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

# Check prerequisites
check_prerequisites() {
    log "Checking prerequisites..."
    
    if ! command -v docker >/dev/null 2>&1; then
        error "Docker is not installed or not in PATH"
        exit 1
    fi
    
    if ! command -v docker compose >/dev/null 2>&1; then
        error "Docker Compose is not installed or not in PATH"
        exit 1
    fi
    
    success "Prerequisites check passed"
}

# Build test images
build_test_images() {
    log "Building test images..."
    
    cd "$TESTS_DIR"
    
    # Build the test runner image
    if ! docker build -f docker/Dockerfile.test -t proompt:test-runner .; then
        error "Failed to build test runner image"
        exit 1
    fi
    
    success "Test images built successfully"
}

# Run integration tests
run_integration_tests() {
    log "Running integration tests..."
    
    cd "$TESTS_DIR"
    
    # Run the test suite
    if docker compose -f docker/compose.test.yaml up --build --abort-on-container-exit; then
        success "Integration tests completed successfully"
        return 0
    else
        error "Integration tests failed"
        return 1
    fi
}

# Cleanup test environment
cleanup() {
    log "Cleaning up test environment..."
    
    cd "$TESTS_DIR"
    
    # Stop and remove containers
    docker compose -f docker/compose.test.yaml down -v >/dev/null 2>&1 || true
    
    # Remove test images (optional)
    if [[ "${CLEANUP_IMAGES:-0}" == "1" ]]; then
        docker rmi proompt:test proompt:test-runner >/dev/null 2>&1 || true
    fi
    
    success "Cleanup completed"
}

# Show test results
show_results() {
    log "Retrieving test results..."
    
    # Try to get test results from the volume
    local results_volume="proompt_test_results"
    
    if docker volume inspect "$results_volume" >/dev/null 2>&1; then
        # Create a temporary container to extract results
        if docker run --rm -v "$results_volume:/results" alpine:3.21 cat /results/test_report.txt 2>/dev/null; then
            echo ""
        else
            warning "Could not retrieve detailed test results"
        fi
    else
        warning "Test results volume not found"
    fi
}

# Print usage
usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  -h, --help          Show this help message"
    echo "  -c, --cleanup       Clean up test environment after running"
    echo "  -v, --verbose       Enable verbose output"
    echo "  --cleanup-images    Remove test images after cleanup"
    echo "  --build-only        Only build test images, don't run tests"
    echo "  --run-only          Only run tests, don't build images"
    echo ""
    echo "Examples:"
    echo "  $0                  Run full test suite"
    echo "  $0 -c               Run tests and cleanup afterwards"
    echo "  $0 --build-only     Only build test images"
    echo "  $0 --run-only       Only run tests (assumes images exist)"
}

# Parse command line arguments
CLEANUP_AFTER=0
VERBOSE=0
BUILD_ONLY=0
RUN_ONLY=0

while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            usage
            exit 0
            ;;
        -c|--cleanup)
            CLEANUP_AFTER=1
            shift
            ;;
        -v|--verbose)
            VERBOSE=1
            export VERBOSE=1
            shift
            ;;
        --cleanup-images)
            export CLEANUP_IMAGES=1
            shift
            ;;
        --build-only)
            BUILD_ONLY=1
            shift
            ;;
        --run-only)
            RUN_ONLY=1
            shift
            ;;
        *)
            error "Unknown option: $1"
            usage
            exit 1
            ;;
    esac
done

# Main execution
main() {
    log "Starting Proompt integration test setup"
    log "Project root: $PROJECT_ROOT"
    log "Tests directory: $TESTS_DIR"
    echo ""
    
    # Set up cleanup trap
    if [[ $CLEANUP_AFTER -eq 1 ]]; then
        trap cleanup EXIT
    fi
    
    # Check prerequisites
    check_prerequisites
    
    # Build images if needed
    if [[ $RUN_ONLY -eq 0 ]]; then
        build_test_images
    fi
    
    # Exit early if build-only
    if [[ $BUILD_ONLY -eq 1 ]]; then
        success "Build completed successfully"
        exit 0
    fi
    
    # Run tests
    if run_integration_tests; then
        success "All tests passed! 🎉"
        show_results
        exit 0
    else
        error "Tests failed! 😞"
        show_results
        exit 1
    fi
}

# Run main function
main "$@"