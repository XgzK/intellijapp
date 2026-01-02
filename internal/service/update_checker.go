package service

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// ReleaseInfo 保存 GitHub Release 信息
type ReleaseInfo struct {
	Version     string        `json:"version"`
	PublishedAt string        `json:"publishedAt"`
	HTMLURL     string        `json:"htmlUrl"`
	Body        string        `json:"body"`
	Assets      []AssetInfo   `json:"assets"`
}

// AssetInfo 保存 Release 资源信息
type AssetInfo struct {
	Name        string `json:"name"`
	DownloadURL string `json:"downloadUrl"`
	Size        int64  `json:"size"`
}

// UpdateCheckResult 保存更新检查结果
type UpdateCheckResult struct {
	HasUpdate bool         `json:"hasUpdate"`
	Release   *ReleaseInfo `json:"release"`
}

const (
	githubAPI   = "https://api.github.com"
	repoOwner   = "XgzK"
	repoName    = "intellijapp"
	httpTimeout = 10 * time.Second
)

// GitHub 镜像站点列表（按优先级排序）
var githubMirrors = []string{
	"https://github.com",      // 官方站点（最优先）
	"https://2git.xyz",        // 镜像站点 1
	"https://lgithub.xyz",     // 镜像站点 2
}

// GitHub API 镜像站点列表
var githubAPIMirrors = []string{
	"https://api.github.com",  // 官方 API（最优先）
}

// gitHubRelease 对应 GitHub API 返回的 Release 结构
type gitHubRelease struct {
	TagName     string `json:"tag_name"`
	Name        string `json:"name"`
	PublishedAt string `json:"published_at"`
	HTMLURL     string `json:"html_url"`
	Body        string `json:"body"`
	Assets      []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
		Size               int64  `json:"size"`
	} `json:"assets"`
}

// fetchLatestRelease 从 GitHub API 获取最新 Release 信息
// 自动尝试官方 API 和镜像站点
func fetchLatestRelease(logger *slog.Logger) (*ReleaseInfo, error) {
	var lastErr error

	// 依次尝试每个 API 镜像
	for i, apiBase := range githubAPIMirrors {
		url := fmt.Sprintf("%s/repos/%s/%s/releases/latest", apiBase, repoOwner, repoName)
		logger.Debug("尝试 GitHub API", slog.String("url", apiBase), slog.Int("attempt", i+1))

		release, err := fetchFromAPI(url, logger)
		if err == nil {
			logger.Info("成功从 GitHub API 获取版本", slog.String("api", apiBase))
			return release, nil
		}

		lastErr = err
		logger.Warn("GitHub API 请求失败，尝试下一个镜像",
			slog.String("api", apiBase),
			slog.Any("error", err))
	}

	return nil, fmt.Errorf("所有 GitHub API 镜像均不可用: %w", lastErr)
}

// fetchFromAPI 从指定的 API URL 获取 Release 信息
func fetchFromAPI(url string, logger *slog.Logger) (*ReleaseInfo, error) {
	client := &http.Client{Timeout: httpTimeout}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		logger.Info("未找到任何 Release")
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API 返回错误状态: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var ghRelease gitHubRelease
	if err := json.Unmarshal(body, &ghRelease); err != nil {
		return nil, fmt.Errorf("解析 JSON 失败: %w", err)
	}

	// 转换为内部结构
	release := &ReleaseInfo{
		Version:     strings.TrimPrefix(ghRelease.TagName, "v"),
		PublishedAt: ghRelease.PublishedAt,
		HTMLURL:     ghRelease.HTMLURL,
		Body:        ghRelease.Body,
		Assets:      make([]AssetInfo, len(ghRelease.Assets)),
	}

	for i, asset := range ghRelease.Assets {
		release.Assets[i] = AssetInfo{
			Name:        asset.Name,
			DownloadURL: asset.BrowserDownloadURL,
			Size:        asset.Size,
		}
	}

	logger.Debug("成功解析 Release 信息", slog.String("version", release.Version))
	return release, nil
}

// GetAccessibleGitHubMirror 测试并返回可访问的 GitHub 镜像站点
// 依次测试所有镜像站点，返回第一个可访问的
func (c *ConfigService) GetAccessibleGitHubMirror() string {
	client := &http.Client{Timeout: 10 * time.Second}

	// 依次测试每个镜像站点
	for _, mirror := range githubMirrors {
		c.logger.Debug("测试 GitHub 镜像", slog.String("mirror", mirror))

		resp, err := client.Head(mirror)
		if err == nil && resp.StatusCode == http.StatusOK {
			c.logger.Info("找到可访问的 GitHub 镜像", slog.String("mirror", mirror))
			return mirror
		}

		c.logger.Debug("镜像站点不可访问",
			slog.String("mirror", mirror),
			slog.Any("error", err))
	}

	// 如果都不可访问，返回官方站点
	c.logger.Warn("所有镜像站点均不可访问，返回官方站点")
	return githubMirrors[0]
}

// ConvertToAccessibleURL 将 GitHub URL 转换为可访问的镜像 URL
// 供前端调用，传入原始 URL，返回可访问的镜像 URL
func (c *ConfigService) ConvertToAccessibleURL(originalURL string) string {
	// 检查是否是 GitHub URL
	if !strings.HasPrefix(originalURL, "https://github.com") {
		return originalURL
	}

	// 获取可访问的镜像站点
	accessibleMirror := c.GetAccessibleGitHubMirror()

	// 如果是官方站点，直接返回原 URL
	if accessibleMirror == "https://github.com" {
		return originalURL
	}

	// 替换为镜像站点 URL
	mirrorURL := strings.Replace(originalURL, "https://github.com", accessibleMirror, 1)
	c.logger.Info("URL 已转换为镜像站点",
		slog.String("original", originalURL),
		slog.String("mirror", mirrorURL))

	return mirrorURL
}

// compareVersions 比较两个版本号
// 返回: 1 如果 v1 > v2, -1 如果 v1 < v2, 0 如果相等
func compareVersions(v1, v2 string) int {
	// 清理版本号前缀
	v1 = strings.TrimSpace(strings.TrimPrefix(v1, "v"))
	v2 = strings.TrimSpace(strings.TrimPrefix(v2, "v"))

	parts1 := strings.Split(v1, ".")
	parts2 := strings.Split(v2, ".")

	maxLen := len(parts1)
	if len(parts2) > maxLen {
		maxLen = len(parts2)
	}

	for i := 0; i < maxLen; i++ {
		var num1, num2 int

		if i < len(parts1) {
			// 提取数字部分
			numStr := extractNumber(parts1[i])
			num1, _ = strconv.Atoi(numStr)
		}

		if i < len(parts2) {
			numStr := extractNumber(parts2[i])
			num2, _ = strconv.Atoi(numStr)
		}

		if num1 > num2 {
			return 1
		}
		if num1 < num2 {
			return -1
		}
	}

	return 0
}

// extractNumber 从版本号部分提取数字
func extractNumber(part string) string {
	var result strings.Builder
	for _, ch := range part {
		if ch >= '0' && ch <= '9' {
			result.WriteRune(ch)
		} else {
			break
		}
	}
	if result.Len() == 0 {
		return "0"
	}
	return result.String()
}

// CheckForUpdates 检查是否有新版本可用
// 自动使用本地版本号进行比较
func (c *ConfigService) CheckForUpdates() (UpdateCheckResult, error) {
	c.logger.Info("开始检查更新", slog.String("currentVersion", Version))

	release, err := fetchLatestRelease(c.logger)
	if err != nil {
		c.logger.Error("检查更新失败", slog.Any("error", err))
		return UpdateCheckResult{HasUpdate: false, Release: nil}, err
	}

	if release == nil {
		c.logger.Info("未找到 Release")
		return UpdateCheckResult{HasUpdate: false, Release: nil}, nil
	}

	hasUpdate := compareVersions(release.Version, Version) > 0

	c.logger.Info("更新检查完成",
		slog.Bool("hasUpdate", hasUpdate),
		slog.String("latestVersion", release.Version),
		slog.String("currentVersion", Version))

	return UpdateCheckResult{
		HasUpdate: hasUpdate,
		Release:   release,
	}, nil
}
