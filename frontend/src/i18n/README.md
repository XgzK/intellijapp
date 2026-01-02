# i18n 国际化使用指南

## 📦 安装依赖

配置文件已经更新，请运行以下命令安装 vue-i18n：

```bash
cd frontend
npm install
```

## 🎯 快速开始

### 1. 基础配置（已完成）

i18n 已经在 `main.ts` 中注册，无需额外配置。

### 2. 在组件中使用

#### 方式一：使用 Composition API（推荐）

```vue
<script setup lang="ts">
import { useI18n } from 'vue-i18n'

const { t, locale } = useI18n()

// 切换语言
const switchLanguage = () => {
  locale.value = locale.value === 'zh-CN' ? 'en-US' : 'zh-CN'
}
</script>

<template>
  <div>
    <!-- 使用翻译 -->
    <h1>{{ t('common.appName') }}</h1>
    <button @click="switchLanguage">
      {{ t('navigation.about') }}
    </button>
  </div>
</template>
```

#### 方式二：使用全局 $t 函数

```vue
<template>
  <div>
    <h1>{{ $t('common.appName') }}</h1>
    <p>{{ $t('mainView.installPath.description') }}</p>
  </div>
</template>
```

### 3. 带参数的翻译

```vue
<script setup lang="ts">
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const count = 5

// 使用参数
const message = t('mainView.applyConfig.successMessage', { count })
</script>
```

### 4. 切换语言

```typescript
import { useI18n } from 'vue-i18n'

const { locale } = useI18n()

// 切换到英文
locale.value = 'en-US'

// 切换到中文
locale.value = 'zh-CN'
```

## 📚 语言包结构

### 中文语言包 (zh-CN.ts)

```typescript
export default {
  common: {
    appName: 'IntelliJ 配置助手',
    loading: '加载中...',
  },
  // ...
}
```

### 英文语言包 (en-US.ts)

```typescript
export default {
  common: {
    appName: 'IntelliJ Config Helper',
    loading: 'Loading...',
  },
  // ...
}
```

## 🔧 添加新语言

1. 在 `locales/` 目录下创建新的语言文件，如 `ja-JP.ts`
2. 在 `i18n/index.ts` 中导入并注册：

```typescript
import jaJP from './locales/ja-JP'

const messages = {
  'zh-CN': zhCN,
  'en-US': enUS,
  'ja-JP': jaJP, // 新增
}
```

3. 更新 `SupportedLocale` 类型：

```typescript
export type SupportedLocale = 'zh-CN' | 'en-US' | 'ja-JP'
```

## 📖 翻译键参考

### 常用翻译键

| 键                                  | 中文                     | 英文                    |
| ----------------------------------- | ------------------------ | ----------------------- |
| `common.appName`                    | IntelliJ 配置助手        | IntelliJ Config Helper  |
| `navigation.main`                   | 主页                     | Main                    |
| `navigation.about`                  | 关于                     | About                   |
| `mainView.applyConfig.submitButton` | 应用配置                 | Apply Config            |
| `validation.emptyPaths`             | 请输入完整的两个路径喵～ | Please enter both paths |

### 完整翻译键列表

请查看 `locales/zh-CN.ts` 或 `locales/en-US.ts` 文件。

## 🎨 最佳实践

1. **始终使用翻译键**：不要在组件中硬编码文本
2. **命名规范**：使用层级结构命名，如 `mainView.applyConfig.submitButton`
3. **参数化消息**：对于包含变量的消息，使用参数替换
4. **回退语言**：设置 `fallbackLocale` 确保缺失翻译时有默认值
5. **类型安全**：使用 TypeScript 确保翻译键的类型安全

## 🚀 高级功能

### 复数形式

```typescript
// 语言包
{
  apple: 'no apples | one apple | {count} apples'
}

// 使用
t('apple', 0) // "no apples"
t('apple', 1) // "one apple"
t('apple', 10) // "10 apples"
```

### 日期和数字格式化

```typescript
import { useI18n } from 'vue-i18n'

const { d, n } = useI18n()

d(new Date(), 'short') // 日期格式化
n(1000.5, 'currency') // 数字格式化
```

## 📝 注意事项

- 语言包文件使用 `.ts` 而非 `.json`，以获得更好的类型提示
- 翻译键区分大小写
- 当前实现使用 Composition API 模式（`legacy: false`）
- 默认语言根据浏览器语言自动检测，回退到中文

---

浮浮酱制作 ฅ'ω'ฅ
