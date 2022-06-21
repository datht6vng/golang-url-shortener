package entity

import (
	json "github.com/bytedance/sonic"
)

type URLData struct {
	URL      string `json:"url" xml:"url" form:"url"`
	ClientID string `json:"client_id" xml:"client_id" form:"client_id"`
}

func (u URLData) MarshalBinary() (data []byte, err error) {
	bytes, err := json.Marshal(u)
	return bytes, err
}

type GenerateCounterKey struct {
	ClientID   string `json:"client_id" xml:"client_id" form:"client_id"`
	CreateDate string `json:"create_date" xml:"create_date" form:"create_date"`
}

func (g GenerateCounterKey) MarshalBinary() (data []byte, err error) {
	bytes, err := json.Marshal(g)
	return bytes, err
}
