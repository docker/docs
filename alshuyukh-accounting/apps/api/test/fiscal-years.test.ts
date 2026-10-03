import { afterAll, beforeAll, describe, expect, it } from 'vitest';
import { chart, client, cr, dr, fiscalYears, register, setupApp, type TestContext } from './helpers.js';

let t: TestContext;
beforeAll(async () => { t = await setupApp(); });
afterAll(async () => { await t.close(); });

describe('fiscal years', () => {
  it('creates the next year and rejects overlaps and invalid starts', async () => {
    const s = await register(t.app);
    const api = client(t.app, s.token);
    const [current] = await fiscalYears(t.app, s.token);
    const nextStart = `${Number(current!.startDate.slice(0, 4)) + 1}-01-01`;

    const next = await api.post('/api/fiscal-years', { startDate: nextStart });
    expect(next.statusCode).toBe(201);
    const years = await fiscalYears(t.app, s.token);
    expect(years).toHaveLength(2);
    expect(years[0]).toMatchObject({ startDate: nextStart, endDate: `${nextStart.slice(0, 4)}-12-31` });
    expect(years[0]!.periods).toHaveLength(12);

    expect((await api.post('/api/fiscal-years', { startDate: nextStart })).statusCode).toBe(409);
    expect((await api.post('/api/fiscal-years', { startDate: `${nextStart.slice(0, 4)}-06-01` })).statusCode).toBe(409);
    expect((await api.post('/api/fiscal-years', { startDate: '2040-01-15' })).statusCode).toBe(400);
    expect((await api.post('/api/fiscal-years', { startDate: '2040-02-30' })).statusCode).toBe(400);
  });

  it('supports a short first year with an explicit end date', async () => {
    const s = await register(t.app);
    const res = await client(t.app, s.token).post('/api/fiscal-years', { startDate: '2040-07-01', endDate: '2040-12-31' });
    expect(res.statusCode).toBe(201);
    const year = (await fiscalYears(t.app, s.token)).find((y) => y.startDate === '2040-07-01');
    expect(year!.periods).toHaveLength(6);
  });

  it('locks the fiscal start month once a fiscal year exists', async () => {
    const s = await register(t.app);
    const res = await client(t.app, s.token).patch('/api/settings/tenant', { fiscalYearStartMonth: 7 });
    expect(res.statusCode).toBe(409);
    expect(res.json().error.code).toBe('FISCAL_YEARS_EXIST');
  });
});

describe('closing a fiscal year', () => {
  it('moves the net profit to retained earnings and closes every period', async () => {
    const s = await register(t.app);
    const api = client(t.app, s.token);
    const c = await chart(t.app, s.token);
    const [year] = await fiscalYears(t.app, s.token);
    const y = year!.startDate.slice(0, 4);
    const post = (entryDate: string, lines: unknown[]) =>
      api.post('/api/journal-entries', { entryDate, description: 'حركة', lines, post: true }).then((r) => { expect(r.statusCode).toBe(201); return r.json(); });

    await post(`${y}-01-10`, [dr(c.byKey.get('CASH')!, '50000.00'), cr(c.byKey.get('CAPITAL')!, '50000.00')]);
    await post(`${y}-02-10`, [dr(c.byKey.get('CASH')!, '11500.00'), cr(c.byKey.get('SALES')!, '10000.00'), cr(c.byKey.get('VAT_OUTPUT')!, '1500.00')]);
    await post(`${y}-03-10`, [dr(c.byKey.get('COGS')!, '4000.00'), cr(c.byKey.get('INVENTORY')!, '4000.00')]);
    await post(`${y}-04-10`, [dr(c.byCode.get('6100')!, '2500.00'), cr(c.byKey.get('BANK')!, '2500.00')]);
    // Net profit = 10000 - 4000 - 2500 = 3500

    const draft = (await api.post('/api/journal-entries', { entryDate: `${y}-05-01`, description: 'مسودة', lines: [dr(c.byCode.get('6200')!, '1.00'), cr(c.byKey.get('CASH')!, '1.00')] })).json();
    const blocked = await api.post(`/api/fiscal-years/${year!.id}/close`);
    expect(blocked.json().error.code).toBe('DRAFTS_EXIST');
    await api.del(`/api/journal-entries/${draft.id}`);

    const res = await api.post(`/api/fiscal-years/${year!.id}/close`);
    expect(res.statusCode).toBe(200);
    const closing = (await api.get(`/api/journal-entries/${res.json().closingEntryId}`)).json();
    expect(closing).toMatchObject({ status: 'POSTED', referenceType: 'YEAR_CLOSING', source: 'SYSTEM', entryDate: year!.endDate });
    const re = closing.lines.find((l: { accountCode: string }) => l.accountCode === '3200');
    expect(re).toMatchObject({ debit: '0.00', credit: '3500.00' });

    for (const k of ['SALES', 'COGS']) {
      expect((await api.get(`/api/accounts/${c.byKey.get(k)}`)).json().balance).toBe('0.00');
    }
    expect((await api.get(`/api/accounts/${c.byCode.get('6100')}`)).json().balance).toBe('0.00');
    expect((await api.get(`/api/accounts/${c.byKey.get('RETAINED_EARNINGS')}`)).json().balance).toBe('-3500.00');
    // Balance-sheet accounts carry forward untouched.
    expect((await api.get(`/api/accounts/${c.byKey.get('CASH')}`)).json().balance).toBe('61500.00');

    const [closed] = await fiscalYears(t.app, s.token);
    expect(closed!.status).toBe('CLOSED');
    expect(closed!.periods.every((p) => p.status === 'CLOSED')).toBe(true);

    const late = await api.post('/api/journal-entries', { entryDate: `${y}-06-01`, description: 'متأخر', post: true, lines: [dr(c.byCode.get('6200')!, '1.00'), cr(c.byKey.get('CASH')!, '1.00')] });
    expect(late.json().error.code).toBe('PERIOD_CLOSED');
    expect((await api.post(`/api/fiscal-periods/${closed!.periods[0]!.id}/reopen`)).statusCode).toBe(409);
    expect((await api.post(`/api/fiscal-years/${year!.id}/close`)).statusCode).toBe(409);
    const reverseClosing = await api.post(`/api/journal-entries/${closing.id}/reverse`, { reason: 'test' });
    expect(reverseClosing.json().error.code).toBe('SYSTEM_ENTRY');

    const tb = (await api.get(`/api/reports/trial-balance?dateFrom=${year!.startDate}&dateTo=${year!.endDate}`)).json();
    expect(tb.balanced).toBe(true);
  });

  it('debits retained earnings for a loss', async () => {
    const s = await register(t.app);
    const api = client(t.app, s.token);
    const c = await chart(t.app, s.token);
    const [year] = await fiscalYears(t.app, s.token);
    await api.post('/api/journal-entries', { entryDate: year!.periods[0]!.startDate, description: 'إيجار', post: true, lines: [dr(c.byCode.get('6100')!, '1200.00'), cr(c.byKey.get('CASH')!, '1200.00')] });
    const res = await api.post(`/api/fiscal-years/${year!.id}/close`);
    const closing = (await api.get(`/api/journal-entries/${res.json().closingEntryId}`)).json();
    expect(closing.lines.find((l: { accountCode: string }) => l.accountCode === '3200')).toMatchObject({ debit: '1200.00', credit: '0.00' });
  });

  it('closes a year with no activity without a closing entry', async () => {
    const s = await register(t.app);
    const [year] = await fiscalYears(t.app, s.token);
    const res = await client(t.app, s.token).post(`/api/fiscal-years/${year!.id}/close`);
    expect(res.json()).toEqual({ closingEntryId: null });
  });

  it('requires earlier years to be closed first', async () => {
    const s = await register(t.app);
    const api = client(t.app, s.token);
    const [current] = await fiscalYears(t.app, s.token);
    await api.post('/api/fiscal-years', { startDate: `${Number(current!.startDate.slice(0, 4)) - 1}-01-01` });
    const res = await api.post(`/api/fiscal-years/${current!.id}/close`);
    expect(res.json().error.code).toBe('EARLIER_YEAR_OPEN');
  });
});
