package util

import (
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	DbDriver            string        `mapstructure:"DB_DRIVER"`
	DbSource            string        `mapstructure:"DB_SOURCE"`
	ServerDriver        string        `mapstructure:"SERVER_ADDRESS"`
	TokenSymmetricKey   string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("tsconfig")
	viper.SetConfigType("json")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}
