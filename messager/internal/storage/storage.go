package storage

type ChatRepository interface {
	CreateChat(userID, userID2 string) error
	DeleteChat(messageID string) error
}
