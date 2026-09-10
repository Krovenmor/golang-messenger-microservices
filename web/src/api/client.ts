import axios, { AxiosError, type InternalAxiosRequestConfig } from 'axios'
import { tokenStorage } from '../utils/tokenStorage'
import { isJwtExpiringSoon } from '../utils/jwt'
import { API_BASE_URL, ensureFreshToken } from './authRefresh'

export { API_BASE_URL }

const BEARER_PREFIX = 'Bearer '

export const apiClient = axios.create({
  baseURL: API_BASE_URL,
})

function bearerTokenOf(header: unknown): string | undefined {
  if (typeof header !== 'string' || !header.startsWith(BEARER_PREFIX)) return undefined
  return header.slice(BEARER_PREFIX.length)
}

apiClient.interceptors.request.use(async (config) => {
  const tokens = tokenStorage.get()
  if (!tokens?.accessToken) return config

  const accessToken = isJwtExpiringSoon(tokens.accessToken)
    ? (await ensureFreshToken()) ?? tokens.accessToken
    : tokens.accessToken

  config.headers = config.headers ?? {}
  config.headers.Authorization = `${BEARER_PREFIX}${accessToken}`
  return config
})

apiClient.interceptors.response.use(
  (response) => response,
  async (error: AxiosError) => {
    const originalRequest = error.config as (InternalAxiosRequestConfig & { _retry?: boolean }) | undefined

    if (error.response?.status !== 401 || !originalRequest || originalRequest._retry) {
      return Promise.reject(error)
    }

    const newToken = await ensureFreshToken(bearerTokenOf(originalRequest.headers?.Authorization))
    if (!newToken) return Promise.reject(error)

    originalRequest._retry = true
    originalRequest.headers = originalRequest.headers ?? {}
    originalRequest.headers.Authorization = `${BEARER_PREFIX}${newToken}`
    return apiClient(originalRequest)
  }
)
