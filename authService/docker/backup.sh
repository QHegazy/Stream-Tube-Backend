#!/bin/bash

# Database backup script
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
BACKUP_NAME="auth_db_backup_${TIMESTAMP}.sql"
BACKUP_PATH="/backups"

# Ensure backup directory exists
mkdir -p $BACKUP_PATH

# Create backup
pg_dump -h localhost \
        -U $POSTGRES_USER \
        -d $POSTGRES_DB \
        -F c \
        -f "${BACKUP_PATH}/${BACKUP_NAME}"

# Remove old backups
find $BACKUP_PATH -type f -name "auth_db_backup_*.sql" -mtime +$BACKUP_RETENTION_DAYS -delete