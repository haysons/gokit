# gokit

[中文](README_CN.md)

[![Go Report Card](https://goreportcard.com/badge/github.com/haysons/gokit)](https://goreportcard.com/report/github.com/haysons/gokit)
[![MIT License](https://img.shields.io/badge/license-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT)
[![Go Reference](https://pkg.go.dev/badge/github.com/haysons/gokit.svg)](https://pkg.go.dev/github.com/haysons/gokit)

**gokit** is a Go language toolkit that provides rich helper functions and components to simplify development processes and improve code efficiency.

## Features

- **Modular design** - Use as needed, no intrusion
- **Practical tools** - Covers common development scenarios
- **Type safety** - Fully utilizes the Go type system
- **Production ready** - Well tested and verified

## Installation

```bash
go get github.com/haysons/gokit
```

## Module Overview

| Module | Description |
|--------|-------------|
| config | Configuration management based on Viper, supporting hot reload |
| log | Structured logging based on slog |
| errors | Enhanced error handling |
| middleware | HTTP/gRPC middleware collection |
| transport | HTTP/gRPC transport layer encapsulation |
| util | Generic utility functions |
| app | Application lifecycle management |
| distributed | Distributed system tools |
| metadata | Context metadata management |
| constraints | Generic programming constraints |

---

## Usage Examples

### Config Configuration Management

Configuration management based on `viper`, supporting multiple configuration formats, environment variables, and hot reload:

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
    
    if err := cfg.Load(); err != nil {
        panic(err)
    }
    
    conf := cfg.Get()
    println(conf.Server.Port)
}
```

### Log Component

Structured logging based on `slog`, supporting JSON/text format and file rotation:

```go
package main

import (
    "github.com/haysons/gokit/log"
)

func main() {
    logger := log.GetDefaultSlog()
    
    logger.Info("server started", "port", 8080)
    logger.Error("connection failed", "error", err)
}
```

### Util Functions

#### Encryption Tools (`util/crypto`)
```go
import "github.com/haysons/gokit/util/crypto"

encrypted, _ := crypto.AESEncrypt(key, plaintext)
decrypted, _ := crypto.AESDecrypt(key, encrypted)
```

#### Unique ID Generation (`util/uid`)
```go
import "github.com/haysons/gokit/util/uid"

id, _ := uid.Snowflake()
xid := uid.XID()
uuid := uid.UUID()
```

### Middleware

#### Logging Middleware
```go
import "github.com/haysons/gokit/middleware/logging"

mux := http.NewServeMux()
handler := logging.HTTPMiddleware(mux)
http.ListenAndServe(":8080", handler)
```

### Errors Handling

```go
import "github.com/haysons/gokit/errors"

err := errors.Wrap(originalErr, "database connection failed")
if errors.IsNotFound(err) {
    // Handle not found case
}
```

---

## Project Structure

```
gokit/
├── app/              # Application framework
├── config/           # Configuration management
├── constraints/      # Generic constraints
├── distributed/      # Distributed tools
├── errors/           # Error handling
├── log/              # Logging component
├── metadata/         # Metadata management
├── middleware/       # Middleware
│   ├── auth/         # Authentication middleware
│   ├── logging/      # Logging middleware
│   ├── metrics/      # Metrics middleware
│   ├── recovery/     # Recovery middleware
│   └── tracing/      # Tracing middleware
├── transport/        # Transport layer
│   ├── grpc/         # gRPC transport
│   └── http/         # HTTP transport
└── util/             # Utility functions
    ├── crypto/       # Encryption tools
    ├── encode/       # Encoding tools
    ├── hash/         # Hash tools
    ├── maps/         # Map utilities
    ├── net/          # Network tools
    ├── slices/       # Slice utilities
    └── uid/          # ID generation
```

---

## Contribution

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for details.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

Thanks to the following projects for inspiration:
- [golang.org/x](https://golang.org/x) - Go official extension packages
- [spf13/viper](https://github.com/spf13/viper) - Configuration management
- [uber-go/zap](https://github.com/uber-go/zap) - Logging library

---

**Made with ❤️ by [@haysons](https://github.com/haysons)**
