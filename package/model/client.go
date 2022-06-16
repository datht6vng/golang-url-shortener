package model

type Client struct {
	ClientID    string `json:"client_id" xml:"client_id" form:"client_id" gorm:"client_id;primaryKey"`
	APIKey      string `json:"api_key" xml:"api_key" form:"api_key" gorm:"api_key"`
	LicienseKey string `json:"liciense_key" xml:"liciense_key" gorm:"liciense_key"`
	MaxLink     int64  `json:"max_link" xml:"max_link" gorm:"max_link"`
}

func (m *Client) TableName() string {
	return `shorten_link_clients`
}
