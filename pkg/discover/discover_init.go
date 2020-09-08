package discover

import (
	"errors"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"log"
	"miaosha/pkg/bootstrap"
	"miaosha/pkg/common"
	"miaosha/pkg/loadbalance"
	"net/http"
	"os"
	"strconv"
)

var (
	ConsulService DiscoveryClient
	LoadBalance loadbalance.LoadBalance
	Logger *log.Logger
	NoInstanceExistedErr = errors.New("no available client")
)

func init()  {
	port,_       := strconv.Atoi(bootstrap.DiscoverConfig.Port)
	ConsulService = NewDiscoveryClient(bootstrap.DiscoverConfig.Host,port)
	LoadBalance   = new(loadbalance.RandomLoadBalance)
	Logger        = log.New(os.Stderr,"",log.LstdFlags)
}

func CheckHealth(writer http.ResponseWriter, reader *http.Request)  {
	Logger.Println("Health check!")
	_, err := fmt.Fprintln(writer,"server is ok")
	if err != nil {
		Logger.Println(err)
	}
}

func DiscoveryService(serviceName string) (*common.ServiceInstance,error) {
	instances := ConsulService.DiscoverServices(serviceName,Logger)
	if len(instances) < 1 {
		Logger.Printf("no avaliable client for %s.",serviceName)
		return nil,NoInstanceExistedErr
	}
	return LoadBalance.SelectService(instances)
}

func Register()  {
	if ConsulService == nil {
		panic(0)
	}
	instanceId := bootstrap.DiscoverConfig.InstanceId
	if instanceId == "" {
		instanceId = bootstrap.DiscoverConfig.ServiceName + "-" + uuid.NewV4().String()
	}

	if !ConsulService.Register(
		instanceId, bootstrap.HttpConfig.Host, "/health",
		bootstrap.HttpConfig.Port,bootstrap.DiscoverConfig.ServiceName,
		bootstrap.DiscoverConfig.Weight, map[string]string{"rpcPort":bootstrap.RpcConfig.Port},
		nil,Logger,
	) {
		Logger.Printf("register service %s failed.",bootstrap.DiscoverConfig.ServiceName)
		panic(0)
	}
	Logger.Printf(bootstrap.DiscoverConfig.ServiceName+"-service for service %s success.",bootstrap.DiscoverConfig.ServiceName )
}

func Deregister()  {
	if ConsulService == nil {
		panic(0)
	}
	instanceId := bootstrap.DiscoverConfig.InstanceId
	if instanceId == "" {
		instanceId = bootstrap.DiscoverConfig.ServiceName + "-" + uuid.NewV4().String()
	}
	if !ConsulService.DeRegister(instanceId,Logger) {
		Logger.Printf("deregister for service %s failed",bootstrap.DiscoverConfig.ServiceName)
		panic(0)
	}
}