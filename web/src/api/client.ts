import axios, { AxiosError, type InternalAxiosRequestConfig } from 'axios'
import { tokenStorage } from '../utils/tokenStorage'
import { API_BASE_URL, ensureFreshToken } from './authRefresh'

export { API_BASE_URL }

export const apiClient = axios.create({
  baseURL: API_BASE_URL,
})

apiClient.interceptors.request.use((config) => {
  const tokens = tokenStorage.get()
  if (tokens?.accessToken) {
    config.headers = config.headers ?? {}
    config.headers.Authorization = `Bearer ${tokens.accessToken}`
  }
  return config
})

apiClient.interceptors.response.use(
  (response) => response,
  async (error: AxiosError) => {
    const originalRequest = error.config as (InternalAxiosRequestConfig & { _retry?: boolean }) | undefined

    if (error.response?.status !== 401 || !originalRequest || originalRequest._retry) {
      return Promise.reject(error)
    }

    const newToken = await ensureFreshToken()
    if (!newToken) return Promise.reject(error)

    originalRequest._retry = true
    originalRequest.headers = originalRequest.headers ?? {}
    originalRequest.headers.Authorization = `Bearer ${newToken}`
    return apiClient(originalRequest)
  }
)
