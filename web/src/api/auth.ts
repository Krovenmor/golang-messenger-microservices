import { apiClient } from './client'
import type { AuthTokens } from '../types'

export interface SendCodeResult {
  expiresIn: number
  retryAfter: number
}

export async function sendCode(email: string): Promise<SendCodeResult> {
  const { data } = await apiClient.post<SendCodeResult>('/auth/send-code', { email })
  return data
}

export async function register(login: string, password: string, email: string, code: string): Promise<void> {
  await apiClient.post('/auth/register', { login, password, email, code })
}

export async function login(login_: string, password: string): Promise<AuthTokens> {
  const { data } = await apiClient.post<AuthTokens>('/auth/login', {
    login: login_,
    password,
  })
  return data
}

export interface AccountInfo {
  login: string
  email: string
}

export async function getAccount(): Promise<AccountInfo> {
  const { data } = await apiClient.get<AccountInfo>('/auth/account')
  return data
}
