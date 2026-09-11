export type SaleStatus = 'draft' | 'confirmed' | 'cancelled'

export interface SaleItem {
  id: string
  productId: string
  productName: string
  productSku: string
  quantity: number
  unitPrice: string
  lineTotal: string
}

export interface Sale {
  id: string
  customerId: string
  customerName: string
  status: SaleStatus
  total: string
  notes?: string
  items?: SaleItem[]
  createdAt: string
  updatedAt?: string
}
