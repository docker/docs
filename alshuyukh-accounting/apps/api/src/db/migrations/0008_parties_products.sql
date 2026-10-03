-- 0008 — Customers, suppliers, products, units, and document numbering.
--
-- Party balances are not stored: they come from journal lines tagged with
-- customer_id / supplier_id, so the ledger stays the single source of truth.

CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- Document numbering ---------------------------------------------------------------
-- One counter per company and document type (CUSTOMER, SUPPLIER, PRODUCT,
-- and later SALES_INVOICE, ...). Incremented under a row lock inside the
-- caller's transaction, so numbers are gap-free.
CREATE TABLE document_sequences (
  tenant_id    uuid NOT NULL,
  company_id   uuid NOT NULL,
  doc_type     text NOT NULL CHECK (doc_type ~ '^[A-Z][A-Z0-9_]{1,62}$'),
  prefix       text NOT NULL CHECK (prefix ~ '^[A-Z0-9-]{0,12}$'),
  padding      smallint NOT NULL DEFAULT 5 CHECK (padding BETWEEN 1 AND 12),
  last_number  integer NOT NULL DEFAULT 0 CHECK (last_number >= 0),
  updated_at   timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY (company_id, doc_type),
  CONSTRAINT document_sequences_company_fk FOREIGN KEY (company_id, tenant_id) REFERENCES companies (id, tenant_id)
);
CREATE TRIGGER document_sequences_updated_at BEFORE UPDATE ON document_sequences
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- Saudi national address parts (ZATCA buyer/seller address fields).
CREATE OR REPLACE FUNCTION valid_saudi_address(building text, postal text, additional text) RETURNS boolean
LANGUAGE sql IMMUTABLE AS $$
  SELECT (building IS NULL OR building ~ '^[0-9]{4}$')
     AND (postal IS NULL OR postal ~ '^[0-9]{5}$')
     AND (additional IS NULL OR additional ~ '^[0-9]{4}$')
$$;

-- Customers ---------------------------------------------------------------------------
CREATE TABLE customers (
  id                       uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id                uuid NOT NULL,
  company_id               uuid NOT NULL,
  code                     text NOT NULL CHECK (code ~ '^[0-9A-Za-z_-]{1,30}$'),
  party_type               text NOT NULL DEFAULT 'BUSINESS' CHECK (party_type IN ('BUSINESS', 'INDIVIDUAL')),
  name_ar                  text NOT NULL CHECK (length(btrim(name_ar)) BETWEEN 2 AND 200),
  name_en                  text,
  vat_number               text CHECK (vat_number ~ '^3[0-9]{13}3$'),
  commercial_registration  text CHECK (commercial_registration ~ '^[0-9]{10}$'),
  national_id              text CHECK (national_id ~ '^[12][0-9]{9}$'),
  email                    citext,
  phone                    text CHECK (phone ~ '^\+?[0-9 ()-]{6,20}$'),
  credit_limit             numeric(18,2) CHECK (credit_limit >= 0),
  payment_terms_days       smallint NOT NULL DEFAULT 0 CHECK (payment_terms_days BETWEEN 0 AND 365),
  -- Overrides the company's ACCOUNTS_RECEIVABLE account when set.
  receivable_account_id    uuid,
  notes                    text,
  is_active                boolean NOT NULL DEFAULT true,
  created_by               uuid REFERENCES users (id),
  updated_by               uuid REFERENCES users (id),
  created_at               timestamptz NOT NULL DEFAULT now(),
  updated_at               timestamptz NOT NULL DEFAULT now(),
  deleted_at               timestamptz,
  CONSTRAINT customers_company_fk FOREIGN KEY (company_id, tenant_id) REFERENCES companies (id, tenant_id),
  CONSTRAINT customers_account_fk FOREIGN KEY (receivable_account_id, company_id) REFERENCES accounts (id, company_id),
  CONSTRAINT customers_id_company_uq UNIQUE (id, company_id)
);
CREATE UNIQUE INDEX customers_code_uq ON customers (company_id, code) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX customers_vat_uq ON customers (company_id, vat_number) WHERE vat_number IS NOT NULL AND deleted_at IS NULL;
CREATE INDEX customers_company_idx ON customers (tenant_id, company_id) WHERE deleted_at IS NULL;
CREATE INDEX customers_name_trgm ON customers USING gin (name_ar gin_trgm_ops);
CREATE INDEX customers_phone_idx ON customers (company_id, phone) WHERE phone IS NOT NULL;
CREATE TRIGGER customers_updated_at BEFORE UPDATE ON customers
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TABLE customer_addresses (
  id                 uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id          uuid NOT NULL,
  company_id         uuid NOT NULL,
  customer_id        uuid NOT NULL,
  address_type       text NOT NULL DEFAULT 'BILLING' CHECK (address_type IN ('BILLING', 'SHIPPING')),
  is_default         boolean NOT NULL DEFAULT false,
  building_number    text,
  street             text,
  district           text,
  city               text,
  postal_code        text,
  additional_number  text,
  country            char(2) NOT NULL DEFAULT 'SA' CHECK (country ~ '^[A-Z]{2}$'),
  created_at         timestamptz NOT NULL DEFAULT now(),
  updated_at         timestamptz NOT NULL DEFAULT now(),
  deleted_at         timestamptz,
  CONSTRAINT customer_addresses_customer_fk FOREIGN KEY (customer_id, company_id) REFERENCES customers (id, company_id),
  CONSTRAINT customer_addresses_company_fk FOREIGN KEY (company_id, tenant_id) REFERENCES companies (id, tenant_id),
  CONSTRAINT customer_addresses_format_ck CHECK (valid_saudi_address(building_number, postal_code, additional_number))
);
CREATE INDEX customer_addresses_customer_idx ON customer_addresses (customer_id) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX customer_addresses_default_uq ON customer_addresses (customer_id, address_type)
  WHERE is_default AND deleted_at IS NULL;
CREATE TRIGGER customer_addresses_updated_at BEFORE UPDATE ON customer_addresses
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- Suppliers ---------------------------------------------------------------------------
CREATE TABLE suppliers (
  id                       uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id                uuid NOT NULL,
  company_id               uuid NOT NULL,
  code                     text NOT NULL CHECK (code ~ '^[0-9A-Za-z_-]{1,30}$'),
  party_type               text NOT NULL DEFAULT 'BUSINESS' CHECK (party_type IN ('BUSINESS', 'INDIVIDUAL')),
  name_ar                  text NOT NULL CHECK (length(btrim(name_ar)) BETWEEN 2 AND 200),
  name_en                  text,
  vat_number               text CHECK (vat_number ~ '^3[0-9]{13}3$'),
  commercial_registration  text CHECK (commercial_registration ~ '^[0-9]{10}$'),
  national_id              text CHECK (national_id ~ '^[12][0-9]{9}$'),
  email                    citext,
  phone                    text CHECK (phone ~ '^\+?[0-9 ()-]{6,20}$'),
  credit_limit             numeric(18,2) CHECK (credit_limit >= 0),
  payment_terms_days       smallint NOT NULL DEFAULT 0 CHECK (payment_terms_days BETWEEN 0 AND 365),
  -- Overrides the company's ACCOUNTS_PAYABLE account when set.
  payable_account_id       uuid,
  notes                    text,
  is_active                boolean NOT NULL DEFAULT true,
  created_by               uuid REFERENCES users (id),
  updated_by               uuid REFERENCES users (id),
  created_at               timestamptz NOT NULL DEFAULT now(),
  updated_at               timestamptz NOT NULL DEFAULT now(),
  deleted_at               timestamptz,
  CONSTRAINT suppliers_company_fk FOREIGN KEY (company_id, tenant_id) REFERENCES companies (id, tenant_id),
  CONSTRAINT suppliers_account_fk FOREIGN KEY (payable_account_id, company_id) REFERENCES accounts (id, company_id),
  CONSTRAINT suppliers_id_company_uq UNIQUE (id, company_id)
);
CREATE UNIQUE INDEX suppliers_code_uq ON suppliers (company_id, code) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX suppliers_vat_uq ON suppliers (company_id, vat_number) WHERE vat_number IS NOT NULL AND deleted_at IS NULL;
CREATE INDEX suppliers_company_idx ON suppliers (tenant_id, company_id) WHERE deleted_at IS NULL;
CREATE INDEX suppliers_name_trgm ON suppliers USING gin (name_ar gin_trgm_ops);
CREATE TRIGGER suppliers_updated_at BEFORE UPDATE ON suppliers
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TABLE supplier_addresses (
  id                 uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id          uuid NOT NULL,
  company_id         uuid NOT NULL,
  supplier_id        uuid NOT NULL,
  address_type       text NOT NULL DEFAULT 'BILLING' CHECK (address_type IN ('BILLING', 'SHIPPING')),
  is_default         boolean NOT NULL DEFAULT false,
  building_number    text,
  street             text,
  district           text,
  city               text,
  postal_code        text,
  additional_number  text,
  country            char(2) NOT NULL DEFAULT 'SA' CHECK (country ~ '^[A-Z]{2}$'),
  created_at         timestamptz NOT NULL DEFAULT now(),
  updated_at         timestamptz NOT NULL DEFAULT now(),
  deleted_at         timestamptz,
  CONSTRAINT supplier_addresses_supplier_fk FOREIGN KEY (supplier_id, company_id) REFERENCES suppliers (id, company_id),
  CONSTRAINT supplier_addresses_company_fk FOREIGN KEY (company_id, tenant_id) REFERENCES companies (id, tenant_id),
  CONSTRAINT supplier_addresses_format_ck CHECK (valid_saudi_address(building_number, postal_code, additional_number))
);
CREATE INDEX supplier_addresses_supplier_idx ON supplier_addresses (supplier_id) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX supplier_addresses_default_uq ON supplier_addresses (supplier_id, address_type)
  WHERE is_default AND deleted_at IS NULL;
CREATE TRIGGER supplier_addresses_updated_at BEFORE UPDATE ON supplier_addresses
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- Party tagging on journal lines (customer / supplier sub-ledgers) ----------------------
ALTER TABLE journal_entry_lines
  ADD COLUMN customer_id uuid,
  ADD COLUMN supplier_id uuid,
  ADD CONSTRAINT journal_lines_customer_fk FOREIGN KEY (customer_id, company_id) REFERENCES customers (id, company_id),
  ADD CONSTRAINT journal_lines_supplier_fk FOREIGN KEY (supplier_id, company_id) REFERENCES suppliers (id, company_id),
  ADD CONSTRAINT journal_lines_one_party_ck CHECK (customer_id IS NULL OR supplier_id IS NULL);
CREATE INDEX journal_lines_customer_idx ON journal_entry_lines (customer_id) WHERE customer_id IS NOT NULL;
CREATE INDEX journal_lines_supplier_idx ON journal_entry_lines (supplier_id) WHERE supplier_id IS NOT NULL;

-- Products --------------------------------------------------------------------------
CREATE TABLE units (
  id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id   uuid NOT NULL,
  company_id  uuid NOT NULL,
  -- UN/ECE Recommendation 20 code where one exists (PCE, KGM, LTR, ...).
  code        text NOT NULL CHECK (code ~ '^[A-Z0-9]{1,10}$'),
  name_ar     text NOT NULL CHECK (length(btrim(name_ar)) BETWEEN 1 AND 50),
  name_en     text,
  is_active   boolean NOT NULL DEFAULT true,
  created_at  timestamptz NOT NULL DEFAULT now(),
  updated_at  timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT units_company_fk FOREIGN KEY (company_id, tenant_id) REFERENCES companies (id, tenant_id),
  CONSTRAINT units_code_uq UNIQUE (company_id, code),
  CONSTRAINT units_id_company_uq UNIQUE (id, company_id)
);
CREATE TRIGGER units_updated_at BEFORE UPDATE ON units
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TABLE product_categories (
  id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id   uuid NOT NULL,
  company_id  uuid NOT NULL,
  parent_id   uuid,
  name_ar     text NOT NULL CHECK (length(btrim(name_ar)) BETWEEN 2 AND 100),
  name_en     text,
  is_active   boolean NOT NULL DEFAULT true,
  created_at  timestamptz NOT NULL DEFAULT now(),
  updated_at  timestamptz NOT NULL DEFAULT now(),
  deleted_at  timestamptz,
  CONSTRAINT product_categories_company_fk FOREIGN KEY (company_id, tenant_id) REFERENCES companies (id, tenant_id),
  CONSTRAINT product_categories_parent_fk FOREIGN KEY (parent_id, company_id) REFERENCES product_categories (id, company_id),
  CONSTRAINT product_categories_id_company_uq UNIQUE (id, company_id),
  CONSTRAINT product_categories_not_own_parent CHECK (parent_id IS NULL OR parent_id <> id)
);
CREATE UNIQUE INDEX product_categories_name_uq ON product_categories (company_id, COALESCE(parent_id, '00000000-0000-0000-0000-000000000000'::uuid), name_ar)
  WHERE deleted_at IS NULL;
CREATE TRIGGER product_categories_updated_at BEFORE UPDATE ON product_categories
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TABLE products (
  id                       uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id                uuid NOT NULL,
  company_id               uuid NOT NULL,
  sku                      text NOT NULL CHECK (sku ~ '^[0-9A-Za-z._/-]{1,40}$'),
  barcode                  text CHECK (barcode ~ '^[0-9A-Za-z-]{4,40}$'),
  name_ar                  text NOT NULL CHECK (length(btrim(name_ar)) BETWEEN 2 AND 200),
  name_en                  text,
  description              text,
  product_type             text NOT NULL DEFAULT 'GOODS' CHECK (product_type IN ('GOODS', 'SERVICE')),
  category_id              uuid,
  unit_id                  uuid NOT NULL,
  -- Unit prices allow 4 decimals; invoice line totals are rounded to 2.
  sale_price               numeric(18,4) NOT NULL DEFAULT 0 CHECK (sale_price >= 0),
  sale_price_includes_vat  boolean NOT NULL DEFAULT false,
  purchase_price           numeric(18,4) NOT NULL DEFAULT 0 CHECK (purchase_price >= 0),
  -- ZATCA VAT category: S standard, Z zero-rated, E exempt, O out of scope.
  -- The rate itself comes from tax settings (Phase 6), never from the product.
  vat_category             char(1) NOT NULL DEFAULT 'S' CHECK (vat_category IN ('S', 'Z', 'E', 'O')),
  track_inventory          boolean NOT NULL DEFAULT true,
  -- Optional overrides of the company's default SALES / purchase accounts.
  sales_account_id         uuid,
  purchase_account_id      uuid,
  is_active                boolean NOT NULL DEFAULT true,
  created_by               uuid REFERENCES users (id),
  updated_by               uuid REFERENCES users (id),
  created_at               timestamptz NOT NULL DEFAULT now(),
  updated_at               timestamptz NOT NULL DEFAULT now(),
  deleted_at               timestamptz,
  CONSTRAINT products_company_fk FOREIGN KEY (company_id, tenant_id) REFERENCES companies (id, tenant_id),
  CONSTRAINT products_category_fk FOREIGN KEY (category_id, company_id) REFERENCES product_categories (id, company_id),
  CONSTRAINT products_unit_fk FOREIGN KEY (unit_id, company_id) REFERENCES units (id, company_id),
  CONSTRAINT products_sales_account_fk FOREIGN KEY (sales_account_id, company_id) REFERENCES accounts (id, company_id),
  CONSTRAINT products_purchase_account_fk FOREIGN KEY (purchase_account_id, company_id) REFERENCES accounts (id, company_id),
  CONSTRAINT products_id_company_uq UNIQUE (id, company_id),
  -- Services never hold stock.
  CONSTRAINT products_service_stock_ck CHECK (product_type = 'GOODS' OR NOT track_inventory)
);
CREATE UNIQUE INDEX products_sku_uq ON products (company_id, sku) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX products_barcode_uq ON products (company_id, barcode) WHERE barcode IS NOT NULL AND deleted_at IS NULL;
CREATE INDEX products_company_idx ON products (tenant_id, company_id) WHERE deleted_at IS NULL;
CREATE INDEX products_category_idx ON products (category_id);
CREATE INDEX products_name_trgm ON products USING gin (name_ar gin_trgm_ops);
CREATE TRIGGER products_updated_at BEFORE UPDATE ON products
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- Row-Level Security --------------------------------------------------------------------
DO $$
DECLARE t text;
BEGIN
  FOREACH t IN ARRAY ARRAY['document_sequences', 'customers', 'customer_addresses', 'suppliers',
                           'supplier_addresses', 'units', 'product_categories', 'products'] LOOP
    EXECUTE format('ALTER TABLE %I ENABLE ROW LEVEL SECURITY', t);
    EXECUTE format('CREATE POLICY %I ON %I USING (tenant_id = app_current_tenant_id()) WITH CHECK (tenant_id = app_current_tenant_id())',
                   t || '_isolation', t);
    EXECUTE format('GRANT SELECT, INSERT, UPDATE ON %I TO alshuyukh_app', t);
  END LOOP;
END
$$;
