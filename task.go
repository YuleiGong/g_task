package g_task

import (
	"github.com/YuleiGong/g_task/client"
	"github.com/YuleiGong/g_task/server"
)

func Server(opts ...server.WorkerOpt) (svr *server.Server) {
	svr = server.NewServer(opts...)

	return svr
}

func Client(opts ...client.ClientOpt) (cli *client.Client, err error) {
	if cli, err = client.GetClient(opts...); err != nil {
		return
	}

	return cli, err
}
