/**
 * i18n 国际化配置
 * 优化：提供完整的多语言支持基础设施
 */
import { createI18n } from 'vue-i18n'
import zhCN from './locales/zh-CN'
import enUS from './locales/en-US'

// 定义支持的语言类型
export type SupportedLocale = 'zh-CN' | 'en-US'

// 语言包映射
const messages = {
  'zh-CN': zhCN,
  'en-US': enUS,
}

// 获取默认语言（优先使用浏览器语言，回退到中文）
function getDefaultLocale(): SupportedLocale {
  const browserLang = navigator.language
  if (browserLang.startsWith('en')) {
    return 'en-US'
  }
  return 'zh-CN'
}

// 创建 i18n 实例
export const i18n = createI18n({
  legacy: false, // 使用 Composition API 模式
  locale: getDefaultLocale(),
  fallbackLocale: 'zh-CN',
  messages,
  globalInjection: true, // 全局注入 $t 函数
})

export default i18n
