package app

import (
	"testing"

	"github.com/cloudwego/eino/schema"

	"nova/internal/session"
)

func TestAppSwitchSessionUsesCurrentSessionHistoryOnly(t *testing.T) {
	store, err := session.NewStore(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	app := &App{sessionStore: store}

	first, err := store.GetOrCreate("default")
	if err != nil {
		t.Fatal(err)
	}
	if err := first.Append(schema.UserMessage("会话 A 消息")); err != nil {
		t.Fatal(err)
	}
	app.session = first

	second, err := app.CreateSession("会话 B")
	if err != nil {
		t.Fatal(err)
	}
	if second.ID == first.ID {
		t.Fatal("新会话 ID 不应复用 default")
	}
	if err := second.Append(schema.UserMessage("会话 B 消息")); err != nil {
		t.Fatal(err)
	}

	history, err := app.SessionMessages("")
	if err != nil {
		t.Fatal(err)
	}
	if len(history) != 1 || history[0].Content != "会话 B 消息" {
		t.Fatalf("当前历史应来自新会话: %#v", history)
	}

	if _, err := app.SwitchSession(first.ID); err != nil {
		t.Fatal(err)
	}
	history, err = app.SessionMessages("")
	if err != nil {
		t.Fatal(err)
	}
	if len(history) != 1 || history[0].Content != "会话 A 消息" {
		t.Fatalf("切换后历史应来自目标会话: %#v", history)
	}
}

func TestAppDeleteActiveSessionSwitchesToRemainingSession(t *testing.T) {
	store, err := session.NewStore(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	first, err := store.GetOrCreate("default")
	if err != nil {
		t.Fatal(err)
	}
	app := &App{sessionStore: store, session: first}
	second, err := app.CreateSession("会话 B")
	if err != nil {
		t.Fatal(err)
	}

	active, err := app.DeleteSession(second.ID)
	if err != nil {
		t.Fatal(err)
	}
	if active.ID != first.ID {
		t.Fatalf("删除当前会话后应切换到剩余会话: want=%s got=%s", first.ID, active.ID)
	}
	metas, err := app.Sessions()
	if err != nil {
		t.Fatal(err)
	}
	if len(metas) != 1 || !metas[0].Active || metas[0].ID != first.ID {
		t.Fatalf("剩余会话列表不符合预期: %#v", metas)
	}
}
