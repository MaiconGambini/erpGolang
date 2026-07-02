import { ref } from 'vue'

export function useCustomerFilters() {
  const search = ref('')
  const status = ref<'all' | 'active' | 'inactive'>('all')

  return { search, status }
}
