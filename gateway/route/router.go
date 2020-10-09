package route

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"
	"sync"

	"miaosha/gateway/config"
	"miaosha/pb"
	"miaosha/pkg/client"
	"miaosha/pkg/discover"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/go-kit/kit/log"
	"github.com/openzipkin/zipkin-go"
	loadbalance2 "miaosha/pkg/loadbalance"
	zipkinhttpsvr "github.com/openzipkin/zipkin-go/middleware/http"
)

type HystrixRouter struct {
	svcMap      *sync.Map
	logger      log.Logger
	fallbackMsg string
	tracer      *zipkin.Tracer
	loadbalance loadbalance2.LoadBalance
}

func Routes(zipkinTracer *zipkin.Tracer,fbMsg string, logger log.Logger) http.Handler {
	return HystrixRouter{
		svcMap:      &sync.Map{},
		logger:      logger,
		fallbackMsg: fbMsg,
		tracer:      zipkinTracer,
		loadbalance: &loadbalance2.RandomLoadBalance{},
	}
}
func preFilter(r *http.Request) bool {
	reqPath := r.URL.Path
	if reqPath == "" {
		return false
	}

	if config.Match(reqPath) {
		return true
	}

	authToken := r.Header.Get("Authorization")
	if authToken == "" {
		return false
	}

	oauthClient, _ := client.NewOauthClient("oauth",nil,nil)
	resp,remoteErr := oauthClient.CheckToken(context.Background(),nil,&pb.CheckTokenRequest{
		Token:                authToken,
	})

	return resp != nil && remoteErr == nil
}

func postFilter()  {
	//todo for custom filter
}

func (router HystrixRouter)ServeHTTP(w http.ResponseWriter,r *http.Request)  {
	reqPath := r.URL.Path
	router.logger.Log("reqPath",reqPath)

	if reqPath == "/health" {
		w.WriteHeader(http.StatusOK)
		return
	}

	var err error
	if !preFilter(r) {
		err = errors.New("illegal request")
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(err.Error()))
		return
	}

	pathArray   := strings.Split(reqPath,"/")
	serviceName := pathArray[1]

	if _,ok := router.svcMap.Load(serviceName);!ok {
		hystrix.ConfigureCommand(serviceName,hystrix.CommandConfig{Timeout:1000})
		router.svcMap.Store(serviceName,serviceName)
	}

	err = hystrix.Do(serviceName, func() error {
		//调用consul api 查询serviceName
		serviceInstance, err := discover.DiscoveryService(serviceName)
		if err != nil {
			return err
		}

		//
		director := func(req *http.Request) {
			destPath := strings.Join(pathArray[2:],"/")

			//随机选择一个服务实例
			router.logger.Log("service id",serviceInstance.Host,serviceInstance.Port)

			//设置代理服务地址信息
			req.URL.Scheme = "http"
			req.URL.Host   = fmt.Sprintf("%s/%d",serviceInstance.Host,serviceInstance.Port)
			req.URL.Path   = "/" + destPath
		}
		var proxyError error = nil
		//为反向代理增加追踪逻辑,使用如下RoundTrip代替默认Transport
		roundTrip, _ := zipkinhttpsvr.NewTransport(router.tracer,zipkinhttpsvr.TransportTrace(true))
		//反向代理失败时错误处理
		errorHandler := func(ew http.ResponseWriter,er *http.Request,err error) {
			proxyError = err
		}
		proxy := &httputil.ReverseProxy{
			Director:       director,
			Transport:      roundTrip,
			ErrorHandler:   errorHandler,
		}

		proxy.ServeHTTP(w,r)

		return proxyError
	}, func(err error) error {
		router.logger.Log("fallback error description",err.Error())
		return errors.New(router.fallbackMsg)
	})
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
	}
}