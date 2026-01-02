<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { GetAboutInfo } from '@/services/configService'
import { AboutInfo } from '../bindings/github.com/XgzK/intellijapp/internal/service/models'
import TitleBar from '@/components/TitleBar.vue'
import MainView from '@/components/MainView.vue'
import AboutView from '@/components/AboutView.vue'
import ErrorBoundary from '@/components/ErrorBoundary.vue'
import UpdateNotification from '@/components/UpdateNotification.vue'
import { useErrorHandler } from '@/composables/useErrorHandler'

const currentView = ref<'main' | 'about'>('main')
const { handleError } = useErrorHandler()

// 优化：使用对象字面量初始化，而非 new AboutInfo()
// AboutInfo 是从 Go 生成的类型定义，应提供合理的默认值
const aboutInfo = ref<AboutInfo>({
  appName: '',
  version: '',
  goVersion: '',
  vueVersion: '',
  wailsVersion: '',
  repoUrl: '',
  developers: [],
})

const switchView = (view: 'main' | 'about') => {
  currentView.value = view
}

// 组件挂载时获取关于信息
onMounted(async () => {
  try {
    aboutInfo.value = await GetAboutInfo()
  } catch (error) {
    handleError(error, 'GetAboutInfo')
  }
})
</script>

<template>
  <div class="window">
    <UpdateNotification />
    <TitleBar :title="$t('titleBar.title')" :current-view="currentView" @switch-view="switchView" />
    <ErrorBoundary>
      <MainView v-if="currentView === 'main'" />
      <AboutView v-else :about-info="aboutInfo" />
    </ErrorBoundary>
  </div>
</template>

<style scoped>
.window {
  /* 间距变量 */
  --space-xs: 0.3rem;
  --space-sm: 0.5rem;
  --space-md: 0.75rem;
  --space-lg: 1.1rem;
  --space-xl: 1.6rem;

  height: 100vh;
  overflow: hidden;
  /* 背景已在 global.css 中统一设置 */
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
