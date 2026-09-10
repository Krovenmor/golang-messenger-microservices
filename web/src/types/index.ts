export interface AuthTokens {
  accessToken: string
  refreshToken: string
}

export interface ProfileAdditional {
  avatars?: string[]
  bio?: string
}

export interface OwnProfile {
  userId: string
  name: string
  userName: string
  pubKey: string
  encryptedPrvKey: string
  kdfSalt: string
  keyNonce: string
  createdAt: string
  additional: ProfileAdditional
}

export interface PublicProfile {
  userId: string
  name: string
  userName: string
  pubKey: string
  createdAt: string
  additional: ProfileAdditional
}

export interface MediaStorageInfo {
  maxSpace: number
  spaceFilled: number
  savedFiles: number
}

export type MediaFileType = 'photo'
export type MediaFileSubType = 'avatar'

export interface MediaFileInfo {
  mediaId: string
  type: MediaFileType
  subType: MediaFileSubType
  size: number
  isPublic: boolean
  addedAt: string
}

export interface Reaction {
  emoji: string
  users: string[]
}

export interface Message {
  messageId: string
  senderId: string
  message: string
  senderKey: string
  receiverKey: string
  nonce: string
  createdAt: string
  redactedAt: string | null
  deletedAt: string | null
  replyToId: string | null
  reactions?: Reaction[]
}

export interface ChatMember {
  userId: string
  joinedAt: string
}

export interface ChatSummary {
  chatId: string
  members: ChatMember[]
  lastMessage: Message | null
}

export interface ChatInfo {
  createdAt: string
  members: ChatMember[]
}

interface WSFrame<T extends string, P> {
  type: T
  payload: P
}

export type WSNewChatEvent = WSFrame<'newChat', { chatId: string }>
export type WSNewMessageEvent = WSFrame<'newMessage', { chatId: string; msgId: string }>
export type WSMessageRedactedEvent = WSFrame<'messageRedacted', { chatId: string; msgId: string }>
export type WSMessageDeletedEvent = WSFrame<'messageDeleted', { chatId: string; msgId: string }>
export type WSNewReactionEvent = WSFrame<'newReaction', { chatId: string; msgId: string; userId: string; emoji: string }>
export type WSDelReactionEvent = WSFrame<'delReaction', { chatId: string; msgId: string; userId: string; emoji: string }>
export type WSStatusEvent = WSFrame<'status', { userId: string; newStatus: ProfileStatus; eventTime: number }>
export type WSResponseEvent = WSFrame<'response', { code: number; msg: string }>

export type WSEvent =
  | WSNewChatEvent
  | WSNewMessageEvent
  | WSMessageRedactedEvent
  | WSMessageDeletedEvent
  | WSNewReactionEvent
  | WSDelReactionEvent
  | WSStatusEvent
  | WSResponseEvent

export type WSRequest =
  | { req: 1; payload: { newStatus: ProfileStatus } }
  | { req: 2; payload: { userId: string } }

export enum ProfileStatus {
  Offline = 1,
  Online = 2,
  Away = 3,
  Typing = 4,
}

export interface StatusInfo {
  status: ProfileStatus
  lastSeen: number
}
