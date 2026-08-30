import { Navigate, Route, BrowserRouter, Routes } from 'react-router-dom'
import { AuthProvider, useAuth } from './context/AuthContext'
import LoginPage from './pages/Login'
import RegisterPage from './pages/Register'
import CreateProfilePage from './pages/CreateProfile'
import UnlockProfilePage from './pages/UnlockProfile'
import ChatsPage from './pages/Chats'
import { ensureNotificationPermission, handleSocketToast, subscribeToasts, type ToastItem } from './utils/notifications'
import { socketManager } from './api/ws'
import { chatStore } from './api/chatStore'
import { useEffect, useState } from 'react'

function LoadingScreen() {
  return (
    <div className="flex min-h-screen items-center justify-center">
      <span className="relative flex h-3 w-3">
        <span className="absolute h-full w-full rounded-full bg-signal animate-pulse-signal" />
      </span>
    </div>
  )
}

function RequireAuth({ children }: { children: JSX.Element }) {
  const { isAuthenticated, loading, hasProfile, cryptoReady } = useAuth()
  if (loading) return <LoadingScreen />
  if (!isAuthenticated) return <Navigate to="/login" replace />
  if (hasProfile === false) return <Navigate to="/setup" replace />
  // A profile can exist without the private key being unlocked in this tab
  // yet (e.g. right after a fresh login) — messages can't be encrypted or
  // decrypted until that happens.
  if (hasProfile === true && !cryptoReady) return <Navigate to="/unlock" replace />
  return children
}

function RequireSetup({ children }: { children: JSX.Element }) {
  const { isAuthenticated, loading } = useAuth()
  if (loading) return <LoadingScreen />
  if (!isAuthenticated) return <Navigate to="/login" replace />
  return children
}

function RequireUnlock({ children }: { children: JSX.Element }) {
  const { isAuthenticated, loading, hasProfile, cryptoReady } = useAuth()
  if (loading) return <LoadingScreen />
  if (!isAuthenticated) return <Navigate to="/login" replace />
  if (hasProfile === false) return <Navigate to="/setup" replace />
  if (cryptoReady) return <Navigate to="/" replace />
  return children
}

function RedirectIfAuthed({ children }: { children: JSX.Element }) {
  const { isAuthenticated, loading } = useAuth()
  if (loading) return <LoadingScreen />
  if (isAuthenticated) return <Navigate to="/" replace />
  return children
}

function ToastView() {
  const [items, setItems] = useState<ToastItem[]>([])
  const [chats, setChats] = useState(() => chatStore.getChats())
  const { profile } = useAuth()

  useEffect(() => {
    const unsub = subscribeToasts(setItems)
    return unsub
  }, [])

  useEffect(() => {
    const unsub = chatStore.subscribeChats(setChats)
    return unsub
  }, [])

  useEffect(() => {
    const enableNotifications = () => {
      ensureNotificationPermission().catch(() => {})
      window.removeEventListener('pointerdown', enableNotifications)
      window.removeEventListener('keydown', enableNotifications)
    }

    window.addEventListener('pointerdown', enableNotifications, { passive: true })
    window.addEventListener('keydown', enableNotifications, { passive: true })

    return () => {
      window.removeEventListener('pointerdown', enableNotifications)
      window.removeEventListener('keydown', enableNotifications)
    }
  }, [])

  useEffect(() => {
    if (!profile) return

    const unsub = socketManager.onEvent((event) => {
      const peerNameForChat = (chatId: string) => {
        const chat = chats.find((c) => c.chatId === chatId)
        const peer = chat?.members.find((m) => m.userId !== profile.userId)
        return peer?.name ?? 'Собеседник'
      }

      handleSocketToast(event, profile.userId, peerNameForChat)
    })

    return unsub
  }, [profile, chats])

  if (items.length === 0) return null

  return (
    <div className="pointer-events-none fixed right-4 top-4 z-50 flex max-w-sm flex-col gap-2">
      {items.map((toast) => (
        <div
          key={toast.id}
          className={`rounded-xl border px-3 py-2 shadow-lg backdrop-blur-sm ${
            toast.kind === 'success'
              ? 'border-emerald-500/40 bg-emerald-500/10 text-emerald-100'
              : toast.kind === 'warn'
                ? 'border-amber-500/40 bg-amber-500/10 text-amber-100'
                : 'border-sky-500/40 bg-sky-500/10 text-sky-100'
          }`}
        >
          <p className="text-xs font-medium uppercase tracking-[0.2em] opacity-80">{toast.title}</p>
          {toast.message && <p className="mt-1 text-sm">{toast.message}</p>}
        </div>
      ))}
    </div>
  )
}

function Router() {
  return (
    <>
      <Routes>
        <Route
          path="/login"
          element={
            <RedirectIfAuthed>
              <LoginPage />
            </RedirectIfAuthed>
          }
        />
        <Route
          path="/register"
          element={
            <RedirectIfAuthed>
              <RegisterPage />
            </RedirectIfAuthed>
          }
        />
        <Route
          path="/setup"
          element={
            <RequireSetup>
              <CreateProfilePage />
            </RequireSetup>
          }
        />
        <Route
          path="/unlock"
          element={
            <RequireUnlock>
              <UnlockProfilePage />
            </RequireUnlock>
          }
        />
        <Route
          path="/"
          element={
            <RequireAuth>
              <ChatsPage />
            </RequireAuth>
          }
        />
        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
      <ToastView />
    </>
  )
}

export default function App() {
  return (
    <BrowserRouter>
      <AuthProvider>
        <Router />
      </AuthProvider>
    </BrowserRouter>
  )
}
