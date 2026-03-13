package budget

import (
	"fmt"
	"sync/atomic"
)

// Config defines the budget constraints for a workflow.
type Config struct {
	MaxTokens      int64
	MaxCostUSD     float64
	AlertAtPercent int
	OnExceeded     string // checkpoint | pause | abort
}

// Tracker provides thread-safe budget tracking.
type Tracker struct {
	cfg            Config
	totalTokens    atomic.Int64
	totalCostMicro atomic.Int64
}

// New creates a new budget Tracker with the given config.
func New(cfg Config) *Tracker {
	return &Tracker{cfg: cfg}
}

// Record adds token usage and cost to the running totals.
func (t *Tracker) Record(tokens int64, costUSD float64) {
	t.totalTokens.Add(tokens)
	t.totalCostMicro.Add(int64(costUSD * 1_000_000))
}

// Exceeded returns true if any budget limit has been reached.
func (t *Tracker) Exceeded() bool {
	if t.cfg.MaxTokens > 0 && t.totalTokens.Load() >= t.cfg.MaxTokens {
		return true
	}
	if t.cfg.MaxCostUSD > 0 {
		spentMicro := t.totalCostMicro.Load()
		limitMicro := int64(t.cfg.MaxCostUSD * 1_000_000)
		if spentMicro >= limitMicro {
			return true
		}
	}
	return false
}

// ShouldAlert returns true if the alert threshold has been reached.
func (t *Tracker) ShouldAlert() bool {
	if t.cfg.AlertAtPercent <= 0 {
		return false
	}
	pct := float64(t.cfg.AlertAtPercent)
	if t.cfg.MaxTokens > 0 {
		used := float64(t.totalTokens.Load())
		if used/float64(t.cfg.MaxTokens)*100 >= pct {
			return true
		}
	}
	if t.cfg.MaxCostUSD > 0 {
		spentMicro := float64(t.totalCostMicro.Load())
		limitMicro := float64(int64(t.cfg.MaxCostUSD * 1_000_000))
		if spentMicro/limitMicro*100 >= pct {
			return true
		}
	}
	return false
}

// Summary returns a human-readable summary of current usage.
func (t *Tracker) Summary() string {
	tokens := t.totalTokens.Load()
	costUSD := float64(t.totalCostMicro.Load()) / 1_000_000
	return fmt.Sprintf("tokens=%d cost=$%.4f", tokens, costUSD)
}
