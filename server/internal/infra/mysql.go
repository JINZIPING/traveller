package infra

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

var MySQLDB *sql.DB

// InitMySQL 初始化 MySQL 连接 & 表结构
func InitMySQL() *sql.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		viper.GetString("mysql.user"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.host"),
		viper.GetInt("mysql.port"),
		viper.GetString("mysql.database"),
	)

	var err error
	for i := 0; i < 5; i++ {
		MySQLDB, err = sql.Open("mysql", dsn)
		if err != nil {
			log.Printf("Failed to connect to MySQL (attempt %d): %v", i+1, err)
			time.Sleep(2 * time.Second)
			continue
		}

		if err = MySQLDB.Ping(); err != nil {
			log.Printf("Failed to ping MySQL (attempt %d): %v", i+1, err)
			time.Sleep(2 * time.Second)
			continue
		}

		log.Println("MySQL initialized successfully")
		break
	}

	if err != nil {
		log.Fatalf("Could not establish connection to MySQL after retries: %v", err)
	}

	// 初始化表结构
	createTableQuery := `
    CREATE TABLE IF NOT EXISTS probe_tasks (
        id INT AUTO_INCREMENT PRIMARY KEY,
        ip VARCHAR(255) NOT NULL,
        count INT NOT NULL,
        port INT DEFAULT 0,
        threshold INT NOT NULL,
        timeout INT NOT NULL,
        status VARCHAR(50) NOT NULL DEFAULT 'pending',
        created_at DATETIME NOT NULL,
        updated_at DATETIME NOT NULL
    );`

	if _, err := MySQLDB.Exec(createTableQuery); err != nil {
		log.Fatalf("Failed to create table probe_tasks: %v", err)
	}
	log.Println("Table probe_tasks ensured.")

	createTcpTableQuery := `
	CREATE TABLE IF NOT EXISTS tcp_probe_tasks (
	    id INT AUTO_INCREMENT PRIMARY KEY,
	    ip VARCHAR(255) NOT NULL,
		port INT DEFAULT 0,
	    timeout INT NOT NULL,
	    status VARCHAR(50) NOT NULL DEFAULT 'pending',
	    created_at DATETIME NOT NULL,
        updated_at DATETIME NOT NULL
	);`
	if _, err := MySQLDB.Exec(createTcpTableQuery); err != nil {
		log.Fatalf("Failed to create table tcp_probe_tasks: %v", err)
	}
	log.Println("Table tcp_probe_tasks ensured.")

	return MySQLDB
}
