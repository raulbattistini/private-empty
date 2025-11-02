package storage

import (
	"github.com/raulbattistini/private-empty/pkg/server/models/enum"
	"github.com/raulbattistini/private-empty/pkg/storage/db_models"
)

type HardStorage interface {
	ToStorage() *db_models.SQLiteStore
	Receive(m db_models.StorageMetric)
	GetLatest() []db_models.StorageMetric
	GetHistory(source enum.RuleSource, limit int) []db_models.StorageMetric
	RowsToMetrics(rows []db_models.MetricRow) []db_models.StorageMetric
	Close() error
}
