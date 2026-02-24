# 基于 go-admin-team 公共代码库做的扩展，适应阿里云、腾讯云服务
### 增加
 - [x] 文件上传（COS、OSS）
 - [x] 短信发送（腾讯云、阿里云）
 - [x] MQ（Rocket、Pulsar）
 - [x] SMTP发件
 - [x] 高德地图地址解析

### 优化
 - [x] Redis增加支持HashSet、HashKeys
 - [x] DTO配置search增加自定义排序模式

### 架构调整（2026-02）
- [x] 移除 `sdk/`、`plugins/` 顶层目录并完成代码归位：
  - `sdk/*` → `core/*`
  - `plugins/logger/zap` → `logger/zap`
  - `plugins/logger/logrus` → `logger/logrus`
- [x] 统一依赖入口，减少版本漂移与 `replace` 链路复杂度
- [x] `runtime` 读路径改为 `RLock`，并修复路由表重复累积问题

#### 为什么这么做
- 单模块更适合当前代码形态：`sdk` 与 `plugins` 已长期强依赖 `core`，独立模块只会放大版本与发布成本。
- 降低维护复杂度：避免“多模块版本不一致 + 本地 replace”的隐性问题，开发与 CI 更稳定。
- 提升并发安全与可预测性：`runtime` 读写锁语义更清晰，路由表不会随着调用次数无上限增长。
