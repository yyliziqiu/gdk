# GDK — Go Development Kit

> A modular utility library for Go server-side development, providing standardized wrappers for databases, caching, logging, HTTP, messaging, and more.

[![Go Version](https://img.shields.io/badge/go-1.24%2B-00ADD8?logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

```bash
go get github.com/yyliziqiu/gdk
```

---

## Packages

| Package           | Description                                                |
|-------------------|------------------------------------------------------------|
| [xboot](#xboot)   | Application lifecycle management                           |
| [xconv](#xconv)   | Primitive type conversions                                 |
| [xcq](#xcq)       | Circular queue with persistence                            |
| [xcsv](#xcsv)     | CSV file generation                                        |
| [xdb](#xdb)       | Database connection management (MySQL / PostgreSQL)        |
| [xes](#xes)       | Elasticsearch client management                            |
| [xerr](#xerr)     | HTTP-aware custom errors                                   |
| [xfile](#xfile)   | File system utilities                                      |
| [xgin](#xgin)     | Gin framework integration                                  |
| [xhttp](#xhttp)   | HTTP client with logging                                   |
| [xif](#xif)       | Generic conditional helpers                                |
| [xkafka](#xkafka) | Kafka producer / consumer management                       |
| [xkvs](#xkvs)     | Typed key-value string map                                 |
| [xlog](#xlog)     | Structured logging with file rotation                      |
| [xredis](#xredis) | Redis client management                                    |
| [xsnap](#xsnap)   | JSON snapshot persistence                                  |
| [xstr](#xstr)     | String utilities                                           |
| [xtask](#xtask)   | Cron task scheduling                                       |
| [xtime](#xtime)   | Time utilities                                             |
| [xtmpl](#xtmpl)   | HTML template management                                   |
| [xuid](#xuid)     | Distributed unique ID generation                           |
| [xutil](#xutil)   | Compression, signing, reflection, load balancing, and more |

---

## Package Details

### xboot

Application lifecycle framework that unifies initialization, startup, and graceful shutdown.

- Register hooks via `InitFunc` / `BootFunc`
- Auto-injects MySQL/PostgreSQL, Redis, Elasticsearch, and Kafka components
- Integrates Viper for config loading; supports `ConfigChecker` interface for validation

```go
app := xboot.New(conf)
app.SetInitFuncs(myInit)
app.SetBootFuncs(myBoot)
app.Run()
```

---

### xconv

Concise conversion functions between common primitive types.

| Function                          | Description                           |
|-----------------------------------|---------------------------------------|
| `S2B` / `S2I` / `S2I64` / `S2F64` | String → bool / int / int64 / float64 |
| `S2T`                             | String → Unix timestamp               |
| `B2S` / `I2S` / `I642S` / `F642S` | Reverse conversions                   |
| `T2S`                             | Unix timestamp → string               |

---

### xcq

Circular queue with dynamic resizing, optional persistence, and a thread-safe wrapper.

- `Push` / `Pop` — enqueue / dequeue; `Pops` — batch pop by predicate
- `Get` / `HeadItem` / `TailItem` — random access
- `New2(path)` — persist queue to disk; content survives process restarts
- `SyncQueue` — mutex-protected wrapper for concurrent use

---

### xcsv

Serialize Go struct slices to CSV files using reflection.

- `Save` — reads struct fields via reflection (supports custom tags) and writes CSV
- `SaveRows` — write raw `[][]string` row data directly

---

### xdb

Connection pool management for MySQL and PostgreSQL, supporting both `database/sql` and GORM.

- `Config` — DSN, pool settings (max open/idle connections, lifetime), ORM options
- `Init` — initialize multiple named connections at once
- `Get` / `GetOrm` — retrieve a connection by ID
- `AddMigration` / `Migrate` — built-in migration support

---

### xes

Lifecycle management for Elasticsearch clients.

- `Config` — hosts, credentials, log level
- `Init` — initialize multiple named clients
- `Get` / `GetDefault` — retrieve a client by ID
- `Finally` — release all client resources

---

### xerr

Custom error type with HTTP semantics for consistent error handling in web services.

- Pre-defined errors: `BadRequest`, `Unauthorized`, `Forbidden`, `NotFound`, `InternalServerError`, …
- `New` / `New2` — create errors with custom code and message
- `Wrap` / `Format` / `With` — chain contextual information onto an error
- `GetStatus` / `Http` — map directly to an HTTP status code and response body
- Code convention: `A`-prefix (client errors), `B`-prefix (server errors), `C`-prefix (third-party errors)

---

### xfile

Lightweight file system helpers.

- `Exist` — check whether a file or directory exists
- `MakeDir` — recursively create a directory (idempotent, equivalent to `mkdir -p`)

---

### xgin

Gin framework integration, split into three sub-packages.

- `Config` — listen address and TLS settings
- `Run` — start the Gin server, automatically wiring xlog for error and access logs
- **`xgin/xmid`** — middleware collection (CORS, etc.)
- **`xgin/xresp`** — helpers for uniform JSON and error responses
- **`xgin/xreq`** — request parameter binding utilities

---

### xhttp

Full-featured HTTP client with built-in request/response logging.

- `New` — functional options: timeout, logger, cookie jar, redirect policy, …
- Supports `Get` / `Post` / `Put` / `Patch` / `Delete` / `Head`
- Auto-detects `Content-Type` and deserializes response bodies
- Configurable body truncation length for log output

---

### xif

Generic conditional helpers that eliminate repetitive `if-else` assignments.

- `If[T](cond, a, b)` — ternary selection
- `Nil` — return the first non-nil value
- `Zero[T]` — return the first non-zero numeric value
- `Empty[T]` — return the first non-empty string / slice / map

---

### xkafka

Lifecycle management for Kafka producers and consumers.

- `Config` — broker addresses, SASL/SSL authentication
- `ProducerConfig` / `ConsumerConfig` — independent producer and consumer settings
- `Init` — batch initialization
- `Produce` / `ProduceModel` — send messages (with optional struct serialization)
- `Push` / `PushModel` — shortcuts using the default producer

---

### xkvs

A type-enriched wrapper around `map[string]string` with automatic value conversion.

| Method style                             | Description                           |
|------------------------------------------|---------------------------------------|
| `String` / `Bool` / `Int` / …            | Returns value + `ok` flag             |
| `S` / `B` / `I` / …                      | Returns value with a fallback default |
| `S2` / `B2` / `I2` / …                   | Case-insensitive key lookup           |
| `Id` / `Type` / `Name` / `Token` / `Url` | Semantic shortcut accessors           |

---

### xlog

Structured logging built on Logrus, with time-based file rotation.

- `Config` — level, format, file path, rotation interval
- `Init` — initialize the global logger
- `New` / `New2` / `New3` — create independent loggers (console or file-backed)
- Built-in globals: `Default` (file) and `Console`

---

### xredis

Redis client management supporting standalone, sentinel, and cluster modes.

- `Config` — mode, addresses, connection pool parameters
- `Init` — initialize multiple named clients
- `Get` / `GetDefault` — retrieve a client by ID
- Mode constants: `ModeStandard` / `ModeCluster` / `ModeSentinel` / `ModeSentinelCluster`

---

### xsnap

Persist arbitrary Go data structures as JSON snapshots, with deduplication and backup support.

- `New` / `New2` / `New3` — create a snapshot with optional deduplication TTL
- `Save` / `Load` — write and read snapshot data
- `Dup` — create a timestamped backup copy
- Implement the `Persistent` interface to integrate custom types

---

### xstr

String manipulation utilities.

- `Truncate` / `TruncateUtf8` — truncate strings (UTF-8 aware)
- `TrimSplit` — split a string and trim whitespace from each part
- `RandomDigit` / `RandomAlphabet` / `RandomString` — generate random digit / alpha / mixed strings
- `Random` — random string from a custom character set

---

### xtask

Second-precision cron scheduling powered by `robfig/cron`.

- `CronTask` — define a task with a name, cron expression, and handler function
- `RunCronTasks` — start the scheduler with a list of tasks
- `RunCronTasksWithConfig` — apply external configuration overrides
- `Once` — run a task exactly once

---

### xtime

Time utility functions for formatting, boundaries, and calendar calculations.

- `DateTime` / `DateTimeIn` — current time as a formatted string
- `Timestamp` — current Unix timestamp
- `ManualDuration` — format a `time.Duration` in human-readable form
- Boundary helpers: `DayBegin` / `DayEnd` / `WeekBegin` / `WeekEnd` / `MonthBegin` / `MonthEnd`, …
- `IsLeap` — leap year check; `DaysOfMonth` — days in a given month

---

### xtmpl

HTML template loading, caching, and rendering with hot-reload support.

- `Manager` — load templates from explicit file paths or glob patterns
- `Html` / `HtmlGin` — render to `http.ResponseWriter` or a Gin `Context`
- `SetDebug` — reload templates on every render (development mode)
- `Reload` — manually reload all templates

---

### xuid

Timestamp-based distributed unique ID generator.

- `Generator` — carries a node ID and sequence counter; produces 16-character hex IDs
- `New2(path)` — persist sequence state so IDs remain unique across restarts
- `Get` — generate an ID; `GetOrFail` — returns an error on clock rollback

---

### xutil

A collection of independent utilities across several domains.

**Compression**

- `GzipEncode` / `GzipDecode` — Gzip compress / decompress
- `ZlibEncode` / `ZlibDecode` — Zlib compress / decompress

**Signatures & Tokens**

- `GenerateSignature` / `VerifySignature` — SHA-256 based HMAC signing
- `GenerateTimeSignature` / `VerifyTimeSignature` — time-bounded tokens with TTL

**Reflection**

- `ReflectFuncName` — get a function's name at runtime
- `ReflectFields` / `ReflectValues` — extract struct field names and values
- `ReflectValueByField` — look up a field value by name

**Load Balancing**

- `Swrr[T]` — smooth weighted round-robin (SWRR)
- `Round` — simple round-robin for equal-weight nodes

**Miscellaneous**

- `Trie` — prefix tree for efficient string prefix matching
- `Percent` / `Round` — numeric helpers
- `Mime` — MIME type utilities

---

## License

[MIT](LICENSE)
