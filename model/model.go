package model

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Url struct {
	User string `json:"user" xml:"user" form:"user"`
	Url  string `json:"url" xml:"url" form:"url"`
}
type UrlRecord struct {
	ID         int       `json:"ID" xml:"ID" form:"ID"`
	User       string    `json:"user" xml:"user" form:"user"`
	ShortUrl   string    `json:"shortUrl" xml:"shortUrl" form:"shortUrl"`
	LongUrl    string    `json:"longUrl"  xml:"longUrl"  form:"longUrl"`
	ExpireTime time.Time `json:"expireTime" xml:"expireTime" form:"expireTime"`
	//UsedCount  int       `json:"usedCount" xml:"usedCount" form:"usedCount"`
}
type Model struct {
	connection *sql.DB
}

func (this *Model) Connect() *Model {
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
	var err error
	this.connection, err = sql.Open("mysql", dbConnectionString)
	if err != nil {
		fmt.Println("Failed to open database", err.Error())
		return this
	}
	this.connection.SetMaxIdleConns(1000)
	this.connection.SetMaxOpenConns(1000)
	this.connection.SetConnMaxLifetime(10 * time.Second)
	if err = this.connection.Ping(); err != nil {
		fmt.Println("Failed to connect to database", err.Error())
		return this
	}
	fmt.Println("Open database")
	return this
}

func (this *Model) Close() error {
	return this.connection.Close()
}
func (this *Model) FindLongUrl(url string) (*UrlRecord, error) {
	result := new(UrlRecord)
	err := this.connection.QueryRow("SELECT * FROM URL WHERE LONG_URL = ?", url).Scan(&result.ID, &result.User, &result.ShortUrl, &result.LongUrl, &result.ExpireTime)
	return result, err
}
func (this *Model) FindShortUrl(url string) (*UrlRecord, error) {
	result := new(UrlRecord)
	err := this.connection.QueryRow("SELECT * FROM URL WHERE SHORT_URL = ?", url).Scan(&result.ID, &result.User, &result.ShortUrl, &result.LongUrl, &result.ExpireTime)
	return result, err
}

func (this *Model) GetMaxID() (string, error) {
	var result interface{}
	err := this.connection.QueryRow("SELECT MAX(ID) FROM URL").Scan(&result)
	if result == nil {
		return "0", err
	}
	return string(result.([]byte)), err
}

func (this *Model) InsertUrl(id string, user string, shortUrl string, longUrl string, expireTime time.Time) error {
	_, err := this.connection.Exec("INSERT INTO URL VALUES (?, ?, ?, ?, ?)", id, user, shortUrl, longUrl, expireTime)
	return err
}
func (this *Model) DeleteUrl(shortUrl string, longUrl string) error {
	var err error
	if shortUrl == "" && longUrl == "" {
		_, err = this.connection.Exec("DELETE FROM URL")
	} else if longUrl == "" {
		_, err = this.connection.Exec("DELETE FROM URL WHERE SHORT_URL = ?", shortUrl)
	} else if shortUrl == "" {
		_, err = this.connection.Exec("DELETE FROM URL WHERE LONG_URL = ?", longUrl)
	} else {
		_, err = this.connection.Exec("DELETE FROM URL WHERE SHORT_URL = ? AND LONG_URL = ?", shortUrl, longUrl)
	}
	return err
}
func (this *Model) DeleteExpiredRecord() error {
	_, err := this.connection.Exec("DELETE FROM URL WHERE EXPIRE_TIME < ?", time.Now())
	return err
}
