-- 0004 — Companies, branches, warehouses.
-- Every company-related row carries tenant_id. Child tables use composite
-- foreign keys (parent_id, tenant_id) so a row can never point at a parent
-- that belongs to another tenant.
--
-- fiscal_years / fiscal_periods are delivered with the Accounting Engine in
-- Phase 2, because posting rules depend on them.

CREATE TABLE companies (
  id                       uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id                uuid NOT NULL REFERENCES tenants (id) ON DELETE RESTRICT,
  name                     text NOT NULL CHECK (length(btrim(name)) BETWEEN 2 AND 200),
  legal_name               text,
  commercial_registration  text CHECK (commercial_registration ~ '^[0-9]{10}$'),
  -- Saudi VAT registration number: 15 digits, starts and ends with 3.
  vat_number               text CHECK (vat_number ~ '^3[0-9]{13}3$'),
  address                  text,
  city                     text,
  country                  char(2) NOT NULL DEFAULT 'SA' CHECK (country ~ '^[A-Z]{2}$'),
  currency                 char(3) NOT NULL DEFAULT 'SAR' CHECK (currency ~ '^[A-Z]{3}$'),
  timezone                 text NOT NULL DEFAULT 'Asia/Riyadh',
  logo_url                 text,
  email                    citext,
  phone                    text,
  is_active                boolean NOT NULL DEFAULT true,
  created_by               uuid REFERENCES users (id),
  updated_by               uuid REFERENCES users (id),
  created_at               timestamptz NOT NULL DEFAULT now(),
  updated_at               timestamptz NOT NULL DEFAULT now(),
  deleted_at               timestamptz,
  CONSTRAINT companies_id_tenant_uq UNIQUE (id, tenant_id)
);
CREATE INDEX companies_tenant_idx ON companies (tenant_id) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX companies_tenant_name_uq ON companies (tenant_id, name) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX companies_tenant_vat_uq ON companies (tenant_id, vat_number)
  WHERE vat_number IS NOT NULL AND deleted_at IS NULL;
CREATE TRIGGER companies_updated_at BEFORE UPDATE ON companies
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TABLE branches (
  id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id   uuid NOT NULL,
  company_id  uuid NOT NULL,
  code        text NOT NULL CHECK (code ~ '^[A-Za-z0-9_-]{1,20}$'),
  name        text NOT NULL CHECK (length(btrim(name)) BETWEEN 2 AND 200),
  address     text,
  city        text,
  phone       text,
  is_main     boolean NOT NULL DEFAULT false,
  is_active   boolean NOT NULL DEFAULT true,
  created_by  uuid REFERENCES users (id),
  updated_by  uuid REFERENCES users (id),
  created_at  timestamptz NOT NULL DEFAULT now(),
  updated_at  timestamptz NOT NULL DEFAULT now(),
  deleted_at  timestamptz,
  CONSTRAINT branches_company_fk FOREIGN KEY (company_id, tenant_id)
    REFERENCES companies (id, tenant_id) ON DELETE RESTRICT,
  CONSTRAINT branches_id_tenant_uq UNIQUE (id, tenant_id)
);
CREATE INDEX branches_company_idx ON branches (tenant_id, company_id) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX branches_code_uq ON branches (tenant_id, company_id, code) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX branches_one_main_uq ON branches (company_id) WHERE is_main AND deleted_at IS NULL;
CREATE TRIGGER branches_updated_at BEFORE UPDATE ON branches
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TABLE warehouses (
  id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id   uuid NOT NULL,
  company_id  uuid NOT NULL,
  branch_id   uuid,
  code        text NOT NULL CHECK (code ~ '^[A-Za-z0-9_-]{1,20}$'),
  name        text NOT NULL CHECK (length(btrim(name)) BETWEEN 2 AND 200),
  address     text,
  is_active   boolean NOT NULL DEFAULT true,
  created_by  uuid REFERENCES users (id),
  updated_by  uuid REFERENCES users (id),
  created_at  timestamptz NOT NULL DEFAULT now(),
  updated_at  timestamptz NOT NULL DEFAULT now(),
  deleted_at  timestamptz,
  CONSTRAINT warehouses_company_fk FOREIGN KEY (company_id, tenant_id)
    REFERENCES companies (id, tenant_id) ON DELETE RESTRICT,
  CONSTRAINT warehouses_branch_fk FOREIGN KEY (branch_id, tenant_id)
    REFERENCES branches (id, tenant_id) ON DELETE RESTRICT,
  CONSTRAINT warehouses_id_tenant_uq UNIQUE (id, tenant_id)
);
CREATE INDEX warehouses_company_idx ON warehouses (tenant_id, company_id) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX warehouses_code_uq ON warehouses (tenant_id, company_id, code) WHERE deleted_at IS NULL;
CREATE TRIGGER warehouses_updated_at BEFORE UPDATE ON warehouses
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();
