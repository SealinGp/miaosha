package bootstrap

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"log"
)

func init() {
	viper.AutomaticEnv()
	initBootstrapConfig()

	log.Println("bootstrap_config.go (local bootstrap.yaml):  http,rpc,discover,config ")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("err:", err.Error())
	}
	if err := SubParse("http", &HttpConfig); err != nil {
		log.Fatal("Fail to parse Http config", err.Error())
	}
	if err := SubParse("rpc", &RpcConfig); err != nil {
		log.Fatal("Fail to parse rpc server", err.Error())
	}
	if err := SubParse("discover", &DiscoverConfig); err != nil {
		log.Fatal("Fail to parse discover config", err.Error())
	}
	if err := SubParse("config", &ConfigServerConfig); err != nil {
		log.Fatal("Fail to parse config server", err.Error())
	}
	if err := SubParse("zookeeper", &ZookeeperConfig); err != nil {
		log.Fatal("Fail to parse zookeeper server:", err.Error())
	}
	log.Printf("[I] bootstrap init")
}

func initBootstrapConfig() {
	viper.SetConfigName("bootstrap")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./")
	viper.AddConfigPath("$GOPATH/src/")
}

func SubParse(key string, value interface{}) error {
	log.Printf("配置文件前缀为:%s", key)
	sub := viper.Sub(key)
	if sub == nil {
		return errors.New("parse sub error")
	}
	sub.AutomaticEnv()
	sub.SetEnvPrefix(key)
	return sub.Unmarshal(value)
}
