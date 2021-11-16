package message

import "encoding/json"

type MessageResult struct {
	TaskID  string   `json:"task_id"`
	Val     []string `json:"val"`
	ErrMsg  string   `json:"err_msg"`
	ErrCode int64    `json:"err_code"`
}

func (m *MessageResult) Serializa() (s string, err error) {
	var b []byte
	if b, err = json.Marshal(m); err != nil {
		return
	}
	return string(b), err
}
