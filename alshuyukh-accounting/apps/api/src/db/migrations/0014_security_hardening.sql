-- Phase 10: security hardening.

-- Login lockout per (account, client IP). A lock earned from one address no
-- longer locks the real owner out from theirs, so failed logins cannot be used
-- to deny someone access. (users.failed_login_attempts/locked_until are no
-- longer used.)
CREATE TABLE login_failures (
  user_id       uuid NOT NULL REFERENCES users (id),
  ip            text NOT NULL,
  failures      integer NOT NULL DEFAULT 0,
  locked_until  timestamptz,
  updated_at    timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY (user_id, ip)
);
GRANT SELECT, INSERT, UPDATE, DELETE ON login_failures TO alshuyukh_app;

-- Set only by the server while creating a new organization (sign-up or by the
-- platform administrator), so the first subscription can be written.
CREATE OR REPLACE FUNCTION app_is_provisioning() RETURNS boolean LANGUAGE sql STABLE AS $$
  SELECT COALESCE(current_setting('app.provisioning', true), '') = 'on'
$$;

-- Billing data: a tenant reads its own; only the platform (or provisioning)
-- writes it. Before this, a tenant context could in principle extend its own
-- period or change its plan.
DROP POLICY subscriptions_isolation ON subscriptions;
CREATE POLICY subscriptions_select ON subscriptions FOR SELECT USING (tenant_id = app_current_tenant_id());
CREATE POLICY subscriptions_insert ON subscriptions FOR INSERT
  WITH CHECK (tenant_id = app_current_tenant_id() AND (app_is_platform_admin() OR app_is_provisioning()));
CREATE POLICY subscriptions_update ON subscriptions FOR UPDATE
  USING (tenant_id = app_current_tenant_id() AND app_is_platform_admin())
  WITH CHECK (tenant_id = app_current_tenant_id() AND app_is_platform_admin());

DROP POLICY subscription_items_isolation ON subscription_items;
CREATE POLICY subscription_items_select ON subscription_items FOR SELECT USING (tenant_id = app_current_tenant_id());
CREATE POLICY subscription_items_insert ON subscription_items FOR INSERT
  WITH CHECK (tenant_id = app_current_tenant_id() AND (app_is_platform_admin() OR app_is_provisioning()));
CREATE POLICY subscription_items_delete ON subscription_items FOR DELETE
  USING (tenant_id = app_current_tenant_id() AND (app_is_platform_admin() OR app_is_provisioning()));

DROP POLICY billing_events_isolation ON billing_events;
CREATE POLICY billing_events_select ON billing_events FOR SELECT USING (tenant_id = app_current_tenant_id());
CREATE POLICY billing_events_insert ON billing_events FOR INSERT
  WITH CHECK (tenant_id = app_current_tenant_id()
              AND (app_is_platform_admin() OR app_is_provisioning() OR event_type = 'PLAN_CHANGE_REQUESTED'));

DROP POLICY tenant_feature_flags_isolation ON tenant_feature_flags;
CREATE POLICY tenant_feature_flags_select ON tenant_feature_flags FOR SELECT USING (tenant_id = app_current_tenant_id());
CREATE POLICY tenant_feature_flags_write ON tenant_feature_flags FOR ALL
  USING (tenant_id = app_current_tenant_id() AND app_is_platform_admin())
  WITH CHECK (tenant_id = app_current_tenant_id() AND app_is_platform_admin());

-- An organization may rename itself but not change its own status, and no
-- user may make themselves a platform administrator. Enforced for the API's
-- role; maintenance through the owner role (CLI) is unaffected.
CREATE OR REPLACE FUNCTION tenants_status_guard() RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
  IF NEW.status IS DISTINCT FROM OLD.status AND current_user = 'alshuyukh_app' AND NOT app_is_platform_admin() THEN
    RAISE EXCEPTION 'Only the platform administrator can change an organization''s status';
  END IF;
  RETURN NEW;
END $$;
CREATE TRIGGER tenants_status_guard BEFORE UPDATE ON tenants FOR EACH ROW EXECUTE FUNCTION tenants_status_guard();

CREATE OR REPLACE FUNCTION users_platform_guard() RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
  IF current_user = 'alshuyukh_app' AND NOT app_is_platform_admin()
     AND (NEW.is_platform_admin IS DISTINCT FROM OLD.is_platform_admin OR NEW.status IS DISTINCT FROM OLD.status) THEN
    RAISE EXCEPTION 'Only the platform administrator can change this account''s status or platform access';
  END IF;
  RETURN NEW;
END $$;
CREATE TRIGGER users_platform_guard BEFORE UPDATE ON users FOR EACH ROW EXECUTE FUNCTION users_platform_guard();
