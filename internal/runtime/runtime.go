package runtime
 
import (
	"context"
	"fmt"
	"sync"
	"time"
 
	"github.com/google/uuid"
	"github.com/Subodh8/ClaudeFlux/internal/budget"
	"github.com/Subodh8/ClaudeFlux/internal/config"
	"github.com/Subodh8/ClaudeFlux/internal/coordinator"
	"github.com/Subodh8/ClaudeFlux/internal/dag"
	"github.com/Subodh8/ClaudeFlux/internal/ipc"
	"github.com/Subodh8/ClaudeFlux/internal/store"
	"github.com/Subodh8/ClaudeFlux/internal/worker"
	"go.uber.org/zap"
)
 
// Options configures the ClaudeFlux runtime.
type Options struct {
	Config           *config.Workflow
	Store            *store.Store
	Logger           *zap.Logger
	DryRun           bool
	DashboardEnabled bool
	DashboardPort    int
}
 
// Runtime orchestrates a workflow execution.
type Runtime struct {
	opts        Options
	dag         *dag.DAG
	pool        *WorkerPool
	coord       *coordinator.Coordinator
	budget      *budget.Tracker
	broker      *ipc.Broker
	mu          sync.RWMutex
	agentStates map[string]*AgentState
}
 
type AgentState struct {
	AgentName  string
	Status     Status
	StartedAt  *time.Time
	EndedAt    *time.Time
	TokensUsed int64
	CostUSD    float64
	Error      error
}
 
type Status string
 
const (
	StatusPending  Status = "pending"
	StatusRunning  Status = "running"
	StatusComplete Status = "complete"
	StatusFailed   Status = "failed"
	StatusBlocked  Status = "blocked"
)
 
func New(opts Options) (*Runtime, error) {
	if opts.Config == nil {
		return nil, fmt.Errorf("config is required")
	}
	if opts.Store == nil {
		return nil, fmt.Errorf("store is required")
	}
 
	d, err := dag.Build(opts.Config.Agents)
	if err != nil {
		return nil, fmt.Errorf("dag build: %w", err)
	}
 
	budgetTracker := budget.New(budget.Config{
		MaxTokens:      opts.Config.Budget.MaxTokens,
		MaxCostUSD:     opts.Config.Budget.MaxCostUSD,
		AlertAtPercent: opts.Config.Budget.AlertAtPercent,
		OnExceeded:     opts.Config.Budget.OnExceeded,
	})
	broker := ipc.NewBroker()
 
	coord, err := coordinator.New(coordinator.Options{
		Logger: opts.Logger,
		Store:  opts.Store,
		Broker: broker,
	})
	if err != nil {
		return nil, fmt.Errorf("coordinator init: %w", err)
	}
 
	pool := NewWorkerPool(WorkerPoolOptions{
		Logger:      opts.Logger,
		Budget:      budgetTracker,
		Broker:      broker,
		Coordinator: coord,
		DryRun:      opts.DryRun,
	})
 
	states := make(map[string]*AgentState, len(opts.Config.Agents))
	for name := range opts.Config.Agents {
		states[name] = &AgentState{AgentName: name, Status: StatusPending}
	}
 
	return &Runtime{
		opts:        opts,
		dag:         d,
		pool:        pool,
		coord:       coord,
		budget:      budgetTracker,
		broker:      broker,
		agentStates: states,
	}, nil
}
 
func (r *Runtime) Run(ctx context.Context) (string, error) {
	runID := uuid.New().String()
 
	r.opts.Logger.Info("starting workflow",
		zap.String("run_id", runID),
		zap.String("workflow", r.opts.Config.Name),
		zap.Int("agents", len(r.opts.Config.Agents)),
		zap.Bool("dry_run", r.opts.DryRun),
	)
 
	if err := r.opts.Store.CreateRun(runID, r.opts.Config.Name); err != nil {
		return runID, fmt.Errorf("store: create run: %w", err)
	}
 
	if r.opts.DashboardEnabled {
		if err := r.startDashboard(ctx, runID); err != nil {
			r.opts.Logger.Warn("dashboard failed to start", zap.Error(err))
		}
	}
 
	layers := r.dag.TopologicalLayers()
 
	for layerIdx, layer := range layers {
		r.opts.Logger.Debug("executing layer",
			zap.Int("layer", layerIdx),
			zap.Strings("agents", layer),
		)
 
		if err := r.executeLayer(ctx, runID, layer); err != nil {
			_ = r.opts.Store.FailRun(runID, err.Error())
			return runID, err
		}
 
		if r.budget.Exceeded() {
			err := fmt.Errorf("workflow budget exceeded: %s", r.budget.Summary())
			_ = r.opts.Store.FailRun(runID, err.Error())
			return runID, err
		}
	}
 
	if err := r.opts.Store.CompleteRun(runID); err != nil {
		r.opts.Logger.Warn("failed to mark run complete", zap.Error(err))
	}
 
	return runID, nil
}
 
func (r *Runtime) executeLayer(ctx context.Context, runID string, agents []string) error {
	var wg sync.WaitGroup
	errCh := make(chan error, len(agents))
 
	for _, agentName := range agents {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			if err := r.runAgent(ctx, runID, name); err != nil {
				errCh <- fmt.Errorf("agent %q: %w", name, err)
			}
		}(agentName)
	}
 
	wg.Wait()
	close(errCh)
 
	var errs []error
	for err := range errCh {
		errs = append(errs, err)
	}
 
	if len(errs) > 0 {
		return fmt.Errorf("layer failed: %v", errs)
	}
	return nil
}
 
func (r *Runtime) runAgent(ctx context.Context, runID, agentName string) error {
	agentCfg := r.opts.Config.Agents[agentName]
 
	r.setAgentStatus(agentName, StatusRunning)
	_ = r.opts.Store.CreateAgentRun(runID, agentName)
 
	w, err := r.pool.Acquire(ctx, worker.Config{
		RunID:     runID,
		AgentName: agentName,
		AgentCfg:  agentCfg,
	})
	if err != nil {
		r.setAgentStatus(agentName, StatusFailed)
		return fmt.Errorf("acquire worker: %w", err)
	}
	defer r.pool.Release(w)
 
	result, err := w.Execute(ctx)
	if err != nil {
		r.setAgentStatus(agentName, StatusFailed)
		_ = r.opts.Store.FailAgentRun(runID, agentName, err.Error())
		return err
	}
 
	r.budget.Record(result.TokensUsed, result.CostUSD)
	r.setAgentStatus(agentName, StatusComplete)
	_ = r.opts.Store.CompleteAgentRun(runID, agentName, result)
 
	return nil
}
 
func (r *Runtime) setAgentStatus(name string, s Status) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if state, ok := r.agentStates[name]; ok {
		state.Status = s
		now := time.Now()
		if s == StatusRunning {
			state.StartedAt = &now
		} else if s == StatusComplete || s == StatusFailed {
			state.EndedAt = &now
		}
	}
}
 
func (r *Runtime) startDashboard(ctx context.Context, runID string) error {
	r.opts.Logger.Info("dashboard starting",
		zap.Int("port", r.opts.DashboardPort),
		zap.String("url", fmt.Sprintf("http://localhost:%d", r.opts.DashboardPort)),
	)
	return r.broker.StartDashboardServer(ctx, runID, r.opts.DashboardPort)
}
