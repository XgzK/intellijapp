/**
 * 更新检查服务
 * 通过后端 Go 服务检查 GitHub Releases 获取最新版本信息
 */

import { CheckForUpdates } from '../../bindings/github.com/XgzK/intellijapp/internal/service/configservice'

export interface ReleaseInfo {
  version: string
  publishedAt: string
  htmlUrl: string
  body: string
  assets: Array<{
    name: string
    downloadUrl: string
    size: number
  }>
}

/**
 * 检查是否有新版本
 * 调用后端 Go 服务进行更新检查（后端自动使用本地版本号）
 */
export async function checkForUpdates(): Promise<{
  hasUpdate: boolean
  release: ReleaseInfo | null
}> {
  try {
    const result = await CheckForUpdates()
    return {
      hasUpdate: result.hasUpdate,
      release: result.release,
    }
  } catch (error) {
    console.error('检查更新失败:', error)
    return { hasUpdate: false, release: null }
  }
}

/**
 * 格式化文件大小
 */
export function formatFileSize(bytes: number): string {
  if (bytes === 0) return '0 B'

  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))

  return `${(bytes / Math.pow(k, i)).toFixed(2)} ${sizes[i]}`
}

/**
 * 格式化日期
 */
export function formatDate(dateString: string): string {
  const date = new Date(dateString)
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  })
}
