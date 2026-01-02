<script setup lang="ts">
/**
 * 更新通知组件
 * 显示新版本可用的通知
 * 后端自动使用本地版本号进行比较
 */
import { ref, onMounted, computed } from 'vue'
import { Browser } from '@wailsio/runtime'
import { checkForUpdates, formatDate, type ReleaseInfo } from '@/services/updateService'
import { ConvertToAccessibleURL } from '@/services/configService'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

interface Props {
  checkOnMount?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  checkOnMount: true,
})

const hasUpdate = ref(false)
const release = ref<ReleaseInfo | null>(null)
const isChecking = ref(false)
const showNotification = ref(false)
const dismissed = ref(false)

const checkUpdate = async () => {
  if (isChecking.value) return

  isChecking.value = true
  try {
    const result = await checkForUpdates()
    hasUpdate.value = result.hasUpdate
    release.value = result.release

    if (result.hasUpdate && !dismissed.value) {
      showNotification.value = true
    }
  } catch (error) {
    console.error(t('update.checkingFailed'), error)
  } finally {
    isChecking.value = false
  }
}

const openReleasePage = async () => {
  if (release.value?.htmlUrl) {
    try {
      // 转换为可访问的镜像 URL
      const accessibleUrl = await ConvertToAccessibleURL(release.value.htmlUrl)
      await Browser.OpenURL(accessibleUrl)
    } catch (error) {
      console.error('打开 Release 页面失败:', error)
    }
  }
}

const dismissNotification = () => {
  showNotification.value = false
  dismissed.value = true
}

const releaseSummary = computed(() => {
  const body = release.value?.body?.trim() || ''
  if (!body) return ''
  return body.length > 100 ? `${body.substring(0, 100)}...` : body
})

onMounted(() => {
  if (props.checkOnMount) {
    checkUpdate()
  }
})

defineExpose({
  checkUpdate,
})
</script>

<template>
  <Transition name="slide-down">
    <div v-if="showNotification && hasUpdate && release" class="update-notification">
      <div class="update-content">
        <div class="update-icon">🎉</div>
        <div class="update-info">
          <h3 class="update-title">{{ $t('update.available') }}</h3>
          <p class="update-version">
            v{{ release.version }}
            <span class="update-date">· {{ formatDate(release.publishedAt) }}</span>
          </p>
          <p v-if="releaseSummary" class="update-description">
            {{ releaseSummary }}
          </p>
        </div>
        <div class="update-actions">
          <button
            class="update-button update-button--primary"
            type="button"
            @click="openReleasePage"
          >
            {{ $t('update.viewDetails') }}
          </button>
          <button
            class="update-button update-button--secondary"
            type="button"
            @click="dismissNotification"
          >
            {{ $t('update.remindLater') }}
          </button>
        </div>
      </div>
      <button class="update-close" type="button" @click="dismissNotification">×</button>
    </div>
  </Transition>
</template>

<style scoped>
.update-notification {
  position: fixed;
  top: 60px;
  left: 50%;
  transform: translateX(-50%);
  width: calc(100% - 40px);
  max-width: 600px;
  background: var(--color-surface-strong);
  border: 1px solid var(--color-border-strong);
  border-radius: 12px;
  padding: 1.25rem;
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.4);
  z-index: 9999;
  backdrop-filter: blur(12px);
}

.update-content {
  display: flex;
  gap: 1rem;
  align-items: flex-start;
}

.update-icon {
  font-size: 2rem;
  flex-shrink: 0;
}

.update-info {
  flex: 1;
  min-width: 0;
}

.update-title {
  font-size: 1.1rem;
  font-weight: 700;
  color: var(--color-text);
  margin-bottom: 0.25rem;
}

.update-version {
  font-size: 0.875rem;
  color: var(--color-accent);
  font-weight: 600;
  margin-bottom: 0.5rem;
}

.update-date {
  color: var(--color-muted);
  font-weight: 400;
}

.update-description {
  font-size: 0.8rem;
  color: var(--color-muted);
  line-height: 1.5;
  margin-top: 0.5rem;
}

.update-actions {
  display: flex;
  gap: 0.5rem;
  margin-top: 1rem;
  flex-shrink: 0;
}

.update-button {
  padding: 0.5rem 1rem;
  border-radius: 8px;
  border: 1px solid transparent;
  font-size: 0.8rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
  white-space: nowrap;
}

.update-button--primary {
  background: var(--color-accent);
  color: #fff;
  box-shadow: 0 4px 12px rgba(66, 184, 131, 0.3);
}

.update-button--primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(66, 184, 131, 0.4);
}

.update-button--secondary {
  background: transparent;
  color: var(--color-muted);
  border-color: var(--color-border);
}

.update-button--secondary:hover {
  background: rgba(255, 255, 255, 0.05);
  color: var(--color-text);
}

.update-close {
  position: absolute;
  top: 0.75rem;
  right: 0.75rem;
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: none;
  background: transparent;
  color: var(--color-muted);
  font-size: 1.5rem;
  line-height: 1;
  cursor: pointer;
  border-radius: 4px;
  transition: all 0.2s ease;
}

.update-close:hover {
  background: rgba(255, 255, 255, 0.1);
  color: var(--color-text);
}

/* 过渡动画 */
.slide-down-enter-active,
.slide-down-leave-active {
  transition: all 0.4s ease;
}

.slide-down-enter-from {
  opacity: 0;
  transform: translate(-50%, -30px);
}

.slide-down-leave-to {
  opacity: 0;
  transform: translate(-50%, -10px);
}

/* 响应式 */
@media (max-width: 640px) {
  .update-content {
    flex-direction: column;
  }

  .update-actions {
    flex-direction: column;
  }

  .update-button {
    width: 100%;
  }
}
</style>
