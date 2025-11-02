package responses

import (
	"github.com/raulbattistini/private-empty/pkg/server/models/enum"
)

type HandleHealthRes struct {
	Status enum.ApiHealthStatus `json:"status"`
	Time   string               `json:"time"`
}
