package queue

import "encoding/json"

type MessageDto struct {
	ID       int    `json:"id"`
	Filename string `json:"filename"`
	Path     string `json:"path"`
}

func (msg *MessageDto) Marshal() ([]byte, error) {
	return json.Marshal(msg)
}

func (msg *MessageDto) Unmarshal(data []byte) error {
	return json.Unmarshal(data, msg)
}
