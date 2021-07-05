package srv_redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	conf "miaosha/pkg/config"
	"miaosha/sk-core/config"
	"time"
)

func RunProcess(ctx context.Context) {
	for i := 0; i < conf.SecKill.CoreReadRedisGoroutineNum; i++ {
		config.SecLayerCtx.WaitGroup.Add(1)
		go HandleReader(ctx)
	}

	for i := 0; i < conf.SecKill.CoreWriteRedisGoroutineNum; i++ {
		config.SecLayerCtx.WaitGroup.Add(1)
		go HandleWrite(ctx)
	}

	for i := 0; i < conf.SecKill.CoreHandleGoroutineNum; i++ {
		config.SecLayerCtx.WaitGroup.Add(1)
		go HandleUser(ctx)
	}

	log.Printf("all process goroutine started")
	config.SecLayerCtx.WaitGroup.Wait()
	log.Printf("wait all goroutine exited")
	return
}

func HandleReader(pctx context.Context) {
	for {
		conn := conf.Redis.RedisConn
		for {
			select {
			case <-pctx.Done():
				return
			default:
			}

			//从Redis队列中取出数据
			data, err := conn.BRPop(time.Second, conf.Redis.Proxy2layerQueueName).Result()
			if err != nil {
				continue
			}
			log.Printf("brpop from proxy to layer queue, data : %s\n", data)

			//转换数据结构
			var req config.SecRequest
			err = json.Unmarshal([]byte(data[1]), &req)
			if err != nil {
				log.Printf("unmarshal to secrequest failed, err : %v", err)
				continue
			}

			//判断是否超时
			nowTime := time.Now().Unix()
			fmt.Println(nowTime, " ", req.SecTime, " ", 100)
			if nowTime-req.SecTime >= int64(conf.SecKill.MaxRequestWaitTimeout) {
				log.Printf("req[%v] is expire", req)
				continue
			}

			//设置超时时间
			ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(conf.SecKill.CoreWaitResultTimeout))
			select {
			case config.SecLayerCtx.Read2HandleChan <- &req:
				cancel()
			case <-ctx.Done():
				cancel()
				log.Printf("send to handle chan timeout, req : %v", req)
				break
			}
		}
	}
}

func HandleWrite(ctx context.Context) {
	log.Println("handle write running")

	for res := range config.SecLayerCtx.Handle2WriteChan {
		select {
		case <-ctx.Done():
			return
		default:
		}

		fmt.Println("===", res)
		err := sendToRedis(res)
		if err != nil {
			log.Printf("send to redis, err : %v, res : %v", err, res)
			continue
		}
	}
}

func sendToRedis(res *config.SecResult) (err error) {
	data, err := json.Marshal(res)
	if err != nil {
		log.Printf("marshal failed, err : %v", err)
		return
	}

	fmt.Printf("推入队列前~~ %v", conf.Redis.Layer2proxyQueueName)
	conn := conf.Redis.RedisConn
	err = conn.LPush(conf.Redis.Layer2proxyQueueName, string(data)).Err()
	fmt.Println("推入队列后~~")
	if err != nil {
		log.Printf("rpush layer to proxy redis queue failed, err : %v", err)
		return
	}
	log.Printf("lpush layer to proxy success. data[%v]", string(data))

	return
}