package config

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	MongoURI string `mapstructure:"MONGO_URI"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Error reading config file", err)
	}
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal("Unable to decode config", err)
	}
	return &config, nil
}
