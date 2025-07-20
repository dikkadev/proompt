#!/bin/bash
# =============================================================================
# Test Environment Cleanup Script
# =============================================================================
# This script cleans up all test-related Docker resources.
# =============================================================================

set -euo pipefail

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

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TESTS_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

# Cleanup functions
cleanup_containers() {
    log "Stopping and removing test containers..."
    
    cd "$TESTS_DIR"
    
    # Stop and remove containers from test compose
    if docker compose -f docker/compose.test.yaml down -v >/dev/null 2>&1; then
        success "Test containers removed"
    else
        warning "No test containers to remove or compose file not found"
    fi
    
    # Also try to clean up any running proompt containers
    local containers
    containers=$(docker ps -a --filter "name=proompt" --format "{{.Names}}" 2>/dev/null || true)
    
    if [[ -n "$containers" ]]; then
        log "Found additional proompt containers: $containers"
        echo "$containers" | xargs -r docker rm -f >/dev/null 2>&1 || true
        success "Additional containers removed"
    fi
}

cleanup_images() {
    log "Removing test images..."
    
    local images_to_remove=(
        "proompt:test"
        "proompt:test-runner"
        "proompt:latest"
    )
    
    for image in "${images_to_remove[@]}"; do
        if docker image inspect "$image" >/dev/null 2>&1; then
            if docker rmi "$image" >/dev/null 2>&1; then
                success "Removed image: $image"
            else
                warning "Failed to remove image: $image (may be in use)"
            fi
        fi
    done
}

cleanup_volumes() {
    log "Removing test volumes..."
    
    local volumes_to_remove=(
        "proompt_test_data"
        "proompt_test_results"
        "proompt_data"
    )
    
    for volume in "${volumes_to_remove[@]}"; do
        if docker volume inspect "$volume" >/dev/null 2>&1; then
            if docker volume rm "$volume" >/dev/null 2>&1; then
                success "Removed volume: $volume"
            else
                warning "Failed to remove volume: $volume (may be in use)"
            fi
        fi
    done
}

cleanup_networks() {
    log "Removing test networks..."
    
    local networks_to_remove=(
        "proompt_test_net"
        "proompt_net"
    )
    
    for network in "${networks_to_remove[@]}"; do
        if docker network inspect "$network" >/dev/null 2>&1; then
            if docker network rm "$network" >/dev/null 2>&1; then
                success "Removed network: $network"
            else
                warning "Failed to remove network: $network (may be in use)"
            fi
        fi
    done
}

cleanup_dangling_resources() {
    log "Cleaning up dangling Docker resources..."
    
    # Remove dangling images
    local dangling_images
    dangling_images=$(docker images -f "dangling=true" -q 2>/dev/null || true)
    if [[ -n "$dangling_images" ]]; then
        echo "$dangling_images" | xargs -r docker rmi >/dev/null 2>&1 || true
        success "Removed dangling images"
    fi
    
    # Remove unused volumes
    if docker volume prune -f >/dev/null 2>&1; then
        success "Removed unused volumes"
    fi
    
    # Remove unused networks
    if docker network prune -f >/dev/null 2>&1; then
        success "Removed unused networks"
    fi
}

show_remaining_resources() {
    log "Checking for remaining proompt-related resources..."
    
    echo ""
    echo "Remaining containers:"
    docker ps -a --filter "name=proompt" --format "table {{.Names}}\t{{.Status}}\t{{.Image}}" 2>/dev/null || echo "None"
    
    echo ""
    echo "Remaining images:"
    docker images --filter "reference=proompt*" --format "table {{.Repository}}\t{{.Tag}}\t{{.Size}}" 2>/dev/null || echo "None"
    
    echo ""
    echo "Remaining volumes:"
    docker volume ls --filter "name=proompt" --format "table {{.Name}}\t{{.Driver}}" 2>/dev/null || echo "None"
    
    echo ""
    echo "Remaining networks:"
    docker network ls --filter "name=proompt" --format "table {{.Name}}\t{{.Driver}}" 2>/dev/null || echo "None"
}

# Print usage
usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  -h, --help          Show this help message"
    echo "  -a, --all           Clean up everything (containers, images, volumes, networks)"
    echo "  -c, --containers    Clean up only containers"
    echo "  -i, --images        Clean up only images"
    echo "  -v, --volumes       Clean up only volumes"
    echo "  -n, --networks      Clean up only networks"
    echo "  -d, --dangling      Clean up only dangling resources"
    echo "  -s, --show          Show remaining resources without cleaning"
    echo ""
    echo "Examples:"
    echo "  $0                  Clean up containers and volumes (default)"
    echo "  $0 -a               Clean up everything"
    echo "  $0 -i               Clean up only images"
    echo "  $0 -s               Show what resources exist"
}

# Parse command line arguments
CLEANUP_CONTAINERS=0
CLEANUP_IMAGES=0
CLEANUP_VOLUMES=0
CLEANUP_NETWORKS=0
CLEANUP_DANGLING=0
SHOW_ONLY=0

# Default behavior if no specific options
if [[ $# -eq 0 ]]; then
    CLEANUP_CONTAINERS=1
    CLEANUP_VOLUMES=1
fi

while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            usage
            exit 0
            ;;
        -a|--all)
            CLEANUP_CONTAINERS=1
            CLEANUP_IMAGES=1
            CLEANUP_VOLUMES=1
            CLEANUP_NETWORKS=1
            CLEANUP_DANGLING=1
            shift
            ;;
        -c|--containers)
            CLEANUP_CONTAINERS=1
            shift
            ;;
        -i|--images)
            CLEANUP_IMAGES=1
            shift
            ;;
        -v|--volumes)
            CLEANUP_VOLUMES=1
            shift
            ;;
        -n|--networks)
            CLEANUP_NETWORKS=1
            shift
            ;;
        -d|--dangling)
            CLEANUP_DANGLING=1
            shift
            ;;
        -s|--show)
            SHOW_ONLY=1
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
    log "Starting Proompt test environment cleanup"
    echo ""
    
    if [[ $SHOW_ONLY -eq 1 ]]; then
        show_remaining_resources
        exit 0
    fi
    
    # Run cleanup functions based on options
    if [[ $CLEANUP_CONTAINERS -eq 1 ]]; then
        cleanup_containers
    fi
    
    if [[ $CLEANUP_IMAGES -eq 1 ]]; then
        cleanup_images
    fi
    
    if [[ $CLEANUP_VOLUMES -eq 1 ]]; then
        cleanup_volumes
    fi
    
    if [[ $CLEANUP_NETWORKS -eq 1 ]]; then
        cleanup_networks
    fi
    
    if [[ $CLEANUP_DANGLING -eq 1 ]]; then
        cleanup_dangling_resources
    fi
    
    echo ""
    success "Cleanup completed!"
    
    # Show what's left
    show_remaining_resources
}

# Run main function
main "$@"