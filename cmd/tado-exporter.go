package main

import (
	tadoprometheus "github.com/johannwagner/tado-exporter-go/internal/prometheus"
	"github.com/johannwagner/tado-exporter-go/internal/tado"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os"
	"sync"
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

	registry := prometheus.NewRegistry()

	tadoCollector := tadoprometheus.NewTadoCollector(tadoClient)
	registry.MustRegister(tadoCollector)
	log.Println("Listening on :9898")

	wg := sync.WaitGroup{}

	wg.Add(2)

	go func() {
		defer wg.Done()
		err := tadoClient.Authorize()
		if err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		defer wg.Done()
		http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
		log.Fatal(http.ListenAndServe(":9898", nil))
	}()

	wg.Wait()
}
