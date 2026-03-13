# gokit

[![Go Report Card](https://goreportcard.com/badge/github.com/haysons/gokit)](https://goreportcard.com/report/github.com/haysons/gokit)
[![MIT License](https://img.shields.io/badge/license-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT)
[![Go Reference](https://pkg.go.dev/badge/github.com/haysons/gokit.svg)](https://pkg.go.dev/github.com/haysons/gokit)

**gokit** 是一个 Go 语言工具包，提供丰富的辅助函数和组件，简化开发流程并提升代码效率。

## 📦 特性

- 🧩 **模块化设计** - 按需使用，无侵入
- 🔧 **实用工具** - 覆盖常用开发场景
- 📝 **类型安全** - 充分利用 Go 类型系统
- 🚀 **生产就绪** - 经过充分测试和验证

## 📥 安装

```bash
go get github.com/haysons/gokit
```

## 📚 模块概览

### 📋 核心模块

| 模块 | 描述 |
|------|------|
| [`config`](#config-配置管理) | 基于 Viper 的配置管理，支持热重载 |
| [`log`](#log-日志组件) | 基于 slog 的结构化日志 |
| [`errors`](#errors-错误处理) | 增强的错误处理 |
| [`middleware`](#middleware-中间件) | HTTP/gRPC 中间件集合 |
| [`transport`](#transport-传输层) | HTTP/gRPC 传输层封装 |
| [`util`](#util-工具函数) | 通用工具函数集合 |
| [`app`](#app-应用框架) | 应用生命周期管理 |
| [`distributed`](#distributed-分布式) | 分布式系统工具 |
| [`metadata`](#metadata-元数据) | 上下文元数据管理 |
| [`constraints`](#constraints-泛型约束) | 泛型编程约束定义 |

---

## 📖 使用示例

### Config 配置管理

基于 `viper` 的配置管理，支持多格式配置文件、环境变量、热重载：

```go
package main

import (
    "github.com/haysons/gokit/config"
)

type AppConfig struct {
    Server struct {
        Port int    `mapstructure:"port"`
        Host string `mapstructure:"host"`
    }
    Database struct {
        DSN string `mapstructure:"dsn"`
    }
}

func main() {
    cfg := config.New[AppConfig]()
    cfg.SetFile("./config.yaml")
    cfg.SetType("yaml")
    
    // 加载配置
    if err := cfg.Load(); err != nil {
        panic(err)
    }
    
    // 获取配置
    conf := cfg.Get()
    println(conf.Server.Port)
    
    // 监听配置变化（热重载）
    cfg.Watch()
    
    // 等待...
    select {}
}
```

**config.yaml 示例：**
```yaml
server:
  port: 8080
  host: localhost
database:
  dsn: "user:pass@tcp(127.0.0.1:3306)/db"
```

### Log 日志组件

基于 `slog` 的结构化日志，支持 JSON/文本格式、文件轮转：

```go
package main

import (
    "github.com/haysons/gokit/log"
)

func main() {
    // 使用默认配置（终端彩色输出）
    logger := log.GetDefaultSlog()
    
    logger.Info("server started", "port", 8080)
    logger.Error("connection failed", "error", err)
    
    // 自定义配置
    log.SetDefaultSlog(&log.Config{
        Level:        "debug",
        Filename:     "/var/log/app.log",
        MaxAge:       30,
        ConsoleFmt:   true,
        ConsoleColor: true,
    })
}
```

### Util 工具函数

#### 加密工具 (`util/crypto`)
```go
import "github.com/haysons/gokit/util/crypto"

// AES 加密/解密
encrypted, _ := crypto.AESEncrypt(key, plaintext)
decrypted, _ := crypto.AESDecrypt(key, encrypted)

// RSA 加密/解密
publicKey, privateKey, _ := crypto.GenerateRSAKey(2048)
encrypted, _ := crypto.RSAEncrypt(publicKey, data)
decrypted, _ := crypto.RSADecrypt(privateKey, encrypted)
```

#### 哈希工具 (`util/hash`)
```go
import "github.com/haysons/gokit/util/hash"

// MD5/SHA1/SHA256
md5 := hash.MD5(data)
sha256 := hash.SHA256(data)

// HMAC
hmac := hash.HMACSHA256(key, data)
```

#### 编码工具 (`util/encode`)
```go
import "github.com/haysons/gokit/util/encode"

// Base64
encoded := encode.Base64Encode(data)
decoded, _ := encode.Base64Decode(encoded)

// Hex
hex := encode.ToHex(data)
```

#### 唯一 ID 生成 (`util/uid`)
```go
import "github.com/haysons/gokit/util/uid"

// Snowflake ID
id, _ := uid.Snowflake()

// XID
xid := uid.XID()

// UUID
uuid := uid.UUID()

// Base58
base58 := uid.Base58()
```

#### 切片工具 (`util/slices`)
```go
import "github.com/haysons/gokit/util/slices"

// 去重
unique := slices.Unique([]int{1, 2, 2, 3})

// 交集
intersection := slices.Intersection(a, b)

// 差集
difference := slices.Difference(a, b)
```

#### Map 工具 (`util/maps`)
```go
import "github.com/haysons/gokit/util/maps"

// 键值反转
reversed := maps.Invert(originalMap)

// 合并
merged := maps.Merge(map1, map2)
```

### Middleware 中间件

#### 日志中间件
```go
import "github.com/haysons/gokit/middleware/logging"

mux := http.NewServeMux()
handler := logging.HTTPMiddleware(mux)
http.ListenAndServe(":8080", handler)
```

#### 恢复中间件
```go
import "github.com/haysons/gokit/middleware/recovery"

handler := recovery.HTTPMiddleware(next)
```

#### 认证中间件
```go
import "github.com/haysons/gokit/middleware/auth"

handler := auth.JWTMiddleware(secret, next)
```

### Errors 错误处理

```go
import "github.com/haysons/gokit/errors"

// 包装错误
err := errors.Wrap(originalErr, "database connection failed")

// 带堆栈的错误
err := errors.NewWithStack("something went wrong")

// 错误分类
if errors.IsNotFound(err) {
    // 处理不存在的情况
}
```

---

## 📁 项目结构

```
gokit/
├── app/              # 应用框架
├── config/           # 配置管理
├── constraints/      # 泛型约束
├── distributed/      # 分布式工具
├── errors/           # 错误处理
├── log/              # 日志组件
├── metadata/         # 元数据管理
├── middleware/       # 中间件
│   ├── auth/         # 认证中间件
│   ├── logging/      # 日志中间件
│   ├── metrics/      # 指标中间件
│   ├── recovery/     # 恢复中间件
│   └── tracing/      # 追踪中间件
├── transport/        # 传输层
│   ├── grpc/         # gRPC 传输
│   └── http/         # HTTP 传输
└── util/             # 工具函数
    ├── crypto/       # 加密工具
    ├── encode/       # 编码工具
    ├── hash/         # 哈希工具
    ├── maps/         # Map 工具
    ├── net/          # 网络工具
    ├── slices/       # 切片工具
    └── uid/          # ID 生成
```

---

## 🤝 贡献

欢迎贡献！请查看 [CONTRIBUTING.md](CONTRIBUTING.md) 了解如何参与。

1. Fork 本项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 提交 Pull Request

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 🙏 致谢

感谢以下项目的启发：
- [golang.org/x](https://golang.org/x) - Go 官方扩展包
- [spf13/viper](https://github.com/spf13/viper) - 配置管理
- [uber-go/zap](https://github.com/uber-go/zap) - 日志库

---

**Made with ❤️ by [@haysons](https://github.com/haysons)**
