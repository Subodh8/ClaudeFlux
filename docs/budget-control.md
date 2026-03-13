# Budget Control

ClaudeFlux enforces token and cost budgets at both the workflow and agent level.

## Workflow-Level Budget

```yaml
budget:
  max_tokens: 100000       # Total tokens across all agents
  max_cost_usd: 2.00       # Total cost cap in USD
  alert_at_percent: 80     # Alert when 80% consumed
  on_exceeded: checkpoint  # What to do when exceeded
```

### `on_exceeded` Actions

| Action | Behaviour |
|--------|-----------|
| `checkpoint` | Save state, stop gracefully, allow resume |
| `pause` | Pause execution, wait for human intervention |
| `abort` | Terminate immediately, mark run as failed |

## Per-Agent Budget

```yaml
agents:
  researcher:
    budget:
      max_tokens: 30000
      max_cost_usd: 0.60
```

Per-agent budgets are enforced independently. An agent that exceeds its own budget is stopped even if the workflow budget has room remaining.

## Cost Tracking

ClaudeFlux tracks costs using the Claude API pricing model. Current rates are defined in `internal/budget/budget.go`.

## Alerts

When the `alert_at_percent` threshold is reached, ClaudeFlux:
1. Logs a warning
2. Sends a notification via configured channels (Slack, webhook)
3. Continues execution (alerts don't stop the workflow)

## Monitoring

Use the dashboard cost heatmap (`/costs`) or the CLI:

```bash
claudeflux status <run-id>
```
