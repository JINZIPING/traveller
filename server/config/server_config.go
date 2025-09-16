package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

// Config 结构体保存配置文件中解析的值
type Config struct {
	MySQL struct {
		User     string
		Password string
		Host     string
		Port     int
		Database string
	}
	ClickHouse struct {
		Host string
		Port int
	}
	LogFile string
}

// 全局变量，用于存储加载的配置
var ServerConfig Config

// InitConfig 初始化配置文件
func InitConfig() {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath != "" {
		viper.SetConfigName("server_config")
		viper.AddConfigPath("/root/config")
		//viper.AddConfigPath("./config")
		viper.SetConfigType("yaml")
	} else {
		viper.SetConfigName("server_config.local")
		viper.AddConfigPath("./server/config")
		viper.SetConfigType("yaml")
	}

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	err = viper.Unmarshal(&ServerConfig)
	if err != nil {
		log.Fatalf("Unable to decode into struct: %v", err)
	}

	log.Println("Config file loaded successfully")
}
