package entity

import (
	"time"
)

type Email struct {
	Id        uint
	FromEmail string
	FromName  string
	Subject   string
	ToEmail   string
	Html      string
	Text      string
	Status    string
	MessageId *string
	Errors    *string
	Created   time.Time
}
