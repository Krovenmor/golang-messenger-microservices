import { useCallback, useState } from 'react'
import type { FieldIssues } from '../api/errors'
import { describeFieldIssues, type FieldMessages } from './errorMessages'

export function findEmptyFields(values: Record<string, string>): FieldIssues {
  const issues: FieldIssues = {}
  for (const [field, value] of Object.entries(values)) {
    if (value.trim() === '') issues[field] = { reason: 'required' }
  }
  return issues
}

export function useFieldErrors() {
  const [fieldErrors, setFieldErrors] = useState<FieldMessages>({})

  const showIssues = useCallback((issues: FieldIssues) => {
    setFieldErrors(describeFieldIssues(issues))
  }, [])

  const setFieldError = useCallback((field: string, message: string) => {
    setFieldErrors((prev) => ({ ...prev, [field]: message }))
  }, [])

  const clearFieldError = useCallback((field: string) => {
    setFieldErrors((prev) => {
      if (!(field in prev)) return prev
      const next = { ...prev }
      delete next[field]
      return next
    })
  }, [])

  const clearAll = useCallback(() => setFieldErrors({}), [])

  return { fieldErrors, showIssues, setFieldError, clearFieldError, clearAll }
}
