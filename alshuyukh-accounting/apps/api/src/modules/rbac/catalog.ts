/**
 * Permission and system-role catalog.
 *
 * This file is the single source of truth. `npm run db:migrate` syncs it into
 * the `permissions`, `roles` (system rows) and `role_permissions` tables.
 * Tenants can create additional custom roles from these permissions.
 *
 * Permissions for modules delivered in later phases are defined now so that
 * role definitions stay stable; the endpoints that check them arrive with
 * their phase.
 */

export interface PermissionDef {
  code: string;
  module: string;
  ar: string;
  en: string;
}

export const PERMISSIONS: readonly PermissionDef[] = [
  // Phase 1 — platform foundations
  { code: 'user.view', module: 'users', ar: 'عرض المستخدمين', en: 'View users' },
  { code: 'user.invite', module: 'users', ar: 'إضافة مستخدمين', en: 'Add users' },
  { code: 'user.manage', module: 'users', ar: 'إدارة المستخدمين وأدوارهم', en: 'Manage users and their roles' },
  { code: 'role.view', module: 'roles', ar: 'عرض الأدوار', en: 'View roles' },
  { code: 'role.manage', module: 'roles', ar: 'إدارة الأدوار والصلاحيات', en: 'Manage roles and permissions' },
  { code: 'company.view', module: 'companies', ar: 'عرض الشركات والفروع', en: 'View companies and branches' },
  { code: 'company.manage', module: 'companies', ar: 'إدارة الشركات والفروع والمستودعات', en: 'Manage companies, branches, and warehouses' },
  { code: 'settings.manage', module: 'settings', ar: 'إدارة الإعدادات', en: 'Manage settings' },
  { code: 'audit.view', module: 'audit', ar: 'عرض سجل التدقيق', en: 'View audit log' },
  { code: 'subscription.manage', module: 'subscriptions', ar: 'إدارة الاشتراك', en: 'Manage subscription' },

  // Phase 2 — accounting
  { code: 'account.view', module: 'accounting', ar: 'عرض دليل الحسابات', en: 'View chart of accounts' },
  { code: 'account.manage', module: 'accounting', ar: 'إدارة دليل الحسابات', en: 'Manage chart of accounts' },
  { code: 'journal.view', module: 'accounting', ar: 'عرض القيود', en: 'View journal entries' },
  { code: 'journal.create', module: 'accounting', ar: 'إنشاء القيود', en: 'Create journal entries' },
  { code: 'journal.post', module: 'accounting', ar: 'ترحيل القيود', en: 'Post journal entries' },
  { code: 'journal.reverse', module: 'accounting', ar: 'عكس القيود', en: 'Reverse journal entries' },
  { code: 'fiscal.manage', module: 'accounting', ar: 'إدارة السنوات والفترات المالية', en: 'Manage fiscal years and periods' },

  // Phase 3–4 — sales, purchases, payments
  { code: 'customer.view', module: 'sales', ar: 'عرض العملاء', en: 'View customers' },
  { code: 'customer.manage', module: 'sales', ar: 'إدارة العملاء', en: 'Manage customers' },
  { code: 'supplier.view', module: 'purchases', ar: 'عرض الموردين', en: 'View suppliers' },
  { code: 'supplier.manage', module: 'purchases', ar: 'إدارة الموردين', en: 'Manage suppliers' },
  { code: 'product.view', module: 'inventory', ar: 'عرض المنتجات', en: 'View products' },
  { code: 'product.manage', module: 'inventory', ar: 'إدارة المنتجات', en: 'Manage products' },
  { code: 'invoice.view', module: 'sales', ar: 'عرض الفواتير', en: 'View invoices' },
  { code: 'invoice.create', module: 'sales', ar: 'إنشاء الفواتير', en: 'Create invoices' },
  { code: 'invoice.edit', module: 'sales', ar: 'تعديل الفواتير (مسودة)', en: 'Edit draft invoices' },
  { code: 'invoice.delete', module: 'sales', ar: 'حذف الفواتير (مسودة)', en: 'Delete draft invoices' },
  { code: 'invoice.post', module: 'sales', ar: 'اعتماد الفواتير', en: 'Issue invoices' },
  { code: 'invoice.cancel', module: 'sales', ar: 'إلغاء الفواتير', en: 'Cancel invoices' },
  { code: 'purchase.view', module: 'purchases', ar: 'عرض المشتريات', en: 'View purchases' },
  { code: 'purchase.create', module: 'purchases', ar: 'إنشاء المشتريات', en: 'Create purchases' },
  { code: 'purchase.post', module: 'purchases', ar: 'اعتماد المشتريات', en: 'Post purchases' },
  { code: 'payment.view', module: 'payments', ar: 'عرض المدفوعات', en: 'View payments' },
  { code: 'payment.create', module: 'payments', ar: 'تسجيل المدفوعات', en: 'Record payments' },

  // Phase 5–6 — inventory, expenses, tax
  { code: 'inventory.view', module: 'inventory', ar: 'عرض المخزون', en: 'View inventory' },
  { code: 'inventory.adjust', module: 'inventory', ar: 'تسوية المخزون', en: 'Adjust inventory' },
  { code: 'inventory.transfer', module: 'inventory', ar: 'تحويل المخزون', en: 'Transfer inventory' },
  { code: 'expense.view', module: 'expenses', ar: 'عرض المصروفات', en: 'View expenses' },
  { code: 'expense.create', module: 'expenses', ar: 'إنشاء المصروفات', en: 'Create expenses' },
  { code: 'expense.post', module: 'expenses', ar: 'اعتماد المصروفات', en: 'Post expenses' },
  { code: 'tax.manage', module: 'tax', ar: 'إدارة الضرائب', en: 'Manage tax rates' },

  // Phase 7–8 — reports, e-invoicing
  { code: 'report.view', module: 'reports', ar: 'عرض التقارير', en: 'View reports' },
  { code: 'zatca.view', module: 'zatca', ar: 'عرض الفوترة الإلكترونية', en: 'View e-invoicing' },
  { code: 'zatca.manage', module: 'zatca', ar: 'إدارة الفوترة الإلكترونية', en: 'Manage e-invoicing' },
] as const;

export const ALL_PERMISSION_CODES = PERMISSIONS.map((p) => p.code);

export interface SystemRoleDef {
  code: string;
  ar: string;
  en: string;
  permissions: readonly string[] | 'ALL';
}

const ADMIN_EXCLUDED = new Set(['subscription.manage']);

/**
 * SUPER_ADMIN is not a tenant role. Platform administration uses
 * users.is_platform_admin and a separate /admin surface (Phase 9).
 */
export const SYSTEM_ROLES: readonly SystemRoleDef[] = [
  { code: 'TENANT_OWNER', ar: 'مالك المنشأة', en: 'Tenant owner', permissions: 'ALL' },
  {
    code: 'COMPANY_ADMIN',
    ar: 'مدير الشركة',
    en: 'Company admin',
    permissions: ALL_PERMISSION_CODES.filter((c) => !ADMIN_EXCLUDED.has(c)),
  },
  {
    code: 'ACCOUNTANT',
    ar: 'محاسب',
    en: 'Accountant',
    permissions: [
      'company.view', 'account.view', 'account.manage', 'journal.view', 'journal.create',
      'journal.post', 'journal.reverse', 'fiscal.manage', 'customer.view', 'supplier.view',
      'product.view', 'invoice.view', 'invoice.create', 'invoice.edit', 'invoice.post',
      'invoice.cancel', 'purchase.view', 'purchase.create', 'purchase.post', 'payment.view',
      'payment.create', 'inventory.view', 'expense.view', 'expense.create', 'expense.post',
      'tax.manage', 'report.view', 'zatca.view',
    ],
  },
  {
    code: 'SALES_MANAGER',
    ar: 'مدير المبيعات',
    en: 'Sales manager',
    permissions: [
      'company.view', 'customer.view', 'customer.manage', 'product.view', 'invoice.view',
      'invoice.create', 'invoice.edit', 'invoice.delete', 'invoice.post', 'invoice.cancel',
      'payment.view', 'payment.create', 'inventory.view', 'report.view', 'zatca.view',
    ],
  },
  {
    code: 'SALES_EMPLOYEE',
    ar: 'موظف مبيعات',
    en: 'Sales employee',
    permissions: ['company.view', 'customer.view', 'product.view', 'invoice.view', 'invoice.create', 'invoice.edit'],
  },
  {
    code: 'PURCHASE_MANAGER',
    ar: 'مدير المشتريات',
    en: 'Purchase manager',
    permissions: [
      'company.view', 'supplier.view', 'supplier.manage', 'product.view', 'purchase.view',
      'purchase.create', 'purchase.post', 'payment.view', 'inventory.view', 'report.view',
    ],
  },
  {
    code: 'WAREHOUSE_MANAGER',
    ar: 'مدير المستودع',
    en: 'Warehouse manager',
    permissions: [
      'company.view', 'product.view', 'product.manage', 'inventory.view', 'inventory.adjust',
      'inventory.transfer', 'report.view',
    ],
  },
  {
    code: 'WAREHOUSE_EMPLOYEE',
    ar: 'موظف مستودع',
    en: 'Warehouse employee',
    permissions: ['company.view', 'product.view', 'inventory.view', 'inventory.transfer'],
  },
  {
    code: 'VIEWER',
    ar: 'مشاهد',
    en: 'Viewer',
    permissions: ALL_PERMISSION_CODES.filter((c) => c.endsWith('.view')),
  },
];

export function resolveRolePermissions(role: SystemRoleDef): string[] {
  return role.permissions === 'ALL' ? [...ALL_PERMISSION_CODES] : [...role.permissions];
}
