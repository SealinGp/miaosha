package setup

import (
	"flag"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"log"
)

func InitServer(host, servicePort string) {
	log.Printf("[I] port:%s", servicePort)
	flag.Parse()

	errChan := make(chan error)
	fieldKeys := []string{"method"}

	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "aoho",
		Subsystem: "user_service",
		Name:      "request_count",
		Help:      "Number of requests received",
	}, fieldKeys)
	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "aoho",
		Subsystem: "user_service",
		Name:      "request_latency",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)

	var ()
}
