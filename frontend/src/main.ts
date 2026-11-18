import Aura from '@primeuix/themes/aura'
import PrimeVue from 'primevue/config'
import ToastService from 'primevue/toastservice'
import { createApp } from 'vue'
import './style.css'
import App from './App.vue'
import { initTheme } from './lib/theme'
import { router } from './router'

initTheme()

const app = createApp(App)

app.use(PrimeVue, {
  theme: {
    preset: Aura,
    options: {
      // Class-based dark mode so a manual toggle can drive it.
      darkModeSelector: '.p-dark',
    },
  },
})

app.use(router)
app.use(ToastService)

// Wait for the initial route to resolve before mounting so a deep link is not
// redirected to the default view on load.
router.isReady().then(() => app.mount('#app'))
