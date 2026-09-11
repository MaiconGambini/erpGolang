export interface Product {
  id: string
  name: string
  sku: string
  price: string
  stock: number
  unit: string
  barcode?: string
  active: boolean
  createdAt: string
  updatedAt?: string
}
