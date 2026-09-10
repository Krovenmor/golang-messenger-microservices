import { useEffect, useMemo, useState } from 'react'
import { Sidebar } from '../components/Sidebar'
import { ChatWindow } from '../components/ChatWindow'
import { NewChatModal } from '../components/NewChatModal'
import { ProfileSettingsModal } from '../components/ProfileSettingsModal'
import { useAuth } from '../context/AuthContext'
import { chatStore } from '../api/chatStore'
import { socketManager } from '../api/ws'
import { primeProfilesFromBatch } from '../api/profileCache'
import type { ChatSummary } from '../types'

export default function ChatsPage() {
  const { profile, logout } = useAuth()
  const [chats, setChats] = useState<ChatSummary[]>(() => chatStore.getChats())
  const [activeChatId, setActiveChatId] = useState<string | null>(null)
  const [connected, setConnected] = useState(false)
  const [showNewChat, setShowNewChat] = useState(false)
  const [showSettings, setShowSettings] = useState(false)

  useEffect(() => {
    const unsub = chatStore.subscribeChats(setChats)
    if (profile) chatStore.init()
    return () => {
      unsub()
    }
  }, [profile])

  useEffect(() => socketManager.onStatus(setConnected), [])

  useEffect(() => {
    if (!profile) return
    const peerIds = chats.flatMap((c) => c.members.map((m) => m.userId)).filter((id) => id !== profile.userId)
    primeProfilesFromBatch(peerIds)
  }, [chats, profile])

  const activeChat = useMemo(() => chats.find((c) => c.chatId === activeChatId) ?? null, [chats, activeChatId])
  const activePeer = useMemo(
    () => activeChat?.members.find((m) => m.userId !== profile?.userId) ?? null,
    [activeChat, profile]
  )

  return (
    <div className="flex h-full w-full overflow-hidden">
      <div className={`h-full w-full sm:flex sm:w-[320px] sm:shrink-0 ${activeChatId ? 'hidden' : 'flex'}`}>
        <Sidebar
          chats={chats}
          activeChatId={activeChatId}
          ownUserId={profile?.userId ?? null}
          profile={profile}
          connected={connected}
          onSelect={setActiveChatId}
          onNewChat={() => setShowNewChat(true)}
          onLogout={logout}
          onOpenSettings={() => setShowSettings(true)}
        />
      </div>

      <div className={`h-full flex-1 sm:flex ${activeChatId ? 'flex' : 'hidden'}`}>
        {activeChatId && profile ? (
          <ChatWindow
            key={activeChatId}
            chatId={activeChatId}
            ownUserId={profile.userId}
            ownPubKey={profile.pubKey}
            peer={activePeer}
            onBack={() => setActiveChatId(null)}
          />
        ) : (
          <div className="hidden flex-1 flex-col items-center justify-center gap-3 text-center sm:flex">
            <span className="relative flex h-3 w-3">
              <span className="absolute h-full w-full rounded-full bg-signal animate-pulse-signal" />
            </span>
            <p className="font-display text-lg text-ink">Выберите диалог</p>
            <p className="max-w-xs text-sm text-faint">
              Или начните новый — нажмите «+» рядом с поиском в списке слева.
            </p>
          </div>
        )}
      </div>

      {showSettings && <ProfileSettingsModal onClose={() => setShowSettings(false)} />}

      {showNewChat && (
        <NewChatModal
          onClose={() => setShowNewChat(false)}
          onCreated={(chatId, peer) => {
            if (!profile) return
            setShowNewChat(false)
            chatStore.addLocalChat(chatId, { userId: profile.userId }, peer)
            setActiveChatId(chatId)
          }}
        />
      )}
    </div>
  )
}
