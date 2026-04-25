export interface User {
  id: number
  username: string
  nickname: string
  email: string
  phone: string
  status: number
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
