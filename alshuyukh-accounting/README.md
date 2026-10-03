# ALSHUYUKH ACCOUNTING — الشيوخ للمحاسبة

نظام محاسبي سحابي متعدد المنشآت (Multi-Tenant SaaS) للسوق السعودي.

**الحالة:** المرحلة 1 مكتملة (تأسيس المشروع، قاعدة البيانات، المصادقة، تعدد المنشآت، الصلاحيات).
الوحدات المحاسبية لم تُبنَ بعد، وتظهر في الواجهة بعلامة «قريبًا» دون أي بيانات تجريبية.

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
│   │   │   ├── lib/                   # الأخطاء، التحقق، كلمات المرور، الرموز
│   │   │   ├── plugins/auth.ts        # المصادقة وحارس الصلاحيات
│   │   │   └── modules/
│   │   │       ├── auth/              # التسجيل، الدخول، التجديد، الخروج، التبديل
│   │   │       ├── users/             # أعضاء المنشأة
│   │   │       ├── rbac/              # catalog.ts (مصدر الصلاحيات) + الأدوار
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

54 اختبار تكامل على PostgreSQL حقيقي باستخدام دور التطبيق المقيد، فتُختبر سياسات RLS فعليًا:

| الملف | ما يغطيه |
|---|---|
| `auth.test.ts` | التسجيل وما ينشئه، تخزين argon2id، التحقق من المدخلات، الدخول، تطابق رسائل الفشل، قفل الحساب، رفض الرموز المعدلة، تدوير الرموز، كشف إعادة الاستخدام، الخروج، تغيير كلمة المرور، ترويسات الأمان |
| `tenant-isolation.test.ts` | محاولات منشأة B الوصول لبيانات A عبر المعرفات (قراءة، تعديل، حذف، فروع، مستخدمون، سجل التدقيق، تبديل، أدوار)؛ ثم مباشرة على قاعدة البيانات: الدور لا يتجاوز RLS، الاستعلام بلا فلتر، بلا سياق، الإدخال في منشأة أخرى، المفتاح المركّب، التعديل عبر المنشآت، أدوار النظام، منع الحذف |
| `rbac.test.ts` | سلامة الكتالوج، المحاسب والمشاهد، سريان تغيير الأدوار فورًا، التعطيل الفوري، منع التصعيد، حماية المالك، أدوار النظام، الأدوار المخصصة |
| `audit.test.ts` | تسجيل العمليات، القيم قبل/بعد، عدم تخزين الأسرار، منع التعديل والحذف، الترقيم |
| `companies.test.ts` | القيم السعودية الافتراضية، التحقق، التكرار، الحذف الناعم، ربط المستودع بفرع شركة أخرى، تجاهل `tenant_id` في جسم الطلب، العضوية في منشأتين والتبديل |

## ملاحظات وقرارات تصميم

- **المنشأة مقابل الشركة:** المنشأة (tenant) هي حدود العزل والاشتراك. يمكن أن تملك المنشأة أكثر من شركة. السجلات المالية في المراحل القادمة ستحمل `tenant_id` و`company_id` معًا.
- **جدول `users` بلا RLS:** الهوية عامة لأن الدخول يحتاج البحث بالبريد قبل معرفة المنشأة. لا يحتوي على بيانات تجارية، والوصول إليه يمر دائمًا عبر `user_tenants`. يمكن لاحقًا نقل البحث إلى دالة `SECURITY DEFINER` وتفعيل RLS عليه.
- **المبالغ المالية:** `pg` مضبوط ليعيد `NUMERIC` كنص، حتى لا يتحول أي مبلغ إلى float. المرحلة 2 ستستخدم `NUMERIC(19,4)` ومكتبة Decimal.
- **الوقت:** كل الاتصالات بتوقيت UTC، والعرض بتوقيت المنشأة (افتراضيًا Asia/Riyadh).

## المتبقي (TODO)

- دعوة المستخدمين بالبريد بدل كلمة المرور المؤقتة (مع نظام الإشعارات).
- تطبيق حدود الباقات (مستخدمون، فروع، شركات) — المرحلة 9.
- لوحة `/admin` لمدير المنصة — المرحلة 9.
- رفع شعار الشركة (حاليًا رابط https فقط).
- نقل ملكية المنشأة.
- منع حذف شركة لها قيود مرحّلة، ومنع تغيير بداية السنة المالية بعد إنشاء سنة — المرحلة 2.
- Rate limiting يعتمد على ذاكرة العملية؛ في الإنتاج بأكثر من نسخة يلزم Redis.
- لم تُطبَّق المصادقة الثنائية (2FA).
