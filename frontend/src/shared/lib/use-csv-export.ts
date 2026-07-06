import { ref } from 'vue'
import { exportList } from '@/shared/api/export'
import { downloadBlob } from '@/shared/lib/download'

export function useCsvExport(path: string, filename: string) {
  const exporting = ref(false)

  async function exportCsv(params: Record<string, string | number | boolean | undefined> = {}) {
    exporting.value = true
    try {
      const blob = await exportList(path, params)
      downloadBlob(filename, blob)
    } finally {
      exporting.value = false
    }
  }

  return { exporting, exportCsv }
}
