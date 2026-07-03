export interface Customer {
  id: string
  name: string
  document?: string
  email?: string
  phone?: string
  active: boolean
  createdAt: string
  updatedAt?: string
}
