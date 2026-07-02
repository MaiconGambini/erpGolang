import type { RouteRecordRaw } from 'vue-router'
import LoginPage from '@/pages/login/ui/LoginPage.vue'
import DashboardPage from '@/pages/dashboard/ui/DashboardPage.vue'
import CustomersPage from '@/pages/customers/ui/CustomersPage.vue'
import NotFoundPage from '@/pages/not-found/ui/NotFoundPage.vue'

export const routes: RouteRecordRaw[] = [
  { path: '/login', component: LoginPage },
  { path: '/', component: DashboardPage },
  { path: '/customers', component: CustomersPage },
  { path: '/:pathMatch(.*)*', component: NotFoundPage },
]
