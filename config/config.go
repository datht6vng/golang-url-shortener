package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type LimitterConfig struct {
	MaxRequest      int `mapstructure:"max_request"`
	LimitterExprire int `mapstructure:"limitter_expire"`
}
type DatabaseConfig struct {
	Host              string `mapstructure:"host"`
	Port              string `mapstructure:"port"`
	User              string `mapstructure:"user"`
	Password          string `mapstructure:"password"`
	Database          string `mapstructure:"database"`
	MaxOpenConnection int    `mapstructure:"max_open_connection"`
	MaxIdleConnection int    `mapstructure:"max_idle_connection"`
}
type RedisConfig struct {
	Address   []string `mapstructure:"address"`
	Password  string   `mapstructure:"password"`
	Database  int      `mapstructure:"database"`
	IsCluster bool     `mapstructure:"is_cluster"`
}
type ServerConfig struct {
	Port   string `mapstructure:"port"`
	Domain string `mapstructure:"domain"`
}
type LoggerConfig struct {
	Path string `mapstructure:"path"`
}
type ViewConfig struct {
	Path string `mapstructure:"path"`
}
type KeyConfig struct {
	Secret string `mapstructure:"secret"`
}
type AppConfig struct {
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	App      ServerConfig   `mapstructure:"app"`
	Limitter LimitterConfig `mapstructure:"limitter"`
	Logger   LoggerConfig   `mapstructure:"logger"`
	View     ViewConfig     `mapstructure:"view"`
	Key      KeyConfig      `mapstructure:"key"`
}

var (
	Config AppConfig
)

func ReadConfig(configPath string) {
	viper := viper.New()
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Cannot load config: %s", err)
		os.Exit(1)
	}
	if err := viper.Unmarshal(&Config); err != nil {
		fmt.Println("Cannot read config file!")
	}
}
