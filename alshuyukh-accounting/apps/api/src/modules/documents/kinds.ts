/**
 * The six commercial documents share one implementation. Each kind declares
 * its tables, party, numbering, permissions, and lifecycle here.
 */

export type DocKey = 'SALES_QUOTE' | 'SALES_INVOICE' | 'SALES_RETURN' | 'PURCHASE_ORDER' | 'PURCHASE_INVOICE' | 'PURCHASE_RETURN';

export interface DocKind {
  key: DocKey;
  path: string;
  table: string;
  itemsTable: string;
  party: 'customer' | 'supplier';
  /** Invoices and returns: number at issue, journal entry, settlement tracking. */
  legal: boolean;
  isReturn: boolean;
  /** For returns: the invoice kind being returned. */
  originalKind?: DocKey;
  prefix: string;
  /** Status set when the document is issued / posted. */
  issuedStatus: string;
  entity: string;
  label: string;
  labelAr: string;
  perm: { view: string; create: string; edit: string; delete: string; issue: string; cancel: string };
}

const SALES_PERM = { view: 'invoice.view', create: 'invoice.create', edit: 'invoice.edit', delete: 'invoice.delete', issue: 'invoice.post', cancel: 'invoice.cancel' };
const PURCHASE_PERM = { view: 'purchase.view', create: 'purchase.create', edit: 'purchase.create', delete: 'purchase.create', issue: 'purchase.post', cancel: 'purchase.cancel' };

export const KINDS: Record<DocKey, DocKind> = {
  SALES_QUOTE: {
    key: 'SALES_QUOTE', path: 'sales-quotes', table: 'sales_quotes', itemsTable: 'sales_quote_items', party: 'customer',
    legal: false, isReturn: false, prefix: 'QT', issuedStatus: 'SENT', entity: 'sales_quote', label: 'Quote', labelAr: 'عرض سعر', perm: SALES_PERM,
  },
  SALES_INVOICE: {
    key: 'SALES_INVOICE', path: 'invoices', table: 'sales_invoices', itemsTable: 'sales_invoice_items', party: 'customer',
    legal: true, isReturn: false, prefix: 'INV', issuedStatus: 'ISSUED', entity: 'sales_invoice', label: 'Invoice', labelAr: 'فاتورة مبيعات', perm: SALES_PERM,
  },
  SALES_RETURN: {
    key: 'SALES_RETURN', path: 'sales-returns', table: 'sales_returns', itemsTable: 'sales_return_items', party: 'customer',
    legal: true, isReturn: true, originalKind: 'SALES_INVOICE', prefix: 'CN', issuedStatus: 'ISSUED', entity: 'sales_return', label: 'Sales return', labelAr: 'مرتجع مبيعات (إشعار دائن)', perm: SALES_PERM,
  },
  PURCHASE_ORDER: {
    key: 'PURCHASE_ORDER', path: 'purchase-orders', table: 'purchase_orders', itemsTable: 'purchase_order_items', party: 'supplier',
    legal: false, isReturn: false, prefix: 'PO', issuedStatus: 'APPROVED', entity: 'purchase_order', label: 'Purchase order', labelAr: 'أمر شراء', perm: PURCHASE_PERM,
  },
  PURCHASE_INVOICE: {
    key: 'PURCHASE_INVOICE', path: 'purchase-invoices', table: 'purchase_invoices', itemsTable: 'purchase_invoice_items', party: 'supplier',
    legal: true, isReturn: false, prefix: 'PINV', issuedStatus: 'POSTED', entity: 'purchase_invoice', label: 'Purchase invoice', labelAr: 'فاتورة مشتريات', perm: PURCHASE_PERM,
  },
  PURCHASE_RETURN: {
    key: 'PURCHASE_RETURN', path: 'purchase-returns', table: 'purchase_returns', itemsTable: 'purchase_return_items', party: 'supplier',
    legal: true, isReturn: true, originalKind: 'PURCHASE_INVOICE', prefix: 'DN', issuedStatus: 'ISSUED', entity: 'purchase_return', label: 'Purchase return', labelAr: 'مرتجع مشتريات (إشعار مدين)', perm: PURCHASE_PERM,
  },
};

export const partyColumn = (k: DocKind) => `${k.party}_id`;
export const partyTable = (k: DocKind) => (k.party === 'customer' ? 'customers' : 'suppliers');
