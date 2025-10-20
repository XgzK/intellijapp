package service

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// ConfigService provides path validation utilities for IntelliJ configuration management.
type ConfigService struct {
	logger *log.Logger
}

// ConfigPaths holds validated configuration paths
type ConfigPaths struct {
	IntelliJBinDir string
	ConfigDir      string
}

// NewConfigService constructs a ConfigService instance ready for binding with Wails.
func NewConfigService() *ConfigService {
	return &ConfigService{
		logger: log.Default(),
	}
}

// SubmitPaths validates the provided paths, modifies vmoptions files, and applies configuration.
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

// ClearConfig removes the added configuration from vmoptions files.
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

// PathExists performs a lightweight existence check for the provided path.
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

func sanitizePath(path string) string {
	trimmed := strings.TrimSpace(path)
	if trimmed == "" {
		return ""
	}
	return filepath.Clean(trimmed)
}

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

// clearVMOptionsFile 清除单个 vmoptions 文件中所有以 --add-opens 和 -javaagent: 开头的配置
func clearVMOptionsFile(filePath string, logger *log.Logger) error {
	// 读取文件内容
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("读取文件失败: %w", err)
	}

	lines := strings.Split(string(content), "\n")
	var newLines []string
	removedCount := 0

	// 处理每一行
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// 跳过所有以 --add-opens 开头的行
		if strings.HasPrefix(trimmed, "--add-opens") {
			logger.Printf("删除行: %s", trimmed)
			removedCount++
			continue
		}

		// 跳过所有以 -javaagent: 开头的行
		if strings.HasPrefix(trimmed, "-javaagent:") {
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
