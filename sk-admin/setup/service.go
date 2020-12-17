package setup

import (
	"context"
	"flag"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	kitzipkin "github.com/go-kit/kit/tracing/zipkin"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"golang.org/x/time/rate"
	"log"
	"miaosha/pkg/config"
	"miaosha/sk-admin/endpoint"
	"miaosha/sk-admin/plugins"
	"miaosha/sk-admin/service"
	"time"
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

	ratebucket := rate.NewLimiter(rate.Every(time.Second), 100)

	getActivityEnd := endpoint.MakeGetActivityEndpoint(activityService)
	getActivityEnd = plugins.NewTokenBucketLimitterWithBuildIn(ratebucket)(getActivityEnd)
	getActivityEnd = kitzipkin.TraceEndpoint(config.ZipkinTracer, "get-activity")(getActivityEnd)

	createActivityEnd := endpoint.MakeCreateActivityEndpoint(activityService)
	createActivityEnd = plugins.NewTokenBucketLimitterWithBuildIn(ratebucket)(createActivityEnd)
	createActivityEnd = kitzipkin.TraceEndpoint(config.ZipkinTracer, "create-activity")(createActivityEnd)

	createProductEnd := endpoint.MakeCreateProductEndpoint(productService)
	createProductEnd = plugins.NewTokenBucketLimitterWithBuildIn(ratebucket)(createProductEnd)
	createProductEnd = kitzipkin.TraceEndpoint(config.ZipkinTracer, "create-product")(createProductEnd)

	healthCheckEnd := endpoint.MakeHealthCheckEndpoint(skAdminService)
	healthCheckEnd = kitzipkin.TraceEndpoint(config.ZipkinTracer, "health-endpoint")(healthCheckEnd)

	endpts := endpoint.SkAdminEndpoints{
		GetActivityEndpoint:    getActivityEnd,
		CreateActivityEndpoint: createActivityEnd,
		CreateProductEndpoint:  createProductEnd,
		HealthCheckEndpoint:    healthCheckEnd,
	}
	ctx := context.Background()

}
