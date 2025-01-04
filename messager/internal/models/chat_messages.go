package models

import "time"

type ChatMessagesModel struct {
	Chat_id      string
	Message_id   string
	Sender_id    string
	Message_text string
	Sent_at      time.Time
}
