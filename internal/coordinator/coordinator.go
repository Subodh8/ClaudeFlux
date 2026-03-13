package coordinator

import "github.com/Subodh8/ClaudeFlux/internal/store"
import "github.com/Subodh8/ClaudeFlux/internal/ipc"
import "go.uber.org/zap"

type Options struct { Logger *zap.Logger; Store *store.Store; Broker *ipc.Broker }

type Coordinator struct{}

func New(opts Options) (*Coordinator, error) { return &Coordinator{}, nil }
