package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/elixir-metatrack/taxhub/internal/model"
	"github.com/elixir-metatrack/taxhub/internal/service"
	"github.com/prometheus/client_golang/prometheus"
)

func newTestHandler(t *testing.T) *Handler {
	t.Helper()
	f, err := os.CreateTemp("", "taxonomy*.csv")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString("tax_id,name\n9606,Homo sapiens\n")
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })

	tx := service.NewTaxonomy()
	if err := tx.LoadCSV(f.Name()); err != nil {
		t.Fatal(err)
	}
	return New(tx, prometheus.NewRegistry())
}

func TestHealthz(t *testing.T) {
	h := newTestHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()
	h.Healthz(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	var resp model.HealthResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.Status != "ok" {
		t.Errorf("expected status ok, got %q", resp.Status)
	}
}

func TestGetTaxon_found(t *testing.T) {
	h := newTestHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/taxon/9606", nil)
	req.SetPathValue("tax_id", "9606")
	w := httptest.NewRecorder()
	h.GetTaxon(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	var resp model.TaxonResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.Name != "Homo sapiens" {
		t.Errorf("expected Homo sapiens, got %q", resp.Name)
	}
}

func TestGetTaxon_notFound(t *testing.T) {
	h := newTestHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/taxon/0", nil)
	req.SetPathValue("tax_id", "0")
	w := httptest.NewRecorder()
	h.GetTaxon(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestGetTaxon_missingID(t *testing.T) {
	h := newTestHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/taxon/", nil)
	w := httptest.NewRecorder()
	h.GetTaxon(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}
