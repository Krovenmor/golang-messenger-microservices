import { useState, type InputHTMLAttributes } from 'react'

interface Props extends InputHTMLAttributes<HTMLInputElement> {
  label: string
}

function EyeIcon() {
  return (
    <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
      <path d="M1 12s4-7 11-7 11 7 11 7-4 7-11 7-11-7-11-7z" strokeLinecap="round" strokeLinejoin="round" />
      <circle cx="12" cy="12" r="3" />
    </svg>
  )
}

function EyeOffIcon() {
  return (
    <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
      <path
        d="M17.94 17.94A10.94 10.94 0 0112 19c-7 0-11-7-11-7a18.7 18.7 0 015.06-5.94M9.9 4.24A10.94 10.94 0 0112 4c7 0 11 7 11 7a18.7 18.7 0 01-2.16 3.19M14.12 14.12a3 3 0 11-4.24-4.24"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
      <line x1="1" y1="1" x2="23" y2="23" strokeLinecap="round" />
    </svg>
  )
}

export function Field({ label, id, type, className, ...rest }: Props) {
  const isPassword = type === 'password'
  const [visible, setVisible] = useState(false)

  return (
    <label className="mb-4 block last:mb-0" htmlFor={id}>
      <span className="mb-1.5 block font-mono text-xs uppercase tracking-wide text-faint">{label}</span>
      <div className="relative">
        <input
          id={id}
          type={isPassword ? (visible ? 'text' : 'password') : type}
          className={`w-full rounded-lg border border-line bg-raised px-3.5 py-2.5 text-sm text-ink placeholder:text-faint outline-none transition-colors focus:border-signal disabled:cursor-not-allowed disabled:opacity-50 ${isPassword ? 'pr-11' : ''} ${className ?? ''}`}
          {...rest}
        />
        {isPassword && (
          <button
            type="button"
            tabIndex={-1}
            onClick={() => setVisible((v) => !v)}
            aria-label={visible ? 'Скрыть пароль' : 'Показать пароль'}
            className="absolute right-1.5 top-1/2 -translate-y-1/2 rounded-md p-1.5 text-faint hover:bg-void hover:text-ink"
          >
            {visible ? <EyeOffIcon /> : <EyeIcon />}
          </button>
        )}
      </div>
    </label>
  )
}
