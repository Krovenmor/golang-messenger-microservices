import { useEffect, useState } from 'react'
import { getAvatarUrl } from '../api/media'
import { fingerprintColor, initials } from '../utils/identity'

interface Props {
  /** photoId of the active avatar, if any (e.g. profile.additional.avatars?.[0]). */
  avatarId?: string | null
  /** Used to derive the fallback color and as alt text. */
  seed: string
  name: string
  className?: string
}

export function Avatar({ avatarId, seed, name, className = '' }: Props) {
  const [errored, setErrored] = useState(false)

  // Reset the error flag if we get handed a different (or new) avatar —
  // otherwise switching from a broken photo to a working one would still
  // show the fallback forever.
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
