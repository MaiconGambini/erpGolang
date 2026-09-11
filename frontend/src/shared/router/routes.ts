import type { RouteRecordRaw } from 'vue-router'
import LoginPage from '@/pages/login/ui/LoginPage.vue'
import DashboardPage from '@/pages/dashboard/ui/DashboardPage.vue'
import CustomersPage from '@/pages/customers/ui/CustomersPage.vue'
import ProductsPage from '@/pages/products/ui/ProductsPage.vue'
import SuppliersPage from '@/pages/suppliers/ui/SuppliersPage.vue'
import SalesPage from '@/pages/sales/ui/SalesPage.vue'
import UsersPage from '@/pages/users/ui/UsersPage.vue'
import AuditPage from '@/pages/audit/ui/AuditPage.vue'
import NotFoundPage from '@/pages/not-found/ui/NotFoundPage.vue'
import { requireAuth, requireAdmin } from './guards'

export const routes: RouteRecordRaw[] = [
  { path: '/login', component: LoginPage, meta: { public: true } },
  { path: '/', component: DashboardPage, beforeEnter: requireAuth },
  { path: '/customers', component: CustomersPage, beforeEnter: requireAuth },
  { path: '/products', component: ProductsPage, beforeEnter: requireAuth },
  { path: '/suppliers', component: SuppliersPage, beforeEnter: requireAuth },
  { path: '/sales', component: SalesPage, beforeEnter: requireAuth },
  { path: '/users', component: UsersPage, beforeEnter: requireAdmin },
  { path: '/audit', component: AuditPage, beforeEnter: requireAdmin },
  { path: '/:pathMatch(.*)*', component: NotFoundPage },
]
