# chathub 设计文档

**日期：** 2026-04-27  
**状态：** 已批准

## 概述

`chathub` 是一个统一的 IM 命令行工具，面向 AI Agent 自动化场景。通过单一 CLI 入口，以 JSON stdout 的方式操作企业微信、飞书、钉钉、微信四个平台的消息收发。

## 架构

单体 CLI，内置所有平台适配器。

```
chathub/
├── cmd/           # CLI 入口，命令定义
├── internal/
│   ├── config/    # 凭证加载（文件 + 环境变量）
│   ├── output/    # 统一 JSON 输出格式
│   └── platform/
│       ├── lark/      # 飞书适配器
│       ├── wecom/     # 企业微信适配器
│       ├── dingtalk/  # 钉钉适配器
│       └── wx/        # 微信适配器
└── main.go
```

所有平台实现统一接口：

```go
type Platform interface {
    SendMessage(to string, text string) error
    ReadMessages(from string, limit int) ([]Message, error)
}
```

## 命令结构

```
chathub <platform> <action> [flags]
```

| 命令 | 说明 |
|------|------|
| `chathub lark send --to <id> --text <msg>` | 飞书发消息 |
| `chathub lark read --from <id> --limit 10` | 飞书读消息 |
| `chathub wecom send --to <id> --text <msg>` | 企微发消息 |
| `chathub wecom read --from <id> --limit 10` | 企微读消息 |
| `chathub dingtalk send --to <id> --text <msg>` | 钉钉发消息 |
| `chathub dingtalk read --from <id> --limit 10` | 钉钉读消息 |
| `chathub wx send --to <id> --text <msg>` | 微信发消息 |
| `chathub wx read --from <id> --limit 10` | 微信读消息 |
| `chathub config set lark.app_id <value>` | 写入凭证 |
| `chathub config show` | 查看配置（脱敏） |

默认输出 JSON（`--output json`），支持 `--output text` 供人工阅读。

## 输出格式

```json
{"ok": true, "data": [...]}
{"ok": false, "error": "invalid token"}
```

错误通过 stderr 输出，exit code 非零；stdout 始终是合法 JSON。

## 凭证管理

配置文件：`~/.config/chathub/config.json`

```json
{
  "lark":     { "app_id": "", "app_secret": "" },
  "wecom":    { "corp_id": "", "secret": "" },
  "dingtalk": { "app_key": "", "app_secret": "" },
  "wx":       { "token": "" }
}
```

环境变量（优先级高于文件）：

```
CHATHUB_LARK_APP_ID / CHATHUB_LARK_APP_SECRET
CHATHUB_WECOM_CORP_ID / CHATHUB_WECOM_SECRET
CHATHUB_DINGTALK_APP_KEY / CHATHUB_DINGTALK_APP_SECRET
CHATHUB_WX_TOKEN
```

加载顺序：读取配置文件 → 环境变量覆盖对应字段。

## 错误处理

| 场景 | 行为 |
|------|------|
| 认证失败 | stderr 输出错误，exit 1 |
| 网络超时 | stderr 输出错误，exit 1 |
| 平台 API 错误 | stdout `{"ok":false,"error":"..."}` |
| 缺少必填参数 | stderr 输出用法提示，exit 2 |

## 测试策略

- 平台适配器：接口 mock 单元测试
- config 加载：集成测试（覆盖文件 + 环境变量优先级）
- 真实平台 API：不纳入自动化 CI（需真实凭证）

## 第一版范围

- 四个平台：飞书、企业微信、钉钉、微信
- 两个操作：发消息、读消息
- 凭证管理：文件 + 环境变量
