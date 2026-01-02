import { createApp } from 'vue'
import App from './App.vue'
import i18n from './i18n'
import { setupGlobalTheme } from './composables/useTheme'
import './styles/themes.css'
import './styles/global.css'

// 初始化主题系统
setupGlobalTheme()

createApp(App).use(i18n).mount('#app')
