import { api } from '../../api';

export interface Warehouse { id: string; code: string; name: string; isActive: boolean }
export interface StockProduct { id: string; sku: string; nameAr: string; trackInventory: boolean; onHand: string }

/**
 * Active warehouses of the (first) company, main warehouse first so it is
 * the default choice everywhere (the server also defaults to it).
 */
export async function loadWarehouses(): Promise<Warehouse[]> {
  const companies = await api<{ data: { id: string }[] }>('GET', '/api/companies');
  const company = companies.data[0];
  if (!company) return [];
  const list = (await api<{ data: Warehouse[] }>('GET', `/api/companies/${company.id}/warehouses`)).data.filter((w) => w.isActive);
  return [...list.filter((w) => w.code === 'MAIN'), ...list.filter((w) => w.code !== 'MAIN')];
}

export async function loadStockedProducts(): Promise<StockProduct[]> {
  return (await api<{ data: StockProduct[] }>('GET', '/api/products?limit=200&productType=GOODS')).data.filter((p) => p.trackInventory);
}

export const qty = (v: string | null | undefined) => (v ? v.replace(/\.?0+$/, '') || '0' : '');

export const MOVEMENT_AR: Record<string, string> = {
  PURCHASE: 'شراء', SALE: 'بيع', SALE_RETURN: 'مرتجع مبيعات', PURCHASE_RETURN: 'مرتجع مشتريات',
  TRANSFER_IN: 'تحويل وارد', TRANSFER_OUT: 'تحويل صادر', ADJUSTMENT_IN: 'تسوية بالزيادة', ADJUSTMENT_OUT: 'تسوية بالنقص',
  CANCELLATION_IN: 'إلغاء (وارد)', CANCELLATION_OUT: 'إلغاء (صادر)',
};
