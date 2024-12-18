package pkg

import (
    "net"

    proto_media_service "github.com/FREEGREAT/protos/gen/go/media"
    "github.com/sirupsen/logrus"
    "gomessage.com/media/internal/service"
    "google.golang.org/grpc"
)

type gRPCServer struct {
    addr string
    mediaService *service.MediaService // Зміна на вказівник
}

func NewGRPCServer(addr string, ms *service.MediaService) *gRPCServer { // Зміна на вказівник
    return &gRPCServer{addr: addr, mediaService: ms}
}

func (s *gRPCServer) Run() error {
    lis, err := net.Listen("tcp", s.addr)
    if err != nil {
        logrus.Fatalf("failed to listen: %v", err)
    }

    grpcServer := grpc.NewServer()
    proto_media_service.RegisterMediaServiceServer(grpcServer, s.mediaService)

    logrus.Println("Starting gRPC server on", s.addr)

    return grpcServer.Serve(lis)
}