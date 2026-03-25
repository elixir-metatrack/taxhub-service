// @title       TaxHub API
// @version     1.0
// @description Maps NCBI Taxonomy IDs to scientific names.
// @host        localhost:8080
// @BasePath    /
package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/elixir-metatrack/taxhub/docs"
	"github.com/elixir-metatrack/taxhub/internal/handler"
	"github.com/elixir-metatrack/taxhub/internal/service"
)

func main() {
	addr    := getenv("ADDR", ":8080")
	csvPath := getenv("CSV_PATH", "scientific_names.csv")

	taxonomy := service.NewTaxonomy()
	if err := taxonomy.LoadCSV(csvPath); err != nil {
		log.Fatalf("failed to load csv: %v", err)
	}

	registry := prometheus.NewRegistry()
	h := handler.New(taxonomy, registry)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", h.Healthz)
	mux.HandleFunc("GET /taxon/{tax_id}", h.GetTaxon)
	mux.Handle("GET /metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	log.Printf("TaxHub listening on %s", addr)
	log.Printf("Loaded %d entries from %s", taxonomy.Count(), csvPath)
	log.Printf("Swagger UI: http://localhost%s/swagger/", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}

func getenv(key, fallback string) string {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return fallback
	}
	return v
}
