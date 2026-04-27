# chathub 实现计划

**日期：** 2026-04-27  
**设计文档：** `docs/superpowers/specs/2026-04-27-chathub-design.md`

## 阶段概览

1. 项目初始化
2. 核心框架（config + output + Platform 接口）
3. 飞书适配器
4. 企业微信适配器
5. 钉钉适配器
6. 微信适配器
7. CLI 命令层
8. 测试

---

## 阶段 1：项目初始化

**目标：** 建立 Go 模块和目录结构

- [ ] `go mod init github.com/<user>/chathub`
- [ ] 创建目录：`cmd/`, `internal/config/`, `internal/output/`, `internal/platform/lark/`, `internal/platform/wecom/`, `internal/platform/dingtalk/`, `internal/platform/wx/`
- [ ] 添加依赖：`github.com/spf13/cobra`（CLI 框架）
- [ ] 创建 `main.go` 入口

**验收：** `go build ./...` 通过

---

## 阶段 2：核心框架

**目标：** Platform 接口、config 加载、JSON 输出

### 2a. Platform 接口 (`internal/platform/platform.go`)

```go
type Message struct {
    ID     string
    From   string
    Text   string
    Time   time.Time
}

type Platform interface {
    SendMessage(to, text string) error
    ReadMessages(from string, limit int) ([]Message, error)
}
```

### 2b. Config 加载 (`internal/config/config.go`)

- 读取 `~/.config/chathub/config.json`
- 环境变量覆盖（`CHATHUB_LARK_APP_ID` 等）
- 暴露 `Load() (*Config, error)`

### 2c. JSON 输出 (`internal/output/output.go`)

- `Success(data any)` → stdout `{"ok":true,"data":...}`
- `Failure(err error)` → stderr 错误信息 + exit 1

**验收：** 单元测试覆盖 config 加载（文件 + 环境变量优先级）

---

## 阶段 3：飞书适配器 (`internal/platform/lark/`)

- 实现 `Platform` 接口
- `SendMessage`：调用飞书发消息 API
- `ReadMessages`：调用飞书拉取消息 API
- 凭证：`app_id` + `app_secret` → access token

**验收：** mock 单元测试通过

---

## 阶段 4：企业微信适配器 (`internal/platform/wecom/`)

- 实现 `Platform` 接口
- `SendMessage`：调用企微发消息 API
- `ReadMessages`：调用企微消息接收 API
- 凭证：`corp_id` + `secret`

**验收：** mock 单元测试通过

---

## 阶段 5：钉钉适配器 (`internal/platform/dingtalk/`)

- 实现 `Platform` 接口
- `SendMessage`：调用钉钉发消息 API
- `ReadMessages`：调用钉钉消息 API
- 凭证：`app_key` + `app_secret`

**验收：** mock 单元测试通过

---

## 阶段 6：微信适配器 (`internal/platform/wx/`)

- 实现 `Platform` 接口
- `SendMessage`：调用微信发消息 API
- `ReadMessages`：调用微信消息 API
- 凭证：`token`

**验收：** mock 单元测试通过

---

## 阶段 7：CLI 命令层 (`cmd/`)

### 文件结构

```
cmd/
├── root.go          # 根命令，--output flag
├── lark.go          # chathub lark send/read
├── wecom.go         # chathub wecom send/read
├── dingtalk.go      # chathub dingtalk send/read
├── wx.go            # chathub wx send/read
└── config.go        # chathub config set/show
```

### 每个平台命令

```
chathub <platform> send --to <id> --text <msg>
chathub <platform> read --from <id> --limit <n>
```

### config 命令

```
chathub config set <key> <value>   # 写入 ~/.config/chathub/config.json
chathub config show                # 打印配置（token/secret 脱敏为 ***）
```

**验收：** `chathub --help` 显示所有子命令；`go build` 通过

---

## 阶段 8：测试

- [ ] `internal/config`: 文件加载 + 环境变量覆盖优先级
- [ ] 每个平台适配器：mock HTTP 客户端，覆盖 send/read 正常路径和错误路径
- [ ] `internal/output`: Success/Failure 输出格式

**验收：** `go test ./...` 全部通过

---

## 依赖清单

| 包 | 用途 |
|----|------|
| `github.com/spf13/cobra` | CLI 框架 |
| `encoding/json` | 标准库，JSON 输出 |
| `net/http` | 标准库，平台 API 调用 |
| `os` | 标准库，环境变量读取 |

不引入额外第三方 HTTP 客户端，使用标准库 `net/http`。

---

## 执行顺序

阶段 1 → 2 → 3 → 4 → 5 → 6 → 7 → 8（顺序执行，每阶段验收后进入下一阶段）
