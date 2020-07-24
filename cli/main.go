package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"log"
)

type ZookeeperConf struct {
	Hosts []string
	SecProductKey string
}

func main() {
	var ZookeeperConfig ZookeeperConf
	viper.AutomaticEnv()
	viper.SetConfigName("bootstrap")
	viper.AddConfigPath("./cli")     //运行时的当前路径的相对路径
	viper.SetConfigType("yaml")

	fmt.Println(viper.ReadInConfig())
	if err := SubParse("zookeeper",&ZookeeperConfig);err != nil {
		log.Fatal("Fail to parse zookeeper server",err.Error())
	}
	b,_ := json.Marshal(ZookeeperConfig)
	fmt.Println(string(b))
}

func SubParse(key string,value interface{}) error {
	log.Printf("配置文件前缀为:%s",key)
	sub := viper.Sub(key)
	if sub == nil {
		return errors.New("?? error")
	}
	sub.AutomaticEnv()
	sub.SetEnvPrefix(key)
	return sub.Unmarshal(value)
}

