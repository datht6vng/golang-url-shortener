package main

import (
	"os"
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
		logger.Infof("Service going down ... ")
		err := recover()
		if err != nil {
			logger.Info(err)
		}
		logger.Info("Service end!")
	}()

	app := http.NewApp()
	port := config.Config.App.Port
	logger.Infof("Start service on port %v (result: %v)", port, app.Listen(":"+port))
}
