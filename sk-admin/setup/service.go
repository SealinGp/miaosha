package setup

import (
	"context"
	"flag"
	"fmt"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	kitzipkin "github.com/go-kit/kit/tracing/zipkin"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"golang.org/x/time/rate"
	"log"
	"miaosha/pkg/config"
	"miaosha/pkg/discover"
	"miaosha/sk-admin/endpoint"
	"miaosha/sk-admin/plugins"
	"miaosha/sk-admin/service"
	"miaosha/sk-admin/transport"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

	GetProductEnd := endpoint.MakeGetProductEndpoint(productService)

	healthCheckEnd := endpoint.MakeHealthCheckEndpoint(skAdminService)
	healthCheckEnd = kitzipkin.TraceEndpoint(config.ZipkinTracer, "health-endpoint")(healthCheckEnd)

	endpts := endpoint.SkAdminEndpoints{
		GetActivityEndpoint:    getActivityEnd,
		CreateActivityEndpoint: createActivityEnd,
		CreateProductEndpoint:  createProductEnd,
		GetProductEndpoint:     GetProductEnd,
		HealthCheckEndpoint:    healthCheckEnd,
	}
	ctx := context.Background()
	r := transport.MakeHttpHandler(ctx, endpts, config.ZipkinTracer, config.Logger)

	go func() {
		fmt.Println("http server start at port:" + servicePort)
		discover.Register()
		handler := r
		errChan <- http.ListenAndServe(":"+servicePort, handler)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	error := <-errChan
	//服务退出取消注册
	discover.Deregister()
	fmt.Println(error)
}
