import { z } from 'zod'

const optionalEmail = z.union([z.string().email(), z.literal('')]).optional()
const optionalText = z.string().optional()

export const supplierFormSchema = z.object({
  name: z.string().min(1, 'Nome é obrigatório'),
  document: optionalText,
  email: optionalEmail,
  phone: optionalText,
  active: z.boolean().default(true),
})

export type SupplierFormValues = z.infer<typeof supplierFormSchema>

export const createSupplierSchema = supplierFormSchema
export const editSupplierSchema = supplierFormSchema
