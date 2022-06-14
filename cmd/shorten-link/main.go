package main

import (
	"os"
	"os/signal"

	"trueid-shorten-link/internal/shorten-link/interface/http"
	"trueid-shorten-link/package/config"
	_log "trueid-shorten-link/package/log"

	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load(".env")
	configPath := os.Getenv("CONFIG_PATH")
	err := config.ReadConfig(configPath)
	if err != nil {
		panic(err)
	}
}

func main() {
	logger := _log.InitLogger(&config.Config.Logger)
	logger.Info("Service start!")

	defer func() {
		err := recover()
		if err != nil {
			logger.Info(err)
		}
		logger.Info("Service end!")
	}()

	// Catch Ctr + C
	go func() {
		signalChannel := make(chan os.Signal)
		signal.Notify(signalChannel, os.Interrupt)
		<-signalChannel
		logger.Info("Service end!")
		os.Exit(0)
	}()

	app := http.NewApp()
	app.Listen(":" + config.Config.App.Port)
}
