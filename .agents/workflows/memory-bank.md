---
description: 任务开始时读取全部记忆库文件，建立项目上下文
---

# /memory-bank 工作流

每次新会话或新任务开始时，执行以下步骤以建立完整的项目上下文。

## 步骤

1. 读取 `memory-bank/projectbrief.md` — 了解项目核心需求和目标
2. 读取 `memory-bank/productContext.md` — 了解项目动机和运作方式
3. 读取 `memory-bank/systemPatterns.md` — 了解系统架构和设计模式
4. 读取 `memory-bank/techContext.md` — 了解技术栈和依赖项
5. 读取 `memory-bank/activeContext.md` — 了解当前工作重点和最近变更
6. 读取 `memory-bank/progress.md` — 了解已完成功能和待办事项
7. 检查 `memory-bank/` 下是否有附加上下文文件，如有则一并读取
8. 读取 `.agents/AGENTS.md` 中的"项目情报"部分

## 注意事项

- 所有文件都必须读取，不可跳过
- 如果某个文件不存在，记录下来并在后续的 Act 模式中创建
- 读取完成后，在对话中简要确认已建立上下文
