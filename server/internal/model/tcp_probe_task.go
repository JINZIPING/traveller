package model

import "time"

type TCPProbeTask struct {
	ID        int64  `gorm:"primaryKey;autoIncrement"`
	TaskID    string `gorm:"type:varchar(64);uniqueIndex"`
	IP        string `gorm:"type:varchar(255);not null"`
	Port      string `gorm:"type:varchar(255);not null"`
	Timeout   int
	Status    string `gorm:"type:varchar(50);default:'pending'"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
