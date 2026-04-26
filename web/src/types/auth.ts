export type Theme = 'light' | 'dark'

export interface User {
  id: number
  username: string
  status: number
  theme: Theme
  created_at: string
  updated_at: string
}

export interface LoginRequest {
  username: string
  password: string
}

export interface LoginResponse {
  token: string
  user: User
}

export interface ChangePasswordRequest {
  old_password: string
  new_password: string
}

export interface UpdateThemeRequest {
  theme: Theme
}
