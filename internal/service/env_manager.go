//go:build windows
// +build windows

package service

import (
	"fmt"
	"log/slog"
	"runtime"
	"strings"

	"golang.org/x/sys/windows/registry"
)

// JetBrains 产品列表（与 VBS 脚本中的列表一致）
var jetbrainsProducts = []string{
	"idea", "clion", "phpstorm", "goland", "pycharm", "webstorm",
	"webide", "rider", "datagrip", "rubymine", "dataspell", "aqua",
	"rustrover", "gateway", "jetbrains_client", "jetbrainsclient",
	"studio", "devecostudio",
}

// buildEnvVarName 构建 JetBrains 产品的环境变量名
func buildEnvVarName(product string) string {
	return strings.ToUpper(product) + "_VM_OPTIONS"
}

// removeJetBrainsEnvVars 删除所有 JetBrains 产品的环境变量
// 在 Windows 上删除用户级和系统级环境变量
// 在其他平台上仅记录警告（因为通常不使用环境变量方式）
// 返回的错误字符串会包含权限提示信息
func removeJetBrainsEnvVars(logger *slog.Logger) (string, error) {
	if runtime.GOOS != "windows" {
		logger.Info("非 Windows 平台，跳过环境变量清除")
		return "", nil
	}

	logger.Info("开始清除 JetBrains 环境变量")

	var warnings []string
	var errors []string
	removedCount := 0
	hasAdminRights := isRunningAsAdmin()

	// 检查用户级和系统级环境变量是否存在
	userVarsExist := checkEnvVarsExist(registry.CURRENT_USER, `Environment`, logger)
	systemVarsExist := checkEnvVarsExist(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Control\Session Manager\Environment`, logger)

	// 清除用户级环境变量
	if userVarsExist {
		userRemoved, err := removeEnvVarsFromRegistry(registry.CURRENT_USER, `Environment`, logger)
		if err != nil {
			errors = append(errors, fmt.Sprintf("用户级: %v", err))
		} else {
			removedCount += userRemoved
			logger.Info("成功清除用户级环境变量", slog.Int("count", userRemoved))
		}
	} else {
		logger.Debug("未检测到用户级环境变量，跳过清除")
	}

	// 清除系统级环境变量（需要管理员权限）
	if systemVarsExist {
		if !hasAdminRights {
			warnings = append(warnings, "检测到系统级环境变量，但当前无管理员权限。请使用管理员权限运行以完全清除配置。")
			logger.Warn("检测到系统级环境变量但无管理员权限")
		} else {
			systemRemoved, err := removeEnvVarsFromRegistry(
				registry.LOCAL_MACHINE,
				`SYSTEM\CurrentControlSet\Control\Session Manager\Environment`,
				logger,
			)
			if err != nil {
				warnings = append(warnings, fmt.Sprintf("清除系统级环境变量失败: %v", err))
				logger.Warn("清除系统级环境变量失败", slog.Any("error", err))
			} else {
				removedCount += systemRemoved
				logger.Info("成功清除系统级环境变量", slog.Int("count", systemRemoved))
			}
		}
	} else {
		logger.Debug("未检测到系统级环境变量，跳过清除")
	}

	// 如果既没有用户级也没有系统级环境变量
	if !userVarsExist && !systemVarsExist {
		logger.Info("未检测到任何 JetBrains 环境变量")
	}

	if len(errors) > 0 && removedCount == 0 {
		return "", fmt.Errorf("清除环境变量失败: %s", strings.Join(errors, "; "))
	}

	logger.Info("环境变量清除完成", slog.Int("removedCount", removedCount))

	warningMsg := ""
	if len(warnings) > 0 {
		warningMsg = strings.Join(warnings, " ")
	}

	return warningMsg, nil
}

// isRunningAsAdmin 检测当前进程是否具有管理员权限
func isRunningAsAdmin() bool {
	if runtime.GOOS != "windows" {
		return false
	}

	// 尝试打开需要管理员权限的注册表键
	key, err := registry.OpenKey(
		registry.LOCAL_MACHINE,
		`SYSTEM\CurrentControlSet\Control\Session Manager\Environment`,
		registry.SET_VALUE,
	)
	if err != nil {
		return false
	}
	defer key.Close()

	return true
}

// checkEnvVarsExist 检查是否存在指定注册表位置的 JetBrains 环境变量
// 这是 checkUserEnvVarsExist 和 checkSystemEnvVarsExist 的通用实现
func checkEnvVarsExist(rootKey registry.Key, path string, logger *slog.Logger) bool {
	key, err := registry.OpenKey(
		rootKey,
		path,
		registry.QUERY_VALUE,
	)
	if err != nil {
		// 无法打开注册表键，假设不存在
		return false
	}
	defer key.Close()

	for _, product := range jetbrainsProducts {
		envVarName := buildEnvVarName(product)
		_, _, err := key.GetStringValue(envVarName)
		if err == nil {
			// 找到至少一个环境变量
			logger.Debug("检测到环境变量", slog.String("name", envVarName))
			return true
		}
	}

	return false
}

// removeEnvVarsFromRegistry 从指定的注册表位置删除 JetBrains 环境变量
func removeEnvVarsFromRegistry(rootKey registry.Key, path string, logger *slog.Logger) (int, error) {
	// 打开注册表键
	key, err := registry.OpenKey(rootKey, path, registry.SET_VALUE|registry.QUERY_VALUE)
	if err != nil {
		return 0, fmt.Errorf("打开注册表键失败: %w", err)
	}
	defer key.Close()

	removedCount := 0
	var errors []string

	// 遍历所有 JetBrains 产品
	for _, product := range jetbrainsProducts {
		envVarName := buildEnvVarName(product)

		// 检查环境变量是否存在
		_, _, err := key.GetStringValue(envVarName)
		if err == registry.ErrNotExist {
			// 环境变量不存在，跳过
			continue
		}
		if err != nil {
			logger.Warn("读取环境变量失败",
				slog.String("name", envVarName),
				slog.Any("error", err))
			continue
		}

		// 删除环境变量
		err = key.DeleteValue(envVarName)
		if err != nil {
			errors = append(errors, fmt.Sprintf("%s: %v", envVarName, err))
			logger.Error("删除环境变量失败",
				slog.String("name", envVarName),
				slog.Any("error", err))
		} else {
			removedCount++
			logger.Debug("成功删除环境变量", slog.String("name", envVarName))
		}
	}

	if len(errors) > 0 {
		return removedCount, fmt.Errorf("部分删除失败 (成功删除: %d, 失败: %d): %s",
			removedCount, len(errors), strings.Join(errors, "; "))
	}

	return removedCount, nil
}

// broadcastEnvironmentChange 广播环境变量变更消息
// 通知 Windows 系统环境变量已更改，使其他程序能够感知到变化
func broadcastEnvironmentChange() {
	if runtime.GOOS != "windows" {
		return
	}

	// 在 Windows 上，删除环境变量后不需要立即广播
	// 系统会在下次登录时自动生效
	// 如果需要立即生效，可以使用 SendMessage API，但这超出了当前需求范围
}
