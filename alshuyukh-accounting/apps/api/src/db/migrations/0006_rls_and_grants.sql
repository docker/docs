-- 0006 — Row-Level Security and runtime privileges.
--
-- The API connects as alshuyukh_app (NOSUPERUSER, NOBYPASSRLS). Every
-- tenant-scoped query runs in a transaction that sets app.tenant_id and
-- app.user_id. These policies make PostgreSQL itself reject cross-tenant
-- access, even if application code forgets a WHERE tenant_id = ... filter.
--
-- The schema owner (migrations) is not subject to these policies.

-- tenants -----------------------------------------------------------------
ALTER TABLE tenants ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenants_select ON tenants FOR SELECT USING (
  id = app_current_tenant_id()
  OR id IN (SELECT ut.tenant_id FROM user_tenants ut
            WHERE ut.user_id = app_current_user_id() AND ut.status = 'ACTIVE')
);
CREATE POLICY tenants_insert ON tenants FOR INSERT WITH CHECK (id = app_current_tenant_id());
CREATE POLICY tenants_update ON tenants FOR UPDATE
  USING (id = app_current_tenant_id()) WITH CHECK (id = app_current_tenant_id());

-- tenant_settings ---------------------------------------------------------
ALTER TABLE tenant_settings ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_settings_isolation ON tenant_settings
  USING (tenant_id = app_current_tenant_id()) WITH CHECK (tenant_id = app_current_tenant_id());

-- user_tenants: a user can see their own memberships (for login and the
-- tenant switcher) plus all memberships of the active tenant.
ALTER TABLE user_tenants ENABLE ROW LEVEL SECURITY;
CREATE POLICY user_tenants_select ON user_tenants FOR SELECT USING (
  tenant_id = app_current_tenant_id() OR user_id = app_current_user_id()
);
CREATE POLICY user_tenants_insert ON user_tenants FOR INSERT
  WITH CHECK (tenant_id = app_current_tenant_id());
CREATE POLICY user_tenants_update ON user_tenants FOR UPDATE
  USING (tenant_id = app_current_tenant_id()) WITH CHECK (tenant_id = app_current_tenant_id());

-- roles: system roles are readable by everyone, never writable at runtime.
ALTER TABLE roles ENABLE ROW LEVEL SECURITY;
CREATE POLICY roles_select ON roles FOR SELECT USING (
  tenant_id IS NULL OR tenant_id = app_current_tenant_id()
);
CREATE POLICY roles_insert ON roles FOR INSERT WITH CHECK (tenant_id = app_current_tenant_id());
CREATE POLICY roles_update ON roles FOR UPDATE
  USING (tenant_id = app_current_tenant_id()) WITH CHECK (tenant_id = app_current_tenant_id());

-- role_permissions follow the visibility and ownership of their role.
ALTER TABLE role_permissions ENABLE ROW LEVEL SECURITY;
CREATE POLICY role_permissions_select ON role_permissions FOR SELECT USING (
  EXISTS (SELECT 1 FROM roles r WHERE r.id = role_id)
);
CREATE POLICY role_permissions_insert ON role_permissions FOR INSERT WITH CHECK (
  EXISTS (SELECT 1 FROM roles r WHERE r.id = role_id AND r.tenant_id = app_current_tenant_id())
);
CREATE POLICY role_permissions_delete ON role_permissions FOR DELETE USING (
  EXISTS (SELECT 1 FROM roles r WHERE r.id = role_id AND r.tenant_id = app_current_tenant_id())
);

-- Plain tenant-scoped tables -----------------------------------------------
ALTER TABLE user_roles ENABLE ROW LEVEL SECURITY;
CREATE POLICY user_roles_isolation ON user_roles
  USING (tenant_id = app_current_tenant_id()) WITH CHECK (tenant_id = app_current_tenant_id());

ALTER TABLE companies ENABLE ROW LEVEL SECURITY;
CREATE POLICY companies_isolation ON companies
  USING (tenant_id = app_current_tenant_id()) WITH CHECK (tenant_id = app_current_tenant_id());

ALTER TABLE branches ENABLE ROW LEVEL SECURITY;
CREATE POLICY branches_isolation ON branches
  USING (tenant_id = app_current_tenant_id()) WITH CHECK (tenant_id = app_current_tenant_id());

ALTER TABLE warehouses ENABLE ROW LEVEL SECURITY;
CREATE POLICY warehouses_isolation ON warehouses
  USING (tenant_id = app_current_tenant_id()) WITH CHECK (tenant_id = app_current_tenant_id());

-- audit_logs: read own tenant; insert for own tenant or pre-tenant events.
ALTER TABLE audit_logs ENABLE ROW LEVEL SECURITY;
CREATE POLICY audit_logs_select ON audit_logs FOR SELECT USING (tenant_id = app_current_tenant_id());
CREATE POLICY audit_logs_insert ON audit_logs FOR INSERT WITH CHECK (
  tenant_id IS NULL OR tenant_id = app_current_tenant_id()
);

-- Runtime privileges --------------------------------------------------------
-- No DELETE on business tables: records are soft-deleted (deleted_at).
-- Join tables (user_roles, role_permissions) allow DELETE; every change is audited.
GRANT USAGE ON SCHEMA public TO alshuyukh_app;
GRANT SELECT ON permissions TO alshuyukh_app;
GRANT SELECT, INSERT, UPDATE ON
  tenants, tenant_settings, users, user_tenants, user_sessions,
  roles, companies, branches, warehouses
TO alshuyukh_app;
GRANT SELECT, INSERT, DELETE ON user_roles, role_permissions TO alshuyukh_app;
GRANT SELECT, INSERT ON audit_logs TO alshuyukh_app;
