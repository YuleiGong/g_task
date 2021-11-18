package main

import (
	"fmt"
	"time"

	"github.com/YuleiGong/g_task"
	"github.com/YuleiGong/g_task/backend"
	"github.com/YuleiGong/g_task/broker"
	"github.com/YuleiGong/g_task/server"
)

var (
	url      = "127.0.0.1:6379"
	db       = 1
	poolSize = 50
	password = ""
)

func add(a, b int) (int, error) {

	return a + b, nil
}
func main() {
	//broker
	brokerCfg := broker.NewRedisConf(url, password, db)
	brokerCfg.SetPoolSize(poolSize)
	brokerCfg.SetExpireTime(1 * time.Hour)

	//backend
	backendCfg := backend.NewRedisConf(url, password, db)
	backendCfg.SetPoolSize(poolSize)
	backendCfg.SetExpireTime(1 * time.Hour)

	opts := []server.WorkerOpt{
		server.WithBroker(broker.NewRedis(brokerCfg)),
		server.WithBackend(backend.NewRedis(backendCfg)),
	}
	svr := g_task.Server(opts...)
	//函数注册
	svr.Reg("add", add)
	if err := svr.Run(10); err != nil {
		fmt.Printf("%v", err)
	}
	//defer svr.ShutDown()

}
