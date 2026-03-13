# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Core runtime with DAG-based workflow execution
- YAML workflow configuration parser
- Git worktree isolation for agent processes
- Token budget enforcement (per-agent and workflow-level)
- Approval gates with timeout support
- SQLite audit log for all events
- Worker pool for concurrent agent execution
- Typed inter-agent messaging via JSON Schema contracts
- Critic/adversarial agent type
- Steward agent for automated review and merge
- CLI commands: `run`, `resume`, `logs`, `status`, `validate`, `estimate`, `approve`, `reject`
- `--dashboard` flag for live DAG visualization
- Docker and Docker Compose support
- GitHub Actions CI pipeline
- GoReleaser configuration for cross-platform builds
- Example workflows: research-write-review, code-review-pipeline, data-pipeline, adversarial-review
