#!/bin/bash
# =============================================================================
# Smoke Tests - Basic functionality verification
# =============================================================================
# These tests verify that the basic server functionality is working.
# =============================================================================

set -euo pipefail

API_URL="${PROOMPT_API_URL:-http://proompt-server:8080}"

# Test functions
test_health_endpoint() {
    echo "Testing health endpoint..."
    
    local response
    response=$(curl -s -w "%{http_code}" "$API_URL/api/health")
    local http_code="${response: -3}"
    local body="${response%???}"
    
    if [[ "$http_code" != "200" ]]; then
        echo "Health check failed: HTTP $http_code"
        return 1
    fi
    
    # Verify JSON response structure
    if ! echo "$body" | jq -e '.status' >/dev/null 2>&1; then
        echo "Health response missing 'status' field"
        return 1
    fi
    
    local status
    status=$(echo "$body" | jq -r '.status')
    if [[ "$status" != "healthy" ]]; then
        echo "Health status is not 'healthy': $status"
        return 1
    fi
    
    echo "✓ Health endpoint working"
    return 0
}

test_server_responds() {
    echo "Testing server responsiveness..."
    
    if ! curl -s -f "$API_URL/api/health" >/dev/null; then
        echo "Server is not responding"
        return 1
    fi
    
    echo "✓ Server is responding"
    return 0
}

test_cors_headers() {
    echo "Testing CORS headers..."
    
    local headers
    headers=$(curl -s -I "$API_URL/api/health")
    
    if ! echo "$headers" | grep -i "access-control-allow-origin" >/dev/null; then
        echo "Warning: CORS headers not found (may be intentional)"
    else
        echo "✓ CORS headers present"
    fi
    
    return 0
}

test_content_type() {
    echo "Testing content type headers..."
    
    local content_type
    content_type=$(curl -s -I "$API_URL/api/health" | grep -i "content-type" | cut -d' ' -f2- | tr -d '\r\n')
    
    if [[ "$content_type" != *"application/json"* ]]; then
        echo "Unexpected content type: $content_type"
        return 1
    fi
    
    echo "✓ Content type is JSON"
    return 0
}

test_invalid_endpoint() {
    echo "Testing invalid endpoint handling..."
    
    local http_code
    http_code=$(curl -s -w "%{http_code}" -o /dev/null "$API_URL/api/nonexistent")
    
    if [[ "$http_code" != "404" ]]; then
        echo "Expected 404 for invalid endpoint, got $http_code"
        return 1
    fi
    
    echo "✓ Invalid endpoints return 404"
    return 0
}

# Run all smoke tests
main() {
    echo "Running smoke tests against $API_URL"
    echo "=================================="
    
    test_server_responds
    test_health_endpoint
    test_content_type
    test_cors_headers
    test_invalid_endpoint
    
    echo ""
    echo "All smoke tests passed! ✓"
}

main "$@"