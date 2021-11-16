//1: go test task_test.go task.go --run="TestServerRun" -v
//2: go test task_test.go task.go --run="TestClientRun" -v
package task

import (
	"runtime/debug"
	"testing"

	"github.com/YuleiGong/g_task/backend"
	"github.com/YuleiGong/g_task/broker"
	"github.com/YuleiGong/g_task/client"
	"github.com/YuleiGong/g_task/server"
)

func TestServerRun(t *testing.T) {
	var err error
	defer func() {
		if err := recover(); err != nil {
			t.Logf("%s", string(debug.Stack()))
			t.Logf("%v", err)
		}
	}()

	var opts []server.WorkerOpt
	var (
		url      = "127.0.0.1:6379"
		db       = 1
		poolSize = 100
		password = ""
	)

	opts = append(opts, server.WithBroker(broker.NewRedis(url, password, db, poolSize)))
	opts = append(opts, server.WithBackend(backend.NewRedis(url, password, db, poolSize)))

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
	var (
		url      = "127.0.0.1:6379"
		db       = 1
		poolSize = 2
		password = ""
	)

	opts = append(opts, client.WithBroker(broker.NewRedis(url, password, db, poolSize)))
	opts = append(opts, client.WithBackend(backend.NewRedis(url, password, db, poolSize)))

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
