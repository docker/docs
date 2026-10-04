-- Phase 9: SaaS subscriptions and platform administration.
--
-- Plans, prices and limits are data edited by the platform administrator,
-- never constants in code. A limit of NULL means unlimited.

-- Platform access ------------------------------------------------------------------------
-- Set by the API (set_config('app.platform_admin', 'on', true)) only inside
-- /admin requests, after re-reading users.is_platform_admin from the database.
CREATE OR REPLACE FUNCTION app_is_platform_admin() RETURNS boolean LANGUAGE sql STABLE AS $$
  SELECT COALESCE(current_setting('app.platform_admin', true), '') = 'on'
$$;

-- Plans ----------------------------------------------------------------------------------
CREATE TABLE plans (
  id                       uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  code                     text NOT NULL UNIQUE CHECK (code ~ '^[A-Z][A-Z0-9_]{1,30}$'),
  name_ar                  text NOT NULL CHECK (length(btrim(name_ar)) BETWEEN 2 AND 100),
  name_en                  text,
  description              text,
  currency                 char(3) NOT NULL DEFAULT 'SAR' CHECK (currency ~ '^[A-Z]{3}$'),
  price_monthly            numeric(18,2) NOT NULL DEFAULT 0 CHECK (price_monthly >= 0),
  price_yearly             numeric(18,2) NOT NULL DEFAULT 0 CHECK (price_yearly >= 0),
  trial_days               integer NOT NULL DEFAULT 0 CHECK (trial_days BETWEEN 0 AND 365),
  -- Days after the period ends during which the tenant keeps full access.
  grace_days               integer NOT NULL DEFAULT 7 CHECK (grace_days BETWEEN 0 AND 90),
  max_users                integer CHECK (max_users > 0),
  max_companies            integer CHECK (max_companies > 0),
  max_branches             integer CHECK (max_branches > 0),
  max_warehouses           integer CHECK (max_warehouses > 0),
  max_products             integer CHECK (max_products > 0),
  max_invoices_per_month   integer CHECK (max_invoices_per_month > 0),
  max_storage_mb           integer CHECK (max_storage_mb > 0),
  max_api_calls_per_month  integer CHECK (max_api_calls_per_month > 0),
  is_public                boolean NOT NULL DEFAULT true,
  -- The plan new sign-ups start on (as a trial when trial_days > 0).
  is_default               boolean NOT NULL DEFAULT false,
  is_active                boolean NOT NULL DEFAULT true,
  sort_order               integer NOT NULL DEFAULT 0,
  created_at               timestamptz NOT NULL DEFAULT now(),
  updated_at               timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT plans_default_active_ck CHECK (NOT is_default OR is_active)
);
CREATE UNIQUE INDEX plans_one_default_uq ON plans (is_default) WHERE is_default;
CREATE TRIGGER plans_updated_at BEFORE UPDATE ON plans FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- The default plan new sign-ups start on. Price 0 and a trial; paid plans are
-- created by the platform administrator.
INSERT INTO plans (code, name_ar, name_en, description, trial_days, max_users, max_companies, max_branches, max_warehouses,
                   max_products, max_invoices_per_month, max_storage_mb, max_api_calls_per_month, is_public, is_default, sort_order)
VALUES ('TRIAL', 'الباقة التجريبية', 'Trial', 'تجربة كاملة المزايا لمدة محدودة', 14, 3, 1, 2, 2, 500, 200, 500, 100000, false, true, 0);

-- Subscriptions --------------------------------------------------------------------------
CREATE TABLE subscriptions (
  id                    uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id             uuid NOT NULL REFERENCES tenants (id),
  plan_id               uuid NOT NULL REFERENCES plans (id),
  -- Stored status. Expiry is derived from the dates (see subscription_state).
  status                text NOT NULL CHECK (status IN ('TRIALING', 'ACTIVE', 'CANCELLED')),
  billing_cycle         text NOT NULL DEFAULT 'MONTHLY' CHECK (billing_cycle IN ('MONTHLY', 'YEARLY')),
  current_period_start  timestamptz NOT NULL DEFAULT now(),
  current_period_end    timestamptz NOT NULL,
  cancelled_at          timestamptz,
  -- Per-tenant limit changes by the platform administrator, e.g. {"max_users": 25}.
  limit_overrides       jsonb NOT NULL DEFAULT '{}' CHECK (jsonb_typeof(limit_overrides) = 'object'),
  notes                 text,
  created_at            timestamptz NOT NULL DEFAULT now(),
  updated_at            timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT subscriptions_period_ck CHECK (current_period_end > current_period_start),
  CONSTRAINT subscriptions_id_tenant_uq UNIQUE (id, tenant_id)
);
-- One current subscription per tenant; history stays as CANCELLED rows.
CREATE UNIQUE INDEX subscriptions_current_uq ON subscriptions (tenant_id) WHERE status <> 'CANCELLED';
CREATE INDEX subscriptions_plan_idx ON subscriptions (plan_id);
CREATE TRIGGER subscriptions_updated_at BEFORE UPDATE ON subscriptions FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- What the subscription is billed for, with the price at the time (a snapshot).
CREATE TABLE subscription_items (
  id               uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id        uuid NOT NULL,
  subscription_id  uuid NOT NULL,
  item_type        text NOT NULL CHECK (item_type IN ('PLAN', 'ADDON')),
  code             text NOT NULL,
  description      text NOT NULL,
  quantity         integer NOT NULL DEFAULT 1 CHECK (quantity > 0),
  unit_price       numeric(18,2) NOT NULL CHECK (unit_price >= 0),
  currency         char(3) NOT NULL DEFAULT 'SAR',
  created_at       timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT subscription_items_sub_fk FOREIGN KEY (subscription_id, tenant_id) REFERENCES subscriptions (id, tenant_id)
);
CREATE INDEX subscription_items_sub_idx ON subscription_items (subscription_id);

-- Metered usage per calendar month (API calls; storage once uploads exist).
CREATE TABLE usage_records (
  tenant_id   uuid NOT NULL REFERENCES tenants (id),
  metric      text NOT NULL CHECK (metric IN ('API_CALLS', 'STORAGE_MB')),
  period      date NOT NULL CHECK (extract(day FROM period) = 1),
  quantity    bigint NOT NULL DEFAULT 0 CHECK (quantity >= 0),
  updated_at  timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY (tenant_id, metric, period)
);

-- Billing history (append-only): trials, plan changes, payments, suspensions…
CREATE TABLE billing_events (
  id               uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id        uuid NOT NULL REFERENCES tenants (id),
  subscription_id  uuid REFERENCES subscriptions (id),
  event_type       text NOT NULL CHECK (event_type IN (
                     'TRIAL_STARTED', 'PLAN_CHANGED', 'PERIOD_EXTENDED', 'PAYMENT_RECORDED', 'LIMITS_CHANGED',
                     'PLAN_CHANGE_REQUESTED', 'SUBSCRIPTION_CANCELLED', 'TENANT_SUSPENDED', 'TENANT_ACTIVATED')),
  amount           numeric(18,2) CHECK (amount >= 0),
  currency         char(3),
  reference        text,
  details          jsonb NOT NULL DEFAULT '{}',
  created_by       uuid REFERENCES users (id),
  created_at       timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX billing_events_tenant_idx ON billing_events (tenant_id, created_at DESC);
CREATE INDEX billing_events_payments_idx ON billing_events (created_at) WHERE event_type = 'PAYMENT_RECORDED';
CREATE TRIGGER billing_events_append_only BEFORE UPDATE OR DELETE ON billing_events FOR EACH ROW EXECUTE FUNCTION prevent_modification();

-- Feature flags: a global default, optionally overridden per tenant.
CREATE TABLE feature_flags (
  key          text PRIMARY KEY CHECK (key ~ '^[a-z][a-z0-9_]{1,62}$'),
  name_ar      text NOT NULL,
  description  text,
  enabled      boolean NOT NULL DEFAULT false,
  created_at   timestamptz NOT NULL DEFAULT now(),
  updated_at   timestamptz NOT NULL DEFAULT now()
);
CREATE TRIGGER feature_flags_updated_at BEFORE UPDATE ON feature_flags FOR EACH ROW EXECUTE FUNCTION set_updated_at();
INSERT INTO feature_flags (key, name_ar, description, enabled) VALUES
  ('zatca_einvoicing', 'الفوترة الإلكترونية', 'ربط المنشأة بمنصة فاتورة وإدارة وحدات الفوترة', true);

CREATE TABLE tenant_feature_flags (
  tenant_id   uuid NOT NULL REFERENCES tenants (id),
  flag_key    text NOT NULL REFERENCES feature_flags (key) ON DELETE CASCADE,
  enabled     boolean NOT NULL,
  updated_at  timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY (tenant_id, flag_key)
);

-- Unexpected server errors (5xx), for the platform's error view.
CREATE TABLE system_errors (
  id           uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id    uuid,
  user_id      uuid,
  request_id   text,
  method       text,
  path         text,
  status_code  integer,
  error_code   text,
  message      text NOT NULL,
  stack        text,
  created_at   timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX system_errors_created_idx ON system_errors (created_at DESC);

-- What platform administrators did (plans, flags, tenants, users).
CREATE TABLE platform_audit_logs (
  id               uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  admin_user_id    uuid NOT NULL REFERENCES users (id),
  action           text NOT NULL CHECK (action ~ '^[A-Z_]+$'),
  entity_type      text NOT NULL,
  entity_id        text,
  target_tenant_id uuid,
  old_values       jsonb,
  new_values       jsonb,
  ip_address       inet,
  user_agent       text,
  created_at       timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX platform_audit_logs_created_idx ON platform_audit_logs (created_at DESC);
CREATE TRIGGER platform_audit_logs_append_only BEFORE UPDATE OR DELETE ON platform_audit_logs FOR EACH ROW EXECUTE FUNCTION prevent_modification();

-- Every existing tenant starts a trial on the default plan.
INSERT INTO subscriptions (tenant_id, plan_id, status, current_period_start, current_period_end)
SELECT t.id, p.id, 'TRIALING', now(), now() + make_interval(days => GREATEST(p.trial_days, 1))
  FROM tenants t CROSS JOIN plans p
 WHERE p.is_default AND t.deleted_at IS NULL
   AND NOT EXISTS (SELECT 1 FROM subscriptions s WHERE s.tenant_id = t.id AND s.status <> 'CANCELLED');

-- Row-Level Security ---------------------------------------------------------------------
DO $$
DECLARE t text;
BEGIN
  FOREACH t IN ARRAY ARRAY['subscriptions', 'subscription_items', 'usage_records', 'billing_events', 'tenant_feature_flags'] LOOP
    EXECUTE format('ALTER TABLE %I ENABLE ROW LEVEL SECURITY', t);
    EXECUTE format('CREATE POLICY %I ON %I USING (tenant_id = app_current_tenant_id()) WITH CHECK (tenant_id = app_current_tenant_id())',
                   t || '_isolation', t);
  END LOOP;
END
$$;

-- Plans and flags are readable by every tenant; only platform administrators change them.
ALTER TABLE plans ENABLE ROW LEVEL SECURITY;
CREATE POLICY plans_read ON plans FOR SELECT USING (true);
CREATE POLICY plans_write ON plans FOR INSERT WITH CHECK (app_is_platform_admin());
CREATE POLICY plans_update ON plans FOR UPDATE USING (app_is_platform_admin()) WITH CHECK (app_is_platform_admin());
ALTER TABLE feature_flags ENABLE ROW LEVEL SECURITY;
CREATE POLICY feature_flags_read ON feature_flags FOR SELECT USING (true);
CREATE POLICY feature_flags_write ON feature_flags FOR INSERT WITH CHECK (app_is_platform_admin());
CREATE POLICY feature_flags_update ON feature_flags FOR UPDATE USING (app_is_platform_admin()) WITH CHECK (app_is_platform_admin());

-- Errors can be logged from any request; only administrators read them.
ALTER TABLE system_errors ENABLE ROW LEVEL SECURITY;
CREATE POLICY system_errors_insert ON system_errors FOR INSERT WITH CHECK (true);
CREATE POLICY system_errors_read ON system_errors FOR SELECT USING (app_is_platform_admin());
ALTER TABLE platform_audit_logs ENABLE ROW LEVEL SECURITY;
CREATE POLICY platform_audit_insert ON platform_audit_logs FOR INSERT WITH CHECK (app_is_platform_admin());
CREATE POLICY platform_audit_read ON platform_audit_logs FOR SELECT USING (app_is_platform_admin());

-- Platform administrators may READ every tenant table (writes still go
-- through the target tenant's own context). Added to every RLS table that has
-- a tenant_id column, plus tenants and user_tenants.
DO $$
DECLARE t text;
BEGIN
  FOR t IN
    SELECT c.relname FROM pg_class c JOIN pg_namespace n ON n.oid = c.relnamespace
     WHERE n.nspname = 'public' AND c.relkind = 'r' AND c.relrowsecurity
       AND (EXISTS (SELECT 1 FROM pg_attribute a WHERE a.attrelid = c.oid AND a.attname = 'tenant_id' AND NOT a.attisdropped)
            OR c.relname = 'tenants')
  LOOP
    EXECUTE format('CREATE POLICY %I ON %I FOR SELECT USING (app_is_platform_admin())', t || '_platform_read', t);
  END LOOP;
END
$$;

GRANT SELECT, INSERT, UPDATE ON plans, feature_flags, subscriptions TO alshuyukh_app;
GRANT SELECT, INSERT, DELETE ON subscription_items TO alshuyukh_app;
GRANT SELECT, INSERT ON billing_events, platform_audit_logs, system_errors TO alshuyukh_app;
GRANT SELECT, INSERT, UPDATE, DELETE ON tenant_feature_flags TO alshuyukh_app;
GRANT SELECT ON usage_records TO alshuyukh_app;
-- For the platform health view (which migrations ran).
GRANT SELECT ON schema_migrations TO alshuyukh_app;

-- API usage is counted in memory and flushed for many tenants at once.
CREATE OR REPLACE FUNCTION usage_increment(p_tenant uuid, p_metric text, p_period date, p_quantity bigint)
RETURNS void LANGUAGE sql SECURITY DEFINER SET search_path = public AS $$
  INSERT INTO usage_records (tenant_id, metric, period, quantity) VALUES (p_tenant, p_metric, p_period, p_quantity)
  ON CONFLICT (tenant_id, metric, period) DO UPDATE SET quantity = usage_records.quantity + EXCLUDED.quantity, updated_at = now()
$$;
REVOKE ALL ON FUNCTION usage_increment(uuid, text, date, bigint) FROM PUBLIC;
GRANT EXECUTE ON FUNCTION usage_increment(uuid, text, date, bigint) TO alshuyukh_app;
