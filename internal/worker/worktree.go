package worker
 
import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)
 
// Worktree manages a git worktree for agent isolation.
type Worktree struct {
	RepoRoot     string
	WorktreePath string
	AgentName    string
	Branch       string
}
 
// NewWorktree creates and checks out a new git worktree for the agent.
func NewWorktree(repoRoot, agentName, runID string) (*Worktree, error) {
	branchName := fmt.Sprintf("cf/%s/%s", runID[:8], agentName)
	worktreeName := fmt.Sprintf("cf-%s-%s", agentName, runID[:8])
	worktreePath := filepath.Join(repoRoot, ".git", "worktrees-cf", worktreeName)
 
	if err := os.MkdirAll(filepath.Dir(worktreePath), 0755); err != nil {
		return nil, fmt.Errorf("mkdir worktree parent: %w", err)
	}
 
	cmd := exec.Command("git", "worktree", "add", "-b", branchName, worktreePath, "HEAD")
	cmd.Dir = repoRoot
	if out, err := cmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("git worktree add: %w\noutput: %s", err, out)
	}
 
	return &Worktree{
		RepoRoot:     repoRoot,
		WorktreePath: worktreePath,
		AgentName:    agentName,
		Branch:       branchName,
	}, nil
}
 
// Remove cleans up the worktree after the agent completes.
func (w *Worktree) Remove() error {
	cmd := exec.Command("git", "worktree", "remove", "--force", w.WorktreePath)
	cmd.Dir = w.RepoRoot
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git worktree remove: %w\noutput: %s", err, out)
	}
 
	del := exec.Command("git", "branch", "-D", w.Branch)
	del.Dir = w.RepoRoot
	_ = del.Run() // Best-effort
 
	return nil
}
 
// Path returns the filesystem path of the worktree.
func (w *Worktree) Path() string {
	return w.WorktreePath
}

// Stubs to resolve compilation errors
type Config struct {
	RunID     string
	AgentName string
	AgentCfg  interface{}
}

type Worker struct{}

func (w *Worker) Execute(ctx context.Context) (*Result, error) {
	return &Result{}, nil
}

type Result struct {
	TokensUsed int64
	CostUSD    float64
}
