package budget

import (
	"fmt"
	"sync/atomic"
)
 
type Config struct {
	MaxTokens      int64
	MaxCostUSD     float64
	AlertAtPercent int
	OnExceeded     string // checkpoint | pause | abort
}
 
type Tracker struct {
	cfg            Config
	totalTokens    atomic.Int64
	totalCostMicro atomic.Int64
}
 
func New(cfg Config) *Tracker {
	return &Tracker{cfg: cfg}
}
 
func (t *Tracker) Record(tokens int64, costUSD float64) {
	t.totalTokens.Add(tokens)
	t.totalCostMicro.Add(int64(costUSD * 1_000_000))
}
 
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
 
func (t *Tracker) ShouldAlert() bool {
	if t.cfg.AlertAtPercent <= 0 {
		return false
	}
	pct := t.cfg.AlertAtPercent
	if t.cfg.MaxTokens > 0 {
		used := t.totalTokens.Load()
		if int(used*100/t.cfg.MaxTokens) >= pct {
			return true
		}
	}
	if t.cfg.MaxCostUSD > 0 {
		spentMicro := t.totalCostMicro.Load()
		limitMicro := int64(t.cfg.MaxCostUSD * 1_000_000)
		if int(spentMicro*100/limitMicro) >= pct {
			return true
		}
	}
	return false
}
 
func (t *Tracker) Summary() string {
	tokens := t.totalTokens.Load()
	costUSD := float64(t.totalCostMicro.Load()) / 1_000_000
	return fmt.Sprintf("tokens=%d cost=$%.4f", tokens, costUSD)
}
