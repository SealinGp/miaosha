package setup

import (
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"miaosha/pkg/bootstrap"
	conf "miaosha/pkg/config"
	"os"
	"time"
)

//初始化zookeeper
func InitZk() {
	conn, _, err := zk.Connect(bootstrap.ZookeeperConfig.Hosts, time.Second*5)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	conf.Zk.ZkConn = conn
	conf.Zk.SecProductKey = bootstrap.ZookeeperConfig.SecProductKey
}
