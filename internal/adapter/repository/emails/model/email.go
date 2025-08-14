package model

import "time"

type Email struct {
	Id        uint      `gorm:"primarykey"`
	FromEmail string    `gorm:"not null"`
	FromName  string    `gorm:"not null"`
	Subject   string    `gorm:"not null"`
	ToEmail   string    `gorm:"not null"`
	Html      string    `gorm:"not null"`
	Text      string    `gorm:"not null"`
	Status    string    `gorm:"not null"`
	MessageId *string   `gorm:""`
	Errors    *string   `gorm:""`
	Created   time.Time `gorm:"not null"`
}
