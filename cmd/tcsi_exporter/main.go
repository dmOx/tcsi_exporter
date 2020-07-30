package main

import (
	"github.com/dmOx/tcsi_exporter"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	config := ConfigFromEnv()

	client, err := tcsi_exporter.NewClient(config.Collector.Token)
	if err != nil {
		log.Fatal(err)
	}

	collectorPortfolio, err := tcsi_exporter.NewPortfolioCollector(client)
	if err != nil {
		log.Fatal(err)
	}

	collectorCash, err := tcsi_exporter.NewCashCollector(client)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Register collectors")
	prometheus.MustRegister(collectorPortfolio)
	prometheus.MustRegister(collectorCash)

	http.Handle("/metrics", promhttp.Handler())

	log.Println("Listening on", config.Http.Addr)
	log.Fatal(http.ListenAndServe(config.Http.Addr, nil))
}
