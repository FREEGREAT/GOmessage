package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"time"

	proto_media_service "github.com/FREEGREAT/protos/gen/go/media"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gomessage.com/users/internal/grpcclient"
	handler "gomessage.com/users/internal/handlers"
	"gomessage.com/users/internal/service"
	repo "gomessage.com/users/internal/storage/postgresql"
	"gomessage.com/users/pkg/postgresql"
	"gomessage.com/users/pkg/utils"
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

	userRepo := repo.NewUserRepository(postgresqlClient)
	userService := service.CreateNewUserService(userRepo)

	router := httprouter.New()
	logrus.Info("Registering handlers")

	if err != nil {
		logrus.Fatalf("Error while creating GRPCClient", err)
	}

	grpcClient, err := grpcclient.NewGRPCClient(viper.GetString("grpc.addr"), &userService)

	if err != nil {
		logrus.Fatalf("Failed to connect to gRPC server: %s. Host:%s ", err, viper.GetString("grpc.host")+viper.GetString("grpc.port"))
	}
	conn, err:=grpcclient.NewGRPCConn(viper.GetString("grpc.addr"))
	client := proto_media_service.NewMediaServiceClient(conn)
	
	userHandler := handler.NewUserHandler(grpcClient,client )
	userHandler.Register(router)
	start(router)

}

func start(router *httprouter.Router) {
	logrus.Info("Starting application")
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	logrus.Info("Server is listening on port :8080")
	logrus.Fatalln(server.Serve(listener))
}
