#!/bin/bash

# Device Authentication Header Generator
# Generates the required headers for device API authentication

DEVICE_SECRET="device-secret-002"

if [ "$#" -eq 0 ]; then
    echo "Usage: $0 <request_body> [device_id]"
    echo ""
    echo "Examples:"
    echo "  $0 '' dev-002                    # For GET requests"
    echo "  $0 '{\"key\":\"value\"}' dev-002   # For POST requests"
    echo ""
    echo "Output format:"
    echo "  -H \"device-id: <device_id>\""
    echo "  -H \"time: <timestamp>\""
    echo "  -H \"digest: <calculated_digest>\""
    exit 1
fi

REQUEST_BODY="$1"
DEVICE_ID="${2:-dev-002}"

# Generate timestamp
TIMESTAMP=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

# Generate digest: SHA256(request_body + device_secret)
DIGEST=$(echo -n "${REQUEST_BODY}${DEVICE_SECRET}" | shasum -a 256 | cut -d' ' -f1)

# Output the headers
echo "-H \"device-id: $DEVICE_ID\" \\"
echo "-H \"time: $TIMESTAMP\" \\"
echo "-H \"digest: $DIGEST\""
