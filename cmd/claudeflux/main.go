package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/Subodh8/ClaudeFlux/internal/config"
	"github.com/Subodh8/ClaudeFlux/internal/runtime"
	"github.com/Subodh8/ClaudeFlux/internal/store"
)

var version = "dev" // injected by goreleaser

func main() {
	if err := rootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func rootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "claudeflux",
		Short: "Production-grade multi-agent orchestration for Claude",
		Long: `ClaudeFlux is a runtime for orchestrating multiple Claude agents
with DAG scheduling, git isolation, token budgets, and approval gates.`,
	}

	root.AddCommand(
		runCmd(),
		resumeCmd(),
		logsCmd(),
		runsCmd(),
		approveCmd(),
		rejectCmd(),
		statusCmd(),
		validateCmd(),
		estimateCmd(),
		versionCmd(),
	)

	return root
}

func runCmd() *cobra.Command {
	var (
		dashboardEnabled bool
		dashboardPort    int
		dryRun           bool
		logLevel         string
		stateDir         string
	)

	cmd := &cobra.Command{
		Use:   "run <workflow.yaml>",
		Short: "Execute a ClaudeFlux workflow",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := signal.NotifyContext(
				context.Background(),
				syscall.SIGINT, syscall.SIGTERM,
			)
			defer cancel()

			logger, err := buildLogger(logLevel)
			if err != nil {
				return fmt.Errorf("logger init: %w", err)
			}
			defer logger.Sync()

			cfg, err := config.Load(args[0])
			if err != nil {
				return fmt.Errorf("config: %w", err)
			}

			db, err := store.Open(stateDir)
			if err != nil {
				return fmt.Errorf("store: %w", err)
			}
			defer db.Close()

			rt, err := runtime.New(runtime.Options{
				Config:           cfg,
				Store:            db,
				Logger:           logger,
				DryRun:           dryRun,
				DashboardEnabled: dashboardEnabled,
				DashboardPort:    dashboardPort,
			})
			if err != nil {
				return fmt.Errorf("runtime init: %w", err)
			}

			runID, err := rt.Run(ctx)
			if err != nil {
				logger.Error("workflow failed", zap.Error(err), zap.String("run_id", runID))
				return err
			}

			logger.Info("workflow completed", zap.String("run_id", runID))
			return nil
		},
	}

	cmd.Flags().BoolVar(&dashboardEnabled, "dashboard", false, "Launch live dashboard")
	cmd.Flags().IntVar(&dashboardPort, "dashboard-port", 7070, "Dashboard port")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Validate and estimate cost without executing")
	cmd.Flags().StringVar(&logLevel, "log-level", "info", "Log level (debug|info|warn|error)")
	cmd.Flags().StringVar(&stateDir, "state-dir", ".claudeflux", "Directory for run state and audit log")

	return cmd
}

func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("claudeflux %s\n", version)
		},
	}
}

func buildLogger(level string) (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	if level == "debug" {
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}
	return cfg.Build()
}

var errNotImplemented = fmt.Errorf("not yet implemented — see the roadmap at https://github.com/Subodh8/ClaudeFlux")

func stubCmd(use, short string) *cobra.Command {
	return &cobra.Command{
		Use:   use,
		Short: short,
		RunE: func(cmd *cobra.Command, args []string) error {
			return errNotImplemented
		},
	}
}

func resumeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "resume <run-id>",
		Short: "Resume a paused or failed workflow run",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("resume not yet implemented")
		},
	}
}

func logsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "logs <run-id>",
		Short: "View audit log for a workflow run",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("logs not yet implemented")
		},
	}
}

func runsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "runs",
		Short: "List all workflow runs",
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("runs not yet implemented")
		},
	}
}

func approveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "approve <gate-id>",
		Short: "Approve a pending approval gate",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("approve not yet implemented")
		},
	}
}

func rejectCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "reject <gate-id>",
		Short: "Reject a pending approval gate",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("reject not yet implemented")
		},
	}
}

func statusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status <run-id>",
		Short: "Show real-time agent status for a run",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("status not yet implemented")
		},
	}
}

func validateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "validate <workflow.yaml>",
		Short: "Validate a workflow file without running it",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("validate not yet implemented")
		},
	}
}

func estimateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "estimate <workflow.yaml>",
		Short: "Estimate token cost before running a workflow",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("estimate not yet implemented")
		},
	}
}
