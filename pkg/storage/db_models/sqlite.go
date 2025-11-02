package db_models

import (
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type SQLiteStore struct {
	Db *sqlx.DB
}

type MetricRow struct {
	ID        int       `db:"id"`
	Timestamp time.Time `db:"timestamp"`
	Source    string    `db:"source"`
	Type      string    `db:"type"`
	Value     float64   `db:"value"`
	Severity  string    `db:"severity"`
	Message   string    `db:"message"`
}
