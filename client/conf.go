package client

import "time"

//send conf
type sendConf struct {
	funcName string
	timeout  time.Duration
	retryNum int64
}

func NewSendConf(funcName string) *sendConf {
	return &sendConf{funcName: funcName}
}

func (s *sendConf) SetTimeout(timeout time.Duration) {
	s.timeout = timeout
}

//retryNum == 0 代表不重试
func (s *sendConf) SetRetryNum(num int64) {
	s.retryNum = num
}
