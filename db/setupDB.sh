#!/bin/bash
set -e


# POSTGRES_USER=cityexplorer
# POSTGRES_PASSWORD=localdevpassword
# POSTGRES_DB=city_explorer
# POSTGRES_PORT=5432

DB_CONTAINER="city-explorer-db"
DB_USER="cityexplorer"
DB_NAME="city_explorer"

echo "==> Setting up database..."

for file in db/migrations/*.up.sql; do
  echo "    Running $file"
  docker exec -i $DB_CONTAINER psql -U $DB_USER -d $DB_NAME < "$file"
done

echo "==> Done."