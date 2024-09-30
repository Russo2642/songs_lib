package config

import "github.com/spf13/viper"

func InitConfig(configPath string, configName string) error {
	viper.AddConfigPath(configPath)
	viper.SetConfigName(configName)
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	return nil
}
