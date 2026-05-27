package api

import (
	"bytes"
	"context"
	"errors"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// handleSystemSelectDirectory GET /api/system/select-directory — 调用本机系统目录选择器。
func (s *Server) handleSystemSelectDirectory(ctx context.Context, c *app.RequestContext) {
	path, cancelled, err := selectDirectory(ctx)
	if err != nil {
		writeError(c, consts.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(c, consts.StatusOK, map[string]any{
		"path":      path,
		"cancelled": cancelled,
	})
}

func selectDirectory(parent context.Context) (string, bool, error) {
	ctx, cancel := context.WithTimeout(parent, 5*time.Minute)
	defer cancel()

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		if _, err := exec.LookPath("osascript"); err != nil {
			return "", false, errors.New("当前系统无法打开文件夹选择器：缺少 osascript")
		}
		cmd = exec.CommandContext(ctx, "osascript", "-e", `POSIX path of (choose folder with prompt "选择书籍工作区")`)
	case "linux":
		if path, ok := lookPath("zenity", "kdialog"); ok {
			if path == "zenity" {
				cmd = exec.CommandContext(ctx, path, "--file-selection", "--directory", "--title=选择书籍工作区")
			} else {
				cmd = exec.CommandContext(ctx, path, "--getexistingdirectory", ".", "选择书籍工作区")
			}
		} else {
			return "", false, errors.New("当前系统无法打开文件夹选择器：请安装 zenity 或 kdialog")
		}
	case "windows":
		if path, ok := lookPath("powershell.exe", "powershell"); ok {
			script := `Add-Type -AssemblyName System.Windows.Forms; $dialog = New-Object System.Windows.Forms.FolderBrowserDialog; $dialog.Description = '选择书籍工作区'; if ($dialog.ShowDialog() -eq [System.Windows.Forms.DialogResult]::OK) { Write-Output $dialog.SelectedPath }`
			cmd = exec.CommandContext(ctx, path, "-NoProfile", "-STA", "-Command", script)
		} else {
			return "", false, errors.New("当前系统无法打开文件夹选择器：缺少 PowerShell")
		}
	default:
		return "", false, errors.New("当前系统暂不支持打开文件夹选择器")
	}

	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	selected := strings.TrimSpace(string(out))
	if ctx.Err() != nil {
		return "", false, errors.New("打开文件夹选择器超时")
	}
	if err != nil {
		if selected == "" {
			return "", true, nil
		}
		detail := strings.TrimSpace(stderr.String())
		if detail == "" {
			detail = err.Error()
		}
		return "", false, errors.New("打开文件夹选择器失败: " + detail)
	}
	if selected == "" {
		return "", true, nil
	}
	return selected, false, nil
}

func lookPath(names ...string) (string, bool) {
	for _, name := range names {
		if _, err := exec.LookPath(name); err == nil {
			return name, true
		}
	}
	return "", false
}
