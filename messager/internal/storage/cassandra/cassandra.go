package cassandra

import (
	"github.com/gocql/gocql"
	"github.com/sirupsen/logrus"
	"gommessage.com/messager/internal/models"
	"gommessage.com/messager/internal/storage"
	database "gommessage.com/messager/pkg/Database"
)

type ChatRepo struct {
	session *gocql.Session
}

// CreateChat implements storage.ChatRepository.
func (c *ChatRepo) CreateChat(userID1 string, userID2 string) error {
	q := `INSERT INTO chats(chat_id ,user_id1, user_id2) VALUES(now(),?,?)`
	if err := database.Exec(q, userID1, userID2); err != nil {
		logrus.Fatalf("Error while creating chat.%w", err)
		return err
	}
	return nil
}

// DeleteChat implements storage.ChatRepository.
func (c *ChatRepo) DeleteChat(userID1 string, userID2 string) error {
	q := `DELETE FROM chats WHERE user_id1 = ? AND user_id2 = ?`
	if err := database.Exec(q, userID1, userID2); err != nil {
		logrus.Fatalf("Error while deleting chat.%w", err)
		return err
	}
	return nil
}

// DeleteMessage implements storage.ChatRepository.
func (c *ChatRepo) DeleteMessage(chatID string, messageID string) error {
	query := `DELETE FROM chats_messages 
              WHERE chat_id = ? AND message_id = ?`

	if err := database.Exec(query, chatID, messageID); err != nil {
		logrus.Errorf("Помилка при видаленні повідомлення: %v", err)
		return err
	}
	return nil
}

// GetMessageHistory implements storage.ChatRepository.
func (c *ChatRepo) GetMessageHistory(chatID string) ([]models.ChatMessagesModel, error) {
	query := `SELECT chat_id, message_id, sender_id, message_text, sent_at 
	FROM chats_messages 
	WHERE chat_id = ? 
	ORDER BY sent_at DESC`

	var messages []models.ChatMessagesModel
	err := database.Exec(query, chatID)
	if err != nil {
		logrus.Error("Blya pomylka")
	}
	// Тут має бути код для виконання запиту до бази даних
	// та заповнення масиву messages

	return messages, nil
}

// GetUserChats implements storage.ChatRepository.
func (c *ChatRepo) GetUserChats(userID string) ([]models.ChatsModel, error) {
	panic("unimplement")
}

// SendMessage implements storage.ChatRepository.
func (c *ChatRepo) SendMessage(message *models.ChatMessagesModel) error {
	query := `INSERT INTO chats_messages(chat_id, message_id,sender_id, message_text, sent_at)
	VALUES(?,?,?,?,?)`
	logrus.Info("Ahuenno querry")
	database.Exec(query, message.Chat_id, message.Message_id, message.Sender_id, message.Message_text, message.Sent_at)
	return nil
}

func NewChatRepository(session *gocql.Session) storage.ChatRepository {
	return &ChatRepo{
		session: session,
	}
}
