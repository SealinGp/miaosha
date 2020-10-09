package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"miaosha/gateway/config"
	"miaosha/gateway/route"
	"miaosha/pkg/bootstrap"
	register "miaosha/pkg/discover"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/go-kit/kit/log"
	"github.com/openzipkin/zipkin-go"
	zipkinhttpsvr "github.com/openzipkin/zipkin-go/middleware/http"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
)

func main()  {
	zipkinURL := config.GetZipkinUrl()
	fmt.Println(zipkinURL)
	return

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger,"ts",log.TimestampFormat(time.Now,time.RFC3339))
		logger = log.With(logger,"caller",log.DefaultCaller)
	}

	var zipkinTracer *zipkin.Tracer
	{
		var (
			useNoopTracer = zipkinURL == ""
			reporter      = zipkinhttp.NewReporter(zipkinURL)
		)
		defer reporter.Close()
		endpoint,err := zipkin.NewEndpoint(bootstrap.DiscoverConfig.ServiceName,net.JoinHostPort(bootstrap.HttpConfig.Host,bootstrap.HttpConfig.Port))
		if err != nil {
			logger.Log("err",err)
			os.Exit(1)
		}


		zipkinTracer,err = zipkin.NewTracer(
			reporter,
			zipkin.WithLocalEndpoint(endpoint),
		)
		if err != nil {
			logger.Log("err,",err)
			os.Exit(1)
		}

		if !useNoopTracer {
			logger.Log("tracer","Zipkin","type","Native","URL","zipkinURL")
		}
	}

	tags := map[string]string{
		"component":"gateway_server",
	}

	hystrixRouter := route.Routes(zipkinTracer,"Circuit Breaker:Service unavailable",logger)
	handler       := zipkinhttpsvr.NewServerMiddleware(
		zipkinTracer,
		zipkinhttpsvr.SpanName(bootstrap.DiscoverConfig.ServiceName),
		zipkinhttpsvr.TagResponseSize(true),
		zipkinhttpsvr.ServerTags(tags),
	)(hystrixRouter)

	errc := make(chan error)

	//启用hystrix实时监控,端口为9090
	hystrixStreamHandler := hystrix.NewStreamHandler()
	hystrixStreamHandler.Start()
	go func() {
		errc <- http.ListenAndServe(net.JoinHostPort(config.HystrixConfig.Host,config.HystrixConfig.Port),hystrixStreamHandler)
	}()
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c,syscall.SIGINT,syscall.SIGTERM)
		errc <- fmt.Errorf("%s",<-c)
	}()

	//网关服务开始监听
	go func() {
		addr := net.JoinHostPort(bootstrap.HttpConfig.Host,bootstrap.HttpConfig.Port)
		logger.Log("transport","HTTP","addr",addr)
		register.Register()
		errc <- http.ListenAndServe(addr,handler)
	}()


	err := <-errc
	register.Deregister()
	logger.Log("exit",err)
}