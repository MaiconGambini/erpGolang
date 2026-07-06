import { LOW_STOCK_THRESHOLD } from '@/shared/config/inventory'
import { apiClient } from '@/shared/api/client'
import type { DashboardSummary } from '../model/types'

interface DashboardSummaryApi {
  active_customers: number
  new_customers_30d: number
  draft_sales: number
  low_stock_alerts: number
  confirmed_sales_count?: number
  confirmed_sales_total?: string | number
}

function mapDashboardSummary(raw: DashboardSummaryApi): DashboardSummary {
  return {
    activeCustomers: raw.active_customers,
    newCustomers30d: raw.new_customers_30d,
    draftSales: raw.draft_sales,
    lowStockAlerts: raw.low_stock_alerts,
    confirmedSalesCount: raw.confirmed_sales_count ?? 0,
    confirmedSalesTotal: String(raw.confirmed_sales_total ?? '0'),
  }
}

export async function getDashboardSummary(threshold = LOW_STOCK_THRESHOLD): Promise<DashboardSummary> {
  const response = await apiClient.get<{ data: DashboardSummaryApi }>('/dashboard/summary', {
    params: { threshold },
  })
  return mapDashboardSummary(response.data.data)
}
