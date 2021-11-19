package message

const (
	Bool = iota
	Int
	Int8
	Int16
	Int32
	Int64
	Uint
	Uint8
	Uint16
	Uint32
	Uint64
	Float32
	Float64
	String
)

const (
	SUCCESS int64 = 0  //成功
	FAILURE int64 = -1 //失败
	PENDING int64 = 1  //排队中/挂起
	RETRY   int64 = 2  //重试
	STARTED int64 = 3  //任务开始执行
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
