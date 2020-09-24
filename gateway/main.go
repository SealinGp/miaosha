package main

import (
	"flag"
	"github.com/go-kit/kit/log"
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"os"
	"time"
)

func main()  {
	var (
		zipkinURL = flag.String("zipkin.url","http://127.0.0.1:9411/api/v2/spans","Zipkin server url")
	)
	flag.Parse()

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger,"ts",log.TimestampFormat(time.Now,time.RFC3339))
		logger = log.With(logger,"caller",log.DefaultCaller)
	}

	var zipkinTracer *zipkin.Tracer
	{
		var (
			err error
			useNoopTracer = *zipkinURL == ""
			reporter = zipkinhttp.NewReporter(*zipkinURL)
		)
		defer reporter.Close()

		zipkin.NewTracer()
	}
}