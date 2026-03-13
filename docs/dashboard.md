# Live Dashboard

ClaudeFlux includes a Next.js dashboard for real-time workflow monitoring.

## Starting the Dashboard

```bash
claudeflux run workflow.yaml --dashboard --dashboard-port 7070
```

Then open `http://localhost:7070`.

## Pages

### DAG Visualization (`/`)
Interactive graph showing the workflow DAG with real-time agent status updates:
- 🟡 Pending
- 🔵 Running
- 🟢 Complete
- 🔴 Failed
- ⚪ Blocked

### Logs (`/logs`)
SSE-powered real-time log streaming from all agents. Filter by agent name or log level.

### Costs (`/costs`)
Cost heatmap showing token usage and USD cost per agent. Updates in real time.

### Approvals (`/approvals`)
Queue of pending approval gates. Approve or reject directly from the dashboard.

## Components

| Component | Description |
|-----------|-------------|
| `DAGVisualization` | React Flow graph of the agent dependency DAG |
| `AgentCard` | Status card for individual agents |
| `CostHeatmap` | Visual cost breakdown by agent |
| `ApprovalModal` | Approve/reject dialog for gates |
| `LogStream` | Real-time log viewer with filtering |

## Architecture

The dashboard connects to ClaudeFlux via SSE (Server-Sent Events) on the IPC port (default: 7071). Events are pushed in real time as agents execute.
