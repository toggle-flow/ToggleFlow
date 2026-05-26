import { api } from './client'

export interface FlagEnvState {
  environment_id: number
  environment_name: string
  environment_slug: string
  enabled: boolean
}

export interface Flag {
  id: number
  project_id: number
  key: string
  name: string
  description: string
  created_at: string
  updated_at: string
  environments: FlagEnvState[]
}

export const flagsApi = {
  list: (projectId: number) => api.get<Flag[]>(`/projects/${projectId}/flags`),
  create: (projectId: number, data: { name: string; key: string; description?: string }) =>
    api.post<Flag>(`/projects/${projectId}/flags`, data),
  toggle: (projectId: number, flagKey: string, environmentId: number, enabled: boolean) =>
    api.patch<{ ok: boolean }>(`/projects/${projectId}/flags/${flagKey}`, { environment_id: environmentId, enabled }),
  delete: (projectId: number, flagKey: string) =>
    api.delete<void>(`/projects/${projectId}/flags/${flagKey}`),
}
