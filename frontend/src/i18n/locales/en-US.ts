/**
 * English (US) Language Pack
 * Optimization: i18n support for future multi-language expansion
 */
export default {
  common: {
    appName: 'IntelliJ Config Helper',
    loading: 'Loading...',
    success: 'Success',
    error: 'Error',
  },

  navigation: {
    main: 'Main',
    about: 'About',
  },

  titleBar: {
    title: 'IntelliJ Config Helper',
  },

  mainView: {
    installPath: {
      title: 'Installation Path',
      description:
        'Point to your IntelliJ installation directory. All subsequent actions will be based on this path.',
      label: 'IntelliJ Installation Directory',
      placeholder: 'e.g., D:/Program Files/JetBrains/IntelliJ IDEA 2024.2',
    },

    applyConfig: {
      title: 'Apply Configuration',
      description:
        'Select the configuration directory to import. It will be merged into the IntelliJ path above.',
      label: 'Configuration Directory',
      placeholder: 'e.g., D:/jetbra',
      submitButton: 'Apply Config',
      submitting: 'Submitting...',
      successMessage:
        'Configuration successfully applied to {count} file(s). Please restart and enter activation code.',
    },

    clearConfig: {
      title: 'Clear Configuration',
      description:
        'Restore IntelliJ to its initial state. Only affects the configuration at the path above.',
      currentPath: 'Current path:',
      noPath: 'Not filled',
      clearButton: 'Clear Config',
      clearing: 'Clearing...',
      successMessage: 'Successfully cleared configuration from {count} file(s)',
    },

    notice: {
      text: 'If the above operations fail, download the archive and execute the following scripts in the <code>scripts</code> folder:',
      windows:
        'Windows: Run <code>uninstall-all-users.vbs</code> and <code>uninstall-current-user.vbs</code>',
      linuxMac: 'Linux / macOS: Run <code>uninstall.sh</code>',
    },
  },

  aboutView: {
    main: {
      title: 'About IntelliJ Config Helper',
      description:
        'This page records the application version, core tech stack, and maintainer information.',
      appName: 'Application Name',
      version: 'Current Version',
      buildTool: 'Build Tool',
    },

    techStack: {
      title: 'Tech Stack',
      description: 'Key technology components supporting IntelliJ Config Helper.',
      frontend: 'Frontend Framework',
      backend: 'Backend Language',
      buildTool: 'Build Tool',
    },

    project: {
      title: 'Project Information',
      description: 'Quick access to repository and maintainers.',
      repo: 'Repository',
      developers: 'Developers',
    },

    notice:
      'All issues arising from running this project are your own responsibility. Developed collaboratively by Claude 4.5 and GPT-5, for educational purposes only.',
    footer: '© 2025 IntelliJ Config Helper · Built with Wails and Vue',
  },

  validation: {
    emptyPaths: 'Please enter both paths',
    invalidConfigPath: 'Config directory path contains unsupported special characters',
    pathNotExist: '{label} does not exist, please check',
    emptyIntellijPath: 'Please enter IntelliJ path',
  },

  errors: {
    unknown: 'Unknown error',
    fetchAboutInfoFailed: 'Failed to fetch about information',
    closeWindowFailed: 'Failed to close window',
    minimizeWindowFailed: 'Failed to minimize window',
    toggleWindowSizeFailed: 'Failed to toggle window size',
    openExternalLinkFailed: 'Failed to open external link',
  },

  theme: {
    dark: 'Dark',
    light: 'Light',
    switchToDark: 'Switch to dark mode',
    switchToLight: 'Switch to light mode',
  },

  language: {
    current: 'EN',
    switchTo: 'Switch to Chinese',
  },

  update: {
    available: 'New version available!',
    viewDetails: 'View Details',
    remindLater: 'Remind Later',
    checkingFailed: 'Failed to check for updates',
  },
}
