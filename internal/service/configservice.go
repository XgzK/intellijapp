package service

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// ConfigService 提供 IntelliJ 配置管理的路径验证工具
type ConfigService struct {
	logger *log.Logger
}

// ConfigPaths 保存已验证的配置路径
type ConfigPaths struct {
	IntelliJBinDir string
	ConfigDir      string
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

// NewConfigService 构造一个准备好与 Wails 绑定的 ConfigService 实例
func NewConfigService() *ConfigService {
	return &ConfigService{
		logger: log.Default(),
	}
}

// SubmitPaths 验证提供的路径，修改 vmoptions 文件并应用配置
func (c *ConfigService) SubmitPaths(projectPath string, configPath string) (string, error) {
	c.logger.Printf("开始验证路径 - IntelliJ路径: %s, 配置路径: %s", projectPath, configPath)

	projectPath = sanitizePath(projectPath)
	configPath = sanitizePath(configPath)

	if projectPath == "" || configPath == "" {
		c.logger.Printf("路径验证失败: 路径为空")
		return "", errors.New("路径不能为空")
	}

	if err := validateConfigPath(configPath); err != nil {
		c.logger.Printf("配置路径验证失败: %v", err)
		return "", err
	}

	binDir, err := validateIntelliJPath(projectPath)
	if err != nil {
		c.logger.Printf("IntelliJ路径验证失败: %v", err)
		return "", err
	}

	// 查找并处理所有 .vmoptions 文件
	vmOptionsFiles, err := findVMOptionsFiles(binDir)
	if err != nil {
		c.logger.Printf("查找vmoptions文件失败: %v", err)
		return "", err
	}

	if len(vmOptionsFiles) == 0 {
		return "", errors.New("未找到任何 .vmoptions 文件")
	}

	c.logger.Printf("找到 %d 个 vmoptions 文件", len(vmOptionsFiles))

	// 规范化配置路径为Unix风格
	normalizedConfigPath := filepath.ToSlash(configPath)

	// 处理每个 vmoptions 文件
	processedCount := 0
	for _, vmFile := range vmOptionsFiles {
		c.logger.Printf("处理文件: %s", vmFile)
		if err := processVMOptionsFile(vmFile, normalizedConfigPath, c.logger); err != nil {
			c.logger.Printf("处理文件 %s 失败: %v", vmFile, err)
			return "", fmt.Errorf("处理文件 %s 失败: %w", filepath.Base(vmFile), err)
		}
		processedCount++
	}

	c.logger.Printf("配置成功应用到 %d 个文件, 请重启需要激活编译器输入激活码", processedCount)
	return fmt.Sprintf("配置成功应用到 %d 个文件, 请重启需要激活编译器输入激活码", processedCount), nil
}

// ClearConfig 从 vmoptions 文件中移除添加的配置
func (c *ConfigService) ClearConfig(projectPath string) (string, error) {
	c.logger.Printf("开始清除配置 - IntelliJ路径: %s", projectPath)

	projectPath = sanitizePath(projectPath)

	if projectPath == "" {
		c.logger.Printf("路径验证失败: 路径为空")
		return "", errors.New("路径不能为空")
	}

	binDir, err := validateIntelliJPath(projectPath)
	if err != nil {
		c.logger.Printf("IntelliJ路径验证失败: %v", err)
		return "", err
	}

	// 查找所有 .vmoptions 文件
	vmOptionsFiles, err := findVMOptionsFiles(binDir)
	if err != nil {
		c.logger.Printf("查找vmoptions文件失败: %v", err)
		return "", err
	}

	if len(vmOptionsFiles) == 0 {
		return "", errors.New("未找到任何 .vmoptions 文件")
	}

	c.logger.Printf("找到 %d 个 vmoptions 文件", len(vmOptionsFiles))

	// 清除每个 vmoptions 文件的配置
	clearedCount := 0
	for _, vmFile := range vmOptionsFiles {
		c.logger.Printf("清除文件: %s", vmFile)
		if err := clearVMOptionsFile(vmFile, c.logger); err != nil {
			c.logger.Printf("清除文件 %s 失败: %v", vmFile, err)
			return "", fmt.Errorf("清除文件 %s 失败: %w", filepath.Base(vmFile), err)
		}
		clearedCount++
	}

	c.logger.Printf("成功清除 %d 个 vmoptions 文件的配置", clearedCount)
	return fmt.Sprintf("成功清除 %d 个文件的配置", clearedCount), nil
}

// PathExists 对提供的路径执行轻量级存在性检查
func (c *ConfigService) PathExists(path string) (bool, error) {
	cleaned := sanitizePath(path)
	if cleaned == "" {
		return false, errors.New("路径不能为空")
	}

	if _, err := os.Stat(cleaned); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("无法访问路径: %w", err)
	}

	return true, nil
}

// GetAboutInfo 返回应用程序关于信息
func (c *ConfigService) GetAboutInfo() AboutInfo {
	return AboutInfo{
		AppName:      "IntelliJ 配置助手",
		Version:      "v1.0.3",
		GoVersion:    runtime.Version(),
		VueVersion:   "3.5.22",
		WailsVersion: "v3",
		RepoURL:      "github.com/XgzK/intellijapp",
		Developers: []Developer{
			{Name: "XgzK", URL: "https://github.com/XgzK"},
			{Name: "Claude (AI)", URL: "https://github.com/anthropics/claude-code"},
			{Name: "Codex (AI)", URL: "https://github.com/openai/codex"},
		},
	}
}

// sanitizePath 清理路径字符串，去除首尾空格并规范化路径
func sanitizePath(path string) string {
	trimmed := strings.TrimSpace(path)
	if trimmed == "" {
		return ""
	}
	return filepath.Clean(trimmed)
}

// validateConfigPath 验证配置目录是否存在且包含必需的 ja-netfilter.jar 文件
func validateConfigPath(configDir string) error {
	info, err := os.Stat(configDir)
	if err != nil {
		if os.IsNotExist(err) {
			return errors.New("配置目录不存在")
		}
		return fmt.Errorf("无法访问配置目录: %w", err)
	}

	if !info.IsDir() {
		return errors.New("配置路径必须是目录")
	}

	// 验证配置文件是否存在
	jarPath := filepath.Join(configDir, "ja-netfilter.jar")
	if _, err := os.Stat(jarPath); err != nil {
		if os.IsNotExist(err) {
			return errors.New("配置目录缺少 ja-netfilter.jar 文件")
		}
		return fmt.Errorf("无法检查配置文件: %w", err)
	}

	return nil
}

// validateIntelliJPath 验证 IntelliJ 软件路径，自动识别并返回 bin 目录路径
func validateIntelliJPath(softwarePath string) (string, error) {
	info, err := os.Stat(softwarePath)
	if err != nil {
		return "", fmt.Errorf("无法访问 IntelliJ 软件路径: %w", err)
	}

	if !info.IsDir() {
		return "", errors.New("IntelliJ 软件路径必须是目录")
	}

	base := strings.ToLower(filepath.Base(softwarePath))
	candidateBin := softwarePath

	if base != "bin" {
		candidateBin = filepath.Join(softwarePath, "bin")
		info, err := os.Stat(candidateBin)
		if err != nil {
			if os.IsNotExist(err) {
				return "", errors.New("非IntelliJ系列软件安装路径")
			}
			return "", fmt.Errorf("无法访问 IntelliJ bin 目录: %w", err)
		}
		if !info.IsDir() {
			return "", errors.New("非IntelliJ系列软件安装路径")
		}
	}

	hasVMOptions, err := directoryHasVMOptions(candidateBin)
	if err != nil {
		return "", err
	}

	if !hasVMOptions {
		return "", errors.New("非IntelliJ系列软件安装路径")
	}

	return candidateBin, nil
}

// directoryHasVMOptions 检查目录是否包含 .vmoptions 文件
func directoryHasVMOptions(dir string) (bool, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false, fmt.Errorf("无法读取目录 %s: %w", dir, err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if strings.HasSuffix(strings.ToLower(entry.Name()), ".vmoptions") {
			return true, nil
		}
	}

	return false, nil
}

// findVMOptionsFiles 查找目录中所有的 .vmoptions 文件
func findVMOptionsFiles(dir string) ([]string, error) {
	// 检查目录读取权限
	if err := checkDirReadPermission(dir); err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("无法读取目录 %s: %w", dir, err)
	}

	var vmOptionsFiles []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if strings.HasSuffix(strings.ToLower(entry.Name()), ".vmoptions") {
			vmOptionsFiles = append(vmOptionsFiles, filepath.Join(dir, entry.Name()))
		}
	}

	return vmOptionsFiles, nil
}

// processVMOptionsFile 处理单个 vmoptions 文件
func processVMOptionsFile(filePath string, configPath string, logger *log.Logger) error {
	// 检查文件读取权限
	if err := checkFileReadPermission(filePath); err != nil {
		return err
	}

	// 检查文件写入权限
	if err := checkFileWritePermission(filePath); err != nil {
		return err
	}

	// 读取文件内容
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("读取文件失败: %w", err)
	}

	lines := strings.Split(string(content), "\n")
	var newLines []string

	// 处理每一行
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// 跳过以 --add-opens 开头的行
		if strings.HasPrefix(trimmed, "--add-opens") {
			logger.Printf("删除行: %s", trimmed)
			continue
		}

		// 跳过以 -javaagent: 开头的行
		if strings.HasPrefix(trimmed, "-javaagent:") {
			logger.Printf("删除行: %s", trimmed)
			continue
		}

		newLines = append(newLines, line)
	}

	// 添加新的配置
	newLines = append(newLines, "--add-opens=java.base/jdk.internal.org.objectweb.asm=ALL-UNNAMED")
	newLines = append(newLines, "--add-opens=java.base/jdk.internal.org.objectweb.asm.tree=ALL-UNNAMED")
	newLines = append(newLines, fmt.Sprintf("-javaagent:%s/ja-netfilter.jar=jetbrains", configPath))

	logger.Printf("添加配置: --add-opens (2行)")
	logger.Printf("添加配置: -javaagent:%s/ja-netfilter.jar=jetbrains", configPath)

	// 写回文件
	newContent := strings.Join(newLines, "\n")
	if err := os.WriteFile(filePath, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("写入文件失败: %w", err)
	}

	logger.Printf("成功更新文件: %s", filepath.Base(filePath))
	return nil
}

// clearVMOptionsFile 清除单个 vmoptions 文件中本工具添加的特定配置（不影响用户自定义配置）
func clearVMOptionsFile(filePath string, logger *log.Logger) error {
	// 检查文件读取权限
	if err := checkFileReadPermission(filePath); err != nil {
		return err
	}

	// 检查文件写入权限
	if err := checkFileWritePermission(filePath); err != nil {
		return err
	}

	// 读取文件内容
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("读取文件失败: %w", err)
	}

	lines := strings.Split(string(content), "\n")
	var newLines []string
	removedCount := 0

	// 本工具添加的特定配置行
	toolAddedLines := map[string]bool{
		"--add-opens=java.base/jdk.internal.org.objectweb.asm=ALL-UNNAMED":      true,
		"--add-opens=java.base/jdk.internal.org.objectweb.asm.tree=ALL-UNNAMED": true,
	}

	// 处理每一行
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// 只删除本工具添加的特定 --add-opens 配置
		if toolAddedLines[trimmed] {
			logger.Printf("删除行: %s", trimmed)
			removedCount++
			continue
		}

		// 只删除包含 ja-netfilter.jar=jetbrains 的 javaagent 配置
		if strings.HasPrefix(trimmed, "-javaagent:") && strings.Contains(trimmed, "ja-netfilter.jar=jetbrains") {
			logger.Printf("删除行: %s", trimmed)
			removedCount++
			continue
		}

		newLines = append(newLines, line)
	}

	// 移除尾部的空行
	for len(newLines) > 0 && strings.TrimSpace(newLines[len(newLines)-1]) == "" {
		newLines = newLines[:len(newLines)-1]
	}

	// 写回文件
	newContent := strings.Join(newLines, "\n")
	if err := os.WriteFile(filePath, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("写入文件失败: %w", err)
	}

	logger.Printf("成功清除文件: %s (删除了 %d 行)", filepath.Base(filePath), removedCount)
	return nil
}

// checkFileReadPermission 检查文件是否有读取权限
func checkFileReadPermission(filePath string) error {
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0)
	if err != nil {
		if os.IsPermission(err) {
			return formatPermissionError(filePath, "读取")
		}
		return err
	}
	file.Close()
	return nil
}

// checkFileWritePermission 检查文件是否有写入权限
func checkFileWritePermission(filePath string) error {
	file, err := os.OpenFile(filePath, os.O_WRONLY, 0)
	if err != nil {
		if os.IsPermission(err) {
			return formatPermissionError(filePath, "写入")
		}
		return err
	}
	file.Close()
	return nil
}

// checkDirReadPermission 检查目录是否有读取权限
func checkDirReadPermission(dirPath string) error {
	_, err := os.ReadDir(dirPath)
	if err != nil {
		if os.IsPermission(err) {
			return formatPermissionError(dirPath, "读取")
		}
		return err
	}
	return nil
}

// formatPermissionError 格式化权限错误信息
func formatPermissionError(path string, operation string) error {
	if runtime.GOOS == "windows" {
		return fmt.Errorf("没有%s权限: %s\n请以管理员身份运行程序", operation, path)
	}
	return fmt.Errorf("没有%s权限: %s\n请使用 sudo 或以 root 身份运行程序", operation, path)
}
