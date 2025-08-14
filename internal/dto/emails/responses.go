package dto

import (
	"time"
)

type EmailResponse struct {
	Id        uint      `json:"id"`
	FromEmail string    `json:"from_email"`
	FromName  string    `json:"from_name"`
	Subject   string    `json:"subject"`
	ToEmail   string    `json:"to_email"`
	Html      string    `json:"html"`
	Text      string    `json:"text"`
	Status    string    `json:"status"`
	MessageId *string   `json:"message_id"`
	Errors    *string   `json:"errors"`
	Created   time.Time `json:"created"`
}
