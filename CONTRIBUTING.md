# Contributing to ClaudeFlux

Thank you for your interest in contributing to ClaudeFlux! This document provides guidelines and information for contributors.

## Development Setup

### Prerequisites

- Go 1.23+
- Node.js 20+ (for dashboard)
- Git
- SQLite3 development libraries

### Getting Started

```bash
git clone https://github.com/Subodh8/ClaudeFlux.git
cd ClaudeFlux
make build
make test
```

### Dashboard Development

```bash
cd dashboard
npm install
npm run dev
```

## Pull Request Process

1. Fork the repository and create your branch from `main`.
2. If you've added code that should be tested, add tests.
3. Ensure the test suite passes: `make test`
4. Ensure the linter passes: `make lint`
5. Update documentation if you've changed APIs.
6. Create the Pull Request.

## Code Style

- Go code should pass `golangci-lint` with our [.golangci.yml](.golangci.yml) config.
- Use `gofmt` and `goimports` for formatting.
- Write meaningful commit messages following [Conventional Commits](https://www.conventionalcommits.org/).

## Commit Message Format

```
<type>(<scope>): <description>

[optional body]
```

Types: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`

Examples:
- `feat(dag): add conditional branching support`
- `fix(budget): correct token counting for streaming responses`
- `docs(readme): update quickstart section`

## Reporting Bugs

Use the [bug report template](.github/ISSUE_TEMPLATE/bug_report.md) on GitHub Issues.

## Requesting Features

Use the [feature request template](.github/ISSUE_TEMPLATE/feature_request.md) on GitHub Issues.

## License

By contributing, you agree that your contributions will be licensed under the Apache 2.0 License.
