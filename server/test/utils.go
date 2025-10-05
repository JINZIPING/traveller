package test

import (
	"testing"

	"github.com/spf13/viper"
)

// InitTestConfig
func InitTestConfig(t *testing.T) {
	viper.SetConfigName("server_config.test")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../../server/config")

	if err := viper.ReadInConfig(); err != nil {
		t.Fatalf("failed to read config file: %v", err)
	}
}
