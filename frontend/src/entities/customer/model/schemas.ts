import { z } from 'zod'

const optionalEmail = z.union([z.string().email(), z.literal('')]).optional()
const optionalText = z.string().optional()

export const customerFormSchema = z.object({
  name: z.string().min(1, 'Nome é obrigatório'),
  document: optionalText,
  email: optionalEmail,
  phone: optionalText,
  active: z.boolean().default(true),
})

export type CustomerFormValues = z.infer<typeof customerFormSchema>

export const createCustomerSchema = customerFormSchema
export const editCustomerSchema = customerFormSchema
