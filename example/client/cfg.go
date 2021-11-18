package main

import (
	"time"

	"github.com/YuleiGong/g_task/client"
)

func cfg() *client.sendConf {
	sendConf := client.NewSendConf("add")

	return sendConf
}

//超时
func cfgWithTimeout() *client.sendConf {
	sendConf := client.NewSendConf("add")
	sendConf.SetTimeout(2 * time.Second)

	return sendConf
}

//超时重试
func cfgWithRetryNum() *client.sendConf {
	sendConf := client.NewSendConf("add")
	sendConf.SetRetryNum(2)

	return sendConf
}
