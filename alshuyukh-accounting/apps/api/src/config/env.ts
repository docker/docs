import { z } from 'zod';

const schema = z.object({
  NODE_ENV: z.enum(['development', 'test', 'production']).default('development'),
  DATABASE_URL: z.string().min(1),
  DATABASE_URL_MIGRATE: z.string().min(1).optional(),
  JWT_SECRET: z.string().min(32, 'JWT_SECRET must be at least 32 characters'),
  ACCESS_TOKEN_TTL_SECONDS: z.coerce.number().int().min(60).max(3600).default(900),
  REFRESH_TOKEN_TTL_DAYS: z.coerce.number().int().min(1).max(90).default(30),
  PORT: z.coerce.number().int().default(3000),
  HOST: z.string().default('0.0.0.0'),
  CORS_ORIGINS: z.string().default('http://localhost:5173'),
  LOGIN_MAX_ATTEMPTS: z.coerce.number().int().min(3).default(5),
  LOGIN_LOCK_MINUTES: z.coerce.number().int().min(1).default(15),
  // 32 random bytes, base64: encrypts ZATCA private keys and CSID secrets at rest.
  ZATCA_ENCRYPTION_KEY: z.string().optional(),
  // Background submission of pending e-invoices (reporting / clearance).
  // Overrides the Fatoora gateway base URL (a proxy, or a local stand-in for testing).
  ZATCA_GATEWAY_URL: z.url().optional(),
  ZATCA_WORKER: z.enum(['on', 'off']).default('on'),
  ZATCA_WORKER_INTERVAL_SECONDS: z.coerce.number().int().min(10).default(60),
});

export type Env = z.infer<typeof schema>;

export function loadEnv(source: NodeJS.ProcessEnv = process.env): Env {
  const parsed = schema.safeParse(source);
  if (!parsed.success) {
    const details = parsed.error.issues.map((i) => `${i.path.join('.')}: ${i.message}`).join('\n');
    throw new Error(`Invalid environment configuration:\n${details}`);
  }
  if (parsed.data.NODE_ENV === 'production' && parsed.data.JWT_SECRET.startsWith('change-me')) {
    throw new Error('JWT_SECRET must be changed in production');
  }
  if (parsed.data.NODE_ENV === 'production' && !parsed.data.ZATCA_ENCRYPTION_KEY) {
    throw new Error('ZATCA_ENCRYPTION_KEY is required in production (openssl rand -base64 32)');
  }
  return parsed.data;
}
