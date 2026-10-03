import { useEffect, useState, type ReactNode } from 'react';
import { ApiError } from './api';

export function useLoad<T>(load: () => Promise<T>, deps: unknown[] = []) {
  const [data, setData] = useState<T | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [version, setVersion] = useState(0);
  useEffect(() => {
    let alive = true;
    setError(null);
    load().then((d) => alive && setData(d)).catch((e: Error) => alive && setError(e.message));
    return () => { alive = false; };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [...deps, version]);
  return { data, error, reload: () => setVersion((v) => v + 1) };
}

export function ErrorBox({ error }: { error: unknown }) {
  if (!error) return null;
  const msg = error instanceof Error ? error.message : String(error);
  // Only field-validation errors carry a list; other details are structured data.
  const details = error instanceof ApiError && Array.isArray(error.details) ? error.details : undefined;
  return (
    <div className="alert alert-error" role="alert">
      {msg}
      {details && <ul>{details.map((d) => <li key={d.path}><code>{d.path}</code>: {d.message}</li>)}</ul>}
    </div>
  );
}

export function PageHeader({ title, children }: { title: string; children?: ReactNode }) {
  return (
    <div className="page-header">
      <h1>{title}</h1>
      <div className="actions">{children}</div>
    </div>
  );
}

export const formatDateTime = (iso: string | null) =>
  iso ? new Intl.DateTimeFormat('ar-SA-u-nu-latn', { dateStyle: 'medium', timeStyle: 'short', timeZone: 'Asia/Riyadh' }).format(new Date(iso)) : '—';
