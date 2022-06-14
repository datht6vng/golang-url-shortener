package model

type Client struct {
	ClientID string `json:"client_id" xml:"client_id" form:"client_id" gorm:"client_id;primaryKey"`
	APIKey   string `json:"api_key" xml:"api_key" form:"api_key" gorm:"api_key"`
}

func (m *Client) TableName() string {
	return `shorten_link_clients`
}
