package discover

import (
	"log"
	"miaosha/pkg/common"
)

type DiscoveryClient interface {
	Register(instanceId,svcHost,healthCheckUrl,svcPort,svcName string,weight int,metaMap map[string]string,tags []string,logger *log.Logger) bool
	DeRegister(instanceId string,logger * log.Logger) bool
	DiscoverServices(svcName string,logger *log.Logger) []*common.ServiceInstance
}