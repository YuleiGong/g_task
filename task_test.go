//1: go test task_test.go task.go --run="TestServerRun" -v
//2: go test task_test.go task.go --run="TestClientRun" -v
package task

import (
	"runtime/debug"
	"testing"
	"time"

	"github.com/YuleiGong/g_task/backend"
	"github.com/YuleiGong/g_task/broker"
	"github.com/YuleiGong/g_task/client"
	"github.com/YuleiGong/g_task/server"
)

var (
	url      = "127.0.0.1:6379"
	db       = 1
	poolSize = 50
	password = ""
)

func TestServerRun(t *testing.T) {
	var err error
	defer func() {
		if err := recover(); err != nil {
			t.Logf("%s", string(debug.Stack()))
			t.Logf("%v", err)
		}
	}()

	brokerCfg := broker.NewRedisConf(url, password, db)
	brokerCfg.SetPoolSize(poolSize)
	brokerCfg.SetExpireTime(1 * time.Hour)

	backendCfg := backend.NewRedisConf(url, password, db)
	backendCfg.SetPoolSize(poolSize)
	backendCfg.SetExpireTime(1 * time.Hour)

	opts := []server.WorkerOpt{
		server.WithBroker(broker.NewRedis(brokerCfg)),
		server.WithBackend(backend.NewRedis(backendCfg)),
	}

	svr := Server(opts...)

	svr.Reg("add", add)
	if err = svr.Run(10); err != nil {
		t.Logf("%v", err)
	}
	defer svr.ShutDown()

}

func add(a, b int) (int, error) {

	return a + b, nil
}

func TestClientRun(t *testing.T) {
	var err error
	var opts []client.ClientOpt

	brokerCfg := broker.NewRedisConf(url, password, db)
	brokerCfg.SetPoolSize(poolSize)
	brokerCfg.SetExpireTime(1 * time.Hour)

	backendCfg := backend.NewRedisConf(url, password, db)
	backendCfg.SetPoolSize(poolSize)
	backendCfg.SetExpireTime(1 * time.Hour)

	opts = append(opts, client.WithBroker(broker.NewRedis(brokerCfg)))
	opts = append(opts, client.WithBackend(backend.NewRedis(backendCfg)))

	cli, err := Client(opts...)
	if err != nil {
		t.Fatal(err)
	}

	sendConf := client.NewSendConf("add")
	for i := 0; i < 1000; i++ {
		var taskID string
		if taskID, err = cli.Send(sendConf, 1, 2); err != nil {
			t.Fatal(err)
		}
		t.Log(taskID)
	}
}
