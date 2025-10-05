package dao

import (
	"log"
	"my_project/server/internal/infra"
)

func StoreClickHouse(timestamp int64, ip string, packetLoss, minRtt, maxRtt, avgRtt float64) error {
	// 开始事务
	tx, err := infra.ClickHouseDB.Begin()
	if err != nil {
		log.Printf("[ERROR CLICKHOUSE]: Failed to begin transaction: %v", err)
		return err
	}

	// 准备插入语句
	stmt, err := tx.Prepare("INSERT INTO my_database.my_table (timestamp, ip, packet_loss, min_rtt, max_rtt, avg_rtt) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Printf("[ERROR CLICKHOUSE]: Error preparing statement: %v", err)
		return err
	}
	defer stmt.Close()

	// 执行插入
	_, err = stmt.Exec(timestamp, ip, packetLoss, minRtt, maxRtt, avgRtt)
	if err != nil {
		log.Printf("[ERROR CLICKHOUSE]: Error executing statement: %v", err)
		// 如果插入失败，回滚事务
		tx.Rollback()
		return err
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		log.Printf("[ERROR CLICKHOUSE]: Failed to commit transaction: %v", err)
		return err
	}

	log.Printf("[INFO CLICKHOUSE]: Successfully inserted into ClickHouse: timestamp=%d, ip=%s, packet_loss=%f, min_rtt=%f, max_rtt=%f, avg_rtt=%f", timestamp, ip, packetLoss, minRtt, maxRtt, avgRtt)
	return nil
}

// 存储TCP结果到clickhouse中
func StoreTCPResult(timestamp int64, ip, port string, rtt float64, success bool) error {
	tx, err := infra.ClickHouseDB.Begin()
	if err != nil {
		log.Printf("[ERROR CLICKHOUSE]: Failed to begin transaction: %v", err)
		return err
	}
	stmt, err := tx.Prepare("INSERT INTO my_database.tcp_results (timestamp, ip, port, rtt, success) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		log.Printf("[ERROR CLICKHOUSE]:Error preparing TCP result statement: %v", err)
		return err
	}
	defer stmt.Close()

	// 执行插入
	_, err = stmt.Exec(timestamp, ip, port, rtt, BoolToInt(success))
	if err != nil {
		log.Printf("[ERROR CLICKHOUSE]: Error executing TCP result statement: %v", err)
		// 如果插入失败，回滚事务
		tx.Rollback()
		return err
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		log.Printf("[ERROR CLICKHOUSE]: Failed to commit TCP result transaction: %v", err)
		return err
	}

	log.Printf("[INFO CLICKHOUSE]: Successfully inserted TCP result into ClickHouse: timestamp=%d, ip=%s, port=%s, rtt=%f", timestamp, ip, port, rtt)
	return nil
}

func BoolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
