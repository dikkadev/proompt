#!/bin/bash
# =============================================================================
# API Tests - Full CRUD operations testing
# =============================================================================
# These tests verify that all API endpoints work correctly with proper
# request/response handling, validation, and error cases.
# =============================================================================

set -euo pipefail

API_URL="${PROOMPT_API_URL:-http://proompt-server:8080}"

# Global variables for test data
CREATED_PROMPT_ID=""
CREATED_SNIPPET_ID=""
CREATED_NOTE_ID=""

# Helper functions
make_request() {
    local method="$1"
    local endpoint="$2"
    local data="${3:-}"
    
    if [[ -n "$data" ]]; then
        curl -s -X "$method" \
            -H "Content-Type: application/json" \
            -d "$data" \
            "$API_URL$endpoint"
    else
        curl -s -X "$method" "$API_URL$endpoint"
    fi
}

make_request_with_code() {
    local method="$1"
    local endpoint="$2"
    local data="${3:-}"
    
    if [[ -n "$data" ]]; then
        curl -s -w "%{http_code}" \
            -X "$method" \
            -H "Content-Type: application/json" \
            -d "$data" \
            "$API_URL$endpoint"
    else
        curl -s -w "%{http_code}" -X "$method" "$API_URL$endpoint"
    fi
}

extract_http_code() {
    local response="$1"
    echo "${response: -3}"
}

extract_body() {
    local response="$1"
    echo "${response%???}"
}

# Prompt API Tests
test_create_prompt() {
    echo "Testing prompt creation..."
    
    local payload='{
        "title": "Test Prompt",
        "content": "This is a test prompt for {{user_name}}",
        "type": "user",
        "use_case": "testing",
        "temperature_suggestion": 0.7
    }'
    
    local response
    response=$(make_request_with_code "POST" "/api/prompts" "$payload")
    local http_code
    http_code=$(extract_http_code "$response")
    local body
    body=$(extract_body "$response")
    
    if [[ "$http_code" != "201" ]]; then
        echo "Create prompt failed: HTTP $http_code"
        echo "Response: $body"
        return 1
    fi
    
    # Extract and store the ID
    CREATED_PROMPT_ID=$(echo "$body" | jq -r '.id')
    if [[ "$CREATED_PROMPT_ID" == "null" || -z "$CREATED_PROMPT_ID" ]]; then
        echo "Created prompt missing ID"
        return 1
    fi
    
    # Verify response structure
    local title
    title=$(echo "$body" | jq -r '.title')
    if [[ "$title" != "Test Prompt" ]]; then
        echo "Created prompt has wrong title: $title"
        return 1
    fi
    
    echo "✓ Prompt created successfully (ID: $CREATED_PROMPT_ID)"
    return 0
}

test_get_prompt() {
    echo "Testing prompt retrieval..."
    
    if [[ -z "$CREATED_PROMPT_ID" ]]; then
        echo "No prompt ID available for testing"
        return 1
    fi
    
    local response
    response=$(make_request_with_code "GET" "/api/prompts/$CREATED_PROMPT_ID")
    local http_code
    http_code=$(extract_http_code "$response")
    local body
    body=$(extract_body "$response")
    
    if [[ "$http_code" != "200" ]]; then
        echo "Get prompt failed: HTTP $http_code"
        return 1
    fi
    
    # Verify the prompt data
    local id
    id=$(echo "$body" | jq -r '.id')
    if [[ "$id" != "$CREATED_PROMPT_ID" ]]; then
        echo "Retrieved prompt has wrong ID: $id"
        return 1
    fi
    
    echo "✓ Prompt retrieved successfully"
    return 0
}

test_update_prompt() {
    echo "Testing prompt update..."
    
    if [[ -z "$CREATED_PROMPT_ID" ]]; then
        echo "No prompt ID available for testing"
        return 1
    fi
    
    local payload='{
        "title": "Updated Test Prompt",
        "content": "This is an updated test prompt"
    }'
    
    local response
    response=$(make_request_with_code "PUT" "/api/prompts/$CREATED_PROMPT_ID" "$payload")
    local http_code
    http_code=$(extract_http_code "$response")
    local body
    body=$(extract_body "$response")
    
    if [[ "$http_code" != "200" ]]; then
        echo "Update prompt failed: HTTP $http_code"
        echo "Response: $body"
        return 1
    fi
    
    # Verify the update
    local title
    title=$(echo "$body" | jq -r '.title')
    if [[ "$title" != "Updated Test Prompt" ]]; then
        echo "Prompt title not updated: $title"
        return 1
    fi
    
    echo "✓ Prompt updated successfully"
    return 0
}

test_list_prompts() {
    echo "Testing prompt listing..."
    
    local response
    response=$(make_request_with_code "GET" "/api/prompts")
    local http_code
    http_code=$(extract_http_code "$response")
    local body
    body=$(extract_body "$response")
    
    if [[ "$http_code" != "200" ]]; then
        echo "List prompts failed: HTTP $http_code"
        return 1
    fi
    
    # Verify response structure
    if ! echo "$body" | jq -e '.data' >/dev/null 2>&1; then
        echo "List response missing 'data' field"
        return 1
    fi
    
    # Check if our created prompt is in the list
    local found
    found=$(echo "$body" | jq -r ".data[] | select(.id == \"$CREATED_PROMPT_ID\") | .id")
    if [[ "$found" != "$CREATED_PROMPT_ID" ]]; then
        echo "Created prompt not found in list"
        return 1
    fi
    
    echo "✓ Prompts listed successfully"
    return 0
}

test_create_snippet() {
    echo "Testing snippet creation..."
    
    local payload='{
        "title": "Test Snippet",
        "content": "echo \"Hello, World!\"",
        "language": "bash",
        "description": "A simple test snippet"
    }'
    
    local response
    response=$(make_request_with_code "POST" "/api/snippets" "$payload")
    local http_code
    http_code=$(extract_http_code "$response")
    local body
    body=$(extract_body "$response")
    
    if [[ "$http_code" != "201" ]]; then
        echo "Create snippet failed: HTTP $http_code"
        echo "Response: $body"
        return 1
    fi
    
    # Extract and store the ID
    CREATED_SNIPPET_ID=$(echo "$body" | jq -r '.id')
    if [[ "$CREATED_SNIPPET_ID" == "null" || -z "$CREATED_SNIPPET_ID" ]]; then
        echo "Created snippet missing ID"
        return 1
    fi
    
    echo "✓ Snippet created successfully (ID: $CREATED_SNIPPET_ID)"
    return 0
}

test_create_note() {
    echo "Testing note creation..."
    
    if [[ -z "$CREATED_PROMPT_ID" ]]; then
        echo "No prompt ID available for note creation"
        return 1
    fi
    
    local payload='{
        "title": "Test Note",
        "body": "This is a test note for the prompt"
    }'
    
    local response
    response=$(make_request_with_code "POST" "/api/prompts/$CREATED_PROMPT_ID/notes" "$payload")
    local http_code
    http_code=$(extract_http_code "$response")
    local body
    body=$(extract_body "$response")
    
    if [[ "$http_code" != "201" ]]; then
        echo "Create note failed: HTTP $http_code"
        echo "Response: $body"
        return 1
    fi
    
    # Extract and store the ID
    CREATED_NOTE_ID=$(echo "$body" | jq -r '.id')
    if [[ "$CREATED_NOTE_ID" == "null" || -z "$CREATED_NOTE_ID" ]]; then
        echo "Created note missing ID"
        return 1
    fi
    
    echo "✓ Note created successfully (ID: $CREATED_NOTE_ID)"
    return 0
}

test_validation_errors() {
    echo "Testing validation errors..."
    
    # Test missing required fields
    local payload='{
        "content": "Missing title"
    }'
    
    local response
    response=$(make_request_with_code "POST" "/api/prompts" "$payload")
    local http_code
    http_code=$(extract_http_code "$response")
    
    if [[ "$http_code" != "400" ]]; then
        echo "Expected 400 for missing title, got $http_code"
        return 1
    fi
    
    # Test invalid JSON
    response=$(curl -s -w "%{http_code}" \
        -X "POST" \
        -H "Content-Type: application/json" \
        -d "invalid json" \
        "$API_URL/api/prompts")
    http_code=$(extract_http_code "$response")
    
    if [[ "$http_code" != "400" ]]; then
        echo "Expected 400 for invalid JSON, got $http_code"
        return 1
    fi
    
    echo "✓ Validation errors handled correctly"
    return 0
}

test_not_found_errors() {
    echo "Testing not found errors..."
    
    local fake_id="00000000-0000-0000-0000-000000000000"
    
    local response
    response=$(make_request_with_code "GET" "/api/prompts/$fake_id")
    local http_code
    http_code=$(extract_http_code "$response")
    
    if [[ "$http_code" != "404" ]]; then
        echo "Expected 404 for non-existent prompt, got $http_code"
        return 1
    fi
    
    echo "✓ Not found errors handled correctly"
    return 0
}

# Cleanup function
cleanup_test_data() {
    echo "Cleaning up test data..."
    
    # Delete created note
    if [[ -n "$CREATED_NOTE_ID" ]]; then
        make_request "DELETE" "/api/notes/$CREATED_NOTE_ID" >/dev/null 2>&1 || true
    fi
    
    # Delete created snippet
    if [[ -n "$CREATED_SNIPPET_ID" ]]; then
        make_request "DELETE" "/api/snippets/$CREATED_SNIPPET_ID" >/dev/null 2>&1 || true
    fi
    
    # Delete created prompt
    if [[ -n "$CREATED_PROMPT_ID" ]]; then
        make_request "DELETE" "/api/prompts/$CREATED_PROMPT_ID" >/dev/null 2>&1 || true
    fi
    
    echo "✓ Test data cleaned up"
}

# Main test execution
main() {
    echo "Running API tests against $API_URL"
    echo "================================="
    
    # Set up cleanup trap
    trap cleanup_test_data EXIT
    
    # Run tests in order
    test_create_prompt
    test_get_prompt
    test_update_prompt
    test_list_prompts
    test_create_snippet
    test_create_note
    test_validation_errors
    test_not_found_errors
    
    echo ""
    echo "All API tests passed! ✓"
}

main "$@"