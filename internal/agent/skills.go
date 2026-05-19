package agent

import (
	"os"
	"path/filepath"
)

// ResolveSkillsDir 解析 Skills 目录路径，返回空字符串表示未配置。
func ResolveSkillsDir(dir string) string {
	if dir == "" {
		return ""
	}
	if abs, err := filepath.Abs(dir); err == nil {
		dir = abs
	}
	fi, err := os.Stat(dir)
	if err != nil || !fi.IsDir() {
		return ""
	}
	return dir
}
