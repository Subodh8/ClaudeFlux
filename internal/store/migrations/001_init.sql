-- ClaudeFlux SQLite Schema
-- Migration 001: Initial tables

CREATE TABLE IF NOT EXISTS runs (
    id          TEXT PRIMARY KEY,
    name        TEXT NOT NULL,
    status      TEXT NOT NULL DEFAULT 'running',
    started_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
    ended_at    DATETIME,
    error       TEXT,
    config_yaml TEXT
);

CREATE TABLE IF NOT EXISTS agent_runs (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    run_id      TEXT NOT NULL REFERENCES runs(id),
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
    run_id     TEXT NOT NULL REFERENCES runs(id),
    agent_name TEXT,
    event_type TEXT NOT NULL,
    payload    TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS approvals (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    run_id      TEXT NOT NULL REFERENCES runs(id),
    agent_name  TEXT NOT NULL,
    status      TEXT NOT NULL DEFAULT 'pending',
    requested_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    decided_at  DATETIME,
    decided_by  TEXT,
    reason      TEXT
);

CREATE INDEX idx_agent_runs_run_id ON agent_runs(run_id);
CREATE INDEX idx_events_run_id ON events(run_id);
CREATE INDEX idx_approvals_run_id ON approvals(run_id);
CREATE INDEX idx_approvals_status ON approvals(status);
