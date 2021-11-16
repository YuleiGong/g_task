package client

import (
	"github.com/YuleiGong/g_task/backend"
	"github.com/YuleiGong/g_task/broker"
	"github.com/YuleiGong/g_task/message"
	"github.com/YuleiGong/g_task/server"
)

type Client struct {
	broker  broker.Broker
	backend backend.Backend
}

var cli *Client

type ClientOpt func(*Client) (err error)

func WithBroker(broker broker.Broker) ClientOpt {
	return func(c *Client) error {
		c.broker = broker
		return c.broker.Activate()
	}
}

func WithBackend(backend backend.Backend) ClientOpt {
	return func(c *Client) error {
		c.backend = backend
		return c.backend.Activate()
	}
}

func GetClient(opts ...ClientOpt) (*Client, error) {
	var err error
	if cli != nil {
		return cli, err
	}
	cli = &Client{}

	for _, opt := range opts {
		if err = opt(cli); err != nil {
			return cli, err
		}
	}

	if server.GetServer() != nil {
		cli.backend = server.GetServer().CloneBackend()
		cli.broker = server.GetServer().CloneBroker()
		if err = cli.backend.Activate(); err != nil {
			return cli, err
		}
		if err = cli.broker.Activate(); err != nil {
			return cli, err
		}
	}

	return cli, err
}

func (c *Client) Send(sendConf *sendConf, args ...interface{}) (taskID string, err error) {
	var m *message.Message
	if m, err = message.NewMessage(sendConf.funcName, args...); err != nil {
		return
	}
	m.SetTimeout(sendConf.timeout)
	m.SetRetryNum(sendConf.retryNum)

	if err = c.broker.Push(m.TaskID, m); err != nil {
		return
	}

	return m.TaskID, err
}
