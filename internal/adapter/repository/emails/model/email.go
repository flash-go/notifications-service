package model

import "time"

type Email struct {
	Id          uint `gorm:"primarykey"`
	FolderId    *uint
	Folder      *EmailFolder `gorm:"foreignKey:FolderId;references:Id"`
	FromEmail   string       `gorm:"not null"`
	FromName    string       `gorm:"not null"`
	Subject     string       `gorm:"not null"`
	Html        string       `gorm:"not null"`
	Text        string       `gorm:"not null"`
	Description string       `gorm:"not null"`
	SystemFlag  bool         `gorm:"not null"`
	Updated     time.Time    `gorm:"not null"`
	Created     time.Time    `gorm:"not null"`
}
