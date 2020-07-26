package client

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"miaosha/pb"
	"miaosha/pkg/discover"
	"miaosha/pkg/loadbalance"
)

type UserClient interface {
	CheckUser(ctx context.Context, tracer opentracing.Tracer, request *pb.UserRequest)(*pb.UserResponse,error)
}

type UserClientImpl struct {
	/**
	 * 可以配置负载均衡策略,重试等机制,也可以配置InvokeAfter 和 invokerBefore
	 */
	manager ClientManager
	serviceName string
	loadBalance loadbalance.LoadBalance
	tracer opentracing.Tracer
}

func (impl *UserClientImpl)CheckUser(ctx context.Context, tracer opentracing.Tracer, request *pb.UserRequest) (*pb.UserResponse,error) {
	response := new(pb.UserResponse)
	if err := impl.manager.DecoratorInvoke("/pb.UserService/Check","user_check",tracer,ctx,request,response); err != nil {
		return nil,err
	}
	return response,nil
}

func NewUserClient(serviceName string,lb loadbalance.LoadBalance,tracer opentracing.Tracer) (UserClient,error) {
	if serviceName == "" {
		serviceName = "user"
	}
	if lb == nil {
		lb = discover.LoadBalance
	}
	return &UserClientImpl{
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