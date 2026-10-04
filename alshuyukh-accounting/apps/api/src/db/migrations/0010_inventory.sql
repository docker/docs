-- 0010 — Inventory: balances, movements, transfers, adjustments.
--
-- Stock is never edited directly. Every change is a row in stock_movements
-- (append-only); inventory_balances holds the running quantity and value per
-- product and warehouse and is updated in the same transaction.
--
-- Costing: weighted average per product and warehouse. Values are kept to
-- the halala (2 decimals) so the stock ledger always equals the Inventory
-- account in the general ledger. An issue takes round(average × qty, 2); the
-- last unit out of a warehouse takes exactly the remaining value.

ALTER TABLE companies
  ADD COLUMN costing_method text NOT NULL DEFAULT 'WEIGHTED_AVERAGE'
  CHECK (costing_method IN ('WEIGHTED_AVERAGE'));  -- FIFO: see inventory/costing.ts

CREATE TABLE inventory_balances (
  tenant_id     uuid NOT NULL,
  company_id    uuid NOT NULL,
  product_id    uuid NOT NULL,
  warehouse_id  uuid NOT NULL,
  quantity      numeric(18,4) NOT NULL DEFAULT 0 CHECK (quantity >= 0),
  value         numeric(18,2) NOT NULL DEFAULT 0 CHECK (value >= 0),
  updated_at    timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY (product_id, warehouse_id),
  CONSTRAINT inventory_balances_company_fk FOREIGN KEY (company_id, tenant_id) REFERENCES companies (id, tenant_id),
  CONSTRAINT inventory_balances_product_fk FOREIGN KEY (product_id, company_id) REFERENCES products (id, company_id),
  CONSTRAINT inventory_balances_warehouse_fk FOREIGN KEY (warehouse_id, company_id) REFERENCES warehouses (id, company_id),
  -- Nothing left means nothing left: no value without quantity.
  CONSTRAINT inventory_balances_empty_ck CHECK (quantity > 0 OR value = 0)
);
CREATE INDEX inventory_balances_warehouse_idx ON inventory_balances (tenant_id, company_id, warehouse_id);
CREATE TRIGGER inventory_balances_updated_at BEFORE UPDATE ON inventory_balances FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TABLE stock_movements (
  id                 uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id          uuid NOT NULL,
  company_id         uuid NOT NULL,
  product_id         uuid NOT NULL,
  warehouse_id       uuid NOT NULL,
  movement_type      text NOT NULL CHECK (movement_type IN (
                       'PURCHASE', 'SALE', 'SALE_RETURN', 'PURCHASE_RETURN', 'TRANSFER_IN', 'TRANSFER_OUT',
                       'ADJUSTMENT_IN', 'ADJUSTMENT_OUT', 'CANCELLATION_IN', 'CANCELLATION_OUT')),
  direction          text NOT NULL CHECK (direction IN ('IN', 'OUT')),
  quantity           numeric(18,4) NOT NULL CHECK (quantity > 0),
  unit_cost          numeric(18,6) NOT NULL CHECK (unit_cost >= 0),
  total_cost         numeric(18,2) NOT NULL CHECK (total_cost >= 0),
  -- Balance right after this movement (stock card).
  balance_quantity   numeric(18,4) NOT NULL,
  balance_value      numeric(18,2) NOT NULL,
  reference_type     text NOT NULL CHECK (reference_type ~ '^[A-Z][A-Z0-9_]{1,62}$'),
  reference_id       uuid NOT NULL,
  reference_line_id  uuid,
  -- For cancellations: the movement being undone.
  reverses_id        uuid REFERENCES stock_movements (id),
  movement_date      date NOT NULL,
  created_by         uuid REFERENCES users (id),
  created_at         timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT stock_movements_company_fk FOREIGN KEY (company_id, tenant_id) REFERENCES companies (id, tenant_id),
  CONSTRAINT stock_movements_product_fk FOREIGN KEY (product_id, company_id) REFERENCES products (id, company_id),
  CONSTRAINT stock_movements_warehouse_fk FOREIGN KEY (warehouse_id, company_id) REFERENCES warehouses (id, company_id),
  CONSTRAINT stock_movements_direction_ck CHECK (
    (direction = 'IN') = (movement_type IN ('PURCHASE', 'SALE_RETURN', 'TRANSFER_IN', 'ADJUSTMENT_IN', 'CANCELLATION_IN'))
  )
);
CREATE INDEX stock_movements_card_idx ON stock_movements (product_id, warehouse_id, created_at);
CREATE INDEX stock_movements_reference_idx ON stock_movements (reference_type, reference_id);
CREATE INDEX stock_movements_line_idx ON stock_movements (reference_line_id) WHERE reference_line_id IS NOT NULL;
CREATE INDEX stock_movements_company_date_idx ON stock_movements (tenant_id, company_id, movement_date);
CREATE UNIQUE INDEX stock_movements_one_reversal_uq ON stock_movements (reverses_id) WHERE reverses_id IS NOT NULL;
CREATE TRIGGER stock_movements_append_only BEFORE UPDATE OR DELETE ON stock_movements
  FOR EACH ROW EXECUTE FUNCTION prevent_modification();

-- Transfers between warehouses of the same company. Posted on creation;
-- corrected with a transfer in the opposite direction.
CREATE TABLE stock_transfers (
  id                 uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id          uuid NOT NULL,
  company_id         uuid NOT NULL,
  transfer_number    text NOT NULL,
  transfer_date      date NOT NULL,
  from_warehouse_id  uuid NOT NULL,
  to_warehouse_id    uuid NOT NULL,
  notes              text,
  total_cost         numeric(18,2) NOT NULL DEFAULT 0,
  created_by         uuid REFERENCES users (id),
  created_at         timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT stock_transfers_company_fk FOREIGN KEY (company_id, tenant_id) REFERENCES companies (id, tenant_id),
  CONSTRAINT stock_transfers_from_fk FOREIGN KEY (from_warehouse_id, company_id) REFERENCES warehouses (id, company_id),
  CONSTRAINT stock_transfers_to_fk FOREIGN KEY (to_warehouse_id, company_id) REFERENCES warehouses (id, company_id),
  CONSTRAINT stock_transfers_distinct_ck CHECK (from_warehouse_id <> to_warehouse_id),
  CONSTRAINT stock_transfers_id_company_uq UNIQUE (id, company_id)
);
CREATE UNIQUE INDEX stock_transfers_number_uq ON stock_transfers (company_id, transfer_number);
CREATE INDEX stock_transfers_date_idx ON stock_transfers (tenant_id, company_id, transfer_date DESC);
CREATE TRIGGER stock_transfers_append_only BEFORE UPDATE OR DELETE ON stock_transfers
  FOR EACH ROW EXECUTE FUNCTION prevent_modification();

CREATE TABLE stock_transfer_items (
  id           uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id    uuid NOT NULL,
  company_id   uuid NOT NULL,
  transfer_id  uuid NOT NULL,
  line_no      smallint NOT NULL,
  product_id   uuid NOT NULL,
  quantity     numeric(18,4) NOT NULL CHECK (quantity > 0),
  total_cost   numeric(18,2) NOT NULL CHECK (total_cost >= 0),
  CONSTRAINT stock_transfer_items_transfer_fk FOREIGN KEY (transfer_id, company_id) REFERENCES stock_transfers (id, company_id),
  CONSTRAINT stock_transfer_items_product_fk FOREIGN KEY (product_id, company_id) REFERENCES products (id, company_id),
  CONSTRAINT stock_transfer_items_line_uq UNIQUE (transfer_id, line_no)
);
CREATE TRIGGER stock_transfer_items_append_only BEFORE UPDATE OR DELETE ON stock_transfer_items
  FOR EACH ROW EXECUTE FUNCTION prevent_modification();

-- Adjustments (stock counts, damage, opening stock). Posted on creation with
-- a journal entry against the offset account.
CREATE TABLE stock_adjustments (
  id                 uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id          uuid NOT NULL,
  company_id         uuid NOT NULL,
  adjustment_number  text NOT NULL,
  adjustment_date    date NOT NULL,
  warehouse_id       uuid NOT NULL,
  reason             text NOT NULL CHECK (length(btrim(reason)) BETWEEN 3 AND 500),
  offset_account_id  uuid NOT NULL,
  total_increase     numeric(18,2) NOT NULL DEFAULT 0,
  total_decrease     numeric(18,2) NOT NULL DEFAULT 0,
  journal_entry_id   uuid,
  created_by         uuid REFERENCES users (id),
  created_at         timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT stock_adjustments_company_fk FOREIGN KEY (company_id, tenant_id) REFERENCES companies (id, tenant_id),
  CONSTRAINT stock_adjustments_warehouse_fk FOREIGN KEY (warehouse_id, company_id) REFERENCES warehouses (id, company_id),
  CONSTRAINT stock_adjustments_account_fk FOREIGN KEY (offset_account_id, company_id) REFERENCES accounts (id, company_id),
  CONSTRAINT stock_adjustments_journal_fk FOREIGN KEY (journal_entry_id, company_id) REFERENCES journal_entries (id, company_id),
  CONSTRAINT stock_adjustments_id_company_uq UNIQUE (id, company_id)
);
CREATE UNIQUE INDEX stock_adjustments_number_uq ON stock_adjustments (company_id, adjustment_number);
CREATE INDEX stock_adjustments_date_idx ON stock_adjustments (tenant_id, company_id, adjustment_date DESC);
CREATE TRIGGER stock_adjustments_append_only BEFORE UPDATE OR DELETE ON stock_adjustments
  FOR EACH ROW EXECUTE FUNCTION prevent_modification();

CREATE TABLE stock_adjustment_items (
  id             uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id      uuid NOT NULL,
  company_id     uuid NOT NULL,
  adjustment_id  uuid NOT NULL,
  line_no        smallint NOT NULL,
  product_id     uuid NOT NULL,
  direction      text NOT NULL CHECK (direction IN ('IN', 'OUT')),
  quantity       numeric(18,4) NOT NULL CHECK (quantity > 0),
  -- For stock counts: what was counted and what the system held before.
  counted_quantity  numeric(18,4),
  system_quantity   numeric(18,4),
  unit_cost      numeric(18,6) NOT NULL CHECK (unit_cost >= 0),
  total_cost     numeric(18,2) NOT NULL CHECK (total_cost >= 0),
  CONSTRAINT stock_adjustment_items_adjustment_fk FOREIGN KEY (adjustment_id, company_id) REFERENCES stock_adjustments (id, company_id),
  CONSTRAINT stock_adjustment_items_product_fk FOREIGN KEY (product_id, company_id) REFERENCES products (id, company_id),
  CONSTRAINT stock_adjustment_items_line_uq UNIQUE (adjustment_id, line_no)
);
CREATE TRIGGER stock_adjustment_items_append_only BEFORE UPDATE OR DELETE ON stock_adjustment_items
  FOR EACH ROW EXECUTE FUNCTION prevent_modification();

-- Inventory differences account for existing companies (new companies get
-- it from the chart template).
INSERT INTO accounts (tenant_id, company_id, code, name_ar, name_en, account_type, parent_id, group_id, is_postable, is_system, system_key)
SELECT c.tenant_id, c.id, '5200', 'فروقات وتسويات المخزون', 'Inventory adjustments', 'COST_OF_GOODS_SOLD',
       h.id, g.id, true, true, 'INVENTORY_ADJUSTMENT'
  FROM companies c
  JOIN accounts h ON h.company_id = c.id AND h.code = '5000' AND h.deleted_at IS NULL AND NOT h.is_postable
  LEFT JOIN account_groups g ON g.company_id = c.id AND g.code = 'COST_OF_SALES'
 WHERE NOT EXISTS (SELECT 1 FROM accounts a WHERE a.company_id = c.id AND a.system_key = 'INVENTORY_ADJUSTMENT')
   AND NOT EXISTS (SELECT 1 FROM accounts a WHERE a.company_id = c.id AND a.code = '5200' AND a.deleted_at IS NULL);

-- Row-Level Security and privileges --------------------------------------------------
DO $$
DECLARE t text;
BEGIN
  FOREACH t IN ARRAY ARRAY['inventory_balances', 'stock_movements', 'stock_transfers', 'stock_transfer_items',
                           'stock_adjustments', 'stock_adjustment_items'] LOOP
    EXECUTE format('ALTER TABLE %I ENABLE ROW LEVEL SECURITY', t);
    EXECUTE format('CREATE POLICY %I ON %I USING (tenant_id = app_current_tenant_id()) WITH CHECK (tenant_id = app_current_tenant_id())',
                   t || '_isolation', t);
    EXECUTE format('GRANT SELECT, INSERT ON %I TO alshuyukh_app', t);
  END LOOP;
END
$$;
-- Balances are the only inventory rows that change after insert.
GRANT UPDATE ON inventory_balances TO alshuyukh_app;
