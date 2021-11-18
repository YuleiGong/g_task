package message

const (
	SUCCESS int64 = 0  //成功
	FAILURE int64 = -1 //失败
	PENDING int64 = 1  //排队中/挂起
	RETRY   int64 = 2
	STARTED int64 = 3 //任务开始执行
)

var FinishMap = map[int64]bool{
	SUCCESS: true,
	FAILURE: true,
}

var StatusMap = map[int64]string{
	SUCCESS: "SUCCESS",
	FAILURE: "FAILURE",
	PENDING: "PENDING",
	RETRY:   "RETRY",
	STARTED: "STARTED",
}
