# PostgreSQL Setup for Nexus Things Platform

This directory contains scripts and configuration for setting up PostgreSQL database for the Nexus Things Platform.

## Quick Start

### Prerequisites

- Docker (for containerized PostgreSQL)
- PostgreSQL client tools (`psql`)
- Just command runner

### Setup Database

1. **Start and initialize PostgreSQL:**
   ```bash
   just pg-setup
   ```
   This command will:
   - Start PostgreSQL in Docker container
   - Wait for it to be ready
   - Run initialization script to create database, user, and tables

2. **Test the connection:**
   ```bash
   just pg-test
   ```

3. **Build and run the application:**
   ```bash
   just run
   ```

## Available Commands

### PostgreSQL Management
- `just pg-start` - Start PostgreSQL container
- `just pg-init` - Initialize database schema
- `just pg-clean` - Stop and remove PostgreSQL container
- `just pg-logs` - View PostgreSQL logs
- `just pg-connect` - Connect as things user
- `just pg-connect-admin` - Connect as postgres superuser
- `just pg-setup` - Complete setup (start + init)
- `just pg-test` - Test connection and show database info

### Application Management
- `just build` - Build the application
- `just run` - Build and run the application
- `just dev-setup` - Setup complete development environment
- `just dev-clean` - Clean up development environment

## Database Configuration

The database configuration is loaded from `conf/app.ini`:

```ini
[database]
Type = postgres
User = things
Password = 123456
Host = 127.0.0.1
Port = 5432
Name = things
TablePrefix = things_
```

## Database Schema

The initialization script creates the following tables:

- `things_iot_devices` - IoT device management
- `things_users` - User accounts
- `things_products` - Product catalog
- `things_applications` - Application management
- `things_buckets` - Storage buckets
- `things_cloud_storages` - Cloud storage configuration
- `things_licenses` - License management
- `things_messages` - Message system
- `things_cloud_recording_plans` - Recording plans
- `things_cloud_recordings` - Recording data

All tables use the `things_` prefix as specified in the configuration.

## Manual Setup

If you prefer to set up PostgreSQL manually:

1. **Install PostgreSQL:**
   ```bash
   # macOS with Homebrew
   brew install postgresql
   brew services start postgresql
   
   # Ubuntu/Debian
   sudo apt-get install postgresql postgresql-contrib
   sudo systemctl start postgresql
   ```

2. **Run initialization script:**
   ```bash
   ./scripts/init_postgres.sh
   ```

3. **Test connection:**
   ```bash
   ./scripts/test_postgres.sh
   ```

## Troubleshooting

### Connection Issues
- Ensure PostgreSQL is running: `just pg-logs`
- Check if port 5432 is available: `lsof -i :5432`
- Verify Docker container is running: `docker ps`

### Permission Issues
- Make sure scripts are executable: `chmod +x scripts/*.sh`
- Check PostgreSQL authentication in `pg_hba.conf`

### Database Issues
- Reset database: `just pg-clean && just pg-setup`
- Check logs: `just pg-logs`
- Connect manually: `just pg-connect-admin`

## Security Notes

- Default passwords are set for development only
- Change passwords in production environments
- Consider using environment variables for sensitive data
- Enable SSL for production databases
