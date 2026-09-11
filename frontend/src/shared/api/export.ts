import { apiClient } from '@/shared/api/client'

export async function exportList(
  path: string,
  params: Record<string, string | number | boolean | undefined> = {},
): Promise<Blob> {
  const response = await apiClient.get(path, {
    params: { ...params, format: 'csv' },
    responseType: 'blob',
  })
  return response.data as Blob
}
