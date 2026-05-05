# go-worker-pool

A simple Go project demonstrating **goroutines** and **channels** using the **worker pool** concurrency pattern.

## Overview

A worker pool keeps a fixed number of goroutines alive and feeds them work through a shared channel.
This avoids the overhead of spawning a new goroutine per task while naturally limiting the number of concurrent operations.

```
Submitter ──► jobs channel ──► [worker 1]
                            ──► [worker 2]  ──► results channel ──► Collector
                            ──► [worker 3]
```

## Features

- ✅ Reusable worker pool implementation
- ✅ Configurable number of concurrent workers
- ✅ Error handling and propagation
- ✅ Graceful shutdown mechanism
- ✅ Task-based architecture (implement `Task` interface)
- ✅ Result collection in any order

## Project Structure

```
go-worker-pool/
├── main.go              # Entry point – squares numbers 1-10 concurrently
├── go.mod               # Go module file
├── README.md            # This file
├── worker/
│   └── pool.go          # Worker pool implementation
├── task/
│   └── task.go          # Task interface and Result type
└── example/
    └── example.go       # Additional example: division with error handling
```

## Key Concepts

| Concept | Where used |
|---------|-----------|
| Goroutines | `worker/pool.go` – each worker runs as a goroutine |
| Channels | `jobs` and `results` channels carry work and outcomes |
| `sync.WaitGroup` | Tracks when all workers have finished |
| Graceful shutdown | `Pool.Stop()` closes the jobs channel and waits for workers |
| Error propagation | `task.Result.Err` carries per-task errors to the collector |

## Requirements

- Go 1.16 or higher

## Quick Start

```bash
git clone https://github.com/hasanulkabir-md/go-worker-pool.git
cd go-worker-pool
go run main.go
```

Expected output (task order may vary):

```
Starting worker pool with 3 workers to process 10 tasks...

Results:
  task 1: 1² = 1
  task 2: 2² = 4
  task 3: 3² = 9
  ...
  task 10: 10² = 100

All tasks completed.
```

## Running the Additional Example

```bash
go run example/example.go
```

## Usage

### 1. Define a Task

Implement the `task.Task` interface:

```go
type myTask struct{ value int }

func (t myTask) Process() (interface{}, error) {
    return t.value * 2, nil
}
```

### 2. Create and Start a Pool

```go
pool := worker.NewPool(5) // 5 concurrent workers
pool.Start()
```

### 3. Submit Work

```go
go func() {
    for i := 0; i < 20; i++ {
        pool.Submit(i, myTask{value: i})
    }
    pool.Stop() // no more tasks; waits for workers to finish
}()
```

### 4. Collect Results

```go
for result := range pool.Results() {
    if result.Err != nil {
        log.Printf("task %d failed: %v", result.TaskID, result.Err)
        continue
    }
    fmt.Printf("task %d → %v\n", result.TaskID, result.Value)
}
```

## When to Use

This pattern is ideal for:
- **API scraping** – limit concurrent requests
- **Data processing** – batch file operations
- **Background jobs** – handle queued tasks
- **Resource-constrained systems** – control resource usage

## License

MIT
