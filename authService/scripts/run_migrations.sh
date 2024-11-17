docker logs auth_postgres_primary
#!/bin/bash

# Load environment variables from the .env file
if [ -f ".env" ]; then
  export $(cat .env | grep -v '^#' | xargs)
else
  echo "Error: .env file not found."
  exit 1
fi

# Check if migrate is installed
if ! command -v migrate &> /dev/null; then
  echo "Error: migrate command not found. Please install the migrate tool."
  exit 1
fi

# Configuration: Use environment variables loaded from .env
MIGRATION_PATH="./db/migrations"  # Path to migration files
DATABASE_URL="postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=require"

# Apply migrations (run the up migration)
apply_migrations() {
  echo "Applying migrations..."
  migrate -path $MIGRATION_PATH -database $DATABASE_URL up
  if [ $? -eq 0 ]; then
    echo "Migrations applied successfully."
  else
    echo "Error applying migrations."
    exit 1
  fi
}

# Rollback migrations (run the down migration)
rollback_migrations() {
  echo "Rolling back migrations..."
  migrate -path $MIGRATION_PATH -database $DATABASE_URL down
  if [ $? -eq 0 ]; then
    echo "Migrations rolled back successfully."
  else
    echo "Error rolling back migrations."
    exit 1
  fi
}

# Show migration version (to see current state of the migrations)
show_version() {
  echo "Current migration version:"
  migrate -path $MIGRATION_PATH -database $DATABASE_URL version
}

# Show usage information
usage() {
  echo "Usage: $0 {apply|rollback|version}"
  echo "  apply   - Apply the up migrations"
  echo "  rollback- Rollback the down migrations"
  echo "  version - Show the current migration version"
}

# Main execution: Check arguments and run the appropriate function
if [ $# -eq 0 ]; then
  usage
  exit 1
fi

case "$1" in
  apply)
    apply_migrations
    ;;
  rollback)
    rollback_migrations
    ;;
  version)
    show_version
    ;;
  *)
    usage
    exit 1
    ;;
esac
