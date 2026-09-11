import { AxiosError, type AxiosResponse } from 'axios'
import { describe, expect, it } from 'vitest'
import type { ApiErrorBody } from './types'
import { getApiErrorMessage } from './errors'

function apiError(code: string): AxiosError<ApiErrorBody> {
  return new AxiosError('request failed', 'ERR_BAD_REQUEST', undefined, undefined, {
    data: { error: { code, message: 'server message' } },
    status: 400,
    statusText: 'Bad Request',
    headers: {},
    config: {} as AxiosResponse['config'],
  })
}

describe('getApiErrorMessage', () => {
  it('maps known API error codes to Portuguese messages', () => {
    expect(getApiErrorMessage(apiError('DUPLICATE_SKU'), 'fallback')).toBe(
      'Já existe um produto com este SKU.',
    )
    expect(getApiErrorMessage(apiError('INSUFFICIENT_STOCK'), 'fallback')).toBe(
      'Estoque insuficiente para esta operação.',
    )
  })

  it('returns fallback for unknown errors', () => {
    expect(getApiErrorMessage(new Error('network'), 'Erro genérico')).toBe('Erro genérico')
    expect(getApiErrorMessage(apiError('INTERNAL_ERROR'), 'Erro genérico')).toBe('Erro genérico')
  })
})
