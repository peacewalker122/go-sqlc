package util

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DBDriver          string        `mapstructure:"DB_DRIVER"`
	DBSource          string        `mapstructure:"DB_SOURCE"`
	HTTPServerAddress string        `mapstructure:"HTTP_SERVER_ADDRESS"`
	GRPCServerAddress string        `mapstructure:"GRPC_SERVER_ADDRESS"`
	SymmectricKey     string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	Duration          time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshToken      time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	// i don't know why viper can't read env file so i tried to change it into json format and it works properly.
	viper.SetConfigType("json")

	log.Println("file: ", viper.ConfigFileUsed())

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}
