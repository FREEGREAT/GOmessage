package wsserver

type wsMessage struct {
	UserID    string `json:"userId"`
	ChatID    string `json:"chatId"`
	Content   string `json:"content"`
	Time      string `json:"time"`
	IPAddress string `json:"ipAddress"`
}
