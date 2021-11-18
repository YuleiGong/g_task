package client

import "time"

//send conf
type SendConf struct {
	funcName string
	timeout  time.Duration
	retryNum int64
}

func NewSendConf(funcName string) *SendConf {
	return &SendConf{funcName: funcName}
}

func (s *SendConf) SetTimeout(timeout time.Duration) {
	s.timeout = timeout
}

//retryNum == 0 代表不重试
func (s *SendConf) SetRetryNum(num int64) {
	s.retryNum = num
}
