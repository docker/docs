import { z } from 'zod';
import { badRequest } from './errors.js';

/** Parses untrusted input. Throws a 400 with field-level details on failure. */
export function parse<T extends z.ZodType>(schema: T, data: unknown): z.infer<T> {
  const result = schema.safeParse(data);
  if (!result.success) {
    throw badRequest(
      'VALIDATION_ERROR',
      'Invalid request',
      result.error.issues.map((i) => ({ path: i.path.join('.'), message: i.message })),
    );
  }
  return result.data;
}

export const uuidParam = z.object({ id: z.uuid() });

const TIMEZONES = new Set(Intl.supportedValuesOf('timeZone'));
export const timezone = z.string().refine((tz) => TIMEZONES.has(tz), 'Unknown time zone');

/** Trims, and turns empty strings into null for optional text fields. */
export const optionalText = (max = 500) =>
  z.preprocess(
    (v) => (typeof v === 'string' && v.trim() === '' ? null : typeof v === 'string' ? v.trim() : v),
    z.string().max(max).nullable().optional(),
  );

export const password = z
  .string()
  .min(10, 'Password must be at least 10 characters')
  .max(200)
  .refine((p) => /[A-Za-z]/.test(p) && /[0-9]/.test(p), 'Password must contain letters and digits');
