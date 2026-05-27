# Ideas
## WIP
- 基础体验
  - 完善的 Skills 系统
  - /brainstorm /spec /plan /execute /review
  - 章节细纲，从大纲到->章节细纲->章节草稿->章节定稿
- Agent能力
  - 支持完善的上下文管理，memory
  - 考虑：自定义Agent
- ide 日志
- Tantivy / MeiliSearch 全局搜索
- 重构：考虑 Code Mirror6 编辑器，实验效果
- 调试模式，开启后可以看到context组成

## 互动模式
- story teller 可以按自己风格规划暗线剧情和事件，根据用户行为动态调整，保证用户有一个连续的互动体验
- 开局配置
- 互动模式支持通过 agent+skill 来初始化世界观&角色&开局
- 优化导入酒馆v2角色卡，优化内容和世界观
- 目前互动模式不是分三栏吗，左侧这一栏我感觉优化空间很大，比如原来的 creator.md 我感觉可以换成 prompt 高级配置
  - 支持规则列表配置 for story teller，thinking，方便注入配置
  - 支持把关键信息注入到每轮的thinking中，比如字数要求

## NEED FIX
- state agent 应该输出快一点，不要thinking了，直接输出


# 规划
- 多语言支持
- 互动创作模式
- 剧情分支系统，允许从特定节点开始，分出不同的剧情线延续，允许对比不同的分支然后选择一个合并
- 版本管理：不用git，自己实现
- 支持导入小说

- prompt 高级自定义
- 支持在diff view中点击accept/reject按钮，确认或拒绝当前diff
