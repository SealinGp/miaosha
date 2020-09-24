package config

import (
	"errors"
	"fmt"
	kitLog "github.com/go-kit/kit/log"
	zk "github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"github.com/spf13/viper"
	"miaosha/pkg/bootstrap"
	"miaosha/pkg/discover"
	"net/http"
	"os"
)

const (
	kConfigType = "CONFIG_TYPE"
)
var ZipkinTracer *zk.Tracer
var Logger kitLog.Logger

func init()  {
	Logger = kitLog.NewLogfmtLogger(os.Stderr)
	Logger = kitLog.With(Logger,"ts",kitLog.DefaultTimestampUTC)
	Logger = kitLog.With(Logger,"caller",kitLog.DefaultCaller)
	viper.AutomaticEnv()
	initDefault()

	Logger.Log("remote","trace")
	if err := LoadRemoteConfig();err != nil {
		Logger.Log("Fail to load remote config:",err.Error())
	}
	if err := SubParse("trace",&TraceConfig);err != nil {
		Logger.Log("Fail to parse trace:",err.Error())
	}
	zipkinUrl := fmt.Sprintf("http://%s:%s%s",TraceConfig.Host,TraceConfig.Port,TraceConfig.Url)
	Logger.Log("zipkin url",zipkinUrl)
	initTracer(zipkinUrl)
}

func initDefault()  {
	viper.SetDefault(kConfigType,"yaml")
}

func initTracer(zipkinURL string)  {
	var (
		err error
		userNoopTracer = zipkinURL == ""
		reporter = zipkinhttp.NewReporter(zipkinURL)
	)
	zEP, _ := zk.NewEndpoint(bootstrap.DiscoverConfig.ServiceName,bootstrap.HttpConfig.Port)
	ZipkinTracer, err = zk.NewTracer(
		reporter,
		zk.WithLocalEndpoint(zEP),
		zk.WithNoopTracer(userNoopTracer),
	)
	if err != nil {
		_ = Logger.Log("err",err.Error())
		os.Exit(1)
	}
	if !userNoopTracer {
		_ = Logger.Log("tracer","Zipkin","type","Native","URL",zipkinURL)
	}
}

//从对应服务的配置中心读取相应的配置
//label       : git分支 master|...
//profile     : 配置文件版本 dev|pro|test...
//application : 应用名称,服务名称
//configType  : 配置文件的类型(后缀名)
func LoadRemoteConfig() error {
	configServerInstance,err := discover.DiscoveryService(bootstrap.ConfigServerConfig.Id)
	if err != nil {
		return err
	}
	configServer := fmt.Sprintf("http://%s:%d",configServerInstance.Host,configServerInstance.Port)

	// http://配置中心地址/label/application-profile.configType
	// http://localhost:8888/master/miaosha-dev.yml
	confAddr     := fmt.Sprintf(
		"%v/%v/%v-%v.%v",
		configServer,bootstrap.ConfigServerConfig.Label,
		bootstrap.DiscoverConfig.ServiceName,bootstrap.ConfigServerConfig.Profile,
		viper.Get(kConfigType),
	)
	resp,err := http.Get(confAddr)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	viper.SetConfigType(viper.GetString(kConfigType))
	if err = viper.ReadConfig(resp.Body); err != nil {
		return err
	}
	Logger.Log("Load config from:",confAddr)
	return nil
}
func SubParse(key string,value interface{}) error {
	Logger.Log("配置文件前缀为:%s",key)
	sub := viper.Sub(key)
	if sub == nil {
		return errors.New("parse sub error")
	}
	sub.AutomaticEnv()
	sub.SetEnvPrefix(key)
	return sub.Unmarshal(value)
}