# g_task 异步任务队列

* [特点](#特点)
* [QuickStart](#QuickStart)
  * [server](#server)
  * [client](#client)
  * [example](#example)
* [Server](#Server)
* [Client](#Server)
* [TimeoutTask](#Timeout)
* [RetryTask](#RetryTask)
* [Broker](#Broker)
* [Backend](#Backend)
* [任务状态标识](#任务状态标识)
 

## 特点
* 简单，无侵入
* 支持任务超时设置
* 支持任务超时重试
* 方便扩展broker backend,目前支持redis


## QuickStart

### server

```go
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
		fmt.Println("%v", err)
	}
	defer svr.ShutDown()

}
```
### client

```go
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
```

### Example
[example](https://github.com/YuleiGong/g_task/tree/main/example) 目录下有更多的样例可供参考


## Server
* 初始化: 通过 ``` Server ``` 函数获取一个服务。

```go
opts := []server.WorkerOpt{
	server.WithBroker(broker.NewRedis(brokerCfg)),
	server.WithBackend(backend.NewRedis(backendCfg)),
}
svr := g_task.Server(opts...)

```

* 配置: 需要配置服务的 __broker__ 和 __backend__ 。 详细见 broker backend 章节

```go
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
```

* 任务注册: 一个任务，就是一个 __函数__ 。注册后，可以作为异步任务执行。注册函数至少要有一个error 返回值

```go
func add(a, b int) (int, error) {

	return a + b, nil
}
svr.Reg("add", add)

```

* 启动和退出

```go
if err := svr.Run(10); err != nil {
		fmt.Printf("%v", err)
	}
defer svr.ShutDown()
```


## Client
* 初始化: 通过``` Client ``` 获取一个客户端。

```go
var cli *client.Client
if cli, err = g_task.Client(opts...); err != nil {
	fmt.Printf("%v", err)
	return
}
```

* 配置: 需要配客户端的 __broker__ 和 __backend__ 。 broker和backend必须和server保持的一致。

```go
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
```
* 配置(实际使用): 实际使用中，不必初始化 __broker__ 和 __backend__ 。client会自动选择server的broker/backend配置。

```go
var cli *client.Client
if cli, err = g_task.Client(opts...); err != nil {
	fmt.Printf("%v", err)
	return
}
```

* 发送任务: 发送任务，需要初始化一个 __sendConf__ 配置。

```go
sendConf := client.NewSendConf("add")
if taskID, err = cli.Send(sendConf, 1, 2); err != nil {
    return
}
```

* 查看任务执行情况

```go
func TaskStatus(task []string) {
	for _, t := range task {
		code, status := cli.Status(t)
		fmt.Printf("code %d status %s \n", code, status)
	}
}
```

```go
func TaskResult(task []string) {
	var err error
	for _, t := range task {
		for !cli.IsFinish(t) {
			time.Sleep(1 * time.Second)
		}
		var res *message.MessageResult
		if res, err = cli.GetTaskResult(t); err != nil {
			fmt.Printf("err %v \n", err)
		}
		fmt.Printf("%+v \n", res)
	}
}
```



## TimeoutTask
* 支持为异步任务设置超时时间，在发送任务的时候，需要配置超时时间。默认情况下，无超时。
* 任务超时后，可以在 __backend__ 中查看任务状态。

```go
func cfgWithTimeout() *client.sendConf {
	sendConf := client.NewSendConf("add")
	sendConf.SetTimeout(2 * time.Second)

	return sendConf
}

```

## RetryTask
* 支持为超时事件，设置重试，并设置重试次数。__注意__: 同时设置了Timeout 和 Retrynum ，才能触发RetryTask。

```go
//超时重试
func cfgWithRetryNum() *client.sendConf {
	sendConf := client.NewSendConf("add")
	sendConf.SetTimeout(2 * time.Second)
	sendConf.SetRetryNum(2)

	return sendConf
}
```

## Broker
* 使用 Broker 与任务队列通信，目前支持的任务队列只有 __redis__

```go
//初始化
brokerCfg := broker.NewRedisConf(url, password, db)
brokerCfg.SetPoolSize(poolSize)
brokerCfg.SetExpireTime(1 * time.Hour)//默认2H

```

* 通过实现以下接口，可以自定义broker

```go
type Broker interface {
	Clone() Broker
	Activate() error
	Push(taskID string, msg *message.Message) (err error)
	Pop() (taskID string, msg *message.Message, err error)
	Del(taskID string) (err error)
	Set(taskID string, msg *message.Message) (err error)
	Get(taskID string) (msg *message.Message, err error)
}
```

## Backend
* 使用 Backend 存储任务结果，目前支持的Backend只有 __redis__

```go
//初始化
//backend
backendCfg := backend.NewRedisConf(url, password, db)
backendCfg.SetPoolSize(poolSize)
backendCfg.SetExpireTime(1 * time.Hour) //消息存放时间，默认2H
```


* 通过实现以下接口，可以自定义backend

```
type Backend interface {
	SetResult(taskID string, msg *message.MessageResult) error
	GetResult(taskID string) (string, error)
	Activate() error
	Clone() Backend
}
```

## 任务状态标识

  * __SUCCES__ int64 = 0  //成功
  * __FAILURE__ int64 = -1 //失败
  * __PENDING__ int64 = 1  //排队中/挂起
  * __RETRY__  int64 = 2  //重试
  * __STARTED__ int64 = 3  //任务开始执行

