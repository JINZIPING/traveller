package dao

import (
	"database/sql"
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

// GetICMPProbeTaskByIP retrieves an ICMP probe task by IP
func GetICMPProbeTaskByIP(ip string) (*model.ICMPProbeTask, error) {
	query := `
        SELECT ip, count, port, threshold, timeout, created_at, updated_at
        FROM probe_tasks WHERE ip = ?
    `
	row := infra.MySQLDB.QueryRow(query, ip)

	var task model.ICMPProbeTask
	err := row.Scan(
		&task.IP,
		&task.Count,
		&task.Port,
		&task.Threshold,
		&task.Timeout,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Printf("Error scanning ICMP task: %v", err)
		return nil, err
	}
	return &task, nil
}

// DeleteICMPProbeTaskByIP deletes an ICMP probe task by IP
func DeleteICMPProbeTaskByIP(ip string) error {
	stmt, err := infra.MySQLDB.Prepare(`DELETE FROM probe_tasks WHERE ip = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(ip)
	if err != nil {
		log.Printf("Error deleting ICMP task for IP %s: %v", ip, err)
		return err
	}
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

// GetTCPProbeTaskByIP retrieves a TCP probe task by IP
func GetTCPProbeTaskByIP(ip string) (*model.TCPProbeTask, error) {
	query := `
        SELECT ip, port, timeout, created_at, updated_at
        FROM tcp_probe_tasks WHERE ip = ?
    `
	row := infra.MySQLDB.QueryRow(query, ip)

	var task model.TCPProbeTask
	err := row.Scan(
		&task.IP,
		&task.Port,
		&task.Timeout,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &task, nil
}

// DeleteTCPProbeTaskByIP deletes a TCP probe task by IP
func DeleteTCPProbeTaskByIP(ip string) error {
	stmt, err := infra.MySQLDB.Prepare(`DELETE FROM tcp_probe_tasks WHERE ip = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(ip)
	return err
}
