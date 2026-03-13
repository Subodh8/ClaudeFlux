# Approval Gates

Approval gates insert human-in-the-loop checkpoints into your workflow DAG.

## How It Works

1. An agent reaches an approval gate
2. ClaudeFlux pauses execution for that branch of the DAG
3. A notification is sent (dashboard, Slack, email)
4. A human approves or rejects
5. Execution continues or the workflow handles the rejection

## Configuration

```yaml
agents:
  deploy:
    depends_on: [test]
    approval:
      required: true
      timeout: 3600          # Auto-reject after 1 hour (0 = no timeout)
      notify:
        - slack: "#deployments"
        - email: "team@yourcompany.com"
```

## CLI Commands

```bash
# Approve a pending gate
claudeflux approve <gate-id>

# Reject a pending gate with reason
claudeflux reject <gate-id> --reason "Output quality too low"
```

## Dashboard

When `--dashboard` is enabled, pending approvals appear in the Approval Queue page at `http://localhost:7070/approvals`.

## Timeout Behaviour

If `timeout` is set and the approval is not acted upon within that time:
- The gate is automatically **rejected**
- The rejection reason is recorded as "timeout"
- Downstream agents are not executed
