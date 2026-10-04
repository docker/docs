-- 0011 — Expenses and the tax ledger.
--
-- tax_transactions is the VAT sub-ledger: one row per (document, VAT
-- category, rate) written when a taxable document posts, and a negated row
-- when it is cancelled. The VAT return is built from it and reconciled
-- against the VAT accounts in the general ledger.
--
-- No backfill: documents posted before this migration only exist in
-- development databases.

-- Expense categories -----------------------------------------------------------------
CREATE TABLE expense_categories (
  id            uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id     uuid NOT NULL,
  company_id    uuid NOT NULL,
  code          text NOT NULL CHECK (code ~ '^[A-Z][A-Z0-9_]{1,30}$'),
  name_ar       text NOT NULL CHECK (length(btrim(name_ar)) BETWEEN 2 AND 100),
  account_id    uuid NOT NULL,
  -- Default VAT category for new lines (rent is often exempt, salaries out of scope).
  vat_category  char(1) NOT NULL DEFAULT 'S' CHECK (vat_category IN ('S', 'Z', 'E', 'O')),
  is_active     boolean NOT NULL DEFAULT true,
  created_at    timestamptz NOT NULL DEFAULT now(),
  updated_at    timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT expense_categories_company_fk FOREIGN KEY (company_id, tenant_id) REFERENCES companies (id, tenant_id),
  CONSTRAINT expense_categories_account_fk FOREIGN KEY (account_id, company_id) REFERENCES accounts (id, company_id),
  CONSTRAINT expense_categories_code_uq UNIQUE (company_id, code),
  CONSTRAINT expense_categories_id_company_uq UNIQUE (id, company_id)
);
CREATE TRIGGER expense_categories_updated_at BEFORE UPDATE ON expense_categories FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- Expenses ---------------------------------------------------------------------------
CREATE TABLE expenses (
  id                  uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id           uuid NOT NULL,
  company_id          uuid NOT NULL,
  branch_id           uuid,
  -- Assigned at posting; gap-free per company.
  expense_number      text,
  expense_date        date NOT NULL,
  -- CASH / BANK: paid now through a payment method. CREDIT: owed to a supplier.
  payment_type        text NOT NULL CHECK (payment_type IN ('CASH', 'BANK', 'CREDIT')),
  method_id           uuid,
  supplier_id         uuid,
  payee_name          text,
  -- The vendor's invoice number, needed to claim input VAT.
  reference           text,
  vendor_vat_number   text CHECK (vendor_vat_number ~ '^3[0-9]{13}3$'),
  prices_include_vat  boolean NOT NULL DEFAULT false,
  subtotal            numeric(18,2) NOT NULL DEFAULT 0 CHECK (subtotal >= 0),
  tax_amount          numeric(18,2) NOT NULL DEFAULT 0 CHECK (tax_amount >= 0),
  total               numeric(18,2) NOT NULL DEFAULT 0 CHECK (total >= 0),
  paid_amount         numeric(18,2) NOT NULL DEFAULT 0 CHECK (paid_amount >= 0),
  remaining_amount    numeric(18,2) GENERATED ALWAYS AS (total - paid_amount) STORED CHECK (remaining_amount >= 0),
  notes               text,
  status              text NOT NULL DEFAULT 'DRAFT' CHECK (status IN ('DRAFT', 'POSTED', 'PARTIALLY_PAID', 'PAID', 'CANCELLED')),
  journal_entry_id    uuid,
  created_by          uuid REFERENCES users (id),
  posted_by           uuid REFERENCES users (id),
  posted_at           timestamptz,
  cancelled_by        uuid REFERENCES users (id),
  cancelled_at        timestamptz,
  cancel_reason       text,
  created_at          timestamptz NOT NULL DEFAULT now(),
  updated_at          timestamptz NOT NULL DEFAULT now(),
  deleted_at          timestamptz,
  CONSTRAINT expenses_company_fk FOREIGN KEY (company_id, tenant_id) REFERENCES companies (id, tenant_id),
  CONSTRAINT expenses_branch_fk FOREIGN KEY (branch_id, company_id) REFERENCES branches (id, company_id),
  CONSTRAINT expenses_method_fk FOREIGN KEY (method_id, company_id) REFERENCES payment_methods (id, company_id),
  CONSTRAINT expenses_supplier_fk FOREIGN KEY (supplier_id, company_id) REFERENCES suppliers (id, company_id),
  CONSTRAINT expenses_journal_fk FOREIGN KEY (journal_entry_id, company_id) REFERENCES journal_entries (id, company_id),
  CONSTRAINT expenses_id_company_uq UNIQUE (id, company_id),
  CONSTRAINT expenses_payment_ck CHECK (
    (payment_type = 'CREDIT' AND supplier_id IS NOT NULL) OR (payment_type <> 'CREDIT' AND method_id IS NOT NULL)
  ),
  CONSTRAINT expenses_totals_ck CHECK (total = subtotal + tax_amount),
  CONSTRAINT expenses_posted_ck CHECK (status IN ('DRAFT') OR (expense_number IS NOT NULL AND journal_entry_id IS NOT NULL AND total > 0)),
  CONSTRAINT expenses_soft_delete_ck CHECK (deleted_at IS NULL OR status = 'DRAFT')
);
CREATE UNIQUE INDEX expenses_number_uq ON expenses (company_id, expense_number) WHERE expense_number IS NOT NULL;
CREATE INDEX expenses_company_date_idx ON expenses (tenant_id, company_id, expense_date DESC) WHERE deleted_at IS NULL;
CREATE INDEX expenses_supplier_idx ON expenses (supplier_id) WHERE supplier_id IS NOT NULL;
CREATE TRIGGER expenses_updated_at BEFORE UPDATE ON expenses FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER expenses_guard BEFORE INSERT OR UPDATE OR DELETE ON expenses
  FOR EACH ROW EXECUTE FUNCTION commercial_doc_guard('status', 'paid_amount', 'remaining_amount', 'cancelled_by', 'cancelled_at', 'cancel_reason');

CREATE TABLE expense_items (
  id               uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id        uuid NOT NULL,
  company_id       uuid NOT NULL,
  document_id      uuid NOT NULL,
  line_no          smallint NOT NULL CHECK (line_no BETWEEN 1 AND 500),
  category_id      uuid NOT NULL,
  account_id       uuid NOT NULL,
  description      text,
  cost_center_id   uuid,
  -- Amount as entered (in the expense's price basis).
  amount           numeric(18,2) NOT NULL CHECK (amount > 0),
  net_amount       numeric(18,2) NOT NULL CHECK (net_amount >= 0),
  vat_category     char(1) NOT NULL CHECK (vat_category IN ('S', 'Z', 'E', 'O')),
  vat_rate         numeric(7,4) NOT NULL CHECK (vat_rate >= 0 AND vat_rate < 1),
  vat_amount       numeric(18,2) NOT NULL CHECK (vat_amount >= 0),
  total_amount     numeric(18,2) NOT NULL CHECK (total_amount >= 0),
  CONSTRAINT expense_items_document_fk FOREIGN KEY (document_id, company_id) REFERENCES expenses (id, company_id),
  CONSTRAINT expense_items_company_fk FOREIGN KEY (company_id, tenant_id) REFERENCES companies (id, tenant_id),
  CONSTRAINT expense_items_category_fk FOREIGN KEY (category_id, company_id) REFERENCES expense_categories (id, company_id),
  CONSTRAINT expense_items_account_fk FOREIGN KEY (account_id, company_id) REFERENCES accounts (id, company_id),
  CONSTRAINT expense_items_cost_center_fk FOREIGN KEY (cost_center_id, company_id) REFERENCES cost_centers (id, company_id),
  CONSTRAINT expense_items_line_uq UNIQUE (document_id, line_no),
  CONSTRAINT expense_items_amounts_ck CHECK (total_amount = net_amount + vat_amount)
);
CREATE INDEX expense_items_document_idx ON expense_items (document_id);
CREATE TRIGGER expense_items_guard BEFORE INSERT OR UPDATE OR DELETE ON expense_items
  FOR EACH ROW EXECUTE FUNCTION commercial_item_guard('expenses');

-- Credit expenses are settled with supplier payments.
ALTER TABLE payment_allocations ADD COLUMN expense_id uuid,
  ADD CONSTRAINT payment_allocations_expense_fk FOREIGN KEY (expense_id, company_id) REFERENCES expenses (id, company_id),
  DROP CONSTRAINT payment_allocations_one_target_ck,
  ADD CONSTRAINT payment_allocations_one_target_ck
    CHECK (num_nonnulls(sales_invoice_id, sales_return_id, purchase_invoice_id, purchase_return_id, expense_id) = 1);
CREATE INDEX payment_allocations_expense_idx ON payment_allocations (expense_id) WHERE expense_id IS NOT NULL;

-- Tax ledger -------------------------------------------------------------------------
CREATE TABLE tax_transactions (
  id                uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id         uuid NOT NULL,
  company_id        uuid NOT NULL,
  -- OUTPUT: VAT charged on sales. INPUT: VAT paid on purchases and expenses.
  direction         text NOT NULL CHECK (direction IN ('OUTPUT', 'INPUT')),
  source_type       text NOT NULL CHECK (source_type IN ('SALES_INVOICE', 'SALES_RETURN', 'PURCHASE_INVOICE', 'PURCHASE_RETURN', 'EXPENSE')),
  source_id         uuid NOT NULL,
  source_number     text,
  journal_entry_id  uuid,
  transaction_date  date NOT NULL,
  vat_category      char(1) NOT NULL CHECK (vat_category IN ('S', 'Z', 'E', 'O')),
  vat_rate          numeric(7,4) NOT NULL,
  -- Signed: returns and cancellations are negative.
  taxable_amount    numeric(18,2) NOT NULL,
  tax_amount        numeric(18,2) NOT NULL,
  -- Returns of the period are reported as adjustments.
  is_adjustment     boolean NOT NULL DEFAULT false,
  -- Set on the negated row written when the source is cancelled.
  reverses_id       uuid REFERENCES tax_transactions (id),
  party_name        text,
  party_vat_number  text,
  created_at        timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT tax_transactions_company_fk FOREIGN KEY (company_id, tenant_id) REFERENCES companies (id, tenant_id),
  CONSTRAINT tax_transactions_journal_fk FOREIGN KEY (journal_entry_id, company_id) REFERENCES journal_entries (id, company_id)
);
CREATE INDEX tax_transactions_period_idx ON tax_transactions (tenant_id, company_id, transaction_date);
CREATE INDEX tax_transactions_source_idx ON tax_transactions (source_type, source_id);
CREATE UNIQUE INDEX tax_transactions_one_reversal_uq ON tax_transactions (reverses_id) WHERE reverses_id IS NOT NULL;
CREATE TRIGGER tax_transactions_append_only BEFORE UPDATE OR DELETE ON tax_transactions
  FOR EACH ROW EXECUTE FUNCTION prevent_modification();

-- Default expense categories for existing companies (new companies get them
-- from setupCompanyAccounting).
INSERT INTO expense_categories (tenant_id, company_id, code, name_ar, account_id, vat_category)
SELECT c.tenant_id, c.id, m.code, m.name_ar, a.id, m.vat
  FROM companies c
  CROSS JOIN (VALUES ('RENT', 'إيجار', '6100', 'S'), ('SALARIES', 'رواتب وأجور', '6200', 'O'),
                     ('MARKETING', 'تسويق وإعلان', '6300', 'S'), ('UTILITIES', 'كهرباء ومياه', '6400', 'S'))
       AS m(code, name_ar, account_code, vat)
  JOIN accounts a ON a.company_id = c.id AND a.code = m.account_code AND a.deleted_at IS NULL AND a.is_postable
 WHERE c.deleted_at IS NULL
ON CONFLICT (company_id, code) DO NOTHING;

-- Row-Level Security and privileges ---------------------------------------------------
DO $$
DECLARE t text;
BEGIN
  FOREACH t IN ARRAY ARRAY['expense_categories', 'expenses', 'expense_items', 'tax_transactions'] LOOP
    EXECUTE format('ALTER TABLE %I ENABLE ROW LEVEL SECURITY', t);
    EXECUTE format('CREATE POLICY %I ON %I USING (tenant_id = app_current_tenant_id()) WITH CHECK (tenant_id = app_current_tenant_id())',
                   t || '_isolation', t);
  END LOOP;
END
$$;
GRANT SELECT, INSERT, UPDATE ON expense_categories, expenses TO alshuyukh_app;
GRANT SELECT, INSERT, UPDATE, DELETE ON expense_items TO alshuyukh_app;
GRANT SELECT, INSERT ON tax_transactions TO alshuyukh_app;
