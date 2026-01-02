//go:build !windows
// +build !windows

package service

import (
	"log/slog"
)

// removeJetBrainsEnvVars 删除所有 JetBrains 产品的环境变量
// 在 Windows 上删除用户级和系统级环境变量
// 在其他平台上仅记录信息（因为通常不使用注册表方式）
func removeJetBrainsEnvVars(logger *slog.Logger) (string, error) {
	logger.Debug("非Windows平台，环境变量清除为空操作")
	return "", nil
}
