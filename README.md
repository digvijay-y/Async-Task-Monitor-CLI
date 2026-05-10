# Async Task Monitor CLI

A real-time terminal dashboard for monitoring asynchronous tasks, workers, and execution pipelines built with Go and Bubble Tea v2.

Designed for developers who spend half their lives staring at terminals waiting for builds, deployments, ETL jobs, scripts, or CI pipelines to either succeed heroically or collapse in flames because someone forgot an environment variable.

---

# Preview

```text
┌──────────────────────────────────────────────────────────────────────┐
│ Async Task Monitor CLI                                              │
├──────────────────────────────────────────────────────────────────────┤
│ Workers: 4        Running: 3        Failed: 1        Queue: 12      │
├──────────────────────────────────────────────────────────────────────┤
│ ID   TASK NAME              STATUS       PROGRESS     DURATION       │
│ ──────────────────────────────────────────────────────────────────── │
│ 01   Build API Service      RUNNING      ██████░░ 65%  00:01:12      │
│ 02   Run Integration Tests  SUCCESS      ████████ 100% 00:02:44      │
│ 03   Deploy Production      FAILED       █████░░░ 48%  00:00:51      │
│ 04   Generate Reports       QUEUED       ░░░░░░░░ 0%   --            │
│ 05   Sync S3 Backups        RUNNING      ███░░░░░ 32%  00:00:29      │
├──────────────────────────────────────────────────────────────────────┤
│ Logs                                                                 │
│ [02] Tests completed successfully                                    │
│ [03] Docker image push failed                                        │
│ [01] Building authentication module                                  │
│ [05] Uploading backup chunks                                         │
├──────────────────────────────────────────────────────────────────────┤
│ q Quit   r Retry   c Cancel   enter View Logs   ↑↓ Navigate          │
└──────────────────────────────────────────────────────────────────────┘
```

Because apparently modern engineering culture decided every human deserves twelve dashboards, six terminals, and a mild stress disorder before lunch.

---

# Features

## Real-Time Task Monitoring

* Live task updates
* Dynamic progress tracking
* Task duration monitoring
* Execution statistics

## Concurrent Task Execution

* Goroutine-based async workers
* Configurable worker pools
* Queue management
* Parallel execution support

## Interactive Terminal UI

Built using:

* Bubble Tea
* Bubbles
* Lip Gloss

Features include:

* Keyboard navigation
* Status highlighting
* Progress bars
* Scrollable logs
* Responsive layouts

## Task Lifecycle Management

Supported states:

* QUEUED
* RUNNING
* SUCCESS
* FAILED
* CANCELLED
* RETRYING

## Log Streaming

* Live logs per task
* Error tracking
* Timestamped execution logs
* Expandable log viewer

## Failure Handling

* Retry failed tasks
* Automatic failure detection
* Worker crash recovery
* Graceful shutdown

---

# Architecture

```text
                    ┌────────────────────┐
                    │  Bubble Tea UI     │
                    └─────────┬──────────┘
                              │
                    ┌─────────▼──────────┐
                    │   Task Manager      │
                    └─────────┬──────────┘
                              │
         ┌────────────────────┼────────────────────┐
         │                    │                    │
┌────────▼────────┐ ┌────────▼────────┐ ┌────────▼────────┐
│ Worker Pool     │ │ Task Scheduler  │ │ Log Streamer    │
└────────┬────────┘ └────────┬────────┘ └────────┬────────┘
         │                   │                   │
         └───────────────────┼───────────────────┘
                             │
                   ┌─────────▼──────────┐
                   │ Async Task Engine   │
                   └─────────────────────┘
```

Tiny distributed systems simulator inside your terminal. Engineers cannot resist building orchestration layers. Give us enough coffee and we reinvent Kubernetes in 400 lines of Go.

---

# Project Structure

```text
async-task-monitor-cli/
│
├── cmd/
│   └── main.go
│
├── internal/
│   ├── manager/
│   │   └── manager.go
│   │
│   ├── worker/
│   │   └── worker.go
│   │
│   ├── model/
│   │   └── task.go
│   │
│   ├── scheduler/
│   │   └── scheduler.go
│   │
│   ├── logger/
│   │   └── logger.go
│   │
│   └── ui/
│       ├── tui.go
│       ├── table.go
│       ├── logs.go
│       └── styles.go
│
├── configs/
│   └── config.yaml
│
├── screenshots/
│
├── go.mod
├── go.sum
├── README.md
└── LICENSE
```

---

# Installation

## Clone Repository

```bash
git clone https://github.com/yourusername/async-task-monitor-cli.git

cd async-task-monitor-cli
```

## Install Dependencies

```bash
go mod tidy
```

## Run Application

```bash
go run ./cmd/main.go
```

## Build Binary

```bash
go build -o atm-cli ./cmd/main.go
```

---

# Keyboard Controls

| Key     | Action                |
| ------- | --------------------- |
| `q`     | Quit application      |
| `↑ / ↓` | Navigate tasks        |
| `enter` | Open detailed logs    |
| `r`     | Retry failed task     |
| `c`     | Cancel running task   |
| `d`     | Delete completed task |
| `tab`   | Switch panels         |
| `h`     | Help menu             |

Because no terminal app is complete until someone presses random keys in panic trying to exit Vim’s emotionally unstable cousin.

---

# Example Use Cases

## CI/CD Monitoring

Track:

* builds
* tests
* deployments
* release pipelines

## DevOps Operations

Monitor:

* Docker jobs
* Kubernetes tasks
* backup jobs
* automation scripts

## Data Pipelines

Track:

* ETL jobs
* ingestion tasks
* ML training pipelines
* report generation

## Background Processing

Monitor:

* email queues
* cron jobs
* distributed workers
* async services

---

# Planned Features

## Phase 2

* SQLite persistence
* Search & filtering
* Multi-tab dashboards
* Metrics aggregation

## Phase 3

* Docker integration
* Kubernetes job monitoring
* SSH remote workers
* WebSocket streaming

## Phase 4

* Prometheus metrics
* Grafana integration
* Distributed task execution
* Cluster-aware scheduling

At some point every useful tool slowly evolves into infrastructure software. Civilization itself is basically layers of automation duct-taped together by sleep-deprived engineers pretending YAML is readable.

---

# Tech Stack

* Language: Go
* TUI Framework: Bubble Tea
* UI Components: Bubbles
* Styling: Lip Gloss

---

# Learning Outcomes

This project demonstrates:

* Concurrent programming in Go
* Goroutines & channels
* Event-driven architecture
* Terminal UI engineering
* Worker pool implementation
* State management
* Async task orchestration
* Real-time log streaming

---

# Future Vision

Long-term, this project can evolve into:

* lightweight job orchestrator
* internal platform engineering tool
* terminal-based CI monitor
* distributed task execution framework

Basically the natural progression from “simple CLI project” to “accidentally building production infrastructure at 2 AM.” A timeless engineering tradition.
