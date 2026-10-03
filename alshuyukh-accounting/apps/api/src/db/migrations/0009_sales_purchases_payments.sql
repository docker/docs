-- 0009 — Sales, purchases, payments, tax rates.
--
-- GENERATED from one template for the six commercial documents (sales quotes,
-- invoices, returns; purchase orders, invoices, returns), which share the
-- same header/line layout. Edit with care and keep the six in sync.
--
-- Money: NUMERIC(18,2). Quantities and unit prices: NUMERIC(18,4).
-- All line amounts are stored VAT-exclusive (see documents/calc.ts).
-- After a document leaves DRAFT, a trigger allows only its status and
-- settlement columns to change; its lines become read-only.

-- Documents reference warehouses within the same company.
ALTER TABLE warehouses ADD CONSTRAINT warehouses_id_company_uq UNIQUE (id, company_id);

-- Tax rates ------------------------------------------------------------------------
-- Only the standard category (S) needs a rate row. Z, E and O are 0 by definition.
CREATE TABLE tax_rates (
  id              uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id       uuid NOT NULL,
  company_id      uuid NOT NULL,
  vat_category    char(1) NOT NULL DEFAULT 'S' CHECK (vat_category = 'S'),
  name_ar         text NOT NULL,
  rate            numeric(7,4) NOT NULL CHECK (rate > 0 AND rate < 1),
  effective_from  date NOT NULL,
  effective_to    date,
  is_active       boolean NOT NULL DEFAULT true,
  created_by      uuid REFERENCES users (id),
  created_at      timestamptz NOT NULL DEFAULT now(),
  updated_at      timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT tax_rates_company_fk FOREIGN KEY (company_id, tenant_id) REFERENCES companies (id, tenant_id),
  CONSTRAINT tax_rates_dates_ck CHECK (effective_to IS NULL OR effective_to >= effective_from),
  CONSTRAINT tax_rates_no_overlap EXCLUDE USING gist (
    company_id WITH =, vat_category WITH =,
    daterange(effective_from, effective_to, '[]') WITH &&
  ) WHERE (is_active)
);
CREATE TRIGGER tax_rates_updated_at BEFORE UPDATE ON tax_rates FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- Payment methods ------------------------------------------------------------------
CREATE TABLE payment_methods (
  id           uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id    uuid NOT NULL,
  company_id   uuid NOT NULL,
  code         text NOT NULL CHECK (code ~ '^[A-Z][A-Z0-9_]{1,30}$'),
  name_ar      text NOT NULL CHECK (length(btrim(name_ar)) BETWEEN 2 AND 100),
  method_type  text NOT NULL CHECK (method_type IN ('CASH', 'BANK', 'CARD', 'STC_PAY', 'TAMARA', 'OTHER')),
  -- Cash, bank, or clearing account the money goes through.
  account_id   uuid NOT NULL,
  is_active    boolean NOT NULL DEFAULT true,
  sort_order   integer NOT NULL DEFAULT 0,
  created_at   timestamptz NOT NULL DEFAULT now(),
  updated_at   timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT payment_methods_company_fk FOREIGN KEY (company_id, tenant_id) REFERENCES companies (id, tenant_id),
  CONSTRAINT payment_methods_account_fk FOREIGN KEY (account_id, company_id) REFERENCES accounts (id, company_id),
  CONSTRAINT payment_methods_code_uq UNIQUE (company_id, code),
  CONSTRAINT payment_methods_id_company_uq UNIQUE (id, company_id)
);
CREATE TRIGGER payment_methods_updated_at BEFORE UPDATE ON payment_methods FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- Document guards --------------------------------------------------------------------
-- TG_ARGV: the columns that may still change after the document leaves DRAFT.
CREATE OR REPLACE FUNCTION commercial_doc_guard() RETURNS trigger
LANGUAGE plpgsql AS $$
DECLARE
  mutable text[] := TG_ARGV || ARRAY['updated_at'];
BEGIN
  IF TG_OP = 'DELETE' THEN
    RAISE EXCEPTION 'documents cannot be deleted' USING ERRCODE = 'insufficient_privilege';
  END IF;
  IF TG_OP = 'INSERT' THEN
    RETURN NEW;
  END IF;
  IF OLD.deleted_at IS NOT NULL THEN
    RAISE EXCEPTION 'a deleted document cannot be changed' USING ERRCODE = 'insufficient_privilege';
  END IF;
  IF OLD.status = 'CANCELLED' THEN
    RAISE EXCEPTION 'a cancelled document cannot be changed' USING ERRCODE = 'insufficient_privilege';
  END IF;
  IF OLD.status <> 'DRAFT' THEN
    IF NEW.status = 'DRAFT' THEN
      RAISE EXCEPTION 'an issued document cannot return to draft' USING ERRCODE = 'insufficient_privilege';
    END IF;
    IF (to_jsonb(NEW) - mutable) IS DISTINCT FROM (to_jsonb(OLD) - mutable) THEN
      RAISE EXCEPTION 'an issued document cannot be modified; cancel it or issue a return' USING ERRCODE = 'insufficient_privilege';
    END IF;
  ELSIF NEW.deleted_at IS NOT NULL AND NEW.status <> 'DRAFT' THEN
    RAISE EXCEPTION 'only drafts can be deleted' USING ERRCODE = 'insufficient_privilege';
  END IF;
  RETURN NEW;
END;
$$;

-- TG_ARGV[0]: the header table. Lines change only while the header is a draft.
CREATE OR REPLACE FUNCTION commercial_item_guard() RETURNS trigger
LANGUAGE plpgsql AS $$
DECLARE
  doc_status text;
  doc_id uuid := CASE WHEN TG_OP = 'DELETE' THEN OLD.document_id ELSE NEW.document_id END;
BEGIN
  EXECUTE format('SELECT status FROM %I WHERE id = $1', TG_ARGV[0]) INTO doc_status USING doc_id;
  IF doc_status IS DISTINCT FROM 'DRAFT' THEN
    RAISE EXCEPTION 'lines of an issued document cannot be changed' USING ERRCODE = 'insufficient_privilege';
  END IF;
  IF TG_OP = 'UPDATE' AND NEW.document_id <> OLD.document_id THEN
    RAISE EXCEPTION 'a line cannot move to another document' USING ERRCODE = 'insufficient_privilege';
  END IF;
  RETURN CASE WHEN TG_OP = 'DELETE' THEN OLD ELSE NEW END;
END;
$$;


-- sales_quotes -----------------------------------------------------------------------
CREATE TABLE sales_quotes (
  id                    uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id             uuid NOT NULL,
  company_id            uuid NOT NULL,
  branch_id             uuid,
  warehouse_id          uuid,
  -- Assigned at creation.
  doc_number            text,
  doc_date              date NOT NULL,
  customer_id           uuid NOT NULL,
  currency              char(3) NOT NULL DEFAULT 'SAR' CHECK (currency ~ '^[A-Z]{3}$'),
  prices_include_vat    boolean NOT NULL DEFAULT false,
  subtotal              numeric(18,2) NOT NULL DEFAULT 0 CHECK (subtotal >= 0),
  discount_total        numeric(18,2) NOT NULL DEFAULT 0 CHECK (discount_total >= 0),
  taxable_amount        numeric(18,2) NOT NULL DEFAULT 0 CHECK (taxable_amount >= 0),
  tax_amount            numeric(18,2) NOT NULL DEFAULT 0 CHECK (tax_amount >= 0),
  total                 numeric(18,2) NOT NULL DEFAULT 0 CHECK (total >= 0),
  valid_until           date,
  converted_invoice_id  uuid,
  notes                 text,
  status                text NOT NULL DEFAULT 'DRAFT' CHECK (status IN ('DRAFT', 'SENT', 'ACCEPTED', 'REJECTED', 'CONVERTED', 'CANCELLED')),
  -- Copy of the customer's name, VAT number and address at issue time.
  party_snapshot        jsonb,
  journal_entry_id      uuid,
  created_by            uuid REFERENCES users (id),
  issued_by             uuid REFERENCES users (id),
  issued_at             timestamptz,
  cancelled_by          uuid REFERENCES users (id),
  cancelled_at          timestamptz,
  cancel_reason         text,
  created_at            timestamptz NOT NULL DEFAULT now(),
  updated_at            timestamptz NOT NULL DEFAULT now(),
  deleted_at            timestamptz,
  CONSTRAINT sales_quotes_company_fk FOREIGN KEY (company_id, tenant_id) REFERENCES companies (id, tenant_id),
  CONSTRAINT sales_quotes_branch_fk FOREIGN KEY (branch_id, company_id) REFERENCES branches (id, company_id),
  CONSTRAINT sales_quotes_warehouse_fk FOREIGN KEY (warehouse_id, company_id) REFERENCES warehouses (id, company_id),
  CONSTRAINT sales_quotes_customer_fk FOREIGN KEY (customer_id, company_id) REFERENCES customers (id, company_id),
  CONSTRAINT sales_quotes_journal_fk FOREIGN KEY (journal_entry_id, company_id) REFERENCES journal_entries (id, company_id),
  CONSTRAINT sales_quotes_id_company_uq UNIQUE (id, company_id),
  CONSTRAINT sales_quotes_totals_ck CHECK (taxable_amount = subtotal - discount_total AND total = taxable_amount + tax_amount),
  CONSTRAINT sales_quotes_soft_delete_ck CHECK (deleted_at IS NULL OR status = 'DRAFT'),
  CONSTRAINT sales_quotes_number_ck CHECK (doc_number IS NOT NULL)
);
CREATE UNIQUE INDEX sales_quotes_number_uq ON sales_quotes (company_id, doc_number) WHERE doc_number IS NOT NULL;
CREATE INDEX sales_quotes_company_date_idx ON sales_quotes (tenant_id, company_id, doc_date DESC) WHERE deleted_at IS NULL;
CREATE INDEX sales_quotes_party_idx ON sales_quotes (customer_id, status);
CREATE INDEX sales_quotes_status_idx ON sales_quotes (company_id, status);
CREATE TRIGGER sales_quotes_updated_at BEFORE UPDATE ON sales_quotes FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER sales_quotes_guard BEFORE INSERT OR UPDATE OR DELETE ON sales_quotes
  FOR EACH ROW EXECUTE FUNCTION commercial_doc_guard('status', 'converted_invoice_id', 'cancelled_by', 'cancelled_at', 'cancel_reason');

CREATE TABLE sales_quote_items (
  id                uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id         uuid NOT NULL,
  company_id        uuid NOT NULL,
  document_id       uuid NOT NULL,
  line_no           smallint NOT NULL CHECK (line_no BETWEEN 1 AND 1000),
  product_id        uuid NOT NULL,
  account_id        uuid,
  description       text,
  quantity          numeric(18,4) NOT NULL CHECK (quantity > 0),
  unit_id           uuid,
  unit_price        numeric(18,4) NOT NULL CHECK (unit_price >= 0),
  gross_amount      numeric(18,2) NOT NULL CHECK (gross_amount >= 0),
  -- Discount as entered (in the document's price basis); used to recalculate.
  discount_basis    numeric(18,2) NOT NULL DEFAULT 0 CHECK (discount_basis >= 0),
  discount_amount   numeric(18,2) NOT NULL DEFAULT 0 CHECK (discount_amount >= 0),
  net_amount        numeric(18,2) NOT NULL CHECK (net_amount >= 0),
  vat_category      char(1) NOT NULL CHECK (vat_category IN ('S', 'Z', 'E', 'O')),
  vat_rate          numeric(7,4) NOT NULL CHECK (vat_rate >= 0 AND vat_rate < 1),
  vat_amount        numeric(18,2) NOT NULL CHECK (vat_amount >= 0),
  total_amount      numeric(18,2) NOT NULL CHECK (total_amount >= 0),
  created_at        timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT sales_quote_items_document_fk FOREIGN KEY (document_id, company_id) REFERENCES sales_quotes (id, company_id),
  CONSTRAINT sales_quote_items_company_fk FOREIGN KEY (company_id, tenant_id) REFERENCES companies (id, tenant_id),
  CONSTRAINT sales_quote_items_product_fk FOREIGN KEY (product_id, company_id) REFERENCES products (id, company_id),
  CONSTRAINT sales_quote_items_account_fk FOREIGN KEY (account_id, company_id) REFERENCES accounts (id, company_id),
  CONSTRAINT sales_quote_items_unit_fk FOREIGN KEY (unit_id, company_id) REFERENCES units (id, company_id),
  CONSTRAINT sales_quote_items_line_uq UNIQUE (document_id, line_no),
  CONSTRAINT sales_quote_items_id_company_uq UNIQUE (id, company_id),
  CONSTRAINT sales_quote_items_amounts_ck CHECK (net_amount = gross_amount - discount_amount AND total_amount = net_amount + vat_amount)
);
CREATE INDEX sales_quote_items_document_idx ON sales_quote_items (document_id);
CREATE INDEX sales_quote_items_product_idx ON sales_quote_items (product_id) WHERE product_id IS NOT NULL;
CREATE TRIGGER sales_quote_items_guard BEFORE INSERT OR UPDATE OR DELETE ON sales_quote_items
  FOR EACH ROW EXECUTE FUNCTION commercial_item_guard('sales_quotes');

-- sales_invoices -----------------------------------------------------------------------
CREATE TABLE sales_invoices (
  id                    uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id             uuid NOT NULL,
  company_id            uuid NOT NULL,
  branch_id             uuid,
  warehouse_id          uuid,
  -- Assigned at issue; gap-free per company.
  doc_number            text,
  doc_date              date NOT NULL,
  customer_id           uuid NOT NULL,
  currency              char(3) NOT NULL DEFAULT 'SAR' CHECK (currency ~ '^[A-Z]{3}$'),
  prices_include_vat    boolean NOT NULL DEFAULT false,
  subtotal              numeric(18,2) NOT NULL DEFAULT 0 CHECK (subtotal >= 0),
  discount_total        numeric(18,2) NOT NULL DEFAULT 0 CHECK (discount_total >= 0),
  taxable_amount        numeric(18,2) NOT NULL DEFAULT 0 CHECK (taxable_amount >= 0),
  tax_amount            numeric(18,2) NOT NULL DEFAULT 0 CHECK (tax_amount >= 0),
  total                 numeric(18,2) NOT NULL DEFAULT 0 CHECK (total >= 0),
  due_date              date,
  -- ZATCA: STANDARD (B2B, buyer VAT number) or SIMPLIFIED (B2C).
  invoice_kind          text NOT NULL DEFAULT 'STANDARD' CHECK (invoice_kind IN ('STANDARD', 'SIMPLIFIED')),
  paid_amount           numeric(18,2) NOT NULL DEFAULT 0 CHECK (paid_amount >= 0),
  returned_amount       numeric(18,2) NOT NULL DEFAULT 0 CHECK (returned_amount >= 0),
  remaining_amount      numeric(18,2) GENERATED ALWAYS AS (total - paid_amount - returned_amount) STORED CHECK (remaining_amount >= 0),
  source_quote_id       uuid,
  notes                 text,
  status                text NOT NULL DEFAULT 'DRAFT' CHECK (status IN ('DRAFT', 'ISSUED', 'PARTIALLY_PAID', 'PAID', 'CANCELLED', 'RETURNED')),
  -- Copy of the customer's name, VAT number and address at issue time.
  party_snapshot        jsonb,
  journal_entry_id      uuid,
  created_by            uuid REFERENCES users (id),
  issued_by             uuid REFERENCES users (id),
  issued_at             timestamptz,
  cancelled_by          uuid REFERENCES users (id),
  cancelled_at          timestamptz,
  cancel_reason         text,
  created_at            timestamptz NOT NULL DEFAULT now(),
  updated_at            timestamptz NOT NULL DEFAULT now(),
  deleted_at            timestamptz,
  CONSTRAINT sales_invoices_company_fk FOREIGN KEY (company_id, tenant_id) REFERENCES companies (id, tenant_id),
  CONSTRAINT sales_invoices_branch_fk FOREIGN KEY (branch_id, company_id) REFERENCES branches (id, company_id),
  CONSTRAINT sales_invoices_warehouse_fk FOREIGN KEY (warehouse_id, company_id) REFERENCES warehouses (id, company_id),
  CONSTRAINT sales_invoices_customer_fk FOREIGN KEY (customer_id, company_id) REFERENCES customers (id, company_id),
  CONSTRAINT sales_invoices_journal_fk FOREIGN KEY (journal_entry_id, company_id) REFERENCES journal_entries (id, company_id),
  CONSTRAINT sales_invoices_id_company_uq UNIQUE (id, company_id),
  CONSTRAINT sales_invoices_totals_ck CHECK (taxable_amount = subtotal - discount_total AND total = taxable_amount + tax_amount),
  CONSTRAINT sales_invoices_soft_delete_ck CHECK (deleted_at IS NULL OR status = 'DRAFT'),
  CONSTRAINT sales_invoices_issued_ck CHECK (status IN ('DRAFT') OR (doc_number IS NOT NULL AND journal_entry_id IS NOT NULL AND issued_at IS NOT NULL AND total > 0))
);
CREATE UNIQUE INDEX sales_invoices_number_uq ON sales_invoices (company_id, doc_number) WHERE doc_number IS NOT NULL;
CREATE INDEX sales_invoices_company_date_idx ON sales_invoices (tenant_id, company_id, doc_date DESC) WHERE deleted_at IS NULL;
CREATE INDEX sales_invoices_party_idx ON sales_invoices (customer_id, status);
CREATE INDEX sales_invoices_status_idx ON sales_invoices (company_id, status);
CREATE TRIGGER sales_invoices_updated_at BEFORE UPDATE ON sales_invoices FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER sales_invoices_guard BEFORE INSERT OR UPDATE OR DELETE ON sales_invoices
  FOR EACH ROW EXECUTE FUNCTION commercial_doc_guard('status', 'paid_amount', 'returned_amount', 'remaining_amount', 'cancelled_by', 'cancelled_at', 'cancel_reason');

CREATE TABLE sales_invoice_items (
  id                uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id         uuid NOT NULL,
  company_id        uuid NOT NULL,
  document_id       uuid NOT NULL,
  line_no           smallint NOT NULL CHECK (line_no BETWEEN 1 AND 1000),
  product_id        uuid NOT NULL,
  account_id        uuid,
  description       text,
  quantity          numeric(18,4) NOT NULL CHECK (quantity > 0),
  unit_id           uuid,
  unit_price        numeric(18,4) NOT NULL CHECK (unit_price >= 0),
  gross_amount      numeric(18,2) NOT NULL CHECK (gross_amount >= 0),
  -- Discount as entered (in the document's price basis); used to recalculate.
  discount_basis    numeric(18,2) NOT NULL DEFAULT 0 CHECK (discount_basis >= 0),
  discount_amount   numeric(18,2) NOT NULL DEFAULT 0 CHECK (discount_amount >= 0),
  net_amount        numeric(18,2) NOT NULL CHECK (net_amount >= 0),
  vat_category      char(1) NOT NULL CHECK (vat_category IN ('S', 'Z', 'E', 'O')),
  vat_rate          numeric(7,4) NOT NULL CHECK (vat_rate >= 0 AND vat_rate < 1),
  vat_amount        numeric(18,2) NOT NULL CHECK (vat_amount >= 0),
  total_amount      numeric(18,2) NOT NULL CHECK (total_amount >= 0),
  created_at        timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT sales_invoice_items_document_fk FOREIGN KEY (document_id, company_id) REFERENCES sales_invoices (id, company_id),
  CONSTRAINT sales_invoice_items_company_fk FOREIGN KEY (company_id, tenant_id) REFERENCES companies (id, tenant_id),
  CONSTRAINT sales_invoice_items_product_fk FOREIGN KEY (product_id, company_id) REFERENCES products (id, company_id),
  CONSTRAINT sales_invoice_items_account_fk FOREIGN KEY (account_id, company_id) REFERENCES accounts (id, company_id),
  CONSTRAINT sales_invoice_items_unit_fk FOREIGN KEY (unit_id, company_id) REFERENCES units (id, company_id),
  CONSTRAINT sales_invoice_items_line_uq UNIQUE (document_id, line_no),
  CONSTRAINT sales_invoice_items_id_company_uq UNIQUE (id, company_id),
  CONSTRAINT sales_invoice_items_amounts_ck CHECK (net_amount = gross_amount - discount_amount AND total_amount = net_amount + vat_amount)
);
CREATE INDEX sales_invoice_items_document_idx ON sales_invoice_items (document_id);
CREATE INDEX sales_invoice_items_product_idx ON sales_invoice_items (product_id) WHERE product_id IS NOT NULL;
CREATE TRIGGER sales_invoice_items_guard BEFORE INSERT OR UPDATE OR DELETE ON sales_invoice_items
  FOR EACH ROW EXECUTE FUNCTION commercial_item_guard('sales_invoices');

-- sales_returns -----------------------------------------------------------------------
CREATE TABLE sales_returns (
  id                    uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id             uuid NOT NULL,
  company_id            uuid NOT NULL,
  branch_id             uuid,
  warehouse_id          uuid,
  -- Assigned at issue; gap-free per company.
  doc_number            text,
  doc_date              date NOT NULL,
  customer_id           uuid NOT NULL,
  currency              char(3) NOT NULL DEFAULT 'SAR' CHECK (currency ~ '^[A-Z]{3}$'),
  prices_include_vat    boolean NOT NULL DEFAULT false,
  subtotal              numeric(18,2) NOT NULL DEFAULT 0 CHECK (subtotal >= 0),
  discount_total        numeric(18,2) NOT NULL DEFAULT 0 CHECK (discount_total >= 0),
  taxable_amount        numeric(18,2) NOT NULL DEFAULT 0 CHECK (taxable_amount >= 0),
  tax_amount            numeric(18,2) NOT NULL DEFAULT 0 CHECK (tax_amount >= 0),
  total                 numeric(18,2) NOT NULL DEFAULT 0 CHECK (total >= 0),
  original_invoice_id   uuid NOT NULL,
  reason                text NOT NULL CHECK (length(btrim(reason)) BETWEEN 3 AND 500),
  -- Part of the credit applied against the original invoice; the rest is owed to the customer.
  applied_amount        numeric(18,2) NOT NULL DEFAULT 0 CHECK (applied_amount >= 0),
  refunded_amount       numeric(18,2) NOT NULL DEFAULT 0 CHECK (refunded_amount >= 0),
  remaining_amount      numeric(18,2) GENERATED ALWAYS AS (total - applied_amount - refunded_amount) STORED CHECK (remaining_amount >= 0),
  notes                 text,
  status                text NOT NULL DEFAULT 'DRAFT' CHECK (status IN ('DRAFT', 'ISSUED', 'CANCELLED')),
  -- Copy of the customer's name, VAT number and address at issue time.
  party_snapshot        jsonb,
  journal_entry_id      uuid,
  created_by            uuid REFERENCES users (id),
  issued_by             uuid REFERENCES users (id),
  issued_at             timestamptz,
  cancelled_by          uuid REFERENCES users (id),
  cancelled_at          timestamptz,
  cancel_reason         text,
  created_at            timestamptz NOT NULL DEFAULT now(),
  updated_at            timestamptz NOT NULL DEFAULT now(),
  deleted_at            timestamptz,
  CONSTRAINT sales_returns_company_fk FOREIGN KEY (company_id, tenant_id) REFERENCES companies (id, tenant_id),
  CONSTRAINT sales_returns_branch_fk FOREIGN KEY (branch_id, company_id) REFERENCES branches (id, company_id),
  CONSTRAINT sales_returns_warehouse_fk FOREIGN KEY (warehouse_id, company_id) REFERENCES warehouses (id, company_id),
  CONSTRAINT sales_returns_customer_fk FOREIGN KEY (customer_id, company_id) REFERENCES customers (id, company_id),
  CONSTRAINT sales_returns_journal_fk FOREIGN KEY (journal_entry_id, company_id) REFERENCES journal_entries (id, company_id),
  CONSTRAINT sales_returns_id_company_uq UNIQUE (id, company_id),
  CONSTRAINT sales_returns_totals_ck CHECK (taxable_amount = subtotal - discount_total AND total = taxable_amount + tax_amount),
  CONSTRAINT sales_returns_soft_delete_ck CHECK (deleted_at IS NULL OR status = 'DRAFT'),
  CONSTRAINT sales_returns_issued_ck CHECK (status IN ('DRAFT') OR (doc_number IS NOT NULL AND journal_entry_id IS NOT NULL AND issued_at IS NOT NULL AND total > 0))
);
ALTER TABLE sales_returns ADD CONSTRAINT sales_returns_original_fk FOREIGN KEY (original_invoice_id, company_id) REFERENCES sales_invoices (id, company_id);
CREATE INDEX sales_returns_original_idx ON sales_returns (original_invoice_id);
CREATE UNIQUE INDEX sales_returns_number_uq ON sales_returns (company_id, doc_number) WHERE doc_number IS NOT NULL;
CREATE INDEX sales_returns_company_date_idx ON sales_returns (tenant_id, company_id, doc_date DESC) WHERE deleted_at IS NULL;
CREATE INDEX sales_returns_party_idx ON sales_returns (customer_id, status);
CREATE INDEX sales_returns_status_idx ON sales_returns (company_id, status);
CREATE TRIGGER sales_returns_updated_at BEFORE UPDATE ON sales_returns FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER sales_returns_guard BEFORE INSERT OR UPDATE OR DELETE ON sales_returns
  FOR EACH ROW EXECUTE FUNCTION commercial_doc_guard('status', 'refunded_amount', 'remaining_amount', 'cancelled_by', 'cancelled_at', 'cancel_reason');

CREATE TABLE sales_return_items (
  id                uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id         uuid NOT NULL,
  company_id        uuid NOT NULL,
  document_id       uuid NOT NULL,
  line_no           smallint NOT NULL CHECK (line_no BETWEEN 1 AND 1000),
  -- The original invoice line being returned.
  source_item_id    uuid NOT NULL,
  product_id        uuid NOT NULL,
  account_id        uuid,
  description       text,
  quantity          numeric(18,4) NOT NULL CHECK (quantity > 0),
  unit_id           uuid,
  unit_price        numeric(18,4) NOT NULL CHECK (unit_price >= 0),
  gross_amount      numeric(18,2) NOT NULL CHECK (gross_amount >= 0),
  -- Discount as entered (in the document's price basis); used to recalculate.
  discount_basis    numeric(18,2) NOT NULL DEFAULT 0 CHECK (discount_basis >= 0),
  discount_amount   numeric(18,2) NOT NULL DEFAULT 0 CHECK (discount_amount >= 0),
  net_amount        numeric(18,2) NOT NULL CHECK (net_amount >= 0),
  vat_category      char(1) NOT NULL CHECK (vat_category IN ('S', 'Z', 'E', 'O')),
  vat_rate          numeric(7,4) NOT NULL CHECK (vat_rate >= 0 AND vat_rate < 1),
  vat_amount        numeric(18,2) NOT NULL CHECK (vat_amount >= 0),
  total_amount      numeric(18,2) NOT NULL CHECK (total_amount >= 0),
  created_at        timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT sales_return_items_document_fk FOREIGN KEY (document_id, company_id) REFERENCES sales_returns (id, company_id),
  CONSTRAINT sales_return_items_company_fk FOREIGN KEY (company_id, tenant_id) REFERENCES companies (id, tenant_id),
  CONSTRAINT sales_return_items_product_fk FOREIGN KEY (product_id, company_id) REFERENCES products (id, company_id),
  CONSTRAINT sales_return_items_account_fk FOREIGN KEY (account_id, company_id) REFERENCES accounts (id, company_id),
  CONSTRAINT sales_return_items_unit_fk FOREIGN KEY (unit_id, company_id) REFERENCES units (id, company_id),
  CONSTRAINT sales_return_items_line_uq UNIQUE (document_id, line_no),
  CONSTRAINT sales_return_items_id_company_uq UNIQUE (id, company_id),
  CONSTRAINT sales_return_items_amounts_ck CHECK (net_amount = gross_amount - discount_amount AND total_amount = net_amount + vat_amount)
);
CREATE INDEX sales_return_items_document_idx ON sales_return_items (document_id);
CREATE INDEX sales_return_items_product_idx ON sales_return_items (product_id) WHERE product_id IS NOT NULL;
CREATE TRIGGER sales_return_items_guard BEFORE INSERT OR UPDATE OR DELETE ON sales_return_items
  FOR EACH ROW EXECUTE FUNCTION commercial_item_guard('sales_returns');

-- purchase_orders -----------------------------------------------------------------------
CREATE TABLE purchase_orders (
  id                    uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id             uuid NOT NULL,
  company_id            uuid NOT NULL,
  branch_id             uuid,
  warehouse_id          uuid,
  -- Assigned at creation.
  doc_number            text,
  doc_date              date NOT NULL,
  supplier_id           uuid NOT NULL,
  currency              char(3) NOT NULL DEFAULT 'SAR' CHECK (currency ~ '^[A-Z]{3}$'),
  prices_include_vat    boolean NOT NULL DEFAULT false,
  subtotal              numeric(18,2) NOT NULL DEFAULT 0 CHECK (subtotal >= 0),
  discount_total        numeric(18,2) NOT NULL DEFAULT 0 CHECK (discount_total >= 0),
  taxable_amount        numeric(18,2) NOT NULL DEFAULT 0 CHECK (taxable_amount >= 0),
  tax_amount            numeric(18,2) NOT NULL DEFAULT 0 CHECK (tax_amount >= 0),
  total                 numeric(18,2) NOT NULL DEFAULT 0 CHECK (total >= 0),
  expected_date         date,
  converted_invoice_id  uuid,
  notes                 text,
  status                text NOT NULL DEFAULT 'DRAFT' CHECK (status IN ('DRAFT', 'APPROVED', 'CONVERTED', 'CANCELLED')),
  -- Copy of the supplier's name, VAT number and address at issue time.
  party_snapshot        jsonb,
  journal_entry_id      uuid,
  created_by            uuid REFERENCES users (id),
  issued_by             uuid REFERENCES users (id),
  issued_at             timestamptz,
  cancelled_by          uuid REFERENCES users (id),
  cancelled_at          timestamptz,
  cancel_reason         text,
  created_at            timestamptz NOT NULL DEFAULT now(),
  updated_at            timestamptz NOT NULL DEFAULT now(),
  deleted_at            timestamptz,
  CONSTRAINT purchase_orders_company_fk FOREIGN KEY (company_id, tenant_id) REFERENCES companies (id, tenant_id),
  CONSTRAINT purchase_orders_branch_fk FOREIGN KEY (branch_id, company_id) REFERENCES branches (id, company_id),
  CONSTRAINT purchase_orders_warehouse_fk FOREIGN KEY (warehouse_id, company_id) REFERENCES warehouses (id, company_id),
  CONSTRAINT purchase_orders_supplier_fk FOREIGN KEY (supplier_id, company_id) REFERENCES suppliers (id, company_id),
  CONSTRAINT purchase_orders_journal_fk FOREIGN KEY (journal_entry_id, company_id) REFERENCES journal_entries (id, company_id),
  CONSTRAINT purchase_orders_id_company_uq UNIQUE (id, company_id),
  CONSTRAINT purchase_orders_totals_ck CHECK (taxable_amount = subtotal - discount_total AND total = taxable_amount + tax_amount),
  CONSTRAINT purchase_orders_soft_delete_ck CHECK (deleted_at IS NULL OR status = 'DRAFT'),
  CONSTRAINT purchase_orders_number_ck CHECK (doc_number IS NOT NULL)
);
CREATE UNIQUE INDEX purchase_orders_number_uq ON purchase_orders (company_id, doc_number) WHERE doc_number IS NOT NULL;
CREATE INDEX purchase_orders_company_date_idx ON purchase_orders (tenant_id, company_id, doc_date DESC) WHERE deleted_at IS NULL;
CREATE INDEX purchase_orders_party_idx ON purchase_orders (supplier_id, status);
CREATE INDEX purchase_orders_status_idx ON purchase_orders (company_id, status);
CREATE TRIGGER purchase_orders_updated_at BEFORE UPDATE ON purchase_orders FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER purchase_orders_guard BEFORE INSERT OR UPDATE OR DELETE ON purchase_orders
  FOR EACH ROW EXECUTE FUNCTION commercial_doc_guard('status', 'converted_invoice_id', 'cancelled_by', 'cancelled_at', 'cancel_reason');

CREATE TABLE purchase_order_items (
  id                uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id         uuid NOT NULL,
  company_id        uuid NOT NULL,
  document_id       uuid NOT NULL,
  line_no           smallint NOT NULL CHECK (line_no BETWEEN 1 AND 1000),
  -- Purchase lines name a product, an account (expense/asset), or both.
  product_id        uuid,
  account_id        uuid,
  description       text,
  quantity          numeric(18,4) NOT NULL CHECK (quantity > 0),
  unit_id           uuid,
  unit_price        numeric(18,4) NOT NULL CHECK (unit_price >= 0),
  gross_amount      numeric(18,2) NOT NULL CHECK (gross_amount >= 0),
  -- Discount as entered (in the document's price basis); used to recalculate.
  discount_basis    numeric(18,2) NOT NULL DEFAULT 0 CHECK (discount_basis >= 0),
  discount_amount   numeric(18,2) NOT NULL DEFAULT 0 CHECK (discount_amount >= 0),
  net_amount        numeric(18,2) NOT NULL CHECK (net_amount >= 0),
  vat_category      char(1) NOT NULL CHECK (vat_category IN ('S', 'Z', 'E', 'O')),
  vat_rate          numeric(7,4) NOT NULL CHECK (vat_rate >= 0 AND vat_rate < 1),
  vat_amount        numeric(18,2) NOT NULL CHECK (vat_amount >= 0),
  total_amount      numeric(18,2) NOT NULL CHECK (total_amount >= 0),
  created_at        timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT purchase_order_items_document_fk FOREIGN KEY (document_id, company_id) REFERENCES purchase_orders (id, company_id),
  CONSTRAINT purchase_order_items_company_fk FOREIGN KEY (company_id, tenant_id) REFERENCES companies (id, tenant_id),
  CONSTRAINT purchase_order_items_product_fk FOREIGN KEY (product_id, company_id) REFERENCES products (id, company_id),
  CONSTRAINT purchase_order_items_account_fk FOREIGN KEY (account_id, company_id) REFERENCES accounts (id, company_id),
  CONSTRAINT purchase_order_items_unit_fk FOREIGN KEY (unit_id, company_id) REFERENCES units (id, company_id),
  CONSTRAINT purchase_order_items_line_uq UNIQUE (document_id, line_no),
  CONSTRAINT purchase_order_items_id_company_uq UNIQUE (id, company_id),
  CONSTRAINT purchase_order_items_amounts_ck CHECK (net_amount = gross_amount - discount_amount AND total_amount = net_amount + vat_amount),
  CONSTRAINT purchase_order_items_target_ck CHECK (product_id IS NOT NULL OR account_id IS NOT NULL)
);
CREATE INDEX purchase_order_items_document_idx ON purchase_order_items (document_id);
CREATE INDEX purchase_order_items_product_idx ON purchase_order_items (product_id) WHERE product_id IS NOT NULL;
CREATE TRIGGER purchase_order_items_guard BEFORE INSERT OR UPDATE OR DELETE ON purchase_order_items
  FOR EACH ROW EXECUTE FUNCTION commercial_item_guard('purchase_orders');

-- purchase_invoices -----------------------------------------------------------------------
CREATE TABLE purchase_invoices (
  id                    uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id             uuid NOT NULL,
  company_id            uuid NOT NULL,
  branch_id             uuid,
  warehouse_id          uuid,
  -- Assigned at issue; gap-free per company.
  doc_number            text,
  doc_date              date NOT NULL,
  supplier_id           uuid NOT NULL,
  currency              char(3) NOT NULL DEFAULT 'SAR' CHECK (currency ~ '^[A-Z]{3}$'),
  prices_include_vat    boolean NOT NULL DEFAULT false,
  subtotal              numeric(18,2) NOT NULL DEFAULT 0 CHECK (subtotal >= 0),
  discount_total        numeric(18,2) NOT NULL DEFAULT 0 CHECK (discount_total >= 0),
  taxable_amount        numeric(18,2) NOT NULL DEFAULT 0 CHECK (taxable_amount >= 0),
  tax_amount            numeric(18,2) NOT NULL DEFAULT 0 CHECK (tax_amount >= 0),
  total                 numeric(18,2) NOT NULL DEFAULT 0 CHECK (total >= 0),
  due_date              date,
  supplier_invoice_number text CHECK (length(btrim(supplier_invoice_number)) BETWEEN 1 AND 60),
  paid_amount           numeric(18,2) NOT NULL DEFAULT 0 CHECK (paid_amount >= 0),
  returned_amount       numeric(18,2) NOT NULL DEFAULT 0 CHECK (returned_amount >= 0),
  remaining_amount      numeric(18,2) GENERATED ALWAYS AS (total - paid_amount - returned_amount) STORED CHECK (remaining_amount >= 0),
  source_order_id       uuid,
  notes                 text,
  status                text NOT NULL DEFAULT 'DRAFT' CHECK (status IN ('DRAFT', 'POSTED', 'PARTIALLY_PAID', 'PAID', 'CANCELLED', 'RETURNED')),
  -- Copy of the supplier's name, VAT number and address at issue time.
  party_snapshot        jsonb,
  journal_entry_id      uuid,
  created_by            uuid REFERENCES users (id),
  issued_by             uuid REFERENCES users (id),
  issued_at             timestamptz,
  cancelled_by          uuid REFERENCES users (id),
  cancelled_at          timestamptz,
  cancel_reason         text,
  created_at            timestamptz NOT NULL DEFAULT now(),
  updated_at            timestamptz NOT NULL DEFAULT now(),
  deleted_at            timestamptz,
  CONSTRAINT purchase_invoices_company_fk FOREIGN KEY (company_id, tenant_id) REFERENCES companies (id, tenant_id),
  CONSTRAINT purchase_invoices_branch_fk FOREIGN KEY (branch_id, company_id) REFERENCES branches (id, company_id),
  CONSTRAINT purchase_invoices_warehouse_fk FOREIGN KEY (warehouse_id, company_id) REFERENCES warehouses (id, company_id),
  CONSTRAINT purchase_invoices_supplier_fk FOREIGN KEY (supplier_id, company_id) REFERENCES suppliers (id, company_id),
  CONSTRAINT purchase_invoices_journal_fk FOREIGN KEY (journal_entry_id, company_id) REFERENCES journal_entries (id, company_id),
  CONSTRAINT purchase_invoices_id_company_uq UNIQUE (id, company_id),
  CONSTRAINT purchase_invoices_totals_ck CHECK (taxable_amount = subtotal - discount_total AND total = taxable_amount + tax_amount),
  CONSTRAINT purchase_invoices_soft_delete_ck CHECK (deleted_at IS NULL OR status = 'DRAFT'),
  CONSTRAINT purchase_invoices_issued_ck CHECK (status IN ('DRAFT') OR (doc_number IS NOT NULL AND journal_entry_id IS NOT NULL AND issued_at IS NOT NULL AND total > 0))
);
CREATE UNIQUE INDEX purchase_invoices_number_uq ON purchase_invoices (company_id, doc_number) WHERE doc_number IS NOT NULL;
CREATE INDEX purchase_invoices_company_date_idx ON purchase_invoices (tenant_id, company_id, doc_date DESC) WHERE deleted_at IS NULL;
CREATE INDEX purchase_invoices_party_idx ON purchase_invoices (supplier_id, status);
CREATE INDEX purchase_invoices_status_idx ON purchase_invoices (company_id, status);
CREATE TRIGGER purchase_invoices_updated_at BEFORE UPDATE ON purchase_invoices FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER purchase_invoices_guard BEFORE INSERT OR UPDATE OR DELETE ON purchase_invoices
  FOR EACH ROW EXECUTE FUNCTION commercial_doc_guard('status', 'paid_amount', 'returned_amount', 'remaining_amount', 'cancelled_by', 'cancelled_at', 'cancel_reason');

CREATE TABLE purchase_invoice_items (
  id                uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id         uuid NOT NULL,
  company_id        uuid NOT NULL,
  document_id       uuid NOT NULL,
  line_no           smallint NOT NULL CHECK (line_no BETWEEN 1 AND 1000),
  -- Purchase lines name a product, an account (expense/asset), or both.
  product_id        uuid,
  account_id        uuid,
  description       text,
  quantity          numeric(18,4) NOT NULL CHECK (quantity > 0),
  unit_id           uuid,
  unit_price        numeric(18,4) NOT NULL CHECK (unit_price >= 0),
  gross_amount      numeric(18,2) NOT NULL CHECK (gross_amount >= 0),
  -- Discount as entered (in the document's price basis); used to recalculate.
  discount_basis    numeric(18,2) NOT NULL DEFAULT 0 CHECK (discount_basis >= 0),
  discount_amount   numeric(18,2) NOT NULL DEFAULT 0 CHECK (discount_amount >= 0),
  net_amount        numeric(18,2) NOT NULL CHECK (net_amount >= 0),
  vat_category      char(1) NOT NULL CHECK (vat_category IN ('S', 'Z', 'E', 'O')),
  vat_rate          numeric(7,4) NOT NULL CHECK (vat_rate >= 0 AND vat_rate < 1),
  vat_amount        numeric(18,2) NOT NULL CHECK (vat_amount >= 0),
  total_amount      numeric(18,2) NOT NULL CHECK (total_amount >= 0),
  created_at        timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT purchase_invoice_items_document_fk FOREIGN KEY (document_id, company_id) REFERENCES purchase_invoices (id, company_id),
  CONSTRAINT purchase_invoice_items_company_fk FOREIGN KEY (company_id, tenant_id) REFERENCES companies (id, tenant_id),
  CONSTRAINT purchase_invoice_items_product_fk FOREIGN KEY (product_id, company_id) REFERENCES products (id, company_id),
  CONSTRAINT purchase_invoice_items_account_fk FOREIGN KEY (account_id, company_id) REFERENCES accounts (id, company_id),
  CONSTRAINT purchase_invoice_items_unit_fk FOREIGN KEY (unit_id, company_id) REFERENCES units (id, company_id),
  CONSTRAINT purchase_invoice_items_line_uq UNIQUE (document_id, line_no),
  CONSTRAINT purchase_invoice_items_id_company_uq UNIQUE (id, company_id),
  CONSTRAINT purchase_invoice_items_amounts_ck CHECK (net_amount = gross_amount - discount_amount AND total_amount = net_amount + vat_amount),
  CONSTRAINT purchase_invoice_items_target_ck CHECK (product_id IS NOT NULL OR account_id IS NOT NULL)
);
CREATE INDEX purchase_invoice_items_document_idx ON purchase_invoice_items (document_id);
CREATE INDEX purchase_invoice_items_product_idx ON purchase_invoice_items (product_id) WHERE product_id IS NOT NULL;
CREATE TRIGGER purchase_invoice_items_guard BEFORE INSERT OR UPDATE OR DELETE ON purchase_invoice_items
  FOR EACH ROW EXECUTE FUNCTION commercial_item_guard('purchase_invoices');

-- purchase_returns -----------------------------------------------------------------------
CREATE TABLE purchase_returns (
  id                    uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id             uuid NOT NULL,
  company_id            uuid NOT NULL,
  branch_id             uuid,
  warehouse_id          uuid,
  -- Assigned at issue; gap-free per company.
  doc_number            text,
  doc_date              date NOT NULL,
  supplier_id           uuid NOT NULL,
  currency              char(3) NOT NULL DEFAULT 'SAR' CHECK (currency ~ '^[A-Z]{3}$'),
  prices_include_vat    boolean NOT NULL DEFAULT false,
  subtotal              numeric(18,2) NOT NULL DEFAULT 0 CHECK (subtotal >= 0),
  discount_total        numeric(18,2) NOT NULL DEFAULT 0 CHECK (discount_total >= 0),
  taxable_amount        numeric(18,2) NOT NULL DEFAULT 0 CHECK (taxable_amount >= 0),
  tax_amount            numeric(18,2) NOT NULL DEFAULT 0 CHECK (tax_amount >= 0),
  total                 numeric(18,2) NOT NULL DEFAULT 0 CHECK (total >= 0),
  original_invoice_id   uuid NOT NULL,
  reason                text NOT NULL CHECK (length(btrim(reason)) BETWEEN 3 AND 500),
  applied_amount        numeric(18,2) NOT NULL DEFAULT 0 CHECK (applied_amount >= 0),
  refunded_amount       numeric(18,2) NOT NULL DEFAULT 0 CHECK (refunded_amount >= 0),
  remaining_amount      numeric(18,2) GENERATED ALWAYS AS (total - applied_amount - refunded_amount) STORED CHECK (remaining_amount >= 0),
  notes                 text,
  status                text NOT NULL DEFAULT 'DRAFT' CHECK (status IN ('DRAFT', 'ISSUED', 'CANCELLED')),
  -- Copy of the supplier's name, VAT number and address at issue time.
  party_snapshot        jsonb,
  journal_entry_id      uuid,
  created_by            uuid REFERENCES users (id),
  issued_by             uuid REFERENCES users (id),
  issued_at             timestamptz,
  cancelled_by          uuid REFERENCES users (id),
  cancelled_at          timestamptz,
  cancel_reason         text,
  created_at            timestamptz NOT NULL DEFAULT now(),
  updated_at            timestamptz NOT NULL DEFAULT now(),
  deleted_at            timestamptz,
  CONSTRAINT purchase_returns_company_fk FOREIGN KEY (company_id, tenant_id) REFERENCES companies (id, tenant_id),
  CONSTRAINT purchase_returns_branch_fk FOREIGN KEY (branch_id, company_id) REFERENCES branches (id, company_id),
  CONSTRAINT purchase_returns_warehouse_fk FOREIGN KEY (warehouse_id, company_id) REFERENCES warehouses (id, company_id),
  CONSTRAINT purchase_returns_supplier_fk FOREIGN KEY (supplier_id, company_id) REFERENCES suppliers (id, company_id),
  CONSTRAINT purchase_returns_journal_fk FOREIGN KEY (journal_entry_id, company_id) REFERENCES journal_entries (id, company_id),
  CONSTRAINT purchase_returns_id_company_uq UNIQUE (id, company_id),
  CONSTRAINT purchase_returns_totals_ck CHECK (taxable_amount = subtotal - discount_total AND total = taxable_amount + tax_amount),
  CONSTRAINT purchase_returns_soft_delete_ck CHECK (deleted_at IS NULL OR status = 'DRAFT'),
  CONSTRAINT purchase_returns_issued_ck CHECK (status IN ('DRAFT') OR (doc_number IS NOT NULL AND journal_entry_id IS NOT NULL AND issued_at IS NOT NULL AND total > 0))
);
ALTER TABLE purchase_returns ADD CONSTRAINT purchase_returns_original_fk FOREIGN KEY (original_invoice_id, company_id) REFERENCES purchase_invoices (id, company_id);
CREATE INDEX purchase_returns_original_idx ON purchase_returns (original_invoice_id);
CREATE UNIQUE INDEX purchase_returns_number_uq ON purchase_returns (company_id, doc_number) WHERE doc_number IS NOT NULL;
CREATE INDEX purchase_returns_company_date_idx ON purchase_returns (tenant_id, company_id, doc_date DESC) WHERE deleted_at IS NULL;
CREATE INDEX purchase_returns_party_idx ON purchase_returns (supplier_id, status);
CREATE INDEX purchase_returns_status_idx ON purchase_returns (company_id, status);
CREATE TRIGGER purchase_returns_updated_at BEFORE UPDATE ON purchase_returns FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER purchase_returns_guard BEFORE INSERT OR UPDATE OR DELETE ON purchase_returns
  FOR EACH ROW EXECUTE FUNCTION commercial_doc_guard('status', 'refunded_amount', 'remaining_amount', 'cancelled_by', 'cancelled_at', 'cancel_reason');

CREATE TABLE purchase_return_items (
  id                uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id         uuid NOT NULL,
  company_id        uuid NOT NULL,
  document_id       uuid NOT NULL,
  line_no           smallint NOT NULL CHECK (line_no BETWEEN 1 AND 1000),
  -- The original invoice line being returned.
  source_item_id    uuid NOT NULL,
  -- Purchase lines name a product, an account (expense/asset), or both.
  product_id        uuid,
  account_id        uuid,
  description       text,
  quantity          numeric(18,4) NOT NULL CHECK (quantity > 0),
  unit_id           uuid,
  unit_price        numeric(18,4) NOT NULL CHECK (unit_price >= 0),
  gross_amount      numeric(18,2) NOT NULL CHECK (gross_amount >= 0),
  -- Discount as entered (in the document's price basis); used to recalculate.
  discount_basis    numeric(18,2) NOT NULL DEFAULT 0 CHECK (discount_basis >= 0),
  discount_amount   numeric(18,2) NOT NULL DEFAULT 0 CHECK (discount_amount >= 0),
  net_amount        numeric(18,2) NOT NULL CHECK (net_amount >= 0),
  vat_category      char(1) NOT NULL CHECK (vat_category IN ('S', 'Z', 'E', 'O')),
  vat_rate          numeric(7,4) NOT NULL CHECK (vat_rate >= 0 AND vat_rate < 1),
  vat_amount        numeric(18,2) NOT NULL CHECK (vat_amount >= 0),
  total_amount      numeric(18,2) NOT NULL CHECK (total_amount >= 0),
  created_at        timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT purchase_return_items_document_fk FOREIGN KEY (document_id, company_id) REFERENCES purchase_returns (id, company_id),
  CONSTRAINT purchase_return_items_company_fk FOREIGN KEY (company_id, tenant_id) REFERENCES companies (id, tenant_id),
  CONSTRAINT purchase_return_items_product_fk FOREIGN KEY (product_id, company_id) REFERENCES products (id, company_id),
  CONSTRAINT purchase_return_items_account_fk FOREIGN KEY (account_id, company_id) REFERENCES accounts (id, company_id),
  CONSTRAINT purchase_return_items_unit_fk FOREIGN KEY (unit_id, company_id) REFERENCES units (id, company_id),
  CONSTRAINT purchase_return_items_line_uq UNIQUE (document_id, line_no),
  CONSTRAINT purchase_return_items_id_company_uq UNIQUE (id, company_id),
  CONSTRAINT purchase_return_items_amounts_ck CHECK (net_amount = gross_amount - discount_amount AND total_amount = net_amount + vat_amount),
  CONSTRAINT purchase_return_items_target_ck CHECK (product_id IS NOT NULL OR account_id IS NOT NULL)
);
CREATE INDEX purchase_return_items_document_idx ON purchase_return_items (document_id);
CREATE INDEX purchase_return_items_product_idx ON purchase_return_items (product_id) WHERE product_id IS NOT NULL;
CREATE TRIGGER purchase_return_items_guard BEFORE INSERT OR UPDATE OR DELETE ON purchase_return_items
  FOR EACH ROW EXECUTE FUNCTION commercial_item_guard('purchase_returns');

ALTER TABLE sales_return_items ADD CONSTRAINT sales_return_items_source_fk
  FOREIGN KEY (source_item_id, company_id) REFERENCES sales_invoice_items (id, company_id);
CREATE INDEX sales_return_items_source_idx ON sales_return_items (source_item_id);
ALTER TABLE purchase_return_items ADD CONSTRAINT purchase_return_items_source_fk
  FOREIGN KEY (source_item_id, company_id) REFERENCES purchase_invoice_items (id, company_id);
CREATE INDEX purchase_return_items_source_idx ON purchase_return_items (source_item_id);

-- Conversion links
ALTER TABLE sales_quotes ADD CONSTRAINT sales_quotes_converted_fk
  FOREIGN KEY (converted_invoice_id, company_id) REFERENCES sales_invoices (id, company_id);
ALTER TABLE sales_invoices ADD CONSTRAINT sales_invoices_quote_fk
  FOREIGN KEY (source_quote_id, company_id) REFERENCES sales_quotes (id, company_id);
ALTER TABLE purchase_orders ADD CONSTRAINT purchase_orders_converted_fk
  FOREIGN KEY (converted_invoice_id, company_id) REFERENCES purchase_invoices (id, company_id);
ALTER TABLE purchase_invoices ADD CONSTRAINT purchase_invoices_order_fk
  FOREIGN KEY (source_order_id, company_id) REFERENCES purchase_orders (id, company_id);

-- A supplier's invoice number is recorded once per supplier.
CREATE UNIQUE INDEX purchase_invoices_supplier_number_uq ON purchase_invoices (supplier_id, supplier_invoice_number)
  WHERE supplier_invoice_number IS NOT NULL AND status <> 'CANCELLED' AND deleted_at IS NULL;

-- Payments ---------------------------------------------------------------------------
CREATE TABLE payments (
  id                 uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id          uuid NOT NULL,
  company_id         uuid NOT NULL,
  branch_id          uuid,
  payment_number     text NOT NULL,
  -- RECEIPT: money in. DISBURSEMENT: money out.
  direction          text NOT NULL CHECK (direction IN ('RECEIPT', 'DISBURSEMENT')),
  customer_id        uuid,
  supplier_id        uuid,
  payment_date       date NOT NULL,
  method_id          uuid NOT NULL,
  amount             numeric(18,2) NOT NULL CHECK (amount > 0),
  allocated_amount   numeric(18,2) NOT NULL DEFAULT 0 CHECK (allocated_amount >= 0),
  unallocated_amount numeric(18,2) GENERATED ALWAYS AS (amount - allocated_amount) STORED CHECK (unallocated_amount >= 0),
  reference          text,
  notes              text,
  status             text NOT NULL DEFAULT 'POSTED' CHECK (status IN ('POSTED', 'VOIDED')),
  journal_entry_id   uuid NOT NULL,
  voided_by          uuid REFERENCES users (id),
  voided_at          timestamptz,
  void_reason        text,
  created_by         uuid REFERENCES users (id),
  created_at         timestamptz NOT NULL DEFAULT now(),
  updated_at         timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT payments_company_fk FOREIGN KEY (company_id, tenant_id) REFERENCES companies (id, tenant_id),
  CONSTRAINT payments_branch_fk FOREIGN KEY (branch_id, company_id) REFERENCES branches (id, company_id),
  CONSTRAINT payments_customer_fk FOREIGN KEY (customer_id, company_id) REFERENCES customers (id, company_id),
  CONSTRAINT payments_supplier_fk FOREIGN KEY (supplier_id, company_id) REFERENCES suppliers (id, company_id),
  CONSTRAINT payments_method_fk FOREIGN KEY (method_id, company_id) REFERENCES payment_methods (id, company_id),
  CONSTRAINT payments_journal_fk FOREIGN KEY (journal_entry_id, company_id) REFERENCES journal_entries (id, company_id),
  CONSTRAINT payments_one_party_ck CHECK ((customer_id IS NULL) <> (supplier_id IS NULL)),
  CONSTRAINT payments_id_company_uq UNIQUE (id, company_id)
);
CREATE UNIQUE INDEX payments_number_uq ON payments (company_id, payment_number);
CREATE INDEX payments_company_date_idx ON payments (tenant_id, company_id, payment_date DESC);
CREATE INDEX payments_customer_idx ON payments (customer_id) WHERE customer_id IS NOT NULL;
CREATE INDEX payments_supplier_idx ON payments (supplier_id) WHERE supplier_id IS NOT NULL;
CREATE TRIGGER payments_updated_at BEFORE UPDATE ON payments FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE OR REPLACE FUNCTION payments_guard() RETURNS trigger
LANGUAGE plpgsql AS $$
DECLARE
  mutable text[] := ARRAY['status', 'allocated_amount', 'unallocated_amount', 'voided_by', 'voided_at', 'void_reason', 'updated_at'];
BEGIN
  IF TG_OP = 'DELETE' THEN
    RAISE EXCEPTION 'payments cannot be deleted; void them instead' USING ERRCODE = 'insufficient_privilege';
  END IF;
  IF OLD.status = 'VOIDED' THEN
    RAISE EXCEPTION 'a voided payment cannot be changed' USING ERRCODE = 'insufficient_privilege';
  END IF;
  IF (to_jsonb(NEW) - mutable) IS DISTINCT FROM (to_jsonb(OLD) - mutable) THEN
    RAISE EXCEPTION 'a payment cannot be modified; void it and record a new one' USING ERRCODE = 'insufficient_privilege';
  END IF;
  RETURN NEW;
END;
$$;
CREATE TRIGGER payments_guard BEFORE UPDATE OR DELETE ON payments FOR EACH ROW EXECUTE FUNCTION payments_guard();

CREATE TABLE payment_allocations (
  id                   uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id            uuid NOT NULL,
  company_id           uuid NOT NULL,
  payment_id           uuid NOT NULL,
  sales_invoice_id     uuid,
  sales_return_id      uuid,
  purchase_invoice_id  uuid,
  purchase_return_id   uuid,
  amount               numeric(18,2) NOT NULL CHECK (amount > 0),
  created_by           uuid REFERENCES users (id),
  created_at           timestamptz NOT NULL DEFAULT now(),
  -- Set when the payment is voided; reversed allocations no longer count.
  reversed_at          timestamptz,
  CONSTRAINT payment_allocations_company_fk FOREIGN KEY (company_id, tenant_id) REFERENCES companies (id, tenant_id),
  CONSTRAINT payment_allocations_payment_fk FOREIGN KEY (payment_id, company_id) REFERENCES payments (id, company_id),
  CONSTRAINT payment_allocations_si_fk FOREIGN KEY (sales_invoice_id, company_id) REFERENCES sales_invoices (id, company_id),
  CONSTRAINT payment_allocations_sr_fk FOREIGN KEY (sales_return_id, company_id) REFERENCES sales_returns (id, company_id),
  CONSTRAINT payment_allocations_pi_fk FOREIGN KEY (purchase_invoice_id, company_id) REFERENCES purchase_invoices (id, company_id),
  CONSTRAINT payment_allocations_pr_fk FOREIGN KEY (purchase_return_id, company_id) REFERENCES purchase_returns (id, company_id),
  CONSTRAINT payment_allocations_one_target_ck CHECK (num_nonnulls(sales_invoice_id, sales_return_id, purchase_invoice_id, purchase_return_id) = 1)
);
CREATE INDEX payment_allocations_payment_idx ON payment_allocations (payment_id);
CREATE INDEX payment_allocations_si_idx ON payment_allocations (sales_invoice_id) WHERE sales_invoice_id IS NOT NULL;
CREATE INDEX payment_allocations_sr_idx ON payment_allocations (sales_return_id) WHERE sales_return_id IS NOT NULL;
CREATE INDEX payment_allocations_pi_idx ON payment_allocations (purchase_invoice_id) WHERE purchase_invoice_id IS NOT NULL;
CREATE INDEX payment_allocations_pr_idx ON payment_allocations (purchase_return_id) WHERE purchase_return_id IS NOT NULL;

-- Seed existing companies: the standard VAT rate and default payment methods.
-- New companies get the same seed from setupCompanyAccounting().
INSERT INTO tax_rates (tenant_id, company_id, name_ar, rate, effective_from)
SELECT tenant_id, id, 'ضريبة القيمة المضافة - النسبة الأساسية', 0.15, DATE '2020-07-01'
  FROM companies WHERE deleted_at IS NULL;

INSERT INTO payment_methods (tenant_id, company_id, code, name_ar, method_type, account_id, sort_order)
SELECT c.tenant_id, c.id, m.code, m.name_ar, m.method_type, a.id, m.sort_order
  FROM companies c
  CROSS JOIN (VALUES ('CASH', 'نقدًا', 'CASH', 'CASH', 1), ('BANK', 'تحويل بنكي', 'BANK', 'BANK', 2),
                     ('CARD', 'بطاقة مدى / ائتمان', 'CARD', 'BANK', 3), ('STC_PAY', 'STC Pay', 'STC_PAY', 'BANK', 4),
                     ('TAMARA', 'تمارا', 'TAMARA', 'BANK', 5)) AS m(code, name_ar, method_type, account_key, sort_order)
  JOIN accounts a ON a.company_id = c.id AND a.system_key = m.account_key AND a.deleted_at IS NULL
 WHERE c.deleted_at IS NULL;

-- Row-Level Security and privileges ---------------------------------------------------
DO $$
DECLARE t text;
BEGIN
  FOREACH t IN ARRAY ARRAY['tax_rates', 'payment_methods', 'payments', 'payment_allocations',
                           'sales_quotes', 'sales_quote_items', 'sales_invoices', 'sales_invoice_items', 'sales_returns', 'sales_return_items', 'purchase_orders', 'purchase_order_items', 'purchase_invoices', 'purchase_invoice_items', 'purchase_returns', 'purchase_return_items'] LOOP
    EXECUTE format('ALTER TABLE %I ENABLE ROW LEVEL SECURITY', t);
    EXECUTE format('CREATE POLICY %I ON %I USING (tenant_id = app_current_tenant_id()) WITH CHECK (tenant_id = app_current_tenant_id())',
                   t || '_isolation', t);
    EXECUTE format('GRANT SELECT, INSERT, UPDATE ON %I TO alshuyukh_app', t);
  END LOOP;
END
$$;
-- Draft lines are replaced with DELETE + INSERT; the item guard refuses that after issue.
GRANT DELETE ON sales_quote_items, sales_invoice_items, sales_return_items, purchase_order_items, purchase_invoice_items, purchase_return_items TO alshuyukh_app;
