package infra

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var MySQLDB *gorm.DB

func InitMySQL() *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		viper.GetString("mysql.user"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.host"),
		viper.GetInt("mysql.port"),
		viper.GetString("mysql.database"),
	)

	var err error
	for i := 0; i < 5; i++ {
		MySQLDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info), // 打印SQL
		})
		if err == nil {
			sqlDB, _ := MySQLDB.DB()
			sqlDB.SetMaxOpenConns(50)
			sqlDB.SetMaxIdleConns(10)
			sqlDB.SetConnMaxLifetime(5 * time.Minute)
			log.Println("MySQL initialized successfully (GORM)")
			break
		}
		log.Printf("Failed to connect to MySQL (attempt %d): %v", i+1, err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatalf("Could not establish connection to MySQL after retries: %v", err)
	}

	return MySQLDB
}
