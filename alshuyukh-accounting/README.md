# ALSHUYUKH ACCOUNTING — الشيوخ للمحاسبة

نظام محاسبي سحابي متعدد المنشآت (Multi-Tenant SaaS) للسوق السعودي.

**الحالة:**

- المرحلة 1 مكتملة: تأسيس المشروع، قاعدة البيانات، المصادقة، تعدد المنشآت، الصلاحيات.
- المرحلة 2 مكتملة: دليل الحسابات، المحرك المحاسبي، القيود اليومية، السنوات والفترات المالية.
- المرحلة 3 مكتملة: العملاء، الموردون، المنتجات والخدمات، الوحدات والتصنيفات.
- المرحلة 4 مكتملة: عروض الأسعار، فواتير المبيعات والمشتريات، المرتجعات، أوامر الشراء، المدفوعات، نسب الضريبة.
- المرحلة 5 مكتملة: حركات المخزون، المتوسط المرجح، تكلفة المبيعات، التحويلات، التسويات والجرد، تقييم المخزون.

بقية الوحدات تظهر في الواجهة بعلامة «قريبًا» دون أي بيانات تجريبية.

## التقنيات المستخدمة

| الطبقة | التقنية | سبب الاختيار |
|---|---|---|
| Backend/API | Node.js 22 + TypeScript + Fastify 5 | أداء عالٍ، أنواع صارمة، Plugins معيارية |
| قاعدة البيانات | PostgreSQL 16 | معاملات ACID، `NUMERIC` للمبالغ، Row-Level Security |
| الوصول للبيانات | `pg` مع SQL صريح ومعاملات (parameterized) | تحكم كامل في المعاملات والأقفال لمحرك المحاسبة |
| التحقق من المدخلات | Zod 4 | تحقق على كل endpoint |
| المصادقة | JWT (HS256، 15 دقيقة) + Refresh Token دوّار في Cookie | جلسات قابلة للإلغاء فورًا |
| كلمات المرور | argon2id (معايير OWASP) | |
| Frontend | React 19 + Vite + React Router، RTL أولًا | |
| الاختبارات | Vitest على قاعدة PostgreSQL حقيقية | |

## هيكل المشروع

```text
alshuyukh-accounting/
├── apps/
│   ├── api/
│   │   ├── src/
│   │   │   ├── app.ts                 # تجميع التطبيق: الأمان، الأخطاء، المسارات
│   │   │   ├── server.ts              # نقطة التشغيل
│   │   │   ├── config/env.ts          # التحقق من متغيرات البيئة
│   │   │   ├── db/
│   │   │   │   ├── migrations/*.sql   # مخطط قاعدة البيانات
│   │   │   │   ├── migrate.ts         # مشغّل الترحيلات + مزامنة الصلاحيات
│   │   │   │   ├── pool.ts            # اتصال UTC، NUMERIC كنص
│   │   │   │   └── tx.ts              # withTx: معاملة + سياق المنشأة
│   │   │   ├── lib/                   # الأخطاء، التحقق، كلمات المرور، الرموز، المبالغ، التواريخ
│   │   │   ├── plugins/auth.ts        # المصادقة وحارس الصلاحيات
│   │   │   └── modules/
│   │   │       ├── auth/              # التسجيل، الدخول، التجديد، الخروج، التبديل
│   │   │       ├── users/             # أعضاء المنشأة
│   │   │       ├── rbac/              # catalog.ts (مصدر الصلاحيات) + الأدوار
│   │   │       ├── accounting/        # المحرك المحاسبي (المرحلة 2)
│   │   │       │   ├── engine.ts      # إنشاء القيود، الترحيل، العكس، إقفال السنة
│   │   │       │   ├── chart-template.ts # دليل الحسابات الافتراضي
│   │   │       │   ├── setup.ts       # تجهيز الشركة: الدليل + السنة المالية
│   │   │       │   └── *.routes.ts    # الحسابات، القيود، السنوات، ميزان المراجعة
│   │   │       ├── documents/         # المستندات التجارية والمدفوعات (المرحلة 4)
│   │   │       │   ├── calc.ts        # حساب السطور والضريبة (دوال نقية)
│   │   │       │   ├── kinds.ts       # إعدادات المستندات الستة
│   │   │       │   ├── drafts.ts      # المسودات والمرتجعات
│   │   │       │   ├── posting.ts     # الإصدار والقيود والإلغاء والتحويل
│   │   │       │   └── payments.ts    # المقبوضات والمدفوعات والتخصيص
│   │   │       ├── inventory/         # المخزون (المرحلة 5)
│   │   │       │   ├── costing.ts     # طريقة التكلفة (المتوسط المرجح؛ واجهة لـ FIFO)
│   │   │       │   ├── engine.ts      # الوارد والصادر والإلغاء مع قفل الرصيد
│   │   │       │   └── routes.ts      # الأرصدة، حركة الصنف، التقييم، التحويلات، التسويات
│   │   │       ├── parties/           # العملاء والموردون (تنفيذ مشترك)
│   │   │       ├── products/          # المنتجات، الوحدات، التصنيفات
│   │   │       ├── companies/         # الشركات، الفروع، المستودعات
│   │   │       ├── settings/          # إعدادات المنشأة
│   │   │       └── audit/             # سجل التدقيق
│   │   └── test/                      # اختبارات تكامل
│   └── web/                           # واجهة React عربية
├── scripts/db-setup.sql               # إنشاء الأدوار وقواعد البيانات
├── docker-compose.yml                 # PostgreSQL للتطوير
└── .env.example
```

## التشغيل محليًا

المتطلبات: Node.js 22 أو أحدث، وPostgreSQL 16 (محليًا أو عبر Docker).

```sh
cd alshuyukh-accounting

# 1. قاعدة البيانات
docker compose up -d                       # أو: psql -U postgres -f scripts/db-setup.sql

# 2. الإعدادات
cp .env.example .env
# غيّر JWT_SECRET إلى قيمة عشوائية: openssl rand -base64 48

# 3. الحزم والترحيلات
npm install
npm run db:migrate

# 4. التشغيل (نافذتان)
npm run dev:api                            # http://localhost:3000
npm run dev:web                            # http://localhost:5173

# 5. الاختبارات (تعيد بناء قاعدة alshuyukh_test في كل تشغيل)
npm test
```

افتح `http://localhost:5173` واختر «أنشئ منشأة جديدة».

## مخطط قاعدة البيانات (ERD)

```mermaid
erDiagram
    tenants ||--|| tenant_settings : has
    tenants ||--o{ user_tenants : has
    users ||--o{ user_tenants : "member of"
    users ||--o{ user_sessions : has
    tenants ||--o{ user_sessions : "active in"
    user_tenants ||--o{ user_roles : "(tenant_id, user_id)"
    roles ||--o{ user_roles : assigned
    roles ||--o{ role_permissions : grants
    permissions ||--o{ role_permissions : in
    tenants ||--o{ roles : "custom roles"
    tenants ||--o{ companies : owns
    companies ||--o{ branches : "(company_id, tenant_id)"
    companies ||--o{ warehouses : "(company_id, tenant_id)"
    branches ||--o{ warehouses : "(branch_id, tenant_id)"
    tenants ||--o{ audit_logs : records
    users ||--o{ audit_logs : performed

    tenants { uuid id PK; text name; citext slug UK; text status; timestamptz deleted_at }
    tenant_settings { uuid tenant_id PK,FK; char default_currency; text timezone; text locale; smallint fiscal_year_start_month }
    users { uuid id PK; citext email UK; text password_hash; text status; bool is_platform_admin; int failed_login_attempts; timestamptz locked_until }
    user_tenants { uuid id PK; uuid tenant_id FK; uuid user_id FK; text status; bool is_owner }
    user_sessions { uuid id PK; uuid user_id FK; uuid tenant_id FK; text token_hash UK; timestamptz expires_at; timestamptz revoked_at; uuid replaced_by FK }
    roles { uuid id PK; uuid tenant_id FK "NULL = system"; text code; bool is_system; timestamptz deleted_at }
    permissions { uuid id PK; text code UK; text module }
    role_permissions { uuid role_id PK,FK; uuid permission_id PK,FK }
    user_roles { uuid id PK; uuid tenant_id; uuid user_id; uuid role_id FK }
    companies { uuid id PK; uuid tenant_id FK; text name; text vat_number; text commercial_registration; char currency; text timezone; timestamptz deleted_at }
    branches { uuid id PK; uuid tenant_id; uuid company_id; text code; bool is_main; timestamptz deleted_at }
    warehouses { uuid id PK; uuid tenant_id; uuid company_id; uuid branch_id; text code; timestamptz deleted_at }
    audit_logs { uuid id PK; uuid tenant_id FK; uuid user_id FK; text action; text entity_type; uuid entity_id; jsonb old_values; jsonb new_values; inet ip_address }
```

### الجداول

| الجدول | الغرض | `tenant_id` | Soft delete | RLS |
|---|---|---|---|---|
| `tenants` | المنشأة (حدود العزل) | — (هو المعرّف) | `deleted_at` | ✓ |
| `tenant_settings` | العملة، المنطقة الزمنية، بداية السنة المالية | ✓ (PK) | — | ✓ |
| `users` | هوية عامة (بريد + كلمة مرور) | — | `deleted_at` | — (انظر الملاحظات) |
| `user_tenants` | عضوية المستخدم في منشأة | ✓ | `status` | ✓ |
| `user_sessions` | جلسات Refresh Token (hash فقط) | ✓ | `revoked_at` | — (وصول خدمة المصادقة فقط) |
| `permissions` | كتالوج الصلاحيات | — | — | — (للقراءة فقط) |
| `roles` | أدوار النظام (`tenant_id` فارغ) والأدوار المخصصة | ✓ أو NULL | `deleted_at` | ✓ |
| `role_permissions` | ربط الأدوار بالصلاحيات | عبر الدور | — | ✓ |
| `user_roles` | أدوار المستخدم داخل المنشأة | ✓ | — | ✓ |
| `companies` | الشركات | ✓ | `deleted_at` | ✓ |
| `branches` | الفروع | ✓ | `deleted_at` | ✓ |
| `warehouses` | المستودعات (الأرصدة في المرحلة 5) | ✓ | `deleted_at` | ✓ |
| `audit_logs` | سجل التدقيق، إضافة فقط | ✓ (NULL قبل معرفة المنشأة) | ممنوع | ✓ |
| `schema_migrations` | الترحيلات المطبقة وبصمتها | — | — | — |

### العلاقات والقيود المهمة

- `branches (company_id, tenant_id)` → `companies (id, tenant_id)`: مفتاح أجنبي مركّب، فلا يمكن ربط فرع بشركة من منشأة أخرى حتى عبر SQL مباشر.
- `warehouses` يرتبط بالشركة وبالفرع بالطريقة نفسها.
- `user_roles (tenant_id, user_id)` → `user_tenants (tenant_id, user_id)`: لا يُسند دور إلا لعضو في المنشأة. ويتحقق Trigger من أن الدور نظامي أو يتبع المنشأة نفسها.
- `roles`: قيد CHECK يفرض أن أدوار النظام بلا `tenant_id` وأن الأدوار المخصصة لها `tenant_id`.
- قيود CHECK على: الرقم الضريبي السعودي (15 رقمًا يبدأ وينتهي بـ 3)، السجل التجاري (10 أرقام)، رموز العملة والدولة، الحالات.
- `audit_logs`: Trigger يمنع UPDATE وDELETE حتى لمالك المخطط.

### الفهارس

| الجدول | الفهرس |
|---|---|
| `tenants` | `slug` (فريد)، `status` |
| `users` | `email` (فريد، غير حساس لحالة الأحرف) |
| `user_tenants` | `(tenant_id, user_id)` (فريد)، `user_id` |
| `user_sessions` | `token_hash` (فريد)، `user_id` للجلسات غير الملغاة |
| `permissions` | `code` (فريد) |
| `roles` | `code` فريد لأدوار النظام؛ `(tenant_id, code)` فريد للأدوار المخصصة غير المحذوفة؛ `tenant_id` |
| `role_permissions` | PK `(role_id, permission_id)`، `permission_id` |
| `user_roles` | `(tenant_id, user_id, role_id)` (فريد)، `role_id` |
| `companies` | `tenant_id`؛ `(tenant_id, name)` فريد؛ `(tenant_id, vat_number)` فريد؛ `(id, tenant_id)` فريد للمفاتيح المركّبة |
| `branches` | `(tenant_id, company_id)`؛ `(tenant_id, company_id, code)` فريد؛ فرع رئيسي واحد لكل شركة |
| `warehouses` | `(tenant_id, company_id)`؛ `(tenant_id, company_id, code)` فريد |
| `audit_logs` | `(tenant_id, created_at DESC)`، `(tenant_id, entity_type, entity_id)`، `(user_id, created_at DESC)`، `(tenant_id, action)` |

الفهارس الفريدة على الجداول ذات الحذف الناعم جزئية (`WHERE deleted_at IS NULL`) حتى يمكن إعادة استخدام الرمز بعد الحذف.

## واجهات API (المرحلة 1)

جميع المدخلات يُتحقق منها بـ Zod. الأخطاء بالشكل `{ "error": { "code", "message", "details" } }`.

| الطريقة | المسار | الصلاحية |
|---|---|---|
| GET | `/api/health` | عام |
| POST | `/api/auth/register` | عام (10/دقيقة) |
| POST | `/api/auth/login` | عام (10/دقيقة) |
| POST | `/api/auth/refresh` | Cookie + ترويسة `X-CSRF-Protection: 1` |
| POST | `/api/auth/logout` | مسجّل دخول + ترويسة CSRF |
| POST | `/api/auth/switch-tenant` | مسجّل دخول وعضو في المنشأة الهدف |
| POST | `/api/auth/change-password` | مسجّل دخول |
| GET | `/api/auth/me` | مسجّل دخول |
| GET / PATCH | `/api/settings/tenant` | عضو / `settings.manage` |
| GET | `/api/users`، `/api/users/:id` | `user.view` |
| POST | `/api/users` | `user.invite` + `user.manage` |
| PATCH | `/api/users/:id` (تفعيل/تعطيل) | `user.manage` |
| PUT | `/api/users/:id/roles` | `user.manage` |
| GET | `/api/permissions`، `/api/roles` | `role.view` |
| POST / PATCH / DELETE | `/api/roles`، `/api/roles/:id` | `role.manage` |
| GET | `/api/companies`، `/api/companies/:id` | `company.view` |
| POST / PATCH / DELETE | `/api/companies`، `/api/companies/:id` | `company.manage` |
| GET / POST | `/api/companies/:id/branches` | `company.view` / `company.manage` |
| PATCH | `/api/branches/:id` | `company.manage` |
| GET / POST | `/api/companies/:id/warehouses` | `company.view` / `company.manage` |
| PATCH | `/api/warehouses/:id` | `company.manage` |
| GET | `/api/audit-logs` (فلاتر: action، entityType، entityId، userId، from، to، cursor) | `audit.view` |

## المرحلة 2: المحرك المحاسبي

### القواعد

كل عملية مالية في النظام (يدوية الآن، ومن الفواتير والمدفوعات لاحقًا) تمر عبر `modules/accounting/engine.ts`. دوال المحرك تعمل داخل معاملة المستدعي، فالمستند وقيده يُحفظان معًا أو يُلغيان معًا.

| القاعدة | التطبيق | قاعدة البيانات |
|---|---|---|
| مجموع المدين = مجموع الدائن، وأكبر من صفر | عند الترحيل (decimal.js) | Trigger يعيد الحساب عند `POSTED` |
| سطران على الأقل | ✓ | ✓ |
| كل سطر مدين **أو** دائن فقط، موجب، بخانتين عشريتين | ✓ (المبالغ نصوص، لا float) | `CHECK` + `NUMERIC(18,2)` |
| الحساب فرعي (يقبل القيود)، نشط، ومن نفس الشركة | ✓ | Trigger + مفتاح أجنبي مركّب |
| التاريخ داخل فترة مالية مفتوحة في سنة مفتوحة | ✓ مع قفل `FOR SHARE` | Trigger |
| القيد المرحّل لا يُعدَّل ولا يُحذف، وكذلك سطوره | ✓ | Trigger يرفض أي تغيير إلا التحول إلى `REVERSED` |
| لا حذف نهائي للقيود | حذف ناعم للمسودات فقط | دور التطبيق بلا صلاحية `DELETE` |
| ترقيم متسلسل بلا فجوات لكل شركة وسنة | يُسند عند الترحيل `JV-2026-000001` | جدول تسلسل بقفل صف، وفهرس فريد |

**دورة حياة القيد:** `DRAFT` ← `POSTED` ← `REVERSED`.

- المسودة تُعدَّل وتُحذف (حذفًا ناعمًا).
- القيد المرحّل يُصحَّح فقط بقيد عكسي (`REVERSAL`) بنفس الحسابات والمبالغ معكوسة، ثم بقيد تصحيح جديد. نقطة `reverse` تنشئ الاثنين في معاملة واحدة إذا أُرسل `correction`.
- القيد الأصلي يبقى في الدفتر بحالة `REVERSED`، والقيدان معًا يصفّران الأثر.
- لا يمكن عكس قيد عكسي، ولا عكس قيد صادر عن مستند (`source = SYSTEM`) يدويًا؛ يُصحَّح من المستند نفسه.

**إقفال السنة المالية:**

1. يتطلب أن تكون السنوات السابقة مقفلة، وألا توجد مسودات داخل السنة.
2. يرحّل قيد إقفال (`YEAR_CLOSING`) بتاريخ آخر يوم في السنة، يصفّر حسابات الإيرادات والتكاليف والمصروفات، ويحوّل صافي الربح أو الخسارة إلى الأرباح المحتجزة.
3. يقفل جميع الفترات ثم السنة. لا يمكن بعدها الترحيل في السنة ولا إعادة فتح فتراتها.

**دليل الحسابات:** يُنشأ تلقائيًا لكل شركة جديدة، مع السنة المالية الحالية (12 فترة شهرية حسب شهر بداية السنة في إعدادات المنشأة).

- يستخدم المحرك «مفتاح النظام» (`system_key`) مثل `ACCOUNTS_RECEIVABLE` و `VAT_OUTPUT` و `RETAINED_EARNINGS`، لا رمز الحساب، فيمكن للمستخدم إعادة ترقيم الحسابات.
- قواعد الشجرة مفروضة في قاعدة البيانات: نوع الفرع = نوع الأصل، الأصل حساب تجميعي، لا دوائر، والمستوى يُحسب تلقائيًا.
- لا يمكن تغيير نوع حساب له قيود أو فروع، ولا إيقاف حساب له رصيد، ولا حذف حساب عليه قيود أو حساب نظامي.

### الجداول الجديدة

| الجدول | الغرض |
|---|---|
| `fiscal_years` | السنوات المالية؛ قيد `EXCLUDE` يمنع التداخل لنفس الشركة |
| `fiscal_periods` | الفترات الشهرية؛ لا تداخل، ولا تُعاد فتح فترة في سنة مقفلة |
| `account_groups` | تصنيف الحسابات للقوائم المالية (أصول متداولة، مصروفات تشغيلية…) |
| `accounts` | دليل الحسابات الشجري |
| `cost_centers` | مراكز التكلفة |
| `journal_entries` | رؤوس القيود |
| `journal_entry_lines` | سطور القيود |
| `journal_sequences` | عداد الترقيم لكل سنة مالية |

```mermaid
erDiagram
    companies ||--o{ fiscal_years : has
    fiscal_years ||--o{ fiscal_periods : "split into"
    fiscal_years ||--|| journal_sequences : numbers
    companies ||--o{ account_groups : has
    companies ||--o{ accounts : has
    accounts ||--o{ accounts : "parent_id"
    account_groups ||--o{ accounts : classifies
    companies ||--o{ cost_centers : has
    companies ||--o{ journal_entries : has
    fiscal_periods ||--o{ journal_entries : "posted in"
    journal_entries ||--o{ journal_entry_lines : has
    accounts ||--o{ journal_entry_lines : "(account_id, company_id)"
    cost_centers ||--o{ journal_entry_lines : tags
    branches ||--o{ journal_entry_lines : tags
    journal_entries ||--o| journal_entries : "reversal_of_id / correction_of_id"

    accounts { uuid id PK; uuid company_id; text code; text account_type; uuid parent_id; smallint level; bool is_postable; text system_key }
    journal_entries { uuid id PK; uuid company_id; text entry_number; date entry_date; text status; text source; numeric total_debit; numeric total_credit; uuid reversal_of_id }
    journal_entry_lines { uuid id PK; uuid journal_entry_id; smallint line_no; uuid account_id; numeric debit; numeric credit; uuid cost_center_id; uuid branch_id }
    fiscal_years { uuid id PK; uuid company_id; date start_date; date end_date; text status; uuid closing_entry_id }
    fiscal_periods { uuid id PK; uuid fiscal_year_id; smallint period_number; date start_date; date end_date; text status }
```

كل الجداول الجديدة عليها `tenant_id` و RLS، وكل المراجع بينها مفاتيح أجنبية مركّبة مع `company_id`، فلا يمكن لسطر قيد أن يشير إلى حساب أو مركز تكلفة أو فرع من شركة أخرى.

**الفهارس الرئيسية:**

- `journal_entries`: `(company_id, fiscal_year_id, entry_number)` فريد؛ `(tenant_id, company_id, entry_date)`؛ `(company_id, status)`؛ `(company_id, reference_type, reference_id)`؛ قيد عكسي واحد فقط لكل قيد.
- `journal_entry_lines`: `account_id`؛ `(tenant_id, company_id, account_id)`؛ `(journal_entry_id, line_no)` فريد.
- `accounts`: `(company_id, code)` و `(company_id, system_key)` فريدان لغير المحذوفة؛ `parent_id`.
- `fiscal_periods`: `(company_id, start_date, end_date)`.

### واجهات API (المرحلة 2)

`companyId` اختياري عندما تملك المنشأة شركة واحدة. المبالغ تُرسل وتُستقبل نصوصًا (`"1150.00"`).

| الطريقة | المسار | الصلاحية |
|---|---|---|
| POST | `/api/accounting/setup` (للشركات القديمة؛ آمن للتكرار) | `account.manage` + `fiscal.manage` |
| GET | `/api/accounts`، `/api/accounts/:id` (مع الرصيد)، `/api/account-groups` | `account.view` |
| POST / PATCH / DELETE | `/api/accounts`، `/api/accounts/:id` | `account.manage` |
| GET / POST / PATCH | `/api/cost-centers`، `/api/cost-centers/:id` | `account.view` / `account.manage` |
| GET / POST | `/api/fiscal-years` | `account.view` / `fiscal.manage` |
| POST | `/api/fiscal-years/:id/close` | `fiscal.manage` |
| POST | `/api/fiscal-periods/:id/close`، `/api/fiscal-periods/:id/reopen` | `fiscal.manage` |
| GET | `/api/journal-entries` (فلاتر: status، dateFrom، dateTo، accountId، referenceType، search) | `journal.view` |
| GET | `/api/journal-entries/:id` | `journal.view` |
| POST | `/api/journal-entries` (مع `post: true` يتطلب `journal.post`) | `journal.create` |
| PATCH / DELETE | `/api/journal-entries/:id` (المسودات فقط) | `journal.create` |
| POST | `/api/journal-entries/:id/post` | `journal.post` |
| POST | `/api/journal-entries/:id/reverse` (`reason`، `date`، `correction`) | `journal.reverse` |
| GET | `/api/reports/trial-balance?dateFrom&dateTo` (رصيد افتتاحي، حركة، ختامي) | `report.view` |

## المرحلة 3: العملاء والموردون والمنتجات

### قرارات التصميم

- **الرصيد من الدفتر وحده.** لا يوجد حقل «رصيد» مخزّن للعميل أو المورد. سطر القيد يحمل `customer_id` أو `supplier_id`، والرصيد = مجموع (مدين − دائن) للسطور المرحّلة الموسومة بالطرف. يمكن إدخال الأرصدة الافتتاحية للعملاء والموردين بقيد يومية عادي يُختار فيه العميل أو المورد على سطر حساب العملاء أو الموردين، وكشف الحساب في المرحلة 7 يُبنى على الأساس نفسه.
- **العملاء والموردون بتنفيذ واحد** (`parties.routes.ts`) مع جداول وصلاحيات وترقيم وحساب رقابة منفصلة لكل منهما.
- **العنوان الوطني** بحقول ZATCA: رقم المبنى (4 أرقام)، الشارع، الحي، المدينة، الرمز البريدي (5 أرقام)، الرقم الإضافي (4 أرقام). الصيغة مفروضة في قاعدة البيانات. عنوان افتراضي واحد لكل نوع (فوترة / شحن). استبدال العناوين يحذفها حذفًا ناعمًا، لأن الفواتير ستحفظ نسختها الخاصة من العنوان.
- **الترقيم التلقائي** (`document_sequences`): عداد لكل شركة ونوع مستند بقفل صف داخل المعاملة، فلا فجوات ولا تكرار. مثل `CUS-00001` و `SUP-00001` و `PRD-00001`. إذا أدخل المستخدم رمزًا يدويًا يطابق رمزًا قادمًا، يتخطاه المولّد. المرحلة 4 ستستخدم الجدول نفسه لأرقام الفواتير.
- **المنتج لا يحمل نسبة الضريبة**، بل فئة ZATCA فقط: `S` أساسية، `Z` صفرية، `E` معفى، `O` خارج النطاق. النسبة تأتي من إعدادات الضريبة في المرحلة 6.
- **أسعار الوحدة بأربع خانات عشرية** `NUMERIC(18,4)`، وإجماليات الفواتير ستُقرَّب لخانتين.
- **الخدمة لا تتتبع مخزونًا**، والقاعدة مفروضة بقيد `CHECK`.
- **الوحدات** بأكواد UN/ECE (PCE، KGM، LTR، HUR…) لأنها مطلوبة في XML الفوترة الإلكترونية، وتُنشأ تلقائيًا لكل شركة.
- **حسابات بديلة اختيارية:** حساب رقابة للعميل (أصل) أو المورد (التزام)، وحساب مبيعات (إيراد) أو مشتريات (أصل / مصروف / تكلفة) للمنتج. النوع يُتحقق منه.

### الجداول الجديدة

| الجدول | الغرض |
|---|---|
| `customers`، `suppliers` | الأطراف: الرمز، الاسم، الرقم الضريبي، السجل التجاري، الهوية، حد الائتمان، مدة السداد، حساب الرقابة |
| `customer_addresses`، `supplier_addresses` | العناوين الوطنية |
| `units` | وحدات القياس |
| `product_categories` | تصنيفات المنتجات (شجرية) |
| `products` | المنتجات والخدمات |
| `document_sequences` | عدادات الترقيم |
| `journal_entry_lines` (تعديل) | عمودا `customer_id` و `supplier_id` مع مفتاح أجنبي مركّب، وطرف واحد كحد أقصى لكل سطر |

**الفهارس:** رمز فريد لكل شركة؛ رقم ضريبي فريد لكل شركة؛ SKU وباركود فريدان لكل شركة (للسجلات غير المحذوفة)؛ فهارس `pg_trgm` على الأسماء العربية لسرعة البحث؛ فهارس على `customer_id` و `supplier_id` في سطور القيود لحساب الأرصدة.

### واجهات API (المرحلة 3)

| الطريقة | المسار | الصلاحية |
|---|---|---|
| GET | `/api/customers` (search، status، limit، offset)، `/api/customers/:id` | `customer.view` |
| POST / PATCH / DELETE | `/api/customers`، `/api/customers/:id` | `customer.manage` |
| PUT | `/api/customers/:id/addresses` | `customer.manage` |
| — | نفس المسارات تحت `/api/suppliers` | `supplier.view` / `supplier.manage` |
| GET / POST / PATCH | `/api/units`، `/api/product-categories` | `product.view` / `product.manage` |
| GET | `/api/products` (search بالاسم أو الرمز أو الباركود، categoryId، productType، status) | `product.view` |
| POST / PATCH / DELETE | `/api/products`، `/api/products/:id` | `product.manage` |

`POST /api/journal-entries` يقبل `customerId` أو `supplierId` على كل سطر.

## المرحلة 4: المبيعات والمشتريات والمدفوعات

### المستندات

ستة مستندات بتنفيذ واحد (`documents/kinds.ts`)، لكل منها جداوله وصلاحياته وترقيمه:

| المستند | الجدول | الترقيم | دورة الحياة | قيد محاسبي |
|---|---|---|---|---|
| عرض سعر | `sales_quotes` | `QT-` عند الإنشاء | مسودة ← مرسل ← مقبول/مرفوض ← محوّل لفاتورة | لا |
| فاتورة مبيعات | `sales_invoices` | `INV-` عند الإصدار | مسودة ← صادرة ← مدفوعة جزئيًا ← مدفوعة / مرتجعة / ملغاة | نعم |
| مرتجع مبيعات (إشعار دائن) | `sales_returns` | `CN-` عند الإصدار | مسودة ← صادر ← ملغى | نعم |
| أمر شراء | `purchase_orders` | `PO-` عند الإنشاء | مسودة ← معتمد ← محوّل لفاتورة | لا |
| فاتورة مشتريات | `purchase_invoices` | `PINV-` عند الترحيل | مسودة ← مرحّلة ← مدفوعة جزئيًا ← مدفوعة / مرتجعة / ملغاة | نعم |
| مرتجع مشتريات (إشعار مدين) | `purchase_returns` | `DN-` عند الإصدار | مسودة ← صادر ← ملغى | نعم |

أرقام الفواتير والمرتجعات تُسند عند الإصدار فقط، فتبقى متسلسلة بلا فجوات (المسودات بلا رقم).

**أسماء الحقول:** المواصفة تسمي حقلي الفاتورة `invoice_number` و `invoice_date`. استخدمت `doc_number` و `doc_date` في الجداول الستة لأن التنفيذ مشترك، وتظهر في الـ API باسم `number` و `date`. بقية الحقول بنفس أسماء المواصفة: `subtotal`، `discount_total`، `taxable_amount`، `tax_amount`، `total`، `paid_amount`، `remaining_amount`، `status`، `branch_id`، `warehouse_id`، `currency`.

### القيود التلقائية

| المستند | مدين | دائن |
|---|---|---|
| فاتورة مبيعات | العملاء (موسوم بالعميل) = الإجمالي | المبيعات (أو حساب المنتج) = الصافي، ضريبة المخرجات = الضريبة |
| مرتجع مبيعات | المبيعات، ضريبة المخرجات | العملاء (موسوم) |
| فاتورة مشتريات | المخزون للسلع المتتبعة، أو حساب الشراء للمنتج، أو حساب السطر؛ ضريبة المدخلات | الموردون (موسوم بالمورد) |
| مرتجع مشتريات | الموردون (موسوم) | الحسابات نفسها، ضريبة المدخلات |
| سند قبض من عميل | حساب طريقة الدفع (صندوق / بنك) | العملاء (موسوم) |
| استرداد لعميل | العملاء (موسوم) | حساب طريقة الدفع |
| سند صرف لمورد | الموردون (موسوم) | حساب طريقة الدفع |
| استرداد من مورد | حساب طريقة الدفع | الموردون (موسوم) |

- المستند وقيده يُحفظان في المعاملة نفسها؛ فشل أي منهما يلغي الاثنين.
- القيد مصدره `SYSTEM` ومرتبط بالمستند (`reference_type`، `reference_id`)، فلا يُعكس يدويًا.
- يُخزَّن حساب الترحيل على كل سطر عند الإصدار للمراجعة.
- **تكلفة البضاعة المباعة وحركة المخزون غير مطبقتين بعد** — المرحلة 5.

### الحساب والضريبة (`documents/calc.ts`)

- كل المبالغ تُحسب في الخادم. الواجهة تطلب `POST /api/<document>/calculate` لعرض الإجماليات أثناء الكتابة، ولا تحسب الضريبة بنفسها.
- كمية وسعر الوحدة حتى 4 خانات؛ مبالغ السطر تُقرَّب لخانتين.
- **ضريبة المستند = تقريب (مجموع الصافي لكل فئة ونسبة × النسبة)** مرة واحدة، متوافقًا مع ZATCA. ضريبة السطر للعرض فقط وقد تختلف عن المجموع بهللة.
- الأسعار الشاملة للضريبة: الصافي = تقريب((الكمية × السعر − الخصم) ÷ (1 + النسبة)). كل المبالغ المخزنة غير شاملة.
- **النسبة ليست في الكود:** جدول `tax_rates` بتواريخ سريان. النسبة الأساسية الحالية (15%) تُنشأ افتراضيًا، وإضافة نسبة جديدة بتاريخ تُغلق السابقة تلقائيًا، وكل مستند يأخذ النسبة السارية في تاريخه. الفئات `Z` و `E` و `O` نسبتها صفر.
- الخصم على مستوى السطر (مبلغ أو نسبة). الخصم المُدخل يُحفظ كما هو (`discount_basis`) لإعادة الحساب بدقة.

### المرتجعات

- تُنشأ من الفاتورة الأصلية بتحديد السطور والكميات، ولا يمكن إرجاع أكثر من الكمية المتبقية.
- المبالغ تُوزَّع بنسبة الكمية. إرجاع آخر كمية من السطر يأخذ الباقي بالضبط، وآخر مرتجع للفاتورة يأخذ الضريبة المتبقية بالضبط، فيصفّر الإرجاع الكامل الفاتورة دون فروق تقريب.
- يُطبَّق المرتجع على الفاتورة حتى مبلغها المفتوح، والزيادة (إن كانت الفاتورة مدفوعة) تبقى مستحقة للعميل وتُسترد بسند صرف مخصص على المرتجع.

### المدفوعات

- طرق الدفع قابلة للتخصيص، وكل طريقة مربوطة بحساب أصل. المنشأة تبدأ بـ: نقدًا، تحويل بنكي، بطاقة، STC Pay، تمارا. يمكن ربط تمارا بحساب تسوية مستقل.
- السند يُخصَّص على فاتورة أو أكثر. **التخصيص لا ينشئ قيدًا**، لأن رصيد العميل في الدفتر صحيح أصلًا؛ هو ربط للمطابقة فقط. المبلغ غير المخصص يبقى دفعة مقدمة على حساب الطرف ويمكن تخصيصه لاحقًا.
- لا يمكن تعديل السند. يُلغى السند فيُعكس قيده وتُفك تخصيصاته وتُعاد حالة الفواتير.
- حالة الفاتورة (صادرة / مدفوعة جزئيًا / مدفوعة / مرتجعة) تُشتق دائمًا من المدفوع والمرتجع.

### قواعد أخرى

- **حد الائتمان:** لا تصدر فاتورة ترفع رصيد العميل في الدفتر فوق حده.
- **نوع الفاتورة الضريبية:** فاتورة ضريبية (`STANDARD`) إذا كان للعميل رقم ضريبي، وإلا مبسطة (`SIMPLIFIED`).
- **لقطة الطرف:** اسم العميل أو المورد ورقمه الضريبي وعنوانه الوطني تُنسخ في المستند عند الإصدار.
- **الإلغاء:** يعكس القيد بتاريخ المستند. لا يُلغى مستند عليه دفعات أو مرتجعات، ولا إذا كانت فترته مقفلة.
- **الحماية في قاعدة البيانات:** بعد الإصدار لا يتغير إلا الحالة وحقول التسوية، والسطور للقراءة فقط، ولا حذف للمستندات أو السندات.
- لا يُحذف عميل أو مورد أو منتج مستخدم في مستند أو سند.

### واجهات API (المرحلة 4)

لكل مستند `<path>` من: `sales-quotes`، `invoices`، `sales-returns`، `purchase-orders`، `purchase-invoices`، `purchase-returns`:

| الطريقة | المسار | الوصف |
|---|---|---|
| GET | `/api/<path>` (status، partyId، dateFrom، dateTo، open، search) | قائمة |
| GET | `/api/<path>/:id` | المستند مع السطور والدفعات والمرتجعات |
| POST | `/api/<path>/calculate` | حساب بدون حفظ |
| POST / PATCH / DELETE | `/api/<path>`، `/api/<path>/:id` | المسودات |
| POST | `/api/<path>/:id/post` | إصدار / ترحيل (الفواتير والمرتجعات) |
| POST | `/api/<path>/:id/status` | حالة عرض السعر أو أمر الشراء |
| POST | `/api/<path>/:id/convert` | تحويل عرض أو أمر إلى مسودة فاتورة |
| POST | `/api/<path>/:id/cancel` | إلغاء مع السبب |

| الطريقة | المسار | الصلاحية |
|---|---|---|
| GET / POST | `/api/payments`، `/api/payments/:id` | `payment.view` / `payment.create` |
| POST | `/api/payments/:id/allocations` | `payment.create` |
| POST | `/api/payments/:id/void` | `payment.void` |
| GET / POST / PATCH | `/api/payment-methods` | `payment.view` / `settings.manage` |
| GET / POST | `/api/tax-rates` | `invoice.view` / `tax.manage` |

صلاحيات المستندات: المبيعات `invoice.view/create/edit/delete/post/cancel`، والمشتريات `purchase.view/create/post/cancel`.

## المرحلة 5: المخزون

### المبادئ

- **لا تُعدَّل الكمية مباشرة.** كل تغيير حركة في `stock_movements` (إضافة فقط، لا تعديل ولا حذف)، ويُحدَّث `inventory_balances` (الكمية والقيمة لكل صنف ومستودع) في المعاملة نفسها بقفل صف.
- **المتوسط المرجح لكل مستودع.** الوارد يضيف قيمته، والصادر يخرج بقيمة تقريب(المتوسط × الكمية، خانتان)، و**آخر وحدة تأخذ القيمة المتبقية بالضبط**، فلا تبقى قيمة بلا كمية (قيد في قاعدة البيانات يمنع ذلك).
- **دفتر المخزون = حساب المخزون في الدفتر العام دائمًا.** القيم بالهللة، وكل حركة لها أثرها في القيد نفسه. تقرير التقييم يعرض الرقمين والفرق، والاختبارات تتحقق من الفرق صفر بعد كل عملية.
- **لا رصيد سالب.** البيع أو التحويل أو التسوية بأكثر من المتاح يُرفض (`INSUFFICIENT_STOCK`) ولا يتغير شيء.
- **طريقة التكلفة قابلة للتوسعة:** `costing.ts` يعرّف واجهة `CostingMethod`، والمتوسط المرجح تنفيذها الحالي. عمود `companies.costing_method` يقبل حاليًا `WEIGHTED_AVERAGE` فقط. FIFO يحتاج جدول طبقات تكلفة ويُضاف كتنفيذ آخر للواجهة نفسها.

### الحركات والقيود

| العملية | الحركة | التكلفة | القيد |
|---|---|---|---|
| فاتورة مشتريات | `PURCHASE` وارد | الصافي قبل الضريبة | المخزون مدين (السلع المتتبعة تُرحَّل دائمًا لحساب المخزون) |
| فاتورة مبيعات | `SALE` صادر | المتوسط | ضمن قيد الفاتورة: تكلفة المبيعات مدين / المخزون دائن |
| مرتجع مبيعات | `SALE_RETURN` وارد | التكلفة التي خرجت بها البضاعة، بالتناسب، وآخر مرتجع يأخذ الباقي بالضبط | المخزون مدين / تكلفة المبيعات دائن |
| مرتجع مشتريات | `PURCHASE_RETURN` صادر | المتوسط | المخزون دائن بالتكلفة، والفرق عن سعر الشراء لتكلفة المبيعات |
| تحويل | `TRANSFER_OUT` ثم `TRANSFER_IN` | المتوسط في المستودع المصدر | لا قيد (القيمة باقية في حساب المخزون) |
| تسوية / جرد | `ADJUSTMENT_IN` أو `ADJUSTMENT_OUT` | الزيادة بالتكلفة المُدخلة أو المتوسط الحالي؛ النقص بالمتوسط | المخزون مقابل حساب الفروقات (5200) أو حساب يختاره المستخدم (رأس المال للرصيد الافتتاحي) |
| إلغاء مستند | `CANCELLATION_IN` / `CANCELLATION_OUT` | نفس تكلفة الحركة الأصلية | يُعكس قيد المستند؛ وإن فرغ المستودع وبقيت قيمة تُحمَّل على تكلفة المبيعات بقيد `INVENTORY_RESIDUAL` |

- المرتجعات تعود إلى (أو تخرج من) مستودع الحركة الأصلية.
- لا يُلغى شراء استُخدمت بضاعته؛ يُستخدم مرتجع المشتريات بدلًا من ذلك.
- مستودع الفاتورة يُختار في النموذج، وإلا فالمستودع الرئيسي.
- الجرد يُدخل الكمية المعدودة ويحسب الخادم الفرق لرصيد المستودع المختار.
- لا يتغير نوع منتج أو تتبع مخزونه ولا يُحذف بعد وجود حركات له.

### الجداول الجديدة

| الجدول | الغرض |
|---|---|
| `inventory_balances` | الكمية والقيمة لكل صنف ومستودع (المفتاح: الصنف + المستودع)؛ لا سالب، ولا قيمة بلا كمية |
| `stock_movements` | كل حركة: النوع، الاتجاه، الكمية، تكلفة الوحدة، الإجمالي، الرصيد بعد الحركة، المرجع (المستند والسطر)، والحركة المعكوسة عند الإلغاء |
| `stock_transfers`، `stock_transfer_items` | التحويلات (`TRF-`) |
| `stock_adjustments`، `stock_adjustment_items` | التسويات والجرد (`ADJ-`) مع الكمية المعدودة والكمية في النظام |
| حساب `5200` | فروقات وتسويات المخزون (مفتاح النظام `INVENTORY_ADJUSTMENT`) |

**الفهارس:** حركة الصنف `(product_id, warehouse_id, created_at)`؛ المرجع `(reference_type, reference_id)`؛ السطر `reference_line_id`؛ التاريخ لكل شركة؛ إلغاء واحد فقط لكل حركة.

### واجهات API (المرحلة 5)

| الطريقة | المسار | الصلاحية |
|---|---|---|
| GET | `/api/inventory/balances` (warehouseId، productId) | `inventory.view` |
| GET | `/api/inventory/movements` (productId، warehouseId، dateFrom، dateTo) — حركة الصنف | `inventory.view` |
| GET | `/api/inventory/valuation` — التقييم مع مطابقة الدفتر | `inventory.view` |
| GET / POST | `/api/stock-transfers`، `/api/stock-transfers/:id` | `inventory.view` / `inventory.transfer` |
| GET / POST | `/api/stock-adjustments`، `/api/stock-adjustments/:id` | `inventory.view` / `inventory.adjust` |

`GET /api/products` يعرض `onHand` (الكمية في كل المستودعات)، والمستندات تقبل `warehouseId`.

## تدفق المصادقة

1. **التسجيل** ينشئ في معاملة واحدة: المستخدم، المنشأة، إعداداتها (SAR، Asia/Riyadh)، الشركة الافتراضية، الفرع الرئيسي، المستودع الرئيسي، العضوية بدور `TENANT_OWNER`، وسجل تدقيق.
2. **الدخول** يتحقق من كلمة المرور بـ argon2id. البريد غير الموجود يمر بنفس زمن التحقق ويعيد الرسالة نفسها. بعد 5 محاولات فاشلة يُقفل الحساب 15 دقيقة.
3. يُصدر **Access Token** (JWT لمدة 15 دقيقة يحمل `sub` و`tid` المنشأة و`sid` الجلسة) و**Refresh Token** عشوائي (32 بايت) يُخزن hash له فقط، ويُرسل في Cookie بخصائص `HttpOnly` و`SameSite=Strict` و`Secure` في الإنتاج، ومساره `/api/auth`.
4. **في كل طلب** يُعاد التحقق من قاعدة البيانات: الجلسة غير ملغاة، المستخدم نشط، العضوية نشطة، المنشأة غير معلقة. وتُحمَّل الصلاحيات من جديد، فتغيير الأدوار أو التعطيل أو الخروج يسري فورًا.
5. **التجديد** يدوّر الرمز. إذا قُدّم رمز سبق تدويره (سرقة أو إعادة استخدام) تُلغى كل جلسات المستخدم ويُسجل `TOKEN_REUSE_DETECTED`.
6. **تبديل المنشأة** ينقل الجلسة للمنشأة الجديدة، فيصبح الـ Access Token القديم غير صالح.
7. **تغيير كلمة المرور** يُخرج المستخدم من جميع الأجهزة الأخرى.

الواجهة تحفظ الـ Access Token في الذاكرة فقط (لا `localStorage`)، وتستعيد الجلسة عند تحديث الصفحة عبر الـ Cookie.

## عزل المنشآت (Tenant Isolation)

ثلاث طبقات مستقلة:

1. **التطبيق:** `tenant_id` يؤخذ دائمًا من الجلسة الموثقة، لا من جسم الطلب أو الرابط. كل استعلام يضيف `WHERE tenant_id = $1`. طلب سجل منشأة أخرى يعيد 404 لا 403، حتى لا يكشف وجوده.
2. **قاعدة البيانات (RLS):** الـ API يتصل بدور `alshuyukh_app` غير المميز (`NOBYPASSRLS`). كل معاملة تضبط `app.tenant_id` و`app.user_id` بـ `set_config(..., true)` فتنتهي بانتهاء المعاملة ولا تتسرب لطلب لاحق على الاتصال نفسه. سياسات RLS ترفض القراءة والكتابة خارج المنشأة حتى لو نسي الكود الفلتر. بلا سياق، النتيجة فارغة.
3. **القيود:** مفاتيح أجنبية مركّبة `(id, tenant_id)` تمنع الربط بين منشأتين. دور التطبيق لا يملك صلاحية `DELETE` على الجداول التجارية (حذف ناعم فقط)، ولا يستطيع تعديل أدوار النظام.

## الصلاحيات (RBAC)

- المصدر الوحيد: `apps/api/src/modules/rbac/catalog.ts`، يُزامَن مع قاعدة البيانات عند `db:migrate`.
- 45 صلاحية (منها صلاحيات المراحل القادمة، حتى تبقى الأدوار ثابتة).
- أدوار النظام: `TENANT_OWNER`، `COMPANY_ADMIN`، `ACCOUNTANT`، `SALES_MANAGER`، `SALES_EMPLOYEE`، `PURCHASE_MANAGER`، `WAREHOUSE_MANAGER`، `WAREHOUSE_EMPLOYEE`، `VIEWER`.
- `SUPER_ADMIN` ليس دورًا داخل منشأة؛ هو علامة منصة `users.is_platform_admin` ولوحته `/admin` في المرحلة 9.
- **منع تصعيد الصلاحيات:** لا يمكن منح (أو سحب) صلاحية لا تملكها، عند إنشاء دور أو تعديله أو إسناده. لا يمكن تغيير أدوارك أو حالتك بنفسك. المالك لا يُعطّل ولا يُجرّد من دور `TENANT_OWNER`. أدوار النظام غير قابلة للتعديل (ومحمية بـ RLS أيضًا).
- يمكن للمنشأة إنشاء أدوار مخصصة.

## سجل التدقيق

يُكتب داخل المعاملة نفسها التي تُجري التغيير، فلا يوجد تغيير بلا سجل. يسجل: `REGISTER`، `LOGIN`، `LOGIN_FAILED`، `LOGOUT`، `TOKEN_REUSE_DETECTED`، `TENANT_SWITCH`، `PASSWORD_CHANGE`، `CREATE`، `UPDATE`، `DELETE`، `SETTINGS_CHANGE`، `PERMISSION_CHANGE`، مع القيم قبل وبعد، وعنوان IP، والـ User-Agent، ومعرّف الطلب. تُحذف كلمات المرور والرموز من القيم قبل الحفظ.

## الأمان

argon2id لكلمات المرور · JWT قصير العمر مع جلسات قابلة للإلغاء · RBAC · عزل بثلاث طبقات · Rate limiting (عام 300/دقيقة، الدخول والتسجيل 10/دقيقة، التجديد 120/دقيقة) · Zod على كل المدخلات · SQL بمعاملات فقط، وأسماء الأعمدة في التحديثات من قائمة بيضاء · Helmet (CSP، HSTS، nosniff، frame-ancestors) · CORS بقائمة أصول محددة · حماية CSRF لمسارات الـ Cookie · قفل الحساب · رسائل أخطاء لا تكشف التفاصيل الداخلية.

## الاختبارات

188 اختبارًا (وحدة وتكامل) على PostgreSQL حقيقي باستخدام دور التطبيق المقيد، فتُختبر سياسات RLS فعليًا:

| الملف | ما يغطيه |
|---|---|
| `auth.test.ts` | التسجيل وما ينشئه، تخزين argon2id، التحقق من المدخلات، الدخول، تطابق رسائل الفشل، قفل الحساب، رفض الرموز المعدلة، تدوير الرموز، كشف إعادة الاستخدام، الخروج، تغيير كلمة المرور، ترويسات الأمان |
| `tenant-isolation.test.ts` | محاولات منشأة B الوصول لبيانات A عبر المعرفات (قراءة، تعديل، حذف، فروع، مستخدمون، سجل التدقيق، تبديل، أدوار)؛ ثم مباشرة على قاعدة البيانات: الدور لا يتجاوز RLS، الاستعلام بلا فلتر، بلا سياق، الإدخال في منشأة أخرى، المفتاح المركّب، التعديل عبر المنشآت، أدوار النظام، منع الحذف |
| `rbac.test.ts` | سلامة الكتالوج، المحاسب والمشاهد، سريان تغيير الأدوار فورًا، التعطيل الفوري، منع التصعيد، حماية المالك، أدوار النظام، الأدوار المخصصة |
| `audit.test.ts` | تسجيل العمليات، القيم قبل/بعد، عدم تخزين الأسرار، منع التعديل والحذف، الترقيم |
| `journal-entries.test.ts` | مثال الفاتورة 1000 + ضريبة 150؛ الترحيل والترقيم المتسلسل؛ رفض القيد غير المتوازن مع التفاصيل؛ التراجع الكامل عند الفشل؛ دقة الكسور (0.10 + 0.20)؛ التحقق من كل سطر؛ الحسابات التجميعية وغير النشطة والأجنبية؛ غياب الفترة؛ ترقيم بلا فجوات تحت التزامن؛ تعديل وحذف المسودات؛ رفض تعديل القيد المرحّل عبر API وعبر SQL مباشر؛ فرض التوازن في قاعدة البيانات؛ العكس والتصحيح؛ الفترات المقفلة؛ الصلاحيات؛ العزل؛ سجل التدقيق؛ ميزان المراجعة |
| `chart-of-accounts.test.ts` | الدليل الافتراضي ومفاتيح النظام؛ السنة المالية التلقائية؛ تجهيز الشركات اللاحقة؛ الحسابات الفرعية وتحويل الحساب إلى تجميعي؛ منع الدوائر والتكرار وعدم تطابق النوع؛ حماية الحسابات النظامية وذات الرصيد؛ الإيقاف والحذف الناعم؛ مراكز التكلفة؛ العزل |
| `fiscal-years.test.ts` | إنشاء السنوات ومنع التداخل؛ السنة القصيرة؛ قفل شهر البداية؛ إقفال السنة بربح وبخسارة وبلا حركة؛ منع الترحيل وإعادة الفتح بعد الإقفال؛ اشتراط إقفال السنوات السابقة |
| `parties.test.ts` | لكل من العملاء والموردين: الترقيم التلقائي وتخطي الرموز المأخوذة؛ التحقق من الرقم الضريبي والسجل والهوية وصيغة العنوان الوطني والعنوان الافتراضي؛ البحث؛ التعديل واستبدال العناوين والإيقاف؛ حساب الرقابة؛ الرصيد من القيود الموسومة ومنع الحذف بعدها؛ العزل بما فيه وسم قيد بطرف من منشأة أخرى؛ الصلاحيات لكل دور؛ سجل التدقيق |
| `products.test.ts` | الوحدات الافتراضية؛ الوحدات والتصنيفات المتداخلة؛ السلع والخدمات؛ الأسعار بأربع خانات؛ التحقق والتكرار؛ الحسابات البديلة؛ البحث بالباركود؛ الحذف الناعم؛ رفض وحدات منشأة أخرى؛ الصلاحيات |
| `calc.test.ts` | مثال 1000 + 150؛ الخصم قبل الضريبة؛ التقريب مرة واحدة لكل فئة؛ الأسعار الشاملة؛ أربع خانات؛ الفئات الصفرية والمعفاة؛ رفض السطور غير الصحيحة |
| `sales.test.ts` | قيد الفاتورة المطلوب (AR 1150 / Sales 1000 / VAT 150)؛ الفاتورة المبسطة؛ الترقيم؛ الحساب دون حفظ؛ الأسعار الشاملة؛ منع التعديل بعد الإصدار عبر API و SQL؛ حد الائتمان؛ تغيير نسبة الضريبة بتاريخ؛ الإلغاء ومنعه مع الدفعات؛ المرتجع الجزئي ثم الكامل بلا فروق تقريب؛ إلغاء المرتجع؛ عروض الأسعار والتحويل؛ الصلاحيات؛ العزل؛ منع حذف المستخدَم؛ توازن الميزان |
| `payments.test.ts` | التحصيل الكامل والجزئي؛ سند واحد لعدة فواتير؛ الدفعة المقدمة ثم التخصيص؛ رفض التخصيص الزائد أو لطرف آخر أو لنوع خاطئ دون أثر؛ إلغاء السند؛ الاسترداد بعد المرتجع؛ منع تعديل السند في قاعدة البيانات؛ دفع المورد؛ طرق الدفع وحساب التسوية؛ الصلاحيات؛ العزل |
| `purchases.test.ts` | قيد فاتورة الشراء (مخزون + ضريبة مدخلات / موردون)؛ رقم فاتورة المورد مرة واحدة؛ سطور المصروفات والمعفاة؛ اشتراط حساب للخدمات؛ مرتجع المشتريات؛ الإلغاء؛ أوامر الشراء والتحويل؛ الصلاحيات |
| `inventory.test.ts` | الرصيد الافتتاحي مقابل رأس المال؛ الشراء والمتوسط؛ تكلفة المبيعات ضمن قيد الفاتورة؛ رفض البيع بلا رصيد دون أثر؛ المرتجع بتكلفة الخروج بالضبط؛ استعادة المخزون عند الإلغاء؛ تقريب 100 على 3 دون بقايا؛ مرتجع المشتريات بالمتوسط وفرق السعر؛ منع إلغاء شراء استُخدم؛ تحميل القيمة المتبقية؛ التحويلات بلا قيد؛ البيع من مستودع محدد؛ الجرد والتحقق؛ قفل نوع المنتج؛ منع تعديل الحركات وسالب الرصيد؛ الصلاحيات؛ العزل؛ مطابقة التقييم مع الدفتر بعد كل خطوة |
| `companies.test.ts` | القيم السعودية الافتراضية، التحقق، التكرار، الحذف الناعم، ربط المستودع بفرع شركة أخرى، تجاهل `tenant_id` في جسم الطلب، العضوية في منشأتين والتبديل |

## ملاحظات وقرارات تصميم

- **المنشأة مقابل الشركة:** المنشأة (tenant) هي حدود العزل والاشتراك. يمكن أن تملك المنشأة أكثر من شركة. السجلات المالية في المراحل القادمة ستحمل `tenant_id` و`company_id` معًا.
- **جدول `users` بلا RLS:** الهوية عامة لأن الدخول يحتاج البحث بالبريد قبل معرفة المنشأة. لا يحتوي على بيانات تجارية، والوصول إليه يمر دائمًا عبر `user_tenants`. يمكن لاحقًا نقل البحث إلى دالة `SECURITY DEFINER` وتفعيل RLS عليه.
- **المبالغ المالية:** مخزنة `NUMERIC(18,2)` (حتى الهللة). `pg` مضبوط ليعيد `NUMERIC` كنص، والحسابات في الخادم بـ decimal.js، والـ API يقبل المبالغ نصوصًا فقط. الواجهة تعرض المجاميع أثناء الكتابة بالهللات (أعداد صحيحة) للعرض فقط، والخادم هو المرجع.
- **القيود العكسية:** الأصل يبقى بحالة `REVERSED` ويُحسب في الأرصدة مع قيده العكسي، فيبقى الدفتر كاملًا والأثر صفرًا.
- **الوقت:** كل الاتصالات بتوقيت UTC، والعرض بتوقيت المنشأة (افتراضيًا Asia/Riyadh).

## المتبقي (TODO)

- دعوة المستخدمين بالبريد بدل كلمة المرور المؤقتة (مع نظام الإشعارات).
- تطبيق حدود الباقات (مستخدمون، فروع، شركات) — المرحلة 9.
- لوحة `/admin` لمدير المنصة — المرحلة 9.
- رفع شعار الشركة (حاليًا رابط https فقط).
- نقل ملكية المنشأة.
- Rate limiting يعتمد على ذاكرة العملية؛ في الإنتاج بأكثر من نسخة يلزم Redis.
- لم تُطبَّق المصادقة الثنائية (2FA).
- **المرحلة 5:**
  - **متغيرات المنتج** (`product_variants`: مقاس، لون) غير مطبقة: تحتاج بُعدًا إضافيًا في سطور المستندات والحركات والأرصدة. البديل الحالي: صنف مستقل لكل متغير.
  - **مواقع المستودع** (`warehouse_locations`: رفوف) غير مطبقة؛ الرصيد على مستوى المستودع.
  - **FIFO** غير مطبق (الواجهة جاهزة).
  - المتوسط المرجح دائم حسب ترتيب الترحيل: إدخال حركة بتاريخ سابق لا يعيد حساب تكلفة الحركات اللاحقة.
  - لا يُلغى تحويل أو تسوية؛ التصحيح بمستند معاكس.
  - السماح بالرصيد السالب كإعداد اختياري غير متاح.
  - تنبيهات حد إعادة الطلب ونقص المخزون — مع نظام الإشعارات.
  - قوائم الأسعار وتسعير العملاء.
- **المرحلة 4:**
  - الخصم على مستوى المستند (لا السطر) غير مطبق.
  - تعدد العملات: المستندات بعملة الشركة فقط.
  - بعد ربط ZATCA (المرحلة 8) يجب منع إلغاء الفاتورة المُبلَّغة والاكتفاء بإشعار دائن، وإضافة رمز سبب الإعفاء للسطور `Z` و `E` و `O`.
  - طباعة الفاتورة وإرسالها بالبريد، ورمز QR — المرحلة 8.
  - تقادم الذمم وكشف الحساب — المرحلة 7.
  - تجاوز حد الائتمان بصلاحية خاصة غير متاح.
- **المرحلة 3:**
  - استيراد العملاء والموردين والمنتجات من ملف Excel/CSV.
- **المرحلة 2:**
  - تعدد العملات وأسعار الصرف: القيود حاليًا بعملة الشركة فقط.
  - إعادة فتح سنة مالية مقفلة (عكس قيد الإقفال) غير مدعومة.
  - استيراد الأرصدة الافتتاحية من ملف: حاليًا تُدخل بقيد يدوي.
  - سير موافقات للقيود (من يُنشئ لا يُرحّل) غير مطبق؛ الفصل الحالي عبر صلاحيتي `journal.create` و `journal.post` فقط.
  - دفتر الأستاذ وبقية التقارير في المرحلة 7. التقارير يجب أن تستثني قيد `YEAR_CLOSING` من قائمة الدخل.
