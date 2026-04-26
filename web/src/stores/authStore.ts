import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import type { User, Theme } from '../types/auth'
import { authApi } from '../api/auth'

interface AuthState {
  user: User | null
  token: string | null
  isAuthenticated: boolean
  theme: Theme
  login: (token: string, user: User) => void
  logout: () => void
  setUser: (user: User) => void
  setTheme: (theme: Theme) => void
  updateTheme: (theme: Theme) => Promise<void>
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      user: null,
      token: null,
      isAuthenticated: false,
      theme: 'light',
      login: (token: string, user: User) => {
        const theme = user.theme || 'light'
        document.documentElement.setAttribute('data-theme', theme)
        set({ token, user, isAuthenticated: true, theme })
      },
      logout: () => {
        document.documentElement.setAttribute('data-theme', 'light')
        set({ token: null, user: null, isAuthenticated: false, theme: 'light' })
      },
      setUser: (user: User) => {
        set({ user })
      },
      setTheme: (theme: Theme) => {
        document.documentElement.setAttribute('data-theme', theme)
        set({ theme })
      },
      updateTheme: async (theme: Theme) => {
        await authApi.updateTheme({ theme })
        document.documentElement.setAttribute('data-theme', theme)
        set({ theme, user: useAuthStore.getState().user ? { ...useAuthStore.getState().user!, theme } : null })
      },
    }),
    {
      name: 'auth-storage',
    }
  )
)
