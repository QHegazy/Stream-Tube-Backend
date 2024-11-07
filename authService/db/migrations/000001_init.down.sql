-- Drop Trigger Function
DROP FUNCTION IF EXISTS update_updated_at_column;

-- Drop all Triggers related to updated_at column
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
        EXECUTE format('DROP TRIGGER IF EXISTS set_%I_%I_updated_at ON %I.%I', 
                       t.table_schema, t.table_name, t.table_schema, t.table_name);
    END LOOP;
END $$;

-- Drop Indexes
DROP INDEX IF EXISTS auth.idx_users_email;
DROP INDEX IF EXISTS auth.idx_users_username;
DROP INDEX IF EXISTS auth.idx_oauth_provider_user;
DROP INDEX IF EXISTS security.idx_access_logs_user_timestamp;
DROP INDEX IF EXISTS session.idx_sessions_user;
DROP INDEX IF EXISTS security.idx_rate_limits_ip_endpoint;
DROP INDEX IF EXISTS notification.idx_notification_history_user;

-- Drop Tables
DROP TABLE IF EXISTS notification.notification_history CASCADE;
DROP TABLE IF EXISTS notification.notification_templates CASCADE;
DROP TABLE IF EXISTS auth.user_preferences CASCADE;
DROP TABLE IF EXISTS security.user_access_logs CASCADE;
DROP TABLE IF EXISTS security.account_security_status CASCADE;
DROP TABLE IF EXISTS security.mfa_settings CASCADE;
DROP TABLE IF EXISTS session.sessions CASCADE;
DROP TABLE IF EXISTS auth.local_users CASCADE;
DROP TABLE IF EXISTS auth.oauth_users CASCADE;
DROP TABLE IF EXISTS auth.contact CASCADE;
DROP TABLE IF EXISTS auth.profiles CASCADE;
DROP TABLE IF EXISTS auth.users CASCADE;

-- Drop Enums
DROP TYPE IF EXISTS notification.notification_type;
DROP TYPE IF EXISTS auth.mfa_type;
DROP TYPE IF EXISTS auth.auth_method;
DROP TYPE IF EXISTS auth.oauth_provider;
DROP TYPE IF EXISTS auth.user_status;
DROP TYPE IF EXISTS security.risk_level;

-- Drop Schemas
DROP SCHEMA IF EXISTS notification CASCADE;
DROP SCHEMA IF EXISTS security CASCADE;
DROP SCHEMA IF EXISTS session CASCADE;
DROP SCHEMA IF EXISTS auth CASCADE;

-- Drop UUID Extension
DROP EXTENSION IF EXISTS "uuid-ossp";
