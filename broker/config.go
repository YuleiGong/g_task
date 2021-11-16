package broker

import "time"

const (
	queueName  = "broker:matrix"
	popTimeout = 2 * time.Second
	expireTime = 2 * time.Hour
)
