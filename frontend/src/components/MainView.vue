<script setup lang="ts">
import { ref } from 'vue'
import { SubmitPaths, ClearConfig, PathExists } from '@/services/configService'

const projectPath = ref('')
const configPath = ref('')
const statusMessage = ref('')
const submitting = ref(false)
const clearing = ref(false)
const isError = ref(false)
const CONFIG_PATH_PATTERN = /^[A-Za-z0-9:\\/\s_\-.]+$/

const extractErrorMessage = (error: unknown): string => {
  const parse = (value: string) => {
    const trimmed = value.trim()
    if (!trimmed) {
      return '未知错误'
    }

    const descIndex = trimmed.lastIndexOf('desc=')
    if (descIndex !== -1) {
      return trimmed.substring(descIndex + 5).trim()
    }

    const colonIndex = trimmed.lastIndexOf(':')
    if (colonIndex !== -1) {
      return trimmed.substring(colonIndex + 1).trim()
    }

    return trimmed
  }

  if (error instanceof Error) {
    return parse(error.message)
  }

  if (typeof error === 'string') {
    return parse(error)
  }

  return parse(String(error))
}

const ensurePathExists = async (path: string, label: string) => {
  try {
    const exists = await PathExists(path)
    if (!exists) {
      statusMessage.value = `${label} 不存在，请检查喵～`
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
    statusMessage.value = '请输入完整的两个路径喵～'
    isError.value = true
    return
  }

  if (!CONFIG_PATH_PATTERN.test(configPath.value)) {
    statusMessage.value = '配置目录路径包含不支持的特殊字符喵～'
    isError.value = true
    return
  }

  statusMessage.value = ''
  isError.value = false

  if (!(await ensurePathExists(projectPath.value, 'IntelliJ 安装路径'))) {
    return
  }

  if (!(await ensurePathExists(configPath.value, '配置目录'))) {
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
    statusMessage.value = '请输入 IntelliJ 路径'
    isError.value = true
    return
  }

  statusMessage.value = ''
  isError.value = false

  if (!(await ensurePathExists(projectPath.value, 'IntelliJ 安装路径'))) {
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
        <h2 class="panel__title">安装路径</h2>
        <p class="panel__description">
          指向 IntelliJ 的安装目录，后续所有动作都会以这里为基准喵。
        </p>
      </header>
      <label class="field field--stack">
        <span class="field__label">IntelliJ 安装目录</span>
        <input
          v-model.trim="projectPath"
          class="input"
          type="text"
          autocomplete="off"
          placeholder="例如：D:/Program Files/JetBrains/IntelliJ IDEA 2024.2"
        />
      </label>
    </section>

    <section class="actions-grid">
      <form class="panel panel--action" @submit.prevent="handleSubmit">
        <header class="panel__header">
          <h2 class="panel__title">应用配置</h2>
          <p class="panel__description">
            选择需要导入的配置目录，浮浮酱会把它融合进上方的 IntelliJ 路径喵。
          </p>
        </header>
        <div class="panel__body">
          <label class="field field--stack">
            <span class="field__label">配置目录</span>
            <input
              v-model.trim="configPath"
              class="input"
              type="text"
              autocomplete="off"
              placeholder="例如：D:/jetbra"
            />
          </label>
        </div>
        <footer class="panel__footer">
          <button class="button button--primary" type="submit" :disabled="submitting || clearing">
            {{ submitting ? '提交中...' : '应用配置' }}
          </button>
        </footer>
      </form>

      <article class="panel panel--action">
        <header class="panel__header">
          <h2 class="panel__title">清除配置</h2>
          <p class="panel__description">
            恢复 IntelliJ 为初始状态，只会动到上方路径对应的配置喵。
          </p>
        </header>
        <div class="panel__body panel__body--compact">
          <p class="note">
            当前路径：
            <span class="note__value">{{ projectPath || '尚未填写' }}</span>
          </p>
        </div>
        <footer class="panel__footer">
          <button
            class="button button--danger"
            type="button"
            :disabled="submitting || clearing"
            @click="handleClear"
          >
            {{ clearing ? '清除中...' : '清除配置' }}
          </button>
        </footer>
      </article>
    </section>

    <p v-if="statusMessage" :class="['feedback', isError ? 'feedback--error' : 'feedback--success']">
      {{ statusMessage }}
    </p>

    <section class="panel panel--notice panel--warning">
      <p class="notice-text">
        若上述操作未成功，请下载压缩包后进入 <code>scripts</code> 文件夹按操作系统执行下列脚本喵：
      </p>
      <ul class="notice-list">
        <li>Windows：运行 <code>uninstall-all-users.vbs</code> 与 <code>uninstall-current-user.vbs</code></li>
        <li>Linux / macOS：运行 <code>uninstall.sh</code></li>
      </ul>
    </section>
  </main>
</template>

<style scoped>
.workspace {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: var(--space-md);
  padding: var(--space-xl) var(--space-lg);
  -webkit-app-region: no-drag;
  overflow-y: auto;
  scrollbar-width: none;
}

.workspace::-webkit-scrollbar {
  display: none;
}

.panel {
  display: flex;
  flex-direction: column;
  gap: var(--space-md);
  padding: calc(var(--space-md) + 0.2rem) var(--space-lg);
  border-radius: 14px;
  background: var(--color-surface);
  box-shadow: 0 18px 36px rgba(0, 0, 0, 0.32);
  border: 1px solid var(--color-border);
  backdrop-filter: blur(12px);
}

.panel--path {
  border-color: var(--color-border-strong);
}

.panel--action {
  justify-content: space-between;
}

.panel--notice {
  background: rgba(255, 255, 255, 0.06);
  border-color: rgba(255, 255, 255, 0.08);
}

.panel--warning {
  background: rgba(255, 161, 76, 0.15);
  border-color: rgba(255, 191, 120, 0.35);
}

.panel__header {
  display: flex;
  flex-direction: column;
  gap: var(--space-xs);
}

.panel__title {
  font-size: 1.05rem;
  font-weight: 700;
  color: var(--color-text);
}

.panel__description {
  font-size: 0.85rem;
  color: var(--color-muted);
}

.panel__body {
  display: flex;
  flex-direction: column;
  gap: var(--space-md);
}

.panel__body--compact {
  gap: var(--space-sm);
}

.panel__footer {
  display: flex;
  justify-content: flex-end;
}

.actions-grid {
  display: grid;
  gap: var(--space-lg);
  grid-template-columns: repeat(auto-fit, minmax(260px, 1fr));
}

.field {
  display: flex;
  flex-direction: column;
  gap: var(--space-xs);
  font-size: 0.875rem;
}

.field--stack {
  width: 100%;
}

.field__label {
  font-weight: 600;
  color: var(--color-text);
}

.input {
  padding: 0.65rem 0.75rem;
  border-radius: 10px;
  border: 1px solid rgba(255, 255, 255, 0.25);
  background: rgba(12, 24, 38, 0.88);
  color: inherit;
  font-size: 0.875rem;
  transition: border-color 0.2s ease, box-shadow 0.2s ease, background 0.2s ease;
}

.input:hover {
  border-color: rgba(255, 255, 255, 0.35);
  background: rgba(12, 24, 38, 0.95);
}

.input:focus {
  outline: none;
  border-color: var(--color-accent);
  box-shadow: 0 0 0 3px rgba(66, 184, 131, 0.18);
  background: rgba(12, 24, 38, 1);
}

.button {
  min-width: 150px;
  padding: 0.65rem 1.1rem;
  border-radius: 999px;
  border: 1px solid transparent;
  font-weight: 700;
  font-size: 0.9rem;
  cursor: pointer;
  transition: transform 0.15s ease, box-shadow 0.2s ease, opacity 0.2s ease, border-color 0.2s ease;
  color: #fff;
}

.button--primary {
  background: linear-gradient(135deg, var(--color-accent), var(--color-accent-strong));
  box-shadow: 0 12px 24px rgba(45, 136, 255, 0.35);
}

.button--primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 16px 30px rgba(45, 136, 255, 0.45);
}

.button--danger {
  background: linear-gradient(135deg, var(--color-danger), #ff7b84);
  box-shadow: 0 12px 24px rgba(238, 90, 111, 0.35);
}

.button--danger:hover {
  transform: translateY(-2px);
  box-shadow: 0 16px 30px rgba(238, 90, 111, 0.45);
}

.button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
  box-shadow: none;
}

.note {
  display: flex;
  align-items: center;
  gap: var(--space-xs);
  color: var(--color-muted);
  font-size: 0.85rem;
}

.note__value {
  font-family: 'Cascadia Code', 'Fira Code', 'Courier New', monospace;
  background: rgba(12, 24, 38, 0.75);
  padding: 0.25rem 0.45rem;
  border-radius: 8px;
  color: #9de3ff;
}

.feedback {
  text-align: center;
  font-size: 0.9rem;
  padding: 0.8rem 1rem;
  border-radius: 10px;
  background: rgba(255, 255, 255, 0.08);
  box-shadow: 0 12px 24px rgba(0, 0, 0, 0.28);
  border: 1px solid rgba(255, 255, 255, 0.1);
  font-weight: 500;
}

.feedback--success {
  color: #8ef5c3;
  border-color: rgba(142, 245, 195, 0.25);
  background: rgba(66, 184, 131, 0.12);
}

.feedback--error {
  color: #ffb8b0;
  border-color: rgba(255, 184, 176, 0.25);
  background: rgba(255, 107, 107, 0.12);
}

.notice-text {
  font-size: 0.85rem;
  color: rgba(255, 230, 200, 0.9);
}

.notice-list {
  margin: 0;
  padding-left: 1.1rem;
  display: flex;
  flex-direction: column;
  gap: 0.35rem;
  color: rgba(255, 235, 210, 0.9);
  font-size: 0.85rem;
}

@media (max-width: 640px) {
  .workspace {
    padding: var(--space-lg);
  }

  .panel {
    padding: var(--space-md);
  }

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
