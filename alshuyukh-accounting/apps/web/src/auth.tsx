import { createContext, useCallback, useContext, useEffect, useState, type ReactNode } from 'react';
import { api, refreshSession, setAccessToken, setSessionExpiredHandler } from './api';

export interface Me {
  user: { id: string; email: string; fullName: string; mustChangePassword: boolean; isPlatformAdmin: boolean };
  tenant: { id: string; name: string; status: string; isOwner: boolean };
  roles: { id: string; code: string; nameAr: string }[];
  permissions: string[];
  memberships: { tenantId: string; tenantName: string; isOwner: boolean }[];
  subscription: { state: 'TRIALING' | 'ACTIVE' | 'GRACE' | 'EXPIRED' | 'CANCELLED' | 'NONE'; writable: boolean; planName: string | null; periodEnd: string | null; graceEnd: string | null };
  features: Record<string, boolean>;
}

interface AuthState {
  me: Me | null;
  loading: boolean;
  can: (permission: string) => boolean;
  login: (email: string, password: string) => Promise<void>;
  register: (data: Record<string, unknown>) => Promise<void>;
  logout: () => Promise<void>;
  switchTenant: (tenantId: string) => Promise<void>;
  reload: () => Promise<void>;
}

const Ctx = createContext<AuthState | null>(null);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [me, setMe] = useState<Me | null>(null);
  const [loading, setLoading] = useState(true);

  const reload = useCallback(async () => setMe(await api<Me>('GET', '/api/auth/me')), []);

  useEffect(() => {
    setSessionExpiredHandler(() => setMe(null));
    // Restore the session from the refresh cookie on page load.
    refreshSession().then(async (ok) => { if (ok) await reload().catch(() => setMe(null)); }).finally(() => setLoading(false));
  }, [reload]);

  const value: AuthState = {
    me,
    loading,
    can: (p) => !!me?.permissions.includes(p),
    login: async (email, password) => {
      const r = await api<{ accessToken: string }>('POST', '/api/auth/login', { email, password });
      setAccessToken(r.accessToken);
      await reload();
    },
    register: async (data) => {
      const r = await api<{ accessToken: string }>('POST', '/api/auth/register', data);
      setAccessToken(r.accessToken);
      await reload();
    },
    logout: async () => {
      await api('POST', '/api/auth/logout').catch(() => undefined);
      setAccessToken(null);
      setMe(null);
    },
    switchTenant: async (tenantId) => {
      const r = await api<{ accessToken: string }>('POST', '/api/auth/switch-tenant', { tenantId });
      setAccessToken(r.accessToken);
      await reload();
    },
    reload,
  };
  return <Ctx.Provider value={value}>{children}</Ctx.Provider>;
}

export function useAuth(): AuthState {
  const v = useContext(Ctx);
  if (!v) throw new Error('useAuth outside AuthProvider');
  return v;
}
