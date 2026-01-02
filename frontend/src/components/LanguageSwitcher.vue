<script setup lang="ts">
/**
 * 语言切换器组件
 * 单按钮切换中文/英文
 */
import { computed } from 'vue'
import { i18n } from '@/i18n'
import type { SupportedLocale } from '@/i18n'

// 使用全局 i18n 实例来修改语言
const locale = computed({
  get: () => i18n.global.locale.value as SupportedLocale,
  set: (val: SupportedLocale) => {
    i18n.global.locale.value = val
  },
})

const toggleLanguage = () => {
  locale.value = (locale.value === 'zh-CN' ? 'en-US' : 'zh-CN') as SupportedLocale
}

const currentLangName = computed(() => (locale.value === 'zh-CN' ? '中文' : 'EN'))

// 导入 useI18n 来访问翻译
import { useI18n } from 'vue-i18n'
const { t } = useI18n()
</script>

<template>
  <button
    class="lang-button"
    type="button"
    :title="t('language.switchTo')"
    @click="toggleLanguage"
  >
    🌐 {{ currentLangName }}
  </button>
</template>

<style scoped>
.lang-button {
  padding: 0.35rem 0.8rem;
  border-radius: 6px;
  border: 1px solid rgba(255, 255, 255, 0.2);
  background: rgba(255, 255, 255, 0.05);
  color: rgba(232, 241, 255, 0.75);
  font-size: 0.8rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}

.lang-button:hover {
  background: rgba(255, 255, 255, 0.12);
  color: var(--color-text);
  transform: translateY(-1px);
}

/* 响应式：小屏幕只显示语言代码 */
@media (max-width: 640px) {
  .lang-button {
    padding: 0.5rem;
  }
}
</style>
