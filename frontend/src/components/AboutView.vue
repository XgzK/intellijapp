<script setup lang="ts">
import { Browser } from '@wailsio/runtime'
import type { AboutInfo } from '../../bindings/github.com/XgzK/intellijapp/internal/service/models'

interface Props {
  aboutInfo: AboutInfo
}

defineProps<Props>()

const openExternal = async (url: string) => {
  try {
    await Browser.OpenURL(url)
  } catch (error) {
    console.error('打开外部链接失败', url, error)
  }
}
</script>

<template>
  <main class="workspace">
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
            <span class="meta-value">{{ aboutInfo.appName }}</span>
          </div>
          <div class="meta-row">
            <span class="meta-label">当前版本</span>
            <span class="meta-value">{{ aboutInfo.version }}</span>
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
              Wails {{ aboutInfo.wailsVersion }}
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
              <span class="meta-value">Vue {{ aboutInfo.vueVersion }}</span>
            </div>
            <div class="meta-row">
              <span class="meta-label">后端语言</span>
              <span class="meta-value">{{ aboutInfo.goVersion }}</span>
            </div>
            <div class="meta-row">
              <span class="meta-label">构建工具</span>
              <span class="meta-value">Wails {{ aboutInfo.wailsVersion }}</span>
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
                :href="`https://${aboutInfo.repoUrl}`"
                target="_blank"
                rel="noopener noreferrer"
                class="meta-link"
                @click.prevent="openExternal(`https://${aboutInfo.repoUrl}`)"
              >
                {{ aboutInfo.repoUrl }}
              </a>
            </div>
            <div class="meta-row">
              <span class="meta-label">开发者</span>
              <span class="meta-value meta-value--developers">
                <template v-for="(dev, index) in aboutInfo.developers" :key="dev.name">
                  <span v-if="index > 0"> · </span>
                  <a
                    class="meta-link"
                    :href="dev.url"
                    target="_blank"
                    rel="noopener noreferrer"
                    @click.prevent="openExternal(dev.url)"
                  >
                    {{ dev.name }}
                  </a>
                </template>
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

.actions-grid {
  display: grid;
  gap: var(--space-lg);
  grid-template-columns: repeat(auto-fit, minmax(260px, 1fr));
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
  .workspace {
    padding: var(--space-lg);
  }

  .panel {
    padding: var(--space-md);
  }
}
</style>
