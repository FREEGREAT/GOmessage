package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"gomessage.com/users/internal/models"
	"gomessage.com/users/internal/service"
)

const (
	usersURL  = "/users"
	userURL   = "/user/:uuid"
	friendURL = "/friend"
)

type handler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *handler {
	return &handler{userService: userService}
}

func (h *handler) Register(router *httprouter.Router) {
	router.GET(usersURL, h.GetList)
	router.GET(userURL, h.GetUserByUUID)
	router.POST(usersURL, h.CreateUser)
	router.PUT(userURL, h.UpdateUser)
	router.DELETE(userURL, h.DeleteUser)
}

func (h *handler) GetList(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	users, err := h.userService.ListOfUsers(r.Context()) // Викликаємо сервіс для отримання списку користувачів
	if err != nil {
		http.Error(w, "Failed to get users: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, "Failed to encode users: "+err.Error(), http.StatusInternalServerError)
	}
}

func (h *handler) GetUserByUUID(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	uuid := param.ByName("uuid")

	// Отримуємо користувача з сервісу
	user, err := h.userService.GetUser(context.TODO(), uuid)
	if err != nil {
		http.Error(w, "Failed to get user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Якщо користувач знайдений, повертаємо його дані в форматі JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 200 OK
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "Failed to encode user data: "+err.Error(), http.StatusInternalServerError)
	}
}

func (h *handler) CreateUser(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	var user models.UserModel
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Inwalid input", http.StatusBadRequest)
	}

	defer r.Body.Close()
	if err := h.userService.CreateUser(context.TODO(), &user); err != nil {
		http.Error(w, "Failed to create user: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User created successfuly"))
}

func (h *handler) UpdateUser(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	var user models.UserModel
	uuid := param.ByName("uuid")
	user.ID = uuid
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
	}
	defer r.Body.Close()

	if err := h.userService.UpdateUser(context.TODO(), user); err != nil {
		http.Error(w, "Failed to update user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(201)
	w.Write([]byte("update user"))
}
func (h *handler) DeleteUser(w http.ResponseWriter, r *http.Request, param httprouter.Params) {

	uuid := param.ByName("uuid")

	info, err := h.userService.DeleteUser(context.TODO(), uuid)
	if err != nil {
		http.Error(w, "Failed to update user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(201)
	w.Write([]byte("delete user" + info))
}
