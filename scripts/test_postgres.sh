#!/bin/bash

# Quick PostgreSQL connection test for Nexus

set -e

# Configuration from app.ini
DB_NAME="things"
DB_USER="things"
DB_PASSWORD="123456"
DB_HOST="127.0.0.1"
DB_PORT="5432"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${YELLOW}Testing PostgreSQL connection...${NC}"

# Test connection
if PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c '\dt' >/dev/null 2>&1; then
    echo -e "${GREEN}✓ Successfully connected to PostgreSQL${NC}"

    # Show table information
    echo -e "${YELLOW}Database information:${NC}"
    echo "Database: $DB_NAME"
    echo "User: $DB_USER"
    echo "Host: $DB_HOST:$DB_PORT"
    echo ""

    echo -e "${YELLOW}Tables in database:${NC}"
    PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "
        SELECT
            table_name,
            (SELECT COUNT(*) FROM information_schema.columns WHERE table_name = t.table_name) as column_count
        FROM information_schema.tables t
        WHERE table_schema = 'public'
        AND table_name LIKE 'things_%'
        ORDER BY table_name;
    "
else
    echo -e "${RED}✗ Failed to connect to PostgreSQL${NC}"
    echo "Please ensure:"
    echo "1. PostgreSQL is running"
    echo "2. Database and user are created"
    echo "3. Run 'just pg-setup' to initialize"
    exit 1
fi
