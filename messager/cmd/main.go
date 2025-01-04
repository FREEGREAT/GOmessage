package main

import (
	"github.com/sirupsen/logrus"
	"gommessage.com/messager/internal/hanlders/kafka"
	wsserver "gommessage.com/messager/internal/server"
	database "gommessage.com/messager/pkg/Database"
)

var addres = []string{
	"localhost:9092",
	"localhost:9093",
}

const topic = "chat-topic"
const consumerGroup = "chat-consumer-group"
const addr = "192.168.0.124:8082"

func main() {

	if err := database.SetupDBConnection(); err != nil {
		logrus.Errorf("Error conn db: %w", err)
	}
	c, err := kafka.NewConsumer(addres, topic, consumerGroup)
	if err != nil {
		logrus.Fatalf("Error while creating consumer: %w", err)
	}
	wsSrv := wsserver.NewWsServer(addr)
	logrus.Info("Server started at: ", addr)
	go wsSrv.Start()
	logrus.Info("Consumer started")
	c.Consuming()

}
