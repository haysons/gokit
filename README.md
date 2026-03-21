# gokit

[中文](README_CN.md)

[![Go Report Card](https://goreportcard.com/badge/github.com/haysons/gokit)](https://goreportcard.com/report/github.com/haysons/gokit)
[![MIT License](https://img.shields.io/badge/license-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT)
[![Go Reference](https://pkg.go.dev/badge/github.com/haysons/gokit.svg)](https://pkg.go.dev/github.com/haysons/gokit)

**gokit** is a Go language toolkit for production-ready services. Modular, type-safe, and thoroughly tested.

## Features

- **Modular** - Import only what you need
- **Type-safe** - Full Go generics support (Go 1.23+)
- **Production-ready** - battle-tested in production environments
- **DI-first** - Built-in dependency injection with uber/fx

## Installation

```bash
go get github.com/haysons/gokit
```

## Quick Start

```go
// Config with hot reload
cfg := config.New[AppConfig]()
cfg.SetFile("./config.yaml")
cfg.Load()

// Structured logging
logger := log.GetDefaultSlog()
logger.Info("server started", "port", 8080)

// Enhanced errors
err := errors.Wrap(dbErr, "database connection failed")
if errors.IsNotFound(err) { /* handle */ }
```

## Modules

| Module | Description |
|--------|-------------|
| `app` | Application lifecycle with uber/fx |
| `config` | Viper-based config with hot reload |
| `log` | slog-based structured logging |
| `errors` | Business codes + stack trace + hints |
| `middleware` | HTTP/gRPC middleware (auth, logging, tracing) |
| `transport` | Unified HTTP/gRPC transport layer |
| `distributed` | etcd-based: lock, election, queue, counter |
| `metadata` | Context metadata for RPC |
| `util` | crypto, uid, hash, slices, maps... |

## Project Structure

```
gokit/
├── app/            # Application framework
├── config/         # Configuration management
├── distributed/    # Distributed tools (etcd)
├── errors/         # Error handling
├── log/            # Structured logging
├── middleware/    # HTTP/gRPC middleware
├── transport/     # Transport layer
├── metadata/      # Context metadata
├── constraints/   # Generics
└── util/          # Utility functions
```

## License

MIT License - see [LICENSE](LICENSE) for details.

---

**Made with ❤️ by [@haysons](https://github.com/haysons)**
