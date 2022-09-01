package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

var config = new(Config)

type Config struct {
	HttpPort     int
	RunMode      string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	err := viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

}

func Get() Config {
	return *config
}
