import type { RouteRecordRaw } from 'vue-router'

import LoginPage from '@/pages/login/ui/LoginPage.vue'

import DashboardPage from '@/pages/dashboard/ui/DashboardPage.vue'

import CustomersPage from '@/pages/customers/ui/CustomersPage.vue'
import ProductsPage from '@/pages/products/ui/ProductsPage.vue'

import NotFoundPage from '@/pages/not-found/ui/NotFoundPage.vue'

import { requireAuth } from './guards'



export const routes: RouteRecordRaw[] = [

  { path: '/login', component: LoginPage, meta: { public: true } },

  { path: '/', component: DashboardPage, beforeEnter: requireAuth },

  { path: '/customers', component: CustomersPage, beforeEnter: requireAuth },
  { path: '/products', component: ProductsPage, beforeEnter: requireAuth },

  { path: '/:pathMatch(.*)*', component: NotFoundPage },

]

