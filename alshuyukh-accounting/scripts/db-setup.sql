-- Run once as a PostgreSQL superuser:
--   psql -U postgres -f scripts/db-setup.sql
--
-- Two roles:
--   alshuyukh_owner : owns the schema, runs migrations.
--   alshuyukh_app   : used by the API at runtime. NOT a superuser and has
--                     NOBYPASSRLS, so PostgreSQL Row-Level Security applies.
-- Change the passwords for any non-local environment.

DO $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'alshuyukh_owner') THEN
    CREATE ROLE alshuyukh_owner LOGIN PASSWORD 'owner_dev_password' NOSUPERUSER NOCREATEROLE;
  END IF;
  IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'alshuyukh_app') THEN
    CREATE ROLE alshuyukh_app LOGIN PASSWORD 'app_dev_password' NOSUPERUSER NOCREATEROLE NOBYPASSRLS;
  END IF;
END
$$;

SELECT 'CREATE DATABASE alshuyukh OWNER alshuyukh_owner'
WHERE NOT EXISTS (SELECT 1 FROM pg_database WHERE datname = 'alshuyukh')\gexec

SELECT 'CREATE DATABASE alshuyukh_test OWNER alshuyukh_owner'
WHERE NOT EXISTS (SELECT 1 FROM pg_database WHERE datname = 'alshuyukh_test')\gexec

GRANT CONNECT ON DATABASE alshuyukh TO alshuyukh_app;
GRANT CONNECT ON DATABASE alshuyukh_test TO alshuyukh_app;
