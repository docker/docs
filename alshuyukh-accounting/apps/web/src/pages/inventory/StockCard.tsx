import { useState } from 'react';
import { api } from '../../api';
import { formatAmount } from '../../money';
import { ErrorBox, useLoad } from '../../ui';
import { MOVEMENT_AR, loadStockedProducts, loadWarehouses, qty } from './shared';

interface Movement {
  id: string; date: string; type: string; direction: string; warehouseName: string; quantity: string; unitCost: string;
  totalCost: string; balanceQuantity: string; balanceValue: string; referenceType: string;
}

export default function StockCard() {
  const products = useLoad(loadStockedProducts);
  const warehouses = useLoad(loadWarehouses);
  const [productId, setProductId] = useState('');
  const [warehouseId, setWarehouseId] = useState('');
  const moves = useLoad(() => (productId
    ? api<{ data: Movement[] }>('GET', `/api/inventory/movements?productId=${productId}${warehouseId ? `&warehouseId=${warehouseId}` : ''}&limit=500`)
    : Promise.resolve({ data: [] as Movement[] })), [productId, warehouseId]);
  return (
    <div className="card">
      <div className="toolbar">
        <select value={productId} onChange={(e) => setProductId(e.target.value)} aria-label="الصنف" className="grow">
          <option value="">— اختر الصنف —</option>
          {products.data?.map((p) => <option key={p.id} value={p.id}>{p.sku} — {p.nameAr}</option>)}
        </select>
        <select value={warehouseId} onChange={(e) => setWarehouseId(e.target.value)} aria-label="المستودع">
          <option value="">كل المستودعات</option>
          {warehouses.data?.map((w) => <option key={w.id} value={w.id}>{w.name}</option>)}
        </select>
      </div>
      <ErrorBox error={moves.error} />
      {productId && (
        <div className="table-wrap">
          <table className="table">
            <thead><tr><th>التاريخ</th><th>الحركة</th><th>المستودع</th><th className="num">وارد</th><th className="num">صادر</th><th className="num">تكلفة الوحدة</th><th className="num">القيمة</th><th className="num">الرصيد</th><th className="num">قيمة الرصيد</th></tr></thead>
            <tbody>
              {moves.data?.data.map((m) => (
                <tr key={m.id}>
                  <td dir="ltr">{m.date}</td><td>{MOVEMENT_AR[m.type] ?? m.type}</td><td>{m.warehouseName}</td>
                  <td className="num" dir="ltr">{m.direction === 'IN' ? qty(m.quantity) : ''}</td>
                  <td className="num" dir="ltr">{m.direction === 'OUT' ? qty(m.quantity) : ''}</td>
                  <td className="num" dir="ltr">{Number(m.unitCost).toFixed(4)}</td>
                  <td className="num" dir="ltr">{formatAmount(m.totalCost)}</td>
                  <td className="num" dir="ltr">{qty(m.balanceQuantity)}</td>
                  <td className="num" dir="ltr">{formatAmount(m.balanceValue)}</td>
                </tr>
              ))}
              {moves.data?.data.length === 0 && <tr><td colSpan={9} className="muted empty-row">لا توجد حركات</td></tr>}
            </tbody>
          </table>
        </div>
      )}
      <p className="muted small">الأحدث أولًا. الرصيد بعد كل حركة في مستودعها.</p>
    </div>
  );
}
