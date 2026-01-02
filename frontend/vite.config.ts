import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import VueI18nPlugin from '@intlify/unplugin-vue-i18n/vite'
import { visualizer } from 'rollup-plugin-visualizer'
import { fileURLToPath, URL } from 'node:url'
import { resolve, dirname } from 'node:path'

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    VueI18nPlugin({
      // 语言包文件路径
      include: resolve(dirname(fileURLToPath(import.meta.url)), './src/i18n/locales/**'),
      // 使用 Composition API 模式
      compositionOnly: true,
      // 启用运行时编译（用于消息格式化）
      runtimeOnly: false,
      // 允许翻译文本中包含 HTML 标签
      strictMessage: false,
    }),
    // 打包分析工具
    visualizer({
      filename: './dist/stats.html',
      open: false,
      gzipSize: true,
      brotliSize: true,
    }),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },
})
