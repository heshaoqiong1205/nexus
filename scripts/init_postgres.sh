#!/bin/bash

# PostgreSQL initialization script for Nexus Things Platform
# This script sets up the database, user, and tables

set -e  # Exit on any error

# Configuration from app.ini
DB_NAME="things"
DB_USER="things"
DB_PASSWORD="123456"
DB_HOST="127.0.0.1"
DB_PORT="5432"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}Nexus Things Platform - PostgreSQL Initialization${NC}"
echo "=================================================="

# Check if PostgreSQL is running
echo -e "${YELLOW}Checking PostgreSQL connection...${NC}"

# Check if we should use Docker or local psql
if command -v psql >/dev/null 2>&1; then
    PSQL_CMD="psql"
    echo -e "${GREEN}Using local psql client${NC}"
else
    PSQL_CMD="docker exec nexus-postgres psql"
    echo -e "${GREEN}Using Docker container psql client${NC}"
fi

# Wait for PostgreSQL to be ready (up to 30 seconds)
echo -e "${YELLOW}Waiting for PostgreSQL to be ready...${NC}"
for i in {1..30}; do
    if $PSQL_CMD -h $DB_HOST -p $DB_PORT -U postgres -c '\q' >/dev/null 2>&1; then
        echo -e "${GREEN}✓ PostgreSQL is ready after ${i} seconds${NC}"
        break
    fi
    if [ $i -eq 30 ]; then
        echo -e "${RED}Error: PostgreSQL is not running or not accessible at $DB_HOST:$DB_PORT after 30 seconds${NC}"
        echo "Please ensure PostgreSQL is installed and running."
        echo ""
        echo "For Docker containers, check the logs:"
        echo "  docker logs nexus-postgres"
        echo ""
        echo "Container may still be starting. Try running:"
        echo "  just pg-init"
        echo "again in a few moments."
        exit 1
    fi
    echo "Waiting... (${i}/30)"
    sleep 1
done

# Check if we can connect as postgres user (for initial setup)
echo -e "${YELLOW}Setting up database and user...${NC}"

# Try to connect as postgres superuser first
if $PSQL_CMD -h $DB_HOST -p $DB_PORT -U postgres -c '\q' >/dev/null 2>&1; then
    echo -e "${GREEN}✓ Connected as postgres superuser${NC}"
    SUPERUSER="postgres"
elif $PSQL_CMD -h $DB_HOST -p $DB_PORT -U "$(whoami)" -c '\q' >/dev/null 2>&1; then
    echo -e "${GREEN}✓ Connected as $(whoami)${NC}"
    SUPERUSER="$(whoami)"
else
    echo -e "${RED}Error: Cannot connect to PostgreSQL${NC}"
    echo "Please ensure you have proper PostgreSQL credentials."
    echo "You may need to:"
    echo "1. Set up PostgreSQL authentication"
    echo "2. Create a superuser account"
    echo "3. Update pg_hba.conf for local connections"
    exit 1
fi

# Run the database initialization SQL script
echo "Running database initialization script..."
if command -v psql >/dev/null 2>&1; then
    # Use local psql client, connect to postgres database first
    PGPASSWORD=$POSTGRES_PASSWORD $PSQL_CMD -d postgres -v ON_ERROR_STOP=1 < scripts/init_postgres.sql
else
    # Use Docker exec with SQL content, connect to postgres database first
    docker exec -i nexus-postgres psql -h $DB_HOST -p $DB_PORT -U postgres -d postgres -v ON_ERROR_STOP=1 < scripts/init_postgres.sql
fi

# Test connection with the new user
echo -e "${YELLOW}Testing connection with new user...${NC}"
if PGPASSWORD=$DB_PASSWORD $PSQL_CMD -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c '\dt' >/dev/null 2>&1; then
    echo -e "${GREEN}✓ Successfully connected as $DB_USER${NC}"
else
    echo -e "${RED}Error: Cannot connect as $DB_USER${NC}"
    exit 1
fi

# Show table count
TABLE_COUNT=$(PGPASSWORD=$DB_PASSWORD $PSQL_CMD -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -t -c "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public' AND table_name LIKE 'things_%';" | tr -d ' ')

echo ""
echo -e "${GREEN}Database setup completed successfully!${NC}"
echo "================================"
echo "Database: $DB_NAME"
echo "User: $DB_USER"
echo "Host: $DB_HOST:$DB_PORT"
echo "Tables created: $TABLE_COUNT"
echo ""
echo -e "${YELLOW}Connection string for application:${NC}"
echo "host=$DB_HOST user=$DB_USER password=$DB_PASSWORD dbname=$DB_NAME port=$DB_PORT sslmode=disable"
echo ""
echo -e "${YELLOW}To connect manually:${NC}"
if command -v psql >/dev/null 2>&1; then
    echo "PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME"
else
    echo "docker exec -it nexus-postgres psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME"
    echo "Or install PostgreSQL client: brew install postgresql"
fi
