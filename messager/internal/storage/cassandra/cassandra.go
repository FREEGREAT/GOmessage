package cassandra

import (
	"github.com/sirupsen/logrus"
	"gommessage.com/messager/internal/models"
	"gommessage.com/messager/internal/storage"
	database "gommessage.com/messager/pkg/Database"
)

type chatRepo struct {
	chRepo storage.ChatRepository
}

func NewStorage(repo storage.ChatRepository) *chatRepo {
	return &chatRepo{chRepo: repo}
}

func sendMessage(message *models.MessageModel) {
	query := `INSERT INTO messages(message_id,user1_id, user2_id, message)
	VALUES(now(),?,?,?)`
	logrus.Info("Ahuenno owou wou querry")
	database.Exec(query, message.User_id1, message.User_id2, message.Message)
}

func CreateChat(userID, userID2 string) error {
	q := `INSERT INTO chats(chat_id ,user_id1, user_id2) VALUES(now(),?,?)`
	if err := database.Exec(q, userID, userID2); err != nil {
		logrus.Fatalf("Error while creating chat.%w", err)
		return err
	}
	return nil

}
