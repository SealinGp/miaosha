package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	localconfig "miaosha/oauth-service/config"
	"miaosha/oauth-service/endpoint"
	"miaosha/oauth-service/plugins"
	"miaosha/oauth-service/service"
	"miaosha/oauth-service/transport"
	"miaosha/pb"
	"miaosha/pkg/bootstrap"
	conf "miaosha/pkg/config"
	register "miaosha/pkg/discover"
	"miaosha/pkg/mysql"

	kitzipkin "github.com/go-kit/kit/tracing/zipkin"
	"github.com/openzipkin/zipkin-go/propagation/b3"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func main()  {
	return
	var (
		servicePort = flag.String("service.port",bootstrap.HttpConfig.Port,"service port")
		grpcAddr    = flag.String("grpc",bootstrap.RpcConfig.Port,"grpc listen address")
	)
	flag.Parse()

	ctx        := context.Background()
	errChan    := make(chan error)
	ratebucket := rate.NewLimiter(rate.Every(time.Second*1),100)

	var (
		tokenService service.TokenService
		tokenGranter service.TokenGranter

		userDetailsService service.UserDetailsService
		clientDetailsService service.ClientDetailsService
		srv service.Service
	)

	//token services
	tokenService         = service.NewDefaultTokenService("secret")
	userDetailsService   = service.NewRemoteUserService()
	tokenGranter         = service.NewComposeTokenGranter(map[string]service.TokenGranter{
		"password"     : service.NewUsernamePasswordTokenGranter("password",userDetailsService,tokenService),
		"refresh_token": service.NewRefreshGranter("refresh_token",tokenService),
	})
	clientDetailsService = service.NewMysqlClientDetailsService()
	srv                  = service.NewCommentService()

	//POST http://url?gran_type=password
	//header: Basic base64(clientId:clientSecret)
	//body: {username:xxx,password:xxx}
	//resp: endpoint.TokenResponse
	tokenEndpoint := endpoint.MakeTokenEndpoint(tokenGranter,clientDetailsService)
	tokenEndpoint  = endpoint.MakeClientAuthorizationMiddleware(localconfig.Logger)(tokenEndpoint)           //认证
	tokenEndpoint  = plugins.NewTokenBucketLimitterWithBuildIn(ratebucket)(tokenEndpoint)                    //限流
	tokenEndpoint  = kitzipkin.TraceEndpoint(localconfig.ZipkinTracer,"token-endpoint")(tokenEndpoint) //日志追踪

	//token解密|检查
	//http
	checkTokenEndpoint := endpoint.MakeCheckTokenEndpoint(tokenService)
	checkTokenEndpoint = endpoint.MakeClientAuthorizationMiddleware(localconfig.Logger)(checkTokenEndpoint)
	checkTokenEndpoint = plugins.NewTokenBucketLimitterWithBuildIn(ratebucket)(checkTokenEndpoint)
	checkTokenEndpoint = kitzipkin.TraceEndpoint(localconfig.ZipkinTracer,"check-endpoint")(checkTokenEndpoint)
	//grpc
	gRPCCheckTokenEndpoint := endpoint.MakeCheckTokenEndpoint(tokenService)
	gRPCCheckTokenEndpoint = plugins.NewTokenBucketLimitterWithBuildIn(ratebucket)(gRPCCheckTokenEndpoint)
	gRPCCheckTokenEndpoint = kitzipkin.TraceEndpoint(localconfig.ZipkinTracer,"grpc-check-endpoint")(gRPCCheckTokenEndpoint)

	healthEndpoint := endpoint.MakeHealthCheckEndpoint(srv)
	healthEndpoint = kitzipkin.TraceEndpoint(localconfig.ZipkinTracer,"health-endpoint")(healthEndpoint)

	endpts := endpoint.OAuth2Endpoints{
		TokenEndpoint:          tokenEndpoint,
		CheckTokenEndpoint:     checkTokenEndpoint,
		GRPCCheckTokenEndpoint: gRPCCheckTokenEndpoint,
		HealthCheckEndpoint:    healthEndpoint,
	}

	r := transport.MakeHttpHandler(ctx,endpts,tokenService,clientDetailsService,localconfig.ZipkinTracer,localconfig.Logger)

	//http
	go func() {
		fmt.Println("Http Server start at port:" + *servicePort)
		mysql.InitMysql(conf.MysqlConfig.Host,conf.MysqlConfig.Port,conf.MysqlConfig.User,conf.MysqlConfig.Pwd,conf.MysqlConfig.Db)
		register.Register()
		handler := r
		errChan <- http.ListenAndServe(":"+*servicePort,handler)
	}()
	//grpc
	go func() {
		fmt.Println("grpc Server start at port:" + *grpcAddr)
		listener,err := net.Listen("tcp",":"+*grpcAddr)
		if err != nil {
			errChan <- err
			return
		}
		serverTracer := kitzipkin.GRPCServerTrace(localconfig.ZipkinTracer,kitzipkin.Name("grpc-transport"))
		tr           := localconfig.ZipkinTracer
		md           := metadata.MD{}
		parentSpan   := tr.StartSpan("test")
		err           = b3.InjectGRPC(&md)(parentSpan.Context())
		if err != nil {
			listener.Close()
			errChan <- err
			return
		}
		ctx        := metadata.NewIncomingContext(context.Background(),md)
		handler    := transport.NewGRPCServer(ctx,endpts,serverTracer)
		gRPCServer := grpc.NewServer()
		pb.RegisterOAuthServiceServer(gRPCServer,handler)
		errChan <- gRPCServer.Serve(listener)
	}()

	go func() {
		c := make(chan os.Signal,1)
		signal.Notify(c,syscall.SIGINT,syscall.SIGTERM)
		errChan <- fmt.Errorf("%s",<-c)
	}()

	err := <-errChan
	register.Deregister()
	fmt.Println(err)
}