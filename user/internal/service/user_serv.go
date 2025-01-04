package service

import (
	"context"

	proto_user_service "github.com/FREEGREAT/protos/gen/go/user"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gomessage.com/users/internal/models"
	"gomessage.com/users/internal/service/kafka"
	"gomessage.com/users/internal/storage"
	"gomessage.com/users/pkg/utils"
)

type userService struct {
	userRepository    storage.UserRepository
	friendsRepository storage.FriendsRepository
	kProducer         kafka.Producer
	proto_user_service.UnimplementedUserServiceServer
}

func CreateNewUserService(usrRepo storage.UserRepository, friendRepo storage.FriendsRepository, producer kafka.Producer) *userService {

	return &userService{
		userRepository:    usrRepo,
		friendsRepository: friendRepo,
		kProducer:         producer,
	}
}

func (u *userService) GetUsers(context.Context, *proto_user_service.GetUsersRequest) (*proto_user_service.GetUsersResponse, error) {

	var users = []models.UserModel{}
	users, err := u.userRepository.FindAll(context.Background())
	if err != nil {
		logrus.Errorf("Error while finding all user. Err %s", err)
		return nil, err
	}
	var protoUsers []*proto_user_service.User
	for _, user := range users {
		var image string
		if user.ImageUrl == nil {
			image = ""
		}
		protoUser := proto_user_service.User{
			Id:           user.ID,
			Username:     user.Nickname,
			PasswordHash: user.PasswordHash,
			Email:        user.Email,
			Age:          uint32(*user.Age),
			ImageUrl:     image,
		}
		protoUsers = append(protoUsers, &protoUser)

	}

	res := &proto_user_service.GetUsersResponse{Status: "Success", Users: protoUsers}
	return res, nil
}

func (u *userService) UpdateUser(ctx context.Context, in *proto_user_service.UpdateUserRequest) (*proto_user_service.UpdateUserResponse, error) {

	existingUser, err := u.userRepository.FindById(ctx, *in.Id)
	if err != nil {
		logrus.Errorf("Error while fetching user: %s", err)
		return &proto_user_service.UpdateUserResponse{Status: "Error"}, err
	}

	if in.Username != nil {
		existingUser.Nickname = *in.Username
	}
	if in.PasswordHash != nil {
		existingUser.PasswordHash = *in.PasswordHash
	}
	if in.ImageUrl != nil {
		existingUser.ImageUrl = in.ImageUrl
	}
	if in.Email != nil {
		existingUser.Email = *in.Email
	}
	if in.Age != nil {
		age := int(*in.Age)
		existingUser.Age = &age
	}

	if err := u.userRepository.Update(ctx, &existingUser); err != nil {
		logrus.Errorf("Error while updating user, error info: %s", err)
		return &proto_user_service.UpdateUserResponse{Status: "Error"}, err
	}
	res := proto_user_service.UpdateUserResponse{Status: "Success"}
	return &res, nil
}
func (u *userService) RegisterUser(ctx context.Context, in *proto_user_service.RegisterUserRequest) (*proto_user_service.RegisterUserResponse, error) {
	pass, err := utils.GeneratePasswordHash(in.Password)
	if err != nil {
		logrus.Errorf("Error hashing password: %w", err)
		return nil, err
	}
	age := int(in.Age)
	user := models.UserModel{
		Nickname:     in.Nickname,
		PasswordHash: pass,
		Email:        in.Email,
		Age:          &age,
		ImageUrl:     &in.ImageUrl,
	}
	id, err := u.userRepository.Create(ctx, &user)
	if err != nil {
		logrus.Errorf("Error creating user: %w", err)
		return nil, err
	}
	req := proto_user_service.RegisterUserResponse{
		UserId: id,
		Status: "Success",
	}

	return &req, nil
}
func (u *userService) LoginUser(ctx context.Context, in *proto_user_service.LoginUserRequest) (*proto_user_service.LoginUserResponse, error) {

	password_hash, err := utils.GeneratePasswordHash(in.Password)
	if err != nil {
		logrus.Errorf("Error while hashing password")
		return &proto_user_service.LoginUserResponse{Status: "error"}, err
	}
	user := models.UserModel{
		Email:        in.Email,
		PasswordHash: password_hash,
	}
	user, err = u.userRepository.FindOne(ctx, &user)
	if err != nil {
		logrus.Errorf("Error while creating user")
		return &proto_user_service.LoginUserResponse{Status: "error"}, nil
	}

	return &proto_user_service.LoginUserResponse{Status: "success", Username: user.Nickname, Id: user.ID, Age: string(*user.Age), ImageUrl: *user.ImageUrl}, nil
}
func (u *userService) DeleteUser(ctx context.Context, req *proto_user_service.DeleteUserRequest) (*proto_user_service.DeleteUserResponse, error) {
	name, err := u.userRepository.Delete(ctx, req.Id)
	if err != nil {
		logrus.Errorf("Cannot delete user by id:%s", req.Id)
	}
	res := proto_user_service.DeleteUserResponse{Status: "Success", Id: req.Id, Nickname: name}
	return &res, nil
}

// If switch fields, querry will create new friends pair
func (u *userService) AddFriends(ctx context.Context, usr *proto_user_service.AddFriendsRequest) (*proto_user_service.AddFriendsResponse, error) {
	friend := models.FriendListModel{
		UserID:   usr.UserId_1,
		FriendID: usr.UserId_2,
	}
	chat := models.ChatsModel{
		User_id1: usr.UserId_1,
		User_id2: usr.UserId_2,
		Action:   "CREATE",
	}
	logrus.Info("Create querry")
	info, err := u.friendsRepository.Create(ctx, &friend)
	logrus.Info(info)
	if err != nil {
		logrus.Errorf("Error while adding friends: %s", err)
		res := &proto_user_service.AddFriendsResponse{Success: false, Message: info}
		return res, nil
	}
	res := proto_user_service.AddFriendsResponse{Success: true, Message: "You add new friend"}

	if err := u.kProducer.Produce(chat, viper.GetString("kafka.topic")); err != nil {
		logrus.Error(err)
	}

	logrus.Info(res.Message)
	logrus.Info(res.Success)
	return &res, nil
}

func (u *userService) DeleteFriend(ctx context.Context, req *proto_user_service.DeleteFriendsRequest) (*proto_user_service.DeleteFriendsResponse, error) {
	friends := models.FriendListModel{
		UserID:   req.UserId_1,
		FriendID: req.UserId_2,
	}
	chat := models.ChatsModel{
		User_id1: req.UserId_1,
		User_id2: req.UserId_2,
		Action:   "DELETE",
	}
	if err := u.friendsRepository.Delete(ctx, &friends); err != nil {
		logrus.Errorf("Error while deleting friend. ErrInfo: %s", err)
		return &proto_user_service.DeleteFriendsResponse{Success: false, Message: "Server error"}, nil
	}

	if err := u.kProducer.Produce(chat, viper.GetString("kafka.topic")); err != nil {
		logrus.Error(err)
	}
	res := proto_user_service.DeleteFriendsResponse{Success: true, Message: "You are not friends anymore"}
	return &res, nil

}

// / NEED TO BE TESTED
func (u *userService) ListOfFriends(ctx context.Context, req *proto_user_service.ListOfFriendsRequest) (*proto_user_service.ListOfFriendsResponse, error) {
	var friends = []models.FriendListModel{}
	friends, err := u.friendsRepository.FindAll(context.Background(), req.UserId)
	if err != nil {
		logrus.Errorf("Error while finding all friends. Err %s", err)
		return nil, err
	}
	var protoFriends []*proto_user_service.Friend
	for _, friend := range friends {
		protoFriend := proto_user_service.Friend{
			UserId:   friend.UserID,
			FirendId: friend.FriendID,
		}
		protoFriends = append(protoFriends, &protoFriend)

	}

	res := &proto_user_service.ListOfFriendsResponse{Friend: protoFriends}
	return res, nil
}
