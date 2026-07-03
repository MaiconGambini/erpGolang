import { z } from 'zod'

export const productFormSchema = z.object({
  name: z.string().min(1, 'Nome é obrigatório'),
  sku: z.string().min(1, 'SKU é obrigatório'),
  price: z
    .string()
    .min(1, 'Preço é obrigatório')
    .refine((v) => {
      const n = Number(v.replace(',', '.'))
      return !Number.isNaN(n) && n > 0
    }, 'Preço deve ser maior que zero'),
  stock: z.coerce.number().int().min(0, 'Estoque não pode ser negativo'),
  active: z.boolean().default(true),
})

export type ProductFormValues = z.infer<typeof productFormSchema>

export const createProductSchema = productFormSchema
export const editProductSchema = productFormSchema
