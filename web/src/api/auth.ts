import { apiClient } from './client'
import type { AuthTokens } from '../types'

export interface SendCodeResult {
  expiresIn: number
  retryAfter: number
}

export interface StartRegistrationResult extends SendCodeResult {
  registrationId: string
}

export async function startRegistration(login: string, password: string, email: string): Promise<StartRegistrationResult> {
  const { data } = await apiClient.post<StartRegistrationResult>('/auth/register', { login, password, email })
  return data
}

export async function resendCode(registrationId: string): Promise<SendCodeResult> {
  const { data } = await apiClient.post<SendCodeResult>('/auth/register/resend', { registrationId })
  return data
}

export async function confirmRegistration(registrationId: string, code: string): Promise<void> {
  await apiClient.post('/auth/register/confirm', { registrationId, code })
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
