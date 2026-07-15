import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { clerkPlugin } from '@clerk/vue'
import App from './App.vue'
import router from './router'

const PUBLISHABLE_KEY = (import.meta as any).env.VITE_CLERK_PUBLISHABLE_KEY || 'pk_test_placeholder'

const app = createApp(App)
app.use(createPinia())
app.use(clerkPlugin, {
  publishableKey: PUBLISHABLE_KEY
})
app.use(router)

// Dummy Header component to satisfy local <Header> template tags in pages
app.component('Header', {
  render() {
    return null
  }
})

app.mount('#app')