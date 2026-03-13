package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Workflow represents a complete ClaudeFlux workflow configuration.
type Workflow struct {
	Version     string                 `yaml:"version"`
	Name        string                 `yaml:"name"`
	Description string                 `yaml:"description,omitempty"`
	Git         GitConfig              `yaml:"git,omitempty"`
	Budget      BudgetConfig           `yaml:"budget,omitempty"`
	Contracts   map[string]interface{} `yaml:"contracts,omitempty"`
	Agents      map[string]AgentConfig `yaml:"agents"`
}

// GitConfig holds git-related workflow settings.
type GitConfig struct {
	Repo            string `yaml:"repo,omitempty"`
	BaseBranch      string `yaml:"base_branch,omitempty"`
	WorktreePrefix  string `yaml:"worktree_prefix,omitempty"`
}

// BudgetConfig holds workflow-level budget constraints.
type BudgetConfig struct {
	MaxTokens      int64   `yaml:"max_tokens,omitempty"`
	MaxCostUSD     float64 `yaml:"max_cost_usd,omitempty"`
	AlertAtPercent int     `yaml:"alert_at_percent,omitempty"`
	OnExceeded     string  `yaml:"on_exceeded,omitempty"`
}

// AgentConfig holds the configuration for a single agent.
type AgentConfig struct {
	Type      string            `yaml:"type,omitempty"`
	Prompt    string            `yaml:"prompt"`
	DependsOn []string          `yaml:"depends_on,omitempty"`
	Budget    AgentBudget       `yaml:"budget,omitempty"`
	Approval  ApprovalConfig    `yaml:"approval,omitempty"`
	Contract  ContractRef       `yaml:"contract,omitempty"`
	Env       map[string]string `yaml:"env,omitempty"`
	Timeout   int               `yaml:"timeout,omitempty"`
	Retry     RetryConfig       `yaml:"retry,omitempty"`
	Targets   []string          `yaml:"targets,omitempty"`
	Rules     map[string]interface{} `yaml:"rules,omitempty"`
}

// AgentBudget holds per-agent budget limits.
type AgentBudget struct {
	MaxTokens  int64   `yaml:"max_tokens,omitempty"`
	MaxCostUSD float64 `yaml:"max_cost_usd,omitempty"`
}

// ApprovalConfig holds approval gate settings.
type ApprovalConfig struct {
	Required bool          `yaml:"required,omitempty"`
	Timeout  int           `yaml:"timeout,omitempty"`
	Notify   []interface{} `yaml:"notify,omitempty"`
}

// ContractRef references input/output JSON Schema contracts.
type ContractRef struct {
	Input  string `yaml:"input,omitempty"`
	Output string `yaml:"output,omitempty"`
}

// RetryConfig holds retry behaviour for agents.
type RetryConfig struct {
	MaxAttempts int    `yaml:"max_attempts,omitempty"`
	Backoff     string `yaml:"backoff,omitempty"`
}

// Load reads and parses a workflow YAML configuration file.
func Load(path string) (*Workflow, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read workflow file: %w", err)
	}

	var wf Workflow
	if err := yaml.Unmarshal(data, &wf); err != nil {
		return nil, fmt.Errorf("parse workflow YAML: %w", err)
	}

	if wf.Name == "" {
		return nil, fmt.Errorf("workflow name is required")
	}
	if len(wf.Agents) == 0 {
		return nil, fmt.Errorf("workflow must define at least one agent")
	}

	// Apply defaults
	if wf.Git.BaseBranch == "" {
		wf.Git.BaseBranch = "main"
	}
	if wf.Git.WorktreePrefix == "" {
		wf.Git.WorktreePrefix = "cf-"
	}
	if wf.Budget.AlertAtPercent == 0 {
		wf.Budget.AlertAtPercent = 80
	}

	return &wf, nil
}

// GetDependsOn returns the dependency list for an agent config.
func (a AgentConfig) GetDependsOn() []string {
	return a.DependsOn
}
