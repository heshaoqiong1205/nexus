# Nexus IoT Platform HTTP API Quick Reference

## Base Configuration
```bash
BASE_URL="http://127.0.0.1:7070"
DEVICE_ID="dev-002"
DEVICE_SECRET="device-secret-002"
```

## Authentication Header Generation

For device APIs, you need three headers:
- `device-id`: The device identifier
- `time`: Current timestamp in RFC3339 format
- `digest`: SHA256 hash of (request_body + device_secret)

### Generate Headers Script
```bash
# Use the generate_auth_headers.sh script:
./scripts/generate_auth_headers.sh "" dev-002           # For GET requests
./scripts/generate_auth_headers.sh '{"key":"value"}' dev-002  # For POST requests
```

## User APIs (No Authentication Required)

### 1. User Sign Up
```bash
curl -X POST http://127.0.0.1:7070/user/sign_up \
  -H "Content-Type: application/json" \
  -d '{
    "account": "newuser123",
    "username": "newuser123",
    "password": "password123",
    "region": "US"
  }'
```

### 2. User Sign In
```bash
curl -X POST http://127.0.0.1:7070/user/sign_in \
  -H "Content-Type: application/json" \
  -d '{
    "account": "admin",
    "password": "password123"
  }'
```

## Device APIs (Authentication Required)

### 3. Device Activation (No Auth)
```bash
curl -X POST http://127.0.0.1:7070/things/device/activate \
  -H "Content-Type: application/json" \
  -d '{
    "license_key": "NEXUS-PLAT-2025-002-EFGH5678",
    "device_id": "dev-002",
    "product_id": "prod-001"
  }'
```

### 4. Get Device Details
```bash
# Generate auth headers first
TIMESTAMP=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
DIGEST=$(echo -n "device-secret-002" | shasum -a 256 | cut -d' ' -f1)

curl -X GET http://127.0.0.1:7070/things/device/details \
  -H "Content-Type: application/json" \
  -H "device-id: dev-002" \
  -H "time: $TIMESTAMP" \
  -H "digest: $DIGEST"
```

### 5. Get Device Config
```bash
curl -X GET http://127.0.0.1:7070/things/device/config \
  -H "Content-Type: application/json" \
  -H "device-id: dev-002" \
  -H "time: $TIMESTAMP" \
  -H "digest: $DIGEST"
```

### 6. Get Storage Credentials
```bash
curl -X GET http://127.0.0.1:7070/things/device/credentials \
  -H "Content-Type: application/json" \
  -H "device-id: dev-002" \
  -H "time: $TIMESTAMP" \
  -H "digest: $DIGEST"
```

### 7. Device Deactivation
```bash
# For POST requests with body, include body in digest calculation
BODY="{}"
DIGEST=$(echo -n "${BODY}device-secret-002" | shasum -a 256 | cut -d' ' -f1)

curl -X POST http://127.0.0.1:7070/things/device/deactivate \
  -H "Content-Type: application/json" \
  -H "device-id: dev-002" \
  -H "time: $TIMESTAMP" \
  -H "digest: $DIGEST" \
  -d "$BODY"
```

## Quick Test Scripts

### Run All Tests
```bash
./scripts/api_test_suite.sh all
```

### Run Specific Test Groups
```bash
./scripts/api_test_suite.sh user      # User APIs only
./scripts/api_test_suite.sh device    # Device APIs only  
./scripts/api_test_suite.sh auth      # Authentication tests only
```

### Individual Tests
```bash
./scripts/test_http_api_simple.sh signin     # Test user sign in
./scripts/test_http_api_simple.sh details    # Test device details
./scripts/test_http_api_simple.sh activate   # Test device activation
./scripts/test_http_api_simple.sh run-all    # Run all quick tests
```

## Expected Response Format

All APIs return JSON in this format:
```json
{
  "success": true,
  "result": { ... },
  "error": "error message if success is false"
}
```

## Common HTTP Status Codes

- `200 OK`: Successful request
- `201 Created`: Resource created successfully (sign up)
- `400 Bad Request`: Invalid request format
- `401 Unauthorized`: Authentication failed
- `404 Not Found`: Resource not found
- `409 Conflict`: Resource already exists
- `500 Internal Server Error`: Server error

## Authentication Digest Calculation

The digest is calculated as: `SHA256(request_body + device_secret)`

Examples:
- GET request (empty body): `SHA256("" + "device-secret-002")`
- POST request: `SHA256('{"key":"value"}' + "device-secret-002")`

## Troubleshooting

### Common Issues

1. **401 Unauthorized**: Check device-id, timestamp, and digest calculation
2. **Missing headers**: Ensure all three auth headers are present
3. **Invalid timestamp**: Use RFC3339 format and current time (±5 minutes)
4. **Wrong digest**: Verify body content and secret key in calculation

### Debug Commands
```bash
# Check server status
curl -s http://127.0.0.1:7070/things/device/config

# Test without authentication (should fail)
curl -s http://127.0.0.1:7070/things/device/details

# Verify timestamp format
date -u +"%Y-%m-%dT%H:%M:%SZ"

# Manual digest calculation
echo -n "request_bodydevice-secret-002" | shasum -a 256
```
