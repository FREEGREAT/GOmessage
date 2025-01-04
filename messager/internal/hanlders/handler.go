package hanlders

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func SetupRoutes() {
	http.HandleFunc("/ws", wsEndpoint)
}

func reader(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			logrus.Error("Error while reading message: %w", err)
		}
		logrus.Infof("Message: %s", p)
		if err := conn.WriteMessage(messageType, p); err != nil {
			logrus.Error("Error while writing message: %w", err)
		}
	}
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logrus.Error("Error while upgrading connection: %w", err)
	}
	reader(ws)
}
