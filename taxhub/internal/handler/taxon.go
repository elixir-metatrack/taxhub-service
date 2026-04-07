package handler

import (
	"encoding/json"
	"net/http"

	"github.com/elixir-metatrack/taxhub/internal/model"
	"github.com/elixir-metatrack/taxhub/internal/service"
	"github.com/prometheus/client_golang/prometheus"
)

type Handler struct {
	taxonomy     *service.Taxonomy
	requestCount *prometheus.CounterVec
}

func New(taxonomy *service.Taxonomy, registry *prometheus.Registry) *Handler {
	requestCount := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "taxhub_requests_total",
		Help: "Total number of requests per endpoint and status",
	}, []string{"endpoint", "status"})

	registry.MustRegister(requestCount)

	return &Handler{taxonomy: taxonomy, requestCount: requestCount}
}

// Healthz godoc
// @Summary     Health check
// @Tags        health
// @Produce     json
// @Success     200 {object} model.HealthResponse
// @Router      /healthz [get]
func (h *Handler) Healthz(w http.ResponseWriter, r *http.Request) {
	h.requestCount.WithLabelValues("healthz", "200").Inc()
	writeJSON(w, http.StatusOK, model.HealthResponse{
		Status:        "ok",
		TaxonCount: h.taxonomy.Count(),
	})
}

// GetTaxon godoc
// @Summary     Get taxon by ID
// @Tags        taxon
// @Produce     json
// @Param       tax_id path string true "NCBI Taxonomy ID"
// @Success     200 {object} model.TaxonResponse
// @Failure     400 {object} model.ErrorResponse
// @Failure     404 {object} model.ErrorResponse
// @Router      /taxon/{tax_id} [get]
func (h *Handler) GetTaxon(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("tax_id")
	if id == "" {
		h.requestCount.WithLabelValues("taxon", "400").Inc()
		writeJSON(w, http.StatusBadRequest, model.ErrorResponse{Error: "missing tax_id"})
		return
	}

	name, ok := h.taxonomy.Get(id)
	if !ok {
		h.requestCount.WithLabelValues("taxon", "404").Inc()
		writeJSON(w, http.StatusNotFound, model.ErrorResponse{Error: "tax_id not found"})
		return
	}

	h.requestCount.WithLabelValues("taxon", "200").Inc()
	writeJSON(w, http.StatusOK, model.TaxonResponse{TaxID: id, Name: name})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
