/**
 * Default chart of accounts created for every new company.
 * Users can rename, renumber, extend, or deactivate accounts. The engine
 * finds accounts by systemKey, never by code.
 */

export type AccountType = 'ASSET' | 'LIABILITY' | 'EQUITY' | 'REVENUE' | 'EXPENSE' | 'COST_OF_GOODS_SOLD';

export const ACCOUNT_TYPES: readonly AccountType[] = ['ASSET', 'LIABILITY', 'EQUITY', 'REVENUE', 'EXPENSE', 'COST_OF_GOODS_SOLD'];

/** Types whose balances close into retained earnings at year end. */
export const PROFIT_AND_LOSS_TYPES: readonly AccountType[] = ['REVENUE', 'EXPENSE', 'COST_OF_GOODS_SOLD'];

/** Debit-normal types; the rest are credit-normal. */
export const DEBIT_NORMAL: ReadonlySet<AccountType> = new Set(['ASSET', 'EXPENSE', 'COST_OF_GOODS_SOLD']);

export interface GroupTemplate { code: string; ar: string; en: string; type: AccountType; sort: number }

export const DEFAULT_GROUPS: readonly GroupTemplate[] = [
  { code: 'CURRENT_ASSETS', ar: 'الأصول المتداولة', en: 'Current assets', type: 'ASSET', sort: 10 },
  { code: 'NON_CURRENT_ASSETS', ar: 'الأصول غير المتداولة', en: 'Non-current assets', type: 'ASSET', sort: 20 },
  { code: 'CURRENT_LIABILITIES', ar: 'الالتزامات المتداولة', en: 'Current liabilities', type: 'LIABILITY', sort: 30 },
  { code: 'NON_CURRENT_LIABILITIES', ar: 'الالتزامات غير المتداولة', en: 'Non-current liabilities', type: 'LIABILITY', sort: 40 },
  { code: 'EQUITY', ar: 'حقوق الملكية', en: 'Equity', type: 'EQUITY', sort: 50 },
  { code: 'OPERATING_REVENUE', ar: 'الإيرادات التشغيلية', en: 'Operating revenue', type: 'REVENUE', sort: 60 },
  { code: 'OTHER_REVENUE', ar: 'إيرادات أخرى', en: 'Other revenue', type: 'REVENUE', sort: 70 },
  { code: 'COST_OF_SALES', ar: 'تكلفة المبيعات', en: 'Cost of sales', type: 'COST_OF_GOODS_SOLD', sort: 80 },
  { code: 'OPERATING_EXPENSES', ar: 'المصروفات التشغيلية', en: 'Operating expenses', type: 'EXPENSE', sort: 90 },
  { code: 'OTHER_EXPENSES', ar: 'مصروفات أخرى', en: 'Other expenses', type: 'EXPENSE', sort: 100 },
];

export interface AccountTemplate {
  code: string;
  ar: string;
  en: string;
  type: AccountType;
  parent?: string;
  group?: string;
  postable: boolean;
  systemKey?: string;
}

export const DEFAULT_ACCOUNTS: readonly AccountTemplate[] = [
  { code: '1000', ar: 'الأصول', en: 'Assets', type: 'ASSET', postable: false },
  { code: '1100', ar: 'النقدية', en: 'Cash', type: 'ASSET', parent: '1000', group: 'CURRENT_ASSETS', postable: true, systemKey: 'CASH' },
  { code: '1200', ar: 'البنك', en: 'Bank', type: 'ASSET', parent: '1000', group: 'CURRENT_ASSETS', postable: true, systemKey: 'BANK' },
  { code: '1300', ar: 'العملاء', en: 'Accounts receivable', type: 'ASSET', parent: '1000', group: 'CURRENT_ASSETS', postable: true, systemKey: 'ACCOUNTS_RECEIVABLE' },
  { code: '1400', ar: 'المخزون', en: 'Inventory', type: 'ASSET', parent: '1000', group: 'CURRENT_ASSETS', postable: true, systemKey: 'INVENTORY' },

  { code: '2000', ar: 'الالتزامات', en: 'Liabilities', type: 'LIABILITY', postable: false },
  { code: '2100', ar: 'الموردون', en: 'Accounts payable', type: 'LIABILITY', parent: '2000', group: 'CURRENT_LIABILITIES', postable: true, systemKey: 'ACCOUNTS_PAYABLE' },
  { code: '2200', ar: 'ضريبة القيمة المضافة', en: 'Value added tax', type: 'LIABILITY', parent: '2000', postable: false },
  { code: '2210', ar: 'ضريبة القيمة المضافة - المخرجات', en: 'VAT output', type: 'LIABILITY', parent: '2200', group: 'CURRENT_LIABILITIES', postable: true, systemKey: 'VAT_OUTPUT' },
  { code: '2220', ar: 'ضريبة القيمة المضافة - المدخلات', en: 'VAT input', type: 'LIABILITY', parent: '2200', group: 'CURRENT_LIABILITIES', postable: true, systemKey: 'VAT_INPUT' },

  { code: '3000', ar: 'حقوق الملكية', en: 'Equity', type: 'EQUITY', postable: false },
  { code: '3100', ar: 'رأس المال', en: 'Capital', type: 'EQUITY', parent: '3000', group: 'EQUITY', postable: true, systemKey: 'CAPITAL' },
  { code: '3200', ar: 'الأرباح المحتجزة', en: 'Retained earnings', type: 'EQUITY', parent: '3000', group: 'EQUITY', postable: true, systemKey: 'RETAINED_EARNINGS' },

  { code: '4000', ar: 'الإيرادات', en: 'Revenue', type: 'REVENUE', postable: false },
  { code: '4100', ar: 'المبيعات', en: 'Sales', type: 'REVENUE', parent: '4000', group: 'OPERATING_REVENUE', postable: true, systemKey: 'SALES' },

  { code: '5000', ar: 'تكلفة المبيعات', en: 'Cost of goods sold', type: 'COST_OF_GOODS_SOLD', postable: false },
  { code: '5100', ar: 'تكلفة البضاعة المباعة', en: 'Cost of goods sold', type: 'COST_OF_GOODS_SOLD', parent: '5000', group: 'COST_OF_SALES', postable: true, systemKey: 'COGS' },
  { code: '5200', ar: 'فروقات وتسويات المخزون', en: 'Inventory adjustments', type: 'COST_OF_GOODS_SOLD', parent: '5000', group: 'COST_OF_SALES', postable: true, systemKey: 'INVENTORY_ADJUSTMENT' },

  { code: '6000', ar: 'المصروفات', en: 'Expenses', type: 'EXPENSE', postable: false },
  { code: '6100', ar: 'الإيجار', en: 'Rent', type: 'EXPENSE', parent: '6000', group: 'OPERATING_EXPENSES', postable: true },
  { code: '6200', ar: 'الرواتب', en: 'Salaries', type: 'EXPENSE', parent: '6000', group: 'OPERATING_EXPENSES', postable: true },
  { code: '6300', ar: 'التسويق', en: 'Marketing', type: 'EXPENSE', parent: '6000', group: 'OPERATING_EXPENSES', postable: true },
  { code: '6400', ar: 'الكهرباء', en: 'Electricity', type: 'EXPENSE', parent: '6000', group: 'OPERATING_EXPENSES', postable: true },
];
