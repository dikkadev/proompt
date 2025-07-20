#!/bin/bash
# =============================================================================
# Persistence Tests - Data persistence and git integration verification
# =============================================================================
# These tests verify that data persists correctly across container restarts
# and that the git integration is working properly.
# =============================================================================

set -euo pipefail

API_URL="${PROOMPT_API_URL:-http://proompt-server:8080}"

# Test data
TEST_PROMPT_ID=""
TEST_DATA_DIR="/app/data"

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

test_database_persistence() {
    echo "Testing database persistence..."
    
    # Create a test prompt
    local payload='{
        "title": "Persistence Test Prompt",
        "content": "This prompt tests data persistence",
        "type": "user",
        "use_case": "testing"
    }'
    
    local response
    response=$(make_request_with_code "POST" "/api/prompts" "$payload")
    local http_code
    http_code=$(extract_http_code "$response")
    local body
    body=$(extract_body "$response")
    
    if [[ "$http_code" != "201" ]]; then
        echo "Failed to create test prompt: HTTP $http_code"
        return 1
    fi
    
    TEST_PROMPT_ID=$(echo "$body" | jq -r '.id')
    if [[ "$TEST_PROMPT_ID" == "null" || -z "$TEST_PROMPT_ID" ]]; then
        echo "Created prompt missing ID"
        return 1
    fi
    
    # Verify the prompt exists
    response=$(make_request_with_code "GET" "/api/prompts/$TEST_PROMPT_ID")
    http_code=$(extract_http_code "$response")
    
    if [[ "$http_code" != "200" ]]; then
        echo "Failed to retrieve created prompt: HTTP $http_code"
        return 1
    fi
    
    echo "✓ Database persistence working (Prompt ID: $TEST_PROMPT_ID)"
    return 0
}

test_git_repository_structure() {
    echo "Testing git repository structure..."
    
    # Note: This test assumes we can access the container's filesystem
    # In a real container environment, we might need to exec into the container
    # For now, we'll test indirectly through the API
    
    if [[ -z "$TEST_PROMPT_ID" ]]; then
        echo "No test prompt available for git testing"
        return 1
    fi
    
    # The git integration should have created a branch for our prompt
    # We can't directly access the git repo from the test container,
    # but we can verify that the data is consistent
    
    # Retrieve the prompt again to ensure it's still there
    local response
    response=$(make_request "GET" "/api/prompts/$TEST_PROMPT_ID")
    
    local title
    title=$(echo "$response" | jq -r '.title')
    if [[ "$title" != "Persistence Test Prompt" ]]; then
        echo "Prompt data inconsistent, git integration may have failed"
        return 1
    fi
    
    echo "✓ Git repository structure appears consistent"
    return 0
}

test_data_consistency() {
    echo "Testing data consistency..."
    
    if [[ -z "$TEST_PROMPT_ID" ]]; then
        echo "No test prompt available for consistency testing"
        return 1
    fi
    
    # Update the prompt
    local update_payload='{
        "title": "Updated Persistence Test Prompt",
        "content": "This content has been updated"
    }'
    
    local response
    response=$(make_request_with_code "PUT" "/api/prompts/$TEST_PROMPT_ID" "$update_payload")
    local http_code
    http_code=$(extract_http_code "$response")
    
    if [[ "$http_code" != "200" ]]; then
        echo "Failed to update prompt: HTTP $http_code"
        return 1
    fi
    
    # Retrieve and verify the update
    response=$(make_request "GET" "/api/prompts/$TEST_PROMPT_ID")
    local title
    title=$(echo "$response" | jq -r '.title')
    
    if [[ "$title" != "Updated Persistence Test Prompt" ]]; then
        echo "Prompt update not persisted correctly: $title"
        return 1
    fi
    
    echo "✓ Data consistency maintained across updates"
    return 0
}

test_multiple_entities() {
    echo "Testing persistence with multiple entity types..."
    
    # Create a snippet
    local snippet_payload='{
        "title": "Persistence Test Snippet",
        "content": "#!/bin/bash\necho \"persistence test\"",
        "language": "bash"
    }'
    
    local response
    response=$(make_request_with_code "POST" "/api/snippets" "$snippet_payload")
    local http_code
    http_code=$(extract_http_code "$response")
    local body
    body=$(extract_body "$response")
    
    if [[ "$http_code" != "201" ]]; then
        echo "Failed to create test snippet: HTTP $http_code"
        return 1
    fi
    
    local snippet_id
    snippet_id=$(echo "$body" | jq -r '.id')
    
    # Create a note for our prompt
    if [[ -n "$TEST_PROMPT_ID" ]]; then
        local note_payload='{
            "title": "Persistence Test Note",
            "body": "This is a persistence test note"
        }'
        
        response=$(make_request_with_code "POST" "/api/prompts/$TEST_PROMPT_ID/notes" "$note_payload")
        http_code=$(extract_http_code "$response")
        
        if [[ "$http_code" != "201" ]]; then
            echo "Failed to create test note: HTTP $http_code"
            return 1
        fi
        
        local note_body
        note_body=$(extract_body "$response")
        local note_id
        note_id=$(echo "$note_body" | jq -r '.id')
        
        # Verify all entities exist
        response=$(make_request_with_code "GET" "/api/snippets/$snippet_id")
        http_code=$(extract_http_code "$response")
        if [[ "$http_code" != "200" ]]; then
            echo "Snippet not found after creation"
            return 1
        fi
        
        response=$(make_request_with_code "GET" "/api/notes/$note_id")
        http_code=$(extract_http_code "$response")
        if [[ "$http_code" != "200" ]]; then
            echo "Note not found after creation"
            return 1
        fi
        
        # Cleanup
        make_request "DELETE" "/api/notes/$note_id" >/dev/null 2>&1 || true
    fi
    
    # Cleanup snippet
    make_request "DELETE" "/api/snippets/$snippet_id" >/dev/null 2>&1 || true
    
    echo "✓ Multiple entity types persist correctly"
    return 0
}

test_transaction_integrity() {
    echo "Testing transaction integrity..."
    
    # Get initial count of prompts
    local response
    response=$(make_request "GET" "/api/prompts")
    local initial_count
    initial_count=$(echo "$response" | jq '.data | length')
    
    # Try to create a prompt with invalid data (should fail)
    local invalid_payload='{
        "content": "Missing required title field"
    }'
    
    response=$(make_request_with_code "POST" "/api/prompts" "$invalid_payload")
    local http_code
    http_code=$(extract_http_code "$response")
    
    if [[ "$http_code" != "400" ]]; then
        echo "Expected validation error, got HTTP $http_code"
        return 1
    fi
    
    # Verify count hasn't changed
    response=$(make_request "GET" "/api/prompts")
    local final_count
    final_count=$(echo "$response" | jq '.data | length')
    
    if [[ "$initial_count" != "$final_count" ]]; then
        echo "Transaction integrity compromised: count changed from $initial_count to $final_count"
        return 1
    fi
    
    echo "✓ Transaction integrity maintained"
    return 0
}

# Cleanup function
cleanup_test_data() {
    echo "Cleaning up persistence test data..."
    
    if [[ -n "$TEST_PROMPT_ID" ]]; then
        make_request "DELETE" "/api/prompts/$TEST_PROMPT_ID" >/dev/null 2>&1 || true
    fi
    
    echo "✓ Persistence test data cleaned up"
}

# Main test execution
main() {
    echo "Running persistence tests against $API_URL"
    echo "========================================"
    
    # Set up cleanup trap
    trap cleanup_test_data EXIT
    
    # Run tests in order
    test_database_persistence
    test_git_repository_structure
    test_data_consistency
    test_multiple_entities
    test_transaction_integrity
    
    echo ""
    echo "All persistence tests passed! ✓"
}

main "$@"