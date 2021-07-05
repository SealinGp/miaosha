package client

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"miaosha/pkg/bootstrap"
	conf "miaosha/pkg/config"
	"miaosha/pkg/discover"
	"miaosha/pkg/loadbalance"


	"github.com/afex/hystrix-go/hystrix"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	zipkinnot "github.com/openzipkin-contrib/zipkin-go-opentracing"
	zipkinGo "github.com/openzipkin/zipkin-go"
	zipkinReporter "github.com/openzipkin/zipkin-go/reporter"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	grpcOri "google.golang.org/grpc"
)
var (
	ErrRPCService = errors.New("no rpc service")
)

//general grpc client manager interface
type ClientManager interface {
	DecoratorInvoke(path string, hystrixName string, tracer opentracing.Tracer,
		ctx context.Context, inputVal interface{}, outVal interface{}) (err error)
}


//default client manager interface
type DefaultClientManager struct {
	serviceName     string
	logger          *log.Logger
	discoveryClient discover.DiscoveryClient
	loadBalance     loadbalance.LoadBalance
	after          []InvokerAfterFunc
	before         []InvokerBeforeFunc
}
type InvokerAfterFunc func()(err error)
type InvokerBeforeFunc func()(err error)

func (manager *DefaultClientManager)DecoratorInvoke(path string, hystrixName string, tracer opentracing.Tracer,
	ctx context.Context, inputVal interface{}, outVal interface{}) error  {

	//调用之前执行的函数
	for _, fn := range manager.before {
		if err := fn(); err != nil {
			return err
		}
	}

	//调用grpc具体逻辑
	run := func() error {
		//服务发现与负载均衡
		instances     := manager.discoveryClient.DiscoverServices(manager.serviceName,manager.logger)
		instance, err := manager.loadBalance.SelectService(instances)
		if err != nil {
			return err
		}
		if instance != nil && instance.GrpcPort <= 0 {
			return ErrRPCService
		}

		//grpc链路追踪zipkin上报设置
		address         := fmt.Sprintf("%s:%s",bootstrap.HttpConfig.Host,bootstrap.HttpConfig.Port)
		tracer,reporter := genTracer(tracer,bootstrap.DiscoverConfig.ServiceName,address)

		//todo上报关闭
		if reporter != nil {
			defer reporter.Close()
		}

		//grpc调用选项,超时返回,链路追踪,熔断保护
		options := []grpcOri.DialOption{
			grpcOri.WithInsecure(),
			grpcOri.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(tracer,otgrpc.LogPayloads())),
		}
		ctxTimout,cf  := context.WithTimeout(ctx,time.Second)
		defer cf()

		//开始连接->执行调用
		conn,err := grpcOri.DialContext(ctxTimout,fmt.Sprintf("%s:%d",instance.Host,instance.GrpcPort),options...)
		if err != nil {
			return err
		}
		if err = conn.Invoke(ctx,path,inputVal,outVal);err != nil {
			return err
		}
		return nil
	}
	//调用失败降级处理
	callBack := func(err error) error {
		return err
	}

	//断路器保护
	if err := hystrix.Do(hystrixName, run, callBack);err != nil {
		return err
	}

	//调用后执行函数
	for _, fn := range manager.after  {
		if err := fn();err != nil {
			return err
		}
	}
	return nil
}

func genTracer(tracer opentracing.Tracer,svcName,address string) (opentracing.Tracer,zipkinReporter.Reporter) {
	if tracer != nil {
		return tracer,nil
	}

	zipkinUrl := fmt.Sprintf("http://%s:%s%s",conf.TraceConfig.Host,conf.TraceConfig.Port,conf.TraceConfig.Url)
	reporter  := zipkinhttp.NewReporter(zipkinUrl)     //上报地址
	ep,err    := zipkinGo.NewEndpoint(svcName,address) //上传
	if err != nil {
		log.Fatal("endpoint error",err.Error())
	}

	tracer1,err := zipkinGo.NewTracer(reporter,zipkinGo.WithLocalEndpoint(ep))
	if err != nil {
		log.Fatal("new tracer error:",err.Error())
	}
	tracer = zipkinnot.Wrap(tracer1)
	return tracer,reporter
}