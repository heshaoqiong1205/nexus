-- PostgreSQL initialization script for Nexus Things Platform
-- This script creates the database, user, and basic tables

-- Create database and user (run as postgres superuser)
-- Note: These commands may need to be run separately with psql as superuser

-- Create user if it doesn't exist
DO
$do$
BEGIN
   IF NOT EXISTS (
      SELECT FROM pg_catalog.pg_roles
      WHERE  rolname = 'things') THEN

      CREATE ROLE things LOGIN PASSWORD '123456';
   END IF;
END
$do$;

-- Create database if it doesn't exist
SELECT 'CREATE DATABASE things OWNER things'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'things')\gexec

-- Grant privileges
GRANT ALL PRIVILEGES ON DATABASE things TO things;

-- Connect to the things database
\c things things

-- Enable extensions if needed
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create tables based on the model structures

-- IoT Devices table
CREATE TABLE IF NOT EXISTS devices (
    id VARCHAR(255) PRIMARY KEY,
    secret_key VARCHAR(255) NOT NULL,
    license_id VARCHAR(255),
    name VARCHAR(255) NOT NULL,
    product_id VARCHAR(255),
    group_id VARCHAR(255),
    features JSONB,
    state JSONB,
    version VARCHAR(100),
    sdk_version VARCHAR(100),
    ip VARCHAR(45),
    online BOOLEAN DEFAULT false,
    location POINT,
    status BOOLEAN DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    active_at TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(255) PRIMARY KEY,
    account VARCHAR(255),
    username VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    region VARCHAR(255),
    location POINT,
    icon VARCHAR(255),
    role VARCHAR(255),
    last_login_time TIMESTAMP,
    status BOOLEAN DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Products table
CREATE TABLE IF NOT EXISTS products (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    required_features JSONB,
    status BOOLEAN DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Applications table
CREATE TABLE IF NOT EXISTS applications (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    secret_key VARCHAR(255) NOT NULL,
    salt VARCHAR(255),
    status BOOLEAN DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Buckets table (for storage)
CREATE TABLE IF NOT EXISTS buckets (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    region VARCHAR(100),
    provider VARCHAR(100),
    endpoint VARCHAR(255),
    status BOOLEAN DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Cloud Storage table
CREATE TABLE IF NOT EXISTS cloud_storages (
    id VARCHAR(255) PRIMARY KEY,
    device_id VARCHAR(255),
    tos VARCHAR(100),
    bucket VARCHAR(255),
    mode VARCHAR(100),
    path VARCHAR(1000),
    status BOOLEAN DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Licenses table
CREATE TABLE IF NOT EXISTS licenses (
    id VARCHAR(255) PRIMARY KEY,
    key VARCHAR(255) UNIQUE NOT NULL,
    status BOOLEAN DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Messages table
CREATE TABLE IF NOT EXISTS messages (
    id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255),
    type VARCHAR(50),
    content JSONB,
    status BOOLEAN DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Cloud Recording Plans table
CREATE TABLE IF NOT EXISTS cloud_recording_plans (
    id VARCHAR(255) PRIMARY KEY,
    device_id VARCHAR(255),
    channel INTEGER,
    enabled BOOLEAN DEFAULT false,
    mode VARCHAR(50),
    mfd INTEGER,
    interval INTEGER,
    storage_id VARCHAR(255),
    status BOOLEAN DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Cloud Recordings table
CREATE TABLE IF NOT EXISTS cloud_recordings (
    id VARCHAR(255) PRIMARY KEY,
    device_id VARCHAR(255),
    channel INTEGER,
    fragment_duration INTEGER,
    bucket_id VARCHAR(255),
    prefix VARCHAR(255),
    fragments JSONB,
    state VARCHAR(50),
    begin_time TIMESTAMP,
    end_time TIMESTAMP,
    status BOOLEAN DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_devices_status ON devices(status);
CREATE INDEX IF NOT EXISTS idx_devices_online ON devices(online);
CREATE INDEX IF NOT EXISTS idx_devices_group_id ON devices(group_id);
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_account ON users(account);
CREATE INDEX IF NOT EXISTS idx_licenses_key ON licenses(key);
CREATE INDEX IF NOT EXISTS idx_messages_user_id ON messages(user_id);
CREATE INDEX IF NOT EXISTS idx_cloud_recordings_device_id ON cloud_recordings(device_id);
CREATE INDEX IF NOT EXISTS idx_cloud_recordings_begin_time ON cloud_recordings(begin_time);
CREATE INDEX IF NOT EXISTS idx_cloud_storages_device_id ON cloud_storages(device_id);

-- Insert some initial data for testing
INSERT INTO products (id, name, description, required_features, status, created_at, updated_at)
VALUES ('prod-001', 'Nexus IoT Platform', 'Core IoT platform for device management', '[]'::JSONB, true,
        NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

INSERT INTO applications (id, name, description, secret_key, salt, status, created_at, updated_at)
VALUES ('app-001', 'Default Application', 'Default application for testing', 'default-secret-key',
        'default-salt', true, NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- Insert initial buckets data
INSERT INTO buckets (id, name, region, provider, endpoint, status, created_at, updated_at)
VALUES
    ('bucket-001', 'nexus-device-data', 'us-west-1', 'mock', 's3.us-west-1.amazonaws.com', true,
     NOW(), NOW()),
    ('bucket-002', 'nexus-media-storage', 'us-east-1', 'mock', 's3.us-east-1.amazonaws.com', true,
     NOW(), NOW()),
    ('bucket-003', 'nexus-backup-storage', 'ap-southeast-1', 'aliyun', 'oss-ap-southeast-1.aliyun.com', true,
     NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- Insert a test user first for license reference
INSERT INTO users (id, account, username, password, region, location, role, last_login_time, status, created_at, updated_at)
VALUES ('user-001', 'test_user', 'test_user',
        'password123456', 'US', POINT(116.4074, 39.9042), 'admin', NOW(), true,
        NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- Insert additional test users
INSERT INTO users (id, account, username, password, region, location, role, last_login_time, status, created_at, updated_at)
VALUES
    ('user-002', 'admin', 'admin',
     'password123456',
     'US', POINT(-74.0059, 40.7128), 'admin', NOW(), true, NOW(), NOW()),
    ('user-003', 'demo_user', 'demo_user',
     'password123456',
     'CN', POINT(116.4074, 39.9042), 'user', NOW(), true, NOW(), NOW()),
    ('user-004', 'operator', 'operator',
     'password123456',
     'EU', POINT(-0.1278, 51.5074), 'operator', NOW(), true, NOW(), NOW()),
    ('user-005', 'guest', 'guest',
     'password123456',
     'US', POINT(-122.4194, 37.7749), 'guest', NOW(), true, NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- Insert initial licenses data
INSERT INTO licenses (id, key, status, created_at, updated_at)
VALUES
    ('lic-001', 'NEXUS-PLAT-2025-001-ABCD1234', true, NOW(), NOW()),
    ('lic-002', 'NEXUS-PLAT-2025-002-EFGH5678', true, NOW(), NOW()),
    ('lic-003', 'NEXUS-TRIAL-2025-001-IJKL9012', true, NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- Insert sample IoT devices with location points
INSERT INTO devices (id, secret_key, license_id, name, product_id, group_id,
                               features, state, version, sdk_version, ip, online, location,
                               status, created_at, active_at, updated_at)
VALUES
    ('dev-001', 'device-secret-001', 'lic-001', 'Temperature Sensor NYC', 'prod-001', 'group-001',
     '{"p2p": "standard", "webrtc": ["SRTP"], "upnp": "enable", "ai": "local", "video_feature": null, "audio_feature": null}'::JSONB,
     '{"video": null, "storage": {"mode": "local", "capacity": 32, "status": true}, "record": null, "motion_detection": null, "decibel_detection": null, "cruise": null, "siren": null, "volume": 50, "privacy_mode": false, "night_vision": false, "motion_tracking": false}'::JSONB,
     '1.0.0', '2.1.0', '192.168.1.10', true, POINT(-74.0059, 40.7128),
     true, NOW(), NOW(), NOW()),
    ('dev-002', 'device-secret-002', 'lic-002', 'Camera Device Beijing', 'prod-001', 'group-002',
     '{"p2p": "enhance", "webrtc": ["SRTP", "DC"], "upnp": "enable", "ai": "remote", "video_feature": {"resolution": ["1080P", "720P"], "codec": ["H264", "H265"], "bitrate": [2048, 4096], "fps": [15, 30]}, "audio_feature": {"codec": ["AAC"], "sample_rate": [8000, 16000], "bitrate": [64, 128]}}'::JSONB,
     '{"video": {"flip": false, "osd": true, "brightness": 50, "sharpness": 50}, "storage": {"mode": "cloud", "capacity": 128, "status": true}, "record": {"mode": "continuous", "duration": 60}, "motion_detection": {"status": true, "sensitivity": 70, "area": [0, 0, 100, 100]}, "decibel_detection": {"status": false, "sensitivity": 50}, "cruise": null, "siren": null, "volume": 80, "privacy_mode": false, "night_vision": true, "motion_tracking": true}'::JSONB,
     '1.1.0', '2.1.0', '192.168.1.11', false, POINT(116.4074, 39.9042),
     true, NOW(), NOW() - INTERVAL '1 hour', NOW()),
    ('dev-003', 'device-secret-003', 'lic-003', 'Smart Lock London', 'prod-001', 'group-001',
     '{"p2p": "standard", "webrtc": ["SRTP"], "upnp": "disable", "ai": null, "video_feature": null, "audio_feature": {"codec": ["AAC"], "sample_rate": [8000], "bitrate": [64]}}'::JSONB,
     '{"video": null, "storage": {"mode": "local", "capacity": 16, "status": true}, "record": null, "motion_detection": null, "decibel_detection": {"status": true, "sensitivity": 80}, "cruise": null, "siren": {"duration": 10, "volume": 90}, "volume": 60, "privacy_mode": false, "night_vision": false, "motion_tracking": false}'::JSONB,
     '1.0.5', '2.0.0', '192.168.1.12', true, POINT(-0.1278, 51.5074),
     true, NOW(), NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

INSERT INTO cloud_storages (id, device_id, tos, bucket, mode, path, status, created_at, updated_at)
VALUES
    ('storage-001', 'dev-001', 'media7', 'nexus-device-data', 'standard', '/data', true, NOW(), NOW()),
    ('storage-002', 'dev-002', 'media7', 'nexus-media-storage', 'enhance', '/media', true, NOW(), NOW()),
    ('storage-003', 'dev-003', 'media7', 'nexus-log-storage', 'standard', '/logs', true, NOW(), NOW());

-- Print success message
\echo 'PostgreSQL database initialization completed successfully!'
\echo 'Database: things'
\echo 'User: things'
\echo 'Tables created without prefix (matching GORM model names)'
