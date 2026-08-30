export interface AuthTokens {
  accessToken: string
  refreshToken: string
}

/** Grows over time as new profile features ship — treat every field as
 * optional/absent rather than assuming today's shape is final. */
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
  additional?: ProfileAdditional
}

export interface PublicProfile {
  userId: string
  name: string
  userName: string
  pubKey: string
  createdAt: string
  additional?: ProfileAdditional
}

export interface MediaStorageInfo {
  /** Bytes. */
  maxSpace: number
  /** Bytes. */
  spaceFilled: number
  savedFiles: number
}

export type MediaFileType = 'photo'
export type MediaFileSubType = 'avatar'

export interface MediaFileInfo {
  mediaId: string
  type: MediaFileType
  subType: MediaFileSubType
  /** Bytes. */
  size: number
  isPublic: boolean
  addedAt: string
}

export interface Reaction {
  emoji: string
  users: string[]
}

/** The wire shape from the backend — message/senderKey/receiverKey are all
 * still ciphertext at this point. Use decryptMessageText (messageCrypto.ts)
 * to get the plaintext for display. */
export interface Message {
  messageId: string
  senderId: string
  /** Base64 AES-256-GCM ciphertext. '' when deletedAt is set. */
  message: string
  /** Base64, 105 bytes: ephemeral pubkey (65B) + AES-KW-wrapped MessageKey (40B), wrapped for the SENDER's static key. */
  senderKey: string
  /** Same shape as senderKey, but wrapped for the RECEIVER's static key. */
  receiverKey: string
  /** Base64, 12-byte AES-GCM nonce used for `message`. */
  nonce: string
  createdAt: string
  /** Set the moment the message text was last edited, otherwise null. */
  redactedAt: string | null
  /** Set the moment the message was deleted, otherwise null. `message` reads as '' once this is set. */
  deletedAt: string | null
  /** messageId this is a reply to, if any. */
  replyToId: string | null
  /** Absent specifically on chats/full's `lastMessage` — treat missing as empty everywhere. */
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

/** Anything the server can push down the socket. `response` is the ack for
 * a request WE sent (status post / track) and isn't dispatched to general
 * event listeners — see SocketManager. */
export type WSEvent =
  | WSNewChatEvent
  | WSNewMessageEvent
  | WSMessageRedactedEvent
  | WSMessageDeletedEvent
  | WSNewReactionEvent
  | WSDelReactionEvent
  | WSStatusEvent
  | WSResponseEvent

/** The two request kinds the client can send over the socket, both wrapped
 * as `{ req, payload }`. */
export type WSRequest =
  | { req: 1; payload: { newStatus: ProfileStatus } } // post your own status
  | { req: 2; payload: { userId: string } } // track someone else's status

export enum ProfileStatus {
  Offline = 1,
  Online = 2,
  Away = 3,
  Typing = 4,
}

export interface StatusInfo {
  status: ProfileStatus
  /** Unix seconds. 0 means the profile doesn't exist or was last online a very long time ago. */
  lastSeen: number
}
