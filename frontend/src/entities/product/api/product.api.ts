import { LOW_STOCK_THRESHOLD } from '@/shared/config/inventory'
import { apiClient } from '@/shared/api/client'
import type { Paginated } from '@/shared/api/types'
import type { Product } from '../model/types'

export interface ListProductsParams {
  limit?: number
  offset?: number
  search?: string
  active?: boolean
}

export async function listProducts(params: ListProductsParams = {}): Promise<Paginated<Product>> {
  const response = await apiClient.get<Paginated<Product>>('/products', { params })
  return response.data
}

export async function listLowStockProducts(threshold = LOW_STOCK_THRESHOLD, limit = 50): Promise<Product[]> {
  const response = await apiClient.get<{ data: Product[] }>('/products/low-stock', {
    params: { threshold, limit },
  })
  return response.data.data
}

export async function getProduct(id: string): Promise<Product> {
  const response = await apiClient.get<{ data: Product }>(`/products/${id}`)
  return response.data.data
}

export interface ProductInput {
  name: string
  sku: string
  price: string
  stock: number
  active: boolean
}

export async function createProduct(data: ProductInput): Promise<Product> {
  const response = await apiClient.post<{ data: Product }>('/products', data)
  return response.data.data
}

export async function updateProduct(id: string, data: ProductInput): Promise<Product> {
  const response = await apiClient.patch<{ data: Product }>(`/products/${id}`, data)
  return response.data.data
}

export async function removeProduct(id: string): Promise<void> {
  await apiClient.delete(`/products/${id}`)
}
