# Step 2: Basic HTTP Server with Graceful Shutdown

**Branch:** [`step2-basic-server`](https://github.com/lentregu/gh-autoresponder/tree/step2-basic-server)  
**Created Files:**
- `cmd/server/main.go`
- `internal/server/server.go`
- `step2-basic-server.md`

## 🧭 Tutorial Navigation

| Previous | Next |
|----------|------|
| [Step 1: Setup and Readme →](step1-readme.md) | [Step 3: Webhook Endpoint →](step3-webhook-endpoint.md) |



## 🏗 Implementation Overview

We've implemented a production-ready HTTP server with:
- Graceful shutdown capability
- Signal handling (SIGTERM, Interrupt)
- Error channel patterns
- Context timeout management

## 📂 Repository Structure

```text
gh-autoresponder/
├── .gitignore          # Standard Go ignore patterns
├── README.md           # Project overview and tutorial index
├── go.mod              # Go module definition
│
├── cmd/
│   └── server/
│       └── main.go     # CLI entry point (minimal logic)
│
├── internal/
│   └── server/
│       └── server.go   # Core server implementation
│
└── doc/
    ├── step1-readme.md # Initial setup documentation
    └── step2-basic-server.md # Current server implementation guide
```

## 🧠 Core Components

### 1. Server Initialization (`internal/server/server.go`)
```go
type Server struct {
    httpServer *http.Server  // Standard HTTP server
    shutdown   chan os.Signal // Graceful shutdown channel
}

func New(port string) *Server {
    return &Server{
        httpServer: &http.Server{Addr: ":" + port},
        shutdown:   make(chan os.Signal, 1),
    }
}
```

### 2. Server Lifecycle Management

```go
func (s *Server) Start() error {
    signal.Notify(s.shutdown, os.Interrupt, syscall.SIGTERM)
    serverErr := make(chan error, 1)
    
    go func() {
        if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            serverErr <- err
        }
    }()
    
    select {
    case err := <-serverErr: return err
    case <-s.shutdown: return s.Stop()
    }
}

func (s *Server) Stop() error {
    ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
    defer cancel()
    return s.httpServer.Shutdown(ctx)
}
```

## 🚀 Startup Sequence (`cmd/server/server.go`)

```go
srv := server.New("8080")
srv.Start()  // Blocks until shutdown
```

## 🛠 Key Features

### Core Functionality Table
| Feature               | Implementation                                | Benefit                                  |
|-----------------------|-----------------------------------------------|------------------------------------------|
| Graceful Shutdown     | `signal.Notify()` + `shutdown` channel        | Clean exit on SIGTERM/CTRL+C             |
| Error Isolation       | Dedicated `serverErr` channel                 | Prevents error swallowing                |
| Context Timeout       | `context.WithTimeout(15*time.Second)`         | Enforces maximum shutdown duration       |
| Concurrent Operation  | `go s.httpServer.ListenAndServe()`            | Non-blocking server start                |

### Signal Handling Details
| Signal            | Trigger Method              | Handler Behavior                     |
|-------------------|-----------------------------|--------------------------------------|
| `os.Interrupt`    | CTRL+C in terminal          | Initiates graceful shutdown          |
| `syscall.SIGTERM` | `kill -15 <PID>`            | Same as interrupt                    |

### Method Comparison Table
| Method    | Parameters      | Return Value         | Thread Safety     |
|-----------|-----------------|----------------------|-------------------|
| `New()`   | `port string`   | `*Server`            | ✅ Yes            |
| `Start()` | None            | `error`              | ❌ No (blocking)  |
| `Stop()`  | None            | `error`              | ✅ Yes            |

### Timeout Configuration
| Scenario          | Recommended Timeout | Rationale                            |
|-------------------|---------------------|--------------------------------------|
| Local Development | 15 seconds          | Balance between fast fail and safety |
| Production        | 30 seconds          | Account for higher load              |