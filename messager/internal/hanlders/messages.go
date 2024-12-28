package hanlders

// import (
// 	"encoding/json"
// 	"net/http"

// 	"github.com/julienschmidt/httprouter"
// 	"github.com/sirupsen/logrus"
// 	"gommessage.com/messager/internal/models"
// 	database "gommessage.com/messager/pkg/Database"
// )

// var (
// 	message models.MessageModel
// )

// func getHandler(w http.ResponseWriter, r *http.Request) {
// 	json.NewEncoder(w).Encode("uuu")
// }

// const nonePhoto = "NULL"
// const SuccessResponse = "Success"
// const (
// 	usersURL = "/send"
// )

// type handler struct {
// }

// func NewGatewayHandler() *handler {
// 	return &handler{}
// }

// func (h *handler) Register(router *httprouter.Router) {
// 	logrus.Info("Ahuenno owou wou register")
// 	router.POST(usersURL, h.postHandler)
// }

// func (h *handler) postHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

// 	json.NewDecoder(r.Body).Decode(&message)
// 	sendMessage(&message)

// 	json.NewEncoder(w).Encode("You send message wow")
// }

// func sendMessage(message *models.MessageModel) {
// 	query := `INSERT INTO messages(message_id,user1_id, user2_id, message)
// 	VALUES(now(),?,?,?)`
// 	logrus.Info("Ahuenno owou wou querry")
// 	database.Exec(query, message.User_id1, message.User_id2, message.Message)
// }
