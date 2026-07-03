import { apiClient } from '@/shared/api/client'
import type { Paginated } from '@/shared/api/types'
import type { Customer } from '../model/types'

export interface ListCustomersParams {
  limit?: number
  offset?: number
  search?: string
  active?: boolean
}

export async function listCustomers(params: ListCustomersParams = {}): Promise<Paginated<Customer>> {
  const response = await apiClient.get<Paginated<Customer>>('/customers', { params })
  return response.data
}

export async function getCustomer(id: string): Promise<Customer> {
  const response = await apiClient.get<{ data: Customer }>(`/customers/${id}`)
  return response.data.data
}

export interface CustomerInput {
  name: string
  document?: string
  email?: string
  phone?: string
  active: boolean
}

export async function createCustomer(data: CustomerInput): Promise<Customer> {
  const response = await apiClient.post<{ data: Customer }>('/customers', data)
  return response.data.data
}

export async function updateCustomer(id: string, data: CustomerInput): Promise<Customer> {
  const response = await apiClient.patch<{ data: Customer }>(`/customers/${id}`, data)
  return response.data.data
}

export async function removeCustomer(id: string): Promise<void> {
  await apiClient.delete(`/customers/${id}`)
}
