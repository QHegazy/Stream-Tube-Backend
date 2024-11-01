CREATE SCHEMA IF NOT EXISTS "auth";
CREATE SCHEMA IF NOT EXISTS "session";
SET search_path TO "auth", public;

-- Enum Types
CREATE TYPE auth.user_status AS ENUM ('active', 'inactive', 'suspended', 'banned');
CREATE TYPE auth.oauth_provider AS ENUM ('google', 'facebook', 'github', 'microsoft');
CREATE TYPE auth.contact_type AS ENUM ('email', 'phone');
CREATE TYPE auth.auth_method AS ENUM ('local', 'oauth');

-- Users Table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(255) UNIQUE NOT NULL,
    auth_method auth.auth_method NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    status auth.user_status DEFAULT 'active',
    email_verified TIMESTAMP DEFAULT NULL,
    phone_verified TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP DEFAULT NULL,
    CONSTRAINT valid_username CHECK (username ~* '^[A-Za-z0-9._-]{3,50}$')
);

-- Profiles Table
CREATE TABLE IF NOT EXISTS profiles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    full_name VARCHAR(255) NOT NULL,
    profile_picture_url VARCHAR(255),
    user_id UUID NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP DEFAULT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- OAuth Users Table
CREATE TABLE IF NOT EXISTS oauth_users (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    provider auth.oauth_provider NOT NULL,
    access_token TEXT NOT NULL,
    refresh_token TEXT,
    expires_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP DEFAULT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE (user_id, provider) -- Ensure unique provider per user
);

-- Local Users Table
CREATE TABLE IF NOT EXISTS local_users (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    password_hash VARCHAR(512) NOT NULL, -- Increased size for hashing flexibility
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP DEFAULT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Trigger for Automatic `updated_at` Updates
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Apply Trigger to All Relevant Tables
CREATE TRIGGER set_timestamp
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON profiles
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON oauth_users
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON local_users
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
