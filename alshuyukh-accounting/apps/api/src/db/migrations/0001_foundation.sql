-- 0001 — Foundation: extensions and shared helper functions.
-- All timestamps are stored as timestamptz (UTC). Presentation converts to the
-- company time zone (default Asia/Riyadh).

CREATE EXTENSION IF NOT EXISTS citext;

-- Keeps updated_at current on every UPDATE.
CREATE OR REPLACE FUNCTION set_updated_at() RETURNS trigger
LANGUAGE plpgsql AS $$
BEGIN
  NEW.updated_at := now();
  RETURN NEW;
END;
$$;

-- Request context. The API sets these with set_config(..., true) at the start
-- of every transaction, so they are scoped to that transaction only.
CREATE OR REPLACE FUNCTION app_current_tenant_id() RETURNS uuid
LANGUAGE sql STABLE AS $$
  SELECT NULLIF(current_setting('app.tenant_id', true), '')::uuid
$$;

CREATE OR REPLACE FUNCTION app_current_user_id() RETURNS uuid
LANGUAGE sql STABLE AS $$
  SELECT NULLIF(current_setting('app.user_id', true), '')::uuid
$$;

-- Blocks UPDATE and DELETE on append-only tables (audit logs, and later
-- posted financial records).
CREATE OR REPLACE FUNCTION prevent_modification() RETURNS trigger
LANGUAGE plpgsql AS $$
BEGIN
  RAISE EXCEPTION '% on table % is not allowed (append-only)', TG_OP, TG_TABLE_NAME
    USING ERRCODE = 'insufficient_privilege';
END;
$$;
