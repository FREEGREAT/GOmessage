package server

import (
	"net"

	"github.com/sirupsen/logrus"
	handlers "gomessage.com/gateway/internal/handlers/grpc"
	"google.golang.org/grpc"
)

type gRPCServer struct {
	addr string
}

func NewGRPCServer(addr string) *gRPCServer {
	return &gRPCServer{addr: addr}
}

func (s *gRPCServer) Run() error {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		logrus.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	handlers.NewGrpcMediaService(grpcServer)
	logrus.Println("Starting gRPC server on", s.addr)

	return grpcServer.Serve(lis)
}
