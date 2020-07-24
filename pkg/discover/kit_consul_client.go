package discover

import (
	"fmt"
	"github.com/go-kit/kit/sd/consul"
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
	"log"
	"miaosha/pkg/common"
	"strconv"
	"sync"
)

type DiscoveryClientInstance struct {
	Host string
	Port int
	config *api.Config
	client consul.Client
	mutex sync.Mutex
	instancesMap sync.Map
}
func NewDiscoveryClient(consulHost string, consulPort int) DiscoveryClient {
	cfg           := api.DefaultConfig()
	cfg.Address    = fmt.Sprintf("%s:%d",consulHost,consulPort)
	cfgClient,err := api.NewClient(cfg)
	if err != nil {
		return nil
	}

	return &DiscoveryClientInstance{
		Host:         consulHost,
		Port:         consulPort,
		config:       cfg,
		client:       consul.NewClient(cfgClient),
	}
}


func (consulClient *DiscoveryClientInstance)Register(instanceId,svcHost,healthCheckUrl,svcPort,svcName string,weight int,metaMap map[string]string,tags []string,logger *log.Logger) bool  {
	port,_ := strconv.Atoi(svcPort)
	cfg    := api.AgentServiceRegistration{
		ID:                instanceId,
		Name:              svcName,
		Address:           svcHost,
		Port:              port,
		Meta:              metaMap,
		Tags:              tags,
		Weights:          &api.AgentWeights{
			Passing: weight,
		},
		Check:&api.AgentServiceCheck{
			DeregisterCriticalServiceAfter: "30s",
			HTTP:                          fmt.Sprintf("http://%s:%s%s",svcHost,svcPort,healthCheckUrl),
			Interval:                       "15s",
		},
	}
	err := consulClient.client.Register(&cfg)
	if err != nil {
		logger.Println("Register Error:",err.Error())
		return false
	}
	logger.Println("Register Success!")
	return true
}
func (consulClient *DiscoveryClientInstance)DeRegister(instanceId string,logger * log.Logger) bool  {
	cfg := api.AgentServiceRegistration{
		ID:instanceId,
	}
	if err := consulClient.client.Deregister(&cfg);err != nil {
		logger.Println("Deregister Error:",err.Error())
		return false
	}
	logger.Println("Deregister success")
	return true
}
func (consulClient *DiscoveryClientInstance)DiscoverServices(svcName string,logger *log.Logger) []*common.ServiceInstance  {
	instanceList, ok := consulClient.instancesMap.Load(svcName)
	if ok {
		return instanceList.([]*common.ServiceInstance)
	}
	consulClient.mutex.Lock()
	defer consulClient.mutex.Unlock()
	instanceList,ok = consulClient.instancesMap.Load(svcName)
	if ok {
		return instanceList.([]*common.ServiceInstance)
	}

	//注册监控
	go func() {
		params := make(map[string]interface{})
		params["type"]    = "service"
		params["service"] = svcName
		plan,_ := watch.Parse(params)
		plan.Handler = func(u uint64, i interface{}) {
			if i == nil {
				return
			}
			v,ok := i.([]*api.ServiceEntry)
			if !ok {
				return
			}
			if len(v) == 0 {
				consulClient.instancesMap.Store(svcName,[]*common.ServiceInstance{})
			}
			var healthSvcs []*common.ServiceInstance
			for _, service := range v {
				if service.Checks.AggregatedStatus() == api.HealthPassing {
					healthSvcs = append(healthSvcs,newServiceInstance(service.Service))
				}
			}
			consulClient.instancesMap.Store(svcName,healthSvcs)
		}
		defer plan.Stop()
		_ = plan.Run(consulClient.config.Address)
	}()

	entries, _, err := consulClient.client.Service(svcName,"",false,nil)
	if err != nil {
		consulClient.instancesMap.Store(svcName,[]*common.ServiceInstance{})
		logger.Println("discovery svc error:",err.Error())
		return nil
	}
	instances := make([]*common.ServiceInstance,len(entries))
	for i := range instances  {
		instances[i] = newServiceInstance(entries[i].Service)
	}
	return instances
}

func newServiceInstance(agentSvc *api.AgentService) *common.ServiceInstance  {
	rpcPort := agentSvc.Port-1
	if agentSvc.Meta != nil {
		if rpcPortString, ok := agentSvc.Meta["rpcPort"];ok {
			rpcPort,_ = strconv.Atoi(rpcPortString)
		}
	}
	return &common.ServiceInstance{
		Host:      agentSvc.Address,
		Port:      agentSvc.Port,
		Weight:    agentSvc.Weights.Passing,
		GrpcPort:  rpcPort,
	}
}