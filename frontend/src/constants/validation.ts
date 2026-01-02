/**
 * 验证相关常量配置
 * 优化：将硬编码的正则表达式和验证规则提取为常量，提高可维护性
 */

/**
 * 配置路径验证正则表达式
 * 允许字母、数字、冒号、斜杠、反斜杠、空格、下划线、连字符和点号
 */
export const CONFIG_PATH_PATTERN = /^[A-Za-z0-9:\\/\s_\-.]+$/

/**
 * 验证错误消息
 */
export const VALIDATION_MESSAGES = {
  EMPTY_PATHS: '请输入完整的两个路径喵～',
  INVALID_CONFIG_PATH: '配置目录路径包含不支持的特殊字符喵～',
  INTELLIJ_PATH_NOT_EXIST: 'IntelliJ 安装路径',
  CONFIG_PATH_NOT_EXIST: '配置目录',
  PATH_NOT_EXIST_SUFFIX: ' 不存在，请检查喵～',
  EMPTY_INTELLIJ_PATH: '请输入 IntelliJ 路径',
} as const

/**
 * 路径标签
 */
export const PATH_LABELS = {
  INTELLIJ_PATH: 'IntelliJ 安装路径',
  CONFIG_PATH: '配置目录',
} as const
