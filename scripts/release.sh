#!/usr/bin/env bash
set -euo pipefail

echo "→ Starting ClaudeFlux release process"

# Ensure we're on main
BRANCH=$(git rev-parse --abbrev-ref HEAD)
if [ "$BRANCH" != "main" ]; then
    echo "✗ Releases must be from the main branch. Currently on: $BRANCH"
    exit 1
fi

# Ensure working directory is clean
if [ -n "$(git status --porcelain)" ]; then
    echo "✗ Working directory is not clean. Commit or stash changes first."
    exit 1
fi

# Get version
VERSION=${1:-}
if [ -z "$VERSION" ]; then
    echo "Usage: ./scripts/release.sh v0.1.0"
    exit 1
fi

echo "→ Creating release $VERSION"

# Run tests
echo "→ Running tests..."
make test

# Tag
git tag -a "$VERSION" -m "Release $VERSION"
git push origin "$VERSION"

echo "✓ Tag $VERSION pushed. GitHub Actions will handle the release."
echo "  Monitor: https://github.com/Subodh8/ClaudeFlux/actions"
