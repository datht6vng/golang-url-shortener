package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

type LimitterConfig struct {
	MaxRequest      string `mapstructure:"max_request"`
	LimitterExprire string `mapstructure:"limitter_expire"`
}
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
}
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Password string `mapstructure:"password"`
}
type ServerConfig struct {
	Port   string `mapstructure:"port"`
	Domain string `mapstructure:"domain"`
}
type AppConfig struct {
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Server   ServerConfig   `mapstructure:"server"`
	Limitter LimitterConfig `mapstructure:"limitter"`
}

var (
	Config AppConfig
)

func ReadConfig() {
	viper := viper.New()
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath("./config")
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Cannot load config: %s", err)
		os.Exit(1)
	}
	if err := viper.Unmarshal(&Config); err != nil {
		log.Println("Cannot read config file!")
	}
}
