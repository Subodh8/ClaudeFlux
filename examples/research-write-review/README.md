# Research → Write → Review

The flagship ClaudeFlux example: a five-agent pipeline that researches a topic, writes a technical article, runs adversarial review, revises the draft, and commits via the steward.

## Agents

| Agent | Type | Role |
|-------|------|------|
| `researcher` | worker | Researches the topic and outputs structured JSON |
| `writer` | worker | Writes a 2000-word technical article |
| `critic` | critic | Adversarial review of the article |
| `revisor` | worker | Addresses all critical/major issues |
| `steward` | steward | Creates a draft PR after human approval |

## Usage

```bash
claudeflux run examples/research-write-review/workflow.yaml --dashboard
```

## DAG

```
researcher → writer → critic → revisor → steward
```

## Output Files

- `research_output.json` — Structured research findings
- `article.md` — First draft
- `critic_review.json` — Adversarial review
- `article_final.md` — Revised article
- `revision_log.json` — What was changed and why
