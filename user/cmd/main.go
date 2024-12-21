package main

import (
	"context"
	"net"
	"os"

	proto_user_service "github.com/FREEGREAT/protos/gen/go/user"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gomessage.com/users/internal/service"
	repo "gomessage.com/users/internal/storage/postgresql"
	"gomessage.com/users/pkg/postgresql"
	"gomessage.com/users/pkg/utils"
	"google.golang.org/grpc"
)

const connectAttempts = 3

func main() {
	logrus.Info("Initializing application")

	if err := utils.InitConfig(); err != nil {
		logrus.Fatalf("error init config: %s", err.Error())
	}
	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading enb var: %s", err.Error())
	}
	// Створення клієнта для PostgreSQL
	logrus.Info("Connecting to database")
	postgresqlClient, err := postgresql.NewClient(context.TODO(), connectAttempts, postgresql.StorageConfig{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		Database: viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})
	if err != nil {
		logrus.Fatalf("Failed to initialize database client: %s", err)
	}

	grpcServer := grpc.NewServer()
	userRepo := repo.NewUserRepository(postgresqlClient)
	friendRepo := repo.NewFriendsRepository(postgresqlClient)
	userService := service.CreateNewUserService(userRepo, friendRepo)

	proto_user_service.RegisterUserServiceServer(grpcServer, userService)

	lis, err := net.Listen("tcp", ":50051")

	if err != nil {

		logrus.Fatalf("failed to listen: %v", err)

	}

	logrus.Println("Media service is running on port 50051")

	if err := grpcServer.Serve(lis); err != nil {

		logrus.Fatalf("failed to serve: %v", err)

		// router := httprouter.New()
		// logrus.Info("Registering handlers")

		// if err != nil {
		// 	logrus.Fatalf("Error while creating GRPCClient", err)
		// }

		// grpcClient, err := grpcclient.NewGRPCClient(viper.GetString("grpc.addr"), &userService)

		// if err != nil {
		// 	logrus.Fatalf("Failed to connect to gRPC server: %s. Host:%s ", err, viper.GetString("grpc.host")+viper.GetString("grpc.port"))
		// }
		// conn, err:=grpcclient.NewGRPCConn(viper.GetString("grpc.addr"))
		// client := proto_media_service.NewMediaServiceClient(conn)

		// userHandler := handler.NewUserHandler(grpcClient,client )
		// userHandler.Register(router)
		// start(router)

	}
}
