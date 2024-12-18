package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	proto_media_service "github.com/FREEGREAT/protos/gen/go/media"
	"github.com/sirupsen/logrus"

	proto_user_service "github.com/FREEGREAT/protos/gen/go/user"
	"github.com/julienschmidt/httprouter"
	"gomessage.com/users/internal/grpcclient"
	"gomessage.com/users/internal/models"
)

const (
	usersURL  = "/users"
	userURL   = "/user/:uuid"
	friendURL = "/friend"
)

type handler struct {
	grpcClient  grpcclient.GRPCClient
	MediaClient proto_media_service.MediaServiceClient
}

func NewUserHandler(grpcClient grpcclient.GRPCClient, mc proto_media_service.MediaServiceClient) *handler {
	return &handler{grpcClient: grpcClient, MediaClient: mc}
}

func (h *handler) Register(router *httprouter.Router) {
	router.GET(usersURL, h.GetList)
	router.POST(usersURL, h.CreateUser)
	router.PUT(userURL, h.UpdateUser)
}

func (h *handler) GetList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	response, err := h.grpcClient.GetUsers(r.Context(), &proto_user_service.GetUsersRequest{})
	if err != nil {
		http.Error(w, "Failed to get users: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response.Users); err != nil {
		http.Error(w, "Failed to encode users: "+err.Error(), http.StatusInternalServerError)
	}
}

func (h *handler) CreateUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	defer r.Body.Close() // Закриваємо тіло запиту

	var user models.UserModel

	r.ParseMultipartForm(10 << 20) // 10 MB
	userPart := r.FormValue("user")

	if err := json.Unmarshal([]byte(userPart), &user); err != nil {
		http.Error(w, "Invalid user data: "+err.Error(), http.StatusBadRequest)
		return
	}

	logrus.Error("Marshal json")
	file, _, err := r.FormFile("photo")
	if err != nil {
		http.Error(w, "Failed to read photo: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()
	logrus.Error("Read file")
	photoBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read photo: "+err.Error(), http.StatusInternalServerError)
		return
	}
	savePhotoReq := proto_media_service.SavePhotoRequest{Photo: photoBytes}
	logrus.Error("Send req")
	res, err := h.MediaClient.SavePhoto(context.Background(), &savePhotoReq)
	if err != nil {
		logrus.Errorf("Error while sending photo to MediaClient: %v", err)
		return
	}
	// Обробка успішної відповіді
	logrus.Infof("Photo saved successfully: %v", res)

	logrus.Error("Get res")
	grpcRequest := &proto_user_service.RegisterUserRequest{
		Nickname: user.Nickname,
		Password: user.PasswordHash,
		Email:    user.Email,
		Age:      int32(*user.Age),
		ImageUrl: res.PhotoLink,
	}

	response, err := h.grpcClient.RegisterUser(r.Context(), grpcRequest)
	if err != nil {
		http.Error(w, "Failed to register user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(response.UserId))
}

func (h *handler) UpdateUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	defer r.Body.Close() // Закриваємо тіло запиту

	uuid := params.ByName("uuid")
	var user models.UserModel

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	grpcRequest := &proto_user_service.UpdateUserRequest{
		Id:       &uuid,
		Username: &user.Nickname,
		Email:    &user.Email,
	}

	_, err := h.grpcClient.UpdateUser(r.Context(), grpcRequest)
	if err != nil {
		http.Error(w, "Failed to update user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User  updated successfully"))
}

// Uncomment and implement the DeleteUser  method if needed
// func (h *handler) DeleteUser (w http.ResponseWriter, r *http.Request, params httprouter.Params) {
// 	uuid := params.ByName("uuid")
// 	grpcRequest := &proto_user_service.DeleteUser Request{ Id: uuid }
// 	_, err := h.grpcClient.DeleteUser (r.Context(), grpcRequest)
// 	if err != nil {
// 		http.Error(w, "Failed to delete user: "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	w.WriteHeader(http.StatusOK)
// 	w.Write([]byte("User  deleted successfully"))
// }
