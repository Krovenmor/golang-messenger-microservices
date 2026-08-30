interface Props {
  connected: boolean
}

export function SignalDot({ connected }: Props) {
  return (
    <span className="inline-flex items-center gap-2 font-mono text-xs text-muted">
      <span className="relative flex h-2 w-2">
        <span
          className={`absolute inline-flex h-full w-full rounded-full ${
            connected ? 'bg-signal animate-pulse-signal' : 'bg-faint'
          }`}
        />
      </span>
      {connected ? 'на связи' : 'переподключение…'}
    </span>
  )
}
