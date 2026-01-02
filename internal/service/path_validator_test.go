package service

import (
	"os"
	"path/filepath"
	"testing"
)

// TestSanitizePath 测试路径清理函数
func TestSanitizePath(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "空路径",
			input:    "",
			expected: "",
		},
		{
			name:     "只有空格的路径",
			input:    "   ",
			expected: "",
		},
		{
			name:     "正常路径",
			input:    "/usr/local/bin",
			expected: filepath.Clean("/usr/local/bin"),
		},
		{
			name:     "带空格的路径",
			input:    "  /usr/local/bin  ",
			expected: filepath.Clean("/usr/local/bin"),
		},
		{
			name:     "Windows路径",
			input:    "C:\\Program Files\\JetBrains",
			expected: filepath.Clean("C:\\Program Files\\JetBrains"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitizePath(tt.input)
			if result != tt.expected {
				t.Errorf("sanitizePath(%q) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestValidateConfigPath 测试配置路径验证
func TestValidateConfigPath(t *testing.T) {
	// 创建临时测试目录
	tempDir := t.TempDir()

	// 创建有效的配置目录
	validConfigDir := filepath.Join(tempDir, "valid-config")
	if err := os.Mkdir(validConfigDir, 0755); err != nil {
		t.Fatalf("无法创建测试目录: %v", err)
	}

	// 创建 ja-netfilter.jar 文件
	jarPath := filepath.Join(validConfigDir, "ja-netfilter.jar")
	if err := os.WriteFile(jarPath, []byte("test"), 0644); err != nil {
		t.Fatalf("无法创建测试文件: %v", err)
	}

	// 创建缺少jar文件的目录
	invalidConfigDir := filepath.Join(tempDir, "invalid-config")
	if err := os.Mkdir(invalidConfigDir, 0755); err != nil {
		t.Fatalf("无法创建测试目录: %v", err)
	}

	tests := []struct {
		name      string
		path      string
		shouldErr bool
		errType   error
	}{
		{
			name:      "有效的配置目录",
			path:      validConfigDir,
			shouldErr: false,
		},
		{
			name:      "不存在的路径",
			path:      filepath.Join(tempDir, "not-exist"),
			shouldErr: true,
			errType:   ErrPathNotExist,
		},
		{
			name:      "缺少jar文件",
			path:      invalidConfigDir,
			shouldErr: true,
			errType:   ErrMissingJarFile,
		},
		{
			name:      "路径是文件而非目录",
			path:      jarPath,
			shouldErr: true,
			errType:   ErrPathNotDir,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateConfigPath(tt.path)
			if tt.shouldErr {
				if err == nil {
					t.Error("期望错误但没有返回错误")
				}
			} else {
				if err != nil {
					t.Errorf("不期望错误但返回了: %v", err)
				}
			}
		})
	}
}

// TestDirectoryHasVMOptions 测试检查目录是否包含 .vmoptions 文件
func TestDirectoryHasVMOptions(t *testing.T) {
	// 创建临时测试目录
	tempDir := t.TempDir()

	// 创建包含 .vmoptions 文件的目录
	withVMOptions := filepath.Join(tempDir, "with-vmoptions")
	if err := os.Mkdir(withVMOptions, 0755); err != nil {
		t.Fatalf("无法创建测试目录: %v", err)
	}
	if err := os.WriteFile(filepath.Join(withVMOptions, "idea.vmoptions"), []byte("test"), 0644); err != nil {
		t.Fatalf("无法创建测试文件: %v", err)
	}

	// 创建不包含 .vmoptions 文件的目录
	withoutVMOptions := filepath.Join(tempDir, "without-vmoptions")
	if err := os.Mkdir(withoutVMOptions, 0755); err != nil {
		t.Fatalf("无法创建测试目录: %v", err)
	}

	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{
			name:     "包含vmoptions文件",
			path:     withVMOptions,
			expected: true,
		},
		{
			name:     "不包含vmoptions文件",
			path:     withoutVMOptions,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := directoryHasVMOptions(tt.path)
			if err != nil {
				t.Errorf("不期望错误: %v", err)
			}
			if result != tt.expected {
				t.Errorf("directoryHasVMOptions(%q) = %v, expected %v", tt.path, result, tt.expected)
			}
		})
	}
}
