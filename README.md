# g_task 异步任务队列
* broker :redis 
* backend :redis
* worker

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
* 初始化
* 配置
* 发送任务

## TimeoutTask

## RetryTask

## Broker

## Backend

