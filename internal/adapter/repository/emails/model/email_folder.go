package model

import "time"

type EmailFolder struct {
	Id          uint `gorm:"primaryKey"`
	ParentId    *uint
	Parent      *EmailFolder `gorm:"foreignKey:ParentId;references:Id"`
	Name        string       `gorm:"not null"`
	Description string       `gorm:"not null"`
	SystemFlag  bool         `gorm:"not null"`
	Updated     time.Time    `gorm:"not null"`
	Created     time.Time    `gorm:"not null"`
}
