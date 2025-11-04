<script setup lang="ts">
import { ref } from 'vue'
import { SubmitPaths, ClearConfig, PathExists } from '@/services/configService'
import { Window, Browser } from '@wailsio/runtime'

const appWindow = Window

const currentView = ref<'main' | 'about'>('main')
const projectPath = ref('')
const configPath = ref('')
const statusMessage = ref('')
const submitting = ref(false)
const clearing = ref(false)
const isError = ref(false)
const CONFIG_PATH_PATTERN = /^[A-Za-z0-9:\\/]+$/

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
    statusMessage.value = '配置目录仅支持字母、数字及:/\\ 字符，其它内容不支持喵～'
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

const handleClose = async () => {
  try {
    await appWindow.Close()
  } catch (error) {
    console.error('关闭窗口失败', error)
  }
}

const handleMinimize = async () => {
  try {
    await appWindow.Minimise()
  } catch (error) {
    console.error('最小化窗口失败', error)
  }
}

const handleZoom = async () => {
  try {
    const maximised = await appWindow.IsMaximised()
    if (maximised) {
      await appWindow.Restore()
    } else {
      await appWindow.Maximise()
    }
  } catch (error) {
    console.error('切换窗口大小失败', error)
  }
}

const switchView = (view: 'main' | 'about') => {
  currentView.value = view
}

const openExternal = async (url: string) => {
  try {
    await Browser.OpenURL(url)
  } catch (error) {
    console.error('打开外部链接失败', url, error)
  }
}

</script>

<template>
  <div class="window">
    <header class="titlebar">
      <div class="traffic-lights">
        <button class="traffic-button close" type="button" @click.stop="handleClose" />
        <button class="traffic-button minimise" type="button" @click.stop="handleMinimize" />
        <button class="traffic-button zoom" type="button" @click.stop="handleZoom" />
      </div>
      <div class="title">IntelliJ 配置助手</div>
      <nav class="nav-buttons">
        <button
          class="nav-button"
          :class="{ 'nav-button--active': currentView === 'main' }"
          type="button"
          @click="switchView('main')"
        >
          主页
        </button>
        <button
          class="nav-button"
          :class="{ 'nav-button--active': currentView === 'about' }"
          type="button"
          @click="switchView('about')"
        >
          关于
        </button>
      </nav>
    </header>
    <main v-if="currentView === 'main'" class="workspace">
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
    <main v-else class="workspace">
      <section class="panel panel--path">
        <header class="panel__header">
          <h2 class="panel__title">关于 IntelliJ 配置助手</h2>
          <p class="panel__description">
            这里记录应用的版本、核心技术栈与维护者信息，布局保持与主页一致喵。
          </p>
        </header>
        <div class="panel__body">
          <div class="meta-list">
            <div class="meta-row">
              <span class="meta-label">应用名称</span>
              <span class="meta-value">IntelliJ 配置助手</span>
            </div>
            <div class="meta-row">
              <span class="meta-label">当前版本</span>
              <span class="meta-value">v1.0.2</span>
            </div>
            <div class="meta-row meta-row--link">
              <span class="meta-label">构建工具</span>
              <a
                href="https://github.com/wailsapp/wails"
                target="_blank"
                rel="noopener noreferrer"
                class="meta-link"
                @click.prevent="openExternal('https://github.com/wailsapp/wails')"
              >
                Wails 3
              </a>
            </div>
          </div>
        </div>
      </section>

      <section class="actions-grid">
        <article class="panel panel--action">
          <header class="panel__header">
            <h2 class="panel__title">技术栈</h2>
            <p class="panel__description">支撑 IntelliJ 配置助手的关键技术组件明细喵。</p>
          </header>
          <div class="panel__body">
            <div class="meta-list">
              <div class="meta-row">
                <span class="meta-label">前端框架</span>
                <span class="meta-value">Vue 3.5.22</span>
              </div>
              <div class="meta-row">
                <span class="meta-label">后端语言</span>
                <span class="meta-value">Go 1.25.3</span>
              </div>
              <div class="meta-row">
                <span class="meta-label">构建工具</span>
                <span class="meta-value">Wails 3</span>
              </div>
            </div>
          </div>
        </article>

        <article class="panel panel--action">
          <header class="panel__header">
            <h2 class="panel__title">项目信息</h2>
            <p class="panel__description">快速找到仓库地址与维护者，方便主人参考喵。</p>
          </header>
          <div class="panel__body">
            <div class="meta-list">
              <div class="meta-row meta-row--link">
                <span class="meta-label">仓库地址</span>
                <a
                  href="https://github.com/XgzK/intellijapp"
                  target="_blank"
                  rel="noopener noreferrer"
                  class="meta-link"
                  @click.prevent="openExternal('https://github.com/XgzK/intellijapp')"
                >
                  github.com/XgzK/intellijapp
                </a>
              </div>
              <div class="meta-row">
                <span class="meta-label">开发者</span>
                <span class="meta-value meta-value--developers">
                  <a
                    class="meta-link"
                    href="https://github.com/XgzK"
                    target="_blank"
                    rel="noopener noreferrer"
                    @click.prevent="openExternal('https://github.com/XgzK')"
                  >
                    XgzK
                  </a>
                  ·
                  <a
                    class="meta-link"
                    href="https://github.com/anthropics/claude-code"
                    target="_blank"
                    rel="noopener noreferrer"
                    @click.prevent="openExternal('https://github.com/anthropics/claude-code')"
                  >
                    Claude (AI)
                  </a>
                  ·
                  <a
                    class="meta-link"
                    href="https://github.com/openai/codex"
                    target="_blank"
                    rel="noopener noreferrer"
                    @click.prevent="openExternal('https://github.com/openai/codex')"
                  >
                    Codex (AI)
                  </a>
                </span>
              </div>
            </div>
          </div>
        </article>
      </section>

      <section class="panel panel--notice">
        <p class="about-notice">
          本项目运行所产生的一切问题需自行承担，项目由 Claude 4.5 和 GPT-5 配合开发，仅限学习使用喵。
        </p>
      </section>

      <p class="about-footer">© 2025 IntelliJ 配置助手 · 使用 Wails 和 Vue 构建</p>
    </main>
  </div>
</template>

<style scoped>
.window {
  --color-background: #0b1624;
  --color-surface: rgba(16, 28, 44, 0.78);
  --color-surface-strong: rgba(24, 40, 60, 0.92);
  --color-border: rgba(255, 255, 255, 0.12);
  --color-border-strong: rgba(66, 184, 131, 0.45);
  --color-text: #e8f1ff;
  --color-muted: rgba(232, 241, 255, 0.72);
  --color-accent: #42b883;
  --color-accent-strong: #2d88ff;
  --color-danger: #ee5a6f;
  --space-xs: 0.3rem;
  --space-sm: 0.5rem;
  --space-md: 0.75rem;
  --space-lg: 1.1rem;
  --space-xl: 1.6rem;
  height: 100vh;
  overflow: hidden;
  background: radial-gradient(circle at top left, #1b2636 0%, #0b1624 45%, #09111d 100%);
  color: var(--color-text);
  font-family: 'Segoe UI', sans-serif;
  display: flex;
  flex-direction: column;
}

:where(h1, h2, h3, p, span) {
  margin: 0;
}

:where(p) {
  line-height: 1.5;
}

.titlebar {
  height: 40px;
  padding: 0 calc(var(--space-lg) * 0.9);
  display: flex;
  align-items: center;
  gap: var(--space-md);
  background: rgba(6, 13, 22, 0.78);
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
  -webkit-app-region: drag;
}

.traffic-lights {
  display: flex;
  gap: 10px;
}

.traffic-button {
  width: 14px;
  height: 14px;
  border-radius: 50%;
  border: none;
  -webkit-app-region: no-drag;
  cursor: pointer;
  transition: transform 0.15s ease;
  position: relative;
  box-shadow: inset 0 1px 1px rgba(255, 255, 255, 0.6),
    inset 0 -1px 1px rgba(0, 0, 0, 0.4);
}

.traffic-button:hover {
  transform: scale(1.05);
}

.traffic-button.close {
  background: linear-gradient(145deg, #ff5f57, #e13c30);
}

.traffic-button.minimise {
  background: linear-gradient(145deg, #febc2e, #e3a400);
}

.traffic-button.zoom {
  background: linear-gradient(145deg, #28c840, #1ea333);
}

.title {
  font-size: 0.95rem;
  letter-spacing: 0.02em;
  opacity: 0.85;
  pointer-events: none;
  user-select: none;
  flex: 1;
}

.nav-buttons {
  display: flex;
  gap: 0.5rem;
  -webkit-app-region: no-drag;
}

.nav-button {
  padding: 0.35rem 1rem;
  border-radius: 999px;
  border: 1px solid transparent;
  background: rgba(255, 255, 255, 0.05);
  color: rgba(232, 241, 255, 0.75);
  font-size: 0.85rem;
  font-weight: 500;
  cursor: pointer;
  transition: background 0.2s ease, color 0.2s ease, border-color 0.2s ease, transform 0.15s ease;
}

.nav-button:hover {
  background: rgba(255, 255, 255, 0.12);
  color: var(--color-text);
  transform: translateY(-1px);
}

.nav-button--active {
  background: linear-gradient(135deg, var(--color-accent), var(--color-accent-strong));
  color: #fff;
  border-color: rgba(255, 255, 255, 0.1);
  box-shadow: 0 8px 18px rgba(45, 136, 255, 0.35);
}

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

.meta-list {
  display: flex;
  flex-direction: column;
  gap: var(--space-sm);
}

.meta-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: var(--space-sm);
  padding: 0.15rem 0;
}

.meta-row--column {
  flex-direction: column;
  align-items: flex-start;
}

.meta-row--link {
  gap: var(--space-sm);
}

.meta-label {
  font-size: 0.85rem;
  color: var(--color-muted);
  font-weight: 500;
}

.meta-value {
  font-size: 0.85rem;
  font-weight: 600;
  color: var(--color-text);
  font-family: 'Cascadia Code', 'Fira Code', 'Courier New', monospace;
  background: rgba(12, 24, 38, 0.6);
  padding: 0.2rem 0.45rem;
  border-radius: 6px;
}

.meta-value--developers {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 0.3rem;
  background: none;
  padding: 0;
}

.meta-link {
  font-size: 0.85rem;
  font-weight: 600;
  color: var(--color-accent-strong);
  text-decoration: none;
  transition: color 0.2s ease, background 0.2s ease;
  font-family: 'Cascadia Code', 'Fira Code', 'Courier New', monospace;
  word-break: break-all;
  background: rgba(12, 24, 38, 0.6);
  padding: 0.35rem 0.5rem;
  border-radius: 6px;
  display: inline-flex;
}

.meta-link:hover {
  color: var(--color-accent);
  background: rgba(12, 24, 38, 0.75);
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

.about-notice {
  font-size: 0.85rem;
  color: rgba(255, 230, 180, 0.92);
  line-height: 1.6;
  text-align: center;
}

.about-footer {
  text-align: center;
  margin-top: var(--space-lg);
  padding-top: var(--space-sm);
  border-top: 1px solid rgba(255, 255, 255, 0.08);
  color: rgba(232, 241, 255, 0.5);
  font-size: 0.75rem;
}

@media (max-width: 640px) {
  .titlebar {
    padding: 0 var(--space-sm);
    gap: var(--space-sm);
  }

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
