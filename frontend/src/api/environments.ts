import { api, type PageResult, type PageParams } from './client'
import { useToastStore } from '@/stores/toast'

export interface Environment {
  id: number
  project_id: number
  name: string
  key: string
  description: string
  protected: boolean
  created_at: string
  updated_at: string
}

export interface SDKKeyRecord {
  id: number
  environment_id: number
  label: string
  key_prefix: string
  expires_at?: string
  created_at: string
}

export interface APIKeyRecord {
  id: number
  project_id: number
  label: string
  key_prefix: string
  expires_at?: string
  created_at: string
}

export const environmentsApi = {
  list: (projectId: number, params?: PageParams) =>
    api.get<PageResult<Environment>>(`/projects/${projectId}/environments`, params),
  create: (
    projectId: number,
    name: string,
    key: string,
    description: string,
    isProtected: boolean
  ) =>
    api
      .post<Environment>(`/projects/${projectId}/environments`, {
        name,
        key,
        description,
        protected: isProtected,
      })
      .then((r) => {
        useToastStore().show('created environment')
        return r
      }),
  update: (
    projectId: number,
    envId: number,
    name: string,
    key: string,
    description: string,
    isProtected: boolean
  ) =>
    api
      .patch<Environment>(`/projects/${projectId}/environments/${envId}`, {
        name,
        key,
        description,
        protected: isProtected,
      })
      .then((r) => {
        useToastStore().show('updated environment')
        return r
      }),
  delete: (projectId: number, envId: number) =>
    api.delete<void>(`/projects/${projectId}/environments/${envId}`).then((r) => {
      useToastStore().show('deleted environment')
      return r
    }),
  sdkKeys: {
    list: (projectId: number, envId: number) =>
      api.get<SDKKeyRecord[]>(`/projects/${projectId}/environments/${envId}/sdk-keys`),
    create: (projectId: number, envId: number, label: string, expiresAt: string | null) =>
      api
        .post<{
          key: string
          record: SDKKeyRecord
        }>(`/projects/${projectId}/environments/${envId}/sdk-keys`, {
          label,
          expires_at: expiresAt || undefined,
        })
        .then((r) => {
          useToastStore().show('created sdk key')
          return r
        }),
    delete: (projectId: number, envId: number, keyId: number) =>
      api
        .delete<void>(`/projects/${projectId}/environments/${envId}/sdk-keys/${keyId}`)
        .then((r) => {
          useToastStore().show('deleted sdk key')
          return r
        }),
  },
  apiKeys: {
    list: (projectId: number) => api.get<APIKeyRecord[]>(`/projects/${projectId}/api-keys`),
    create: (projectId: number, label: string, expiresAt: string | null) =>
      api
        .post<{ key: string; record: APIKeyRecord }>(`/projects/${projectId}/api-keys`, {
          label,
          expires_at: expiresAt || undefined,
        })
        .then((r) => {
          useToastStore().show('created api key')
          return r
        }),
    delete: (projectId: number, keyId: number) =>
      api.delete<void>(`/projects/${projectId}/api-keys/${keyId}`).then((r) => {
        useToastStore().show('deleted api key')
        return r
      }),
  },
}
