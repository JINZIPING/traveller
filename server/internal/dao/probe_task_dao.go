package dao

import (
	"log"
	"my_project/pkg/model"
	"my_project/server/internal/infra"
)

// StoreICMPProbeTask 存储探测任务的元数据到MySQL
func StoreICMPProbeTask(task *model.ICMPProbeTask) error {
	stmt, err := infra.MySQLDB.Prepare("INSERT INTO my_database.probe_tasks (ip, count, port, threshold, timeout, status, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Printf("Error preparing statement: %v", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(task.IP, task.Count, task.Port, task.Threshold, task.Timeout, "pending", task.CreatedAt, task.UpdatedAt)
	if err != nil {
		log.Printf("Error executing ICMP statement: %v", err)
		return err
	}
	log.Printf("Successfully inserted ICMP probe task: %+v", task)
	return nil
}

// StoreTCPProbeTask 存储探测任务的元数据到MySQL
func StoreTCPProbeTask(task *model.TCPProbeTask) error {
	stmt, err := infra.MySQLDB.Prepare("INSERT INTO my_database.tcp_probe_tasks (ip, port, timeout, status, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Printf("Error preparing TCP statement: %v", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(task.IP, task.Port, task.Timeout, "pending", task.CreatedAt, task.UpdatedAt)
	if err != nil {
		log.Printf("Error executing TCP statement: %v", err)
		return err
	}
	log.Printf("Successfully inserted TCP probe task: %+v", task)
	return nil
}
