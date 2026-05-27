import { api } from './client'
import { useToastStore } from '@/stores/toast'

export type Role = 'superuser' | 'admin' | 'owner' | 'editor' | 'viewer'

export interface User {
  id: number
  uuid: string
  name: string
  email: string
  role: Role
  locale: string
  activated_at?: string
  created_by?: number
  created_at: string
  updated_at: string
}

export interface InviteResult {
  user: User
  welcome_token: string
}

export interface ResetResult {
  user: User
  reset_token: string
}

export const usersApi = {
  list: () => api.get<User[]>('/users'),
  create: (name: string, email: string, role: Role, expiryDays: number) =>
    api.post<InviteResult>('/users', { name, email, role, expiry_days: expiryDays }).then((r) => {
      useToastStore().show('invited user')
      return r
    }),
  update: (id: number, data: Partial<{ name: string; email: string; role: Role }>) =>
    api.patch<User>(`/users/${id}`, data).then((r) => {
      useToastStore().show('updated user')
      return r
    }),
  reinvite: (id: number, expiryDays: number) =>
    api.post<InviteResult>(`/users/${id}/reinvite`, { expiry_days: expiryDays }).then((r) => {
      useToastStore().show('regenerated invite')
      return r
    }),
  delete: (id: number) =>
    api.delete<void>(`/users/${id}`).then((r) => {
      useToastStore().show('removed user')
      return r
    }),
  generateResetLink: (id: number, expiryDays = 1) =>
    api.post<ResetResult>(`/users/${id}/reset-link`, { expiry_days: expiryDays }).then((r) => {
      useToastStore().show('generated reset link')
      return r
    }),
}
