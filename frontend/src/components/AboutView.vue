<script setup lang="ts">
import { Browser } from '@wailsio/runtime'
import { ConvertToAccessibleURL } from '@/services/configService'
import type { AboutInfo } from '../../bindings/github.com/XgzK/intellijapp/internal/service'

interface Props {
  aboutInfo: AboutInfo
}

defineProps<Props>()

const normalizeUrl = (url: string) => {
  if (/^https?:\/\//i.test(url)) {
    return url
  }
  return `https://${url}`
}

const openExternal = async (url: string) => {
  try {
    const normalizedUrl = normalizeUrl(url)
    // 如果是 GitHub 链接，转换为可访问的镜像 URL
    const accessibleUrl = await ConvertToAccessibleURL(normalizedUrl)
    await Browser.OpenURL(accessibleUrl)
  } catch (error) {
    console.error('打开外部链接失败', url, error)
  }
}
</script>

<template>
  <main class="workspace">
    <section class="panel panel--path">
      <header class="panel__header">
        <h2 class="panel__title">{{ $t('aboutView.main.title') }}</h2>
      </header>
      <div class="panel__body">
        <div class="meta-list">
          <div class="meta-row">
            <span class="meta-label">{{ $t('aboutView.main.appName') }}</span>
            <span class="meta-value">{{ aboutInfo.appName }}</span>
          </div>
          <div class="meta-row">
            <span class="meta-label">{{ $t('aboutView.main.version') }}</span>
            <span class="meta-value">{{ aboutInfo.version }}</span>
          </div>
          <div class="meta-row meta-row--link">
            <span class="meta-label">{{ $t('aboutView.main.buildTool') }}</span>
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
        <div class="panel__body">
          <div class="meta-list">
            <div class="meta-row">
              <span class="meta-label">{{ $t('aboutView.techStack.frontend') }}</span>
              <span class="meta-value">Vue {{ aboutInfo.vueVersion }}</span>
            </div>
            <div class="meta-row">
              <span class="meta-label">{{ $t('aboutView.techStack.backend') }}</span>
              <span class="meta-value">{{ aboutInfo.goVersion }}</span>
            </div>
            <div class="meta-row">
              <span class="meta-label">{{ $t('aboutView.techStack.buildTool') }}</span>
              <span class="meta-value">Wails {{ aboutInfo.wailsVersion }}</span>
            </div>
          </div>
        </div>
      </article>

      <article class="panel panel--action">
        <div class="panel__body">
          <div class="meta-list">
            <div class="meta-row meta-row--link">
              <span class="meta-label">{{ $t('aboutView.project.repo') }}</span>
              <a
                :href="normalizeUrl(aboutInfo.repoUrl)"
                target="_blank"
                rel="noopener noreferrer"
                class="meta-link"
                @click.prevent="openExternal(aboutInfo.repoUrl)"
              >
                {{ aboutInfo.repoUrl }}
              </a>
            </div>
            <div class="meta-row">
              <span class="meta-label">{{ $t('aboutView.project.developers') }}</span>
              <span class="meta-value meta-value--developers">
                <template v-for="(dev, index) in aboutInfo.developers" :key="dev.name">
                  <span v-if="index > 0">·</span>
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
        {{ $t('aboutView.notice') }}
      </p>
    </section>

    <p class="about-footer">{{ $t('aboutView.footer') }}</p>
  </main>
</template>

<style scoped>
/* 组件特定样式 - 公共样式已移至 global.css */

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
  background: var(--input-background);
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
  transition:
    color 0.2s ease,
    background 0.2s ease,
    transform 0.2s ease,
    box-shadow 0.2s ease;
  font-family: 'Cascadia Code', 'Fira Code', 'Courier New', monospace;
  word-break: break-all;
  background: var(--input-background);
  padding: 0.35rem 0.5rem;
  border-radius: 6px;
  display: inline-flex;
}

.meta-link:hover {
  color: var(--color-accent);
  background: var(--input-background-hover);
  transform: translateY(-1px);
  box-shadow: 0 6px 14px rgba(45, 136, 255, 0.2);
}

.meta-link:active {
  transform: translateY(0);
}

.about-notice {
  font-size: 0.85rem;
  color: var(--about-notice-text);
  line-height: 1.6;
  text-align: center;
}

.about-footer {
  text-align: center;
  margin-top: var(--space-lg);
  padding-top: var(--space-sm);
  border-top: 1px solid rgba(255, 255, 255, 0.08);
  color: var(--about-footer-text);
  font-size: 0.75rem;
}
</style>
