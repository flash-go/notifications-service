package model

import "time"

type EmailLog struct {
	Id        uint   `gorm:"primarykey"`
	FromEmail string `gorm:"not null"`
	FromName  string `gorm:"not null"`
	Subject   string `gorm:"not null"`
	ToEmail   string `gorm:"not null"`
	Html      string `gorm:"not null"`
	Text      string `gorm:"not null"`
	Status    string `gorm:"not null"`
	MessageId *string
	Errors    *string
	Created   time.Time `gorm:"not null"`
}
