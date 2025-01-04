package wsserver

type wsMessage struct {
	Content     string `json:"content"`
	FromUser    string `json:"from_user"`
	ToUser      string `json:"to_user"`
	IPAddress   string `json:"ip_address"`
	Time        string `json:"time"`
	MessageType string `json:"message_type"`
}
