import { useState, type ReactNode } from 'react';
import { Link } from 'react-router-dom';
import { formatAmount } from '../../money';
import { PageHeader } from '../../ui';

export const today = () => new Intl.DateTimeFormat('en-CA', { timeZone: 'Asia/Riyadh' }).format(new Date());
const pad = (n: number) => String(n).padStart(2, '0');
const lastDay = (y: number, m: number) => new Date(Date.UTC(y, m, 0)).getUTCDate();

/** Date-range presets in Riyadh time. Months are 1-based. */
export function preset(kind: 'month' | 'lastMonth' | 'quarter' | 'year'): { from: string; to: string } {
  const [y, m] = today().split('-').map(Number) as [number, number];
  if (kind === 'month') return { from: `${y}-${pad(m)}-01`, to: `${y}-${pad(m)}-${lastDay(y, m)}` };
  if (kind === 'lastMonth') {
    const [py, pm] = m === 1 ? [y - 1, 12] : [y, m - 1];
    return { from: `${py}-${pad(pm)}-01`, to: `${py}-${pad(pm)}-${lastDay(py, pm)}` };
  }
  if (kind === 'quarter') {
    const qs = Math.floor((m - 1) / 3) * 3 + 1;
    return { from: `${y}-${pad(qs)}-01`, to: `${y}-${pad(qs + 2)}-${lastDay(y, qs + 2)}` };
  }
  return { from: `${y}-01-01`, to: `${y}-12-31` };
}

export function useRange(initial: Parameters<typeof preset>[0] = 'year') {
  const [range, setRange] = useState(() => preset(initial));
  return { range, setRange, query: `dateFrom=${range.from}&dateTo=${range.to}` };
}

export function RangeBar({ range, setRange, children }: { range: { from: string; to: string }; setRange: (r: { from: string; to: string }) => void; children?: ReactNode }) {
  const presets: [Parameters<typeof preset>[0], string][] = [['month', 'هذا الشهر'], ['lastMonth', 'الشهر السابق'], ['quarter', 'هذا الربع'], ['year', 'هذه السنة']];
  return (
    <div className="toolbar no-print">
      <label>من<input type="date" value={range.from} onChange={(e) => setRange({ ...range, from: e.target.value })} /></label>
      <label>إلى<input type="date" value={range.to} onChange={(e) => setRange({ ...range, to: e.target.value })} /></label>
      <div className="chips">
        {presets.map(([k, label]) => <button key={k} type="button" className="btn btn-small" onClick={() => setRange(preset(k))}>{label}</button>)}
      </div>
      {children}
    </div>
  );
}

/** Report page frame: title, back link, print and CSV export. */
export function ReportFrame({ title, subtitle, csv, children }: { title: string; subtitle?: string; csv?: () => (string | number | null | undefined)[][]; children: ReactNode }) {
  return (
    <>
      <PageHeader title={title}>
        <Link className="btn btn-ghost no-print" to="/reports">كل التقارير</Link>
        {csv && <button className="btn no-print" onClick={() => downloadCsv(`alshuyukh-${location.pathname.split('/').pop()}-${today()}`, csv())}>تصدير CSV</button>}
        <button className="btn no-print" onClick={() => window.print()}>طباعة</button>
      </PageHeader>
      {subtitle && <p className="muted print-only">{subtitle}</p>}
      {children}
    </>
  );
}

/**
 * Downloads rows as CSV with a BOM so Excel opens Arabic text correctly. Values are exported as the server sent them.
 * File names stay ASCII: some browsers drop non-ASCII download names.
 */
export function downloadCsv(name: string, rows: (string | number | null | undefined)[][]) {
  const esc = (v: string | number | null | undefined) => {
    let s = v === null || v === undefined ? '' : String(v);
    // Spreadsheet formula injection: a cell starting with = + - @ (or tab/CR) would
    // be evaluated by Excel. Prefix text with ' — numbers such as -150.00 stay numbers.
    if (/^[=+\-@\t\r]/.test(s) && !/^-?\d+(\.\d+)?$/.test(s)) s = `'${s}`;
    return /[",\n\r]/.test(s) ? `"${s.replace(/"/g, '""')}"` : s;
  };
  const blob = new Blob(['﻿' + rows.map((r) => r.map(esc).join(',')).join('\r\n')], { type: 'text/csv;charset=utf-8' });
  const a = document.createElement('a');
  a.href = URL.createObjectURL(blob);
  a.download = `${name}.csv`;
  document.body.append(a);
  a.click();
  a.remove();
  // Revoking immediately can cancel the download in some browsers.
  setTimeout(() => URL.revokeObjectURL(a.href), 1000);
}

export const Amount = ({ v, strong }: { v: string | null | undefined; strong?: boolean }) =>
  <td className={`num ${strong ? 'strong' : ''} ${v?.startsWith('-') ? 'negative' : ''}`} dir="ltr">{formatAmount(v)}</td>;

export const REF_AR: Record<string, string> = {
  MANUAL: 'قيد يدوي', REVERSAL: 'قيد عكسي', YEAR_CLOSING: 'إقفال سنة', SALES_INVOICE: 'فاتورة مبيعات', SALES_RETURN: 'مرتجع مبيعات',
  PURCHASE_INVOICE: 'فاتورة مشتريات', PURCHASE_RETURN: 'مرتجع مشتريات', PAYMENT_RECEIPT: 'سند قبض', PAYMENT_DISBURSEMENT: 'سند صرف',
  EXPENSE: 'مصروف', STOCK_ADJUSTMENT: 'تسوية مخزون', INVENTORY_RESIDUAL: 'فرق تقييم مخزون',
};

/** Link to the document behind a ledger reference, when the app has a page for it. */
export function docLink(type: string, id: string | null | undefined): string | null {
  if (!id) return null;
  const base: Record<string, string> = {
    SALES_INVOICE: '/sales/invoices', SALES_RETURN: '/sales/returns', PURCHASE_INVOICE: '/purchases/invoices',
    PURCHASE_RETURN: '/purchases/returns', EXPENSE: '/expenses/list',
  };
  return base[type] ? `${base[type]}/${id}` : null;
}
