package worker

import (
	"context"

	"github.com/Subodh8/ClaudeFlux/internal/config"
)

// Config holds the configuration for a worker process.
type Config struct {
	RunID     string
	AgentName string
	AgentCfg  config.AgentConfig
}

// Result holds the output of a worker execution.
type Result struct {
	TokensUsed int64
	CostUSD    float64
	Output     string
	ExitCode   int
}

// Worker represents a single agent worker process.
type Worker struct {
	cfg      Config
	worktree *Worktree
}

// NewWorker creates a new worker with the given configuration.
func NewWorker(cfg Config) *Worker {
	return &Worker{cfg: cfg}
}

// Execute runs the worker's agent task and returns the result.
func (w *Worker) Execute(ctx context.Context) (*Result, error) {
	// TODO: Implement full Claude agent execution
	// 1. Create git worktree for isolation
	// 2. Spawn claude CLI process with the agent prompt
	// 3. Monitor token usage and enforce budget
	// 4. Parse output and collect results
	// 5. Clean up worktree on completion
	return &Result{}, nil
}
