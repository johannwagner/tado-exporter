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
	TADO_USERNAME := os.Getenv("TADO_USERNAME")
	TADO_PASSWORD := os.Getenv("TADO_PASSWORD")

	tadoClient := tado.NewTadoClient(
		TADO_USERNAME,
		TADO_PASSWORD,
	)

	registry := prometheus.NewRegistry()

	tadoCollector := tadoprometheus.NewTadoCollector(tadoClient)
	registry.MustRegister(tadoCollector)
	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	log.Fatal(http.ListenAndServe(":9898", nil))
}
