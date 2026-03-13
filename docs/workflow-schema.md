# Workflow Schema Reference

Complete reference for ClaudeFlux workflow YAML files.

## Top-Level Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `version` | string | ✅ | Schema version. Currently `"1"`. |
| `name` | string | ✅ | Workflow name. |
| `description` | string | | Optional description. |
| `git` | object | | Git configuration. |
| `budget` | object | | Workflow-level budget limits. |
| `contracts` | object | | JSON Schema definitions for inter-agent data. |
| `agents` | object | ✅ | Agent definitions. |

## Git Configuration

```yaml
git:
  repo: "."              # Repository path. Default: current directory
  base_branch: "main"    # Default: main
  worktree_prefix: "cf-" # Default: cf-
```

## Budget Configuration

```yaml
budget:
  max_tokens: 100000       # Workflow-level hard cap
  max_cost_usd: 2.00       # Workflow-level cost cap
  alert_at_percent: 80     # Alert threshold (default: 80)
  on_exceeded: checkpoint  # checkpoint | pause | abort
```

## Agent Configuration

```yaml
agents:
  <name>:
    prompt: string          # Required. Agent system prompt.
    type: string            # worker | steward | critic (default: worker)
    depends_on: [string]    # Agent names this agent waits for.
    budget:
      max_tokens: int
      max_cost_usd: float
    approval:
      required: bool
      timeout: int          # Seconds. 0 = no timeout.
      notify: [object]
    contract:
      input: string         # JSON Schema ref for input validation
      output: string        # JSON Schema ref for output validation
    env: {string: string}   # Env vars injected into the worker process
    timeout: int            # Seconds. Agent-level execution timeout.
    retry:
      max_attempts: int     # Default: 1
      backoff: string       # linear | exponential
    targets: [string]       # For critic agents: which agents to review
    rules: object           # For steward agents: merge/PR rules
```
