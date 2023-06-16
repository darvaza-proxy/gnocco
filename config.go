package gnocco

import (
	"fmt"

	"github.com/spf13/viper"
)

// NewConfig creates a new config to be run
func NewConfig() {
	viper.SetConfigName("gnocco.conf")
	viper.AddConfigPath(".")
	viper.SetConfigType("toml")

	setDefaults()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err)
	}
}

func setDefaults() {
	viper.SetDefault("Host", "")
	viper.SetDefault("Port", 53)
	viper.SetDefault("Daemon", true)
}
