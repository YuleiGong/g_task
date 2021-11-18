package message

import "encoding/json"

type MessageResult struct {
	TaskID  string   `json:"task_id"`
	Val     []string `json:"val"`
	ErrMsg  string   `json:"err_msg"`
	ErrCode int64    `json:"err_code"`
}

var DefaultMessageResult = &MessageResult{}

func (m *MessageResult) Serializa() (s string, err error) {
	var b []byte
	if b, err = json.Marshal(m); err != nil {
		return
	}
	return string(b), err
}

func (m *MessageResult) Deserialize(bs []byte) (*MessageResult, error) {
	var (
		err error
		res MessageResult
	)

	if err = json.Unmarshal(bs, &res); err != nil {
		return nil, err
	}

	return &res, err

}
