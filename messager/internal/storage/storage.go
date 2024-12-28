package storage

import "gommessage.com/messager/internal/models"

type ChatRepository interface {
	CreateChat(message *models.ChatsModel) error
	DeleteChat(messageID string) error
}
