import type { DocumentType } from '@/shared/lib/document'

export interface PartyAddressFields {
  documentType?: DocumentType
  document?: string
  postalCode?: string
  street?: string
  streetNumber?: string
  city?: string
  state?: string
}

export interface Customer extends PartyAddressFields {
  id: string
  name: string
  email?: string
  phone?: string
  active: boolean
  createdAt: string
  updatedAt?: string
}
