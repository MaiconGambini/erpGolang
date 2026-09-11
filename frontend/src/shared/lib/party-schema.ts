import { z } from 'zod'
import { stripDocument } from '@/shared/lib/document'

const optionalEmail = z.union([z.string().email(), z.literal('')]).optional()
const optionalText = z.string().optional()

export const partyFormSchema = z
  .object({
    name: z.string().min(1, 'Nome é obrigatório'),
    documentType: z.union([z.enum(['cpf', 'cnpj']), z.literal('')]).optional(),
    document: optionalText,
    email: optionalEmail,
    phone: optionalText,
    postalCode: optionalText,
    street: optionalText,
    streetNumber: optionalText,
    city: optionalText,
    state: z
      .union([z.literal(''), z.string().length(2, 'UF deve ter 2 letras')])
      .optional(),
    active: z.boolean().default(true),
  })
  .superRefine((data, ctx) => {
    if (data.documentType && data.document) {
      const digits = stripDocument(data.document)
      const expected = data.documentType === 'cpf' ? 11 : 14
      if (digits.length !== expected) {
        ctx.addIssue({
          code: 'custom',
          message: data.documentType === 'cpf' ? 'CPF inválido' : 'CNPJ inválido',
          path: ['document'],
        })
      }
    }
  })

export type PartyFormValues = z.infer<typeof partyFormSchema>
