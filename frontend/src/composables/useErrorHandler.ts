/**
 * 全局错误处理 Composable
 * 优化：统一的错误处理逻辑，提供更好的用户体验
 */
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'

export interface AppError {
  message: string
  code?: string
  timestamp: number
  stack?: string
}

const errors = ref<AppError[]>([])
const maxErrors = 10 // 最多保留10个错误

/**
 * 错误处理 Hook
 */
export function useErrorHandler() {
  const t = (() => {
    try {
      return useI18n().t
    } catch {
      return (key: string) => (key === 'errors.unknown' ? '未知错误' : key)
    }
  })()

  /**
   * 处理错误
   */
  const handleError = (error: unknown, context?: string): string => {
    const errorMessage = extractErrorMessage(error, context)

    // 记录错误
    const appError: AppError = {
      message: errorMessage,
      timestamp: Date.now(),
      stack: error instanceof Error ? error.stack : undefined,
    }

    // 添加到错误列表（保持最多 maxErrors 个）
    errors.value = [appError, ...errors.value].slice(0, maxErrors)

    // 在开发环境输出详细错误
    if (import.meta.env.DEV) {
      console.error(`[Error Handler]${context ? ` [${context}]` : ''}:`, error)
    }

    return errorMessage
  }

  /**
   * 提取错误消息
   */
  const extractErrorMessage = (error: unknown, _context?: string): string => {
    // 处理字符串错误
    if (typeof error === 'string') {
      return parseErrorString(error)
    }

    // 处理 Error 对象
    if (error instanceof Error) {
      return parseErrorString(error.message)
    }

    // 处理其他类型
    return parseErrorString(String(error))
  }

  /**
   * 解析错误字符串
   */
  const parseErrorString = (errorStr: string): string => {
    const trimmed = errorStr.trim()

    if (!trimmed) {
      return t('errors.unknown')
    }

    // 提取 desc= 后的内容
    const descIndex = trimmed.lastIndexOf('desc=')
    if (descIndex !== -1) {
      return trimmed.substring(descIndex + 5).trim()
    }

    // 提取最后一个冒号后的内容
    const colonIndex = trimmed.lastIndexOf(':')
    if (colonIndex !== -1) {
      const extracted = trimmed.substring(colonIndex + 1).trim()
      if (extracted) {
        return extracted
      }
    }

    return trimmed
  }

  /**
   * 清除所有错误
   */
  const clearErrors = () => {
    errors.value = []
  }

  /**
   * 清除特定错误
   */
  const clearError = (timestamp: number) => {
    errors.value = errors.value.filter(e => e.timestamp !== timestamp)
  }

  /**
   * 获取最新错误
   */
  const getLatestError = (): AppError | null => {
    return errors.value[0] || null
  }

  return {
    errors,
    handleError,
    clearErrors,
    clearError,
    getLatestError,
  }
}

/**
 * 创建全局错误处理器实例
 */
let globalErrorHandler: ReturnType<typeof useErrorHandler> | null = null

export function setupGlobalErrorHandler() {
  if (!globalErrorHandler) {
    globalErrorHandler = useErrorHandler()
  }
  return globalErrorHandler
}

export function getGlobalErrorHandler() {
  if (!globalErrorHandler) {
    throw new Error('Global error handler not initialized. Call setupGlobalErrorHandler first.')
  }
  return globalErrorHandler
}
