package config

import (
	"github.com/go-redis/redis"
	"github.com/samuel/go-zookeeper/zk"
	"miaosha/sk-core/service/srv_limit"
	"sync"
)

var (
	Redis       RedisConf
	SecKill     SecKillConf
	MysqlConfig MysqlConf
	TraceConfig TraceConf
	Zk          ZookeeperConf
)

type ZookeeperConf struct {
	ZkConn *zk.Conn
	SecProductKey string //商品键
}


type TraceConf struct {
	Host string
	Port string
	Url  string
}
type MysqlConf struct {
	Host string
	Port string
	User string
	Pwd  string
	Db   string
}

type RedisConf struct {
	RedisConn *redis.Client
	Proxy2layerQueueName string //队列名称
	Layer2proxyQueueName string //队列名称

	IdBlackListHash      string //用户黑名单hash表
	IpBlackListHash      string //ip黑名单hash表

	IdBlackListQueue     string //用户黑名单队列
	IpBlackListQueue     string //ip黑名单队列
	Host                 string
	Password             string
	Db                   int
}

type SecKillConf struct {
	RedisConf *RedisConf
	CookieSecretKey string
	ReferWhiteList []string //白名单

	AccessLimitConf AccessLimitConf

	RWBlackLock sync.RWMutex //读共享,写独占锁
	WriteProxy2LayerGoroutineNum int
	ReadProxy2LayerGoroutineNum int

	IPBlackMap map[string]bool
	IDBlackMap map[int]bool

	SecProductInfoMap map[int]*SecProductInfoConf

	AppWriteToHandleGoroutineNum int
	AppReadToHandleGoroutineNum int

	CoreReadRedisGoroutineNum  int
	CoreWriteRedisGoroutineNum int
	CoreHandleGoroutineNum     int

	AppWaitResultTimeout int
	CoreWaitResultTimeout int
	MaxRequestWaitTimeout int
	SendToWriteChanTimeout int
	SendToHandleChanTimeout int
	TokenPasswd string
}

type SecProductInfoConf struct {
	ProductId int `json:"product_id"` //商品ID
	StartTime int64 `json:"start_time"` //开始时间
	EndTime   int64 `json:"end_time"` //结束时间
	Status    int `json:"status"`//状态
	Total     int `json:"total"`//商品总数量
	Left      int `json:"left"`//商品剩余数量
	OnePersonBuyLimit int `json:""`//单个用户限制购买数量
	BuyRate   float64 `json:"buy_rate"`//购买频率限制
	SoldMaxLimit int `json:"sold_max_limit"`
	//todo error
	SecLimit  *srv_limit.SecLimit `json:"sec_limit"`//限速控制
}

//访问限制
type AccessLimitConf struct {
	IpSecAccessLimit   int //IP每秒访问限制
	UserSecAccessLimit int //用户每秒访问限制
	IPMinAccessLimit   int //IP每分钟访问限制
	UserMinAccessLimit int //用户每分钟访问限制
}