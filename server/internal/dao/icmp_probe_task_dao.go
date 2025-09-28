package dao

import (
	"my_project/server/internal/infra"
	"time"

	pkgModel "my_project/pkg/model"
	dbModel "my_project/server/internal/model"
)

// StoreICMPProbeTask 插入 ICMP 任务
func StoreICMPProbeTask(dto *pkgModel.ICMPProbeTaskDTO) error {
	task := &dbModel.ICMPProbeTask{
		IP:        dto.IP,
		Count:     dto.Count,
		Threshold: dto.Threshold,
		Timeout:   dto.Timeout,
		Status:    "pending",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return infra.MySQLDB.Create(task).Error
}

// GetICMPProbeTaskByIP 根据 IP 查询 ICMP 任务
func GetICMPProbeTaskByIP(ip string) (*dbModel.ICMPProbeTask, error) {
	var task dbModel.ICMPProbeTask
	if err := infra.MySQLDB.Where("ip = ?", ip).First(&task).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

// DeleteICMPProbeTaskByIP 根据 IP 删除 ICMP 任务
func DeleteICMPProbeTaskByIP(ip string) error {
	return infra.MySQLDB.Where("ip = ?", ip).Delete(&dbModel.ICMPProbeTask{}).Error
}
