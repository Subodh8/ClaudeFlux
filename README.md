<div align="center">
 
# ClaudeFlux
 
### Production-grade multi-agent orchestration runtime for Claude.
### DAG scheduling · git isolation · token budgets · approval gates · live dashboard.
 
[![CI](https://github.com/yourusername/claudeflux/actions/workflows/ci.yml/badge.svg)]
[![Go Report Card](https://goreportcard.com/badge/github.com/yourusername/claudeflux)]
[![License: Apache 2.0](https://img.shields.io/badge/License-Apache%202.0-blue.svg)]
[![Stars](https://img.shields.io/github/stars/yourusername/claudeflux?style=social)]
 
[Quickstart] · [Documentation] · [Examples] · [Discord]
 
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

```ascii
┌─────────────────────────────────────────────────────────────────┐
│                        ClaudeFlux Runtime                       │
│                                                                 │
│   ┌───────────┐    ┌──────────────┐    ┌───────────────────┐   │
│   │ YAML DAG  │───▶│  DAG Parser  │───▶│  Dispatch Daemon  │   │
│   │  Workflow │    │  + Validator │    │  (dependency res) │   │
│   └───────────┘    └──────────────┘    └────────┬──────────┘   │
│                                                 │              │
│                    ┌────────────────────────────▼──────────┐   │
│                    │           Worker Pool Manager          │   │
│   ┌────────────────▼──────────────────────────────────┐    │   │
│   │                  Worker Process (per agent)        │    │   │
│   │  ┌─────────────┐  ┌──────────┐  ┌─────────────┐  │    │   │
│   │  │ git worktree│  │  Claude  │  │   Token $   │  │    │   │
│   │  │  isolation  │  │  Code    │  │   Budget    │  │    │   │
│   │  │  (isolated) │  │  (proc)  │  │  Enforcer   │  │    │   │
│   │  └─────────────┘  └──────────┘  └─────────────┘  │    │   │
│   └────────────────────────────────────────────────────┘    │   │
│                                                              │   │
│   ┌──────────────┐    ┌──────────────┐    ┌─────────────┐   │   │
│   │ Coordinator  │    │   Steward    │    │  Rate Limit │   │   │
│   │  (routing +  │    │  (reviewer + │    │ Coordinator │   │   │
│   │   approval)  │    │  auto-merge) │    │             │   │   │
│   └──────────────┘    └──────────────┘    └─────────────┘   │   │
│                                                              │   │
│   ┌──────────────────────────────────────────────────────┐   │   │
│   │              SQLite Audit Log (every event)           │   │   │
│   └──────────────────────────────────────────────────────┘   │   │
└─────────────────────────────────────────────────────────────────┘
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
brew tap yourusername/claudeflux
brew install claudeflux
```
 
### Option 3: Docker
 
```bash
docker run --rm -v $(pwd):/workspace ghcr.io/yourusername/claudeflux:latest run workflow.yaml
```
 
### Option 4: Build from source
 
```bash
git clone https://github.com/yourusername/claudeflux
cd claudeflux
make build
./bin/claudeflux version
```
 
Prerequisite: claude CLI must be installed and authenticated. See https://docs.anthropic.com/claude-code

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

`⚠️ = partial / requires custom implementation`
