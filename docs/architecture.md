# Architecture

ClaudeFlux is built as a modular Go application with clear separation of concerns.

## Components

### Runtime (`internal/runtime/`)
The core orchestration loop. Receives a parsed workflow, builds the DAG, and executes agents layer by layer. Each layer runs agents in parallel using goroutines.

### DAG Engine (`internal/dag/`)
Parses agent dependency graphs and produces topological layers for execution ordering. Validates for cycles and missing dependencies.

### Worker Pool (`internal/runtime/worker_pool.go`)
Manages concurrent worker processes. Each worker gets an isolated git worktree and spawns a Claude CLI process.

### Worker (`internal/worker/`)
Individual agent execution. Handles:
- Git worktree creation and cleanup
- Claude CLI process spawning
- Token counting and budget enforcement
- Output parsing and validation

### Coordinator (`internal/coordinator/`)
Routes messages between agents and manages approval gates. When an agent reaches an approval gate, execution pauses until a human approves or rejects.

### Steward (`internal/steward/`)
Post-workflow automation: reviewing agent output, merging git branches, and creating pull requests.

### Budget (`internal/budget/`)
Thread-safe token and cost tracking. Supports per-agent and workflow-level limits with configurable actions on exceeded (checkpoint, pause, abort).

### Store (`internal/store/`)
SQLite-backed audit log. Records every event: agent starts, completions, failures, budget alerts, and approval decisions.

### IPC (`internal/ipc/`)
Inter-process communication layer. Provides:
- Unix socket for local agent-to-agent messaging
- SSE (Server-Sent Events) for dashboard real-time updates

### Config (`internal/config/`)
YAML workflow parser with full validation against the JSON Schema.

## Data Flow

```
workflow.yaml → Config Parser → DAG Builder → Runtime
                                                  ↓
                                        ┌── Worker Pool ──┐
                                        │  Worker 1       │
                                        │  Worker 2       │
                                        │  Worker N       │
                                        └─────────────────┘
                                                  ↓
                                        Coordinator → Steward
                                                  ↓
                                            SQLite Audit Log
```
