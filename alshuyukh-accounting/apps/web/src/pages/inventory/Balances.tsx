import { useState } from 'react';
import { api } from '../../api';
import { formatAmount } from '../../money';
import { ErrorBox, useLoad } from '../../ui';
import { loadWarehouses, qty } from './shared';

interface Row { productId: string; sku: string; productName: string; unitCode: string; warehouseName: string; quantity: string; value: string; averageCost: string }

export default function Balances() {
  const [warehouseId, setWarehouseId] = useState('');
  const warehouses = useLoad(loadWarehouses);
  const rows = useLoad(() => api<{ data: Row[] }>('GET', `/api/inventory/balances${warehouseId ? `?warehouseId=${warehouseId}` : ''}`), [warehouseId]);
  return (
    <div className="card">
      <div className="toolbar">
        <select value={warehouseId} onChange={(e) => setWarehouseId(e.target.value)} aria-label="المستودع">
          <option value="">كل المستودعات</option>
          {warehouses.data?.map((w) => <option key={w.id} value={w.id}>{w.name}</option>)}
        </select>
      </div>
      <ErrorBox error={rows.error} />
      <div className="table-wrap">
        <table className="table">
          <thead><tr><th>الرمز</th><th>الصنف</th><th>المستودع</th><th className="num">الكمية</th><th className="num">متوسط التكلفة</th><th className="num">القيمة</th></tr></thead>
          <tbody>
            {!rows.data && !rows.error && <tr><td colSpan={6} className="muted empty-row">جارٍ التحميل…</td></tr>}
            {rows.data?.data.map((r) => (
              <tr key={`${r.productId}-${r.warehouseName}`}>
                <td dir="ltr" className="code">{r.sku}</td><td>{r.productName}</td><td>{r.warehouseName}</td>
                <td className="num" dir="ltr">{qty(r.quantity)} {r.unitCode}</td>
                <td className="num" dir="ltr">{formatAmount(Number(r.averageCost).toFixed(2))}</td>
                <td className="num" dir="ltr">{formatAmount(r.value)}</td>
              </tr>
            ))}
            {rows.data?.data.length === 0 && <tr><td colSpan={6} className="muted empty-row">لا يوجد مخزون</td></tr>}
          </tbody>
        </table>
      </div>
    </div>
  );
}
