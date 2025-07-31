#!/bin/bash

# Nexus IoT Platform API Test Suite
# Complete API tests with proper authentication

set -e

# Configuration
BASE_URL="http://127.0.0.1:7070"
DEVICE_ID="dev-002"
DEVICE_SECRET="device-secret-002"
TOKEN=""  # Will be set by sign in test

# Device activation response storage
ACTIVATED_DEVICE_ID=""
ACTIVATED_DEVICE_SECRET=""

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

log() {
    echo -e "${BLUE}[TEST]${NC} $1"
}

success() {
    echo -e "${GREEN}[PASS]${NC} $1"
}

error() {
    echo -e "${RED}[FAIL]${NC} $1"
}

# Generate auth headers
generate_auth_headers() {
    local body="$1"
    local device_id="${2:-$DEVICE_ID}"
    local digest
    digest=$(echo -n "$device_id" | openssl dgst -sha256 -hmac "$DEVICE_SECRET" -binary | base64)
    echo "-H \"device-id: $device_id\" -H \"digest: $digest\""
}

# Test 1: User Sign Up
test_user_sign_up() {
    log "User Sign Up"
    local username="testuser_$(date +%s)"
    curl -s -X POST "$BASE_URL/cloud/account/sign_up" \
        -H "Content-Type: application/json" \
        -d "{
            \"account\": \"$username\",
            \"username\": \"$username\",
            \"password\": \"password123\",
            \"region\": \"US\"
        }" | jq .
}

# Test 2: User Sign In and get token
test_user_sign_in() {
    log "User Sign In"
    local response
    response=$(curl -s -X POST "$BASE_URL/cloud/account/sign_in" \
        -H "Content-Type: application/json" \
        -d '{
            "account": "admin",
            "password": "password123456"
        }')
    echo "$response" | jq .

    # Extract token from response if successful
    if echo "$response" | jq -e '.success == true' > /dev/null; then
        TOKEN=$(echo "$response" | jq -r '.result.token // empty')
        if [ -n "$TOKEN" ] && [ "$TOKEN" != "null" ]; then
            success "Token obtained: ${TOKEN:0:20}..."
        else
            error "No token in response"
        fi
    else
        error "Sign in failed"
    fi
}


# Test 3: Get Device List
test_device_list() {
    log "Get Device List"
    if [ -z "$TOKEN" ]; then
        error "No token available, please run sign in first"
        return 1
    fi
    local response
    response=$(curl -s -X GET "$BASE_URL/cloud/devices" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $TOKEN" \
        -d '{ "group_list": ["group-001"], "page": 1, "page_size": 10}')
    echo "$response" | jq .
}

# Test 4: User Details
test_user_details() {
    log "User Details"
    if [ -z "$TOKEN" ]; then
        error "No token available, please run sign in first"
        return 1
    fi

    curl -s -X GET "$BASE_URL/cloud/user/details" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $TOKEN" | jq .
}

# Test 5: Get Device Details
test_device_details() {
    log "Get Device Details"
    if [ -z "$TOKEN" ]; then
        error "No token available, please run sign in first"
        return 1
    fi

    local response
    response=$(curl -s -X GET "$BASE_URL/cloud/device/$DEVICE_ID" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $TOKEN")
    echo "$response" | jq .
}


# Test 6: Device Activation
test_device_activation() {
    log "Device Activation"
    local license_id="lic-001"
    local license_key="NEXUS-PLAT-2025-001-ABCD1234"
    local auth_digest
    auth_digest=$(echo -n "$license_id" | openssl dgst -sha256 -hmac "$license_key" -binary | base64)
    echo "Calculated digest: $auth_digest"

    local response
    response=$(curl -s -X POST "$BASE_URL/things/activate" \
        -H "Content-Type: application/json" \
        -d "{
            \"license_id\": \"$license_id\",
            \"authenticate\": \"$auth_digest\",
            \"product_id\": \"prod-001\",
            \"version\": \"v0.0.1\",
            \"sdk_version\": \"v1.0.0\",
            \"ip\": \"192.168.1.100\",
            \"state\": {
                \"siren\": null,
                \"video\": {
                    \"osd\": true,
                    \"flip\": false,
                    \"sharpness\": 50,
                    \"brightness\": 50
                },
                \"cruise\": null,
                \"record\": {
                    \"mode\": \"continuous\",
                    \"duration\": 60
                },
                \"volume\": 80,
                \"storage\": {
                    \"mode\": \"cloud\",
                    \"status\": true,
                    \"capacity\": 128
                },
                \"night_vision\": true,
                \"privacy_mode\": false,
                \"motion_tracking\": true,
                \"motion_detection\": {
                    \"area\": [0, 0, 100, 100],
                    \"status\": true,
                    \"sensitivity\": 70
                },
                \"decibel_detection\": {
                    \"status\": false,
                    \"sensitivity\": 50
                }
            },
            \"features\": {
                \"ai\": \"remote\",
                \"p2p\": \"enhance\",
                \"upnp\": \"IDGV1\",
                \"webrtc\": [\"SRTP\", \"DC\"],
                \"audio_feature\": {
                    \"codecs\": [\"AAC\"],
                    \"channels\": 2,
                    \"sampling_rate\": 16000,
                    \"bits\": 16
                },
                \"video_feature\": {
                    \"num\": 2,
                    \"resolution_ratios\": [{\"width\": 1920, \"height\": 1080}],
                    \"codecs\": [\"H264\", \"H265\"],
                    \"streams\": [1, 2]
                }
            }
        }")

    echo "$response" | jq .

    # Extract device information if activation was successful
    if echo "$response" | jq -e '.success == true' > /dev/null; then
        ACTIVATED_DEVICE_ID=$(echo "$response" | jq -r '.result.id // empty')
        ACTIVATED_DEVICE_SECRET=$(echo "$response" | jq -r '.result.secret_key // empty')
        if [ -n "$ACTIVATED_DEVICE_ID" ] && [ "$ACTIVATED_DEVICE_ID" != "null" ]; then
            success "Device activated with ID: $ACTIVATED_DEVICE_ID"
        else
            error "No device ID in activation response"
        fi
    else
        error "Device activation failed"
    fi
}

# Test 9: Device Deactivation
test_device_deactivation() {
    log "Device Deactivation"
    local body="{}"
    local auth_headers

    local digest
    digest=$(echo -n "$ACTIVATED_DEVICE_ID" | openssl dgst -sha256 -hmac "$ACTIVATED_DEVICE_SECRET" -binary | base64)

    local response
    response=$(curl -s -X POST "$BASE_URL/things/deactivate" \
        -H "Content-Type: application/json" \
        -H "device-id: $ACTIVATED_DEVICE_ID" \
        -H "digest: $digest" \
        -d "'$body'")

    echo "$response" | jq .
}

# Test 7: Get Device Config
test_device_config() {
    log "Get Device Config"
    local auth_headers
    auth_headers=$(generate_auth_headers "")
    eval curl -s -X GET "$BASE_URL/things/device/config" \
        -H "Content-Type: application/json" \
        $auth_headers | jq .
}

# Test 8: Get Storage Credentials
test_storage_credentials() {
    log "Get Storage Credentials"
    local auth_headers
    auth_headers=$(generate_auth_headers "")
    eval curl -s -X GET "$BASE_URL/things/device/credentials" \
        -H "Content-Type: application/json" \
        $auth_headers | jq .
}



# Test 10: Authentication Error (Missing Headers)
test_auth_error() {
    log "Authentication Error Test (Missing Headers)"
    curl -s -X GET "$BASE_URL/things/device/config" \
        -H "Content-Type: application/json" | jq .
}

# Test 11: Invalid Device ID
test_invalid_device() {
    log "Invalid Device ID Test"
    local auth_headers
    auth_headers=$(generate_auth_headers "" "invalid-device")
    eval curl -s -X GET "$BASE_URL/things/device/config" \
        -H "Content-Type: application/json" \
        $auth_headers | jq .
}

# Main execution
main() {
    echo "=========================================="
    echo "Nexus IoT Platform API Test Suite"
    echo "=========================================="
    echo

    # Check if server is running
    if ! curl -s --connect-timeout 5 "$BASE_URL/health" > /dev/null; then
        error "Server is not running at $BASE_URL"
        exit 1
    fi
    success "Server is running"
    echo

    # Run tests based on parameter
    case "${1:-all}" in
        "user_api")
            test_user_sign_up
            echo
            test_user_sign_in
            echo
            test_user_details
            echo
            test_device_list
            echo
            test_device_details
            ;;
        "device_api")
            test_device_activation
            echo
            test_device_deactivation
            echo
            test_device_config
            echo
            test_storage_credentials
            ;;
        "auth")
            test_auth_error
            echo
            test_invalid_device
            ;;
        "all")
            echo "=== User Tests ==="
            test_user_sign_up
            echo
            test_user_sign_in
            echo
            test_user_details
            echo
            test_device_list
            echo
            test_device_details
            echo
            echo "=== Device Tests ==="
            test_device_activation
            echo
            test_device_deactivation
            echo
            test_device_config
            echo
            test_storage_credentials
            echo
            echo "=== Authentication Tests ==="
            test_auth_error
            echo
            test_invalid_device
            ;;
        *)
            echo "Usage: $0 [user|device|auth|all]"
            echo "  user   - Test user APIs only"
            echo "  device - Test device APIs only"
            echo "  auth   - Test authentication only"
            echo "  all    - Run all tests (default)"
            ;;
    esac

    echo
    echo "=========================================="
    echo "Test suite completed"
    echo "=========================================="
}

# Execute main function
main "$@"
