import { describe, expect, it } from 'vitest'
import { supplierFormSchema } from './schemas'

describe('supplierFormSchema', () => {
  it('accepts valid supplier', () => {
    const result = supplierFormSchema.safeParse({
      name: 'Distribuidora Norte',
      document: '12345678000199',
      email: 'contato@norte.com',
      phone: '11999999999',
      active: true,
    })
    expect(result.success).toBe(true)
  })

  it('rejects empty name', () => {
    const result = supplierFormSchema.safeParse({ name: '', active: true })
    expect(result.success).toBe(false)
  })

  it('rejects invalid email', () => {
    const result = supplierFormSchema.safeParse({ name: 'Test', email: 'not-email', active: true })
    expect(result.success).toBe(false)
  })
})
