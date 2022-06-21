package entity

import json "github.com/bytedance/sonic"

type UpdateType string

const (
	UpdateLimit UpdateType = "update_limit"
)

type UpdateClientRequest struct {
	ClientId   string `json:"client_id"`
	UpdateType string `json:"update_type"`
	Limit      int64  `json:"limit"`
}

func (c UpdateClientRequest) MarshalBinary() (data []byte, err error) {
	bytes, err := json.Marshal(c)
	return bytes, err
}
