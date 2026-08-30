import type { ReactNode } from 'react'

export function AuthShell({
  eyebrow,
  title,
  subtitle,
  children,
  footer,
}: {
  eyebrow: string
  title: string
  subtitle: string
  children: ReactNode
  footer?: ReactNode
}) {
  return (
    <div className="flex min-h-screen items-center justify-center px-4">
      <div className="w-full max-w-[420px] animate-rise">
        <div className="mb-8 flex items-center gap-3">
          <span className="relative flex h-2.5 w-2.5">
            <span className="absolute h-full w-full rounded-full bg-signal animate-pulse-signal" />
          </span>
          <span className="font-display text-lg font-medium tracking-tight text-ink">Signalline</span>
        </div>

        <div className="rounded-2xl border border-line bg-surface p-8 shadow-2xl shadow-black/40">
          <p className="mb-1 font-mono text-xs uppercase tracking-[0.2em] text-signal">{eyebrow}</p>
          <h1 className="mb-2 font-display text-2xl font-semibold text-ink">{title}</h1>
          <p className="mb-6 text-sm text-muted">{subtitle}</p>
          {children}
        </div>

        {footer && <div className="mt-5 text-center text-sm text-muted">{footer}</div>}
      </div>
    </div>
  )
}
