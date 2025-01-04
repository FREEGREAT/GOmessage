package wsserver

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

const (
	templatedir = "./web/templates/html"
)

type Client struct {
	UserID string
	Conn   *websocket.Conn
}
type WSServer interface {
	Start() error
	Stop() error
}

type wsServ struct {
	mux       *http.ServeMux
	srv       *http.Server
	wsUpg     *websocket.Upgrader
	wsClients map[*websocket.Conn]struct{}
	mutex     *sync.RWMutex
	broadcast chan *wsMessage
}

func NewWsServer(addr string) WSServer {
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
		},
		wsClients: map[*websocket.Conn]struct{}{},
		mutex:     &sync.RWMutex{},
		broadcast: make(chan *wsMessage),
	}
}

func (ws *wsServ) Start() error {
	ws.mux.Handle("/", http.FileServer(http.Dir(templatedir)))
	ws.mux.HandleFunc("/ws", ws.wsHandler)
	ws.mux.HandleFunc("/test", ws.testHandler)
	go ws.writeToClientsBroadcast()
	return ws.srv.ListenAndServe()
}

func (ws *wsServ) wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := ws.wsUpg.Upgrade(w, r, nil)
	if err != nil {
		logrus.Errorf("Error upgrading connection: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	logrus.Infof(conn.RemoteAddr().String())
	ws.mutex.Lock()
	ws.wsClients[conn] = struct{}{}
	ws.mutex.Unlock()
	ws.readFromClient(conn)
}

func (ws *wsServ) readFromClient(conn *websocket.Conn) {
	for {
		msg := new(wsMessage)
		err := conn.ReadJSON(msg)
		if err != nil {
			wsErr, ok := err.(*websocket.CloseError)
			if !ok || wsErr.Code != websocket.CloseGoingAway {
				logrus.Errorf("Error with reading from ws client: %v", err)
			}
			break
		}
		msg.IPAddress = conn.RemoteAddr().String()
		msg.Time = time.Now().GoString()
		ws.broadcast <- msg
	}
	ws.mutex.Lock()
	delete(ws.wsClients, conn)
	ws.mutex.Unlock()

}

func (ws *wsServ) writeToClientsBroadcast() {
	for msg := range ws.broadcast {
		ws.mutex.RLock()
		for client := range ws.wsClients {
			go func() {
				if err := client.WriteJSON(msg); err != nil {
					logrus.Errorf("Error with writing message: %v", err)
				}
			}()
		}
		ws.mutex.RUnlock()
	}
}

func (ws *wsServ) testHandler(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("Hello World"))
}
func (ws *wsServ) Stop() error {
	close(ws.broadcast)
	ws.mutex.Lock()
	for conn := range ws.wsClients {
		conn.Close()
		delete(ws.wsClients, conn)
	}
	ws.mutex.Unlock()
	return ws.srv.Shutdown(context.Background())
}
