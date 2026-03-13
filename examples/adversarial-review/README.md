# Adversarial Review

A red team / blue team review pattern:

1. **Author** — Writes a security policy document
2. **Red Team** (Critic) — Attacks the document, finding gaps and weaknesses
3. **Author v2** — Revises the document addressing all critical findings
4. **Steward** — Creates a draft PR after human approval

## Usage

```bash
claudeflux run examples/adversarial-review/workflow.yaml --dashboard
```

## When to Use

This pattern is ideal for any document that needs adversarial validation:
- Security policies
- Architecture decisions
- Legal documents
- API specifications
