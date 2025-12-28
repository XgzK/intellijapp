<script setup lang="ts">
import { Window } from '@wailsio/runtime'

interface Props {
  title: string
  currentView: 'main' | 'about'
}

interface Emits {
  (e: 'switch-view', view: 'main' | 'about'): void
}

defineProps<Props>()
const emit = defineEmits<Emits>()

const handleClose = async () => {
  try {
    await Window.Close()
  } catch (error) {
    console.error('关闭窗口失败', error)
  }
}

const handleMinimize = async () => {
  try {
    await Window.Minimise()
  } catch (error) {
    console.error('最小化窗口失败', error)
  }
}

const handleZoom = async () => {
  try {
    const maximised = await Window.IsMaximised()
    if (maximised) {
      await Window.Restore()
    } else {
      await Window.Maximise()
    }
  } catch (error) {
    console.error('切换窗口大小失败', error)
  }
}

const switchView = (view: 'main' | 'about') => {
  emit('switch-view', view)
}
</script>

<template>
  <header class="titlebar">
    <div class="traffic-lights">
      <button class="traffic-button close" type="button" @click.stop="handleClose" />
      <button class="traffic-button minimise" type="button" @click.stop="handleMinimize" />
      <button class="traffic-button zoom" type="button" @click.stop="handleZoom" />
    </div>
    <div class="title">{{ title }}</div>
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
</template>

<style scoped>
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
</style>
