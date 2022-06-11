package entity

type Response struct {
	Code    int         `json:"code" xml:"code" form:"code"`
	Message string      `json:"message" xml:"message" form:"message"`
	Data    interface{} `json:"data" xml:"data" form:"data"`
}
