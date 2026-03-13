<div align="center">

# ClaudeFlux

### Production-grade multi-agent orchestration runtime for Claude.
### DAG scheduling · git isolation · token budgets · approval gates · live dashboard.

[![CI](https://github.com/Subodh8/ClaudeFlux/actions/workflows/ci.yml/badge.svg)](https://github.com/Subodh8/ClaudeFlux/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/Subodh8/ClaudeFlux)](https://goreportcard.com/report/github.com/Subodh8/ClaudeFlux)
[![License: Apache 2.0](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Stars](https://img.shields.io/github/stars/Subodh8/ClaudeFlux?style=social)](https://github.com/Subodh8/ClaudeFlux)

[Quickstart](#quickstart) · [Documentation](docs/) · [Examples](examples/) · [Discord](#)

</div>

---

> "We rewrote our 2,000-line custom orchestrator in a 40-line YAML workflow.
> ClaudeFlux handles the rest."

---

## The Problem

Running multiple Claude agents in parallel is surprisingly hard:

| Problem | What happens without ClaudeFlux |
|---------|----------------------------------|
| Shared filesystem | Agent A overwrites Agent B's files mid-task |
| No cost control | A runaway agent burns $200 in one session |
| No approval gates | Agents auto-merge broken code to main |
| Context exhaustion | Orchestrator hits token limit, silent failure |
| No observability | You have no idea what 6 agents are doing |
| Duplicated work | Two agents solve the same subproblem |
| Custom IPC | Every team rebuilds inter-agent messaging |

ClaudeFlux is the runtime that fixes all of this — in a single binary.

## Architecture Diagram

```
┌────────────────────────────────────────────────────────────────┐
│                        ClaudeFlux Runtime                      │
│                                                                │
│   ┌───────────┐    ┌──────────────┐    ┌───────────────────┐   │
│   │ YAML DAG  │───▶│  DAG Parser  │───▶│  Dispatch Daemon │   │
│   │  Workflow │    │  + Validator │    │  (dependency res) │   │
│   └───────────┘    └──────────────┘    └────────┬──────────┘   │
│                                                 │              │
│                    ┌────────────────────────────▼──────────┐   │
│                    │           Worker Pool Manager         │   │
│   ┌────────────────▼──────────────────────────────────┐    │   │
│   │                  Worker Process (per agent)       │    │   │
│   │  ┌─────────────┐  ┌──────────┐  ┌─────────────┐   │    │   │
│   │  │ git worktree│  │  Claude  │  │   Token $   │   │    │   │
│   │  │  isolation  │  │  Code    │  │   Budget    │   │    │   │
│   │  │  (isolated) │  │  (proc)  │  │  Enforcer   │   │    │   │
│   │  └─────────────┘  └──────────┘  └─────────────┘   │    │   │
│   └───────────────────────────────────────────────────┘    │   │
│                                                            │   │
│   ┌──────────────┐    ┌──────────────┐    ┌─────────────┐  │   │
│   │ Coordinator  │    │   Steward    │    │  Rate Limit │  │   │
│   │  (routing +  │    │  (reviewer + │    │ Coordinator │  │   │
│   │   approval)  │    │  auto-merge) │    │             │  │   │
│   └──────────────┘    └──────────────┘    └─────────────┘  │   │
│                                                            │   │
│   ┌──────────────────────────────────────────────────────┐ │   │
│   │              SQLite Audit Log (every event)          │ │   │
│   └──────────────────────────────────────────────────────┘ │   │
└────────────────────────────────────────────────────────────────┘
              │                               │
              ▼                               ▼
  ┌──────────────────────┐       ┌────────────────────────┐
  │   Unix Socket / SSE  │       │   Next.js Dashboard    │
  │   (IPC layer)        │       │   DAG · Logs · Costs   │
  └──────────────────────┘       │   Approval Queue       │
                                 └────────────────────────┘
```

## Quickstart

### Option 1: Single binary (recommended)

```bash
# macOS / Linux
curl -sSL https://claudeflux.dev/install.sh | sh

# Verify
claudeflux version
```

### Option 2: Homebrew

```bash
brew tap Subodh8/ClaudeFlux
brew install claudeflux
```

### Option 3: Docker

```bash
docker run --rm -v $(pwd):/workspace ghcr.io/Subodh8/ClaudeFlux:latest run workflow.yaml
```

### Option 4: Build from source

```bash
git clone https://github.com/Subodh8/ClaudeFlux
cd ClaudeFlux
make build
./bin/claudeflux version
```

> **Prerequisite:** claude CLI must be installed and authenticated. See https://docs.anthropic.com/claude-code

## Your First Workflow

Create `workflow.yaml`:

```yaml
version: "1"
name: research-write-review

budget:
  max_tokens: 100000
  max_cost_usd: 2.00
  alert_at_percent: 80

agents:
  researcher:
    prompt: |
      Research the top 5 open-source vector databases in 2025.
      For each: name, license, performance benchmarks, best use case.
      Output as structured JSON to researcher_output.json.
    budget:
      max_tokens: 30000

  writer:
    depends_on: [researcher]
    prompt: |
      Read researcher_output.json. Write a 1500-word technical blog post
      comparing the databases. Save to post.md. Audience: senior engineers.
    budget:
      max_tokens: 40000

  reviewer:
    depends_on: [writer]
    prompt: |
      Review post.md for technical accuracy, readability, and completeness.
      Output a JSON review to review.json with fields:
      score (0-10), issues (array), approved (bool).
    budget:
      max_tokens: 20000

  steward:
    depends_on: [reviewer]
    type: steward
    rules:
      require_reviewer_approval: true
      min_review_score: 7
      on_approved: commit_and_pr
      on_rejected: notify_slack
```

Run it:

```bash
claudeflux run workflow.yaml --dashboard
```

Open `http://localhost:7070` to watch the DAG execute in real time.

## CLI Reference

```bash
# Run a workflow
claudeflux run workflow.yaml

# Run with live dashboard
claudeflux run workflow.yaml --dashboard

# Dry run (validate + cost estimate, no execution)
claudeflux run workflow.yaml --dry-run

# Resume a failed workflow from checkpoint
claudeflux resume <run-id>

# View audit log for a run
claudeflux logs <run-id>

# List all runs
claudeflux runs list

# Approve a pending gate
claudeflux approve <gate-id>

# Reject a pending gate
claudeflux reject <gate-id> --reason "Output quality too low"

# Real-time agent status
claudeflux status <run-id>

# Validate workflow YAML without running
claudeflux validate workflow.yaml

# Cost estimate (calls Claude API for token estimation)
claudeflux estimate workflow.yaml
```

## Core Features

### Git Worktree Isolation

Every agent gets its own git worktree — a fully isolated working directory backed by the same repository. Agents can read/write freely without colliding. The steward agent merges cleanly when the workflow completes.

```bash
# ClaudeFlux creates these automatically:
.git/worktrees/cf-researcher-a3f2/
.git/worktrees/cf-writer-b1e9/
.git/worktrees/cf-reviewer-c7d4/
```

### Token Budget Enforcement

Define per-agent and workflow-level token and cost limits. Agents that exceed their budget are gracefully stopped with their partial output preserved.

```yaml
budget:
  max_tokens: 50000        # Hard stop at 50k tokens
  max_cost_usd: 1.00       # Hard stop at $1.00
  alert_at_percent: 75     # Slack/webhook alert at 75%
  on_exceeded: checkpoint  # Options: checkpoint | pause | abort
```

### Approval Gates

Insert human-in-the-loop checkpoints anywhere in the DAG. ClaudeFlux pauses and surfaces the approval request in the dashboard (and optionally Slack).

```yaml
agents:
  deploy:
    depends_on: [test]
    approval:
      required: true
      timeout: 3600          # Auto-reject after 1 hour
      notify:
        - slack: "#deployments"
        - email: "team@yourcompany.com"
```

### Typed Inter-Agent Messaging

Agents communicate via typed JSON contracts, validated against schemas you define. No more parsing freeform text between agents.

```yaml
contracts:
  researcher_output:
    type: object
    required: [databases]
    properties:
      databases:
        type: array
        items:
          type: object
          required: [name, license, benchmark_qps]
```

### Critic Agents (Adversarial Validation)

Add a `type: critic` agent that automatically challenges another agent's output. Built-in adversarial review without custom prompting.

```yaml
  skeptic:
    type: critic
    targets: [writer]
    prompt: |
      Find factual errors, unsupported claims, and missing context
      in the draft. Be harsh. Output structured critique.
```

### Live Dashboard

A Next.js dashboard with React Flow DAG visualization, SSE-powered log streaming, cost heatmap, and an approval queue UI.

```bash
claudeflux run workflow.yaml --dashboard --dashboard-port 7070
```

## Comparison Table: ClaudeFlux vs Alternatives

| Feature | ClaudeFlux | LangGraph | CrewAI | Airflow |
|---------|------------|-----------|--------|---------|
| Claude-native | ✅ | ⚠️ | ⚠️ | ❌ |
| Git worktree isolation | ✅ | ❌ | ❌ | ❌ |
| Token budget enforcement | ✅ | ❌ | ❌ | ❌ |
| Approval gates | ✅ | ⚠️ | ❌ | ✅ |
| Typed inter-agent messages | ✅ | ⚠️ | ❌ | ❌ |
| Critic / adversarial agents | ✅ | ❌ | ⚠️ | ❌ |
| Single binary deployment | ✅ | ❌ | ❌ | ❌ |
| SQLite audit log | ✅ | ❌ | ❌ | ✅ |
| Live DAG dashboard | ✅ | ⚠️ | ❌ | ✅ |
| YAML declarative config | ✅ | ❌ | ❌ | ✅ |
| MCP-server aware | ✅ | ❌ | ❌ | ❌ |

> `⚠️ = partial / requires custom implementation`

## Roadmap

- [ ] MCP server registry integration
- [ ] Conditional DAG branches (if/else)
- [ ] Slack + PagerDuty approval notifications
- [ ] `claudeflux init` interactive scaffolder
- [ ] Remote worker mode (gRPC)
- [ ] PostgreSQL backend
- [ ] OpenTelemetry traces
- [ ] Anthropic API direct mode
- [ ] Agent memory store
- [ ] Dynamic agent spawning
- [ ] VS Code extension

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

Apache 2.0 — see [LICENSE](LICENSE) for details.
