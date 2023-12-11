package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

const (
	DEFAULT_CONFIG = ".env"
)

type Config struct {
	Database struct {
		Host     string `mapstructure:"HOST"`
		Port     int    `mapstructure:"PORT"`
		User     string `mapstructure:"USER"`
		Password string `mapstructure:"PASSWORD"`
		DBName   string `mapstructure:"DBNAME"`
		SSLMode  string `mapstructure:"SSLMODE"`
	} `mapstructure:"DATABASE"`
}

func LoadConfig(path string) (*Config, error) {
	viper.SetConfigFile(getEnv(path))

	viper.AutomaticEnv()
	viper.SetConfigType("env")

	// Read the config file
	if err := viper.ReadInConfig(); err != nil {
		log.Println("Error reading config file:", err, "Path:", viper.ConfigFileUsed())
		return nil, err
	}

	// Unmarshal the config data into the struct
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Println("Error unmarshaling config data:", err, "Path:", viper.ConfigFileUsed())
		return nil, err
	}

	return &config, nil
}

func getEnv(path string) string {
	if path == "" {
		if err := os.Chdir(".."); err != nil {
			log.Println("Error changing working directory:", err)
			return DEFAULT_CONFIG
		}
		return DEFAULT_CONFIG
	}
	return path
}
