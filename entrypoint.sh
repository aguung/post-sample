#!/bin/sh

# Construct DATABASE_URL if not set
if [ -z "$DATABASE_URL" ]; then
  # Default values if env vars are missing
  DB_HOST=${DB_HOST:-localhost}
  DB_PORT=${DB_PORT:-5432}
  DB_USER=${DB_USER:-postgres}
  DB_PASSWORD=${DB_PASSWORD:-postgres}
  DB_NAME=${DB_NAME:-post_db}
  DB_SSLMODE=${DB_SSLMODE:-disable}

  export DATABASE_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}"
fi

# Run migrations
echo "Running migrations..."
/usr/local/bin/atlas migrate apply \
  --url "$DATABASE_URL" \
  --dir "file://migrations" \
  --revisions-schema "public"

# Check migration status
if [ $? -ne 0 ]; then
  echo "Migration failed!"
  exit 1
fi

echo "Starting application..."
exec "$@"
