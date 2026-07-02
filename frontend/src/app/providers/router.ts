import { createRouter, createWebHistory } from 'vue-router'
import { routes } from '@/shared/router/routes'

export const router = createRouter({
  history: createWebHistory(),
  routes,
})
