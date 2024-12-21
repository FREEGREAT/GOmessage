package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	proto_media_service "github.com/FREEGREAT/protos/gen/go/media"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	proto_user_service "github.com/FREEGREAT/protos/gen/go/user"
	"github.com/julienschmidt/httprouter"
	"gomessage.com/gateway/internal/models"
)

const nonePhoto = "NULL"
const SuccessResponse = "Success"
const (
	usersURL  = "/users"
	userURL   = "/user/:uuid"
	loginURL  = "/login"
	signupURL = "/signup"
	deleteUrl = "/delete/:uuid"
	friendURL = "/friend"
)

type handler struct {
	MediaGrpcClient proto_media_service.MediaServiceClient
	UserGrpcClient  proto_user_service.UserServiceClient
}

func NewGatewayHandler(grpcClient proto_user_service.UserServiceClient, mc proto_media_service.MediaServiceClient) *handler {
	return &handler{UserGrpcClient: grpcClient, MediaGrpcClient: mc}
}

func (h *handler) Register(router *httprouter.Router) {
	router.GET(usersURL, h.ListOfUsers)
	router.GET(friendURL, h.ListOfFriends)

	router.POST(signupURL, h.CreateUser)
	router.POST(friendURL, h.AddFriend)
	router.POST(loginURL, h.LoginUser)

	router.PUT(userURL, h.UpdateUser)

	router.DELETE(deleteUrl, h.DeleteUser)
	router.DELETE(friendURL, h.DeleteFriend)
}

func (h *handler) ListOfUsers(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	response, err := h.UserGrpcClient.GetUsers(r.Context(), &proto_user_service.GetUsersRequest{})
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

func (h *handler) ListOfFriends(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	defer r.Body.Close()

	var friendData models.FriendListModel

	if err := json.NewDecoder(r.Body).Decode(&friendData); err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}
	if friendData.UserID == "" {
		http.Error(w, "Missing required fields: userID", http.StatusBadRequest)
		return
	}

	response, err := h.UserGrpcClient.ListOfFriends(r.Context(), &proto_user_service.ListOfFriendsRequest{
		UserId: friendData.UserID,
	})
	if err != nil {
		http.Error(w, "Failed to get users: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response.Friend); err != nil {
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

	grpcRequest := &proto_user_service.RegisterUserRequest{
		Nickname: user.Nickname,
		Password: user.PasswordHash,
		Email:    user.Email,
		Age:      int32(*user.Age),
		ImageUrl: nonePhoto,
	}
	logrus.Info("Register user without photo")
	response, err := h.UserGrpcClient.RegisterUser(r.Context(), grpcRequest)
	if err != nil {
		http.Error(w, "Failed to register user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	logrus.Info("Check response status")
	if response.Status == SuccessResponse {
		file, _, err := r.FormFile("photo")
		if err != nil {
			if err == http.ErrMissingFile {
				logrus.Info("No photo uploaded, proceeding without it.")
			} else {
				http.Error(w, "Failed to read photo: "+err.Error(), http.StatusBadRequest)
				return
			}
		} else {
			defer file.Close()
			photoBytes, err := io.ReadAll(file)
			if err != nil {
				http.Error(w, "Failed to read photo: "+err.Error(), http.StatusInternalServerError)
				return
			}
			savePhotoReq := proto_media_service.SavePhotoRequest{Photo: photoBytes}
			res, err := h.MediaGrpcClient.SavePhoto(context.Background(), &savePhotoReq)
			if err != nil {
				logrus.Errorf("Error while sending photo to MediaClient: %v", err)
				return
			}
			logrus.Info(res.PhotoLink)

			grpcUpdateRequest := &proto_user_service.UpdateUserRequest{
				Id:       &response.UserId,
				ImageUrl: &res.PhotoLink,
			}
			_, err = h.UserGrpcClient.UpdateUser(context.Background(), grpcUpdateRequest)
			if err != nil {
				logrus.Errorf("Error while saving photo link: %v", err)
				return
			}
		}
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(response.UserId))
}

func (h *handler) LoginUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	defer r.Body.Close()

	var loginData models.UserModel

	if err := json.NewDecoder(r.Body).Decode(&loginData); err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}
	if loginData.Email == "" || loginData.PasswordHash == "" {
		http.Error(w, "Missing required fields: nickname or password", http.StatusBadRequest)
		return
	}
	grpcRequest := &proto_user_service.LoginUserRequest{
		Email:    loginData.Email,
		Password: loginData.PasswordHash,
	}

	_, err := h.UserGrpcClient.LoginUser(context.Background(), grpcRequest)
	if err != nil {
		http.Error(w, "Failed to login: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User logged in successfully"))

}

func (h *handler) UpdateUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	defer r.Body.Close() // Закриваємо тіло запиту

	_, err := uuid.Parse(params.ByName("uuid"))
	if err != nil {
		http.Error(w, "Invalid UUID: "+err.Error(), http.StatusBadRequest)
		return
	}
	uuidSTR := params.ByName("uuid")
	var user models.UserModel

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	if user.Nickname == "" || user.Email == "" {
		http.Error(w, "Missing required fields: nickname or email", http.StatusBadRequest)
		return
	}

	grpcRequest := &proto_user_service.UpdateUserRequest{
		Id:           &uuidSTR,
		Username:     &user.Nickname,
		Email:        &user.Email,
		PasswordHash: &user.PasswordHash,
	}

	_, err = h.UserGrpcClient.UpdateUser(r.Context(), grpcRequest)
	if err != nil {
		http.Error(w, "Failed to update user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User  updated successfully"))
}

func (h *handler) DeleteUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	uuid := params.ByName("uuid")
	grpcRequest := proto_user_service.DeleteUserRequest{Id: uuid}
	_, err := h.UserGrpcClient.DeleteUser(r.Context(), &grpcRequest)
	if err != nil {
		http.Error(w, "Failed to delete user: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User  deleted successfully"))
}

func (h *handler) AddFriend(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var friends models.FriendListModel

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	logrus.Info("Decoder")
	r.Body = io.NopCloser(bytes.NewBuffer(body))
	if err := json.NewDecoder(r.Body).Decode(&friends); err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}
	logrus.Info("Nil checker")
	if friends.UserID == "" || friends.FriendID == "" {
		http.Error(w, "Missing required fields: uid, fid"+err.Error(), http.StatusBadRequest)
	}
	logrus.Info("grpcreq creating")
	grpcRequest := &proto_user_service.AddFriendsRequest{
		UserId_1: friends.UserID,
		UserId_2: friends.FriendID,
	}
	logrus.Info("Grpc querry")
	resp, err := h.UserGrpcClient.AddFriends(r.Context(), grpcRequest)
	if err != nil {
		http.Error(w, "Failed to add friend: "+err.Error(), http.StatusBadRequest)
		return
	}
	if !resp.Success {
		http.Error(w, "Error: "+resp.Message, http.StatusAccepted)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("You are friends"))

}

func (h *handler) DeleteFriend(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var friends models.FriendListModel

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	r.Body = io.NopCloser(bytes.NewBuffer(body))
	if err := json.NewDecoder(r.Body).Decode(&friends); err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}
	if friends.UserID == "" || friends.FriendID == "" {
		http.Error(w, "Missing required fields: uid, fid"+err.Error(), http.StatusBadRequest)
	}
	grpcRequest := &proto_user_service.AddFriendsRequest{
		UserId_1: friends.UserID,
		UserId_2: friends.FriendID,
	}
	resp, err := h.UserGrpcClient.DeleteFriend(r.Context(), (*proto_user_service.DeleteFriendsRequest)(grpcRequest))
	if err != nil {
		http.Error(w, "Failed to delete friend: "+err.Error(), http.StatusBadRequest)
		return
	}
	if !resp.Success {
		http.Error(w, "Error: "+resp.Message, http.StatusAccepted)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("You you are not friends anymore"))

}
