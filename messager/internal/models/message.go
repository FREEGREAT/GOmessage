package models

import (
	"time"
)

type MessageModel struct {
	Message_id string    `json:"message_id"`
	User_id1   string    `json:"user_1id"`
	User_id2   string    `json:"user_2id"`
	Message    string    `json:"message"`
	Sent_time  time.Time `json:"sent_time"`
	Is_edited  bool      `json:"is_edited"`
}
