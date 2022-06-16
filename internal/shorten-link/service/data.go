package service

type URLData struct {
	URL      string `json:"url" xml:"url" form:"url"`
	ClientID string `json:"client_id" xml:"client_id" form:"client_id"`
}

type GenerateCounterKey struct {
	ClientID   string `json:"client_id" xml:"client_id" form:"client_id"`
	CreateDate string `json:"create_date" xml:"create_date" form:"create_date"`
}
