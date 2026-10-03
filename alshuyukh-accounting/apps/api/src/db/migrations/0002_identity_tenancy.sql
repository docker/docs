-- 0002 — Identity and multi-tenancy.
-- A tenant is the isolation boundary (one subscribing organization).
-- Users are global identities; membership in a tenant is in user_tenants.

CREATE TABLE tenants (
  id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  name        text NOT NULL CHECK (length(btrim(name)) BETWEEN 2 AND 200),
  slug        citext NOT NULL CHECK (slug ~ '^[a-z0-9][a-z0-9-]{1,62}$'),
  status      text NOT NULL DEFAULT 'ACTIVE' CHECK (status IN ('ACTIVE', 'SUSPENDED', 'CANCELLED')),
  created_at  timestamptz NOT NULL DEFAULT now(),
  updated_at  timestamptz NOT NULL DEFAULT now(),
  deleted_at  timestamptz
);
CREATE UNIQUE INDEX tenants_slug_uq ON tenants (slug);
CREATE INDEX tenants_status_idx ON tenants (status);
CREATE TRIGGER tenants_updated_at BEFORE UPDATE ON tenants
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TABLE tenant_settings (
  tenant_id               uuid PRIMARY KEY REFERENCES tenants (id) ON DELETE RESTRICT,
  default_currency        char(3) NOT NULL DEFAULT 'SAR' CHECK (default_currency ~ '^[A-Z]{3}$'),
  timezone                text NOT NULL DEFAULT 'Asia/Riyadh',
  locale                  text NOT NULL DEFAULT 'ar' CHECK (locale IN ('ar', 'en')),
  fiscal_year_start_month smallint NOT NULL DEFAULT 1 CHECK (fiscal_year_start_month BETWEEN 1 AND 12),
  date_format             text NOT NULL DEFAULT 'YYYY-MM-DD',
  extra                   jsonb NOT NULL DEFAULT '{}'::jsonb,
  created_at              timestamptz NOT NULL DEFAULT now(),
  updated_at              timestamptz NOT NULL DEFAULT now()
);
CREATE TRIGGER tenant_settings_updated_at BEFORE UPDATE ON tenant_settings
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TABLE users (
  id                    uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  email                 citext NOT NULL CHECK (position('@' IN email) > 1),
  password_hash         text NOT NULL,
  full_name             text NOT NULL CHECK (length(btrim(full_name)) BETWEEN 2 AND 200),
  phone                 text,
  status                text NOT NULL DEFAULT 'ACTIVE' CHECK (status IN ('ACTIVE', 'DISABLED')),
  -- Platform (Super Admin) access is a platform flag, not a tenant role.
  -- The /admin panel that uses it is delivered in Phase 9.
  is_platform_admin     boolean NOT NULL DEFAULT false,
  must_change_password  boolean NOT NULL DEFAULT false,
  failed_login_attempts integer NOT NULL DEFAULT 0 CHECK (failed_login_attempts >= 0),
  locked_until          timestamptz,
  last_login_at         timestamptz,
  password_changed_at   timestamptz NOT NULL DEFAULT now(),
  created_at            timestamptz NOT NULL DEFAULT now(),
  updated_at            timestamptz NOT NULL DEFAULT now(),
  deleted_at            timestamptz
);
CREATE UNIQUE INDEX users_email_uq ON users (email);
CREATE TRIGGER users_updated_at BEFORE UPDATE ON users
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TABLE user_tenants (
  id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id   uuid NOT NULL REFERENCES tenants (id) ON DELETE RESTRICT,
  user_id     uuid NOT NULL REFERENCES users (id) ON DELETE RESTRICT,
  status      text NOT NULL DEFAULT 'ACTIVE' CHECK (status IN ('ACTIVE', 'INVITED', 'DISABLED')),
  is_owner    boolean NOT NULL DEFAULT false,
  joined_at   timestamptz NOT NULL DEFAULT now(),
  created_at  timestamptz NOT NULL DEFAULT now(),
  updated_at  timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT user_tenants_tenant_user_uq UNIQUE (tenant_id, user_id)
);
CREATE INDEX user_tenants_user_idx ON user_tenants (user_id);
CREATE TRIGGER user_tenants_updated_at BEFORE UPDATE ON user_tenants
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- Refresh-token sessions. Only the SHA-256 hash of the token is stored.
CREATE TABLE user_sessions (
  id            uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id       uuid NOT NULL REFERENCES users (id) ON DELETE RESTRICT,
  tenant_id     uuid NOT NULL REFERENCES tenants (id) ON DELETE RESTRICT,
  token_hash    text NOT NULL,
  expires_at    timestamptz NOT NULL,
  revoked_at    timestamptz,
  replaced_by   uuid REFERENCES user_sessions (id),
  ip_address    inet,
  user_agent    text,
  created_at    timestamptz NOT NULL DEFAULT now(),
  last_used_at  timestamptz
);
CREATE UNIQUE INDEX user_sessions_token_hash_uq ON user_sessions (token_hash);
CREATE INDEX user_sessions_user_idx ON user_sessions (user_id) WHERE revoked_at IS NULL;
