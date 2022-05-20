package model

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Url struct {
	Url string `json:"url" xml:"url" form:"url"`
}
type UrlRecord struct {
	ID         int       `json:"ID" xml:"ID" form:"ID"`
	ShortUrl   string    `json:"shortUrl" xml:"shortUrl" form:"shortUrl"`
	LongUrl    string    `json:"longUrl"  xml:"longUrl"  form:"longUrl"`
	ExpireTime time.Time `json:"expireTime" xml:"expireTime" form:"expireTime"`
	UsedCount  int       `json:"usedCount" xml:"usedCount" form:"usedCount"`
}
type Model struct {
	connection *sql.DB
	currentID  int64
}

func (this *Model) Connect() {
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}
	user := os.Getenv("DB_USER")
	if user == "" {
		user = "root"
	}
	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "123456"
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "3306"
	}
	database := os.Getenv("DB_DATABASE")
	if database == "" {
		database = "URL_SHORTENER"
	}
	dbConnectionString := fmt.Sprintf(
		"%s:%s@%s(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&multiStatements=true",
		user,
		password,
		"tcp",
		host,
		port,
		database,
	)
	db, err := sql.Open("mysql", dbConnectionString)
	if err != nil {
		fmt.Println("Failed to open database", err.Error())
		return
	}
	if err = db.Ping(); err != nil {
		fmt.Println("Failed to connect to database", err.Error())
		return
	}
	fmt.Println("Open database")
	db.SetMaxOpenConns(9999)
	this.connection = db
	this.currentID, err = this.GetMaxID()

}
func (this *Model) Close() error {
	return this.connection.Close()
}
func (this *Model) CreateModel() error {
	_, err := this.connection.Query(`
		CREATE TABLE IF NOT EXISTS URL (
			ID BIGINT PRIMARY KEY,
			SHORT_URL VARCHAR(500),
			LONG_URL VARCHAR(500),
			EXPIRE_TIME TIMESTAMP,
			USED_COUNT INT
		);
		ALTER TABLE URL ADD INDEX INDEX_ID USING BTREE (ID);
		ALTER TABLE URL ADD INDEX INDEX_SHORT_URL USING HASH (SHORT_URL);
		ALTER TABLE URL ADD INDEX INDEX_LONG_URL USING HASH (LONG_URL);
	`)
	return err
}
func (this *Model) FindLongUrl(url string) (*UrlRecord, error) {
	result := new(UrlRecord)
	err := this.connection.QueryRow("SELECT * FROM URL WHERE LONG_URL = ?", url).Scan(&result.ID, &result.ShortUrl, &result.LongUrl, &result.ExpireTime, &result.UsedCount)
	return result, err
}
func (this *Model) FindShortUrl(url string) (*UrlRecord, error) {
	result := new(UrlRecord)
	err := this.connection.QueryRow("SELECT * FROM URL WHERE SHORT_URL = ?", url).Scan(&result.ID, &result.ShortUrl, &result.LongUrl, &result.ExpireTime, &result.UsedCount)
	return result, err
}
func (this *Model) GetMaxID() (int64, error) {
	var result int64
	err := this.connection.QueryRow("SELECT MAX(ID) FROM URL").Scan(&result)
	return result, err
}
func (this *Model) GetNextID() int64 {
	this.currentID += 1
	return this.currentID
}
func (this *Model) InsertUrl(id int64, shortUrl string, longUrl string, expireTime time.Time, usedCount int) error {
	_, err := this.connection.Query("INSERT INTO URL VALUES (?, ?, ?, ?, ?)", id, shortUrl, longUrl, expireTime, usedCount)
	return err
}
func (this *Model) DeleteUrl(shortUrl string, longUrl string) error {
	var err error
	if shortUrl == "" && longUrl == "" {
		_, err = this.connection.Query("DELETE FROM URL")
	} else if longUrl == "" {
		_, err = this.connection.Query("DELETE FROM URL WHERE SHORT_URL = ?", shortUrl)
	} else if shortUrl == "" {
		_, err = this.connection.Query("DELETE FROM URL WHERE LONG_URL = ?", longUrl)
	} else {
		_, err = this.connection.Query("DELETE FROM URL WHERE SHORT_URL = ? AND LONG_URL = ?", shortUrl, longUrl)
	}
	return err
}
