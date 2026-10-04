#!/bin/sh
# Runs once, when the PostgreSQL volume is first created.
# Creates the schema owner (migrations) and the restricted application role
# (NOBYPASSRLS, so row-level security isolates tenants).
set -eu
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname postgres \
  -v owner_pw="$OWNER_DB_PASSWORD" -v app_pw="$APP_DB_PASSWORD" <<'SQL'
CREATE ROLE alshuyukh_owner LOGIN PASSWORD :'owner_pw' NOSUPERUSER NOCREATEROLE NOCREATEDB;
CREATE ROLE alshuyukh_app LOGIN PASSWORD :'app_pw' NOSUPERUSER NOCREATEROLE NOCREATEDB NOBYPASSRLS;
CREATE DATABASE alshuyukh OWNER alshuyukh_owner;
REVOKE ALL ON DATABASE alshuyukh FROM PUBLIC;
GRANT CONNECT ON DATABASE alshuyukh TO alshuyukh_app;
SQL
