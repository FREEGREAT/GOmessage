package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	proto_media_service "github.com/FREEGREAT/protos/gen/go/media"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	proto_geoip_service "github.com/FREEGREAT/protos/gen/go/geoip"
	proto_user_service "github.com/FREEGREAT/protos/gen/go/user"
	"github.com/julienschmidt/httprouter"
	"gomessage.com/gateway/internal/handler/middleware"
	"gomessage.com/gateway/internal/models"
	"gomessage.com/gateway/internal/service"
	"gomessage.com/gateway/pkg/utils"
)

type IPLocation struct {
	Status      string  `json:"status"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	Region      string  `json:"region"`
	RegionName  string  `json:"regionName"`
	City        string  `json:"city"`
	Zip         string  `json:"zip"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Timezone    string  `json:"timezone"`
	ISP         string  `json:"isp"`
	Query       string  `json:"query"`
}

const (
	nonePhoto       = "NULL"
	SuccessResponse = "Success"
)
const (
	usersURL  = "/users"
	userURL   = "/user/:uuid"
	loginURL  = "/login"
	logoutURL = "/logout"
	signupURL = "/signup"
	deleteUrl = "/delete/:uuid"
	friendURL = "/friend"
)

type handler struct {
	MediaGrpcClient proto_media_service.MediaServiceClient
	UserGrpcClient  proto_user_service.UserServiceClient
	GeoIpGrpcClient proto_geoip_service.GeoIpServiceClient
	jwtService      *service.JWTService
	jwtMidlleware   *middleware.JWTMiddleware
}

func NewGatewayHandler(grpcClient proto_user_service.UserServiceClient, mc proto_media_service.MediaServiceClient, gc proto_geoip_service.GeoIpServiceClient, jwtService *service.JWTService) *handler {
	return &handler{
		UserGrpcClient:  grpcClient,
		MediaGrpcClient: mc,
		GeoIpGrpcClient: gc,
		jwtService:      jwtService,
		jwtMidlleware:   middleware.NewJWTMiddleware(jwtService)}
}

func (h *handler) Register(router *httprouter.Router) {
	router.GET(usersURL, h.jwtMidlleware.Middleware(h.ListOfUsers))
	router.GET(friendURL, h.jwtMidlleware.Middleware(h.ListOfFriends))

	router.POST(signupURL, h.CreateUser)
	router.POST(friendURL, h.jwtMidlleware.Middleware(h.AddFriend))
	router.POST(loginURL, h.LoginUser)

	router.PUT(userURL, h.jwtMidlleware.Middleware(h.UpdateUser))

	router.DELETE(deleteUrl, h.jwtMidlleware.Middleware(h.DeleteUser))
	router.DELETE(friendURL, h.jwtMidlleware.Middleware(h.DeleteFriend))
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

func (h handler) CreateUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	defer r.Body.Close()

	deviceInfo := utils.GetDeviceInfo(r.UserAgent())

	locationInfo := utils.GetIP(r)

	grpcRequest := &proto_geoip_service.GetLocationRequest{
		Ip: locationInfo,
	}
	response, err := h.GeoIpGrpcClient.GetLocationByIP(r.Context(), grpcRequest)
	if err != nil {
		logrus.Errorf("Failed to parse user IP: %v", err)
	}

	user, err := parseUserData(r)
	if err != nil {
		logrus.Errorf("Failed to parse user data: %v", err)
		http.Error(w, "Invalid user data: "+err.Error(), http.StatusBadRequest)
		return
	}

	userId, err := registerUser(r.Context(), h, user, response.Location, deviceInfo["deviceType"])
	if err != nil {
		logrus.Errorf("Failed to register user: %v", err)
		http.Error(w, "Failed to register user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = handleUserPhoto(r, h, userId)
	if err != nil {
		logrus.Errorf("Failed to handle user photo: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func parseUserData(r *http.Request) (*models.UserModel, error) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		return nil, fmt.Errorf("failed to parse form: %v", err)
	}

	userPart := r.FormValue("user")
	if userPart == "" {
		return nil, fmt.Errorf("user data is empty")
	}

	var user models.UserModel
	if err := json.Unmarshal([]byte(userPart), &user); err != nil {
		return nil, fmt.Errorf("invalid user data: %v", err)
	}

	return &user, nil
}

func registerUser(ctx context.Context, h handler, user *models.UserModel, location string, device string) (string, error) {
	grpcRequest := &proto_user_service.RegisterUserRequest{
		Nickname: user.Nickname,
		Password: user.PasswordHash,
		Email:    user.Email,
		Age:      int32(*user.Age),
		ImageUrl: nonePhoto,
		Location: location,
		Device:   device,
	}

	response, err := h.UserGrpcClient.RegisterUser(ctx, grpcRequest)
	if err != nil {
		return "", err
	}

	if response.Status != SuccessResponse {
		return "", fmt.Errorf("registration failed with status: %s", response.Status)
	}

	return response.UserId, nil
}

func handleUserPhoto(r *http.Request, h handler, userId string) (string, error) {
	file, _, err := r.FormFile("photo")
	if err != nil {
		if err == http.ErrMissingFile {
			return nonePhoto, nil
		}
		return "", err
	}
	defer file.Close()

	photoBytes, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	savePhotoReq := &proto_media_service.SavePhotoRequest{Photo: photoBytes}
	res, err := h.MediaGrpcClient.SavePhoto(context.Background(), savePhotoReq)
	if err != nil {
		return "", err
	}

	updateReq := &proto_user_service.UpdateUserRequest{
		Id:       &userId,
		ImageUrl: &res.PhotoLink,
	}

	_, err = h.UserGrpcClient.UpdateUser(context.Background(), updateReq)
	if err != nil {
		return "", err
	}

	return res.PhotoLink, nil
}

func (h *handler) LoginUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	defer r.Body.Close()

	var loginData models.UserModel

	if err := json.NewDecoder(r.Body).Decode(&loginData); err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	if loginData.Email == "" || loginData.PasswordHash == "" {
		http.Error(w, "Missing required fields: email or password", http.StatusBadRequest)
		return
	}

	grpcRequest := &proto_user_service.LoginUserRequest{
		Email:    loginData.Email,
		Password: loginData.PasswordHash,
	}

	grpcResp, err := h.UserGrpcClient.LoginUser(context.Background(), grpcRequest)
	if err != nil {
		http.Error(w, "Failed to login: "+err.Error(), http.StatusInternalServerError)
		return
	}
	age, err := strconv.Atoi(grpcResp.Age)
	if err != nil {
		logrus.Error("Failed to get user age")
		age = 0
	}
	token, err := h.jwtService.GenerateToken(grpcResp.Id, grpcResp.Username, grpcResp.ImageUrl, age)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User logged in successfully      " + token))
}

func (h *handler) UpdateUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	defer r.Body.Close()

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

func getLocation(ip string) (*IPLocation, error) {
	url := fmt.Sprintf("http://ip-api.com/json/%s", ip)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var location IPLocation
	if err := json.Unmarshal(body, &location); err != nil {
		return nil, err
	}

	return &location, nil
}
