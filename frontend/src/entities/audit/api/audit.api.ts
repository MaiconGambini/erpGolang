import { apiClient } from '@/shared/api/client'
import type { Paginated } from '@/shared/api/types'
import type { AuditLog } from '../model/types'

export interface ListAuditLogsParams {
  limit?: number
  offset?: number
}

export async function listAuditLogs(params: ListAuditLogsParams = {}): Promise<Paginated<AuditLog>> {
  const response = await apiClient.get<Paginated<AuditLog>>('/audit-logs', { params })
  return response.data
}
