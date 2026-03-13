# Data Pipeline

An ETL-style data processing workflow:

1. **Extractor** — Pulls data from APIs and normalizes responses
2. **Transformer** — Deduplicates, normalizes, and computes derived metrics
3. **Validator** (Critic) — Checks data integrity and flags anomalies

## Usage

```bash
claudeflux run examples/data-pipeline/workflow.yaml --dashboard
```
