package config

import (
	"github.com/spf13/viper"
)

type BasicConfig struct {
	basic *viper.Viper
}

func NewBasicConfig() *BasicConfig {
	basicViper := viper.New()
	basicViper.SetDefault("cloud.enable", false)
	basicViper.SetDefault("cloud.backend_url", "")
	basicViper.SetDefault("cloud.app_id", "")
	basicViper.SetDefault("cloud.app_secret", "")
	basicViper.SetDefault("cloud.heartbeat_interval", 180)
	bc := &BasicConfig{
		basic: basicViper,
	}
	return bc
}

func (bc *BasicConfig) GetBool(key string) bool {
	return bc.basic.GetBool(key)
}

func (bc *BasicConfig) GetString(key string) string {
	return bc.basic.GetString(key)
}

func (bc *BasicConfig) LoadProperties(path string) error {
	bc.basic.SetConfigType("yaml")
	bc.basic.SetConfigFile(path)
	err := bc.basic.ReadInConfig()
	if err != nil {
		return err
	} else {
		return nil
	}
}
