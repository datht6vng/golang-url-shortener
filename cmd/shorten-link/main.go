package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"trueid-shorten-link/config"
	"trueid-shorten-link/internal/shorten-link/interface/http"
	_log "trueid-shorten-link/log"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load("./cmd/shorten-link/.env"); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	configPath := os.Getenv("CONFIG_PATH")
	config.ReadConfig(configPath)

	// Set logger file
	_log.InitLogger()
}
func main() {
	log.Println("Service start!")
	defer func() {
		err := recover()
		if err != nil {
			log.Println(err)
		}
		log.Println("Service end!")
		_log.Logger.Close()
	}()
	// Catch Ctr + C
	go func() {
		signalChannel := make(chan os.Signal)
		signal.Notify(signalChannel, os.Interrupt)
		<-signalChannel
		log.Println("Service end!")
		_log.Logger.Close()
		os.Exit(0)
	}()
	app := http.NewApp()
	app.Listen(":" + config.Config.App.Port)
}
