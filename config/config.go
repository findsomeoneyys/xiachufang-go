package config

import (
	"log"
	"strings"
	"time"

	"github.com/spf13/viper"
)

var config *Config

type Config struct {
	Server Server
}

type Server struct {
	RunMode      string
	HttpPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func NewConfig() *Config {
	return &Config{
		Server: Server{},
	}
}

// viper环境变量覆盖时需要全大写 eg: SERVER_RUNMODE="release"
func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath("./config")
	viper.SetConfigType("yaml")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	config = NewConfig()
	err := viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

}

func Get() *Config {
	return config
}
