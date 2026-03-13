package coordinator

import (
	"github.com/Subodh8/ClaudeFlux/internal/ipc"
	"github.com/Subodh8/ClaudeFlux/internal/store"
	"go.uber.org/zap"
)

// Options configures the Coordinator.
type Options struct {
	Logger *zap.Logger
	Store  *store.Store
	Broker *ipc.Broker
}

// Coordinator routes messages between agents and manages approval gates.
type Coordinator struct {
	opts Options
}

// New creates a new Coordinator with the given options.
func New(opts Options) (*Coordinator, error) {
	return &Coordinator{opts: opts}, nil
}
