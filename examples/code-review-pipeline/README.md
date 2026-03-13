# Code Review Pipeline

An automated code review workflow using three agents:

1. **Analyzer** — Scans the codebase for bugs, security issues, and anti-patterns
2. **Reviewer** (Critic) — Reviews findings for false positives and prioritizes by severity
3. **Fixer** — Generates patches for critical and major issues
4. **Steward** — Creates a draft PR with the fixes after human approval

## Usage

```bash
claudeflux run examples/code-review-pipeline/workflow.yaml --dashboard
```

## Customization

Edit the `analyzer` prompt to focus on specific languages, frameworks, or coding standards relevant to your project.
