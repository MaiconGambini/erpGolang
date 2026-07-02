import { createApp } from 'vue'
import App from './app/App.vue'
import { createAppPinia } from './app/providers/pinia'
import { createAppPrimeVue } from './app/providers/primevue'
import { createAppQuery } from './app/providers/query'
import { router } from './app/providers/router'
import './app/styles/index.scss'

const app = createApp(App)

app.use(createAppPinia())
app.use(createAppQuery())
app.use(createAppPrimeVue())
app.use(router)

app.mount('#app')
