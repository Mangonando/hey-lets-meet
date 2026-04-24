import { createContext, useEffect, useMemo, useState } from 'react'
import { api, type User } from './api'

type AuthState =
  | { status: 'loading'; user: null }
  | { status: 'authed'; user: User }
  | { status: 'anon'; user: null }

type AuthContextValue = {
  state: AuthState
  refresh: () => Promise<void>
  logout: () => Promise<void>
}

const AuthContext = createContext<AuthContextValue | null>(null)
export { AuthContext }

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [state, setState] = useState<AuthState>({
    status: 'loading',
    user: null,
  })

  async function refresh() {
    try {
      const user = await api<User>('/auth/me')
      setState({ status: 'authed', user })
    } catch {
      setState({ status: 'anon', user: null })
    }
  }

  async function logout() {
    try {
      await api<{ message: string }>('/auth/logout', {
        method: 'POST',
        body: JSON.stringify({}),
      })
    } finally {
      setState({ status: 'anon', user: null })
    }
  }

  useEffect(() => {
    void refresh()
  }, [])

  const value = useMemo<AuthContextValue>(
    () => ({ state, refresh, logout }),
    [state],
  )

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}
