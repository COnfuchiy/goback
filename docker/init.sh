#!/bin/bash
set -e
export PGPASSWORD="$DB_PASS"

psql -v ON_ERROR_STOP=1 --username "$DB_USER" --dbname="$DB_NAME"<<-EOSQL
   CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
EOSQL