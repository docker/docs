import type { FastifyInstance } from 'fastify';
import { client, dateInYear, type Session } from './helpers.js';

export interface Commerce {
  api: ReturnType<typeof client>;
  date: string;
  customer(extra?: Record<string, unknown>): Promise<{ id: string }>;
  supplier(extra?: Record<string, unknown>): Promise<{ id: string }>;
  /** A service by default (no stock); use goods() for stocked items. */
  product(extra?: Record<string, unknown>): Promise<{ id: string }>;
  goods(extra?: Record<string, unknown>): Promise<{ id: string }>;
  invoice(customerId: string, lines: unknown[], extra?: Record<string, unknown>): Promise<Record<string, any>>;
  postInvoice(customerId: string, lines: unknown[], extra?: Record<string, unknown>): Promise<Record<string, any>>;
  purchase(supplierId: string, lines: unknown[], extra?: Record<string, unknown>): Promise<Record<string, any>>;
  method(code: string): Promise<string>;
  journal(id: string): Promise<{ accountCode: string; debit: string; credit: string; customerId: string | null; supplierId: string | null }[]>;
}

let seq = 0;
export async function commerce(app: FastifyInstance, s: Session): Promise<Commerce> {
  const api = client(app, s.token);
  const date = await dateInYear(app, s.token, 3, 10);
  const ok = async (p: Promise<{ statusCode: number; json(): any; body: string }>, status = 201) => {
    const r = await p;
    if (r.statusCode !== status) throw new Error(`expected ${status}, got ${r.statusCode}: ${r.body}`);
    return r.json();
  };
  return {
    api,
    date,
    customer: (extra = {}) => ok(api.post('/api/customers', { nameAr: `عميل ${++seq}`, ...extra })),
    supplier: (extra = {}) => ok(api.post('/api/suppliers', { nameAr: `مورد ${++seq}`, ...extra })),
    product: (extra = {}) => ok(api.post('/api/products', { nameAr: `خدمة ${++seq}`, productType: 'SERVICE', salePrice: '100', purchasePrice: '60', ...extra })),
    goods: (extra = {}) => ok(api.post('/api/products', { nameAr: `صنف ${++seq}`, productType: 'GOODS', salePrice: '100', purchasePrice: '60', ...extra })),
    invoice: (customerId, lines, extra = {}) => ok(api.post('/api/invoices', { partyId: customerId, docDate: date, lines, ...extra })),
    async postInvoice(customerId, lines, extra = {}) {
      const inv = await ok(api.post('/api/invoices', { partyId: customerId, docDate: date, lines, ...extra }));
      return ok(api.post(`/api/invoices/${inv.id}/post`), 200);
    },
    async purchase(supplierId, lines, extra = {}) {
      const inv = await ok(api.post('/api/purchase-invoices', { partyId: supplierId, docDate: date, lines, ...extra }));
      return ok(api.post(`/api/purchase-invoices/${inv.id}/post`), 200);
    },
    async method(code) {
      const m = (await api.get('/api/payment-methods')).json().data.find((x: { code: string }) => x.code === code);
      return m.id;
    },
    async journal(id) {
      return (await api.get(`/api/journal-entries/${id}`)).json().lines;
    },
  };
}

export const balanceOf = async (api: Commerce['api'], path: 'customers' | 'suppliers', id: string) =>
  (await api.get(`/api/${path}/${id}`)).json().balance as string;
