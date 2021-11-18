package main

import (
	"time"

	"github.com/YuleiGong/g_task/client"
)

func Cfg() *client.SendConf {
	sendConf := client.NewSendConf("add")

	return sendConf
}

//超时
func CfgWithTimeout() *client.SendConf {
	sendConf := client.NewSendConf("add")
	sendConf.SetTimeout(2 * time.Second)

	return sendConf
}

//超时重试
func CfgWithRetryNum() *client.SendConf {
	sendConf := client.NewSendConf("add")
	sendConf.SetTimeout(2 * time.Second)
	sendConf.SetRetryNum(2)

	return sendConf
}
