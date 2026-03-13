package worker

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Subodh8/ClaudeFlux/internal/config"
)

type Config struct {
	RunID     string
	AgentName string
	AgentCfg  config.AgentConfig
}

type Result struct {
	TokensUsed int64
	CostUSD    float64
	Output     string
	ExitCode   int
}

type Worker struct {
	cfg      Config
	worktree *Worktree
}

func NewWorker(cfg Config) *Worker {
	return &Worker{cfg: cfg}
}

// outputFilename picks the right filename for each agent's output
func outputFilename(agentName string) string {
	names := map[string]string{
		"researcher":  "research_output.json",
		"writer":      "article.md",
		"critic":      "critic_review.json",
		"revisor":     "article_final.md",
		"steward":     "steward_report.md",
		"analyzer":    "analysis.json",
		"reviewer":    "reviewed_analysis.json",
		"fixer":       "fix_summary.md",
		"extractor":   "raw_data.json",
		"transformer": "transformed_data.json",
		"validator":   "validation_report.json",
		"author":      "security_policy.md",
		"red_team":    "red_team_report.json",
		"author_v2":   "security_policy_v2.md",
	}
	if name, ok := names[agentName]; ok {
		return name
	}
	return agentName + "_output.txt"
}

func (w *Worker) Execute(ctx context.Context) (*Result, error) {
	prompt := w.cfg.AgentCfg.Prompt
	if prompt == "" {
		prompt = "You are a steward agent. Review the workflow outputs and confirm completion."
	}

	instruction := "IMPORTANT: Do not attempt to write any files. Do not ask for file write permissions. Just output the content directly as plain text or JSON in your response. Your response will be saved automatically. After incorporating all revisions, remove all <!-- CF-REVISION: ... --> comment tags from the final output. The delivered article should be clean prose with no editorial markers visible.\n\n"

	args := []string{
		"--print",
		instruction + prompt,
	}

	cmd := exec.CommandContext(ctx, "claude", args...)
	cmd.Dir = "."

	cmd.Env = os.Environ()
	for k, v := range w.cfg.AgentCfg.Env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("claude process failed: %w\ndetails: %s", err, stderr.String())
	}

	outputStr := strings.TrimSpace(string(output))

	// Save the output to a file named after the agent
	filename := outputFilename(w.cfg.AgentName)
	outputPath := filepath.Join(".", filename)
	if writeErr := os.WriteFile(outputPath, []byte(outputStr), 0644); writeErr != nil {
		fmt.Printf("warning: could not save output file %s: %v\n", filename, writeErr)
	} else {
		fmt.Printf("agent %s wrote output to %s\n", w.cfg.AgentName, filename)
	}

	return &Result{
		Output:   outputStr,
		ExitCode: 0,
	}, nil
}
