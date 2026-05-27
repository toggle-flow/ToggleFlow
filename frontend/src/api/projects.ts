import { api, type PageResult, type PageParams } from './client'
import { useToastStore } from '@/stores/toast'

export interface Project {
  id: number
  name: string
  key: string
  description: string
  created_by: number | null
  created_by_name: string
  created_at: string
  updated_at: string
}

export const projectsApi = {
  list: (params?: PageParams) => api.get<PageResult<Project>>('/projects', params),
  create: (name: string, key: string, description: string) =>
    api.post<Project>('/projects', { name, key, description }).then((r) => {
      useToastStore().show('created project')
      return r
    }),
  update: (id: number, name: string, key: string, description: string) =>
    api.patch<Project>(`/projects/${id}`, { name, key, description }).then((r) => {
      useToastStore().show('updated project')
      return r
    }),
  delete: (id: number) =>
    api.delete<void>(`/projects/${id}`).then((r) => {
      useToastStore().show('deleted project')
      return r
    }),
}
