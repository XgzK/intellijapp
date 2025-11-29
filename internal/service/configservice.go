package service

import (
	"errors"
	"fmt"
	"io/fs"
	"iter"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
)

// Sentinel errors - 定义可复用的错误类型喵～
var (
	ErrEmptyPath        = errors.New("路径不能为空")
	ErrPathNotExist     = errors.New("路径不存在")
	ErrPathNotDir       = errors.New("路径必须是目录")
	ErrNotIntelliJDir   = errors.New("非IntelliJ系列软件安装路径")
	ErrNoVMOptions      = errors.New("未找到任何 .vmoptions 文件")
	ErrMissingJarFile   = errors.New("配置目录缺少 ja-netfilter.jar 文件")
	ErrPermissionDenied = errors.New("权限不足")
)

// ConfigService 提供 IntelliJ 配置管理的路径验证工具
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

// NewConfigService 构造一个准备好与 Wails 绑定的 ConfigService 实例
func NewConfigService() *ConfigService {
	return &ConfigService{
		logger: slog.Default(),
	}
}

// SubmitPaths 验证提供的路径，修改 vmoptions 文件并应用配置
func (c *ConfigService) SubmitPaths(projectPath, configPath string) (string, error) {
	c.logger.Info("开始验证路径",
		slog.String("intellijPath", projectPath),
		slog.String("configPath", configPath))

	projectPath = sanitizePath(projectPath)
	configPath = sanitizePath(configPath)

	if projectPath == "" || configPath == "" {
		c.logger.Warn("路径验证失败: 路径为空")
		return "", ErrEmptyPath
	}

	if err := validateConfigPath(configPath); err != nil {
		c.logger.Error("配置路径验证失败", slog.Any("error", err))
		return "", err
	}

	binDir, err := validateIntelliJPath(projectPath)
	if err != nil {
		c.logger.Error("IntelliJ路径验证失败", slog.Any("error", err))
		return "", err
	}

	// 查找并处理所有 .vmoptions 文件
	vmOptionsFiles, err := findVMOptionsFiles(binDir)
	if err != nil {
		c.logger.Error("查找vmoptions文件失败", slog.Any("error", err))
		return "", err
	}

	if len(vmOptionsFiles) == 0 {
		return "", ErrNoVMOptions
	}

	c.logger.Info("找到vmoptions文件", slog.Int("count", len(vmOptionsFiles)))

	// 规范化配置路径为Unix风格
	normalizedConfigPath := filepath.ToSlash(configPath)

	// 处理每个 vmoptions 文件
	processedCount := 0
	for _, vmFile := range vmOptionsFiles {
		c.logger.Debug("处理文件", slog.String("file", vmFile))
		if err := processVMOptionsFile(vmFile, normalizedConfigPath, c.logger); err != nil {
			c.logger.Error("处理文件失败",
				slog.String("file", vmFile),
				slog.Any("error", err))
			return "", fmt.Errorf("处理文件 %s 失败: %w", filepath.Base(vmFile), err)
		}
		processedCount++
	}

	c.logger.Info("配置应用成功", slog.Int("processedCount", processedCount))
	return fmt.Sprintf("配置成功应用到 %d 个文件, 请重启需要激活编译器输入激活码", processedCount), nil
}

// ClearConfig 从 vmoptions 文件中移除添加的配置
func (c *ConfigService) ClearConfig(projectPath string) (string, error) {
	c.logger.Info("开始清除配置", slog.String("intellijPath", projectPath))

	projectPath = sanitizePath(projectPath)

	if projectPath == "" {
		c.logger.Warn("路径验证失败: 路径为空")
		return "", ErrEmptyPath
	}

	binDir, err := validateIntelliJPath(projectPath)
	if err != nil {
		c.logger.Error("IntelliJ路径验证失败", slog.Any("error", err))
		return "", err
	}

	// 查找所有 .vmoptions 文件
	vmOptionsFiles, err := findVMOptionsFiles(binDir)
	if err != nil {
		c.logger.Error("查找vmoptions文件失败", slog.Any("error", err))
		return "", err
	}

	if len(vmOptionsFiles) == 0 {
		return "", ErrNoVMOptions
	}

	c.logger.Info("找到vmoptions文件", slog.Int("count", len(vmOptionsFiles)))

	// 清除每个 vmoptions 文件的配置
	clearedCount := 0
	for _, vmFile := range vmOptionsFiles {
		c.logger.Debug("清除文件", slog.String("file", vmFile))
		if err := clearVMOptionsFile(vmFile, c.logger); err != nil {
			c.logger.Error("清除文件失败",
				slog.String("file", vmFile),
				slog.Any("error", err))
			return "", fmt.Errorf("清除文件 %s 失败: %w", filepath.Base(vmFile), err)
		}
		clearedCount++
	}

	c.logger.Info("配置清除成功", slog.Int("clearedCount", clearedCount))
	return fmt.Sprintf("成功清除 %d 个文件的配置", clearedCount), nil
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
		AppName:      "IntelliJ 配置助手",
		Version:      "v1.0.4",
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
		if errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("%w: %s", ErrPathNotExist, configDir)
		}
		return fmt.Errorf("无法访问配置目录: %w", err)
	}

	if !info.IsDir() {
		return fmt.Errorf("%w: %s", ErrPathNotDir, configDir)
	}

	// 验证配置文件是否存在
	jarPath := filepath.Join(configDir, "ja-netfilter.jar")
	if _, err := os.Stat(jarPath); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return ErrMissingJarFile
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
		return "", fmt.Errorf("%w: %s", ErrPathNotDir, softwarePath)
	}

	base := strings.ToLower(filepath.Base(softwarePath))
	candidateBin := softwarePath

	if base != "bin" {
		candidateBin = filepath.Join(softwarePath, "bin")
		info, err := os.Stat(candidateBin)
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				return "", ErrNotIntelliJDir
			}
			return "", fmt.Errorf("无法访问 IntelliJ bin 目录: %w", err)
		}
		if !info.IsDir() {
			return "", ErrNotIntelliJDir
		}
	}

	hasVMOptions, err := directoryHasVMOptions(candidateBin)
	if err != nil {
		return "", err
	}

	if !hasVMOptions {
		return "", ErrNotIntelliJDir
	}

	return candidateBin, nil
}

// directoryHasVMOptions 检查目录是否包含 .vmoptions 文件
func directoryHasVMOptions(dir string) (bool, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false, fmt.Errorf("无法读取目录 %s: %w", dir, err)
	}

	// 使用 slices.ContainsFunc 进行函数式查找
	return slices.ContainsFunc(entries, func(entry os.DirEntry) bool {
		return !entry.IsDir() && strings.HasSuffix(strings.ToLower(entry.Name()), ".vmoptions")
	}), nil
}

// findVMOptionsFiles 查找目录中所有的 .vmoptions 文件
// 使用 Go 1.23 iter.Seq 迭代器和 slices.Collect 收集结果
func findVMOptionsFiles(dir string) ([]string, error) {
	// 检查目录读取权限
	if err := checkDirReadPermission(dir); err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("无法读取目录 %s: %w", dir, err)
	}

	// 使用 Go 1.23 slices.Collect 配合自定义迭代器收集匹配的文件路径
	vmOptionsFiles := slices.Collect(mapEntriesToPaths(filterVMOptionsEntries(entries), dir))

	return vmOptionsFiles, nil
}

// processVMOptionsFile 处理单个 vmoptions 文件
func processVMOptionsFile(filePath, configPath string, logger *slog.Logger) error {
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

	// 使用 slices.DeleteFunc 过滤需要删除的行
	newLines := slices.DeleteFunc(lines, func(line string) bool {
		trimmed := strings.TrimSpace(line)
		shouldDelete := strings.HasPrefix(trimmed, "--add-opens") ||
			strings.HasPrefix(trimmed, "-javaagent:")
		if shouldDelete {
			logger.Debug("删除行", slog.String("line", trimmed))
		}
		return shouldDelete
	})

	// 添加新的配置
	newLines = append(newLines,
		"--add-opens=java.base/jdk.internal.org.objectweb.asm=ALL-UNNAMED",
		"--add-opens=java.base/jdk.internal.org.objectweb.asm.tree=ALL-UNNAMED",
		fmt.Sprintf("-javaagent:%s/ja-netfilter.jar=jetbrains", configPath),
	)

	logger.Debug("添加配置",
		slog.Int("addOpensCount", 2),
		slog.String("javaagent", configPath+"/ja-netfilter.jar"))

	// 写回文件
	newContent := strings.Join(newLines, "\n")
	if err := os.WriteFile(filePath, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("写入文件失败: %w", err)
	}

	logger.Debug("成功更新文件", slog.String("file", filepath.Base(filePath)))
	return nil
}

// toolAddedLines 定义本工具添加的特定配置行（使用包级变量避免重复创建）
var toolAddedLines = map[string]struct{}{
	"--add-opens=java.base/jdk.internal.org.objectweb.asm=ALL-UNNAMED":      {},
	"--add-opens=java.base/jdk.internal.org.objectweb.asm.tree=ALL-UNNAMED": {},
}

// clearVMOptionsFile 清除单个 vmoptions 文件中本工具添加的特定配置（不影响用户自定义配置）
func clearVMOptionsFile(filePath string, logger *slog.Logger) error {
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
	removedCount := 0

	// 使用 slices.DeleteFunc 过滤需要删除的行
	newLines := slices.DeleteFunc(lines, func(line string) bool {
		trimmed := strings.TrimSpace(line)

		// 只删除本工具添加的特定 --add-opens 配置
		if _, exists := toolAddedLines[trimmed]; exists {
			logger.Debug("删除行", slog.String("line", trimmed))
			removedCount++
			return true
		}

		// 只删除包含 ja-netfilter.jar=jetbrains 的 javaagent 配置
		if strings.HasPrefix(trimmed, "-javaagent:") && strings.Contains(trimmed, "ja-netfilter.jar=jetbrains") {
			logger.Debug("删除行", slog.String("line", trimmed))
			removedCount++
			return true
		}

		return false
	})

	// 移除尾部的空行 - 使用 Go 1.23 slices.Backward 反向迭代
	newLines = trimTrailingEmptyLines(newLines)

	// 写回文件
	newContent := strings.Join(newLines, "\n")
	if err := os.WriteFile(filePath, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("写入文件失败: %w", err)
	}

	logger.Debug("成功清除文件",
		slog.String("file", filepath.Base(filePath)),
		slog.Int("removedCount", removedCount))
	return nil
}

// checkFileReadPermission 检查文件是否有读取权限
func checkFileReadPermission(filePath string) error {
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0)
	if err != nil {
		if errors.Is(err, fs.ErrPermission) {
			return formatPermissionError(filePath, "读取")
		}
		return err
	}
	defer file.Close()
	return nil
}

// checkFileWritePermission 检查文件是否有写入权限
func checkFileWritePermission(filePath string) error {
	file, err := os.OpenFile(filePath, os.O_WRONLY, 0)
	if err != nil {
		if errors.Is(err, fs.ErrPermission) {
			return formatPermissionError(filePath, "写入")
		}
		return err
	}
	defer file.Close()
	return nil
}

// checkDirReadPermission 检查目录是否有读取权限
func checkDirReadPermission(dirPath string) error {
	_, err := os.ReadDir(dirPath)
	if err != nil {
		if errors.Is(err, fs.ErrPermission) {
			return formatPermissionError(dirPath, "读取")
		}
		return err
	}
	return nil
}

// formatPermissionError 格式化权限错误信息
func formatPermissionError(path, operation string) error {
	if runtime.GOOS == "windows" {
		return fmt.Errorf("%w: 没有%s权限 %s\n请以管理员身份运行程序", ErrPermissionDenied, operation, path)
	}
	return fmt.Errorf("%w: 没有%s权限 %s\n请使用 sudo 或以 root 身份运行程序", ErrPermissionDenied, operation, path)
}

// trimTrailingEmptyLines 使用 Go 1.23 slices.Backward 移除尾部空行
func trimTrailingEmptyLines(lines []string) []string {
	trimCount := 0
	for _, line := range slices.Backward(lines) {
		if strings.TrimSpace(line) != "" {
			break
		}
		trimCount++
	}
	if trimCount > 0 {
		return lines[:len(lines)-trimCount]
	}
	return lines
}

// filterVMOptionsEntries 使用 Go 1.23 iter.Seq 创建过滤迭代器
func filterVMOptionsEntries(entries []os.DirEntry) iter.Seq[os.DirEntry] {
	return func(yield func(os.DirEntry) bool) {
		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(strings.ToLower(entry.Name()), ".vmoptions") {
				if !yield(entry) {
					return
				}
			}
		}
	}
}

// mapEntriesToPaths 使用 Go 1.23 iter.Seq 将目录条目映射为完整文件路径
func mapEntriesToPaths(entries iter.Seq[os.DirEntry], dir string) iter.Seq[string] {
	return func(yield func(string) bool) {
		for entry := range entries {
			if !yield(filepath.Join(dir, entry.Name())) {
				return
			}
		}
	}
}
