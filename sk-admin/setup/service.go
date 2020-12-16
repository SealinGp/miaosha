package setup

import (
	"flag"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"log"
	"miaosha/pkg/config"
	"miaosha/sk-admin/plugins"
	"miaosha/sk-admin/service"
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

	var (
		activityService service.ActivityService
		productService  service.ProductService
		skAdminService  service.Service
	)

	skAdminService = service.SkAdminService{}
	activityService = &service.ActivityServiceImpl{}
	productService = &service.ProductServiceImpl{}

	skAdminService = plugins.SkAdminLoggingMiddleware(config.Logger)(skAdminService)
	skAdminService = plugins.SkAdminMetrics(requestCount, requestLatency)(skAdminService)

	activityService = plugins.ActivityLoggingMiddleware(config.Logger)(activityService)
	activityService = plugins.ActivityMetrics(requestCount, requestLatency)(activityService)

	productService = plugins.ProductLoggingMiddleware(config.Logger)(productService)
	productService = plugins.ProductMetrics(requestCount, requestLatency)(productService)

}
