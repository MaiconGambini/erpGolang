import { z } from 'zod'

const saleItemSchema = z.object({
  productId: z.string().min(1, 'Produto é obrigatório'),
  quantity: z.coerce.number().int().min(1, 'Quantidade mínima é 1'),
})

export const createSaleSchema = z.object({
  customerId: z.string().min(1, 'Cliente é obrigatório'),
  notes: z.string().optional(),
  items: z.array(saleItemSchema).min(1, 'Adicione ao menos um item'),
})

export type CreateSaleFormValues = z.infer<typeof createSaleSchema>
