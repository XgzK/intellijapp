package service

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

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
