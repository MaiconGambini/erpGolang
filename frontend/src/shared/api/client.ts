import axios from 'axios'
import { env } from '@/shared/config/env'
import { useSessionStore } from '@/entities/session/model/session.store'

let refreshPromise: Promise<string | null> | null = null

export const apiClient = axios.create({
  baseURL: env.apiBaseUrl,
  withCredentials: true,
})

apiClient.interceptors.request.use((config) => {
  const session = useSessionStore()
  if (session.accessToken) {
    config.headers.Authorization = `Bearer ${session.accessToken}`
  }
  return config
})

apiClient.interceptors.response.use(
  (response) => response,
  async (error) => {
    const original = error.config
    if (!original || original._retry || original.url?.includes('/auth/refresh')) {
      return Promise.reject(error)
    }
    if (error.response?.status !== 401) {
      return Promise.reject(error)
    }
    original._retry = true
    if (!refreshPromise) {
      refreshPromise = apiClient
        .post('/auth/refresh')
        .then((res) => {
          const session = useSessionStore()
          const token = res.data?.data?.accessToken as string
          const user = res.data?.data?.user
          session.setAccess(token)
          if (user) session.setUser(user)
          return token
        })
        .catch(() => {
          useSessionStore().clear()
          return null
        })
        .finally(() => {
          refreshPromise = null
        })
    }
    const token = await refreshPromise
    if (!token) {
      return Promise.reject(error)
    }
    original.headers.Authorization = `Bearer ${token}`
    return apiClient(original)
  },
)
