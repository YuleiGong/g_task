package backend

import "github.com/YuleiGong/g_task/message"

type Backend interface {
	SetResult(taskID string, msg *message.MessageResult) error
	GetResult(taskID string) (msg *message.MessageResult, err error)
	Activate() error
	Clone() Backend
}
