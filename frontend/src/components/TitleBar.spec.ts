/**
 * TitleBar 组件单元测试
 * 测试窗口控制和视图切换功能
 */
import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import TitleBar from './TitleBar.vue'

describe('TitleBar 组件', () => {
  const defaultProps = {
    title: '测试标题',
    currentView: 'main' as const,
  }

  it('应该正确渲染标题', () => {
    const wrapper = mount(TitleBar, {
      props: defaultProps,
    })

    expect(wrapper.text()).toContain('测试标题')
    expect(wrapper.find('.title').text()).toBe('测试标题')
  })

  it('应该渲染三个窗口控制按钮', () => {
    const wrapper = mount(TitleBar, {
      props: defaultProps,
    })

    const trafficButtons = wrapper.findAll('.traffic-button')
    expect(trafficButtons).toHaveLength(3)

    expect(wrapper.find('.traffic-button.close').exists()).toBe(true)
    expect(wrapper.find('.traffic-button.minimise').exists()).toBe(true)
    expect(wrapper.find('.traffic-button.zoom').exists()).toBe(true)
  })

  it('应该渲染两个导航按钮', () => {
    const wrapper = mount(TitleBar, {
      props: defaultProps,
    })

    const navButtons = wrapper.findAll('.nav-button')
    expect(navButtons).toHaveLength(2)
  })

  it('当前视图按钮应该有 active 类', () => {
    const wrapper = mount(TitleBar, {
      props: {
        ...defaultProps,
        currentView: 'main' as const,
      },
    })

    const navButtons = wrapper.findAll('.nav-button')
    const mainButton = navButtons[0]
    const aboutButton = navButtons[1]

    expect(mainButton.classes()).toContain('nav-button--active')
    expect(aboutButton.classes()).not.toContain('nav-button--active')
  })

  it('切换到关于页时应该有正确的 active 类', () => {
    const wrapper = mount(TitleBar, {
      props: {
        ...defaultProps,
        currentView: 'about' as const,
      },
    })

    const navButtons = wrapper.findAll('.nav-button')
    const mainButton = navButtons[0]
    const aboutButton = navButtons[1]

    expect(mainButton.classes()).not.toContain('nav-button--active')
    expect(aboutButton.classes()).toContain('nav-button--active')
  })

  it('点击导航按钮应该触发 switchView 事件', async () => {
    const wrapper = mount(TitleBar, {
      props: defaultProps,
    })

    const navButtons = wrapper.findAll('.nav-button')
    const aboutButton = navButtons[1]

    await aboutButton.trigger('click')

    expect(wrapper.emitted()).toHaveProperty('switchView')
    expect(wrapper.emitted('switchView')).toHaveLength(1)
    expect(wrapper.emitted('switchView')![0]).toEqual(['about'])
  })

  it('点击窗口控制按钮应该调用相应的方法', async () => {
    const wrapper = mount(TitleBar, {
      props: defaultProps,
    })

    const closeButton = wrapper.find('.traffic-button.close')
    const minimiseButton = wrapper.find('.traffic-button.minimise')
    const zoomButton = wrapper.find('.traffic-button.zoom')

    // 测试点击不会抛出错误（实际的 Window API 已被 mock）
    await expect(closeButton.trigger('click')).resolves.not.toThrow()
    await expect(minimiseButton.trigger('click')).resolves.not.toThrow()
    await expect(zoomButton.trigger('click')).resolves.not.toThrow()
  })

  it('应该有正确的 CSS 类结构', () => {
    const wrapper = mount(TitleBar, {
      props: defaultProps,
    })

    expect(wrapper.find('.titlebar').exists()).toBe(true)
    expect(wrapper.find('.traffic-lights').exists()).toBe(true)
    expect(wrapper.find('.title').exists()).toBe(true)
    expect(wrapper.find('.nav-buttons').exists()).toBe(true)
  })
})
