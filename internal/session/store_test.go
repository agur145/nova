package session

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cloudwego/eino/schema"
)

func TestClearMarkerKeepsHistoryAndLimitsEffectiveContext(t *testing.T) {
	store, err := NewStore(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	sess, err := store.GetOrCreate("default")
	if err != nil {
		t.Fatal(err)
	}
	if err := sess.Append(schema.UserMessage("清理前用户")); err != nil {
		t.Fatal(err)
	}
	if err := sess.Append(schema.AssistantMessage("清理前助手", nil)); err != nil {
		t.Fatal(err)
	}
	if err := sess.Clear(); err != nil {
		t.Fatal(err)
	}
	if err := sess.Append(schema.UserMessage("清理后用户")); err != nil {
		t.Fatal(err)
	}

	all := sess.GetMessages()
	if len(all) != 3 {
		t.Fatalf("clear 不应删除历史消息，实际消息数: %d", len(all))
	}
	effective := sess.GetEffectiveMessages()
	if len(effective) != 1 || effective[0].Content != "清理后用户" {
		t.Fatalf("有效上下文应只包含 clear 后消息: %#v", effective)
	}
	history := sess.History()
	if len(history) != 4 || history[2].Type != "clear" {
		t.Fatalf("历史中应保留 clear 分界: %#v", history)
	}
}

func TestLoadLegacyJSONLWithoutClearMarkerUsesFullHistory(t *testing.T) {
	dir := t.TempDir()
	legacy := strings.Join([]string{
		`{"type":"session","id":"legacy","created_at":"2026-01-01T00:00:00Z"}`,
		`{"role":"user","content":"旧问题"}`,
		`{"role":"assistant","content":"旧回答"}`,
		"",
	}, "\n")
	if err := os.WriteFile(filepath.Join(dir, "legacy.jsonl"), []byte(legacy), 0o644); err != nil {
		t.Fatal(err)
	}
	store, err := NewStore(dir)
	if err != nil {
		t.Fatal(err)
	}
	sess, err := store.Get("legacy")
	if err != nil {
		t.Fatal(err)
	}

	effective := sess.GetEffectiveMessages()
	if len(effective) != 2 {
		t.Fatalf("旧文件无 clear 标记时应全部作为有效上下文: %d", len(effective))
	}
	if got := sess.Title(); got != "旧问题" {
		t.Fatalf("旧文件应从首条用户消息推导标题: %s", got)
	}
}

func TestMultipleSessionsAreIsolatedAndActiveSessionPersists(t *testing.T) {
	store, err := NewStore(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	first, err := store.GetOrCreate("default")
	if err != nil {
		t.Fatal(err)
	}
	if err := first.Append(schema.UserMessage("会话 A")); err != nil {
		t.Fatal(err)
	}
	second, err := store.Create("会话 B")
	if err != nil {
		t.Fatal(err)
	}
	if err := second.Append(schema.UserMessage("会话 B")); err != nil {
		t.Fatal(err)
	}
	if err := store.SetActiveID(second.ID); err != nil {
		t.Fatal(err)
	}

	reloaded, err := NewStore(store.dir)
	if err != nil {
		t.Fatal(err)
	}
	active, err := reloaded.GetActiveOrCreate()
	if err != nil {
		t.Fatal(err)
	}
	if active.ID != second.ID {
		t.Fatalf("应恢复最近激活会话: want=%s got=%s", second.ID, active.ID)
	}
	if active.GetMessages()[0].Content != "会话 B" {
		t.Fatalf("激活会话上下文不应串到其他会话: %#v", active.GetMessages())
	}

	metas, err := reloaded.List(active.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(metas) != 2 {
		t.Fatalf("应列出两个会话: %#v", metas)
	}
	if !metas[0].Active {
		t.Fatalf("会话列表应标记当前激活会话: %#v", metas)
	}
}

func TestDeleteRejectsOnlySession(t *testing.T) {
	store, err := NewStore(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.GetOrCreate("default"); err != nil {
		t.Fatal(err)
	}
	if err := store.Delete("default"); err == nil {
		t.Fatal("删除唯一会话应失败")
	}
}
