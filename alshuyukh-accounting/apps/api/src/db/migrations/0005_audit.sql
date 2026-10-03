-- 0005 — Audit log (append-only).
-- tenant_id is NULL only for events that happen before a tenant is known,
-- such as a failed login for an unknown e-mail address.

CREATE TABLE audit_logs (
  id           uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id    uuid REFERENCES tenants (id) ON DELETE RESTRICT,
  user_id      uuid REFERENCES users (id) ON DELETE RESTRICT,
  action       text NOT NULL CHECK (action ~ '^[A-Z_]+$'),
  entity_type  text,
  entity_id    uuid,
  old_values   jsonb,
  new_values   jsonb,
  ip_address   inet,
  user_agent   text,
  request_id   text,
  created_at   timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX audit_logs_tenant_created_idx ON audit_logs (tenant_id, created_at DESC);
CREATE INDEX audit_logs_entity_idx ON audit_logs (tenant_id, entity_type, entity_id);
CREATE INDEX audit_logs_user_idx ON audit_logs (user_id, created_at DESC);
CREATE INDEX audit_logs_action_idx ON audit_logs (tenant_id, action);

CREATE TRIGGER audit_logs_append_only BEFORE UPDATE OR DELETE ON audit_logs
  FOR EACH ROW EXECUTE FUNCTION prevent_modification();
