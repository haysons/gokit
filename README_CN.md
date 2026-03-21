# gokit

[English](README.md)

[![Go Report Card](https://goreportcard.com/badge/github.com/haysons/gokit)](https://goreportcard.com/report/github.com/haysons/gokit)
[![MIT License](https://img.shields.io/badge/license-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT)
[![Go Reference](https://pkg.go.dev/badge/github.com/haysons/gokit.svg)](https://pkg.go.dev/github.com/haysons/gokit)

**gokit** 是一个 Go 语言工具库，专为生产级服务设计。模块化、类型安全、经过充分测试。

## 特性

- **模块化** - 按需引入，无侵入
- **类型安全** - 完整支持 Go 泛型（Go 1.23+）
- **生产就绪** - 已在生产环境验证
- **依赖注入** - 内置 uber/fx 支持

## 安装

```bash
go get github.com/haysons/gokit
```

## 快速开始

```go
// 配置管理（支持热重载）
cfg := config.New[AppConfig]()
cfg.SetFile("./config.yaml")
cfg.Load()

// 结构化日志
logger := log.GetDefaultSlog()
logger.Info("server started", "port", 8080)

// 增强错误处理
err := errors.Wrap(dbErr, "数据库连接失败")
if errors.IsNotFound(err) { /* 处理不存在 */ }
```

## 模块

| 模块 | 描述 |
|------|------|
| `app` | 应用生命周期管理（uber/fx） |
| `config` | Viper 配置管理，支持热重载 |
| `log` | slog 结构化日志 |
| `errors` | 业务码 + 堆栈 + 用户提示 |
| `middleware` | HTTP/gRPC 中间件（认证、日志、追踪） |
| `distributed` | etcd 分布式工具（锁、选举、队列、计数器） |
| `transport` | 统一传输层 |
| `metadata` | 上下文元数据 |
| `util` | 工具函数（加密、ID、哈希、切片） |

## 项目结构

```
gokit/
├── app/            # 应用框架
├── config/         # 配置管理
├── distributed/   # 分布式工具
├── errors/        # 错误处理
├── log/           # 日志组件
├── middleware/    # 中间件
├── transport/     # 传输层
├── metadata/      # 元数据
├── constraints/   # 泛型约束
└── util/          # 工具函数
```

## License

MIT 许可证 - 详见 [LICENSE](LICENSE)。

---

**Made with ❤️ by [@haysons](https://github.com/haysons)**
