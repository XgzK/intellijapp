<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { GetAboutInfo } from '@/services/configService'
import { AboutInfo } from '../bindings/github.com/XgzK/intellijapp/internal/service/models'
import TitleBar from '@/components/TitleBar.vue'
import MainView from '@/components/MainView.vue'
import AboutView from '@/components/AboutView.vue'

const currentView = ref<'main' | 'about'>('main')
const aboutInfo = ref<AboutInfo>(new AboutInfo())

const switchView = (view: 'main' | 'about') => {
  currentView.value = view
}

// 组件挂载时获取关于信息
onMounted(async () => {
  try {
    aboutInfo.value = await GetAboutInfo()
  } catch (error) {
    console.error('获取关于信息失败', error)
  }
})
</script>

<template>
  <div class="window">
    <TitleBar
      title="IntelliJ 配置助手"
      :current-view="currentView"
      @switch-view="switchView"
    />
    <MainView v-if="currentView === 'main'" />
    <AboutView v-else :about-info="aboutInfo" />
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
</style>
