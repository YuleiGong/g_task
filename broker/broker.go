package broker

import "github.com/YuleiGong/g_task/message"

type Broker interface {
	Clone() Broker
	Activate() error
	Push(taskID string, msg *message.Message) (err error)
	Pop() (taskID string, msg *message.Message, err error)
	Del(taskID string) (err error)
	Set(taskID string, msg *message.Message) (err error)
}
