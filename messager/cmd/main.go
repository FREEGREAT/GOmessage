package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"gommessage.com/messager/internal/hanlders/kafka"
	database "gommessage.com/messager/pkg/Database"
)

var addres = []string{
	"localhost:9092",
	"localhost:9093",
}

var topic = "chat-topic"
var consumerGroup = "chat-consumer-group"

func main() {
	if err := database.SetupDBConnection(); err != nil {
		logrus.Errorf("Error conn db: %w", err)
	}

	c, err := kafka.NewConsumer(addres, topic, consumerGroup)
	if err != nil {
		logrus.Fatal(err)
	}
	go c.Consuming()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
}

// func start(router *httprouter.Router) {
// 	logrus.Info("Starting application")
// 	listener, err := net.Listen("tcp", ":8082")
// 	if err != nil {
// 		panic(err)
// 	}
// 	server := &http.Server{
// 		Handler:      router,
// 		WriteTimeout: 15 * time.Second,
// 		ReadTimeout:  15 * time.Second,
// 	}
// 	logrus.Info("Server is listening on port :8081")
// 	logrus.Fatalln(server.Serve(listener))
// }
