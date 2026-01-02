package service

import (
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
)

// ConfigService 提供 IntelliJ 配置管理的路径验证工具
// 优化：将大文件拆分为多个功能模块，遵循单一职责原则
type ConfigService struct {
	logger *slog.Logger
}

// Developer 保存开发者信息
type Developer struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// AboutInfo 保存应用程序关于信息
type AboutInfo struct {
	AppName      string      `json:"appName"`
	Version      string      `json:"version"`
	GoVersion    string      `json:"goVersion"`
	VueVersion   string      `json:"vueVersion"`
	WailsVersion string      `json:"wailsVersion"`
	RepoURL      string      `json:"repoUrl"`
	Developers   []Developer `json:"developers"`
}

// VMOptionsOperation 定义对 vmoptions 文件的操作
type VMOptionsOperation func(filePath string) error

// NewConfigService 构造一个准备好与 Wails 绑定的 ConfigService 实例
func NewConfigService() *ConfigService {
	return &ConfigService{
		logger: slog.Default(),
	}
}

// processVMOptionsFilesGeneric 通用的 vmoptions 文件处理流程
// 提取 SubmitPaths 和 ClearConfig 中的共同逻辑，避免代码重复
func (c *ConfigService) processVMOptionsFilesGeneric(projectPath string, operation VMOptionsOperation, operationName string) (int, error) {
	projectPath = sanitizePath(projectPath)

	if projectPath == "" {
		c.logger.Warn("路径验证失败: 路径为空")
		return 0, ErrEmptyPath
	}

	binDir, err := validateIntelliJPath(projectPath)
	if err != nil {
		c.logger.Error("IntelliJ路径验证失败", slog.Any("error", err))
		return 0, err
	}

	// 查找所有 .vmoptions 文件
	vmOptionsFiles, err := findVMOptionsFiles(binDir)
	if err != nil {
		c.logger.Error("查找vmoptions文件失败", slog.Any("error", err))
		return 0, err
	}

	if len(vmOptionsFiles) == 0 {
		return 0, ErrNoVMOptions
	}

	c.logger.Info("找到vmoptions文件", slog.Int("count", len(vmOptionsFiles)))

	// 处理每个 vmoptions 文件
	count := 0
	for _, vmFile := range vmOptionsFiles {
		c.logger.Debug(operationName+"文件", slog.String("file", vmFile))
		if err := operation(vmFile); err != nil {
			c.logger.Error(operationName+"文件失败",
				slog.String("file", vmFile),
				slog.Any("error", err))
			return 0, fmt.Errorf(operationName+"文件 %s 失败: %w", filepath.Base(vmFile), err)
		}
		count++
	}

	return count, nil
}

// SubmitPaths 验证提供的路径，修改 vmoptions 文件并应用配置
func (c *ConfigService) SubmitPaths(projectPath, configPath string) (string, error) {
	c.logger.Info("开始验证路径",
		slog.String("intellijPath", projectPath),
		slog.String("configPath", configPath))

	configPath = sanitizePath(configPath)

	if configPath == "" {
		c.logger.Warn("路径验证失败: 配置路径为空")
		return "", ErrEmptyPath
	}

	if err := validateConfigPath(configPath); err != nil {
		c.logger.Error("配置路径验证失败", slog.Any("error", err))
		return "", err
	}

	// 先清除已有的环境变量，避免旧配置干扰
	c.logger.Info("开始清除旧的环境变量")
	envWarningMsg, err := removeJetBrainsEnvVars(c.logger)
	if err != nil {
		c.logger.Warn("清除环境变量时出现警告", slog.Any("error", err))
		// 环境变量清除失败不影响整体操作，继续执行
	}

	// 规范化配置路径为Unix风格
	normalizedConfigPath := filepath.ToSlash(configPath)

	// 处理 vmoptions 文件
	operation := func(vmFile string) error {
		return processVMOptionsFile(vmFile, normalizedConfigPath, c.logger)
	}

	processedCount, err := c.processVMOptionsFilesGeneric(projectPath, operation, "处理")
	if err != nil {
		return "", err
	}

	c.logger.Info("配置应用成功", slog.Int("processedCount", processedCount))

	// 构建返回消息
	resultMsg := fmt.Sprintf("配置成功应用到 %d 个文件, 请重启需要激活编译器输入激活码", processedCount)
	if envWarningMsg != "" {
		resultMsg += "\n⚠️ " + envWarningMsg
	}

	return resultMsg, nil
}

// ClearConfig 从 vmoptions 文件中移除添加的配置
func (c *ConfigService) ClearConfig(projectPath string) (string, error) {
	c.logger.Info("开始清除配置", slog.String("intellijPath", projectPath))

	// 处理 vmoptions 文件
	operation := func(vmFile string) error {
		return clearVMOptionsFile(vmFile, c.logger)
	}

	clearedCount, err := c.processVMOptionsFilesGeneric(projectPath, operation, "清除")
	if err != nil {
		return "", err
	}

	// 清除环境变量
	warningMsg, err := removeJetBrainsEnvVars(c.logger)
	if err != nil {
		c.logger.Warn("清除环境变量时出现警告", slog.Any("error", err))
		// 环境变量清除失败不影响整体操作，继续执行
	}

	c.logger.Info("配置清除成功", slog.Int("clearedCount", clearedCount))

	// 构建返回消息
	resultMsg := fmt.Sprintf("成功清除 %d 个文件的配置", clearedCount)
	if warningMsg != "" {
		resultMsg += "\n⚠️ " + warningMsg
	}

	return resultMsg, nil
}

// PathExists 对提供的路径执行轻量级存在性检查
func (c *ConfigService) PathExists(path string) (bool, error) {
	cleaned := sanitizePath(path)
	if cleaned == "" {
		return false, ErrEmptyPath
	}

	if _, err := os.Stat(cleaned); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return false, nil
		}
		return false, fmt.Errorf("无法访问路径: %w", err)
	}

	return true, nil
}

// GetAboutInfo 返回应用程序关于信息
func (c *ConfigService) GetAboutInfo() AboutInfo {
	return AboutInfo{
		AppName:      AppName,
		Version:      Version,
		GoVersion:    runtime.Version(),
		VueVersion:   "3.5.22",
		WailsVersion: "v3",
		RepoURL:      RepoURL,
		Developers: []Developer{
			{Name: "XgzK", URL: "https://github.com/XgzK"},
			{Name: "Claude (AI)", URL: "https://github.com/anthropics/claude-code"},
			{Name: "Codex (AI)", URL: "https://github.com/openai/codex"},
		},
	}
}
