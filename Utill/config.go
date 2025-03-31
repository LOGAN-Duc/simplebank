package Utill

import (
	"github.com/spf13/viper"
)

type Config struct {
	DbDriver     string `mapstructure:"DB_DRIVER"`
	DbSource     string `mapstructure:"DB_SOURCE"`
	ServerDriver string `mapstructure:"SERVER_ADDRESS"`
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
