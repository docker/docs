import { useState } from 'react';
import { formatAmount } from '../money';

/**
 * Small SVG charts for the dashboard. Values arrive as decimal strings from the
 * server; Number() is used only to place marks, never for displayed amounts.
 * Colors come from CSS variables (--series-1, --series-2) validated for contrast
 * and colour-vision deficiency.
 */
const W = 640;
const H = 220;
const PAD = { top: 12, right: 12, bottom: 28, left: 64 };

function scale(values: number[]) {
  // With no data at all, show a round empty scale instead of fractions of a riyal.
  const empty = values.every((v) => v === 0);
  const max = empty ? 100 : Math.max(0, ...values);
  const min = Math.min(0, ...values);
  const span = max - min || 1;
  const step = niceStep(span / 4);
  const top = Math.ceil(max / step) * step || step;
  const bottom = Math.floor(min / step) * step;
  const ticks: number[] = [];
  for (let v = bottom; v <= top + step / 2; v += step) ticks.push(v);
  const y = (v: number) => PAD.top + (top - v) / (top - bottom) * (H - PAD.top - PAD.bottom);
  return { y, ticks };
}
function niceStep(raw: number) {
  const p = 10 ** Math.floor(Math.log10(raw));
  return [1, 2, 2.5, 5, 10].map((m) => m * p).find((s) => s >= raw) ?? 10 * p;
}
const compact = (v: number) => new Intl.NumberFormat('en', { notation: 'compact', maximumFractionDigits: 1 }).format(v);
const monthLabel = (m: string) => new Intl.DateTimeFormat('ar-SA-u-nu-latn-ca-gregory', { month: 'short', timeZone: 'UTC' }).format(new Date(`${m}-01T00:00:00Z`));

interface Series { key: string; label: string; color: string }

export function MonthlyBars({ data, series }: { data: ({ month: string } & Record<string, string>)[]; series: Series[] }) {
  const [hover, setHover] = useState<number | null>(null);
  const { y, ticks } = scale(data.flatMap((d) => series.map((s) => Number(d[s.key]))));
  const band = (W - PAD.left - PAD.right) / data.length;
  const barW = Math.min(14, (band - 10) / series.length);
  return (
    <div className="chart">
      <Legend series={series} />
      <div className="chart-plot" dir="ltr">
        <svg viewBox={`0 0 ${W} ${H}`} role="img" aria-label={series.map((s) => s.label).join(' و')}>
          <Axes ticks={ticks} y={y} />
          {data.map((d, i) => {
            const x0 = PAD.left + i * band + (band - barW * series.length - 2 * (series.length - 1)) / 2;
            return (
              <g key={d.month}>
                {series.map((s, j) => {
                  const v = Number(d[s.key]);
                  const top = Math.min(y(v), y(0));
                  const h = Math.max(Math.abs(y(v) - y(0)), v === 0 ? 0 : 1);
                  return <rect key={s.key} x={x0 + j * (barW + 2)} y={top} width={barW} height={h} rx={Math.min(4, barW / 2)} fill={s.color} opacity={hover === null || hover === i ? 1 : 0.45} />;
                })}
                <text x={PAD.left + i * band + band / 2} y={H - 8} textAnchor="middle" className="chart-tick">{monthLabel(d.month)}</text>
                <rect x={PAD.left + i * band} y={PAD.top} width={band} height={H - PAD.top - PAD.bottom} fill="transparent"
                  onMouseEnter={() => setHover(i)} onMouseLeave={() => setHover(null)} onFocus={() => setHover(i)} onBlur={() => setHover(null)} tabIndex={0} />
              </g>
            );
          })}
        </svg>
        {hover !== null && <Tooltip x={(PAD.left + hover * band + band / 2) / W} month={data[hover]!.month} rows={series.map((s) => [s.label, data[hover]![s.key]!, s.color])} />}
      </div>
    </div>
  );
}

export function MonthlyLine({ data, field, label, color }: { data: ({ month: string } & Record<string, string>)[]; field: string; label: string; color: string }) {
  const [hover, setHover] = useState<number | null>(null);
  const { y, ticks } = scale(data.map((d) => Number(d[field])));
  const step = (W - PAD.left - PAD.right) / Math.max(data.length - 1, 1);
  const x = (i: number) => PAD.left + i * step;
  const path = data.map((d, i) => `${i ? 'L' : 'M'}${x(i)},${y(Number(d[field]))}`).join(' ');
  return (
    <div className="chart">
      <div className="chart-plot" dir="ltr">
        <svg viewBox={`0 0 ${W} ${H}`} role="img" aria-label={label}>
          <Axes ticks={ticks} y={y} />
          <path d={path} fill="none" stroke={color} strokeWidth={2} strokeLinejoin="round" strokeLinecap="round" />
          {hover !== null && <line x1={x(hover)} x2={x(hover)} y1={PAD.top} y2={H - PAD.bottom} className="chart-crosshair" />}
          {data.map((d, i) => (
            <g key={d.month}>
              {(hover === i || i === data.length - 1) && <circle cx={x(i)} cy={y(Number(d[field]))} r={4} fill={color} stroke="var(--surface)" strokeWidth={2} />}
              <text x={x(i)} y={H - 8} textAnchor="middle" className="chart-tick">{monthLabel(d.month)}</text>
              <rect x={x(i) - step / 2} y={PAD.top} width={step} height={H - PAD.top - PAD.bottom} fill="transparent"
                onMouseEnter={() => setHover(i)} onMouseLeave={() => setHover(null)} onFocus={() => setHover(i)} onBlur={() => setHover(null)} tabIndex={0} />
            </g>
          ))}
        </svg>
        {hover !== null && <Tooltip x={x(hover) / W} month={data[hover]!.month} rows={[[label, data[hover]![field]!, color]]} />}
      </div>
    </div>
  );
}

function Axes({ ticks, y }: { ticks: number[]; y: (v: number) => number }) {
  return (
    <g>
      {ticks.map((t) => (
        <g key={t}>
          <line x1={PAD.left} x2={W - PAD.right} y1={y(t)} y2={y(t)} className={t === 0 ? 'chart-baseline' : 'chart-grid'} />
          <text x={PAD.left - 8} y={y(t) + 4} textAnchor="end" className="chart-tick">{compact(t)}</text>
        </g>
      ))}
    </g>
  );
}

function Legend({ series }: { series: Series[] }) {
  return (
    <div className="chart-legend">
      {series.map((s) => <span key={s.key}><i style={{ background: s.color }} />{s.label}</span>)}
    </div>
  );
}

function Tooltip({ x, month, rows }: { x: number; month: string; rows: [string, string, string][] }) {
  return (
    <div className="chart-tooltip" style={{ left: `${Math.min(Math.max(x * 100, 12), 88)}%` }} role="status">
      <strong dir="ltr">{month}</strong>
      {rows.map(([l, v, c]) => <div key={l}><i style={{ background: c }} />{l}<span dir="ltr">{formatAmount(v)}</span></div>)}
    </div>
  );
}
