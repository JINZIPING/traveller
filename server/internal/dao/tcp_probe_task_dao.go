package dao

import (
	"my_project/server/internal/infra"
	"time"

	pkgModel "my_project/pkg/model"
	dbModel "my_project/server/internal/model"
)

// StoreTCPProbeTask 插入 TCP 任务（DTO → ORM）
func StoreTCPProbeTask(dto *pkgModel.TCPProbeTaskDTO) error {
	task := &dbModel.TCPProbeTask{
		IP:        dto.IP,
		Port:      dto.Port,
		Timeout:   dto.Timeout,
		Status:    "pending",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return infra.MySQLDB.Create(task).Error
}

// GetTCPProbeTaskByIP 根据 IP 查询 TCP 任务
func GetTCPProbeTaskByIP(ip string) (*dbModel.TCPProbeTask, error) {
	var task dbModel.TCPProbeTask
	if err := infra.MySQLDB.Where("ip = ?", ip).First(&task).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

// DeleteTCPProbeTaskByIP 根据 IP 删除 TCP 任务
func DeleteTCPProbeTaskByIP(ip string) error {
	return infra.MySQLDB.Where("ip = ?", ip).Delete(&dbModel.TCPProbeTask{}).Error
}
