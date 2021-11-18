package client

import (
	"errors"

	"github.com/YuleiGong/g_task/backend"
	"github.com/YuleiGong/g_task/broker"
	"github.com/YuleiGong/g_task/log"
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
	if cli.backend == nil {
		err = errors.New("please run server or set backend and broker")
	}

	return cli, err
}

func (c *Client) Send(sendConf *SendConf, args ...interface{}) (taskID string, err error) {
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

func (c *Client) Status(taskID string) (code int64, status string) {
	var msg *message.Message
	var err error
	if msg, err = c.broker.Get(taskID); err != nil {
		if !errors.Is(err, broker.ErrBrokerNil) {
			log.Error("%v", err)
		}
		return
	}
	return msg.Status, message.StatusMap[msg.Status]

}

func (c *Client) IsFinish(taskID string) bool {
	var msg *message.Message
	var err error
	if msg, err = c.broker.Get(taskID); err != nil {
		if !errors.Is(err, broker.ErrBrokerNil) {
			log.Error("%v", err)
		}
		return false
	}

	return message.FinishMap[msg.Status]
}

func (c *Client) GetTaskResult(taskID string) (result *message.MessageResult, err error) {
	return c.backend.GetResult(taskID)
}
