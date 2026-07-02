import { z } from 'zod'

export const createCustomerSchema = z.object({
  name: z.string().min(1),
  document: z.string().optional(),
  email: z.string().email().optional(),
  phone: z.string().optional(),
  active: z.boolean().default(true),
})
