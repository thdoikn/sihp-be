#!/bin/sh
set -e

DB_HOST="${DB_HOST:-db}"
DB_PORT="${DB_PORT:-5432}"
DB_USER="${DB_USER:-postgres}"
MIGRATE_ENV="${MIGRATE_ENV:-docker}"
SKIP_DB_SETUP="${SKIP_DB_SETUP:-false}"

if [ "$SKIP_DB_SETUP" = "true" ]; then
  echo "SKIP_DB_SETUP=true, skipping migrations and seed-admin"
  exec "$@"
fi

echo "Waiting for PostgreSQL at ${DB_HOST}:${DB_PORT}..."
until pg_isready -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -q; do
  sleep 1
done

echo "Running database migrations (env=${MIGRATE_ENV})..."
cd /app/migration
cat > dbconfig.runtime.yml <<EOF
docker:
  dialect: postgres
  datasource: host=${DB_HOST} dbname=${DB_NAME:-sihp} user=${DB_USER} password=${DB_PASSWORD} sslmode=${DB_SSL_MODE:-disable} port=${DB_PORT}
  dir: files
  table: migrations
EOF
sql-migrate up -config=dbconfig.runtime.yml -env=docker

echo "Seeding admin user..."
cd /app
/usr/local/bin/seed-admin

echo "Starting application..."
exec "$@"
