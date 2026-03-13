package ipc

import "context"

type Broker struct{}

func NewBroker() *Broker { return &Broker{} }

func (b *Broker) StartDashboardServer(ctx context.Context, runID string, port int) error { return nil }
