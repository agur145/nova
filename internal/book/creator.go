package book

import (
	"fmt"
	"os"
	"path/filepath"
)

// InitCreatorPrompt 在 workspace 下创建 CREATOR.md 模板。
func (s *Service) InitCreatorPrompt() string {
	creatorPath := filepath.Join(s.workspace, "CREATOR.md")

	if _, err := os.Stat(creatorPath); err == nil {
		return fmt.Sprintf("CREATOR.md 已存在: %s\n直接编辑该文件即可自定义创作指令", creatorPath)
	}

	template := `# 创作者指令

<!-- 
  这是你的自定义创作指令文件，Nova 每次对话都会读取此文件。
  在这里写下你对 AI 创作的全局要求，它具备最高优先级。
  删除示例内容，写入你自己的指令即可。
-->

## 写作风格

- 文风偏好：（如：简洁有力 / 华丽细腻 / 幽默轻松 / 冷峻克制）
- 叙事视角：（如：第一人称 / 第三人称有限 / 全知视角）
- 对话风格：（如：简短精炼 / 富有个性 / 方言口语化）

## 创作约束

- 每章字数：（如：3000-5000字）
- 禁止内容：（如：不要出现说教性质的旁白）
- 必须遵守：（如：每章结尾留悬念，角色说话要有口癖区分）

## 作品信息

- 类型：（如：都市悬疑 / 玄幻仙侠 / 科幻硬核 / 言情甜文）
- 目标读者：（如：男频 / 女频 / 通用）
- 整体基调：（如：暗黑压抑 / 热血燃向 / 温馨治愈）

## 其他要求

- （写下任何你希望 AI 始终遵守的规则）
`

	if err := os.WriteFile(creatorPath, []byte(template), 0o644); err != nil {
		return fmt.Sprintf("创建 CREATOR.md 失败: %v", err)
	}

	return fmt.Sprintf("已创建 CREATOR.md: %s\n请编辑该文件写入你的创作偏好，Nova 每次对话都会自动读取", creatorPath)
}
