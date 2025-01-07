package wsserver

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	middleware "gommessage.com/messager/internal/hanlders/midleware"
	"gommessage.com/messager/internal/models"
	"gommessage.com/messager/internal/storage"
)

const (
	templatedir         = "./web/templates/html"
	MaxChatParticipants = 2
)

type Client struct {
	UserID string
	Conn   *websocket.Conn
}

type Chat struct {
	ID       string
	Clients  map[string]*Client // map[userID]*Client
	Messages chan *wsMessage
	mutex    sync.RWMutex
}

type WSServer interface {
	Start() error
	Stop() error
}

type wsServ struct {
	mux   *http.ServeMux
	srv   *http.Server
	wsUpg *websocket.Upgrader
	chats map[string]*Chat // map[chatID]*Chat
	mutex sync.RWMutex
	repo  storage.ChatRepository
}

func NewWsServer(addr string, repo storage.ChatRepository) WSServer {
	m := http.NewServeMux()
	return &wsServ{
		mux: m,
		srv: &http.Server{
			Addr:    addr,
			Handler: m,
		},
		wsUpg: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true // You might want to implement proper origin checking
			},
		},
		chats: make(map[string]*Chat),
		repo:  repo,
	}
}

func (ws *wsServ) Start() error {
	ws.mux.Handle("/", http.FileServer(http.Dir(templatedir)))
	ws.mux.HandleFunc("/ws/", middleware.WSAuthMiddleware(ws.wsHandler))
	ws.mux.HandleFunc("/chat/list", ws.getUserChatsHandler)
	return ws.srv.ListenAndServe()
}

func (ws *wsServ) getUserChatsHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userId")
	if userID == "" {
		http.Error(w, "userId is required", http.StatusBadRequest)
		return
	}
	logrus.Infof("Fetching chats for userId: %s", userID)
	chats, err := ws.repo.GetUserChats(userID)
	if err != nil {
		logrus.Errorf("Error getting user chats: %v", err)
		http.Error(w, "Failed to get user chats", http.StatusInternalServerError)
		return
	}
	logrus.Infof("Retrieved chats: %+v", chats)
	// Перетворюємо чати в формат для відповіді
	type chatResponse struct {
		ID      string `json:"id"`
		UserID1 string `json:"userId1"`
		UserID2 string `json:"userId2"`
	}

	response := make([]chatResponse, len(chats))
	for i, chat := range chats {
		response[i] = chatResponse{
			ID:      chat.Chat_id,
			UserID1: chat.User_id1,
			UserID2: chat.User_id2,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logrus.Errorf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
func (ws *wsServ) wsHandler(w http.ResponseWriter, r *http.Request) {
	chatID := r.URL.Query().Get("chatId")
	userID := r.URL.Query().Get("userId")

	if chatID == "" || userID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	chatInfo, err := ws.repo.GetChatInfo(chatID)
	if err != nil {
		logrus.Errorf("Error getting chat info: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if userID != chatInfo.User_id1 && userID != chatInfo.User_id2 {
		logrus.Warnf("Unauthorized access attempt. User %s not in chat %s", userID, chatID)
		w.WriteHeader(http.StatusForbidden)
		return
	}
	//Load chat history
	messages, err := ws.repo.GetMessageHistory(chatID)
	if err != nil {
		logrus.Errorf("Error loading chat history: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ws.mutex.Lock()
	chat, exists := ws.chats[chatID]
	if !exists {
		chat = &Chat{
			ID:       chatID,
			Clients:  make(map[string]*Client),
			Messages: make(chan *wsMessage, 100),
		}
		ws.chats[chatID] = chat
		go ws.handleChatMessages(chat)
	}
	ws.mutex.Unlock()

	chat.mutex.Lock()
	if len(chat.Clients) >= MaxChatParticipants {
		chat.mutex.Unlock()
		w.WriteHeader(http.StatusForbidden)
		return
	}

	conn, err := ws.wsUpg.Upgrade(w, r, nil)
	if err != nil {
		chat.mutex.Unlock()
		logrus.Errorf("Error upgrading connection: %v", err)
		return
	}

	client := &Client{
		UserID: userID,
		Conn:   conn,
	}
	chat.Clients[userID] = client
	chat.mutex.Unlock()

	for _, msg := range messages {
		wsMsg := &wsMessage{
			UserID:    msg.Sender_id,
			ChatID:    msg.Chat_id,
			Content:   msg.Message_text,
			Time:      msg.Sent_at.Format(time.RFC3339),
			IPAddress: conn.RemoteAddr().String(),
		}
		if err := client.Conn.WriteJSON(wsMsg); err != nil {
			logrus.Errorf("Error sending history message: %v", err)
			continue
		}
	}

	ws.readFromClient(client, chat)
}

func (ws *wsServ) readFromClient(client *Client, chat *Chat) {
	defer func() {
		chat.mutex.Lock()
		delete(chat.Clients, client.UserID)
		client.Conn.Close()

		if len(chat.Clients) == 0 {
			ws.mutex.Lock()
			delete(ws.chats, chat.ID)
			close(chat.Messages)
			ws.mutex.Unlock()
		}
		chat.mutex.Unlock()
	}()

	for {
		msg := struct {
			MessageText string `json:"message_text"`
			Action      string `json:"action"`     // Додано поле для дії
			MessageID   string `json:"message_id"` // Додано поле для ID повідомлення
		}{}

		err := client.Conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logrus.Errorf("Error reading from client: %v", err)
			}
			break
		}

		switch msg.Action {
		case "send":
			// Обробка дії "send"
			dbMsg := &models.ChatMessagesModel{
				Chat_id:      chat.ID,
				Message_id:   uuid.New().String(),
				Sender_id:    client.UserID,
				Message_text: msg.MessageText,
				Sent_at:      time.Now(),
			}

			if err := ws.repo.SendMessage(dbMsg); err != nil {
				logrus.Errorf("Error saving message: %v", err)
				continue
			}

			broadcastMsg := struct {
				ChatID      string    `json:"chat_id"`
				MessageID   string    `json:"message_id"`
				SenderID    string    `json:"sender_id"`
				MessageText string    `json:"message_text"`
				SentAt      time.Time `json:"sent_at"`
			}{
				ChatID:      dbMsg.Chat_id,
				MessageID:   dbMsg.Message_id,
				SenderID:    dbMsg.Sender_id,
				MessageText: dbMsg.Message_text,
				SentAt:      dbMsg.Sent_at,
			}

			chat.mutex.RLock()
			for _, c := range chat.Clients {
				go func(c *Client) {
					if err := c.Conn.WriteJSON(broadcastMsg); err != nil {
						logrus.Errorf("Error writing message to client: %v", err)
					}
				}(c)
			}
			chat.mutex.RUnlock()

		case "delete":
			// Обробка дії "delete"
			if msg.MessageID == "" {
				logrus.Warn("Message ID is required for delete action")
				continue
			}

			// Виклик методу для видалення повідомлення
			if err := ws.repo.DeleteMessage(chat.ID, msg.MessageID); err != nil {
				logrus.Errorf("Error deleting message: %v", err)
				continue
			}

			// Підготовка повідомлення для всіх клієнтів про видалення
			deleteMsg := struct {
				ChatID    string `json:"chat_id"`
				MessageID string `json:"message_id"`
			}{
				ChatID:    chat.ID,
				MessageID: msg.MessageID,
			}

			chat.mutex.RLock()
			for _, c := range chat.Clients {
				go func(c *Client) {
					if err := c.Conn.WriteJSON(deleteMsg); err != nil {
						logrus.Errorf("Error writing delete message to client: %v", err)
					}
				}(c)
			}
			chat.mutex.RUnlock()
		default:
			logrus.Warnf("Unknown action: %s", msg.Action)
		}
	}
}

func (ws *wsServ) handleChatMessages(chat *Chat) {
	for msg := range chat.Messages {
		chat.mutex.RLock()
		for _, client := range chat.Clients {
			go func(c *Client) {
				if err := c.Conn.WriteJSON(msg); err != nil {
					logrus.Errorf("Error writing message: %v", err)
				}
			}(client)
		}
		chat.mutex.RUnlock()
	}
}

func (ws *wsServ) Stop() error {
	ws.mutex.Lock()
	for _, chat := range ws.chats {
		chat.mutex.Lock()
		for _, client := range chat.Clients {
			client.Conn.Close()
		}
		close(chat.Messages)
		chat.mutex.Unlock()
	}
	ws.mutex.Unlock()
	return ws.srv.Shutdown(context.Background())
}
