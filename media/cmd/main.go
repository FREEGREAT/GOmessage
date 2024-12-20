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

	logrus.Println("Media service is running on port 50051")

	if err := grpcServer.Serve(lis); err != nil {

		logrus.Fatalf("failed to serve: %v", err)

	}
	// mediaService := service.NewMediaService()
	// srv := pkg.NewGRPCServer("localhost:50051", mediaService)
	// srv.Run()

	// var buf bytes.Buffer
	// buf.WriteString("This is a test image content.")
	// if err_cfg != nil {
	// 	logrus.Fatalf("error init config: %s", err_cfg.Error())
	// }
	// _, err := os.Create("test_image.txt")

	// if err != nil {

	// 	logrus.Fatalf("Failed to create test file: %v", err)

	// }
	// defer os.Remove("test_image.txt")

	// var imgFile multipart.File // This should be initialized with your file

	// var imgSize int64
	// // Set this to the size of your file

	// contentType := "image/jpeg" // Set the appropriate content type

	// err = storage.UploadImgFile("test_image.jpg", imgFile, imgSize, contentType)

	// if err != nil {

	// 	logrus.Fatalf("Failed to upload image: %v", err)

	// }

	logrus.Info(viper.GetString("minio.endpoint"))
	logrus.Info(bucket_name)

}
