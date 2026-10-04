/**
 * API client. The access token is kept in memory only (never localStorage).
 * The refresh token lives in an HttpOnly cookie the browser sends to
 * /api/auth/refresh; on a 401 the client refreshes once and retries.
 */
let accessToken: string | null = null;
let refreshing: Promise<boolean> | null = null;
let onSessionExpired: () => void = () => {};

export const setAccessToken = (t: string | null) => { accessToken = t; };
export const setSessionExpiredHandler = (fn: () => void) => { onSessionExpired = fn; };

export class ApiError extends Error {
  constructor(public status: number, public code: string, message: string, public details?: unknown) {
    super(message);
  }
}

export async function refreshSession(): Promise<boolean> {
  refreshing ??= (async () => {
    try {
      const res = await fetch('/api/auth/refresh', { method: 'POST', credentials: 'include', headers: { 'X-CSRF-Protection': '1' } });
      if (!res.ok) return false;
      setAccessToken((await res.json()).accessToken);
      return true;
    } catch {
      return false;
    } finally {
      setTimeout(() => { refreshing = null; }, 0);
    }
  })();
  return refreshing;
}

export async function api<T = unknown>(method: string, path: string, body?: unknown, retry = true): Promise<T> {
  const headers: Record<string, string> = { 'X-CSRF-Protection': '1' };
  if (body !== undefined) headers['Content-Type'] = 'application/json';
  if (accessToken) headers.Authorization = `Bearer ${accessToken}`;
  const res = await fetch(path, { method, headers, credentials: 'include', body: body === undefined ? undefined : JSON.stringify(body) });

  if (res.status === 401 && retry && !path.startsWith('/api/auth/login') && !path.startsWith('/api/auth/register')) {
    if (await refreshSession()) return api<T>(method, path, body, false);
    setAccessToken(null);
    onSessionExpired();
  }
  if (res.status === 204) return undefined as T;
  const data = await res.json().catch(() => ({}));
  if (!res.ok) {
    const e = data?.error ?? {};
    throw new ApiError(res.status, e.code ?? 'ERROR', translateError(e.code, e.message), e.details);
  }
  return data as T;
}

/** Downloads a file from an authenticated endpoint and hands it to the browser. */
export async function downloadFile(path: string, filename: string, retry = true): Promise<void> {
  const res = await fetch(path, { headers: { 'X-CSRF-Protection': '1', ...(accessToken ? { Authorization: `Bearer ${accessToken}` } : {}) }, credentials: 'include' });
  if (res.status === 401 && retry && await refreshSession()) return downloadFile(path, filename, false);
  if (!res.ok) {
    const e = (await res.json().catch(() => ({})))?.error ?? {};
    throw new ApiError(res.status, e.code ?? 'ERROR', translateError(e.code, e.message), e.details);
  }
  const a = document.createElement('a');
  a.href = URL.createObjectURL(await res.blob());
  a.download = filename;
  document.body.append(a);
  a.click();
  a.remove();
  setTimeout(() => URL.revokeObjectURL(a.href), 1000);
}

const MESSAGES: Record<string, string> = {
  UNAUTHORIZED: 'البريد الإلكتروني أو كلمة المرور غير صحيحة',
  ACCOUNT_LOCKED: 'تم إيقاف الحساب مؤقتًا بسبب محاولات دخول فاشلة متكررة. حاول لاحقًا.',
  EMAIL_TAKEN: 'يوجد حساب مسجل بهذا البريد الإلكتروني',
  FORBIDDEN: 'لا تملك صلاحية تنفيذ هذا الإجراء',
  NOT_FOUND: 'السجل غير موجود',
  VALIDATION_ERROR: 'تحقق من البيانات المدخلة',
  DUPLICATE: 'توجد قيمة مكررة لحقل يجب أن يكون فريدًا',
  TENANT_SUSPENDED: 'تم تعليق هذه المنشأة',
  ALREADY_MEMBER: 'المستخدم عضو في المنشأة بالفعل',
  LAST_COMPANY: 'يجب أن تبقى شركة واحدة على الأقل',
  ROLE_IN_USE: 'الدور مُسند لمستخدمين؛ أزله منهم أولًا',
  INVALID_CURRENT_PASSWORD: 'كلمة المرور الحالية غير صحيحة',
  UNBALANCED_ENTRY: 'القيد غير متوازن: مجموع المدين لا يساوي مجموع الدائن',
  TOO_FEW_LINES: 'القيد يحتاج سطرين على الأقل',
  INVALID_LINE: 'كل سطر يجب أن يحتوي على مبلغ مدين أو دائن فقط',
  INVALID_AMOUNT: 'مبلغ غير صحيح: يجب أن يكون موجبًا وبخانتين عشريتين كحد أقصى',
  INVALID_ACCOUNT: 'الحساب غير موجود في هذه الشركة',
  ACCOUNT_NOT_POSTABLE: 'لا يمكن الترحيل على حساب تجميعي',
  ACCOUNT_INACTIVE: 'الحساب غير نشط',
  NO_FISCAL_PERIOD: 'لا توجد فترة مالية تغطي هذا التاريخ. أنشئ السنة المالية أولًا.',
  PERIOD_CLOSED: 'الفترة المالية لهذا التاريخ مقفلة',
  ENTRY_NOT_DRAFT: 'لا يمكن تعديل قيد مرحّل؛ استخدم العكس',
  ALREADY_REVERSED: 'تم عكس هذا القيد مسبقًا',
  CANNOT_REVERSE_REVERSAL: 'لا يمكن عكس قيد عكسي؛ أنشئ قيدًا جديدًا',
  SYSTEM_ENTRY: 'هذا القيد صادر عن مستند؛ صحّحه من المستند نفسه',
  ACCOUNT_HAS_BALANCE: 'لا يمكن إيقاف حساب له رصيد',
  SYSTEM_ACCOUNT: 'حساب نظامي يستخدمه المحرك المحاسبي ولا يمكن إيقافه أو حذفه',
  ACCOUNT_IN_USE: 'الحساب عليه قيود؛ يمكن إيقافه بدل حذفه',
  DRAFTS_EXIST: 'توجد مسودات قيود في هذه السنة؛ رحّلها أو احذفها أولًا',
  EARLIER_YEAR_OPEN: 'أقفل السنوات المالية السابقة أولًا',
  YEAR_CLOSED: 'السنة المالية مقفلة',
  FISCAL_YEAR_OVERLAP: 'السنة المالية تتداخل مع سنة موجودة',
  FISCAL_YEARS_EXIST: 'لا يمكن تغيير بداية السنة المالية بعد إنشاء سنوات مالية',
  DEVICE_STATE: 'لا يمكن تنفيذ هذه الخطوة في المرحلة الحالية لوحدة الفوترة',
  COMPLIANCE_INCOMPLETE: 'شغّل فحوص الامتثال حتى تنجح جميعها قبل التفعيل',
  NOT_PENDING: 'هذه الفاتورة الإلكترونية ليست بانتظار الإرسال (أو يجري إرسالها الآن)',
  ZATCA_NOT_CONFIGURED: 'مفتاح تشفير الفوترة الإلكترونية غير مضبوط في الخادم',
  OWN_TENANT: 'لا يمكنك إيقاف المنشأة التي سجلت الدخول بها',
  SELF: 'لا يمكنك تنفيذ هذا الإجراء على حسابك',
  DEFAULT_PLAN: 'الباقة الافتراضية يجب أن تبقى متاحة؛ اختر باقة افتراضية أخرى أولًا',
  NO_SUBSCRIPTION: 'اختر باقة لهذه المنشأة أولًا',
  PLAN_INACTIVE: 'هذه الباقة لم تعد متاحة',
};
const translateError = (code?: string, fallback?: string) => (code && MESSAGES[code]) || fallback || 'حدث خطأ غير متوقع';
