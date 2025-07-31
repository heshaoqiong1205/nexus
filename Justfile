DOCKER_CLI := `echo ${DOCKER_CLI:-docker}`

# PostgreSQL commands (Docker-based)
pg-start:
    @echo "Starting PostgreSQL with configuration matching app.ini..."
    @{{DOCKER_CLI}} run -d --name nexus-postgres \
        -e POSTGRES_USER=postgres \
        -e POSTGRES_PASSWORD=postgres \
        -e POSTGRES_DB=postgres \
        -p 5432:5432 \
        postgres:15

pg-init:
    @echo "Initializing PostgreSQL database for Nexus..."
    @./scripts/init_postgres.sh

pg-clean:
    @echo "Cleaning up PostgreSQL..."
    @{{DOCKER_CLI}} stop nexus-postgres || true
    @{{DOCKER_CLI}} rm nexus-postgres || true

pg-logs:
    @{{DOCKER_CLI}} logs nexus-postgres

pg-connect:
    @echo "Connecting to PostgreSQL as things user..."
    @PGPASSWORD=123456 psql -h 127.0.0.1 -p 5432 -U things -d things

pg-connect-admin:
    @echo "Connecting to PostgreSQL as postgres superuser..."
    @PGPASSWORD=postgres psql -h 127.0.0.1 -p 5432 -U postgres

# Complete PostgreSQL setup (start + initialize)
pg-setup: pg-start
    @echo "Waiting for PostgreSQL to be ready..."
    @sleep 10
    @just pg-init

pg-test:
    @echo "Testing PostgreSQL connection and setup..."
    @./scripts/test_postgres.sh

# Build and run application
build:
    @echo "Building Nexus application..."
    @go build -o nexus .

run: build
    @echo "Starting Nexus application..."
    @./nexus

# Development commands
dev-setup: pg-setup
    @echo "Setting up development environment..."
    @just build
    @echo "Development environment ready!"

dev-clean: pg-clean
    @echo "Cleaning up development environment..."
    @rm -f nexus

# Show database configuration
show-db-config:
    @echo "Displaying database configuration from app.ini..."
    @go run examples/db-config/main.go

docker-clean-all:
    @echo "Stopping all running containers..."
    @{{DOCKER_CLI}} stop $({{DOCKER_CLI}} ps -q) || true
    @echo "Removing all containers..."
    @{{DOCKER_CLI}} rm $({{DOCKER_CLI}} ps -aq) || true
    @echo "All Docker containers cleaned up"

docker-status:
    @echo "Current Docker containers:"
    @{{DOCKER_CLI}} ps -a

rabbitmq-start:
    @{{DOCKER_CLI}} run -d --name rabbitmq -p 5672:5672 -p 15672:15672 -e RABBITMQ_DEFAULT_USER=admin -e RABBITMQ_DEFAULT_PASS=admin rabbitmq:3.12-management
    @{{DOCKER_CLI}} cp deps/rabbitmq_delayed_message_exchange-3.12.0.ez rabbitmq:/plugins/.
    @{{DOCKER_CLI}} exec rabbitmq rabbitmq-plugins enable rabbitmq_delayed_message_exchange
    @{{DOCKER_CLI}} restart rabbitmq

rabbitmq-clean:
    @{{DOCKER_CLI}} stop rabbitmq
    @{{DOCKER_CLI}} rm rabbitmq

emqx-start:
    @{{DOCKER_CLI}} run -d --name emqx -p 1883:1883 -p 8081:8081 -p 8083:8083 -p 8084:8084 emqx:5.8.6

emqx-clean:
    @{{DOCKER_CLI}} stop emqx
    @{{DOCKER_CLI}} rm emqx
