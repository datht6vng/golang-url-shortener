package model

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var Connection *sql.DB

type Url struct {
	Url string `json:"url" xml:"url" form:"url"`
}
type UrlRecord struct {
	ShortUrl   string    `json:"shortUrl" xml:"shortUrl" form:"shortUrl"`
	LongUrl    string    `json:"longUrl"  xml:"longUrl"  form:"longUrl"`
	ExpireTime time.Time `json:"expireTime" xml:"expireTime" form:"expireTime"`
}

func Connect() {
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "127.0.0.1"
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
		database = "mysql"
	}
	dbConneectionString := fmt.Sprintf(
		"%s:%s@%s(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user,
		password,
		"tcp",
		host,
		port,
		database,
	)
	db, err := sql.Open("mysql", dbConneectionString)
	if err != nil {
		fmt.Println("Failed to open database", err.Error())
		return
	}
	if err = db.Ping(); err != nil {
		fmt.Println("Failed to connect to database", err.Error())
		return
	}
	fmt.Println("Open database")
	Connection = db
}
func CreateModel(connection *sql.DB) error {
	_, err := connection.Query(`
		CREATE TABLE IF NOT EXISTS URL (
			SHORT_URL VARCHAR(500) PRIMARY KEY,
			LONG_URL VARCHAR(500),
			EXPIRE_TIME TIMESTAMP
		);
	`)
	return err
}
func FindLongUrl(connection *sql.DB, url string) (UrlRecord, error) {
	result := UrlRecord{}
	err := connection.QueryRow("SELECT * FROM URL WHERE LONG_URL = ?", url).Scan(&result.ShortUrl, &result.LongUrl, &result.ExpireTime)
	return result, err
}
func FindShortUrl(connection *sql.DB, url string) (UrlRecord, error) {
	result := UrlRecord{}
	err := connection.QueryRow("SELECT * FROM URL WHERE LONG_URL = ?", url).Scan(&result.ShortUrl, &result.LongUrl, &result.ExpireTime)
	return result, err
}
func InsertUrl(connection *sql.DB, shortUrl string, longUrl string, expiryTime time.Time) error {
	_, err := connection.Query("INSERT INTO URL VALUES (?, ?, ?)", shortUrl, longUrl, expiryTime)
	return err
}
