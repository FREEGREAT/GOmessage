package main

import (
	"context"
	"net"

	proto_user_service "github.com/FREEGREAT/protos/gen/go/user"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gomessage.com/users/internal/service"
	"gomessage.com/users/internal/service/kafka"
	repo "gomessage.com/users/internal/storage/postgresql"
	"gomessage.com/users/pkg/postgresql"
	"gomessage.com/users/pkg/utils"
	"google.golang.org/grpc"
)

const connectAttempts = 3

var addres = []string{
	"localhost:9092",
	"localhost:9093",
}

func main() {
	logrus.Info("Initializing application")

	if err := utils.InitConfig(); err != nil {
		logrus.Fatalf("error init config: %s", err.Error())
	}
	// Створення клієнта для PostgreSQL
	logrus.Info("Connecting to database")
	postgresqlClient, err := postgresql.NewClient(context.TODO(), connectAttempts, postgresql.StorageConfig{
		Host:     viper.GetString("postgre.host"),
		Port:     viper.GetString("postgre.port"),
		Username: viper.GetString("postgre.username"),
		Password: viper.GetString("postgre.password"),
		Database: viper.GetString("postgre.dbname"),
		SSLMode:  viper.GetString("postgre.sslmode"),
	})
	if err != nil {
		logrus.Fatalf("Failed to initialize database client: %s", err)
	}
	p, err := kafka.NewProducer(addres)
	if err != nil {
		logrus.Fatal("Error while creating producer; Err: %w", err)
	}
	grpcServer := grpc.NewServer()
	userRepo := repo.NewUserRepository(postgresqlClient)
	friendRepo := repo.NewFriendsRepository(postgresqlClient)
	userService := service.CreateNewUserService(userRepo, friendRepo, *p)

	proto_user_service.RegisterUserServiceServer(grpcServer, userService)

	lis, err := net.Listen("tcp", viper.GetString("grpc.addr"))
	if err != nil {
		logrus.Fatalf("failed to listen: %v", err)
	}

	if err := grpcServer.Serve(lis); err != nil {
		logrus.Fatalf("failed to serve: %v", err)
	}
}
