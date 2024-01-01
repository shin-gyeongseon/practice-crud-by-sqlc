package util

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

// Config stores all configuration of the application.
// The Values are read by viper from a config file or environment variables.
type Config struct {
	DBDriver            string        `mapstructure:"DB_DRIVER"`
	DBSource            string        `mapstructure:"DB_SOURCE"`
	ServerAddress       string        `mapstructure:"SERVER_ADDRESS"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	TokenSymmenticKey   string        `mapstructure:"TOKEN_SYMMENTIC_KEY"`
}

// LoadConfig reads configuration from file or environment variables
func LoadConfig(path string) (config Config, err error) {
	viper.SetConfigName("app")
	viper.AddConfigPath(path)
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	err = viper.Unmarshal(&config)
	return
}
