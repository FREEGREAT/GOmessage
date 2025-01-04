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

var err_cfg = pkg.InitConfig()

var bucket_name = viper.GetString("minio.storage")

func main() {
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		logrus.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	minioRepo := miniorepo.NewMediaRepository()
	mediaService := service.NewMediaService(minioRepo)

	proto_media_service.RegisterMediaServiceServer(grpcServer, mediaService)

	if err := grpcServer.Serve(lis); err != nil {
		logrus.Fatalf("failed to serve: %v", err)
	}
	logrus.Info(viper.GetString("minio.endpoint"))
	logrus.Info(bucket_name)

}
