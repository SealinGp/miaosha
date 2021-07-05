package config

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/spf13/viper"
	conf "miaosha/pkg/config"
	"os"
)

const (
	kConfigType = "CONFIG_TYPE"
)

var Logger log.Logger

func init()  {
	Logger = log.NewLogfmtLogger(os.Stderr)
	Logger = log.With(Logger,"ts",log.DefaultTimestampUTC)
	Logger = log.With(Logger,"caller",log.DefaultCaller)
	viper.AutomaticEnv()
	initDefault()

	if err := conf.LoadRemoteConfig(); err != nil {
		Logger.Log("fail to load remote config ",err)
	}
	if err := conf.SubParse("auth",&AuthPermitConfig);err != nil {
		Logger.Log("fail to parse config ",err)
	}
	if err := conf.SubParse("hystrix",&HystrixConfig);err != nil {
		Logger.Log("fail to parse config ",err)
	}

}
func initDefault()  {
	viper.SetDefault(kConfigType,"yaml")
}

func GetZipkinUrl() string {
	return fmt.Sprintf("http://%s:%s%s",conf.TraceConfig.Host,conf.TraceConfig.Port,conf.TraceConfig.Url)
}