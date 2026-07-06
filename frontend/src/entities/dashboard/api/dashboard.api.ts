import { LOW_STOCK_THRESHOLD } from '@/shared/config/inventory'
import { apiClient } from '@/shared/api/client'
import type { DashboardSummary } from '../model/types'

export async function getDashboardSummary(threshold = LOW_STOCK_THRESHOLD): Promise<DashboardSummary> {
  const response = await apiClient.get<{ data: DashboardSummary }>('/dashboard/summary', {
    params: { threshold },
  })
  return response.data.data
}
