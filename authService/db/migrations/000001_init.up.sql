-- Enable UUID extension if not already enabled
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE SCHEMA IF NOT EXISTS "auth";
CREATE SCHEMA IF NOT EXISTS "session";
CREATE SCHEMA IF NOT EXISTS "security";
CREATE SCHEMA IF NOT EXISTS "notification";

SET search_path TO "auth", "security", "session", "notification", public;

-- Enhanced Enum Types
CREATE TYPE auth.user_status AS ENUM ('active', 'inactive', 'suspended', 'banned', 'pending_verification');
CREATE TYPE auth.oauth_provider AS ENUM ('google', 'facebook', 'github', 'microsoft');
CREATE TYPE auth.auth_method AS ENUM ('local', 'oauth', 'mfa');
CREATE TYPE auth.mfa_type AS ENUM ('authenticator', 'sms', 'email', 'security_key');
CREATE TYPE notification.notification_type AS ENUM ('security', 'account', 'marketing', 'system');
CREATE TYPE security.risk_level AS ENUM ('low', 'medium', 'high', 'critical');

-- Enhanced Users Table
CREATE TABLE IF NOT EXISTS auth.users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(255) UNIQUE NOT NULL,
    auth_method auth.auth_method NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    status auth.user_status DEFAULT 'pending_verification',
    email_verified TIMESTAMP DEFAULT NULL,
    last_active_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP DEFAULT NULL,
    CONSTRAINT valid_username CHECK (username ~* '^[A-Za-z0-9._-]{3,50}$'),
    CONSTRAINT valid_email CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$')
);

-- Enhanced Profiles Table
CREATE TABLE IF NOT EXISTS auth.profiles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL UNIQUE,
    full_name VARCHAR(255) NOT NULL,
    profile_picture_url VARCHAR(255),
    cover_photo_url VARCHAR(255),
    birth_date DATE NOT NULL,
    gender VARCHAR(50),
    location VARCHAR(255),
    timezone VARCHAR(50),
    language VARCHAR(10) DEFAULT 'en',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP DEFAULT NULL,
    FOREIGN KEY (user_id) REFERENCES auth.users(id) ON DELETE CASCADE,
    CONSTRAINT valid_age CHECK (birth_date <= CURRENT_DATE - INTERVAL '13 years')
);

-- Enhanced Contact Table
CREATE TABLE IF NOT EXISTS auth.contact (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    phone_number VARCHAR(20),
    secondary_email VARCHAR(255),
    emergency_contact_name VARCHAR(255),
    emergency_contact_phone VARCHAR(20),
    address_line1 VARCHAR(255),
    address_line2 VARCHAR(255),
    city VARCHAR(100),
    state VARCHAR(100),
    postal_code VARCHAR(20),
    country VARCHAR(100),
    phone_verified TIMESTAMP DEFAULT NULL,
    secondary_email_verified TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP DEFAULT NULL,
    FOREIGN KEY (user_id) REFERENCES auth.users(id) ON DELETE CASCADE,
    CONSTRAINT valid_phone CHECK (phone_number ~* '^\+?[1-9]\d{1,14}$'),
    CONSTRAINT valid_secondary_email CHECK (secondary_email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$')
);

-- Enhanced OAuth Users Table
CREATE TABLE IF NOT EXISTS auth.oauth_users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    provider auth.oauth_provider NOT NULL,
    provider_user_id VARCHAR(255) NOT NULL,
    access_token TEXT NOT NULL,
    refresh_token TEXT,
    expires_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP DEFAULT NULL,
    FOREIGN KEY (user_id) REFERENCES auth.users(id) ON DELETE CASCADE,
);

-- Enhanced Local Users Table
CREATE TABLE IF NOT EXISTS auth.local_users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL UNIQUE,
    password VARCHAR(512) NOT NULL,
    last_password_change TIMESTAMP WITH TIME ZONE,
    password_history JSONB[], -- Store previous password hashes
    force_password_change BOOLEAN DEFAULT FALSE,
    password_expires_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP DEFAULT NULL,
    FOREIGN KEY (user_id) REFERENCES auth.users(id) ON DELETE CASCADE
);

-- New MFA Table
CREATE TABLE IF NOT EXISTS security.mfa_settings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL UNIQUE,
    mfa_enabled BOOLEAN DEFAULT FALSE,
    mfa_type auth.mfa_type,
    mfa_secret TEXT,
    backup_codes TEXT[],
    last_mfa_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES auth.users(id) ON DELETE CASCADE
);

-- Enhanced Security Status Table
CREATE TABLE IF NOT EXISTS security.account_security_status (
    status_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL UNIQUE,
    risk_level security.risk_level DEFAULT 'low',
    lock_reason VARCHAR(50),
    locked_until TIMESTAMP WITH TIME ZONE,
    failed_login_attempts INTEGER DEFAULT 0,
    failed_mfa_attempts INTEGER DEFAULT 0,
    suspicious_activity_count INTEGER DEFAULT 0,
    failed_login_reset_at TIMESTAMP WITH TIME ZONE,
    last_login_at TIMESTAMP WITH TIME ZONE,
    last_login_ip INET,
    known_devices JSONB[],
    trusted_locations JSONB[],
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES auth.users(id) ON DELETE CASCADE
);

-- Enhanced User Access Logs
CREATE TABLE IF NOT EXISTS security.user_access_logs (
    log_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    session_id UUID,
    auth_method auth.auth_method NOT NULL,
    auth_provider auth.oauth_provider,
    ip_address INET,
    user_agent TEXT,
    device_info JSONB,
    location_info JSONB,
    login_successful BOOLEAN NOT NULL,
    mfa_used BOOLEAN DEFAULT FALSE,
    mfa_type auth.mfa_type,
    risk_score INTEGER,
    metadata JSONB DEFAULT '{}',
    failure_reason VARCHAR(255),
    login_timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES auth.users(id) ON DELETE CASCADE,
    FOREIGN KEY (session_id) REFERENCES session.sessions(id) ON DELETE SET NULL
);

-- Enhanced Session Management
CREATE TABLE IF NOT EXISTS session.sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    refresh_token VARCHAR(512),
    device_id VARCHAR(255),
    device_type VARCHAR(50),
    ip_address INET,
    user_agent TEXT,
    is_mfa_completed BOOLEAN DEFAULT FALSE,
    last_activity TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    blacklisted_at TIMESTAMP WITH TIME ZONE,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES auth.users(id) ON DELETE CASCADE
);

-- User Preferences Table
CREATE TABLE IF NOT EXISTS auth.user_preferences (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL UNIQUE,
    email_notifications BOOLEAN DEFAULT TRUE,
    sms_notifications BOOLEAN DEFAULT FALSE,
    push_notifications BOOLEAN DEFAULT TRUE,
    two_factor_auth_enabled BOOLEAN DEFAULT FALSE,
    preferred_communication_method VARCHAR(50) DEFAULT 'email',
    marketing_preferences JSONB DEFAULT '{}',
    privacy_settings JSONB DEFAULT '{}',
    theme_preference VARCHAR(20) DEFAULT 'light',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES auth.users(id) ON DELETE CASCADE
);

-- Notification Settings and History
CREATE TABLE IF NOT EXISTS notification.notification_templates (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    type notification.notification_type NOT NULL,
    name VARCHAR(100) NOT NULL,
    subject VARCHAR(255),
    content TEXT NOT NULL,
    variables JSONB,
    active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS notification.notification_history (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    template_id UUID NOT NULL,
    type notification.notification_type NOT NULL,
    content TEXT NOT NULL,
    sent_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    delivered_at TIMESTAMP WITH TIME ZONE,
    read_at TIMESTAMP WITH TIME ZONE,
    status VARCHAR(50),
    metadata JSONB,
    FOREIGN KEY (user_id) REFERENCES auth.users(id) ON DELETE CASCADE,
    FOREIGN KEY (template_id) REFERENCES notification.notification_templates(id)
);

-- Rate Limiting Table
CREATE TABLE IF NOT EXISTS security.rate_limits (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID,
    ip_address INET NOT NULL,
    endpoint VARCHAR(255) NOT NULL,
    request_count INTEGER DEFAULT 1,
    window_start TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES auth.users(id) ON DELETE CASCADE
);

-- Create indexes for performance
CREATE INDEX idx_users_email ON auth.users(email);
CREATE INDEX idx_users_username ON auth.users(username);
CREATE INDEX idx_oauth_provider_user ON auth.oauth_users(provider, provider_user_id);
CREATE INDEX idx_access_logs_user_timestamp ON security.user_access_logs(user_id, login_timestamp);
CREATE INDEX idx_sessions_user ON session.sessions(user_id);
CREATE INDEX idx_rate_limits_ip_endpoint ON security.rate_limits(ip_address, endpoint);
CREATE INDEX idx_notification_history_user ON notification.notification_history(user_id);

-- Add triggers for automatic timestamp updates
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Apply timestamp triggers to all tables with updated_at
DO $$ 
DECLARE 
    t record;
BEGIN
    FOR t IN 
        SELECT table_schema, table_name 
        FROM information_schema.columns 
        WHERE column_name = 'updated_at'
        AND table_schema IN ('auth', 'security', 'session', 'notification')
    LOOP
        EXECUTE format('CREATE TRIGGER set_%I_%I_updated_at
                       BEFORE UPDATE ON %I.%I
                       FOR EACH ROW
                       EXECUTE FUNCTION update_updated_at_column()',
                       t.table_schema, t.table_name, t.table_schema, t.table_name);
    END LOOP;
END $$;