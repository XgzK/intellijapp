package service

import "errors"

// Sentinel errors - 定义可复用的错误类型
// 优化：将错误定义独立到单独文件，提高代码组织性
var (
	ErrEmptyPath        = errors.New("路径不能为空")
	ErrPathNotExist     = errors.New("路径不存在")
	ErrPathNotDir       = errors.New("路径必须是目录")
	ErrNotIntelliJDir   = errors.New("非IntelliJ系列软件安装路径")
	ErrNoVMOptions      = errors.New("未找到任何 .vmoptions 文件")
	ErrMissingJarFile   = errors.New("配置目录缺少 ja-netfilter.jar 文件")
	ErrPermissionDenied = errors.New("权限不足")
)

// toolAddedLines 定义本工具添加的特定配置行（使用包级变量避免重复创建）
var toolAddedLines = map[string]struct{}{
	"--add-opens=java.base/jdk.internal.org.objectweb.asm=ALL-UNNAMED":      {},
	"--add-opens=java.base/jdk.internal.org.objectweb.asm.tree=ALL-UNNAMED": {},
}
