/**
 * 主题管理 Composable
 * 默认暗色，支持手动切换亮色/暗色
 */
import { ref } from 'vue'

export type Theme = 'light' | 'dark'

const STORAGE_KEY = 'app-theme'
// 默认强制暗色
const currentTheme = ref<Theme>('dark')
const isDark = ref(true)

/**
 * 应用主题到 DOM
 */
function applyTheme(theme: 'light' | 'dark') {
  const root = document.documentElement
  root.setAttribute('data-theme', theme)
  isDark.value = theme === 'dark'
}

/**
 * 主题管理 Hook
 */
export function useTheme() {
  /**
   * 设置主题
   */
  const setTheme = (theme: Theme) => {
    currentTheme.value = theme
    localStorage.setItem(STORAGE_KEY, theme)
    applyTheme(theme)
  }

  /**
   * 切换主题 (亮色 ↔ 暗色)
   */
  const toggleTheme = () => {
    const newTheme: Theme = currentTheme.value === 'dark' ? 'light' : 'dark'
    setTheme(newTheme)
  }

  /**
   * 初始化主题
   */
  const initTheme = () => {
    // 从 localStorage 读取保存的主题
    const savedTheme = localStorage.getItem(STORAGE_KEY) as Theme | null
    if (savedTheme && ['light', 'dark'].includes(savedTheme)) {
      currentTheme.value = savedTheme
    } else {
      // 如果没有保存的配置，强制默认暗色
      currentTheme.value = 'dark'
    }

    // 应用主题
    applyTheme(currentTheme.value)
  }

  return {
    currentTheme,
    isDark,
    setTheme,
    toggleTheme,
    initTheme,
  }
}

// 创建全局主题管理器实例
let globalThemeManager: ReturnType<typeof useTheme> | null = null

export function setupGlobalTheme() {
  if (!globalThemeManager) {
    globalThemeManager = useTheme()
    globalThemeManager.initTheme()
  }
  return globalThemeManager
}

export function getGlobalTheme() {
  if (!globalThemeManager) {
    return setupGlobalTheme()
  }
  return globalThemeManager
}
