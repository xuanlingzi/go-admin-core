---
description: 审查并更新所有记忆库文件，确保文档与项目当前状态一致
---

# /update-memory-bank 工作流

当需要更新记忆库时（用户请求、重大变更后、发现新模式时），执行以下步骤。

## 步骤

1. **审查所有核心文件**（即使不需要更新，也必须全部审查）
   - `memory-bank/projectbrief.md`
   - `memory-bank/productContext.md`
   - `memory-bank/systemPatterns.md`
   - `memory-bank/techContext.md`
   - `memory-bank/activeContext.md`
   - `memory-bank/progress.md`

2. **记录当前状态**
   - 重点更新 `activeContext.md`：当前工作重点、最近变更、下一步计划
   - 重点更新 `progress.md`：已完成功能、待办事项、已知问题

3. **澄清下一步**
   - 在 `activeContext.md` 中明确记录下一步行动计划

4. **更新 .agents 规则**
   - 如果发现新的项目模式、用户偏好或关键见解，更新 `.agents/AGENTS.md` 中的"项目情报"部分

## 注意事项

- 记录"事实 + 决策 + 原因"，避免空泛描述
- 如果某个核心文件不存在，立即创建并填充当前已知内容
- 附加上下文文件也应一并审查和更新
