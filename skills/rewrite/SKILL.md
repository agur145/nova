---
name: rewrite
description: 重写或修改已有章节。当用户要求修改某章、改写某段、换视角重写、调整对话时使用。
---

# 章节重写与修改

根据作者要求修改已有的章节内容。

## 工作流程

1. 使用 read_file 读取目标章节文件 `chapters/chXX-章节名.md`
2. 使用 read_file 读取 `setting/outline.md`、`setting/characters.md` 和 `setting/progress.md` 了解上下文
3. 根据作者要求进行修改（改对话、调情节、换视角、扩写、缩写等）
4. 使用 write_file 写回修改后的章节
5. 更新 `setting/progress.md` 中的章节摘要和最近事件
6. 如果角色状态有变化，同步更新 `setting/characters.md`
7. 除非作者明确要求调整大纲，不更新 `setting/outline.md`

## 修改类型

- **改对话**：调整角色对白，保持性格一致
- **调情节**：修改事件走向，注意前后文连贯
- **换视角**：用不同角色的视角重写，保持信息一致
- **扩写**：在保持主线不变的前提下丰富细节
- **缩写**：精简冗余描写，保留核心情节
- **风格调整**：改变叙事语调或文风

## 注意事项

- 修改后必须与前后章节保持连贯
- 重大情节变化需提醒作者检查后续章节是否需要调整
- 保持角色性格和说话方式的一致性
- 修改后主要更新 `setting/progress.md` 和 `setting/characters.md` 来追踪进展
- 不要把已完成章节复盘写进 `setting/outline.md`，也不要把未来章节规划写进 `setting/characters.md`
- 除非作者明确提出重构主线、调整卷章或修改章节目标，不轻易更新 `setting/outline.md`
