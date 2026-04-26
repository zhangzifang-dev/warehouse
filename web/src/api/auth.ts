import api from './client'
import type { LoginRequest, LoginResponse, ChangePasswordRequest, UpdateThemeRequest, User } from '../types/auth'

export const authApi = {
  login: async (data: LoginRequest): Promise<LoginResponse> => {
    const response = await api.post<LoginResponse>('/auth/login', data)
    return response.data
  },

  getProfile: async (): Promise<User> => {
    const response = await api.get<User>('/auth/profile')
    return response.data
  },

  changePassword: async (data: ChangePasswordRequest): Promise<void> => {
    await api.put('/auth/password', data)
  },

  updateTheme: async (data: UpdateThemeRequest): Promise<void> => {
    await api.put('/auth/theme', data)
  },
}
