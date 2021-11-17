package main

import (
	"fmt"
	"time"

	"github.com/YuleiGong/g_task"
	"github.com/YuleiGong/g_task/backend"
	"github.com/YuleiGong/g_task/broker"
	"github.com/YuleiGong/g_task/client"
)

var (
	url      = "127.0.0.1:6379"
	db       = 1
	poolSize = 50
	password = ""
)

func main() {
	var err error
	//broker
	brokerCfg := broker.NewRedisConf(url, password, db)
	brokerCfg.SetPoolSize(poolSize)
	brokerCfg.SetExpireTime(1 * time.Hour)
	//backend
	backendCfg := backend.NewRedisConf(url, password, db)
	backendCfg.SetPoolSize(poolSize)
	backendCfg.SetExpireTime(1 * time.Hour)
	opts := []client.ClientOpt{
		client.WithBroker(broker.NewRedis(brokerCfg)),
		client.WithBackend(backend.NewRedis(backendCfg)),
	} //实际使用中，不需要初始化broker broker, client会自动复用server的配置
	var cli *client.Client
	if cli, err = g_task.Client(opts...); err != nil {
		fmt.Printf("%v", err)
		return
	}

	sendConf := client.NewSendConf("add")
	var taskID string
	if taskID, err = cli.Send(sendConf, 1, 2); err != nil {
		fmt.Printf("%s", err.Error())
		return
	}
	fmt.Printf("%s", taskID)

}
