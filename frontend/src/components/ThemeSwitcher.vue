<script setup lang="ts">
/**
 * 主题切换器组件
 * 单按钮切换亮色/暗色主题
 */
import { computed } from 'vue'
import { useTheme } from '@/composables/useTheme'

const { currentTheme, toggleTheme } = useTheme()

const themeIcon = computed(() => (currentTheme.value === 'dark' ? '🌙' : '☀️'))

// 导入 i18n 来访问翻译
import { useI18n } from 'vue-i18n'
const { t } = useI18n()
</script>

<template>
  <button
    class="theme-button"
    type="button"
    :title="currentTheme === 'dark' ? t('theme.switchToLight') : t('theme.switchToDark')"
    @click="toggleTheme"
  >
    <span class="theme-icon">{{ themeIcon }}</span>
    <span class="theme-label">{{ currentTheme === 'dark' ? t('theme.dark') : t('theme.light') }}</span>
  </button>
</template>

<style scoped>
.theme-button {
  display: flex;
  align-items: center;
  gap: 0.35rem;
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

.theme-button:hover {
  background: rgba(255, 255, 255, 0.12);
  color: var(--color-text);
  transform: translateY(-1px);
}

.theme-icon {
  font-size: 0.9rem;
  line-height: 1;
}

.theme-label {
  line-height: 1;
}

/* 响应式：小屏幕只显示图标 */
@media (max-width: 640px) {
  .theme-label {
    display: none;
  }

  .theme-button {
    padding: 0.5rem;
  }
}
</style>
