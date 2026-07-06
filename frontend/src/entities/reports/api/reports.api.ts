import { apiClient } from '@/shared/api/client'

export interface SalesByDayPoint {
  date: string
  count: number
  total: string
}

export interface TopProductPoint {
  productName: string
  productSku: string
  quantity: number
  revenue: string
}

export interface ReportsRangeParams {
  from?: string
  to?: string
}

export async function getSalesByDay(params: ReportsRangeParams = {}): Promise<SalesByDayPoint[]> {
  const response = await apiClient.get<{ data: SalesByDayPoint[] }>('/reports/sales-by-day', { params })
  return response.data.data
}

export async function getTopProducts(
  params: ReportsRangeParams & { limit?: number } = {},
): Promise<TopProductPoint[]> {
  const response = await apiClient.get<{ data: TopProductPoint[] }>('/reports/top-products', { params })
  return response.data.data
}

export async function downloadSalePdf(id: string): Promise<Blob> {
  const response = await apiClient.get(`/reports/sales/${id}/pdf`, { responseType: 'blob' })
  return response.data as Blob
}

export async function downloadSalesSummaryPdf(from: string, to: string): Promise<Blob> {
  const response = await apiClient.get('/reports/sales-summary.pdf', {
    params: { from, to },
    responseType: 'blob',
  })
  return response.data as Blob
}
