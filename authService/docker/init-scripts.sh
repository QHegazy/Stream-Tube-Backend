#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "auth" --dbname "auth" <<-EOSQL
    DO
    \$do\$
    BEGIN
        IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'auth') THEN
            CREATE USER auth WITH REPLICATION PASSWORD 'auth';
        END IF;
    END
    \$do\$;
EOSQL