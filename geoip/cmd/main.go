package main

import (
	"net"

	proto_geoip_service "github.com/FREEGREAT/protos/gen/go/geoip"
	"github.com/sirupsen/logrus"
	"gomessace.com/geoip/internal/service"
	"gomessace.com/geoip/pkg"
	"google.golang.org/grpc"
)

func main() {
	if err := pkg.InitConfig(); err != nil {
		panic("Error init config geoIp main")
	}
	grpcServer := grpc.NewServer()

	geoIpService := service.CreateNewGeoIpService()

	proto_geoip_service.RegisterGeoIpServiceServer(grpcServer, geoIpService)

	lis, err := net.Listen("tcp", ":9024")
	if err != nil {
		logrus.Fatalf("failed to listen: %v", err)
	}
	logrus.Info("Starting Grpc:")
	if err := grpcServer.Serve(lis); err != nil {
		logrus.Fatalf("failed to serve: %v", err)
	}

}
