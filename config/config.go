package config

import (
	"github.com/spf13/viper"
)

func Init(configPath string) error {
	loadDefaultConfig(configPath)
	return viper.MergeInConfig()
}

func loadDefaultConfig(configPath string) error {
	viper.SetConfigName("default")
	viper.SetConfigType("json")
	viper.AddConfigPath(configPath)
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	return nil
}
