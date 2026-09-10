import { Navigate, Route, BrowserRouter, Routes } from 'react-router-dom'
import { AuthProvider, useAuth } from './context/AuthContext'
import LoginPage from './pages/Login'
import RegisterPage from './pages/Register'
import CreateProfilePage from './pages/CreateProfile'
import UnlockProfilePage from './pages/UnlockProfile'
import ChatsPage from './pages/Chats'

function LoadingScreen() {
  return (
    <div className="flex min-h-full items-center justify-center">
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

function Router() {
  return (
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
