<script setup lang="ts">
import { ref } from 'vue'
import { SubmitPaths, ClearConfig } from '@/services/configService'
import { Window } from '@wailsio/runtime'

const appWindow = Window

const projectPath = ref('')
const configPath = ref('')
const statusMessage = ref('')
const submitting = ref(false)
const clearing = ref(false)
const isError = ref(false)

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

const handleSubmit = async () => {
  if (!projectPath.value || !configPath.value) {
    statusMessage.value = '请输入完整的两个路径喵～'
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
</script>

<template>
  <div class="window">
    <header class="titlebar">
      <div class="traffic-lights">
        <button class="traffic-button close" type="button" @click.stop="handleClose" />
        <button class="traffic-button minimise" type="button" @click.stop="handleMinimize" />
        <button class="traffic-button zoom" type="button" @click.stop="handleZoom" />
      </div>
      <div class="title">IntelliJ 软件路径提交助手</div>
    </header>
    <main class="container">
      <h1>提交 IntelliJ 软件路径喵</h1>
      <form class="form" @submit.prevent="handleSubmit">
        <label class="field">
          <span>IntelliJ 软件路径地址</span>
          <input
            v-model.trim="projectPath"
            class="input"
            type="text"
            placeholder="例如：D:/Program Files/JetBrains/IntelliJ IDEA 2024.2"
          />
        </label>

        <label class="field">
          <span>配置文件目录</span>
          <input
            v-model.trim="configPath"
            class="input"
            type="text"
            placeholder="例如：D:/jetbra"
          />
        </label>

        <div class="button-group">
          <button class="submit" type="submit" :disabled="submitting || clearing">
            {{ submitting ? '提交中...' : '应用配置' }}
          </button>
          <button
            class="clear-button"
            type="button"
            :disabled="submitting || clearing"
            @click="handleClear"
          >
            {{ clearing ? '清除中...' : '清除配置' }}
          </button>
        </div>
      </form>

      <p
        v-if="statusMessage"
        :class="['status', isError ? 'status--error' : 'status--success']"
      >
        {{ statusMessage }}
      </p>
    </main>
  </div>
</template>

<style scoped>
.window {
  height: 100vh;
  overflow: hidden;
  background: radial-gradient(circle at top left, #1b2636 0%, #0b1624 45%, #09111d 100%);
  color: #e8f1ff;
  font-family: 'Segoe UI', sans-serif;
  display: flex;
  flex-direction: column;
}

.titlebar {
  height: 38px;
  padding: 0 16px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: rgba(6, 13, 22, 0.75);
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
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
  text-align: center;
  flex: 1;
}

.container {
  max-width: 520px;
  margin: 0 auto;
  padding: 2.5rem 1.5rem;
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
  -webkit-app-region: no-drag;
  flex: 1;
  overflow-y: auto;
  scrollbar-width: none;
}

.container::-webkit-scrollbar {
  display: none;
}

h1 {
  font-size: 1.6rem;
  text-align: center;
}

.form {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  padding: 1.75rem;
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.08);
  box-shadow: 0 18px 36px rgba(0, 0, 0, 0.38);
  backdrop-filter: blur(12px);
}

.field {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  font-weight: 600;
}

.input {
  padding: 0.7rem 0.85rem;
  border-radius: 10px;
  border: 1px solid rgba(255, 255, 255, 0.3);
  background: rgba(12, 24, 38, 0.85);
  color: inherit;
  transition: border-color 0.2s ease, box-shadow 0.2s ease;
}

.input:focus {
  outline: none;
  border-color: #42b883;
  box-shadow: 0 0 0 3px rgba(66, 184, 131, 0.2);
}

.button-group {
  display: flex;
  gap: 0.75rem;
}

.submit,
.clear-button {
  flex: 1;
  padding: 0.8rem;
  border-radius: 10px;
  border: none;
  font-weight: 700;
  cursor: pointer;
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.submit {
  background: linear-gradient(135deg, #42b883, #2d88ff);
  color: #0b1729;
}

.clear-button {
  background: linear-gradient(135deg, #ff6b6b, #ee5a6f);
  color: #ffffff;
}

.submit:disabled,
.clear-button:disabled {
  opacity: 0.7;
  cursor: not-allowed;
  box-shadow: none;
}

.submit:not(:disabled):hover {
  transform: translateY(-1px);
  box-shadow: 0 10px 20px rgba(45, 136, 255, 0.35);
}

.clear-button:not(:disabled):hover {
  transform: translateY(-1px);
  box-shadow: 0 10px 20px rgba(238, 90, 111, 0.35);
}

.status {
  text-align: center;
  font-size: 0.95rem;
  padding: 0.75rem 1rem;
  border-radius: 10px;
  background: rgba(255, 255, 255, 0.08);
  box-shadow: 0 10px 20px rgba(0, 0, 0, 0.25);
}

.status--success {
  color: #8ef5c3;
}

.status--error {
  color: #ffb8b0;
}
</style>
