import { LOW_STOCK_THRESHOLD } from '@/shared/config/inventory'
import { apiClient } from '@/shared/api/client'
import type { DashboardSummary } from '../model/types'

interface DashboardSummaryApi {
  activeCustomers: number
  newCustomers30d: number
  draftSales: number
  lowStockAlerts: number
  confirmedSalesCount?: number
  confirmedSalesTotal?: string | number
}

function mapDashboardSummary(raw: DashboardSummaryApi): DashboardSummary {
  return {
    activeCustomers: raw.activeCustomers,
    newCustomers30d: raw.newCustomers30d,
    draftSales: raw.draftSales,
    lowStockAlerts: raw.lowStockAlerts,
    confirmedSalesCount: raw.confirmedSalesCount ?? 0,
    confirmedSalesTotal: String(raw.confirmedSalesTotal ?? '0'),
  }
}

export async function getDashboardSummary(threshold = LOW_STOCK_THRESHOLD): Promise<DashboardSummary> {
  const response = await apiClient.get<{ data: DashboardSummaryApi }>('/dashboard/summary', {
    params: { threshold },
  })
  return mapDashboardSummary(response.data.data)
}
