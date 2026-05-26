import { api } from './client'

export interface Member {
  user_id: number
  name: string
  email: string
  role: string
  joined_at: string
}

export const membersApi = {
  list: (projectId: number) => api.get<Member[]>(`/projects/${projectId}/members`),
  add: (projectId: number, userId: number) =>
    api.post<Member>(`/projects/${projectId}/members`, { user_id: userId }),
  remove: (projectId: number, userId: number) =>
    api.delete<void>(`/projects/${projectId}/members/${userId}`),
}
