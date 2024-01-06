package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func ReadConfig() {
	viper.AddConfigPath("./config")
	viper.SetConfigName("config") // Register config file name (no extension)
	viper.SetConfigType("yaml")   // Look for specific type
	err := viper.ReadInConfig()

	if err != nil {
		fmt.Println(err.Error())
	}
}
