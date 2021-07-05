package bootstrap


var (
	HttpConfig HttpConf
	DiscoverConfig DiscoverConf
	ConfigServerConfig ConfigServerConf
	RpcConfig RpcConf
	ZookeeperConfig ZookeeperConf
)

//Http配置
type HttpConf struct {
	Host string
	Port string
}

// RPC配置
type RpcConf struct {
	Port string
}

//服务注册与发现配置
type DiscoverConf struct {
	Host string
	Port string
	ServiceName string //应用注册的svc name
	Weight int
	InstanceId string  //svc id
}

//配置中心
type ConfigServerConf struct {
	Id string       //服务注册发现中心的配置中心注册的serviceName
	Profile string  //
	Label string    //git 分支
}
type ZookeeperConf struct {
	Hosts []string
	SecProductKey string
}