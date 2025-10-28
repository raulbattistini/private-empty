package exporter

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/raulbattistini/private-empty/internal/event"
)

type HTTPExporter struct {
	endpoint string
}

func NewHTTPExporter(endpoint string) *HTTPExporter {
	return &HTTPExporter{endpoint: endpoint}
}

func (e *HTTPExporter) Export(ctx context.Context, ev event.Event) error {
	body, err := json.Marshal(ev)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", e.endpoint, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	_, err = http.DefaultClient.Do(req)
	// simplifying the most -- not meant for critical usage
	return err
}
