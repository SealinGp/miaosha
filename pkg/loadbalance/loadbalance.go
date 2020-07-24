package loadbalance

import (
	"errors"
	"math/rand"
	"miaosha/pkg/common"
)

// 负载均衡器
type LoadBalance interface {
	SelectService(service []*common.ServiceInstance) (*common.ServiceInstance, error)
}

//随机负载均衡器(适用于每台机器配置相同的业务场景)
type RandomLoadBalance struct {}
func (loadBalance *RandomLoadBalance)SelectService(services []*common.ServiceInstance) (*common.ServiceInstance,error) {
	if services == nil || len(services) == 0 {
		return nil,errors.New("service instances are not exists")
	}
	return services[rand.Intn(len(services))],nil
}

//权重负载均衡(适用于每台机器的配置不同的业务场景)
type WeightRoundRobinLoadBalance struct {}
func (loadBalance *WeightRoundRobinLoadBalance)SelectService(services []*common.ServiceInstance) (*common.ServiceInstance,error) {
	if services == nil || len(services) == 0 {
		return nil,errors.New("service instances are not exists")
	}

	var best *common.ServiceInstance
	total := 0
	for  i := range services {
		w := services[i]
		if w == nil {
			continue
		}
		w.CurWeight += w.Weight
		total       += w.Weight
		if best == nil || w.CurWeight > best.CurWeight {
			best = w
		}
	}
	if best == nil {
		return nil,errors.New("unknown error: can not select instance")
	}
	return best,nil
}