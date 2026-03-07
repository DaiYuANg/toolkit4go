// sync-versions.go
// 版本文档同步工具
// 用途：从 git tags 自动创建版本文档目录和配置文件

package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

// Version 代表一个版本文档
type Version struct {
	Name    string
	Release string
	Path    string
	Current bool
}

func main() {
	fmt.Println("========================================")
	fmt.Println("   ArcGo 版本文档同步工具")
	fmt.Println("========================================")
	fmt.Println()

	// 获取项目根目录
	projectRoot, err := getProjectRoot()
	if err != nil {
		fmt.Printf("❌ 错误：%v\n", err)
		os.Exit(1)
	}

	docsDir := filepath.Join(projectRoot, "docs")
	contentDir := filepath.Join(docsDir, "content")
	versionsFile := filepath.Join(docsDir, "data", "versions.yaml")

	// 获取所有 git tags
	fmt.Println("[1/4] 获取 git tags...")
	tags, err := getGitTags()
	if err != nil {
		fmt.Printf("❌ 错误：%v\n", err)
		os.Exit(1)
	}

	if len(tags) == 0 {
		fmt.Println("❌ 没有找到任何 git tags")
		os.Exit(1)
	}

	fmt.Println("✅ 找到以下 tags:")
	for _, tag := range tags {
		fmt.Printf("   - %s\n", tag)
	}
	fmt.Println()

	// 获取最新 tag
	latestTag := tags[0]
	fmt.Printf("[2/4] 当前版本：%s\n", latestTag)
	fmt.Println()

	// 创建版本配置
	fmt.Println("[3/4] 创建版本配置文件...")
	versions := createVersionsConfig(tags)
	if err := writeVersionsFile(versionsFile, versions); err != nil {
		fmt.Printf("❌ 错误：%v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✅ 版本文档配置已更新到：%s\n", versionsFile)
	fmt.Println()

	// 创建版本文档目录
	fmt.Println("[4/4] 创建版本文档目录...")
	if err := createVersionedDirs(contentDir, versions, latestTag); err != nil {
		fmt.Printf("❌ 错误：%v\n", err)
		os.Exit(1)
	}
	fmt.Println()

	// 输出统计
	fmt.Println("========================================")
	fmt.Println("   版本统计")
	fmt.Println("========================================")
	fmt.Printf("   当前版本：%s\n", latestTag)
	fmt.Printf("   历史版本数：%d 个\n", len(tags))
	fmt.Println("========================================")
	fmt.Println()
	fmt.Println("💡 提示：运行 'go tool hugo server -D' 预览版本文档")
	fmt.Println()
}

// getProjectRoot 获取项目根目录
func getProjectRoot() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("无法获取项目根目录：%w", err)
	}
	return strings.TrimSpace(string(output)), nil
}

// getGitTags 获取所有 git tags，按版本号降序排序
func getGitTags() ([]string, error) {
	cmd := exec.Command("git", "tag", "--list", "--sort=-version:refname")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("无法获取 git tags：%w", err)
	}

	var tags []string
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		tag := strings.TrimSpace(scanner.Text())
		if tag != "" {
			tags = append(tags, tag)
		}
	}

	return tags, scanner.Err()
}

// createVersionsConfig 创建版本配置
func createVersionsConfig(tags []string) []Version {
	var versions []Version

	for i, tag := range tags {
		version := Version{
			Name:    tag,
			Release: tag,
			Current: i == 0, // 第一个版本为当前版本
		}

		if i == 0 {
			// 当前版本路径为空
			version.Path = ""
		} else {
			// 历史版本路径
			version.Path = fmt.Sprintf("/versioned/%s", tag)
		}

		versions = append(versions, version)
	}

	return versions
}

// writeVersionsFile 写入 versions.yaml 文件
func writeVersionsFile(filename string, versions []Version) error {
	// 确保目录存在
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("无法创建目录：%w", err)
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("无法创建文件：%w", err)
	}
	defer file.Close()

	// 写入文件头
	header := `# 版本文档配置
# 此文件定义了文档的版本列表
# 版本按时间倒序排列，第一个为当前版本

versions:
`
	if _, err := file.WriteString(header); err != nil {
		return fmt.Errorf("写入文件头失败：%w", err)
	}

	// 写入版本配置
	for i, v := range versions {
		var section string
		if v.Current {
			section = "  # 当前版本（最新版本）\n"
		} else {
			section = "\n  # 历史版本\n"
		}

		if _, err := file.WriteString(section); err != nil {
			return fmt.Errorf("写入章节失败：%w", err)
		}

		lines := []string{
			fmt.Sprintf("  - name: \"%s\"\n", v.Name),
			fmt.Sprintf("    release: \"%s\"\n", v.Release),
			fmt.Sprintf("    path: \"%s\"\n", v.Path),
			fmt.Sprintf("    current: %t\n", v.Current),
		}

		for _, line := range lines {
			if _, err := file.WriteString(line); err != nil {
				return fmt.Errorf("写入行失败：%w", err)
			}
		}

		// 如果不是最后一个，添加空行
		if i < len(versions)-1 {
			if _, err := file.WriteString("\n"); err != nil {
				return fmt.Errorf("写入空行失败：%w", err)
			}
		}
	}

	return nil
}

// createVersionedDirs 创建版本文档目录
func createVersionedDirs(contentDir string, versions []Version, latestTag string) error {
	versionedDir := filepath.Join(contentDir, "versioned")

	// 获取源文档目录
	sourceDocsDir := filepath.Join(contentDir, "docs")
	sourceRootFiles := []string{
		filepath.Join(contentDir, "_index.md"),
		filepath.Join(contentDir, "_index.en.md"),
		filepath.Join(contentDir, "_index.zh.md"),
	}

	for _, version := range versions {
		// 跳过当前版本（已经是最新的）
		if version.Current {
			continue
		}

		versionDir := filepath.Join(versionedDir, version.Name)
		versionDocsDir := filepath.Join(versionDir, "docs")

		// 如果目录已存在，跳过
		if _, err := os.Stat(versionDir); err == nil {
			fmt.Printf("   ⏭️  跳过已存在的版本目录：%s\n", version.Name)
			continue
		}

		fmt.Printf("   📁 创建版本文档目录：%s\n", version.Name)

		// 创建目录结构
		if err := os.MkdirAll(versionDocsDir, 0755); err != nil {
			return fmt.Errorf("无法创建目录 %s：%w", versionDir, err)
		}

		// 复制 docs 目录
		if err := copyDir(sourceDocsDir, versionDocsDir); err != nil {
			fmt.Printf("      ⚠️  复制 docs 目录失败：%v\n", err)
		}

		// 复制根文件
		for _, srcFile := range sourceRootFiles {
			if _, err := os.Stat(srcFile); os.IsNotExist(err) {
				continue
			}

			dstFile := filepath.Join(versionDir, filepath.Base(srcFile))
			if err := copyFile(srcFile, dstFile); err != nil {
				fmt.Printf("      ⚠️  复制文件 %s 失败：%v\n", filepath.Base(srcFile), err)
			}
		}
	}

	return nil
}

// copyDir 递归复制目录
func copyDir(src, dst string) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		return fmt.Errorf("读取目录 %s 失败：%w", src, err)
	}

	for _, entry := range entries {
		// 跳过 versioned 目录
		if entry.Name() == "versioned" {
			continue
		}

		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := os.MkdirAll(dstPath, 0755); err != nil {
				return fmt.Errorf("创建目录 %s 失败：%w", dstPath, err)
			}
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// copyFile 复制文件
func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("读取文件 %s 失败：%w", src, err)
	}

	if err := os.WriteFile(dst, data, 0644); err != nil {
		return fmt.Errorf("写入文件 %s 失败：%w", dst, err)
	}

	return nil
}

// 按语义版本号排序（可选工具函数）
func sortVersions(tags []string) {
	sort.Slice(tags, func(i, j int) bool {
		return compareVersions(tags[i], tags[j]) > 0
	})
}

// compareVersions 比较两个版本号
// 返回 1 表示 v1 > v2，-1 表示 v1 < v2，0 表示相等
func compareVersions(v1, v2 string) int {
	// 移除 'v' 前缀
	v1 = strings.TrimPrefix(v1, "v")
	v2 = strings.TrimPrefix(v2, "v")

	parts1 := strings.Split(v1, ".")
	parts2 := strings.Split(v2, ".")

	maxLen := len(parts1)
	if len(parts2) > maxLen {
		maxLen = len(parts2)
	}

	for i := 0; i < maxLen; i++ {
		var n1, n2 int

		if i < len(parts1) {
			fmt.Sscanf(parts1[i], "%d", &n1)
		}
		if i < len(parts2) {
			fmt.Sscanf(parts2[i], "%d", &n2)
		}

		if n1 > n2 {
			return 1
		} else if n1 < n2 {
			return -1
		}
	}

	return 0
}
