package config

import (
	"github.com/spf13/viper"
	"os"
)

type Config struct {
	Port                string   `json:"port"`
	HealthCheckInterval string   `json:"healthCheckInterval"`
	Servers             []string `json:"servers"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigType("yaml")
	if os.Getenv("DOCKER_ENV") == "true" {
		viper.AddConfigPath("/app/config/")
		viper.SetConfigName("testconfig")
	} else {
		viper.AddConfigPath("./config/")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
