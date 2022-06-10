package log

import (
	"fmt"
	"log"
	"os"
	"trueid-shorten-link/config"
)

var (
	Logger *os.File
)

func InitLogger() {
	var err error
	Logger, err = os.OpenFile(config.Config.Logger.Path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	log.SetOutput(Logger)
}
