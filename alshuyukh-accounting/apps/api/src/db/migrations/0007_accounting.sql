-- 0007 — Accounting core: fiscal calendar, chart of accounts, cost centers,
-- journal entries.
--
-- Money is NUMERIC(18,2) (SAR to the halala). Never float.
-- Every rule that protects the ledger is enforced here in the database, not
-- only in the API: balanced posting, open-period posting, immutability of
-- posted entries, and same-company references.

CREATE EXTENSION IF NOT EXISTS btree_gist;

-- Branches need (id, company_id) as a composite FK target for journal lines.
ALTER TABLE branches ADD CONSTRAINT branches_id_company_uq UNIQUE (id, company_id);

-- Fiscal calendar ------------------------------------------------------------
CREATE TABLE fiscal_years (
  id                uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id         uuid NOT NULL,
  company_id        uuid NOT NULL,
  name              text NOT NULL CHECK (length(btrim(name)) BETWEEN 1 AND 100),
  start_date        date NOT NULL,
  end_date          date NOT NULL,
  status            text NOT NULL DEFAULT 'OPEN' CHECK (status IN ('OPEN', 'CLOSED')),
  closing_entry_id  uuid,
  closed_at         timestamptz,
  closed_by         uuid REFERENCES users (id),
  created_by        uuid REFERENCES users (id),
  created_at        timestamptz NOT NULL DEFAULT now(),
  updated_at        timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT fiscal_years_dates_ck CHECK (end_date > start_date AND end_date - start_date < 550),
  CONSTRAINT fiscal_years_company_fk FOREIGN KEY (company_id, tenant_id) REFERENCES companies (id, tenant_id),
  CONSTRAINT fiscal_years_id_tenant_uq UNIQUE (id, tenant_id),
  CONSTRAINT fiscal_years_id_company_uq UNIQUE (id, company_id),
  CONSTRAINT fiscal_years_no_overlap EXCLUDE USING gist
    (company_id WITH =, daterange(start_date, end_date, '[]') WITH &&)
);
CREATE INDEX fiscal_years_company_idx ON fiscal_years (tenant_id, company_id, start_date);
CREATE TRIGGER fiscal_years_updated_at BEFORE UPDATE ON fiscal_years
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TABLE fiscal_periods (
  id              uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id       uuid NOT NULL,
  company_id      uuid NOT NULL,
  fiscal_year_id  uuid NOT NULL,
  period_number   smallint NOT NULL CHECK (period_number BETWEEN 1 AND 18),
  name            text NOT NULL,
  start_date      date NOT NULL,
  end_date        date NOT NULL,
  status          text NOT NULL DEFAULT 'OPEN' CHECK (status IN ('OPEN', 'CLOSED')),
  closed_at       timestamptz,
  closed_by       uuid REFERENCES users (id),
  created_at      timestamptz NOT NULL DEFAULT now(),
  updated_at      timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT fiscal_periods_dates_ck CHECK (end_date >= start_date),
  CONSTRAINT fiscal_periods_year_fk FOREIGN KEY (fiscal_year_id, company_id) REFERENCES fiscal_years (id, company_id),
  CONSTRAINT fiscal_periods_company_fk FOREIGN KEY (company_id, tenant_id) REFERENCES companies (id, tenant_id),
  CONSTRAINT fiscal_periods_number_uq UNIQUE (fiscal_year_id, period_number),
  CONSTRAINT fiscal_periods_id_tenant_uq UNIQUE (id, tenant_id),
  CONSTRAINT fiscal_periods_id_company_uq UNIQUE (id, company_id),
  CONSTRAINT fiscal_periods_no_overlap EXCLUDE USING gist
    (company_id WITH =, daterange(start_date, end_date, '[]') WITH &&)
);
CREATE INDEX fiscal_periods_lookup_idx ON fiscal_periods (company_id, start_date, end_date);
CREATE TRIGGER fiscal_periods_updated_at BEFORE UPDATE ON fiscal_periods
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- A period of a closed year cannot be reopened.
CREATE OR REPLACE FUNCTION fiscal_periods_guard() RETURNS trigger
LANGUAGE plpgsql AS $$
BEGIN
  IF NEW.status = 'OPEN' AND OLD.status = 'CLOSED' AND EXISTS (
    SELECT 1 FROM fiscal_years y WHERE y.id = NEW.fiscal_year_id AND y.status = 'CLOSED'
  ) THEN
    RAISE EXCEPTION 'cannot reopen a period of a closed fiscal year' USING ERRCODE = 'check_violation';
  END IF;
  IF NEW.start_date <> OLD.start_date OR NEW.end_date <> OLD.end_date OR NEW.fiscal_year_id <> OLD.fiscal_year_id THEN
    RAISE EXCEPTION 'fiscal period dates cannot change' USING ERRCODE = 'check_violation';
  END IF;
  RETURN NEW;
END;
$$;
CREATE TRIGGER fiscal_periods_guard BEFORE UPDATE ON fiscal_periods
  FOR EACH ROW EXECUTE FUNCTION fiscal_periods_guard();

-- Chart of accounts -----------------------------------------------------------
-- Groups classify accounts for financial statements (current assets, etc.).
CREATE TABLE account_groups (
  id            uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id     uuid NOT NULL,
  company_id    uuid NOT NULL,
  code          text NOT NULL CHECK (code ~ '^[A-Z][A-Z0-9_]{1,62}$'),
  name_ar       text NOT NULL,
  name_en       text NOT NULL,
  account_type  text NOT NULL CHECK (account_type IN ('ASSET', 'LIABILITY', 'EQUITY', 'REVENUE', 'EXPENSE', 'COST_OF_GOODS_SOLD')),
  sort_order    integer NOT NULL DEFAULT 0,
  is_system     boolean NOT NULL DEFAULT false,
  created_at    timestamptz NOT NULL DEFAULT now(),
  updated_at    timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT account_groups_company_fk FOREIGN KEY (company_id, tenant_id) REFERENCES companies (id, tenant_id),
  CONSTRAINT account_groups_code_uq UNIQUE (company_id, code),
  CONSTRAINT account_groups_id_company_uq UNIQUE (id, company_id)
);
CREATE TRIGGER account_groups_updated_at BEFORE UPDATE ON account_groups
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TABLE accounts (
  id            uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id     uuid NOT NULL,
  company_id    uuid NOT NULL,
  code          text NOT NULL CHECK (code ~ '^[0-9A-Za-z.-]{1,20}$'),
  name_ar       text NOT NULL CHECK (length(btrim(name_ar)) BETWEEN 1 AND 200),
  name_en       text,
  account_type  text NOT NULL CHECK (account_type IN ('ASSET', 'LIABILITY', 'EQUITY', 'REVENUE', 'EXPENSE', 'COST_OF_GOODS_SOLD')),
  parent_id     uuid,
  level         smallint NOT NULL DEFAULT 1 CHECK (level BETWEEN 1 AND 10),
  group_id      uuid,
  -- Only postable (leaf) accounts accept journal lines; others are headers.
  is_postable   boolean NOT NULL DEFAULT true,
  is_active     boolean NOT NULL DEFAULT true,
  is_system     boolean NOT NULL DEFAULT false,
  -- Stable role used by the engine (AR, AP, VAT_OUTPUT, ...). Codes may be
  -- renumbered by the user; system keys never change.
  system_key    text CHECK (system_key ~ '^[A-Z][A-Z0-9_]{1,62}$'),
  description   text,
  created_by    uuid REFERENCES users (id),
  updated_by    uuid REFERENCES users (id),
  created_at    timestamptz NOT NULL DEFAULT now(),
  updated_at    timestamptz NOT NULL DEFAULT now(),
  deleted_at    timestamptz,
  CONSTRAINT accounts_company_fk FOREIGN KEY (company_id, tenant_id) REFERENCES companies (id, tenant_id),
  CONSTRAINT accounts_parent_fk FOREIGN KEY (parent_id, company_id) REFERENCES accounts (id, company_id),
  CONSTRAINT accounts_group_fk FOREIGN KEY (group_id, company_id) REFERENCES account_groups (id, company_id),
  CONSTRAINT accounts_id_company_uq UNIQUE (id, company_id),
  CONSTRAINT accounts_not_own_parent CHECK (parent_id IS NULL OR parent_id <> id)
);
CREATE UNIQUE INDEX accounts_code_uq ON accounts (company_id, code) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX accounts_system_key_uq ON accounts (company_id, system_key)
  WHERE system_key IS NOT NULL AND deleted_at IS NULL;
CREATE INDEX accounts_parent_idx ON accounts (parent_id);
CREATE INDEX accounts_company_idx ON accounts (tenant_id, company_id) WHERE deleted_at IS NULL;
CREATE TRIGGER accounts_updated_at BEFORE UPDATE ON accounts
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- Keeps the tree consistent: same type as parent, parent is a header,
-- no cycles, level derived from the parent, group type matches.
CREATE OR REPLACE FUNCTION accounts_validate() RETURNS trigger
LANGUAGE plpgsql AS $$
DECLARE
  p record;
  g_type text;
BEGIN
  IF NEW.parent_id IS NULL THEN
    NEW.level := 1;
  ELSE
    SELECT account_type, is_postable, level, deleted_at INTO p FROM accounts WHERE id = NEW.parent_id;
    IF NOT FOUND OR p.deleted_at IS NOT NULL THEN
      RAISE EXCEPTION 'parent account not found' USING ERRCODE = 'foreign_key_violation';
    END IF;
    IF p.account_type <> NEW.account_type THEN
      RAISE EXCEPTION 'account type must match its parent (%)', p.account_type USING ERRCODE = 'check_violation';
    END IF;
    IF p.is_postable THEN
      RAISE EXCEPTION 'parent account must be a header (non-postable) account' USING ERRCODE = 'check_violation';
    END IF;
    IF TG_OP = 'UPDATE' AND NEW.parent_id IS DISTINCT FROM OLD.parent_id AND EXISTS (
      WITH RECURSIVE descendants AS (
        SELECT id FROM accounts WHERE parent_id = NEW.id
        UNION ALL
        SELECT a.id FROM accounts a JOIN descendants d ON a.parent_id = d.id
      ) SELECT 1 FROM descendants WHERE id = NEW.parent_id
    ) THEN
      RAISE EXCEPTION 'an account cannot be moved under its own descendant' USING ERRCODE = 'check_violation';
    END IF;
    NEW.level := p.level + 1;
  END IF;

  IF NEW.group_id IS NOT NULL THEN
    SELECT account_type INTO g_type FROM account_groups WHERE id = NEW.group_id;
    IF g_type IS DISTINCT FROM NEW.account_type THEN
      RAISE EXCEPTION 'account group type must match account type' USING ERRCODE = 'check_violation';
    END IF;
  END IF;

  IF TG_OP = 'UPDATE' THEN
    IF NEW.account_type <> OLD.account_type AND (
      EXISTS (SELECT 1 FROM journal_entry_lines l WHERE l.account_id = NEW.id)
      OR EXISTS (SELECT 1 FROM accounts c WHERE c.parent_id = NEW.id AND c.deleted_at IS NULL)
    ) THEN
      RAISE EXCEPTION 'account type cannot change once the account has children or journal lines' USING ERRCODE = 'check_violation';
    END IF;
    IF NEW.is_postable AND NOT OLD.is_postable
       AND EXISTS (SELECT 1 FROM accounts c WHERE c.parent_id = NEW.id AND c.deleted_at IS NULL) THEN
      RAISE EXCEPTION 'an account with children must stay a header' USING ERRCODE = 'check_violation';
    END IF;
    IF NOT NEW.is_postable AND OLD.is_postable
       AND EXISTS (SELECT 1 FROM journal_entry_lines l WHERE l.account_id = NEW.id) THEN
      RAISE EXCEPTION 'an account with journal lines cannot become a header' USING ERRCODE = 'check_violation';
    END IF;
    IF NEW.deleted_at IS NOT NULL AND OLD.deleted_at IS NULL AND (
      OLD.is_system
      OR EXISTS (SELECT 1 FROM journal_entry_lines l WHERE l.account_id = NEW.id)
      OR EXISTS (SELECT 1 FROM accounts c WHERE c.parent_id = NEW.id AND c.deleted_at IS NULL)
    ) THEN
      RAISE EXCEPTION 'system accounts and accounts with children or journal lines cannot be deleted' USING ERRCODE = 'check_violation';
    END IF;
    IF OLD.is_system AND NEW.system_key IS DISTINCT FROM OLD.system_key THEN
      RAISE EXCEPTION 'system key of a system account cannot change' USING ERRCODE = 'check_violation';
    END IF;
  END IF;
  RETURN NEW;
END;
$$;

-- Cost centers ------------------------------------------------------------------
CREATE TABLE cost_centers (
  id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id   uuid NOT NULL,
  company_id  uuid NOT NULL,
  code        text NOT NULL CHECK (code ~ '^[0-9A-Za-z_-]{1,20}$'),
  name        text NOT NULL CHECK (length(btrim(name)) BETWEEN 2 AND 200),
  is_active   boolean NOT NULL DEFAULT true,
  created_by  uuid REFERENCES users (id),
  updated_by  uuid REFERENCES users (id),
  created_at  timestamptz NOT NULL DEFAULT now(),
  updated_at  timestamptz NOT NULL DEFAULT now(),
  deleted_at  timestamptz,
  CONSTRAINT cost_centers_company_fk FOREIGN KEY (company_id, tenant_id) REFERENCES companies (id, tenant_id),
  CONSTRAINT cost_centers_id_company_uq UNIQUE (id, company_id)
);
CREATE UNIQUE INDEX cost_centers_code_uq ON cost_centers (company_id, code) WHERE deleted_at IS NULL;
CREATE TRIGGER cost_centers_updated_at BEFORE UPDATE ON cost_centers
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- Journal entries -----------------------------------------------------------------
CREATE TABLE journal_entries (
  id                    uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id             uuid NOT NULL,
  company_id            uuid NOT NULL,
  -- Assigned at posting. Gap-free per company and fiscal year.
  fiscal_year_id        uuid,
  fiscal_period_id      uuid,
  entry_number          text,
  entry_date            date NOT NULL,
  description           text NOT NULL CHECK (length(btrim(description)) BETWEEN 1 AND 1000),
  -- MANUAL, REVERSAL, YEAR_CLOSING, and later SALES_INVOICE, PURCHASE_INVOICE, ...
  reference_type        text NOT NULL DEFAULT 'MANUAL' CHECK (reference_type ~ '^[A-Z][A-Z0-9_]{1,62}$'),
  reference_id          uuid,
  source                text NOT NULL DEFAULT 'MANUAL' CHECK (source IN ('MANUAL', 'SYSTEM')),
  status                text NOT NULL DEFAULT 'DRAFT' CHECK (status IN ('DRAFT', 'POSTED', 'REVERSED')),
  currency              char(3) NOT NULL DEFAULT 'SAR' CHECK (currency ~ '^[A-Z]{3}$'),
  total_debit           numeric(18,2) NOT NULL DEFAULT 0,
  total_credit          numeric(18,2) NOT NULL DEFAULT 0,
  reversal_of_id        uuid,
  reversed_by_entry_id  uuid,
  correction_of_id      uuid,
  created_by            uuid REFERENCES users (id),
  posted_by             uuid REFERENCES users (id),
  posted_at             timestamptz,
  reversed_by           uuid REFERENCES users (id),
  reversed_at           timestamptz,
  created_at            timestamptz NOT NULL DEFAULT now(),
  updated_at            timestamptz NOT NULL DEFAULT now(),
  deleted_at            timestamptz,
  CONSTRAINT journal_entries_company_fk FOREIGN KEY (company_id, tenant_id) REFERENCES companies (id, tenant_id),
  CONSTRAINT journal_entries_year_fk FOREIGN KEY (fiscal_year_id, company_id) REFERENCES fiscal_years (id, company_id),
  CONSTRAINT journal_entries_period_fk FOREIGN KEY (fiscal_period_id, company_id) REFERENCES fiscal_periods (id, company_id),
  CONSTRAINT journal_entries_id_company_uq UNIQUE (id, company_id),
  CONSTRAINT journal_entries_reversal_fk FOREIGN KEY (reversal_of_id, company_id) REFERENCES journal_entries (id, company_id),
  CONSTRAINT journal_entries_reversed_by_fk FOREIGN KEY (reversed_by_entry_id, company_id) REFERENCES journal_entries (id, company_id),
  CONSTRAINT journal_entries_correction_fk FOREIGN KEY (correction_of_id, company_id) REFERENCES journal_entries (id, company_id),
  CONSTRAINT journal_entries_posted_ck CHECK (
    status = 'DRAFT' OR (entry_number IS NOT NULL AND posted_at IS NOT NULL
                         AND fiscal_period_id IS NOT NULL AND fiscal_year_id IS NOT NULL
                         AND total_debit = total_credit AND total_debit > 0)
  ),
  CONSTRAINT journal_entries_reversed_ck CHECK (
    status <> 'REVERSED' OR (reversed_by_entry_id IS NOT NULL AND reversed_at IS NOT NULL)
  ),
  CONSTRAINT journal_entries_soft_delete_ck CHECK (deleted_at IS NULL OR status = 'DRAFT')
);
CREATE UNIQUE INDEX journal_entries_number_uq ON journal_entries (company_id, fiscal_year_id, entry_number)
  WHERE entry_number IS NOT NULL;
CREATE UNIQUE INDEX journal_entries_one_reversal_uq ON journal_entries (reversal_of_id) WHERE reversal_of_id IS NOT NULL;
CREATE INDEX journal_entries_date_idx ON journal_entries (tenant_id, company_id, entry_date DESC);
CREATE INDEX journal_entries_status_idx ON journal_entries (company_id, status);
CREATE INDEX journal_entries_reference_idx ON journal_entries (company_id, reference_type, reference_id);
CREATE TRIGGER journal_entries_updated_at BEFORE UPDATE ON journal_entries
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TABLE journal_entry_lines (
  id                uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id         uuid NOT NULL,
  company_id        uuid NOT NULL,
  journal_entry_id  uuid NOT NULL,
  line_no           smallint NOT NULL CHECK (line_no BETWEEN 1 AND 1000),
  account_id        uuid NOT NULL,
  debit             numeric(18,2) NOT NULL DEFAULT 0,
  credit            numeric(18,2) NOT NULL DEFAULT 0,
  description       text,
  cost_center_id    uuid,
  branch_id         uuid,
  created_at        timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT journal_lines_amounts_ck CHECK (
    debit >= 0 AND credit >= 0 AND ((debit > 0 AND credit = 0) OR (credit > 0 AND debit = 0))
  ),
  CONSTRAINT journal_lines_company_fk FOREIGN KEY (company_id, tenant_id) REFERENCES companies (id, tenant_id),
  CONSTRAINT journal_lines_entry_fk FOREIGN KEY (journal_entry_id, company_id) REFERENCES journal_entries (id, company_id),
  CONSTRAINT journal_lines_account_fk FOREIGN KEY (account_id, company_id) REFERENCES accounts (id, company_id),
  CONSTRAINT journal_lines_cost_center_fk FOREIGN KEY (cost_center_id, company_id) REFERENCES cost_centers (id, company_id),
  CONSTRAINT journal_lines_branch_fk FOREIGN KEY (branch_id, company_id) REFERENCES branches (id, company_id),
  CONSTRAINT journal_lines_line_no_uq UNIQUE (journal_entry_id, line_no)
);
CREATE INDEX journal_lines_account_idx ON journal_entry_lines (account_id);
CREATE INDEX journal_lines_tenant_account_idx ON journal_entry_lines (tenant_id, company_id, account_id);
CREATE INDEX journal_lines_cost_center_idx ON journal_entry_lines (cost_center_id) WHERE cost_center_id IS NOT NULL;
CREATE INDEX journal_lines_branch_idx ON journal_entry_lines (branch_id) WHERE branch_id IS NOT NULL;

-- Now that journal_entry_lines exists, attach the account validation trigger.
CREATE TRIGGER accounts_validate BEFORE INSERT OR UPDATE ON accounts
  FOR EACH ROW EXECUTE FUNCTION accounts_validate();

ALTER TABLE fiscal_years ADD CONSTRAINT fiscal_years_closing_entry_fk
  FOREIGN KEY (closing_entry_id, company_id) REFERENCES journal_entries (id, company_id);

-- Gap-free entry numbering per company and fiscal year.
CREATE TABLE journal_sequences (
  fiscal_year_id  uuid PRIMARY KEY,
  tenant_id       uuid NOT NULL,
  company_id      uuid NOT NULL,
  last_number     integer NOT NULL DEFAULT 0 CHECK (last_number >= 0),
  CONSTRAINT journal_sequences_year_fk FOREIGN KEY (fiscal_year_id, company_id) REFERENCES fiscal_years (id, company_id),
  CONSTRAINT journal_sequences_company_fk FOREIGN KEY (company_id, tenant_id) REFERENCES companies (id, tenant_id)
);

-- Ledger protection triggers ------------------------------------------------------

-- Lines can change only while their entry is a draft, and may reference only
-- active postable accounts.
CREATE OR REPLACE FUNCTION journal_lines_guard() RETURNS trigger
LANGUAGE plpgsql AS $$
DECLARE
  entry_status text;
  acc record;
BEGIN
  IF TG_OP IN ('UPDATE', 'DELETE') THEN
    SELECT status INTO entry_status FROM journal_entries WHERE id = OLD.journal_entry_id;
    IF entry_status IS DISTINCT FROM 'DRAFT' THEN
      RAISE EXCEPTION 'lines of a posted journal entry cannot be changed' USING ERRCODE = 'insufficient_privilege';
    END IF;
    IF TG_OP = 'DELETE' THEN RETURN OLD; END IF;
  END IF;

  SELECT status INTO entry_status FROM journal_entries WHERE id = NEW.journal_entry_id;
  IF entry_status IS DISTINCT FROM 'DRAFT' THEN
    RAISE EXCEPTION 'lines of a posted journal entry cannot be changed' USING ERRCODE = 'insufficient_privilege';
  END IF;
  SELECT is_postable, is_active, deleted_at INTO acc FROM accounts WHERE id = NEW.account_id;
  IF NOT FOUND OR acc.deleted_at IS NOT NULL OR NOT acc.is_active OR NOT acc.is_postable THEN
    RAISE EXCEPTION 'account % is not an active postable account', NEW.account_id USING ERRCODE = 'check_violation';
  END IF;
  RETURN NEW;
END;
$$;
CREATE TRIGGER journal_lines_guard BEFORE INSERT OR UPDATE OR DELETE ON journal_entry_lines
  FOR EACH ROW EXECUTE FUNCTION journal_lines_guard();

-- Entries are created as drafts. Posting re-validates everything. Posted
-- entries are immutable except for being marked REVERSED.
CREATE OR REPLACE FUNCTION journal_entries_guard() RETURNS trigger
LANGUAGE plpgsql AS $$
DECLARE
  v_debit   numeric(18,2);
  v_credit  numeric(18,2);
  v_lines   integer;
  v_bad     integer;
  v_period  record;
  mutable   text[] := ARRAY['status', 'reversed_at', 'reversed_by', 'reversed_by_entry_id', 'updated_at'];
BEGIN
  IF TG_OP = 'INSERT' THEN
    IF NEW.status <> 'DRAFT' THEN
      RAISE EXCEPTION 'journal entries must be created as DRAFT and then posted' USING ERRCODE = 'check_violation';
    END IF;
    RETURN NEW;
  END IF;

  IF TG_OP = 'DELETE' THEN
    RAISE EXCEPTION 'journal entries cannot be deleted' USING ERRCODE = 'insufficient_privilege';
  END IF;

  -- UPDATE
  IF OLD.status = 'REVERSED' THEN
    RAISE EXCEPTION 'a reversed journal entry cannot be changed' USING ERRCODE = 'insufficient_privilege';
  END IF;

  IF OLD.status = 'POSTED' THEN
    IF NEW.status <> 'REVERSED'
       OR (to_jsonb(NEW) - mutable) IS DISTINCT FROM (to_jsonb(OLD) - mutable) THEN
      RAISE EXCEPTION 'a posted journal entry cannot be modified; reverse it instead' USING ERRCODE = 'insufficient_privilege';
    END IF;
    RETURN NEW;
  END IF;

  -- OLD.status = 'DRAFT'
  IF OLD.deleted_at IS NOT NULL THEN
    RAISE EXCEPTION 'a deleted draft cannot be changed' USING ERRCODE = 'insufficient_privilege';
  END IF;
  IF NEW.status = 'REVERSED' THEN
    RAISE EXCEPTION 'a draft cannot be reversed' USING ERRCODE = 'check_violation';
  END IF;

  IF NEW.status = 'POSTED' THEN
    SELECT count(*), COALESCE(sum(debit), 0), COALESCE(sum(credit), 0)
      INTO v_lines, v_debit, v_credit
      FROM journal_entry_lines WHERE journal_entry_id = NEW.id;
    IF v_lines < 2 THEN
      RAISE EXCEPTION 'a journal entry needs at least two lines' USING ERRCODE = 'check_violation';
    END IF;
    IF v_debit <> v_credit OR v_debit = 0 THEN
      RAISE EXCEPTION 'journal entry is not balanced (debit %, credit %)', v_debit, v_credit USING ERRCODE = 'check_violation';
    END IF;

    SELECT count(*) INTO v_bad
      FROM journal_entry_lines l JOIN accounts a ON a.id = l.account_id
     WHERE l.journal_entry_id = NEW.id
       AND (a.deleted_at IS NOT NULL OR NOT a.is_active OR NOT a.is_postable);
    IF v_bad > 0 THEN
      RAISE EXCEPTION 'journal entry uses inactive or non-postable accounts' USING ERRCODE = 'check_violation';
    END IF;

    SELECT p.id, p.status AS period_status, y.id AS year_id, y.status AS year_status INTO v_period
      FROM fiscal_periods p JOIN fiscal_years y ON y.id = p.fiscal_year_id
     WHERE p.company_id = NEW.company_id AND NEW.entry_date BETWEEN p.start_date AND p.end_date;
    IF NOT FOUND THEN
      RAISE EXCEPTION 'no fiscal period covers %', NEW.entry_date USING ERRCODE = 'check_violation';
    END IF;
    IF v_period.period_status <> 'OPEN' OR v_period.year_status <> 'OPEN' THEN
      RAISE EXCEPTION 'fiscal period for % is closed', NEW.entry_date USING ERRCODE = 'check_violation';
    END IF;
    IF NEW.fiscal_period_id IS DISTINCT FROM v_period.id OR NEW.fiscal_year_id IS DISTINCT FROM v_period.year_id THEN
      RAISE EXCEPTION 'fiscal period does not match the entry date' USING ERRCODE = 'check_violation';
    END IF;

    NEW.total_debit := v_debit;
    NEW.total_credit := v_credit;
  END IF;
  RETURN NEW;
END;
$$;
CREATE TRIGGER journal_entries_guard BEFORE INSERT OR UPDATE OR DELETE ON journal_entries
  FOR EACH ROW EXECUTE FUNCTION journal_entries_guard();

-- Row-Level Security --------------------------------------------------------------
ALTER TABLE fiscal_years ENABLE ROW LEVEL SECURITY;
CREATE POLICY fiscal_years_isolation ON fiscal_years
  USING (tenant_id = app_current_tenant_id()) WITH CHECK (tenant_id = app_current_tenant_id());
ALTER TABLE fiscal_periods ENABLE ROW LEVEL SECURITY;
CREATE POLICY fiscal_periods_isolation ON fiscal_periods
  USING (tenant_id = app_current_tenant_id()) WITH CHECK (tenant_id = app_current_tenant_id());
ALTER TABLE account_groups ENABLE ROW LEVEL SECURITY;
CREATE POLICY account_groups_isolation ON account_groups
  USING (tenant_id = app_current_tenant_id()) WITH CHECK (tenant_id = app_current_tenant_id());
ALTER TABLE accounts ENABLE ROW LEVEL SECURITY;
CREATE POLICY accounts_isolation ON accounts
  USING (tenant_id = app_current_tenant_id()) WITH CHECK (tenant_id = app_current_tenant_id());
ALTER TABLE cost_centers ENABLE ROW LEVEL SECURITY;
CREATE POLICY cost_centers_isolation ON cost_centers
  USING (tenant_id = app_current_tenant_id()) WITH CHECK (tenant_id = app_current_tenant_id());
ALTER TABLE journal_entries ENABLE ROW LEVEL SECURITY;
CREATE POLICY journal_entries_isolation ON journal_entries
  USING (tenant_id = app_current_tenant_id()) WITH CHECK (tenant_id = app_current_tenant_id());
ALTER TABLE journal_entry_lines ENABLE ROW LEVEL SECURITY;
CREATE POLICY journal_entry_lines_isolation ON journal_entry_lines
  USING (tenant_id = app_current_tenant_id()) WITH CHECK (tenant_id = app_current_tenant_id());
ALTER TABLE journal_sequences ENABLE ROW LEVEL SECURITY;
CREATE POLICY journal_sequences_isolation ON journal_sequences
  USING (tenant_id = app_current_tenant_id()) WITH CHECK (tenant_id = app_current_tenant_id());

-- Runtime privileges: no DELETE on the ledger. Draft lines are replaced by
-- DELETE + INSERT; the line trigger refuses that once the entry is posted.
GRANT SELECT, INSERT, UPDATE ON
  fiscal_years, fiscal_periods, account_groups, accounts, cost_centers,
  journal_entries, journal_sequences
TO alshuyukh_app;
GRANT SELECT, INSERT, UPDATE, DELETE ON journal_entry_lines TO alshuyukh_app;
