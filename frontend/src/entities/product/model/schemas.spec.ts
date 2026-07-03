import { describe, expect, it } from 'vitest'
import { productFormSchema } from './schemas'

describe('productFormSchema', () => {
  it('accepts valid product', () => {
    const result = productFormSchema.safeParse({
      name: 'Camiseta',
      sku: 'CAM-001',
      price: '29.90',
      stock: 10,
      active: true,
    })
    expect(result.success).toBe(true)
  })

  it('rejects empty sku', () => {
    const result = productFormSchema.safeParse({
      name: 'Test',
      sku: '',
      price: '10.00',
      stock: 0,
      active: true,
    })
    expect(result.success).toBe(false)
  })

  it('rejects non-positive price', () => {
    const result = productFormSchema.safeParse({
      name: 'Test',
      sku: 'SKU-1',
      price: '0',
      stock: 0,
      active: true,
    })
    expect(result.success).toBe(false)
  })

  it('rejects negative stock', () => {
    const result = productFormSchema.safeParse({
      name: 'Test',
      sku: 'SKU-1',
      price: '10.00',
      stock: -1,
      active: true,
    })
    expect(result.success).toBe(false)
  })
})
