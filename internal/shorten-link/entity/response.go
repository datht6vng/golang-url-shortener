package entity

import json "github.com/bytedance/sonic"

type Response struct {
	Code    int         `json:"code" xml:"code" form:"code"`
	Message string      `json:"message" xml:"message" form:"message"`
	Data    interface{} `json:"data" xml:"data" form:"data"`
}

func (c Response) MarshalBinary() (data []byte, err error) {
	bytes, err := json.Marshal(c)
	return bytes, err
}
