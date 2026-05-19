package agent

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
	"unicode/utf8"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/schema"

	"nova/internal/book"
	"nova/internal/session"
)

const (
	maxReferenceFileBytes       = 80 * 1024
	maxReferenceTotalBytes      = 200 * 1024
	maxStyleReferenceFileBytes  = 80 * 1024
	maxStyleReferenceTotalBytes = 200 * 1024
)

// Event 表示 Agent 输出的传输无关事件。
type Event struct {
	Type string
	Data interface{}
}

// ChatRequest 表示一次聊天请求的传输无关参数。
type ChatRequest struct {
	Message         string             `json:"message"`
	References      []string           `json:"references"`
	StyleReferences []string           `json:"style_references"`
	Selections      []TextSelectionRef `json:"selections"`
	PlanMode        bool               `json:"plan_mode"`
}

// TextSelectionRef 表示用户在编辑器中选中的一段文本引用。
type TextSelectionRef struct {
	FileName  string `json:"file_name"`
	StartLine int    `json:"start_line"`
	EndLine   int    `json:"end_line"`
	Content   string `json:"content"`
}

// ChatService 编排会话历史、文件引用和 Agent 流式响应。
type ChatService struct{}

// NewChatService 创建聊天服务。
func NewChatService() *ChatService {
	return &ChatService{}
}

// Run 运行一次聊天请求，并通过 emit 输出流式事件。
func (s *ChatService) Run(
	ctx context.Context,
	runner *adk.Runner,
	sess *session.Session,
	bookService *book.Service,
	req ChatRequest,
	emit func(Event),
) {
	agentMessage := req.Message
	if req.PlanMode {
		agentMessage = appendPlanModeInstruction(agentMessage)
	}
	if len(req.References) > 0 {
		agentMessage = appendReferenceContext(bookService, req.Message, req.References)
	}
	if len(req.StyleReferences) > 0 {
		agentMessage = appendStyleReferenceContext(bookService, agentMessage, req.StyleReferences)
	}
	if len(req.Selections) > 0 {
		agentMessage = appendSelectionContext(agentMessage, req.Selections)
	}

	_ = sess.Append(schema.UserMessage(req.Message))
	history := append([]*schema.Message(nil), sess.GetEffectiveMessages()...)
	if len(history) > 0 {
		history[len(history)-1] = schema.UserMessage(agentMessage)
	}

	events := runner.Run(ctx, history)
	var fullContent strings.Builder
	log.Printf("[agent-run] started history=%d message_len=%d agent_message_len=%d plan_mode=%v", len(history), len(req.Message), len(agentMessage), req.PlanMode)

	for {
		if err := ctx.Err(); err != nil {
			log.Printf("[agent-run] interrupted reason=context err=%v generated_bytes=%d", err, fullContent.Len())
			appendAssistantIfAny(sess, &fullContent)
			emit(Event{Type: "aborted", Data: map[string]string{}})
			return
		}
		event, ok := events.Next()
		if !ok {
			break
		}
		if event.Err != nil {
			log.Printf("[agent-run] interrupted reason=runner_error err=%v generated_bytes=%d", event.Err, fullContent.Len())
			appendAssistantIfAny(sess, &fullContent)
			emit(Event{Type: "error", Data: map[string]string{"message": event.Err.Error()}})
			return
		}

		if event.Output == nil || event.Output.MessageOutput == nil {
			log.Printf("[agent-run] skip invalid_output output_nil=%v message_output_nil=%v", event.Output == nil, event.Output != nil && event.Output.MessageOutput == nil)
			continue
		}

		mv := event.Output.MessageOutput
		if mv.Role == schema.Tool {
			if mv.Message == nil {
				continue
			}
			content := drainContent(mv)
			if content == "" {
				content = "(无返回内容)"
			}
			if len(content) > 300 {
				content = content[:300] + "..."
			}
			logToolResult(mv.Message.ToolName, mv.Message.ToolCallID, content)
			emit(Event{Type: "tool_result", Data: map[string]string{
				"id":      mv.Message.ToolCallID,
				"name":    mv.Message.ToolName,
				"content": content,
			}})
			continue
		}

		if mv.Role != schema.Assistant && mv.Role != "" {
			continue
		}
		if mv.IsStreaming && mv.MessageStream != nil {
			if !processStreamingEvent(mv, &fullContent, emit) {
				appendAssistantIfAny(sess, &fullContent)
				return
			}
			continue
		}
		if mv.Message != nil {
			processNonStreamingEvent(mv, &fullContent, emit)
		}
	}

	appendAssistantIfAny(sess, &fullContent)
	log.Printf("[agent-run] completed")
	emit(Event{Type: "done", Data: map[string]string{}})
}

// appendAssistantIfAny 将已生成的正文持久化，避免异常中断后刷新丢失输出。
func appendAssistantIfAny(sess *session.Session, content *strings.Builder) {
	if content == nil || content.Len() == 0 {
		return
	}
	_ = sess.Append(schema.AssistantMessage(content.String(), nil))
	log.Printf("[agent-run] persisted assistant message bytes=%d", content.Len())
	content.Reset()
}

// appendReferenceContext 将用户引用的文件内容追加到本次 Agent 输入。
func appendReferenceContext(bookService *book.Service, message string, references []string) string {
	var sb strings.Builder
	sb.WriteString(message)
	sb.WriteString("\n\n---\n以下是用户引用的文件：\n")

	total := 0
	seen := make(map[string]bool)
	for _, ref := range references {
		ref = strings.TrimSpace(ref)
		if ref == "" || seen[ref] {
			continue
		}
		seen[ref] = true

		sb.WriteString("\n## @")
		sb.WriteString(ref)
		sb.WriteString("\n")

		if total >= maxReferenceTotalBytes {
			sb.WriteString("引用内容总量已超过限制，后续文件未读取。\n")
			continue
		}

		content, n, err := readReferencedFile(bookService, ref, maxReferenceFileBytes, maxReferenceTotalBytes-total)
		total += n
		if err != nil {
			sb.WriteString("读取失败：")
			sb.WriteString(err.Error())
			sb.WriteString("\n")
			continue
		}

		sb.WriteString("```markdown\n")
		sb.WriteString(content)
		sb.WriteString("\n```\n")
	}

	return sb.String()
}

// appendStyleReferenceContext 将本轮指定的风格参考追加到 Agent 输入。
func appendStyleReferenceContext(bookService *book.Service, message string, styleReferences []string) string {
	var sb strings.Builder
	sb.WriteString(message)
	sb.WriteString("\n\n---\n以下是用户本轮指定的风格参考。请只把它们作为文风、节奏、叙述方式、句式和氛围参考，不要照搬内容、人物、情节或设定：\n")

	total := 0
	seen := make(map[string]bool)
	for _, ref := range styleReferences {
		ref = strings.TrimSpace(ref)
		if ref == "" || seen[ref] {
			continue
		}
		seen[ref] = true

		sb.WriteString("\n## #")
		sb.WriteString(ref)
		sb.WriteString("\n")

		if total >= maxStyleReferenceTotalBytes {
			sb.WriteString("风格参考内容总量已超过限制，后续文件未读取。\n")
			continue
		}

		content, n, err := readStyleReferencedFile(bookService, ref, maxStyleReferenceFileBytes, maxStyleReferenceTotalBytes-total)
		total += n
		if err != nil {
			sb.WriteString("读取失败：")
			sb.WriteString(err.Error())
			sb.WriteString("\n")
			continue
		}

		sb.WriteString("```markdown\n")
		sb.WriteString(content)
		sb.WriteString("\n```\n")
	}

	return sb.String()
}

// appendSelectionContext 将用户在编辑器中选中的文本片段追加到消息上下文。
func appendSelectionContext(message string, selections []TextSelectionRef) string {
	var sb strings.Builder
	sb.WriteString(message)
	sb.WriteString("\n\n---\n以下是用户在编辑器中选中的文本片段，请针对这些内容进行操作：\n")

	for _, sel := range selections {
		sb.WriteString("\n## 选中内容来自 ")
		sb.WriteString(sel.FileName)
		sb.WriteString(fmt.Sprintf(":L%d-L%d\n", sel.StartLine, sel.EndLine))
		sb.WriteString("```\n")
		sb.WriteString(sel.Content)
		sb.WriteString("\n```\n")
	}

	return sb.String()
}

// readReferencedFile 安全读取引用文件，并按单文件和总大小限制截断。
func readReferencedFile(bookService *book.Service, relPath string, fileLimit, remainLimit int) (string, int, error) {
	limit := fileLimit
	if remainLimit < limit {
		limit = remainLimit
	}
	if limit <= 0 {
		return "", 0, errors.New("引用内容总量已超过限制")
	}

	content, err := bookService.ReadFile(relPath)
	if err != nil {
		return "", 0, err
	}

	data := []byte(content)
	truncated := false
	if len(data) > limit {
		data = data[:limit]
		truncated = true
	}

	result := string(data)
	if truncated {
		result += "\n\n[内容已截断]"
	}
	return result, len(data), nil
}

// readStyleReferencedFile 安全读取风格参考文件，并按单文件和总大小限制截断。
func readStyleReferencedFile(bookService *book.Service, stylePath string, fileLimit, remainLimit int) (string, int, error) {
	limit := fileLimit
	if remainLimit < limit {
		limit = remainLimit
	}
	if limit <= 0 {
		return "", 0, errors.New("风格参考内容总量已超过限制")
	}

	content, err := bookService.ReadStyleFile(stylePath)
	if err != nil {
		return "", 0, err
	}

	data := []byte(content)
	truncated := false
	if len(data) > limit {
		data = data[:limit]
		truncated = true
	}

	result := string(data)
	if truncated {
		result += "\n\n[内容已截断]"
	}
	return result, len(data), nil
}

// processStreamingEvent 处理流式助手消息，输出领域事件。
// 工具调用在流中一检测到名称就立即 emit，让前端尽早展示 running 卡片。
// 参数在流中逐帧 emit tool_args_delta，前端可实时展示 write_file 内容。
func processStreamingEvent(mv *adk.MessageVariant, fullContent *strings.Builder, emit func(Event)) bool {
	mv.MessageStream.SetAutomaticClose()
	var accumulatedToolCalls []schema.ToolCall
	emittedTools := make(map[int]bool) // 按 index 记录已 emit tool_call 的工具
	lastArgsLen := make(map[int]int)   // 记录上次已发送的参数长度

	for {
		frame, err := mv.MessageStream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			log.Printf("[agent-run] interrupted reason=stream_recv_error err=%v generated_bytes=%d", err, fullContent.Len())
			emit(Event{Type: "error", Data: map[string]string{"message": err.Error()}})
			return false
		}
		if frame == nil {
			continue
		}
		if frame.ReasoningContent != "" {
			emit(Event{Type: "thinking", Data: map[string]string{"content": frame.ReasoningContent}})
		}
		if frame.Content != "" {
			fullContent.WriteString(frame.Content)
			emit(Event{Type: "chunk", Data: map[string]string{"content": frame.Content}})
		}
		if len(frame.ToolCalls) > 0 {
			accumulatedToolCalls = mergeToolCalls(accumulatedToolCalls, frame.ToolCalls)
			for i, tc := range accumulatedToolCalls {
				if tc.Function.Name == "" {
					continue
				}
				// 首次检测到工具名称，emit tool_call
				if !emittedTools[i] {
					emittedTools[i] = true
					lastArgsLen[i] = 0
					logToolCall(tc.Function.Name, tc.ID, len(tc.Function.Arguments), "streaming")
					data := map[string]interface{}{
						"id":   tc.ID,
						"name": tc.Function.Name,
						"args": "",
					}
					if tc.Index != nil {
						data["index"] = *tc.Index
					}
					emit(Event{Type: "tool_call", Data: data})
				}
				// 参数有增量时 emit tool_args_delta
				currentLen := len(tc.Function.Arguments)
				if currentLen > lastArgsLen[i] {
					delta := tc.Function.Arguments[lastArgsLen[i]:currentLen]
					lastArgsLen[i] = currentLen
					data := map[string]interface{}{
						"id":    tc.ID,
						"name":  tc.Function.Name,
						"delta": delta,
					}
					if tc.Index != nil {
						data["index"] = *tc.Index
					}
					emit(Event{Type: "tool_args_delta", Data: data})
				}
			}
		}
	}
	return true
}

// processNonStreamingEvent 处理非流式助手消息，输出领域事件。
func processNonStreamingEvent(mv *adk.MessageVariant, fullContent *strings.Builder, emit func(Event)) {
	if mv.Message.ReasoningContent != "" {
		emit(Event{Type: "thinking", Data: map[string]string{"content": mv.Message.ReasoningContent}})
	}
	if mv.Message.Content != "" {
		fullContent.WriteString(mv.Message.Content)
		emit(Event{Type: "chunk", Data: map[string]string{"content": mv.Message.Content}})
	}
	for _, tc := range mv.Message.ToolCalls {
		name := tc.Function.Name
		if name == "" {
			continue
		}
		args := tc.Function.Arguments
		logToolCall(name, tc.ID, len(args), "non_streaming")
		if len(args) > 200 {
			args = args[:200] + "..."
		}
		data := map[string]interface{}{
			"id":   tc.ID,
			"name": name,
			"args": args,
		}
		if tc.Index != nil {
			data["index"] = *tc.Index
		}
		emit(Event{Type: "tool_call", Data: data})
	}
}

// drainContent 从 MessageVariant 中提取完整内容。
func drainContent(mv *adk.MessageVariant) string {
	if mv.IsStreaming && mv.MessageStream != nil {
		mv.MessageStream.SetAutomaticClose()
		var sb strings.Builder
		for {
			chunk, err := mv.MessageStream.Recv()
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				break
			}
			if chunk != nil && chunk.Content != "" {
				sb.WriteString(chunk.Content)
			}
		}
		return sb.String()
	}
	if mv.Message != nil {
		return mv.Message.Content
	}
	return ""
}

func logToolCall(name, id string, argsBytes int, source string) {
	log.Printf("[agent-tool] call source=%s name=%s id=%s args_bytes=%d", source, name, id, argsBytes)
}

func logToolResult(name, id, content string) {
	if looksLikeToolFailure(content) {
		log.Printf("[agent-tool] result suspected_failure=true name=%s id=%s bytes=%d preview=%q", name, id, len(content), safeLogPreview(content, 300))
		return
	}
	log.Printf("[agent-tool] result name=%s id=%s bytes=%d", name, id, len(content))
}

func looksLikeToolFailure(content string) bool {
	text := strings.ToLower(content)
	failureKeywords := []string{
		"error", "failed", "failure", "panic", "exception", "traceback",
		"permission denied", "not found", "timeout", "timed out",
		"失败", "错误", "异常", "拒绝", "超时", "不存在",
	}
	for _, keyword := range failureKeywords {
		if strings.Contains(text, keyword) {
			return true
		}
	}
	return false
}

func safeLogPreview(content string, limit int) string {
	content = strings.ReplaceAll(content, "\n", "\\n")
	content = strings.ReplaceAll(content, "\r", "\\r")
	if len(content) <= limit {
		return content
	}
	for limit > 0 && !utf8.RuneStart(content[limit]) {
		limit--
	}
	return content[:limit] + "..."
}

// mergeToolCalls 合并流式 frame 中分散的 tool call 信息。
func mergeToolCalls(existing []schema.ToolCall, incoming []schema.ToolCall) []schema.ToolCall {
	for _, tc := range incoming {
		idx := tc.Index
		if idx == nil {
			if tc.Function.Name != "" {
				existing = append(existing, tc)
			}
			continue
		}

		i := *idx
		for len(existing) <= i {
			existing = append(existing, schema.ToolCall{})
		}
		if tc.Function.Name != "" {
			existing[i].Function.Name = tc.Function.Name
		}
		existing[i].Function.Arguments += tc.Function.Arguments
		if tc.ID != "" {
			existing[i].ID = tc.ID
		}
		existing[i].Index = tc.Index
	}
	return existing
}

// appendPlanModeInstruction 在用户消息前追加规划模式指令，允许读取文件但禁止写操作，只输出结构化计划。
func appendPlanModeInstruction(message string) string {
	return `[规划模式] 请你先制定计划，不要执行任何写操作。

要求：
1. 你可以使用 read_file 工具读取文件内容来了解当前状态
2. 分析用户的需求，列出需要完成的步骤
3. 说明每一步涉及哪些文件、要做什么操作
4. 如果有多种方案，列出利弊供用户选择
5. 禁止使用 write_file、edit_file、delete_file 等任何写操作工具
6. 等待用户确认或调整计划后再执行

用户需求：
` + message
}

// EventError 创建标准错误事件。
func EventError(err error) Event {
	return Event{Type: "error", Data: map[string]string{"message": fmt.Sprint(err)}}
}
