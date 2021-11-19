package message

import (
	"encoding/json"
	"time"

	"github.com/twinj/uuid"
)

type Message struct {
	TaskID      string        `json:"task_id"`
	FuncName    string        `json:"func_name"`
	Args        []Signature   `json:"args"`
	Status      int64         `json:"status"`
	RetryNum    int64         `json:"retry_num"`
	MaxRetryNum int64         `json:"max_retry_num"`
	Timeout     time.Duration `json:"timeout"`
}
type Signature struct {
	Type int         `json:"Type"`
	Val  interface{} `json:"val"`
}

var DefaultMessage = &Message{}

func NewMessage(funcName string, args ...Signature) (m *Message, err error) {
	m = &Message{}
	m.Args = args
	m.FuncName = funcName
	m.TaskID = uuid.NewV4().String()
	m.Status = PENDING

	return m, err
}

func (m *Message) SetTimeout(timeout time.Duration) {
	m.Timeout = timeout
}

func (m *Message) SetRetryNum(num int64) {
	m.MaxRetryNum = num
}

func (m *Message) Failure() {
	m.Status = FAILURE
}

func (m *Message) Success() {
	m.Status = SUCCESS
}

func (m *Message) Started() {
	m.Status = STARTED
}

func (m *Message) Retry() {
	m.Status = RETRY
}

func (m *Message) AddRetry() {
	m.RetryNum = m.RetryNum + 1
}

func (m *Message) IsRetryOpt() bool {
	return m.MaxRetryNum > 0
}

func (m *Message) IsRetry() bool {
	return m.RetryNum < m.MaxRetryNum && m.MaxRetryNum > 0
}

func (m *Message) IsTimeoutOpt() bool {
	return m.Timeout >= 1
}

func (m *Message) Serialize() (s string, err error) {
	var b []byte
	if b, err = json.Marshal(m); err != nil {
		return
	}
	return string(b), err
}

func (m *Message) Deserialize(bs []byte) (*Message, error) {
	var (
		err error
		msg Message
	)

	if err = json.Unmarshal(bs, &msg); err != nil {
		return nil, err
	}

	return &msg, err

}
