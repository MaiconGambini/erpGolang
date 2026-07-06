import { describe, expect, it } from 'vitest'
import { customerFormSchema } from './schemas'

describe('customerFormSchema', () => {
  it('accepts valid customer', () => {
    const result = customerFormSchema.safeParse({
      name: 'Padaria Central',
      documentType: 'cnpj',
      document: '12345678000199',
      email: 'contato@padaria.com',
      phone: '11999999999',
      active: true,
    })
    expect(result.success).toBe(true)
  })

  it('rejects empty name', () => {
    const result = customerFormSchema.safeParse({ name: '', active: true })
    expect(result.success).toBe(false)
  })

  it('rejects invalid email', () => {
    const result = customerFormSchema.safeParse({ name: 'Test', email: 'not-email', active: true })
    expect(result.success).toBe(false)
  })
})
