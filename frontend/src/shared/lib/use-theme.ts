import { computed, ref, watch } from 'vue'

export type ThemePreference = 'light' | 'dark' | 'system'

const STORAGE_KEY = 'goerp-theme'

function readStored(): ThemePreference {
  try {
    const value = localStorage.getItem(STORAGE_KEY)
    return value === 'light' || value === 'dark' ? value : 'system'
  } catch {
    return 'system'
  }
}

function readSystemDark(): boolean {
  return typeof matchMedia !== 'undefined' && matchMedia('(prefers-color-scheme: dark)').matches
}

// Module-level singleton: every consumer shares one source of truth.
const preference = ref<ThemePreference>(readStored())
const systemDark = ref(readSystemDark())

if (typeof matchMedia !== 'undefined') {
  matchMedia('(prefers-color-scheme: dark)').addEventListener('change', (event) => {
    systemDark.value = event.matches
  })
}

const isDark = computed(
  () => preference.value === 'dark' || (preference.value === 'system' && systemDark.value),
)

watch(
  isDark,
  (dark) => {
    document.documentElement.classList.toggle('dark', dark)
  },
  { immediate: true },
)

export function useTheme() {
  function toggle() {
    preference.value = isDark.value ? 'light' : 'dark'
    try {
      localStorage.setItem(STORAGE_KEY, preference.value)
    } catch {
      // private mode: preference lives for the session only
    }
  }

  return { preference, isDark, toggle }
}
