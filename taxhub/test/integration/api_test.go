//go:build integration

package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/elixir-metatrack/taxhub/internal/handler"
	"github.com/elixir-metatrack/taxhub/internal/model"
	"github.com/elixir-metatrack/taxhub/internal/service"
	"github.com/prometheus/client_golang/prometheus"
)

func newTestServer(t *testing.T) *httptest.Server {
	t.Helper()

	csvPath := os.Getenv("CSV_PATH")
	if csvPath == "" {
		t.Skip("CSV_PATH not set, skipping integration test")
	}

	tx := service.NewTaxonomy()
	if err := tx.LoadCSV(csvPath); err != nil {
		t.Fatalf("failed to load csv: %v", err)
	}

	h := handler.New(tx, prometheus.NewRegistry())

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", h.Healthz)
	mux.HandleFunc("GET /taxon/{tax_id}", h.GetTaxon)

	return httptest.NewServer(mux)
}

func TestIntegration_Healthz(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()

	resp, err := http.Get(fmt.Sprintf("%s/healthz", srv.URL))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	var body model.HealthResponse
	json.NewDecoder(resp.Body).Decode(&body)
	if body.Status != "ok" {
		t.Errorf("expected status ok, got %q", body.Status)
	}
	if body.Count == 0 {
		t.Error("expected non-zero count")
	}
}

func TestIntegration_GetTaxon_found(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()

	resp, err := http.Get(fmt.Sprintf("%s/taxon/9606", srv.URL))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	var body model.TaxonResponse
	json.NewDecoder(resp.Body).Decode(&body)
	if body.Name != "Homo sapiens" {
		t.Errorf("expected Homo sapiens, got %q", body.Name)
	}
}

func TestIntegration_GetTaxon_notFound(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()

	resp, err := http.Get(fmt.Sprintf("%s/taxon/0", srv.URL))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404, got %d", resp.StatusCode)
	}
}
