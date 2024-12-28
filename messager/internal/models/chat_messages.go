package models

import "time"

type ChatMessagesModel struct{
	chat_id string
	message_id string
	sender_id string
	message_text string
	sent_at time.Time
}