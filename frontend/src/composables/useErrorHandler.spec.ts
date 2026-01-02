/**
 * useErrorHandler Composable 单元测试
 * 测试错误处理、错误提取和错误历史功能
 */
import { describe, it, expect, beforeEach, vi } from 'vitest'
import { useErrorHandler } from './useErrorHandler'

// 模拟 vue-i18n
vi.mock('vue-i18n', () => ({
  useI18n: () => ({
    t: (key: string) => key,
  }),
}))

describe('useErrorHandler', () => {
  let errorHandler: ReturnType<typeof useErrorHandler>

  beforeEach(() => {
    errorHandler = useErrorHandler()
    errorHandler.clearErrors()
  })

  describe('handleError', () => {
    it('应该处理字符串错误', () => {
      const errorMessage = '测试错误'
      const result = errorHandler.handleError(errorMessage)

      expect(result).toBe(errorMessage)
      expect(errorHandler.errors.value).toHaveLength(1)
      expect(errorHandler.errors.value[0].message).toBe(errorMessage)
    })

    it('应该处理 Error 对象', () => {
      const error = new Error('Error 对象测试')
      const result = errorHandler.handleError(error)

      expect(result).toBe('Error 对象测试')
      expect(errorHandler.errors.value).toHaveLength(1)
      expect(errorHandler.errors.value[0].stack).toBeDefined()
    })

    it('应该提取 desc= 后的内容', () => {
      const errorStr = 'Error: code=500 desc=服务器错误'
      const result = errorHandler.handleError(errorStr)

      expect(result).toBe('服务器错误')
    })

    it('应该提取最后一个冒号后的内容', () => {
      const errorStr = 'Error: Network: Connection failed: 连接超时'
      const result = errorHandler.handleError(errorStr)

      expect(result).toBe('连接超时')
    })

    it('应该保持最多10个错误', () => {
      // 添加 11 个错误
      for (let i = 0; i < 11; i++) {
        errorHandler.handleError(`错误 ${i}`)
      }

      expect(errorHandler.errors.value).toHaveLength(10)
      expect(errorHandler.errors.value[0].message).toBe('错误 10') // 最新的错误
    })

    it('应该处理空字符串错误', () => {
      const result = errorHandler.handleError('')

      expect(result).toBe('errors.unknown')
    })
  })

  describe('clearErrors', () => {
    it('应该清除所有错误', () => {
      errorHandler.handleError('错误 1')
      errorHandler.handleError('错误 2')

      expect(errorHandler.errors.value).toHaveLength(2)

      errorHandler.clearErrors()

      expect(errorHandler.errors.value).toHaveLength(0)
    })
  })

  describe('clearError', () => {
    it('应该清除指定时间戳的错误', async () => {
      // 确保从空状态开始
      errorHandler.clearErrors()

      errorHandler.handleError('错误 1')
      // 添加小延迟确保 timestamp 不同
      await new Promise(resolve => setTimeout(resolve, 10))
      errorHandler.handleError('错误 2')

      // 验证有两个错误
      expect(errorHandler.errors.value).toHaveLength(2)

      // 最新的错误在数组的第一位，最旧的在最后
      const timestampToRemove = errorHandler.errors.value[1].timestamp

      errorHandler.clearError(timestampToRemove)

      expect(errorHandler.errors.value).toHaveLength(1)
      expect(errorHandler.errors.value[0].message).toBe('错误 2')
    })
  })

  describe('getLatestError', () => {
    it('应该返回最新的错误', () => {
      errorHandler.handleError('旧错误')
      errorHandler.handleError('新错误')

      const latestError = errorHandler.getLatestError()

      expect(latestError).not.toBeNull()
      expect(latestError?.message).toBe('新错误')
    })

    it('没有错误时应该返回 null', () => {
      const latestError = errorHandler.getLatestError()

      expect(latestError).toBeNull()
    })
  })
})
