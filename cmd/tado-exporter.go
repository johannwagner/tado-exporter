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
	exporterUsername, hasUsername := os.LookupEnv("EXPORTER_USERNAME")
	exporterPassword, hasPassword := os.LookupEnv("EXPORTER_PASSWORD")

	var tadoClient tado.TadoClient

	if hasUsername && hasPassword {
		tadoClient = tado.NewTadoWebClient(
			exporterUsername,
			exporterPassword,
		)
	} else {
		tadoClient = tado.NewTadoAPIClient()
	}

	err := tadoClient.Authorize()
	if err != nil {
		log.Fatal(err)
	}

	registry := prometheus.NewRegistry()

	tadoCollector := tadoprometheus.NewTadoCollector(tadoClient)
	registry.MustRegister(tadoCollector)
	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	log.Println("Listening on :9898")
	log.Fatal(http.ListenAndServe(":9898", nil))
}
