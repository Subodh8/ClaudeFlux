package config

type AgentConfig struct{}

type Budget struct {
	MaxTokens      int64
	MaxCostUSD     float64
	AlertAtPercent int
	OnExceeded     string
}

type Workflow struct {
	Name   string
	Agents map[string]AgentConfig
	Budget Budget
}

func Load(path string) (*Workflow, error) {
	return &Workflow{}, nil
}
