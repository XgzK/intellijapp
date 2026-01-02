/**
 * 中文（简体）语言包
 * 优化：添加国际化支持，便于未来扩展多语言
 */
export default {
  common: {
    appName: 'IntelliJ 配置助手',
    loading: '加载中...',
    success: '成功',
    error: '错误',
  },

  navigation: {
    main: '主页',
    about: '关于',
  },

  titleBar: {
    title: 'IntelliJ 配置助手',
  },

  mainView: {
    installPath: {
      title: '安装路径',
      description: '指向 IntelliJ 的安装目录，后续所有动作都会以这里为基准喵。',
      label: 'IntelliJ 安装目录',
      placeholder: '例如：D:/Program Files/JetBrains/IntelliJ IDEA 2024.2',
    },

    applyConfig: {
      title: '应用配置',
      description: '选择需要导入的配置目录，浮浮酱会把它融合进上方的 IntelliJ 路径喵。',
      label: '配置目录',
      placeholder: '例如：D:/jetbra',
      submitButton: '应用配置',
      submitting: '提交中...',
      successMessage: '配置成功应用到 {count} 个文件, 请重启需要激活编译器输入激活码',
    },

    clearConfig: {
      title: '清除配置',
      description: '恢复 IntelliJ 为初始状态，只会动到上方路径对应的配置喵。',
      currentPath: '当前路径：',
      noPath: '尚未填写',
      clearButton: '清除配置',
      clearing: '清除中...',
      successMessage: '成功清除 {count} 个文件的配置',
    },

    notice: {
      text: '若上述操作未成功，请下载压缩包后进入 <code>scripts</code> 文件夹按操作系统执行下列脚本喵：',
      windows:
        'Windows：运行 <code>uninstall-all-users.vbs</code> 与 <code>uninstall-current-user.vbs</code>',
      linuxMac: 'Linux / macOS：运行 <code>uninstall.sh</code>',
    },
  },

  aboutView: {
    main: {
      title: '关于 IntelliJ 配置助手',
      description: '这里记录应用的版本、核心技术栈与维护者信息，布局保持与主页一致喵。',
      appName: '应用名称',
      version: '当前版本',
      buildTool: '构建工具',
    },

    techStack: {
      title: '技术栈',
      description: '支撑 IntelliJ 配置助手的关键技术组件明细喵。',
      frontend: '前端框架',
      backend: '后端语言',
      buildTool: '构建工具',
    },

    project: {
      title: '项目信息',
      description: '快速找到仓库地址与维护者，方便主人参考喵。',
      repo: '仓库地址',
      developers: '开发者',
    },

    notice:
      '本项目运行所产生的一切问题需自行承担，项目由 Claude 4.5 和 GPT-5 配合开发，仅限学习使用喵。',
    footer: '© 2025 IntelliJ 配置助手 · 使用 Wails 和 Vue 构建',
  },

  validation: {
    emptyPaths: '请输入完整的两个路径喵～',
    invalidConfigPath: '配置目录路径包含不支持的特殊字符喵～',
    pathNotExist: '{label} 不存在，请检查喵～',
    emptyIntellijPath: '请输入 IntelliJ 路径',
  },

  errors: {
    unknown: '未知错误',
    fetchAboutInfoFailed: '获取关于信息失败',
    closeWindowFailed: '关闭窗口失败',
    minimizeWindowFailed: '最小化窗口失败',
    toggleWindowSizeFailed: '切换窗口大小失败',
    openExternalLinkFailed: '打开外部链接失败',
  },

  theme: {
    dark: '暗色',
    light: '亮色',
    switchToDark: '切换到暗色模式',
    switchToLight: '切换到亮色模式',
  },

  language: {
    current: '中文',
    switchTo: '切换到英文',
  },

  update: {
    available: '新版本可用！',
    viewDetails: '查看详情',
    remindLater: '稍后提醒',
    checkingFailed: '检查更新失败',
  },
}
