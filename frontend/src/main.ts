import { createApp } from 'vue'
import App from './app/App.vue'
import { createAppPinia } from './app/providers/pinia'
import { createAppPrimeVue } from './app/providers/primevue'
import { createAppQuery } from './app/providers/query'
import { router } from './app/providers/router'
import { bootAuth } from '@/processes/auth/boot-auth'
import './app/styles/index.scss'

const app = createApp(App)
const pinia = createAppPinia()

app.use(pinia)
app.use(createAppQuery())
app.use(createAppPrimeVue())

bootAuth().finally(() => {
  app.use(router)
  app.mount('#app')
})
