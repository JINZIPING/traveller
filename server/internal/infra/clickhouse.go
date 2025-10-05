package infra

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/ClickHouse/clickhouse-go"
	"github.com/spf13/viper"
)

var ClickHouseDB *sql.DB

// InitClickHouse 初始化 ClickHouse 连接
func InitClickHouse() *sql.DB {
	dsn := fmt.Sprintf("tcp://%s:%d?debug=true",
		viper.GetString("clickhouse.host"),
		viper.GetInt("clickhouse.port"))

	var err error
	for i := 0; i < 5; i++ {
		ClickHouseDB, err = sql.Open("clickhouse", dsn)
		if err != nil {
			log.Fatalf("[ERROR CLICKHOUSE]: Failed to connect to ClickHouse: %v", err)
		} else {
			err = ClickHouseDB.Ping()
			if err == nil {
				log.Println("[INIT]: Successfully connected to ClickHouse")
				break
			}
			log.Printf("[ERROR CLICKHOUSE]: Failed to ping ClickHouse, attempt (%d/5): %v", i+1, err)
		}
		time.Sleep(1 * time.Second) // 等待 1 秒后重试
	}

	if err != nil {
		log.Fatalf("[ERROR CLICKHOUSE]: Failed to connect to ClickHouse after retries: %v", err)
	}

	// 确保数据库存在
	_, err = ClickHouseDB.Exec("CREATE DATABASE IF NOT EXISTS my_database")
	if err != nil {
		log.Fatalf("[ERROR CLICKHOUSE]: Error creating database: %v", err)
	}
	log.Println("[INIT]: ClickHouse database created successfully or already exists")

	// 初始化表结构（ICMP & TCP）
	_, err = ClickHouseDB.Exec(`
        CREATE TABLE IF NOT EXISTS my_database.my_table (
            timestamp DateTime,
            ip String,
            packet_loss Float64,
            min_rtt Float64,
            max_rtt Float64,
            avg_rtt Float64,
            threshold   Int32,
            success     UInt8 
        ) ENGINE = MergeTree()
        ORDER BY timestamp
    `)
	if err != nil {
		log.Fatalf("[ERROR CLICKHOUSE]: Failed to create table icmp_results: %v", err)
	}
	log.Println("[INIT]: ClickHouse ICMP results table created successfully")

	// 创建表用于TCP探测
	_, err = ClickHouseDB.Exec(`
		CREATE TABLE IF NOT EXISTS my_database.tcp_results(
		      timestamp DateTime,
		      ip String,
		      port String,
		      rtt Float64,
		      success UInt8
		) ENGINE = MergeTree()
		      ORDER BY timestamp
	`)
	if err != nil {
		log.Fatalf("[ERROR CLICKHOUSE]: Failed to create table tcp_results: %v", err)
	}

	log.Println("[INIT]: ClickHouse initialized successfully")
	return ClickHouseDB
}
