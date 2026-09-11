export function debounce<T extends (...args: unknown[]) => void>(fn: T, waitMs: number) {
  let timeout: ReturnType<typeof setTimeout> | undefined

  return (...args: Parameters<T>) => {
    if (timeout) clearTimeout(timeout)
    timeout = setTimeout(() => fn(...args), waitMs)
  }
}
