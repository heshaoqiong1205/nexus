#!/bin/bash

# Nexus IoT Platform HTTP API Test Script
# This script tests all HTTP endpoints with proper authentication

set -e

# Configuration
BASE_URL="http://127.0.0.1:7070"
DEVICE_ID="dev-002"
DEVICE_SECRET="device-secret-002"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Helper functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

# Generate timestamp
get_timestamp() {
    date -u +"%Y-%m-%dT%H:%M:%SZ"
}

# Generate digest for device authentication
generate_digest() {
    local body="$1"
    local secret="$2"
    echo -n "${body}${secret}" | shasum -a 256 | cut -d' ' -f1
}

# Test API endpoint with device authentication
test_device_endpoint() {
    local method="$1"
    local endpoint="$2"
    local body="$3"
    local description="$4"
    
    log_info "Testing: $description"
    
    local timestamp=$(get_timestamp)
    local digest=$(generate_digest "$body" "$DEVICE_SECRET")
    
    local response
    if [ "$method" = "GET" ]; then
        response=$(curl -s -w "HTTPSTATUS:%{http_code}" \
            -H "Content-Type: application/json" \
            -H "device-id: $DEVICE_ID" \
            -H "time: $timestamp" \
            -H "digest: $digest" \
            "$BASE_URL$endpoint")
    else
        response=$(curl -s -w "HTTPSTATUS:%{http_code}" \
            -X "$method" \
            -H "Content-Type: application/json" \
            -H "device-id: $DEVICE_ID" \
            -H "time: $timestamp" \
            -H "digest: $digest" \
            -d "$body" \
            "$BASE_URL$endpoint")
    fi
    
    local http_code=$(echo "$response" | grep -o "HTTPSTATUS:[0-9]*" | cut -d: -f2)
    local body_response=$(echo "$response" | sed -E 's/HTTPSTATUS:[0-9]*$//')
    
    echo "  HTTP Status: $http_code"
    echo "  Response: $body_response"
    
    if [ "$http_code" -ge 200 ] && [ "$http_code" -lt 300 ]; then
        log_success "$description - PASSED"
    else
        log_error "$description - FAILED (HTTP $http_code)"
    fi
    echo
}

# Test API endpoint without authentication
test_public_endpoint() {
    local method="$1"
    local endpoint="$2"
    local body="$3"
    local description="$4"
    
    log_info "Testing: $description"
    
    local response
    if [ "$method" = "GET" ]; then
        response=$(curl -s -w "HTTPSTATUS:%{http_code}" \
            -H "Content-Type: application/json" \
            "$BASE_URL$endpoint")
    else
        response=$(curl -s -w "HTTPSTATUS:%{http_code}" \
            -X "$method" \
            -H "Content-Type: application/json" \
            -d "$body" \
            "$BASE_URL$endpoint")
    fi
    
    local http_code=$(echo "$response" | grep -o "HTTPSTATUS:[0-9]*" | cut -d: -f2)
    local body_response=$(echo "$response" | sed -E 's/HTTPSTATUS:[0-9]*$//')
    
    echo "  HTTP Status: $http_code"
    echo "  Response: $body_response"
    
    if [ "$http_code" -ge 200 ] && [ "$http_code" -lt 300 ]; then
        log_success "$description - PASSED"
    else
        log_error "$description - FAILED (HTTP $http_code)"
    fi
    echo
}

# Check if server is running
check_server() {
    log_info "Checking if server is running on $BASE_URL"
    if curl -s --connect-timeout 5 "$BASE_URL/things/device/config" > /dev/null; then
        log_success "Server is running"
    else
        log_error "Server is not running or not accessible"
        exit 1
    fi
    echo
}

# Main test execution
main() {
    echo "=================================================="
    echo "Nexus IoT Platform HTTP API Test Suite"
    echo "=================================================="
    echo
    
    check_server
    
    # Test User APIs (Public endpoints)
    echo "=== User API Tests ==="
    
    test_public_endpoint "POST" "/user/sign_up" '{
        "account": "test_user_'$(date +%s)'",
        "username": "testuser'$(date +%s)'",
        "password": "password123",
        "region": "US"
    }' "User Sign Up"
    
    test_public_endpoint "POST" "/user/sign_in" '{
        "account": "admin",
        "password": "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi"
    }' "User Sign In"
    
    # Test Device APIs (Authenticated endpoints)
    echo "=== Device API Tests ==="
    
    test_public_endpoint "POST" "/things/device/activate" '{
        "license_key": "NEXUS-PLAT-2025-002-EFGH5678",
        "device_id": "'$DEVICE_ID'",
        "product_id": "prod-001"
    }' "Device Activation"
    
    test_device_endpoint "GET" "/things/device/details" "" "Get Device Details"
    
    test_device_endpoint "GET" "/things/device/config" "" "Get Device Config"
    
    test_device_endpoint "GET" "/things/device/credentials" "" "Get Storage Credentials"
    
    test_device_endpoint "POST" "/things/device/deactivate" '{}' "Device Deactivation"
    
    # Test Authentication Failures
    echo "=== Authentication Tests ==="
    
    log_info "Testing: Missing device-id header"
    response=$(curl -s -w "HTTPSTATUS:%{http_code}" \
        -H "Content-Type: application/json" \
        "$BASE_URL/things/device/details")
    http_code=$(echo "$response" | grep -o "HTTPSTATUS:[0-9]*" | cut -d: -f2)
    if [ "$http_code" = "401" ]; then
        log_success "Missing device-id header - PASSED (Expected 401)"
    else
        log_error "Missing device-id header - FAILED (Expected 401, got $http_code)"
    fi
    echo
    
    log_info "Testing: Invalid device-id"
    timestamp=$(get_timestamp)
    digest=$(generate_digest "" "invalid-secret")
    response=$(curl -s -w "HTTPSTATUS:%{http_code}" \
        -H "Content-Type: application/json" \
        -H "device-id: invalid-device" \
        -H "time: $timestamp" \
        -H "digest: $digest" \
        "$BASE_URL/things/device/details")
    http_code=$(echo "$response" | grep -o "HTTPSTATUS:[0-9]*" | cut -d: -f2)
    if [ "$http_code" = "401" ]; then
        log_success "Invalid device-id - PASSED (Expected 401)"
    else
        log_error "Invalid device-id - FAILED (Expected 401, got $http_code)"
    fi
    echo
    
    # Test Error Cases
    echo "=== Error Handling Tests ==="
    
    test_public_endpoint "POST" "/user/sign_in" '{
        "invalid": "json"
    }' "Invalid JSON for Sign In"
    
    test_device_endpoint "GET" "/things/device/details" "" "Valid Device Details Request"
    
    echo "=================================================="
    echo "Test Suite Completed"
    echo "=================================================="
}

# Script entry point
if [ "$#" -eq 0 ]; then
    main
else
    case "$1" in
        "user")
            echo "=== User API Tests Only ==="
            test_public_endpoint "POST" "/user/sign_up" '{
                "account": "test_user_'$(date +%s)'",
                "username": "testuser'$(date +%s)'",
                "password": "password123",
                "region": "US"
            }' "User Sign Up"
            
            test_public_endpoint "POST" "/user/sign_in" '{
                "account": "admin",
                "password": "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi"
            }' "User Sign In"
            ;;
        "device")
            echo "=== Device API Tests Only ==="
            check_server
            test_public_endpoint "POST" "/things/device/activate" '{
                "license_key": "NEXUS-PLAT-2025-002-EFGH5678",
                "device_id": "'$DEVICE_ID'",
                "product_id": "prod-001"
            }' "Device Activation"
            
            test_device_endpoint "GET" "/things/device/details" "" "Get Device Details"
            test_device_endpoint "GET" "/things/device/config" "" "Get Device Config"
            test_device_endpoint "GET" "/things/device/credentials" "" "Get Storage Credentials"
            ;;
        "auth")
            echo "=== Authentication Tests Only ==="
            check_server
            log_info "Testing: Missing headers"
            curl -s "$BASE_URL/things/device/details" | jq . || echo "No JSON response"
            ;;
        *)
            echo "Usage: $0 [user|device|auth]"
            echo "  user   - Test only user APIs"
            echo "  device - Test only device APIs"  
            echo "  auth   - Test only authentication"
            echo "  (no args) - Run all tests"
            ;;
    esac
fi
