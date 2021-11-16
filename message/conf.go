package message

const (
	SUCCESS int64 = 0  //成功
	FAILURE int64 = -1 //失败
	PENDING int64 = 1  //排队中/挂起
	RETRY   int64 = 2
	STARTED int64 = 3 //任务开始执行
)
