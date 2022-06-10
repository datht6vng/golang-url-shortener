package _const

type UrlExpiredError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (this UrlExpiredError) Error() string {
	return this.Message
}
