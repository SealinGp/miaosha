package config

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"github.com/spf13/viper"
	"miaosha/pkg/bootstrap"
	conf "miaosha/pkg/config"
	"os"
)

const (
	kConfigType = "CONFIG_TYPE"
)
var ZipkinTracer *zipkin.Tracer
var Logger log.Logger

func init()  {
	Logger = log.NewLogfmtLogger(os.Stderr)
	Logger = log.With(Logger,"ts",log.DefaultTimestamp)
	Logger = log.With(Logger,"caller",log.DefaultCaller)
	viper.AutomaticEnv()
	initDefault()

	if err := conf.LoadRemoteConfig(); err != nil {
		Logger.Log("Fail to load remote config",err)
	}
	if err := conf.SubParse("mysql",&conf.MysqlConfig);err != nil {
		Logger.Log("Fail to parse mysql",err)
	}
	if err := conf.SubParse("trace",&conf.TraceConfig); err != nil {
		Logger.Log("Fail to parse trace",err)
	}

	zipkinUrl := fmt.Sprintf("http://%s:%s%s",conf.TraceConfig.Host,conf.TraceConfig.Port,conf.TraceConfig.Url)
	Logger.Log("zipkin url",zipkinUrl)
	initTracer(zipkinUrl)
}
func initDefault()  {
	viper.SetDefault(kConfigType,"yaml")
}

func initTracer(zipkinUrl string)  {
	var (
		err error
		useNoopTracer = zipkinUrl == ""
		reporter = zipkinhttp.NewReporter(zipkinUrl)
	)
	zEP, _ := zipkin.NewEndpoint(bootstrap.DiscoverConfig.ServiceName,bootstrap.HttpConfig.Port)
	ZipkinTracer,err = zipkin.NewTracer(
		reporter,zipkin.WithLocalEndpoint(zEP),zipkin.WithNoopTracer(useNoopTracer),
	)
	if err != nil {
		Logger.Log("tracer init error",err)
		os.Exit(1)
	}
	if !useNoopTracer {
		Logger.Log("tracker","Zipkin","type","Native","URL",zipkinUrl)
	}
}