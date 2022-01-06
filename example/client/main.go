package main

import (
	"fmt"
	"time"

	"github.com/YuleiGong/g_task"
	"github.com/YuleiGong/g_task/backend"
	"github.com/YuleiGong/g_task/broker"
	"github.com/YuleiGong/g_task/client"
	"github.com/YuleiGong/g_task/message"
)

var (
	url      = "127.0.0.1:6379"
	db       = 1
	poolSize = 50
	password = ""
)

var cli *client.Client

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
	}

	if cli, err = g_task.Client(opts...); err != nil {
		fmt.Printf("%v", err)
		return
	}

	c1 := Cfg()
	//c2 := CfgWithTimeout()
	//c3 := CfgWithRetryNum()

	cfgs := []*client.SendConf{c1, c1, c1, c1}
	for i := 0; i < 100000; i++ {
		cfgs = append(cfgs, c1)
	}

	sig := []message.Signature{
		{Type: message.Int64, Val: 1},
		{Type: message.Int64, Val: 2},
	}
	var task []string
	for _, c := range cfgs {
		var taskID string
		if taskID, err = cli.Send(c, sig...); err != nil {
			fmt.Printf("%s\n", err.Error())
			return
		}
		task = append(task, taskID)
	}
	//TaskStatus(task)
	TaskResult(task)

}

func TaskStatus(task []string) {
	for _, t := range task {
		code, status := cli.Status(t)
		fmt.Printf("code %d status %s \n", code, status)
	}
}

func TaskResult(task []string) {
	var err error
	for _, t := range task {
		for !cli.IsFinish(t) {
			time.Sleep(1 * time.Second)
			continue
		}
		var res *message.MessageResult
		if res, err = cli.GetTaskResult(t); err != nil {
			fmt.Printf("err %v \n", err)
		}
		fmt.Printf("%+v \n", res)
	}
}
