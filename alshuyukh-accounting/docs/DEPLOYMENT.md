# دليل النشر — الشيوخ للمحاسبة

هذا الدليل يشرح تشغيل النظام على خادم واحد باستخدام Docker Compose. جُرِّبت هذه الخطوات كاملة: بناء الصور، والترحيلات، والتسجيل، واختبار المتصفح، والنسخ الاحتياطي والاستعادة.

## المكونات

| الخدمة | الصورة | الدور |
|---|---|---|
| `postgres` | `postgres:16` | قاعدة البيانات. لا تنشر منفذًا خارج الشبكة الداخلية |
| `migrate` | `alshuyukh-api` | مهمة تعمل مرة واحدة: تطبق الترحيلات بدور مالك المخطط ثم تنتهي |
| `api` | `alshuyukh-api` | واجهة API بدور التطبيق المقيد (RLS مفعّل). تعمل بمستخدم غير root ونظام ملفات للقراءة فقط |
| `web` | `alshuyukh-web` | nginx بمستخدم غير root: يخدم الواجهة ويمرر `/api` إلى الـ API، مع ترويسات الأمان و CSP |

لا تبدأ `api` إلا بعد نجاح `migrate`، ولا تبدأ `web` إلا بعد أن تصبح `api` سليمة.

## المتطلبات

- خادم Linux عليه Docker Engine و Docker Compose v2.
- اسم نطاق وشهادة TLS. يُنهى TLS قبل حاوية `web` (موازن أحمال، أو Caddy، أو nginx مع certbot)، ثم يُفعَّل سطر `Strict-Transport-Security` في `deploy/nginx.conf`.
- نسخ احتياطي خارج الخادم (تخزين كائنات أو خادم آخر).

## التشغيل لأول مرة

```sh
cp deploy/.env.example deploy/.env
# املأ الأسرار:
#   openssl rand -hex 32      لكل من POSTGRES_PASSWORD و OWNER_DB_PASSWORD و APP_DB_PASSWORD
#   openssl rand -base64 48   لـ JWT_SECRET
#   openssl rand -base64 32   لـ ZATCA_ENCRYPTION_KEY
chmod 600 deploy/.env
docker compose -f docker-compose.prod.yml --env-file deploy/.env up -d --build
```

تحقق من الجاهزية:

```sh
curl -s http://localhost:8080/api/ready     # {"status":"ready"}
```

- `/api/health`: العملية تعمل وتصل لقاعدة البيانات (فحص الحياة).
- `/api/ready`: قاعدة البيانات متاحة ولا ترحيلات معلقة (فحص الجاهزية). يرد 503 قبل ذلك.

ثم امنح حسابك صلاحية مدير المنصة بعد التسجيل من الواجهة:

```sh
docker compose -f docker-compose.prod.yml --env-file deploy/.env run --rm migrate node dist/admin-cli.js grant you@example.com
```

يعمل الأمر عبر خدمة `migrate` لأنها تتصل بدور مالك المخطط. دور التطبيق لا يستطيع منح صلاحية المنصة، حتى لو اختُرق الـ API.

## إعدادات مهمة

| المتغير | الملاحظة |
|---|---|
| `JWT_SECRET` | يرفض التشغيل في الإنتاج إن كان أقصر من 32 حرفًا أو القيمة الافتراضية |
| `ZATCA_ENCRYPTION_KEY` | إلزامي في الإنتاج. يشفّر مفاتيح ZATCA الخاصة وأسرار CSID. **احتفظ بنسخة منه خارج الخادم**: فقده يعني إعادة ربط كل وحدات الفوترة |
| `TRUST_PROXY` | `1` في ملف Compose لأن nginx أمام الـ API. إن أضفت موازن أحمال أمام nginx فارفعه إلى `2` أو اكتب عناوين الوكلاء. لا تضبطه إن كان الـ API مكشوفًا مباشرة، فيمكن حينها تزوير عنوان العميل وتجاوز حد المحاولات |
| `PUBLIC_ORIGIN` | العنوان الذي يفتحه المستخدمون؛ يُستخدم لـ CORS |
| `DB_STATEMENT_TIMEOUT_MS` | يلغي الاستعلام الأطول من ذلك (افتراضيًا 15 ثانية) ويرد 503 |
| `DB_POOL_MAX` | حجم مجمع الاتصالات (افتراضيًا 20) |

## التحديث

```sh
git pull
docker compose -f docker-compose.prod.yml --env-file deploy/.env up -d --build
```

تُطبَّق الترحيلات الجديدة تلقائيًا قبل بدء الـ API. الترحيل الذي عُدّل بعد تطبيقه يوقف التشغيل عمدًا: اكتب ترحيلًا جديدًا بدل تعديل القديم.

## النسخ الاحتياطي والاستعادة

```sh
scripts/backup.sh /var/backups/alshuyukh          # نسخة مضغوطة، ويحذف ما مضى عليه 14 يومًا
scripts/restore.sh /var/backups/alshuyukh/alshuyukh-20260101T021500Z.dump
```

- النسخة بصيغة `pg_dump --format=custom` وتحفظ الملكية والصلاحيات، فتبقى سياسات RLS سارية على دور التطبيق بعد الاستعادة. يتحقق السكربت من إمكانية قراءة النسخة قبل اعتمادها.
- جدولة يومية عبر cron مثلًا: `15 2 * * * cd /srv/alshuyukh && scripts/backup.sh /var/backups/alshuyukh`.
- الاستعادة تطلب كتابة `restore` للتأكيد، وتوقف الـ API أثناءها ثم تطبق أي ترحيلات أحدث من النسخة.
- انسخ الملفات خارج الخادم، وجرّب الاستعادة دوريًا على خادم تجريبي.

## البناء خلف وكيل يفحص TLS

إن كان الخادم خلف وكيل شركة يعيد توقيع اتصالات TLS، مرّر شهادته أثناء البناء:

```sh
docker build -f deploy/api.Dockerfile -t alshuyukh-api --secret id=npm_ca,src=/path/to/proxy-ca.crt --build-arg HTTPS_PROXY .
docker build -f deploy/web.Dockerfile -t alshuyukh-web --secret id=npm_ca,src=/path/to/proxy-ca.crt --build-arg HTTPS_PROXY .
```

## اختبار ما بعد النشر

```sh
npm ci
E2E_BASE_URL=https://staging.example.com npm run e2e
```

يفتح متصفحًا حقيقيًا، وينشئ منشأة جديدة، وعميلًا وخدمة، ويصدر فاتورة بقيمة 1000 ريال، ثم يتحقق من القيد (مدين 1150 / دائن 1000 / دائن 150) والتحصيل وتوازن ميزان المراجعة. يفشل الاختبار عند أي خطأ في المتصفح، ومنه مخالفات CSP. شغّله على بيئة تجريبية لأنه ينشئ بيانات. لتحديد مسار المتصفح استخدم `CHROMIUM_PATH`.

## حدود هذا الإعداد

- خادم واحد. مع أكثر من نسخة من الـ API تحتاج حدود المحاولات وعدادات الاستخدام إلى Redis، ويُفضَّل فصل مُرسل ZATCA في عملية مستقلة.
- لا مراقبة أو تنبيهات مدمجة: اربط `/api/ready` بأداة مراقبة، وراجع لوحة «صحة النظام» في `/admin`.
- السجلات تخرج على stdout بصيغة JSON؛ اجمعها بأداة السجلات المعتمدة لديك.
