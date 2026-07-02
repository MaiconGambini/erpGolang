import { defineStore } from 'pinia'
import type { SessionUser } from './types'

export const useSessionStore = defineStore('session', {
  state: () => ({
    accessToken: null as string | null,
    user: null as SessionUser | null,
  }),
  getters: {
    isAuthenticated: (state) => Boolean(state.accessToken),
  },
})
