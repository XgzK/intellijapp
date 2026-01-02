<script setup lang="ts">
/**
 * 错误边界组件
 * 优化：捕获组件渲染错误，防止整个应用崩溃
 */
import { ref, onErrorCaptured } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const hasError = ref(false)
const errorMessage = ref('')
const errorStack = ref('')

// 捕获子组件错误
onErrorCaptured((err, instance, info) => {
  hasError.value = true
  errorMessage.value = err.message
  errorStack.value = err.stack || ''

  console.error('[ErrorBoundary] Caught error:', {
    error: err,
    component: instance,
    info,
  })

  // 返回 false 阻止错误继续向上传播
  return false
})

const resetError = () => {
  hasError.value = false
  errorMessage.value = ''
  errorStack.value = ''
}
</script>

<template>
  <div v-if="hasError" class="error-boundary">
    <div class="error-content">
      <div class="error-icon">⚠️</div>
      <h2 class="error-title">{{ t('errors.unknown') }}</h2>
      <p class="error-message">{{ errorMessage }}</p>
      <details v-if="errorStack" class="error-details">
        <summary>技术详情</summary>
        <pre>{{ errorStack }}</pre>
      </details>
      <button class="error-reset-button" @click="resetError">重新加载</button>
    </div>
  </div>
  <slot v-else />
</template>

<style scoped>
.error-boundary {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 400px;
  padding: var(--space-xl);
}

.error-content {
  max-width: 600px;
  text-align: center;
}

.error-icon {
  font-size: 4rem;
  margin-bottom: var(--space-md);
}

.error-title {
  font-size: 1.5rem;
  font-weight: 700;
  color: var(--color-danger);
  margin-bottom: var(--space-sm);
}

.error-message {
  font-size: 1rem;
  color: var(--color-muted);
  margin-bottom: var(--space-lg);
  line-height: 1.6;
}

.error-details {
  text-align: left;
  margin-bottom: var(--space-lg);
  padding: var(--space-md);
  background: rgba(0, 0, 0, 0.3);
  border-radius: 8px;
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.error-details summary {
  cursor: pointer;
  font-weight: 600;
  color: var(--color-accent);
  margin-bottom: var(--space-sm);
}

.error-details pre {
  font-size: 0.75rem;
  color: var(--color-muted);
  overflow-x: auto;
  margin: 0;
  padding: var(--space-sm);
  background: rgba(0, 0, 0, 0.2);
  border-radius: 4px;
}

.error-reset-button {
  padding: 0.75rem 1.5rem;
  border-radius: 999px;
  border: none;
  background: linear-gradient(135deg, var(--color-accent), var(--color-accent-strong));
  color: #fff;
  font-weight: 700;
  font-size: 0.95rem;
  cursor: pointer;
  transition:
    transform 0.2s ease,
    box-shadow 0.2s ease;
  box-shadow: 0 12px 24px rgba(45, 136, 255, 0.35);
}

.error-reset-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 16px 30px rgba(45, 136, 255, 0.45);
}
</style>
