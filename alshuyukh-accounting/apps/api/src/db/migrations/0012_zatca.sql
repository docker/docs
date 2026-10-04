-- Phase 8: ZATCA e-invoicing (Fatoora, phase 2 "integration").
--
-- One EGS unit (e-invoice generation solution) per company or branch holds
-- the cryptographic stamp: an EC secp256k1 key pair, its CSR and the CSIDs
-- (certificates) issued by ZATCA. Every simplified/standard invoice and
-- credit note it produces gets the next ICV (counter) and the hash of the
-- previous document (PIH), forming a chain that cannot fork.

-- Structured national address of the seller, required in the e-invoice.
ALTER TABLE companies
  ADD COLUMN building_number   text CHECK (building_number ~ '^[0-9]{4}$'),
  ADD COLUMN street            text CHECK (length(btrim(street)) BETWEEN 1 AND 200),
  ADD COLUMN district          text CHECK (length(btrim(district)) BETWEEN 1 AND 200),
  ADD COLUMN postal_code       text CHECK (postal_code ~ '^[0-9]{5}$'),
  ADD COLUMN additional_number text CHECK (additional_number ~ '^[0-9]{4}$');

-- Zero-rated and exempt supplies must carry a ZATCA exemption reason code.
ALTER TABLE products
  ADD COLUMN vat_exemption_code   text CHECK (vat_exemption_code ~ '^VATEX-SA-[A-Z0-9-]{2,12}$'),
  ADD COLUMN vat_exemption_reason text CHECK (length(btrim(vat_exemption_reason)) BETWEEN 1 AND 300);

CREATE TABLE zatca_devices (
  id                     uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id              uuid NOT NULL,
  company_id             uuid NOT NULL,
  branch_id              uuid,
  name                   text NOT NULL CHECK (length(btrim(name)) BETWEEN 2 AND 100),
  -- DEVELOPER: developer portal (sandbox), SIMULATION, PRODUCTION (core).
  environment            text NOT NULL CHECK (environment IN ('DEVELOPER', 'SIMULATION', 'PRODUCTION')),
  -- CSR subject / SAN fields required by ZATCA.
  egs_serial             text NOT NULL,             -- 1-<solution>|2-<model>|3-<serial>
  common_name            text NOT NULL,
  organization_unit      text NOT NULL,
  invoice_types          text NOT NULL DEFAULT '1100' CHECK (invoice_types ~ '^[01]{4}$' AND invoice_types <> '0000'),
  registered_address     text NOT NULL,
  business_category      text NOT NULL,
  -- NEW → CSR generated; COMPLIANCE → compliance CSID received;
  -- ACTIVE → production CSID received, signs invoices; REVOKED.
  status                 text NOT NULL DEFAULT 'NEW' CHECK (status IN ('NEW', 'COMPLIANCE', 'ACTIVE', 'REVOKED')),
  -- Secrets are encrypted by the application (AES-256-GCM); never returned by the API.
  private_key_enc        text NOT NULL,
  public_key             text NOT NULL,
  csr                    text NOT NULL,
  compliance_csid        text,
  compliance_secret_enc  text,
  compliance_request_id  text,
  compliance_checks      jsonb NOT NULL DEFAULT '{}',
  production_csid        text,
  production_secret_enc  text,
  certificate_serial     text,
  certificate_issuer     text,
  certificate_expires_at timestamptz,
  -- Chain state: last ICV used and the hash of the last document.
  icv                    bigint NOT NULL DEFAULT 0 CHECK (icv >= 0),
  last_hash              text NOT NULL DEFAULT 'NWZlY2ViNjZmZmM4NmYzOGQ5NTI3ODZjNmQ2OTZjNzljMmRiYzIzOWRkNGU5MWI0NjcyOWQ3M2EyN2ZiNTdlOQ==',
  activated_at           timestamptz,
  revoked_at             timestamptz,
  created_by             uuid REFERENCES users (id),
  created_at             timestamptz NOT NULL DEFAULT now(),
  updated_at             timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT zatca_devices_company_fk FOREIGN KEY (company_id, tenant_id) REFERENCES companies (id, tenant_id),
  CONSTRAINT zatca_devices_branch_fk FOREIGN KEY (branch_id, company_id) REFERENCES branches (id, company_id),
  CONSTRAINT zatca_devices_id_company_uq UNIQUE (id, company_id),
  CONSTRAINT zatca_devices_active_ck CHECK (status <> 'ACTIVE' OR (production_csid IS NOT NULL AND production_secret_enc IS NOT NULL))
);
-- At most one live unit per company (branch_id NULL) and per branch.
CREATE UNIQUE INDEX zatca_devices_scope_uq ON zatca_devices (company_id, COALESCE(branch_id, '00000000-0000-0000-0000-000000000000'::uuid))
  WHERE status <> 'REVOKED';
CREATE TRIGGER zatca_devices_updated_at BEFORE UPDATE ON zatca_devices FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TABLE zatca_invoices (
  id                   uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id            uuid NOT NULL,
  company_id           uuid NOT NULL,
  device_id            uuid NOT NULL,
  document_type        text NOT NULL CHECK (document_type IN ('SALES_INVOICE', 'SALES_RETURN', 'SALES_DEBIT_NOTE')),
  document_id          uuid NOT NULL,
  document_number      text NOT NULL,
  invoice_kind         text NOT NULL CHECK (invoice_kind IN ('STANDARD', 'SIMPLIFIED')),
  type_code            text NOT NULL CHECK (type_code IN ('388', '381', '383')),
  subtype              text NOT NULL CHECK (subtype ~ '^0[12][01]{5}$'),
  uuid                 uuid NOT NULL UNIQUE,
  icv                  bigint NOT NULL CHECK (icv > 0),
  invoice_hash         text NOT NULL,
  previous_hash        text NOT NULL,
  issue_date           date NOT NULL,
  issue_time           text NOT NULL CHECK (issue_time ~ '^[0-9]{2}:[0-9]{2}:[0-9]{2}$'),
  qr                   text NOT NULL,
  -- PENDING: waiting for clearance (standard) or reporting (simplified).
  status               text NOT NULL DEFAULT 'PENDING' CHECK (status IN ('PENDING', 'REPORTED', 'CLEARED', 'REJECTED')),
  has_warnings         boolean NOT NULL DEFAULT false,
  attempts             integer NOT NULL DEFAULT 0,
  next_attempt_at      timestamptz NOT NULL DEFAULT now(),
  last_error           text,
  submitted_at         timestamptz,
  created_at           timestamptz NOT NULL DEFAULT now(),
  updated_at           timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT zatca_invoices_company_fk FOREIGN KEY (company_id, tenant_id) REFERENCES companies (id, tenant_id),
  CONSTRAINT zatca_invoices_device_fk FOREIGN KEY (device_id, company_id) REFERENCES zatca_devices (id, company_id),
  CONSTRAINT zatca_invoices_document_uq UNIQUE (document_type, document_id),
  CONSTRAINT zatca_invoices_icv_uq UNIQUE (device_id, icv)
);
CREATE INDEX zatca_invoices_company_idx ON zatca_invoices (tenant_id, company_id, created_at DESC);
CREATE INDEX zatca_invoices_due_idx ON zatca_invoices (next_attempt_at) WHERE status = 'PENDING';
CREATE TRIGGER zatca_invoices_updated_at BEFORE UPDATE ON zatca_invoices FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- The identity of a generated e-invoice never changes once it is in the chain.
CREATE OR REPLACE FUNCTION zatca_invoices_guard() RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
  IF TG_OP = 'DELETE' THEN RAISE EXCEPTION 'zatca_invoices rows cannot be deleted'; END IF;
  IF (NEW.uuid, NEW.icv, NEW.invoice_hash, NEW.previous_hash, NEW.device_id, NEW.document_id, NEW.document_type, NEW.qr, NEW.issue_date, NEW.issue_time)
     IS DISTINCT FROM
     (OLD.uuid, OLD.icv, OLD.invoice_hash, OLD.previous_hash, OLD.device_id, OLD.document_id, OLD.document_type, OLD.qr, OLD.issue_date, OLD.issue_time) THEN
    RAISE EXCEPTION 'A generated e-invoice cannot be modified';
  END IF;
  IF OLD.status IN ('REPORTED', 'CLEARED') AND NEW.status <> OLD.status THEN
    RAISE EXCEPTION 'A reported or cleared e-invoice cannot change status';
  END IF;
  RETURN NEW;
END $$;
CREATE TRIGGER zatca_invoices_guard BEFORE UPDATE OR DELETE ON zatca_invoices FOR EACH ROW EXECUTE FUNCTION zatca_invoices_guard();

-- XML artefacts: the signed invoice we generated and, for standard invoices,
-- the cleared invoice returned by ZATCA (which is the legal document).
CREATE TABLE zatca_documents (
  id                uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id         uuid NOT NULL,
  zatca_invoice_id  uuid NOT NULL REFERENCES zatca_invoices (id),
  kind              text NOT NULL CHECK (kind IN ('SIGNED', 'CLEARED')),
  content           text NOT NULL,
  sha256            text NOT NULL,
  created_at        timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT zatca_documents_kind_uq UNIQUE (zatca_invoice_id, kind)
);
CREATE TRIGGER zatca_documents_append_only BEFORE UPDATE OR DELETE ON zatca_documents FOR EACH ROW EXECUTE FUNCTION prevent_modification();

-- Every call to the ZATCA API (onboarding, compliance checks, reporting, clearance).
CREATE TABLE zatca_submissions (
  id                uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id         uuid NOT NULL,
  device_id         uuid NOT NULL REFERENCES zatca_devices (id),
  zatca_invoice_id  uuid REFERENCES zatca_invoices (id),
  operation         text NOT NULL CHECK (operation IN ('COMPLIANCE_CSID', 'COMPLIANCE_CHECK', 'PRODUCTION_CSID', 'REPORTING', 'CLEARANCE')),
  http_status       integer,
  outcome           text NOT NULL CHECK (outcome IN ('SUCCESS', 'WARNING', 'REJECTED', 'TRANSPORT_ERROR')),
  -- Response without secrets or certificates.
  response          jsonb,
  duration_ms       integer,
  created_by        uuid REFERENCES users (id),
  created_at        timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX zatca_submissions_invoice_idx ON zatca_submissions (zatca_invoice_id, created_at);
CREATE INDEX zatca_submissions_device_idx ON zatca_submissions (device_id, created_at);
CREATE TRIGGER zatca_submissions_append_only BEFORE UPDATE OR DELETE ON zatca_submissions FOR EACH ROW EXECUTE FUNCTION prevent_modification();

CREATE TABLE zatca_errors (
  id                uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id         uuid NOT NULL,
  submission_id     uuid NOT NULL REFERENCES zatca_submissions (id),
  zatca_invoice_id  uuid REFERENCES zatca_invoices (id),
  level             text NOT NULL CHECK (level IN ('ERROR', 'WARNING', 'INFO')),
  code              text,
  category          text,
  message           text NOT NULL,
  created_at        timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX zatca_errors_invoice_idx ON zatca_errors (zatca_invoice_id);
CREATE TRIGGER zatca_errors_append_only BEFORE UPDATE OR DELETE ON zatca_errors FOR EACH ROW EXECUTE FUNCTION prevent_modification();

-- The hash chain. One row per ICV; a previous hash can be followed only once,
-- so the chain of a device cannot fork.
CREATE TABLE invoice_hashes (
  id                uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id         uuid NOT NULL,
  company_id        uuid NOT NULL,
  device_id         uuid NOT NULL,
  icv               bigint NOT NULL CHECK (icv > 0),
  invoice_hash      text NOT NULL,
  previous_hash     text NOT NULL,
  zatca_invoice_id  uuid NOT NULL UNIQUE REFERENCES zatca_invoices (id),
  created_at        timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT invoice_hashes_device_fk FOREIGN KEY (device_id, company_id) REFERENCES zatca_devices (id, company_id),
  CONSTRAINT invoice_hashes_icv_uq UNIQUE (device_id, icv),
  CONSTRAINT invoice_hashes_no_fork_uq UNIQUE (device_id, previous_hash)
);
CREATE TRIGGER invoice_hashes_append_only BEFORE UPDATE OR DELETE ON invoice_hashes FOR EACH ROW EXECUTE FUNCTION prevent_modification();

-- Sales returns record which kind of invoice they credit (standard or simplified).
ALTER TABLE sales_returns ADD COLUMN invoice_kind text CHECK (invoice_kind IN ('STANDARD', 'SIMPLIFIED'));

-- Row-Level Security and privileges ---------------------------------------------------
DO $$
DECLARE t text;
BEGIN
  FOREACH t IN ARRAY ARRAY['zatca_devices', 'zatca_invoices', 'zatca_documents', 'zatca_submissions', 'zatca_errors', 'invoice_hashes'] LOOP
    EXECUTE format('ALTER TABLE %I ENABLE ROW LEVEL SECURITY', t);
    EXECUTE format('CREATE POLICY %I ON %I USING (tenant_id = app_current_tenant_id()) WITH CHECK (tenant_id = app_current_tenant_id())',
                   t || '_isolation', t);
  END LOOP;
END
$$;
GRANT SELECT, INSERT, UPDATE ON zatca_devices, zatca_invoices TO alshuyukh_app;
GRANT SELECT, INSERT ON zatca_documents, zatca_submissions, zatca_errors, invoice_hashes TO alshuyukh_app;

-- The background submitter needs to find due e-invoices across tenants. This
-- function returns only identifiers; the work itself runs in each tenant's context.
CREATE OR REPLACE FUNCTION zatca_due_invoices(max_rows integer)
RETURNS TABLE (tenant_id uuid, id uuid) LANGUAGE sql STABLE SECURITY DEFINER SET search_path = public AS $$
  SELECT z.tenant_id, z.id FROM zatca_invoices z
   WHERE z.status = 'PENDING' AND z.next_attempt_at <= now()
   ORDER BY z.next_attempt_at LIMIT LEAST(max_rows, 100)
$$;
REVOKE ALL ON FUNCTION zatca_due_invoices(integer) FROM PUBLIC;
GRANT EXECUTE ON FUNCTION zatca_due_invoices(integer) TO alshuyukh_app;
