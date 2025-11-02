package storage

import (
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/raulbattistini/private-empty/pkg/metrics"
	"github.com/raulbattistini/private-empty/pkg/server/models/enum"
	"github.com/raulbattistini/private-empty/pkg/storage/db_models"
)

type SQLiteStore struct {
	db   *sqlx.DB
	rows sqlx.Rows
}

func (sqlModels *SQLiteStore) ToStorage() *db_models.SQLiteStore {
	return &db_models.SQLiteStore{}
}

type metricRow struct {
	ID        int       `db:"id"`
	Timestamp time.Time `db:"timestamp"`
	Source    string    `db:"source"`
	Type      string    `db:"type"`
	Value     float64   `db:"value"`
	Severity  string    `db:"severity"`
	Message   string    `db:"message"`
}

func NewSQLiteStore(dbPath string) (*SQLiteStore, error) {
	db, err := sqlx.Connect("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Create schema
	schema := `
	CREATE TABLE IF NOT EXISTS metrics (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		timestamp DATETIME NOT NULL,
		source TEXT NOT NULL,
		type TEXT NOT NULL,
		value REAL NOT NULL,
		severity TEXT NOT NULL,
		message TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_timestamp ON metrics(timestamp);
	CREATE INDEX IF NOT EXISTS idx_source ON metrics(source);
	CREATE INDEX IF NOT EXISTS idx_severity ON metrics(severity);
	`

	if _, err := db.Exec(schema); err != nil {
		return nil, err
	}

	return &SQLiteStore{db: db}, nil
}

func (s *SQLiteStore) Receive(m metrics.Metric) {
	query := `
		INSERT INTO metrics (timestamp, source, type, value, severity, message)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	_, err := s.db.Exec(query, m.Timestamp, m.Source, m.Type, m.Value, m.Severity, m.Message)
	if err != nil {
		// Log but don't crash
	}
}

func (s *SQLiteStore) GetLatest() []metrics.Metric {
	query := `
		SELECT DISTINCT ON (source) 
			timestamp, source, type, value, severity, message
		FROM metrics
		ORDER BY source, timestamp DESC
		LIMIT 100
	`

	var rows []metricRow
	if err := s.db.Select(&rows, query); err != nil {
		return nil
	}

	return s.rowsToMetrics(rows)
}

func (s *SQLiteStore) GetHistory(source enum.RuleSource, limit int) []metrics.Metric {
	query := `
		SELECT timestamp, source, type, value, severity, message
		FROM metrics
		WHERE source = ?
		ORDER BY timestamp DESC
		LIMIT ?
	`

	var rows []metricRow
	if err := s.db.Select(&rows, query, source, limit); err != nil {
		return nil
	}

	return s.rowsToMetrics(rows)
}

func (s *SQLiteStore) rowsToMetrics(rows []metricRow) []metrics.Metric {
	result := make([]metrics.Metric, len(rows))
	for i, row := range rows {
		currRes := result[i]
		src, err := enum.ParseMetrics(row.Source)
		if err != nil {
			currRes.Source = currRes.Source
		}
		currRes = metrics.Metric{
			Timestamp: row.Timestamp,
			Source:    string(*src),
			Type:      row.Type,
			Value:     row.Value,
			Severity:  row.Severity,
			Message:   row.Message,
		}
	}
	return result
}

func (s *SQLiteStore) Close() error {
	return s.db.Close()
}
