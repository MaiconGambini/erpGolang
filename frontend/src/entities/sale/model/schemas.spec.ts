import { describe, expect, it } from 'vitest'
import { createSaleSchema } from './schemas'

describe('createSaleSchema', () => {
  it('accepts valid sale', () => {
    const result = createSaleSchema.safeParse({
      customerId: 'uuid-1',
      items: [{ productId: 'uuid-2', quantity: 2 }],
    })
    expect(result.success).toBe(true)
  })

  it('rejects empty items', () => {
    const result = createSaleSchema.safeParse({
      customerId: 'uuid-1',
      items: [],
    })
    expect(result.success).toBe(false)
  })
})
