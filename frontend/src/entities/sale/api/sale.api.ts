import { apiClient } from '@/shared/api/client'
import type { Paginated } from '@/shared/api/types'
import type { Sale, SaleStatus } from '../model/types'

export interface ListSalesParams {
  limit?: number
  offset?: number
  search?: string
  status?: SaleStatus
  from?: string
  to?: string
}

export interface SaleItemInput {
  productId: string
  quantity: number
}

export interface SaleInput {
  customerId: string
  notes?: string
  items: SaleItemInput[]
}

export async function listSales(params: ListSalesParams = {}): Promise<Paginated<Sale>> {
  const response = await apiClient.get<Paginated<Sale>>('/sales', { params })
  return response.data
}

export async function getSale(id: string): Promise<Sale> {
  const response = await apiClient.get<{ data: Sale }>(`/sales/${id}`)
  return response.data.data
}

export async function createSale(data: SaleInput): Promise<Sale> {
  const response = await apiClient.post<{ data: Sale }>('/sales', data)
  return response.data.data
}

export async function updateSale(id: string, data: SaleInput): Promise<Sale> {
  const response = await apiClient.patch<{ data: Sale }>(`/sales/${id}`, data)
  return response.data.data
}

export async function confirmSale(id: string): Promise<Sale> {
  const response = await apiClient.post<{ data: Sale }>(`/sales/${id}/confirm`)
  return response.data.data
}

export async function cancelSale(id: string): Promise<Sale> {
  const response = await apiClient.post<{ data: Sale }>(`/sales/${id}/cancel`)
  return response.data.data
}

export async function removeSale(id: string): Promise<void> {
  await apiClient.delete(`/sales/${id}`)
}
