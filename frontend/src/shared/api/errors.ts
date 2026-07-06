import axios from 'axios'
import type { ApiErrorBody } from './types'

const API_ERROR_MESSAGES: Record<string, string> = {
  DUPLICATE_SKU: 'Já existe um produto com este SKU.',
  DUPLICATE_DOCUMENT: 'Já existe um cadastro com este documento.',
  VALIDATION_ERROR: 'Verifique os dados informados.',
  INSUFFICIENT_STOCK: 'Estoque insuficiente para esta operação.',
  INVALID_STATUS: 'Esta ação não é permitida para o status atual.',
  CUSTOMER_HAS_SALES: 'Não é possível excluir um cliente com vendas vinculadas.',
  NOT_FOUND: 'Registro não encontrado.',
}

function extractApiErrorCode(error: unknown): string | undefined {
  if (!axios.isAxiosError(error)) return undefined
  const body = error.response?.data as ApiErrorBody | undefined
  return body?.error?.code
}

export function getApiErrorMessage(error: unknown, fallback: string): string {
  const code = extractApiErrorCode(error)
  if (code && API_ERROR_MESSAGES[code]) {
    return API_ERROR_MESSAGES[code]
  }
  return fallback
}
