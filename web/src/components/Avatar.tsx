import { useEffect, useState } from 'react'
import { getAvatarUrl } from '../api/media'
import { fingerprintColor, initials } from '../utils/identity'

interface Props {
  avatarId?: string | null
  seed: string
  name: string
  className?: string
}

export function Avatar({ avatarId, seed, name, className = '' }: Props) {
  const [errored, setErrored] = useState(false)

  useEffect(() => {
    setErrored(false)
  }, [avatarId])

  if (avatarId && !errored) {
    return (
      <img
        src={getAvatarUrl(avatarId)}
        alt={name}
        onError={() => setErrored(true)}
        className={`shrink-0 rounded-full object-cover ${className}`}
      />
    )
  }

  return (
    <div
      className={`flex shrink-0 items-center justify-center rounded-full font-display font-semibold text-void ${className}`}
      style={{ backgroundColor: fingerprintColor(seed) }}
    >
      {initials(name)}
    </div>
  )
}
