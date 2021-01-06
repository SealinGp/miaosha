package setup

import (
	"fmt"
	"miaosha/pkg/discover"
	"miaosha/sk-core/service/srv_redis"
	"os"
	"os/signal"
	"syscall"
)

func RunService() {
	srv_redis.RunProcess()
	errCh := make(chan error)
	go func() {
		discover.Register()
	}()
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		errCh <- fmt.Errorf("%s", <-sigCh)
	}()

	err := <-errCh
	discover.Deregister()
	fmt.Printf("[E] err :%s", err)
}
