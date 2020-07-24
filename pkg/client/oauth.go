package client

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"miaosha/pb"
	"miaosha/pkg/discover"
	"miaosha/pkg/loadbalance"
)

type OAuthClient interface {
	CheckToken(ctx context.Context, tracer opentracing.Tracer,request *pb.CheckTokenRequest) (*pb.CheckTokenResponse, error)
}

type OAuthClientImpl struct {
	manager ClientManager
	serviceName string
	loadBalance loadbalance.LoadBalance
	tracer opentracing.Tracer
}
func (impl *OAuthClientImpl)CheckToken(ctx context.Context, tracer opentracing.Tracer,request *pb.CheckTokenRequest) (*pb.CheckTokenResponse, error)  {
	response := new(pb.CheckTokenResponse)
	err      := impl.manager.DecoratorInvoke("/pb.OauthService/CheckToken","token_check",tracer,ctx,request,response)
	if err != nil {
		return nil,err
	}
	return response,nil
}

func NewOauthClient(serviceName string, lb loadbalance.LoadBalance,tracer opentracing.Tracer) (OAuthClient, error) {
	if serviceName == "" {
		serviceName = "oauth"
	}
	if lb == nil {
		lb = discover.LoadBalance
	}
	return &OAuthClientImpl{
		manager:&DefaultClientManager{
			serviceName:serviceName,
			loadBalance:lb,
			discoveryClient:discover.ConsulService,
			logger:discover.Logger,
		},
		serviceName:serviceName,
		loadBalance:lb,
		tracer:tracer,
	},nil
}