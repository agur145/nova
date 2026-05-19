package agent

import (
	"context"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cloudwego/eino/adk"

	"nova/internal/book"
)

// TestHandleUnknownTool 验证 LLM 幻觉调用不存在工具时，处理器返回引导性
// ToolMessage 而不是抛出错误，从而让 Agent 自行修正。
func TestHandleUnknownTool(t *testing.T) {
	result, err := handleUnknownTool(context.Background(), "write_todo", `{"todos":[]}`)
	if err != nil {
		t.Fatalf("处理未知工具不应返回错误: %v", err)
	}
	if !strings.Contains(result, "write_todo") {
		t.Fatalf("结果应包含工具名: %s", result)
	}
	if !strings.Contains(result, "[tool error]") {
		t.Fatalf("结果应携带 [tool error] 前缀以提示模型自我修复: %s", result)
	}
}

func TestShouldSnapshotBeforeChapterWrite(t *testing.T) {
	workspace := t.TempDir()
	toolCtx := &adk.ToolContext{Name: "write_file"}

	cases := []struct {
		name string
		ctx  *adk.ToolContext
		args string
		want bool
	}{
		{
			name: "相对路径章节",
			ctx:  toolCtx,
			args: `{"file_path":"chapters/ch01.md","content":"正文"}`,
			want: true,
		},
		{
			name: "绝对路径章节",
			ctx:  toolCtx,
			args: `{"file_path":"` + filepath.Join(workspace, "chapters/ch02.md") + `","content":"正文"}`,
			want: true,
		},
		{
			name: "非章节文件",
			ctx:  toolCtx,
			args: `{"file_path":"setting/progress.md","content":"进度"}`,
			want: false,
		},
		{
			name: "非 write_file 工具",
			ctx:  &adk.ToolContext{Name: "edit_file"},
			args: `{"file_path":"chapters/ch01.md","content":"正文"}`,
			want: false,
		},
		{
			name: "非法参数",
			ctx:  toolCtx,
			args: `{bad json}`,
			want: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := shouldSnapshotBeforeChapterWrite(workspace, tc.ctx, tc.args)
			if got != tc.want {
				t.Fatalf("判断结果不符合预期: got=%v want=%v", got, tc.want)
			}
		})
	}
}

func TestAutoCommitBeforeChapterWriteInitializesAndCommitsDirtyWorkspace(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("当前环境未安装 git")
	}
	workspace := t.TempDir()
	mustWriteTestFile(t, workspace, "setting/progress.md", "原始进度\n")

	if err := autoCommitBeforeChapterWrite(context.Background(), workspace); err != nil {
		t.Fatalf("自动快照提交失败: %v", err)
	}

	gitService := book.NewGitService(workspace)
	status, err := gitService.Status(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if !status.Initialized || !status.Clean {
		t.Fatalf("自动快照后应初始化且干净: %#v", status)
	}
	history, err := gitService.History(context.Background(), 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(history) != 1 || !strings.HasPrefix(history[0].Subject, "Nova 自动快照：写章节前") {
		t.Fatalf("自动快照提交信息不符合预期: %#v", history)
	}
}
