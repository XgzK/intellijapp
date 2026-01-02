/**
 * MainView 组件单元测试
 * 测试路径验证、表单提交和配置清除功能
 */
import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import MainView from './MainView.vue'
import { CONFIG_PATH_PATTERN } from '@/constants/validation'

// 模拟服务
vi.mock('@/services/configService', () => ({
  SubmitPaths: vi.fn().mockResolvedValue('配置应用成功'),
  ClearConfig: vi.fn().mockResolvedValue('配置清除成功'),
  PathExists: vi.fn().mockResolvedValue(true),
}))

// 模拟 vue-i18n
vi.mock('vue-i18n', () => ({
  useI18n: () => ({
    t: (key: string) => key,
  }),
}))

describe('MainView 组件', () => {
  let wrapper: ReturnType<typeof mount>

  beforeEach(() => {
    wrapper = mount(MainView)
  })

  it('应该正确渲染', () => {
    expect(wrapper.exists()).toBe(true)
    expect(wrapper.find('.workspace').exists()).toBe(true)
  })

  it('应该包含两个输入框', () => {
    const inputs = wrapper.findAll('input')
    expect(inputs).toHaveLength(2)
  })

  it('应该包含两个按钮（应用配置和清除配置）', () => {
    const buttons = wrapper.findAll('button')
    expect(buttons.length).toBeGreaterThanOrEqual(2)
  })

  describe('路径验证', () => {
    it('CONFIG_PATH_PATTERN 应该匹配有效路径', () => {
      const validPaths = [
        'D:/jetbra',
        'C:/Program Files/Config',
        'D:/test_folder-123',
        '/home/user/config',
      ]

      validPaths.forEach(path => {
        expect(CONFIG_PATH_PATTERN.test(path)).toBe(true)
      })
    })

    it('CONFIG_PATH_PATTERN 应该拒绝无效路径', () => {
      const invalidPaths = ['D:/folder@invalid', 'C:/test#123', 'D:/test&path']

      invalidPaths.forEach(path => {
        expect(CONFIG_PATH_PATTERN.test(path)).toBe(false)
      })
    })
  })

  describe('表单交互', () => {
    it('应该能够输入项目路径', async () => {
      const input = wrapper.findAll('input')[0]
      await input.setValue('D:/IntelliJ IDEA')

      expect((input.element as HTMLInputElement).value).toBe('D:/IntelliJ IDEA')
    })

    it('应该能够输入配置路径', async () => {
      const input = wrapper.findAll('input')[1]
      await input.setValue('D:/jetbra')

      expect((input.element as HTMLInputElement).value).toBe('D:/jetbra')
    })

    it('空路径时应该显示错误', async () => {
      const form = wrapper.find('form')
      await form.trigger('submit')

      // 等待 DOM 更新
      await wrapper.vm.$nextTick()

      // 应该显示错误消息
      expect(wrapper.text()).toContain('路径')
    })
  })

  describe('样式类', () => {
    it('应该有正确的 CSS 类', () => {
      expect(wrapper.find('.workspace').exists()).toBe(true)
      expect(wrapper.find('.panel').exists()).toBe(true)
      expect(wrapper.find('.panel--path').exists()).toBe(true)
    })

    it('按钮应该有正确的样式类', () => {
      const primaryButton = wrapper.find('.button--primary')
      const dangerButton = wrapper.find('.button--danger')

      expect(primaryButton.exists()).toBe(true)
      expect(dangerButton.exists()).toBe(true)
    })
  })
})
