-- Revert the triggers first
DROP TRIGGER IF EXISTS set_timestamp ON users;
DROP TRIGGER IF EXISTS set_timestamp ON profiles;
DROP TRIGGER IF EXISTS set_timestamp ON oauth_users;
DROP TRIGGER IF EXISTS set_timestamp ON local_users;

DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop Tables in Reverse Order of Dependencies
DROP TABLE IF EXISTS local_users;
DROP TABLE IF EXISTS oauth_users;
DROP TABLE IF EXISTS profiles;
DROP TABLE IF EXISTS users;

-- Drop Enums
DROP TYPE IF EXISTS auth.user_status;
DROP TYPE IF EXISTS auth.oauth_provider;
DROP TYPE IF EXISTS auth.contact_type;
DROP TYPE IF EXISTS auth.auth_method;

-- Drop Schemas
DROP SCHEMA IF EXISTS auth CASCADE;
DROP SCHEMA IF EXISTS session CASCADE;
