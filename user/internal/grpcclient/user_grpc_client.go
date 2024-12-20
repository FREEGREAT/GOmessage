package grpcclient

// import (
// 	"context"
// 	"fmt"

// 	proto_media_service "github.com/FREEGREAT/protos/gen/go/media"
// 	proto_user_service "github.com/FREEGREAT/protos/gen/go/user"
// 	"github.com/sirupsen/logrus"
// 	"gomessage.com/users/internal/models"
// 	"gomessage.com/users/internal/service"
// 	"google.golang.org/grpc"
// 	"google.golang.org/grpc/credentials/insecure"
// )

// type UserGRPCClient struct {
// 	conn        *grpc.ClientConn
// 	client      proto_media_service.MediaServiceClient
// 	userService service.UserService
// }

// // DeleteUser implements GRPCClient.
// func (c *UserGRPCClient) DeleteUser(ctx context.Context, req *proto_user_service.DeleteUserRequest) (*proto_user_service.DeleteUserResponse, error) {
// 	panic("unimplemented")
// }

// // GetUsers implements GRPCClient.
// func (c *UserGRPCClient) GetUsers(ctx context.Context, req *proto_user_service.GetUsersRequest) (*proto_user_service.GetUsersResponse, error) {
// 	panic("unimplemented")
// }

// // UpdateUser implements GRPCClient.
// func (c *UserGRPCClient) UpdateUser(ctx context.Context, req *proto_user_service.UpdateUserRequest) (*proto_user_service.UpdateUserResponse, error) {
// 	panic("unimplemented")
// }

// // NewGRPCClient повертає інтерфейс GRPCClient
// func NewGRPCClient(address string, userService *service.UserService) (GRPCClient, error) {
// 	conn, err := grpc.Dial(
// 		address,
// 		grpc.WithTransportCredentials(insecure.NewCredentials()),
// 	)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to connect: %v", err)
// 	}

// 	client := proto_media_service.NewMediaServiceClient(conn)

// 	return &UserGRPCClient{
// 		conn:        conn,
// 		client:      client,
// 		userService: *userService,
// 	}, nil
// }

// func NewGRPCConn(address string) (*grpc.ClientConn, error) {
// 	conn, err := grpc.Dial(
// 		address,
// 		grpc.WithTransportCredentials(insecure.NewCredentials()),
// 	)
// 	if err != nil {
// 		panic("failed to connect:")
// 	}
// 	return conn, nil
// }

// // Реалізація методів інтерфейсу GRPCClient
// func (c *UserGRPCClient) RegisterUser(ctx context.Context, req *proto_user_service.RegisterUserRequest) (*proto_user_service.RegisterUserResponse, error) {
// 	age := 9
// 	user := &models.UserModel{
// 		Nickname:     req.Nickname,
// 		PasswordHash: req.Password,
// 		Email:        req.Email,
// 		Age:          &age,
// 		ImageUrl:     &req.ImageUrl,
// 	}

// 	err := c.userService.RegisterUser(ctx, )
// 	if err != nil {
// 		logrus.Fatalf("Error grpc create user: %s", err)
// 	}

// 	res := &proto_user_service.RegisterUserResponse{
// 		UserId: "q23",
// 	}
// 	logrus.Info("Register user complete work")
// 	return res, nil
// }

// func (c *UserGRPCClient) Close() error {
// 	return c.conn.Close()
// }
