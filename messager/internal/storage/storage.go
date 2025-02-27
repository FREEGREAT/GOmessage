package storage

import "gommessage.com/messager/internal/models"

type ChatRepository interface {
	SendMessage(message *models.ChatMessagesModel) error
	GetChatInfo(id string) (*models.ChatsModel, error)
	GetMessageById(chatID string, messageID string) (message *models.ChatMessagesModel, err error)
	CreateChat(userID1, userID2 string) error
	DeleteChat(userID1, userID2 string) error
	GetMessageHistory(chatID string) ([]models.ChatMessagesModel, error)
	DeleteMessage(chatID, messageID string) error
	GetUserChats(userID string) ([]models.ChatsModel, error)
}
