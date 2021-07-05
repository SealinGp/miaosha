package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"miaosha/pkg/discover"
	"miaosha/sk-core/service/srv_redis"
	"miaosha/sk-core/setup"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	//remote + local 配置初始化, zk初始化
	if err := setup.InitZk(); err != nil {
		log.Printf("[E] init zk failed. err:%s", err)
		return
	}
	if err := setup.InitRedis(); err != nil {
		log.Printf("[E] init redis failed. err:%s", err)
		return
	}

	log.Printf("[I] start run service...")
	ctx, cancel := context.WithCancel(context.Background())
	srv_redis.RunProcess(ctx)

	errCh := make(chan error)
	regErr := make(chan error)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				regErr <- errors.New(fmt.Sprintf("[E] recover from register. msg:%v", err))
			}
		}()
		discover.Register()
	}()

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		errCh <- fmt.Errorf("%s", <-sigCh)
	}()

	var err error
	select {
	case err = <-errCh:
		cancel()
		discover.Deregister()
	case err = <-regErr:
		cancel()
	}
	fmt.Printf("[E] err :%s", err)
}
