package model

import "time"

type ICMPProbeTask struct {
	ID        int64  `gorm:"primaryKey;autoIncrement"`
	IP        string `gorm:"type:varchar(255);not null"`
	Count     int    `gorm:"not null"`
	Threshold int    `gorm:"not null"`
	Timeout   int    `gorm:"not null"`
	Status    string `gorm:"type:varchar(50);default:'pending'"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
