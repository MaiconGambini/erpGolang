import type { PartyFormValues } from '@/shared/lib/party-schema'
import type { CustomerInput } from '@/entities/customer/api/customer.api'

export function toPartyInput(data: PartyFormValues): CustomerInput {
  const documentType =
    data.documentType === 'cpf' || data.documentType === 'cnpj' ? data.documentType : undefined
  return {
    name: data.name,
    documentType,
    document: data.document || undefined,
    email: data.email || undefined,
    phone: data.phone || undefined,
    postalCode: data.postalCode || undefined,
    street: data.street || undefined,
    streetNumber: data.streetNumber || undefined,
    city: data.city || undefined,
    state: data.state ? data.state.toUpperCase() : undefined,
    active: data.active,
  }
}
