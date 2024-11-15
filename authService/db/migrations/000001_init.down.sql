DO $$ 
DECLARE 
    t record;
BEGIN
    FOR t IN 
        SELECT trigger_name, event_object_schema, event_object_table
        FROM information_schema.triggers 
        WHERE trigger_name LIKE 'set_%_updated_at'
    LOOP
        EXECUTE format('DROP TRIGGER IF EXISTS %I ON %I.%I CASCADE',
                      t.trigger_name,
                      t.event_object_schema,
                      t.event_object_table);
    END LOOP;
END $$;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column() CASCADE;

-- Drop indexes
DROP INDEX IF EXISTS auth.idx_users_email CASCADE;
DROP INDEX IF EXISTS auth.idx_users_username CASCADE;
DROP INDEX IF EXISTS auth.idx_oauth_provider_user CASCADE;
DROP INDEX IF EXISTS security.idx_access_logs_user_timestamp CASCADE;
DROP INDEX IF EXISTS session.idx_sessions_user CASCADE;
DROP INDEX IF EXISTS security.idx_rate_limits_ip_endpoint CASCADE;
DROP INDEX IF EXISTS notification.idx_notification_history_user CASCADE;

-- Drop tables in reverse order of dependencies
DROP TABLE IF EXISTS security.rate_limits CASCADE;
DROP TABLE IF EXISTS notification.notification_history CASCADE;
DROP TABLE IF EXISTS notification.notification_templates CASCADE;
DROP TABLE IF EXISTS auth.user_preferences CASCADE;
DROP TABLE IF EXISTS security.user_access_logs CASCADE;
DROP TABLE IF EXISTS security.account_security_status CASCADE;
DROP TABLE IF EXISTS security.mfa_settings CASCADE;
DROP TABLE IF EXISTS session.sessions CASCADE;
DROP TABLE IF EXISTS auth.local_users CASCADE;
DROP TABLE IF EXISTS auth.oauth_users CASCADE;
DROP TABLE IF EXISTS profile.contact CASCADE;
DROP TABLE IF EXISTS profile.profiles CASCADE;
DROP TABLE IF EXISTS auth.users CASCADE;

-- Drop custom types
DROP TYPE IF EXISTS auth.user_status CASCADE;
DROP TYPE IF EXISTS auth.oauth_provider CASCADE;
DROP TYPE IF EXISTS auth.auth_method CASCADE;
DROP TYPE IF EXISTS auth.mfa_type CASCADE;
DROP TYPE IF EXISTS notification.notification_type CASCADE;
DROP TYPE IF EXISTS security.risk_level CASCADE;
DROP TYPE IF EXISTS profile.gender CASCADE;

-- Drop schemas
DROP SCHEMA IF EXISTS auth CASCADE;
DROP SCHEMA IF EXISTS session CASCADE;
DROP SCHEMA IF EXISTS security CASCADE;
DROP SCHEMA IF EXISTS notification CASCADE;
DROP SCHEMA IF EXISTS profile CASCADE;

-- Drop UUID Extension
DROP EXTENSION IF EXISTS "uuid-ossp" CASCADE;