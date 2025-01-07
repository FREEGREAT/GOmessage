package cassandra

import (
	"errors"
	"fmt"
	"net/http"
	"sort"

	"github.com/gocql/gocql"
	"github.com/sirupsen/logrus"
	"gommessage.com/messager/internal/models"
	"gommessage.com/messager/internal/storage"
	database "gommessage.com/messager/pkg/Database"
)

type ChatRepo struct {
	session *gocql.Session
}

// GetChatInfo implements storage.ChatRepository.
func (c *ChatRepo) GetChatInfo(id string) (*models.ChatsModel, error) {

	if id == "" {
		return nil, errors.New("chat ID cannot be empty")
	}

	query := `
        SELECT chat_id, user_id1, user_id2 
        FROM chats 
        WHERE chat_id = ?
		ALLOW FILTERING
    `
	var chatInfo models.ChatsModel

	err := c.session.Query(query, id).Scan(
		&chatInfo.Chat_id,
		&chatInfo.User_id1,
		&chatInfo.User_id2,
	)

	if err != nil {
		if err == gocql.ErrNotFound {
			logrus.Infof("Chat not found: %s", id)
			return nil, fmt.Errorf("chat with ID %s not found", id)
		}
		logrus.Errorf("Error retrieving chat info: %v", err)
		return nil, fmt.Errorf("failed to retrieve chat info: %w", err)

	}
	if chatInfo.Chat_id == "" {
		return nil, errors.New("retrieved chat info is invalid")
	}

	logrus.Infof("Retrieved chat info: %+v", chatInfo)
	return &chatInfo, nil

}

// GetMessageById implements storage.ChatRepository.
func (c *ChatRepo) GetMessageById(chatID string, messageID string) (message *models.ChatMessagesModel, err error) {
	panic("unimplemented")
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
	query := `DELETE FROM chat_messages
              WHERE chat_id = ? AND message_id = ?`
	chatUUID, err := gocql.ParseUUID(chatID)
	if err != nil {
		logrus.Error("Error while parsing chat id to uuid")
		return err
	}
	messageUUID, err := gocql.ParseUUID(messageID)
	if err != nil {
		logrus.Error("Error while parsing message id to uuid")
		return err
	}
	if err := database.Exec(query, chatUUID, messageUUID); err != nil {
		logrus.Errorf("Error while deleting message: %v", err)
		return err
	}
	return nil
}

func (c *ChatRepo) GetMessageHistory(chatID string) ([]models.ChatMessagesModel, error) {
	if chatID == "" {
		return nil, errors.New("chat ID cannot be empty")
	}

	query := `SELECT chat_id, message_id, sender_id, message_text, sent_at 
    FROM chat_messages 
    WHERE chat_id = ?`
	chatUUID, err := gocql.ParseUUID(chatID)
	if err != nil {
		logrus.Error("Error while parsing uuid for chat")
		return nil, err
	}
	iter := c.session.Query(query, chatUUID).Iter()

	var messages []models.ChatMessagesModel

	var message models.ChatMessagesModel
	for iter.Scan(
		&message.Chat_id,
		&message.Message_id,
		&message.Sender_id,
		&message.Message_text,
		&message.Sent_at,
	) {
		messages = append(messages, message)
	}
	logrus.Infof("Retrieved %d messages for chat ID %s", len(messages), chatID)
	if err := iter.Close(); err != nil {
		logrus.Errorf("Error retrieving message history: %v", err)
		return nil, err
	}
	sort.Slice(messages, func(i, j int) bool {
		return messages[i].Sent_at.Before(messages[j].Sent_at)
	})

	return messages, nil

}

// GetUserChats implements storage.ChatRepository.
func (c *ChatRepo) GetUserChats(userID string) ([]models.ChatsModel, error) {
	var chats []models.ChatsModel

	userUUID, err := gocql.ParseUUID(userID)
	if err != nil {
		return nil, http.ErrServerClosed
	}
	query1 := `SELECT chat_id FROM chats WHERE user_id1 = ? ALLOW FILTERING`
	iter1 := c.session.Query(query1, userUUID).Iter()
	var chatID string
	for iter1.Scan(&chatID) {
		chats = append(chats, models.ChatsModel{Chat_id: chatID})
	}
	iter1.Close()

	// Другий запит для user_id2
	query2 := `SELECT chat_id FROM chats WHERE user_id2 = ? ALLOW FILTERING`
	iter2 := c.session.Query(query2, userUUID).Iter()
	for iter2.Scan(&chatID) {
		chats = append(chats, models.ChatsModel{Chat_id: chatID})
	}
	iter2.Close()

	return chats, nil
}

// SendMessage implements storage.ChatRepository.
func (c *ChatRepo) SendMessage(message *models.ChatMessagesModel) error {
	query := `INSERT INTO chat_messages(chat_id, message_id,sender_id, message_text, sent_at)
	VALUES(?,?,?,?,?)`
	database.Exec(query, message.Chat_id, message.Message_id, message.Sender_id, message.Message_text, message.Sent_at)
	return nil
}

func NewChatRepository(session *gocql.Session) storage.ChatRepository {
	return &ChatRepo{
		session: session,
	}
}
