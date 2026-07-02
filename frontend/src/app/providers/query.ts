import { VueQueryPlugin, QueryClient } from '@tanstack/vue-query'

export function createAppQuery() {
  return {
    install(app: import('vue').App) {
      app.use(VueQueryPlugin, { queryClient: new QueryClient() })
    },
  }
}
