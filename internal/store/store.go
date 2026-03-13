package store

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

type Store struct {
	db *sql.DB
}

func Open(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("create state dir: %w", err)
	}
	dbPath := filepath.Join(dir, "claudeflux.db")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	if err := runMigrations(db); err != nil {
		return nil, fmt.Errorf("migrations: %w", err)
	}
	return &Store{db: db}, nil
}

func (s *Store) Close() error { return s.db.Close() }

func (s *Store) CreateRun(runID, name string) error {
	_, err := s.db.Exec(
		`INSERT INTO runs (id, name, status) VALUES (?, ?, 'running')`,
		runID, name,
	)
	return err
}

func (s *Store) FailRun(runID, errMsg string) error {
	_, err := s.db.Exec(
		`UPDATE runs SET status='failed', error=?, ended_at=CURRENT_TIMESTAMP WHERE id=?`,
		errMsg, runID,
	)
	return err
}

func (s *Store) CompleteRun(runID string) error {
	_, err := s.db.Exec(
		`UPDATE runs SET status='complete', ended_at=CURRENT_TIMESTAMP WHERE id=?`,
		runID,
	)
	return err
}

func (s *Store) CreateAgentRun(runID, name string) error {
	_, err := s.db.Exec(
		`INSERT OR IGNORE INTO agent_runs (run_id, agent_name, status, started_at)
		 VALUES (?, ?, 'running', CURRENT_TIMESTAMP)`,
		runID, name,
	)
	return err
}

func (s *Store) FailAgentRun(runID, name, errMsg string) error {
	_, err := s.db.Exec(
		`UPDATE agent_runs SET status='failed', error=?, ended_at=CURRENT_TIMESTAMP
		 WHERE run_id=? AND agent_name=?`,
		errMsg, runID, name,
	)
	return err
}

func (s *Store) CompleteAgentRun(runID, name string, result any) error {
	// Extract token/cost info if result is a worker.Result
	type resultInfo interface {
		GetTokensUsed() int64
		GetCostUSD() float64
		GetOutput() string
	}

	var tokens int64
	var cost float64
	var output string

	if r, ok := result.(resultInfo); ok {
		tokens = r.GetTokensUsed()
		cost = r.GetCostUSD()
		output = r.GetOutput()
	}

	_, err := s.db.Exec(
		`UPDATE agent_runs SET status='complete', ended_at=CURRENT_TIMESTAMP,
		 tokens_used=?, cost_usd=?, output=?
		 WHERE run_id=? AND agent_name=?`,
		tokens, cost, output, runID, name,
	)
	return err
}

func runMigrations(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS runs (
			id          TEXT PRIMARY KEY,
			name        TEXT NOT NULL,
			status      TEXT NOT NULL DEFAULT 'running',
			started_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
			ended_at    DATETIME,
			error       TEXT
		);
		CREATE TABLE IF NOT EXISTS agent_runs (
			id          INTEGER PRIMARY KEY AUTOINCREMENT,
			run_id      TEXT NOT NULL,
			agent_name  TEXT NOT NULL,
			status      TEXT NOT NULL DEFAULT 'pending',
			started_at  DATETIME,
			ended_at    DATETIME,
			tokens_used INTEGER DEFAULT 0,
			cost_usd    REAL DEFAULT 0.0,
			output      TEXT,
			error       TEXT,
			UNIQUE(run_id, agent_name)
		);
		CREATE TABLE IF NOT EXISTS events (
			id         INTEGER PRIMARY KEY AUTOINCREMENT,
			run_id     TEXT NOT NULL,
			agent_name TEXT,
			event_type TEXT NOT NULL,
			payload    TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
		CREATE TABLE IF NOT EXISTS approvals (
			id           INTEGER PRIMARY KEY AUTOINCREMENT,
			run_id       TEXT NOT NULL,
			agent_name   TEXT NOT NULL,
			status       TEXT NOT NULL DEFAULT 'pending',
			requested_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			decided_at   DATETIME,
			decided_by   TEXT,
			reason       TEXT
		);
	`)
	return err
}