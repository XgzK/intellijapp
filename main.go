package main

import (
	"embed"
	"log/slog"
	"os"

	"github.com/XgzK/intellijapp/internal/service"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// Wails 使用 Go 的 `embed` 包将前端文件嵌入到二进制文件中
// frontend/dist 文件夹中的所有文件都将被嵌入到二进制文件中
// 并可供前端使用
// 更多信息请参阅 https://pkg.go.dev/embed

//go:embed all:frontend/dist
var assets embed.FS

// main 函数作为应用程序的入口点，它初始化应用程序、创建窗口
// 然后运行应用程序并记录可能发生的任何错误
func main() {

	// 通过提供必要的选项创建一个新的 Wails 应用程序
	// 变量 'Name' 和 'Description' 用于应用程序元数据
	// 'Assets' 配置资产服务器，'FS' 变量指向前端文件
	// 'Services' 是 Go 服务实例列表，前端可以访问这些实例的方法
	// 'Mac' 选项用于在 macOS 上运行应用程序时进行定制
	app := application.New(application.Options{
		Name:        "intellijapp",
		Description: "IntelliJ Configuration Helper",
		Services: []application.Service{
			application.NewService(service.NewConfigService()),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	// 使用必要的选项创建一个新窗口
	// 'Title' 是窗口的标题
	// 'Mac' 选项用于在 macOS 上运行时定制窗口
	// 'BackgroundColour' 是窗口的背景颜色
	// 'URL' 是将加载到 webview 中的 URL
	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:     "IntelliJ Config Helper",
		Frameless: true,
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		BackgroundColour: application.NewRGB(27, 38, 54),
		URL:              "/",
	})

	// 运行应用程序。这会阻塞直到应用程序退出
	err := app.Run()

	// 如果运行应用程序时发生错误，使用 slog 记录错误并退出
	if err != nil {
		slog.Error("应用程序运行失败", slog.Any("error", err))
		os.Exit(1)
	}
}
