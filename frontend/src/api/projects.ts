import { api } from './client'

export interface Project {
  id: number
  name: string
  slug: string
  created_at: string
  updated_at: string
}

export const projectsApi = {
  list: () => api.get<Project[]>('/projects'),
  create: (name: string, slug: string) => api.post<Project>('/projects', { name, slug }),
}
