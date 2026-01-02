<script setup lang="ts">
import { ref } from 'vue'
import { SubmitPaths, ClearConfig, PathExists } from '@/services/configService'
import { CONFIG_PATH_PATTERN, VALIDATION_MESSAGES, PATH_LABELS } from '@/constants/validation'
import { useErrorHandler } from '@/composables/useErrorHandler'

const projectPath = ref('')
const configPath = ref('')
const statusMessage = ref('')
const submitting = ref(false)
const clearing = ref(false)
const isError = ref(false)

const { handleError: handleGlobalError } = useErrorHandler()

// 优化：使用统一的错误处理器
const extractErrorMessage = (error: unknown, context?: string): string => {
  return handleGlobalError(error, context)
}

const ensurePathExists = async (path: string, label: string) => {
  try {
    const exists = await PathExists(path)
    if (!exists) {
      statusMessage.value = `${label}${VALIDATION_MESSAGES.PATH_NOT_EXIST_SUFFIX}`
      isError.value = true
      return false
    }
    return true
  } catch (error) {
    statusMessage.value = extractErrorMessage(error)
    isError.value = true
    return false
  }
}

const handleSubmit = async () => {
  if (!projectPath.value || !configPath.value) {
    statusMessage.value = VALIDATION_MESSAGES.EMPTY_PATHS
    isError.value = true
    return
  }

  if (!CONFIG_PATH_PATTERN.test(configPath.value)) {
    statusMessage.value = VALIDATION_MESSAGES.INVALID_CONFIG_PATH
    isError.value = true
    return
  }

  statusMessage.value = ''
  isError.value = false

  if (!(await ensurePathExists(projectPath.value, PATH_LABELS.INTELLIJ_PATH))) {
    return
  }

  if (!(await ensurePathExists(configPath.value, PATH_LABELS.CONFIG_PATH))) {
    return
  }

  submitting.value = true
  statusMessage.value = ''
  isError.value = false
  try {
    const response = await SubmitPaths(projectPath.value, configPath.value)
    statusMessage.value = response
    isError.value = false
  } catch (error) {
    statusMessage.value = extractErrorMessage(error)
    isError.value = true
  } finally {
    submitting.value = false
  }
}

const handleClear = async () => {
  if (!projectPath.value) {
    statusMessage.value = VALIDATION_MESSAGES.EMPTY_INTELLIJ_PATH
    isError.value = true
    return
  }

  statusMessage.value = ''
  isError.value = false

  if (!(await ensurePathExists(projectPath.value, PATH_LABELS.INTELLIJ_PATH))) {
    return
  }

  clearing.value = true
  statusMessage.value = ''
  isError.value = false
  try {
    const response = await ClearConfig(projectPath.value)
    statusMessage.value = response
    isError.value = false
  } catch (error) {
    statusMessage.value = extractErrorMessage(error)
    isError.value = true
  } finally {
    clearing.value = false
  }
}
</script>

<template>
  <main class="workspace">
    <section class="panel panel--path">
      <header class="panel__header">
        <h2 class="panel__title">{{ $t('mainView.installPath.title') }}</h2>
      </header>
      <div class="panel__body">
        <label class="field field--inline">
          <span class="field__label">{{ $t('mainView.installPath.label') }}</span>
          <input
            v-model.trim="projectPath"
            class="input"
            type="text"
            autocomplete="off"
            :placeholder="$t('mainView.installPath.placeholder')"
          />
        </label>
        <label class="field field--inline">
          <span class="field__label">{{ $t('mainView.applyConfig.label') }}</span>
          <input
            v-model.trim="configPath"
            class="input"
            type="text"
            autocomplete="off"
            :placeholder="$t('mainView.applyConfig.placeholder')"
          />
        </label>
      </div>
    </section>

    <section class="actions-grid">
      <form class="panel panel--action" @submit.prevent="handleSubmit">
        <footer class="panel__footer panel__footer--actions">
          <button class="button button--primary" type="submit" :disabled="submitting || clearing">
            {{ submitting ? $t('mainView.applyConfig.submitting') : $t('mainView.applyConfig.submitButton') }}
          </button>
          <button
            class="button button--danger"
            type="button"
            :disabled="submitting || clearing"
            @click="handleClear"
          >
            {{ clearing ? $t('mainView.clearConfig.clearing') : $t('mainView.clearConfig.clearButton') }}
          </button>
        </footer>
      </form>
    </section>

    <p
      v-if="statusMessage"
      :class="['feedback', isError ? 'feedback--error' : 'feedback--success']"
    >
      {{ statusMessage }}
    </p>

    <section class="panel panel--notice panel--warning">
      <p class="notice-text" v-html="$t('mainView.notice.text')"></p>
      <ul class="notice-list">
        <li v-html="$t('mainView.notice.windows')"></li>
        <li v-html="$t('mainView.notice.linuxMac')"></li>
      </ul>
    </section>
  </main>
</template>

<style scoped>
/* 组件特定样式 - 公共样式已移至 global.css */

.field {
  display: flex;
  flex-direction: column;
  gap: var(--space-xs);
  font-size: 0.875rem;
}

.field--stack {
  width: 100%;
}

.field--inline {
  flex-direction: row;
  align-items: center;
  gap: var(--space-sm);
}

.field--inline .field__label {
  width: 140px;
  text-align: right;
  margin-bottom: 0;
}

.field--inline .input {
  flex: 1;
}

.field__label {
  font-weight: 600;
  color: var(--color-text);
}

.input {
  padding: 0.65rem 0.75rem;
  border-radius: 10px;
  border: 1px solid var(--color-border);
  background: var(--input-background);
  color: inherit;
  font-size: 0.875rem;
  transition:
    border-color 0.2s ease,
    box-shadow 0.2s ease,
    background 0.2s ease,
    transform 0.2s ease;
}

.input:hover {
  border-color: var(--color-border-strong);
  background: var(--input-background-hover);
  transform: translateY(-1px);
}

.input:focus {
  outline: none;
  border-color: var(--color-accent);
  box-shadow: 0 0 0 3px rgba(66, 184, 131, 0.18);
  background: var(--input-background-focus);
  transform: translateY(-1px);
}

.button {
  min-width: 150px;
  padding: 0.65rem 1.1rem;
  border-radius: 999px;
  border: 1px solid transparent;
  font-weight: 700;
  font-size: 0.9rem;
  cursor: pointer;
  transition:
    transform 0.15s ease,
    box-shadow 0.2s ease,
    opacity 0.2s ease,
    border-color 0.2s ease;
  color: #fff;
}

.button--primary {
  background: linear-gradient(135deg, var(--color-accent), var(--color-accent-strong));
  box-shadow: var(--shadow-primary);
}

.button--primary:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-primary-hover);
}

.button--primary:active {
  transform: translateY(0);
}

.button--danger {
  background: linear-gradient(135deg, var(--color-danger), #ff7b84);
  box-shadow: var(--shadow-danger);
}

.button--danger:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-danger-hover);
}

.button--danger:active {
  transform: translateY(0);
}

.button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
  box-shadow: none;
}

.panel__footer--actions {
  gap: var(--space-sm);
  flex-wrap: wrap;
  justify-content: center;
}

.feedback {
  text-align: center;
  font-size: 0.9rem;
  padding: 0.8rem 1rem;
  border-radius: 10px;
  background: rgba(255, 255, 255, 0.08);
  box-shadow: var(--shadow-feedback);
  border: 1px solid rgba(255, 255, 255, 0.1);
  font-weight: 500;
}

.feedback--success {
  color: var(--feedback-success-color);
  border-color: var(--feedback-success-border);
  background: var(--feedback-success-bg);
}

.feedback--error {
  color: var(--feedback-error-color);
  border-color: var(--feedback-error-border);
  background: var(--feedback-error-bg);
}

.notice-text {
  font-size: 0.85rem;
  color: var(--warning-text);
}

.notice-list {
  margin: 0;
  padding-left: 1.1rem;
  display: flex;
  flex-direction: column;
  gap: 0.35rem;
  color: var(--warning-text-alt);
  font-size: 0.85rem;
}

/* 组件特定的响应式样式 */
@media (max-width: 640px) {
  .button {
    width: 100%;
    justify-content: center;
  }

  .panel__footer {
    justify-content: stretch;
  }

  .panel__footer .button {
    width: 100%;
  }
}
</style>
