package ws

import (
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"gommessage.com/messager/internal/models"
	"gommessage.com/messager/internal/storage"
)

type Client struct {
	hub      *Hub
	conn     *websocket.Conn
	send     chan []byte
	userID   string
	chatID   string
	chatRepo storage.ChatRepository
}

type WSMessage struct {
	Type      string    `json:"type"`
	ChatID    string    `json:"chat_id"`
	MessageID string    `json:"message_id"`
	SenderID  string    `json:"sender_id"`
	Text      string    `json:"text"`
	SentAt    time.Time `json:"sent_at"`
}

// NewClient створює нового клієнта
func NewClient(hub *Hub, conn *websocket.Conn, userID, chatID string, repo storage.ChatRepository) *Client {
	return &Client{
		hub:      hub,
		conn:     conn,
		send:     make(chan []byte, 256),
		userID:   userID,
		chatID:   chatID,
		chatRepo: repo,
	}
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logrus.Error("Помилка при читанні повідомлення:", err)
			}
			break
		}

		var wsMessage WSMessage
		if err := json.Unmarshal(message, &wsMessage); err != nil {
			logrus.Error("Помилка при розкодуванні повідомлення:", err)
			continue
		}

		// Обробка різних типів повідомлень
		switch wsMessage.Type {
		case "send_message":
			chatMessage := &models.ChatMessagesModel{
				Chat_id:      wsMessage.ChatID,
				Message_id:   wsMessage.MessageID,
				Sender_id:    wsMessage.SenderID,
				Message_text: wsMessage.Text,
				Sent_at:      time.Now(),
			}

			if err := c.chatRepo.SendMessage(chatMessage); err != nil {
				logrus.Errorf("Помилка при відправці повідомлення: %v", err)
				continue
			}

			// Відправляємо повідомлення всім клієнтам в чаті
			messageBytes, _ := json.Marshal(wsMessage)
			c.hub.broadcast <- messageBytes

		case "get_history":
			// Отримання історії чату
			history, err := c.chatRepo.GetMessageHistory(wsMessage.ChatID)
			if err != nil {
				logrus.Errorf("Помилка при отриманні історії: %v", err)
				continue
			}

			historyBytes, err := json.Marshal(history)
			if err != nil {
				logrus.Errorf("Помилка при серіалізації історії: %v", err)
				continue
			}
			c.send <- historyBytes

		case "delete_message":
			if err := c.chatRepo.DeleteMessage(wsMessage.ChatID, wsMessage.MessageID); err != nil {
				logrus.Errorf("Помилка при видаленні повідомлення: %v", err)
				continue
			}

			messageBytes, _ := json.Marshal(wsMessage)
			c.hub.broadcast <- messageBytes
		}
	}
}

func (c *Client) writePump() {
	defer func() {
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)
			if err := w.Close(); err != nil {
				return
			}
		}
	}
}
