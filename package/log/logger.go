package logger

import (
	"io"
	"log"
	"os"
	"sync"
	config2 "trueid-shorten-link/package/config"

	"github.com/evalphobia/logrus_sentry"
	"github.com/sirupsen/logrus"
)

var (
	Log  *logrus.Logger
	once sync.Once
)

const (
	timeFormat = "2006-01-02T15:04:05.000"
)

type Fields map[string]interface{}

func InitLogger(cfg *config2.LogConfig) *logrus.Logger {
	once.Do(func() {
		Log = logrus.New()
		Log.SetFormatter(&TextFormatter{
			Prefix:          os.Getenv("HOSTNAME"),
			ForceFormatting: true,
			FullTimestamp:   true,
			TimestampFormat: timeFormat,
		})

		logLevel := logrus.DebugLevel
		switch cfg.Level {
		case "trace":
			logLevel = logrus.TraceLevel
		case "debug":
			logLevel = logrus.DebugLevel
		case "info":
			logLevel = logrus.InfoLevel
		case "warn":
			logLevel = logrus.WarnLevel
		case "error":
			logLevel = logrus.ErrorLevel
		}

		Log.Level = logLevel

		file := &Logger{
			Filename:   cfg.FilePath,
			MaxSize:    cfg.MaxSize, // megabytes
			MaxBackups: cfg.MaxBackup,
			MaxAge:     cfg.MaxAge, //days
		}
		// defer file.Close()

		multi := io.MultiWriter(file, os.Stdout)
		log.SetOutput(multi)

		if cfg.SentryURL != "" {
			hook, err := logrus_sentry.NewSentryHook(cfg.SentryURL, []logrus.Level{
				logrus.PanicLevel,
				logrus.FatalLevel,
				logrus.ErrorLevel,
				logrus.WarnLevel,
			})

			if err == nil {
				Log.Hooks.Add(hook)
			}
		}
	})
	return Log
}

func GetLog() *logrus.Logger {
	return Log
}
