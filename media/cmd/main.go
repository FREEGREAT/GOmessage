package main

import (
	"net"

	proto_media_service "github.com/FREEGREAT/protos/gen/go/media"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gomessage.com/media/internal/service"
	miniorepo "gomessage.com/media/internal/storage/minio"
	"gomessage.com/media/pkg"
	"google.golang.org/grpc"
)

func main() {
	logrus.Info("Init config")
	if err_cfg := pkg.InitConfig(); err_cfg != nil {
		logrus.Errorf("Error while init config. %v", err_cfg)
	}

	lis, err := net.Listen("tcp", ":9023")
	if err != nil {
		logrus.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	minioRepo := miniorepo.NewMediaRepository()
	mediaService := service.NewMediaService(minioRepo)

	proto_media_service.RegisterMediaServiceServer(grpcServer, mediaService)
	logrus.Info(viper.GetString("minio.endpoint"))
	logrus.Info(viper.GetString("minio.storage"))
	logrus.Info("ServeGRPC")
	if err := grpcServer.Serve(lis); err != nil {
		logrus.Fatalf("failed to serve: %v", err)
	}

}
