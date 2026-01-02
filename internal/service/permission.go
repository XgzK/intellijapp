package service

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"runtime"
)

// checkFilePermission 检查文件是否有指定的权限（读取或写入）
// 优化：合并原有的 checkFileReadPermission 和 checkFileWritePermission，遵循 DRY 原则
func checkFilePermission(filePath string, mode int, operation string) error {
	file, err := os.OpenFile(filePath, mode, 0)
	if err != nil {
		if errors.Is(err, fs.ErrPermission) {
			return formatPermissionError(filePath, operation)
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

// checkFileReadPermission 检查文件是否有读取权限
// 便捷函数，避免重复传递 operation 参数
func checkFileReadPermission(filePath string) error {
	return checkFilePermission(filePath, os.O_RDONLY, "读取")
}

// checkFileWritePermission 检查文件是否有写入权限
// 便捷函数，避免重复传递 operation 参数
func checkFileWritePermission(filePath string) error {
	return checkFilePermission(filePath, os.O_WRONLY, "写入")
}
