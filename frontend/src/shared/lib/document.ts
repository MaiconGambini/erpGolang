export type DocumentType = 'cpf' | 'cnpj'

export function stripDocument(value: string): string {
  return value.replace(/\D/g, '')
}

export function formatDocument(value: string, type?: DocumentType | ''): string {
  const digits = stripDocument(value)
  if (type === 'cpf') {
    return digits
      .slice(0, 11)
      .replace(/(\d{3})(\d)/, '$1.$2')
      .replace(/(\d{3})(\d)/, '$1.$2')
      .replace(/(\d{3})(\d{1,2})$/, '$1-$2')
  }
  if (type === 'cnpj') {
    return digits
      .slice(0, 14)
      .replace(/^(\d{2})(\d)/, '$1.$2')
      .replace(/^(\d{2})\.(\d{3})(\d)/, '$1.$2.$3')
      .replace(/\.(\d{3})(\d)/, '.$1/$2')
      .replace(/(\d{4})(\d{1,2})$/, '$1-$2')
  }
  return digits
}

export function formatPostalCode(value: string): string {
  const digits = stripDocument(value).slice(0, 8)
  if (digits.length <= 5) return digits
  return `${digits.slice(0, 5)}-${digits.slice(5)}`
}
