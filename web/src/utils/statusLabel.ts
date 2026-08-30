import { ProfileStatus, type StatusInfo } from '../types'

export function statusLabel(status: StatusInfo | null): string {
  if (!status) return ''

  switch (status.status) {
    case ProfileStatus.Online:
      return 'в сети'
    case ProfileStatus.Typing:
      return 'печатает…'
    case ProfileStatus.Away:
      return 'отошёл(ла)'
    case ProfileStatus.Offline:
    default: {
      if (!status.lastSeen) return 'давно не был(а) в сети'
      const diffMs = Date.now() - status.lastSeen * 1000
      const diffMin = Math.round(diffMs / 60_000)
      if (diffMin < 1) return 'только что был(а) в сети'
      if (diffMin < 60) return `был(а) в сети ${diffMin} мин назад`
      const diffHours = Math.round(diffMin / 60)
      if (diffHours < 24) return `был(а) в сети ${diffHours} ч назад`
      const diffDays = Math.round(diffHours / 24)
      return `был(а) в сети ${diffDays} дн назад`
    }
  }
}

export function isStatusActive(status: StatusInfo | null): boolean {
  return status?.status === ProfileStatus.Online || status?.status === ProfileStatus.Typing
}
