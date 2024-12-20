package main

import (
	"net"
	"net/http"
	"time"

	proto_media_service "github.com/FREEGREAT/protos/gen/go/media"
	proto_user_service "github.com/FREEGREAT/protos/gen/go/user"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gomessage.com/gateway/internal/handler"
	"gomessage.com/gateway/pkg"
)

var err_cfg = pkg.InitConfig()

func main() {
	if err_cfg != nil {
		panic("Error loading config at main.go")
	}
	router := httprouter.New()
	logrus.Info("Registering handlers")

	user_conn, err := pkg.NewGRPCConn(viper.GetString("grpc.user_addr"))
	if err != nil {
		logrus.Errorf("Error while creating User GRPC connect: %w", err)
	}
	media_conn, err := pkg.NewGRPCConn(viper.GetString("grpc.media_addr"))
	if err != nil {
		logrus.Errorf("Error while creating Media GRPC connect: %w", err)
	}
	clientMediaService := proto_media_service.NewMediaServiceClient(media_conn)
	clientUserService := proto_user_service.NewUserServiceClient(user_conn)
	gatehand := handler.NewGatewayHandler(clientUserService, clientMediaService)
	gatehand.Register(router)
	start(router)
}

func start(router *httprouter.Router) {
	logrus.Info("Starting application")
	listener, err := net.Listen("tcp", ":8081")
	if err != nil {
		panic(err)
	}
	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	logrus.Info("Server is listening on port :8081")
	logrus.Fatalln(server.Serve(listener))
}
