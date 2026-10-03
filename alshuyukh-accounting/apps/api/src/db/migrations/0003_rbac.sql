-- 0003 — Role-based access control.
-- System roles (tenant_id IS NULL) are shared templates managed by the
-- application catalog. Tenants can add custom roles (tenant_id = tenant).

CREATE TABLE permissions (
  id              uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  code            text NOT NULL CHECK (code ~ '^[a-z_]+\.[a-z_]+$'),
  module          text NOT NULL,
  description_ar  text NOT NULL,
  description_en  text NOT NULL,
  created_at      timestamptz NOT NULL DEFAULT now()
);
CREATE UNIQUE INDEX permissions_code_uq ON permissions (code);

CREATE TABLE roles (
  id           uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id    uuid REFERENCES tenants (id) ON DELETE RESTRICT,
  code         text NOT NULL CHECK (code ~ '^[A-Z][A-Z0-9_]{1,62}$'),
  name_ar      text NOT NULL,
  name_en      text NOT NULL,
  description  text,
  is_system    boolean NOT NULL DEFAULT false,
  created_at   timestamptz NOT NULL DEFAULT now(),
  updated_at   timestamptz NOT NULL DEFAULT now(),
  deleted_at   timestamptz,
  CONSTRAINT roles_system_scope_ck CHECK (
    (is_system AND tenant_id IS NULL) OR (NOT is_system AND tenant_id IS NOT NULL)
  )
);
-- A role code is unique among system roles, and unique per tenant.
CREATE UNIQUE INDEX roles_system_code_uq ON roles (code) WHERE tenant_id IS NULL;
CREATE UNIQUE INDEX roles_tenant_code_uq ON roles (tenant_id, code)
  WHERE tenant_id IS NOT NULL AND deleted_at IS NULL;
CREATE INDEX roles_tenant_idx ON roles (tenant_id);
CREATE TRIGGER roles_updated_at BEFORE UPDATE ON roles
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TABLE role_permissions (
  role_id        uuid NOT NULL REFERENCES roles (id) ON DELETE CASCADE,
  permission_id  uuid NOT NULL REFERENCES permissions (id) ON DELETE RESTRICT,
  created_at     timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY (role_id, permission_id)
);
CREATE INDEX role_permissions_permission_idx ON role_permissions (permission_id);

CREATE TABLE user_roles (
  id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id   uuid NOT NULL,
  user_id     uuid NOT NULL,
  role_id     uuid NOT NULL REFERENCES roles (id) ON DELETE RESTRICT,
  created_by  uuid REFERENCES users (id),
  created_at  timestamptz NOT NULL DEFAULT now(),
  -- The user must be a member of the same tenant.
  CONSTRAINT user_roles_membership_fk FOREIGN KEY (tenant_id, user_id)
    REFERENCES user_tenants (tenant_id, user_id) ON DELETE RESTRICT,
  CONSTRAINT user_roles_uq UNIQUE (tenant_id, user_id, role_id)
);
CREATE INDEX user_roles_role_idx ON user_roles (role_id);

-- The assigned role must be a system role or a role of the same tenant.
CREATE OR REPLACE FUNCTION check_user_role_tenant() RETURNS trigger
LANGUAGE plpgsql AS $$
DECLARE
  role_tenant uuid;
  role_found  boolean;
BEGIN
  SELECT true, tenant_id INTO role_found, role_tenant
    FROM roles WHERE id = NEW.role_id AND deleted_at IS NULL;
  IF role_found IS NULL OR (role_tenant IS NOT NULL AND role_tenant <> NEW.tenant_id) THEN
    RAISE EXCEPTION 'role % is not available in tenant %', NEW.role_id, NEW.tenant_id
      USING ERRCODE = 'foreign_key_violation';
  END IF;
  RETURN NEW;
END;
$$;
CREATE TRIGGER user_roles_tenant_check BEFORE INSERT OR UPDATE ON user_roles
  FOR EACH ROW EXECUTE FUNCTION check_user_role_tenant();
