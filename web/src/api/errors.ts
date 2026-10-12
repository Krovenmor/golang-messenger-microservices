import axios from 'axios'

export interface FieldIssue {
  reason: string
  param?: string
}

export type FieldIssues = Record<string, FieldIssue>

export interface ApiError {
  isNetworkError: boolean
  status: number | null
  message: string | null
  fields: FieldIssues
}

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === 'object' && value !== null && !Array.isArray(value)
}

function parseFields(value: unknown): FieldIssues {
  if (!isRecord(value)) return {}
  const fields: FieldIssues = {}
  for (const [name, issue] of Object.entries(value)) {
    if (isRecord(issue) && typeof issue.reason === 'string') {
      fields[name] = {
        reason: issue.reason,
        param: typeof issue.param === 'string' ? issue.param : undefined,
      }
    }
  }
  return fields
}

export function parseApiError(err: unknown): ApiError {
  if (!axios.isAxiosError(err)) {
    return { isNetworkError: false, status: null, message: null, fields: {} }
  }
  if (!err.response) {
    return { isNetworkError: true, status: null, message: null, fields: {} }
  }
  const { status, data } = err.response
  if (!isRecord(data)) {
    return { isNetworkError: false, status, message: null, fields: {} }
  }
  return {
    isNetworkError: false,
    status,
    message: typeof data.error === 'string' ? data.error : null,
    fields: parseFields(data.fields),
  }
}
