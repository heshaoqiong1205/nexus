-- Simplified database schema initialization
-- Compatible with standard PostgreSQL installations

-- Enable UUID extension if available
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- IoT Devices table
CREATE TABLE IF NOT EXISTS devices (
    id VARCHAR(255) PRIMARY KEY,
    secret_key VARCHAR(255) NOT NULL,
    license_id VARCHAR(255),
    name VARCHAR(255) NOT NULL,
    product_id VARCHAR(255),
    group_id VARCHAR(255),
    features TEXT, -- Using TEXT instead of JSONB for compatibility
    state TEXT, -- Using TEXT instead of JSONB for compatibility
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
    required_features TEXT, -- Using TEXT instead of JSONB for compatibility
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

-- Buckets table
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
    content TEXT, -- Using TEXT instead of JSONB for compatibility
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
    fragments TEXT, -- Using TEXT instead of JSONB for compatibility
    state VARCHAR(50),
    begin_time TIMESTAMP,
    end_time TIMESTAMP,
    status BOOLEAN DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Desired States table for device state management
CREATE TABLE IF NOT EXISTS desired_states (
    id VARCHAR(255) PRIMARY KEY, -- DeviceID as primary key (one desired state per device)
    state TEXT NOT NULL, -- JSON format desired state
    version INTEGER DEFAULT 1, -- Version for optimistic locking
    status BOOLEAN DEFAULT true, -- Whether this desired state is valid/active
    last_desired_id INTEGER, -- Last desired ID reference
    confirmed_at TIMESTAMP, -- When device confirmed the desired state
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (id) REFERENCES devices(id) ON DELETE CASCADE
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
CREATE INDEX IF NOT EXISTS idx_desired_states_status ON desired_states(status);
CREATE INDEX IF NOT EXISTS idx_desired_states_version ON desired_states(version);

-- Insert initial test data
INSERT INTO products (id, name, description, required_features, status, created_at, updated_at)
VALUES ('prod-001', 'Nexus IoT Platform', 'Core IoT platform for device management', '[]', true,
        NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

INSERT INTO applications (id, name, description, secret_key, salt, status, created_at, updated_at)
VALUES ('app-001', 'Default Application', 'Default application for testing', 'default-secret-key',
        'default-salt', true, NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- Insert sample users with location points
INSERT INTO users (id, account, username, password, region, location, role, status, created_at, updated_at)
VALUES
    ('user-001', 'admin', 'admin',
     '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi',
     'US', POINT(-74.0059, 40.7128), 'admin', true, NOW(), NOW()),
    ('user-002', 'demo_user', 'demo_user',
     '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi',
     'CN', POINT(116.4074, 39.9042), 'user', true, NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- Insert sample IoT devices with location points
INSERT INTO devices (id, secret_key, license_id, name, product_id, group_id,
                               features, state, desired, version, sdk_version, ip, online, location,
                               status, created_at, active_at, updated_at)
VALUES
    ('dev-001', 'device-secret-001', 'lic-001', 'Temperature Sensor NYC', 'prod-001', 'group-001',
     '{"p2p": "standard", "webrtc": ["SRTP"], "upnp": "enable", "ai": "local", "video_feature": null, "audio_feature": null}',
     '{"video": null, "storage": {"mode": "local", "capacity": 32, "status": true}, "record": null, "motion_detection": null, "decibel_detection": null, "cruise": null, "siren": null, "volume": 50, "privacy_mode": false, "night_vision": false, "motion_tracking": false}',
     '{"volume": 100, "record": {"status": true, "mode": 1}, "timestamp": 1633072800000}',
     '1.0.0', '2.1.0', '192.168.1.10', true, POINT(-74.0059, 40.7128),
     true, NOW(), NOW(), NOW()),
    ('dev-002', 'device-secret-002', 'lic-002', 'Camera Device Beijing', 'prod-001', 'group-002',
     '{"p2p": "enhance", "webrtc": ["SRTP", "DC"], "upnp": "enable", "ai": "remote", "video_feature": {"resolution": ["1080P", "720P"], "codec": ["H264", "H265"], "bitrate": [2048, 4096], "fps": [15, 30]}, "audio_feature": {"codec": ["AAC"], "sample_rate": [8000, 16000], "bitrate": [64, 128]}}',
     '{"video": {"flip": false, "osd": true, "brightness": 50, "sharpness": 50}, "storage": {"mode": "cloud", "capacity": 128, "status": true}, "record": {"mode": "continuous", "duration": 60}, "motion_detection": {"status": true, "sensitivity": 70, "area": [0, 0, 100, 100]}, "decibel_detection": {"status": false, "sensitivity": 50}, "cruise": null, "siren": null, "volume": 80, "privacy_mode": false, "night_vision": true, "motion_tracking": true}',
     '{"volume": 100, "record": {"status": true, "mode": 1}, "timestamp": 1633072800000}',
     '1.1.0', '2.1.0', '192.168.1.11', false, POINT(116.4074, 39.9042),
     true, NOW(), NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
