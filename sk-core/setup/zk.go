package setup

import (
	"encoding/json"
	"github.com/samuel/go-zookeeper/zk"
	"log"
	"miaosha/pkg/bootstrap"
	remoteCfg "miaosha/pkg/config"
	"time"
)

type EtcdCfg struct {
	Host              string
	EtcdSecProductKey string
}

func InitZk() error {
	etcdcfg := &EtcdCfg{}
	err := bootstrap.SubParse("etcd", etcdcfg)
	if err != nil {
		return err
	}
	remoteCfg.Etcd.EtcdSecProductKey = etcdcfg.EtcdSecProductKey
	remoteCfg.Etcd.Host = etcdcfg.Host

	conn, _, err := zk.Connect(bootstrap.ZookeeperConfig.Hosts, time.Second*5, zk.WithEventCallback(waitSecProductEvent))
	if err != nil {
		return err
	}

	remoteCfg.Zk.ZkConn = conn
	remoteCfg.Zk.SecProductKey = "/product"
	go loadSecConf(conn)
	return nil
}

func loadSecConf(conn *zk.Conn) {
	log.Printf("[E] connect zk success. etcdKey:%s", remoteCfg.Etcd.EtcdSecProductKey)
	v, _, err := conn.Get(remoteCfg.Zk.SecProductKey)
	if err != nil {
		log.Printf("[E] zk get secProductKey failed. err:%s", err)
		return
	}
	var secProductInfo []*remoteCfg.SecProductInfoConf
	err = json.Unmarshal(v, &secProductInfo)
	if err != nil {
		log.Printf("[E] zk parse secProductKey failed. err:%s", err)
		return
	}

	updateSecProductInfo(secProductInfo)
}

func waitSecProductEvent(event zk.Event) {
	if event.Path == remoteCfg.Zk.SecProductKey {

	}
}

func updateSecProductInfo(secProductInfo []*remoteCfg.SecProductInfoConf) {
	tmp := make(map[int]*remoteCfg.SecProductInfoConf, 1024)
	for _, v := range secProductInfo {
		tmp[v.ProductId] = v
	}
	remoteCfg.SecKill.RWBlackLock.Lock()
	remoteCfg.SecKill.SecProductInfoMap = tmp
	remoteCfg.SecKill.RWBlackLock.Unlock()
}
