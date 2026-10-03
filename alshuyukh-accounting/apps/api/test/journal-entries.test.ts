import { afterAll, beforeAll, describe, expect, it } from 'vitest';
import { withTx } from '../src/db/tx.js';
import { postEntry } from '../src/modules/accounting/engine.js';
import {
  addMember, chart, client, cr, dateInYear, dr, fiscalYears, register, setupApp,
  type Chart, type Session, type TestContext,
} from './helpers.js';

let t: TestContext;
let owner: Session;
let c: Chart;
let date: string;

beforeAll(async () => {
  t = await setupApp();
  owner = await register(t.app);
  c = await chart(t.app, owner.token);
  date = await dateInYear(t.app, owner.token, 3, 15);
});
afterAll(async () => { await t.close(); });

const api = () => client(t.app, owner.token);
const key = (k: string) => c.byKey.get(k)!;
const code = (k: string) => c.byCode.get(k)!;
const create = (lines: unknown[], extra: Record<string, unknown> = {}) =>
  api().post('/api/journal-entries', { entryDate: date, description: 'قيد اختبار', lines, ...extra });

describe('the required example: invoice 1000 SAR + VAT 150 SAR', () => {
  it('posts AR 1150 debit, Sales 1000 credit, VAT 150 credit through the engine', async () => {
    const entry = await withTx(t.pool, { tenantId: owner.tenantId, userId: owner.userId }, async (db) => {
      const { rows: [company] } = await db.query<{ id: string }>('SELECT id FROM companies WHERE tenant_id = $1', [owner.tenantId]);
      return postEntry(db, { tenantId: owner.tenantId, userId: owner.userId }, {
        companyId: company!.id, entryDate: date, description: 'فاتورة مبيعات INV-0001',
        referenceType: 'SALES_INVOICE', source: 'SYSTEM',
        lines: [
          dr(key('ACCOUNTS_RECEIVABLE'), '1150.00'),
          cr(key('SALES'), '1000.00'),
          cr(key('VAT_OUTPUT'), '150.00'),
        ],
      });
    });

    expect(entry.status).toBe('POSTED');
    expect(entry.totalDebit).toBe('1150.00');
    expect(entry.totalCredit).toBe('1150.00');
    expect(entry.lines.map((l) => [l.accountCode, l.debit, l.credit])).toEqual([
      ['1300', '1150.00', '0.00'],
      ['4100', '0.00', '1000.00'],
      ['2210', '0.00', '150.00'],
    ]);

    // Account balances come only from posted journal lines.
    const ar = (await api().get(`/api/accounts/${key('ACCOUNTS_RECEIVABLE')}`)).json();
    const sales = (await api().get(`/api/accounts/${key('SALES')}`)).json();
    const vat = (await api().get(`/api/accounts/${key('VAT_OUTPUT')}`)).json();
    expect(ar.balance).toBe('1150.00');
    expect(sales.balance).toBe('-1000.00');
    expect(vat.balance).toBe('-150.00');
  });

  it('refuses to reverse a document-generated entry manually', async () => {
    const list = (await api().get('/api/journal-entries?referenceType=SALES_INVOICE')).json().data;
    const res = await api().post(`/api/journal-entries/${list[0].id}/reverse`, { reason: 'خطأ' });
    expect(res.statusCode).toBe(409);
    expect(res.json().error.code).toBe('SYSTEM_ENTRY');
  });
});

describe('posting rules', () => {
  it('posts a balanced entry and numbers it sequentially', async () => {
    const [year] = await fiscalYears(t.app, owner.token);
    const a = (await create([dr(code('6100'), '2500.00'), cr(key('BANK'), '2500.00')], { post: true })).json();
    const b = (await create([dr(code('6200'), '8000.00'), cr(key('BANK'), '8000.00')], { post: true })).json();
    const prefix = `JV-${year!.startDate.slice(0, 4)}-`;
    expect(a.entryNumber).toMatch(new RegExp(`^${prefix}\\d{6}$`));
    expect(Number(b.entryNumber.slice(-6))).toBe(Number(a.entryNumber.slice(-6)) + 1);
    expect(a).toMatchObject({ status: 'POSTED', totalDebit: '2500.00', totalCredit: '2500.00' });
  });

  it('rejects an unbalanced entry at posting and keeps it as a draft', async () => {
    const draft = (await create([dr(code('6300'), '100.00'), cr(key('CASH'), '90.00')])).json();
    expect(draft).toMatchObject({ status: 'DRAFT', entryNumber: null, totalDebit: '100.00', totalCredit: '90.00' });
    const res = await api().post(`/api/journal-entries/${draft.id}/post`);
    expect(res.statusCode).toBe(400);
    expect(res.json().error).toMatchObject({ code: 'UNBALANCED_ENTRY', details: { totalDebit: '100.00', totalCredit: '90.00', difference: '10.00' } });
    expect((await api().get(`/api/journal-entries/${draft.id}`)).json().status).toBe('DRAFT');
  });

  it('rolls back completely when create-and-post fails', async () => {
    const before = (await api().get('/api/journal-entries?limit=200')).json().total;
    const res = await create([dr(code('6300'), '100.00'), cr(key('CASH'), '99.99')], { post: true, description: 'should not exist' });
    expect(res.statusCode).toBe(400);
    expect((await api().get('/api/journal-entries?limit=200')).json().total).toBe(before);
  });

  it('uses exact decimal arithmetic (0.10 + 0.20 = 0.30)', async () => {
    const res = await create([dr(code('6300'), '0.10'), dr(code('6400'), '0.20'), cr(key('CASH'), '0.30')], { post: true });
    expect(res.statusCode).toBe(201);
    expect(res.json().totalDebit).toBe('0.30');
  });

  it('validates every line', async () => {
    const cases = [
      [{ accountId: code('6300'), debit: '10.00', credit: '10.00' }, cr(key('CASH'), '10.00')], // both sides
      [{ accountId: code('6300'), debit: '0' }, cr(key('CASH'), '0')],                          // zero
      [dr(code('6300'), '-5.00'), cr(key('CASH'), '-5.00')],                                     // negative
      [dr(code('6300'), '1.005'), cr(key('CASH'), '1.005')],                                     // 3 decimals
      [{ accountId: code('6300'), debit: 100 }, { accountId: key('CASH'), credit: 100 }],       // number, not string
      [dr(code('6300'), '10.00')],                                                               // single line
    ];
    for (const lines of cases) {
      const res = await create(lines);
      expect(res.statusCode, JSON.stringify(lines)).toBe(400);
    }
  });

  it('rejects header, inactive and foreign accounts', async () => {
    const header = await create([dr(code('6000'), '10.00'), cr(key('CASH'), '10.00')]);
    expect(header.json().error.code).toBe('ACCOUNT_NOT_POSTABLE');

    const temp = (await api().post('/api/accounts', { code: '6950', nameAr: 'غير نشط', parentId: code('6000') })).json();
    await api().patch(`/api/accounts/${temp.id}`, { isActive: false });
    expect((await create([dr(temp.id, '10.00'), cr(key('CASH'), '10.00')])).json().error.code).toBe('ACCOUNT_INACTIVE');

    const other = await register(t.app);
    const otherChart = await chart(t.app, other.token);
    const foreign = await create([dr(otherChart.byKey.get('CASH')!, '10.00'), cr(key('CASH'), '10.00')]);
    expect(foreign.json().error.code).toBe('INVALID_ACCOUNT');
  });

  it('requires an open fiscal period that covers the date', async () => {
    const res = await create([dr(code('6300'), '10.00'), cr(key('CASH'), '10.00')], { entryDate: '2010-01-01', post: true });
    expect(res.statusCode).toBe(400);
    expect(res.json().error.code).toBe('NO_FISCAL_PERIOD');
  });

  it('assigns unique gap-free numbers under concurrent posting', async () => {
    const drafts = await Promise.all(Array.from({ length: 8 }, () =>
      create([dr(code('6400'), '1.00'), cr(key('CASH'), '1.00')]).then((r) => r.json().id as string)));
    const posted = await Promise.all(drafts.map((id) => api().post(`/api/journal-entries/${id}/post`).then((r) => r.json())));
    const numbers = posted.map((e) => Number(e.entryNumber.slice(-6))).sort((x, y) => x - y);
    expect(new Set(numbers).size).toBe(8);
    expect(numbers[7]! - numbers[0]!).toBe(7);
  });
});

describe('drafts', () => {
  it('edits and soft-deletes a draft', async () => {
    const draft = (await create([dr(code('6300'), '10.00'), cr(key('CASH'), '10.00')])).json();
    const edited = await api().patch(`/api/journal-entries/${draft.id}`, {
      description: 'معدل', lines: [dr(code('6300'), '25.00'), cr(key('BANK'), '25.00')],
    });
    expect(edited.json()).toMatchObject({ description: 'معدل', totalDebit: '25.00' });
    expect(edited.json().lines[1].accountCode).toBe('1200');
    expect((await api().del(`/api/journal-entries/${draft.id}`)).statusCode).toBe(204);
    expect((await api().get(`/api/journal-entries/${draft.id}`)).statusCode).toBe(404);
  });
});

describe('immutability of posted entries', () => {
  let entry: { id: string; lines: { id: string }[] };

  beforeAll(async () => {
    entry = (await create([dr(code('6100'), '300.00'), cr(key('CASH'), '300.00')], { post: true })).json();
  });

  it('refuses edits and deletion through the API', async () => {
    expect((await api().patch(`/api/journal-entries/${entry.id}`, { description: 'تعديل' })).statusCode).toBe(409);
    expect((await api().del(`/api/journal-entries/${entry.id}`)).statusCode).toBe(409);
    expect((await api().post(`/api/journal-entries/${entry.id}/post`)).statusCode).toBe(409);
  });

  it('refuses direct SQL changes from the application role', async () => {
    const ctx = { tenantId: owner.tenantId, userId: owner.userId };
    const attempts = [
      `UPDATE journal_entries SET description = 'tampered' WHERE id = '${entry.id}'`,
      `UPDATE journal_entries SET status = 'DRAFT' WHERE id = '${entry.id}'`,
      `UPDATE journal_entry_lines SET debit = 1 WHERE id = '${entry.lines[0]!.id}'`,
      `DELETE FROM journal_entry_lines WHERE id = '${entry.lines[0]!.id}'`,
      `INSERT INTO journal_entry_lines (tenant_id, company_id, journal_entry_id, line_no, account_id, debit)
         SELECT tenant_id, company_id, id, 99, '${code('6100')}', 5 FROM journal_entries WHERE id = '${entry.id}'`,
      `DELETE FROM journal_entries WHERE id = '${entry.id}'`,
    ];
    for (const sql of attempts) {
      await expect(withTx(t.pool, ctx, (db) => db.query(sql)), sql).rejects.toThrow();
    }
    const after = (await api().get(`/api/journal-entries/${entry.id}`)).json();
    expect(after).toMatchObject({ status: 'POSTED', totalDebit: '300.00' });
  });

  it('enforces the balance rule in the database even if the API is bypassed', async () => {
    const ctx = { tenantId: owner.tenantId, userId: owner.userId };
    const draft = (await create([dr(code('6100'), '50.00'), cr(key('CASH'), '40.00')])).json();
    const [year] = await fiscalYears(t.app, owner.token);
    const period = year!.periods[2]!;
    await expect(withTx(t.pool, ctx, (db) => db.query(
      `UPDATE journal_entries SET status = 'POSTED', entry_number = 'FAKE-1', fiscal_year_id = $2, fiscal_period_id = $3, posted_at = now() WHERE id = $1`,
      [draft.id, year!.id, period.id]))).rejects.toThrow(/not balanced/);
  });
});

describe('reversal and correction', () => {
  it('reverses a posted entry with a mirror entry and links both', async () => {
    const original = (await create([dr(code('6300'), '700.00'), cr(key('BANK'), '700.00')], { post: true, description: 'حملة إعلانية' })).json();
    const res = await api().post(`/api/journal-entries/${original.id}/reverse`, { reason: 'تسجيل في حساب خاطئ' });
    expect(res.statusCode).toBe(200);
    const { original: o, reversal } = res.json();
    expect(o).toMatchObject({ status: 'REVERSED', reversedByEntryId: reversal.id });
    expect(reversal).toMatchObject({ status: 'POSTED', referenceType: 'REVERSAL', referenceId: original.id, reversalOfId: original.id, totalDebit: '700.00' });
    expect(reversal.lines.map((l: { accountCode: string; debit: string; credit: string }) => [l.accountCode, l.debit, l.credit]))
      .toEqual([['6300', '0.00', '700.00'], ['1200', '700.00', '0.00']]);

    const again = await api().post(`/api/journal-entries/${original.id}/reverse`, { reason: 'مرة ثانية' });
    expect(again.json().error.code).toBe('ALREADY_REVERSED');
    const ofReversal = await api().post(`/api/journal-entries/${reversal.id}/reverse`, { reason: 'عكس العكس' });
    expect(ofReversal.json().error.code).toBe('CANNOT_REVERSE_REVERSAL');
  });

  it('creates the correcting entry as a linked draft in the same step', async () => {
    const original = (await create([dr(code('6300'), '400.00'), cr(key('CASH'), '400.00')], { post: true })).json();
    const res = await api().post(`/api/journal-entries/${original.id}/reverse`, {
      reason: 'الحساب الصحيح هو الكهرباء',
      correction: { entryDate: date, description: 'قيد تصحيح', lines: [dr(code('6400'), '400.00'), cr(key('CASH'), '400.00')] },
    });
    expect(res.statusCode).toBe(200);
    const correction = res.json().correction;
    expect(correction).toMatchObject({ status: 'DRAFT', correctionOfId: original.id });
    expect((await api().post(`/api/journal-entries/${correction.id}/post`)).json().status).toBe('POSTED');
  });

  it('cannot reverse a draft', async () => {
    const draft = (await create([dr(code('6300'), '10.00'), cr(key('CASH'), '10.00')])).json();
    expect((await api().post(`/api/journal-entries/${draft.id}/reverse`, { reason: 'test' })).json().error.code).toBe('ENTRY_NOT_POSTED');
  });
});

describe('fiscal periods', () => {
  it('blocks posting into a closed period and allows it after reopening', async () => {
    const session = await register(t.app);
    const sc = await chart(t.app, session.token);
    const sApi = client(t.app, session.token);
    const [year] = await fiscalYears(t.app, session.token);
    const period = year!.periods[4]!;
    const d = period.startDate;
    const draft = (await sApi.post('/api/journal-entries', { entryDate: d, description: 'x', lines: [dr(sc.byCode.get('6100')!, '10.00'), cr(sc.byKey.get('CASH')!, '10.00')] })).json();

    expect((await sApi.post(`/api/fiscal-periods/${period.id}/close`)).statusCode).toBe(200);
    const blocked = await sApi.post(`/api/journal-entries/${draft.id}/post`);
    expect(blocked.statusCode).toBe(409);
    expect(blocked.json().error.code).toBe('PERIOD_CLOSED');

    expect((await sApi.post(`/api/fiscal-periods/${period.id}/reopen`)).statusCode).toBe(200);
    expect((await sApi.post(`/api/journal-entries/${draft.id}/post`)).json().status).toBe('POSTED');
  });
});

describe('permissions and tenant isolation', () => {
  it('applies journal permissions per role', async () => {
    const lines = [dr(code('6300'), '10.00'), cr(key('CASH'), '10.00')];
    const sales = await addMember(t.app, owner, 'SALES_EMPLOYEE');
    expect((await client(t.app, sales.token).post('/api/journal-entries', { entryDate: date, description: 'x', lines })).statusCode).toBe(403);
    expect((await client(t.app, sales.token).get('/api/journal-entries')).statusCode).toBe(403);

    const viewer = await addMember(t.app, owner, 'VIEWER');
    expect((await client(t.app, viewer.token).get('/api/journal-entries')).statusCode).toBe(200);
    expect((await client(t.app, viewer.token).post('/api/journal-entries', { entryDate: date, description: 'x', lines })).statusCode).toBe(403);

    const accountant = await addMember(t.app, owner, 'ACCOUNTANT');
    const res = await client(t.app, accountant.token).post('/api/journal-entries', { entryDate: date, description: 'x', lines, post: true });
    expect(res.statusCode).toBe(201);

    const role = await api().post('/api/roles', { code: 'JOURNAL_CLERK', nameAr: 'مدخل قيود', nameEn: 'Clerk', permissions: ['journal.view', 'journal.create', 'account.view'] });
    const clerk = await addMember(t.app, owner, 'VIEWER');
    await api().put(`/api/users/${clerk.userId}/roles`, { roleIds: [role.json().id] });
    const cApi = client(t.app, clerk.token);
    expect((await cApi.post('/api/journal-entries', { entryDate: date, description: 'x', lines, post: true })).statusCode).toBe(403);
    const draft = (await cApi.post('/api/journal-entries', { entryDate: date, description: 'x', lines })).json();
    expect((await cApi.post(`/api/journal-entries/${draft.id}/post`)).statusCode).toBe(403);
  });

  it('hides entries from other tenants', async () => {
    const entry = (await create([dr(code('6300'), '10.00'), cr(key('CASH'), '10.00')], { post: true })).json();
    const other = await register(t.app);
    const oApi = client(t.app, other.token);
    expect((await oApi.get(`/api/journal-entries/${entry.id}`)).statusCode).toBe(404);
    expect((await oApi.post(`/api/journal-entries/${entry.id}/reverse`, { reason: 'hack' })).statusCode).toBe(404);
    expect((await oApi.get('/api/journal-entries')).json().total).toBe(0);
  });

  it('records posting and reversal in the audit log', async () => {
    const actions = (await api().get('/api/audit-logs?entityType=journal_entry&limit=200')).json().data.map((l: { action: string }) => l.action);
    expect(actions).toEqual(expect.arrayContaining(['CREATE', 'POST', 'REVERSE', 'UPDATE', 'DELETE']));
  });
});

describe('trial balance', () => {
  it('balances and nets reversed entries to zero', async () => {
    const [year] = await fiscalYears(t.app, owner.token);
    const res = await api().get(`/api/reports/trial-balance?dateFrom=${year!.startDate}&dateTo=${year!.endDate}`);
    expect(res.statusCode).toBe(200);
    const tb = res.json();
    expect(tb.balanced).toBe(true);
    expect(tb.totals.periodDebit).toBe(tb.totals.periodCredit);
    expect(tb.totals.closingDebit).toBe(tb.totals.closingCredit);
    const marketing = tb.data.find((r: { code: string }) => r.code === '6300');
    // 700 and 400 were reversed, 100 was never posted; remaining: 0.10 + 10 (accountant) + 10 (isolation test).
    expect(marketing.closingDebit).toBe('20.10');
  });

  it('splits opening and period movements by date', async () => {
    const [year] = await fiscalYears(t.app, owner.token);
    const tb = (await api().get(`/api/reports/trial-balance?dateFrom=${date.slice(0, 8)}16&dateTo=${year!.endDate}`)).json();
    const ar = tb.data.find((r: { code: string }) => r.code === '1300');
    expect(ar).toMatchObject({ openingDebit: '1150.00', periodDebit: '0.00', closingDebit: '1150.00' });
    expect(tb.balanced).toBe(true);
  });
});
