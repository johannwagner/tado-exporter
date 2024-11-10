package main

import (
	tadoprometheus "github.com/johannwagner/tado-exporter-go/internal/prometheus"
	"github.com/johannwagner/tado-exporter-go/internal/tado"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os"
)

func main() {

	// Note: We may want to support a config file.
	EXPORTER_USERNAME := os.Getenv("EXPORTER_USERNAME")
	EXPORTER_PASSWORD := os.Getenv("EXPORTER_PASSWORD")

	tadoClient := tado.NewTadoClient(
		EXPORTER_USERNAME,
		EXPORTER_PASSWORD,
	)

	registry := prometheus.NewRegistry()

	tadoCollector := tadoprometheus.NewTadoCollector(tadoClient)
	registry.MustRegister(tadoCollector)
	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	log.Fatal(http.ListenAndServe(":9898", nil))
}
