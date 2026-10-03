export interface DocConfig {
  key: 'SALES_QUOTE' | 'SALES_INVOICE' | 'SALES_RETURN' | 'PURCHASE_ORDER' | 'PURCHASE_INVOICE' | 'PURCHASE_RETURN';
  api: string;          // API path segment
  route: string;        // UI route, relative to the section
  title: string;
  singular: string;
  newLabel: string;
  party: 'customer' | 'supplier';
  legal: boolean;
  isReturn: boolean;
  perm: { view: string; create: string; issue: string; cancel: string };
}

const SALES = { view: 'invoice.view', create: 'invoice.create', issue: 'invoice.post', cancel: 'invoice.cancel' };
const PURCHASE = { view: 'purchase.view', create: 'purchase.create', issue: 'purchase.post', cancel: 'purchase.cancel' };

export const DOCS: Record<DocConfig['key'], DocConfig> = {
  SALES_INVOICE: { key: 'SALES_INVOICE', api: 'invoices', route: 'invoices', title: 'فواتير المبيعات', singular: 'فاتورة', newLabel: 'فاتورة جديدة', party: 'customer', legal: true, isReturn: false, perm: SALES },
  SALES_QUOTE: { key: 'SALES_QUOTE', api: 'sales-quotes', route: 'quotes', title: 'عروض الأسعار', singular: 'عرض سعر', newLabel: 'عرض سعر جديد', party: 'customer', legal: false, isReturn: false, perm: SALES },
  SALES_RETURN: { key: 'SALES_RETURN', api: 'sales-returns', route: 'returns', title: 'مرتجعات المبيعات', singular: 'إشعار دائن', newLabel: '', party: 'customer', legal: true, isReturn: true, perm: SALES },
  PURCHASE_INVOICE: { key: 'PURCHASE_INVOICE', api: 'purchase-invoices', route: 'invoices', title: 'فواتير المشتريات', singular: 'فاتورة مشتريات', newLabel: 'فاتورة مشتريات جديدة', party: 'supplier', legal: true, isReturn: false, perm: PURCHASE },
  PURCHASE_ORDER: { key: 'PURCHASE_ORDER', api: 'purchase-orders', route: 'orders', title: 'أوامر الشراء', singular: 'أمر شراء', newLabel: 'أمر شراء جديد', party: 'supplier', legal: false, isReturn: false, perm: PURCHASE },
  PURCHASE_RETURN: { key: 'PURCHASE_RETURN', api: 'purchase-returns', route: 'returns', title: 'مرتجعات المشتريات', singular: 'إشعار مدين', newLabel: '', party: 'supplier', legal: true, isReturn: true, perm: PURCHASE },
};

export const STATUS_AR: Record<string, string> = {
  DRAFT: 'مسودة', ISSUED: 'صادرة', POSTED: 'مرحّلة', PARTIALLY_PAID: 'مدفوعة جزئيًا', PAID: 'مدفوعة', CANCELLED: 'ملغاة',
  RETURNED: 'مرتجعة', SENT: 'مرسل', ACCEPTED: 'مقبول', REJECTED: 'مرفوض', CONVERTED: 'محوّل لفاتورة', APPROVED: 'معتمد',
  VOIDED: 'ملغى',
};

export const statusClass = (s: string) =>
  ({ DRAFT: 'status-draft', PAID: 'status-posted', ISSUED: 'status-open', POSTED: 'status-open', PARTIALLY_PAID: 'status-partial',
     CANCELLED: 'status-reversed', VOIDED: 'status-reversed', RETURNED: 'status-reversed' } as Record<string, string>)[s] ?? 'status-open';

export interface DocLine {
  id: string; lineNo: number; productId: string | null; productName: string | null; sku: string | null; accountId: string | null;
  sourceItemId: string | null; description: string | null; quantity: string; unitCode: string | null; unitPrice: string;
  grossAmount: string; discountBasis: string; discountAmount: string; netAmount: string; vatCategory: string; vatRate: string;
  vatAmount: string; totalAmount: string;
}

export interface Doc {
  id: string; number: string | null; date: string; partyId: string; partyName: string; partyCode: string; status: string;
  pricesIncludeVat: boolean; subtotal: string; discountTotal: string; taxableAmount: string; taxAmount: string; total: string;
  notes: string | null; journalEntryId: string | null; cancelReason: string | null;
  dueDate?: string | null; paidAmount?: string; returnedAmount?: string; remainingAmount?: string; invoiceKind?: string;
  supplierInvoiceNumber?: string | null; originalInvoiceId?: string; originalInvoiceNumber?: string; reason?: string;
  appliedAmount?: string; refundedAmount?: string; convertedInvoiceId?: string | null; validUntil?: string | null;
  lines: DocLine[];
  allocations?: { id: string; amount: string; paymentId: string; paymentNumber: string; paymentDate: string }[];
  returns?: { id: string; number: string | null; date: string; status: string; total: string }[];
}
