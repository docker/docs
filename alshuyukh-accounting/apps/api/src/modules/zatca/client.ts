/**
 * HTTP client for the ZATCA e-invoicing (Fatoora) API.
 * NOTE: written against ZATCA's published API description; it has not yet been
 * run against the live developer portal, simulation or production gateways.
 */
export type Environment = 'DEVELOPER' | 'SIMULATION' | 'PRODUCTION';

export const BASE_URLS: Record<Environment, string> = {
  DEVELOPER: 'https://gw-fatoora.zatca.gov.sa/e-invoicing/developer-portal',
  SIMULATION: 'https://gw-fatoora.zatca.gov.sa/e-invoicing/simulation',
  PRODUCTION: 'https://gw-fatoora.zatca.gov.sa/e-invoicing/core',
};

export interface Transport { (url: string, init: { method: string; headers: Record<string, string>; body: string }): Promise<{ status: number; text(): Promise<string> }> }

export interface ZatcaResponse {
  status: number;
  body: Record<string, any> | null;
  durationMs: number;
  /** Network failure or timeout: nothing reached ZATCA (or we cannot tell). */
  transportError?: string;
}

export interface Credentials { username: string; password: string }

export class ZatcaClient {
  constructor(private readonly env: Environment, private readonly transport: Transport = defaultTransport, private readonly baseUrl = BASE_URLS[env]) {}

  private async call(path: string, body: unknown, extra: Record<string, string>, auth?: Credentials): Promise<ZatcaResponse> {
    const headers: Record<string, string> = { 'Content-Type': 'application/json', Accept: 'application/json', 'Accept-Language': 'en', 'Accept-Version': 'V2', ...extra };
    if (auth) headers.Authorization = `Basic ${Buffer.from(`${auth.username}:${auth.password}`).toString('base64')}`;
    const started = Date.now();
    try {
      const res = await this.transport(`${this.baseUrl}${path}`, { method: 'POST', headers, body: JSON.stringify(body) });
      const text = await res.text();
      let parsed: Record<string, any> | null = null;
      try { parsed = text ? JSON.parse(text) : null; } catch { parsed = { raw: text.slice(0, 2000) }; }
      return { status: res.status, body: parsed, durationMs: Date.now() - started };
    } catch (e) {
      return { status: 0, body: null, durationMs: Date.now() - started, transportError: (e as Error).message };
    }
  }

  /** Step 1 of onboarding: CSR + one-time password from the Fatoora portal → compliance CSID. */
  complianceCsid(csrPem: string, otp: string) {
    return this.call('/compliance', { csr: Buffer.from(csrPem).toString('base64') }, { OTP: otp });
  }
  /** Step 2: sample documents signed with the compliance CSID. */
  complianceCheck(auth: Credentials, doc: { invoiceHash: string; uuid: string; invoice: string }) {
    return this.call('/compliance/invoices', doc, {}, auth);
  }
  /** Step 3: production CSID, once the compliance checks pass. */
  productionCsid(auth: Credentials, complianceRequestId: string) {
    return this.call('/production/csids', { compliance_request_id: complianceRequestId }, {}, auth);
  }
  /** Simplified invoices and their notes: report within 24 hours. */
  report(auth: Credentials, doc: { invoiceHash: string; uuid: string; invoice: string }) {
    return this.call('/invoices/reporting/single', doc, { 'Clearance-Status': '0' }, auth);
  }
  /** Standard invoices and their notes: must be cleared before they are given to the buyer. */
  clear(auth: Credentials, doc: { invoiceHash: string; uuid: string; invoice: string }) {
    return this.call('/invoices/clearance/single', doc, { 'Clearance-Status': '1' }, auth);
  }
}

const defaultTransport: Transport = async (url, init) => {
  const res = await fetch(url, { ...init, signal: AbortSignal.timeout(30_000) });
  return { status: res.status, text: () => res.text() };
};

/** The transport used by the app; tests replace it. */
let currentTransport: Transport = defaultTransport;
let gatewayOverride: string | null = null;
export const setTransport = (t: Transport | null) => { currentTransport = t ?? defaultTransport; };
/** Points every environment at another gateway, e.g. a local stand-in. */
export const setGatewayUrl = (url: string | null) => { gatewayOverride = url; };
export const clientFor = (env: Environment) =>
  new ZatcaClient(env, (url, init) => currentTransport(url, init), gatewayOverride ? `${gatewayOverride.replace(/\/$/, '')}/${env.toLowerCase()}` : BASE_URLS[env]);

export interface ValidationMessage { type?: string; code?: string; category?: string; message?: string; status?: string }

/** Validation messages from a reporting/clearance/compliance response. */
export function validationMessages(body: Record<string, any> | null): { level: 'ERROR' | 'WARNING' | 'INFO'; code: string | null; category: string | null; message: string }[] {
  const v = body?.validationResults;
  if (!v) {
    // Errors outside validation (authentication, malformed request).
    if (body?.message || body?.errors) {
      const list = Array.isArray(body.errors) ? body.errors : [{ message: body.message }];
      return list.map((e: ValidationMessage | string) => ({ level: 'ERROR' as const, code: typeof e === 'string' ? null : e.code ?? null, category: null, message: typeof e === 'string' ? e : String(e.message ?? JSON.stringify(e)) }));
    }
    return [];
  }
  const map = (list: ValidationMessage[] | undefined, level: 'ERROR' | 'WARNING' | 'INFO') =>
    (list ?? []).map((m) => ({ level, code: m.code ?? null, category: m.category ?? null, message: String(m.message ?? '') }));
  return [...map(v.errorMessages, 'ERROR'), ...map(v.warningMessages, 'WARNING'), ...map(v.infoMessages, 'INFO')];
}
