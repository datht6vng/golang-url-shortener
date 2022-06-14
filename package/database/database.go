package database

import (
	"fmt"
	"trueid-shorten-link/package/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect() *gorm.DB {
	dsn := fmt.Sprintf(
		"%s:%s@%s(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&multiStatements=true",
		config.Config.Database.User,
		config.Config.Database.Password,
		"tcp",
		config.Config.Database.Host,
		config.Config.Database.Port,
		config.Config.Database.Database,
	)

	connection, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	db, err := connection.DB()
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	fmt.Println("Connected to Database!")
	db.SetMaxOpenConns(config.Config.Database.MaxOpenConnection)
	db.SetMaxIdleConns(config.Config.Database.MaxIdleConnection)
	return connection
}
