package main

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gommessage.com/messager/internal/hanlders/kafka"
	wsserver "gommessage.com/messager/internal/server"
	"gommessage.com/messager/internal/storage/cassandra"
	"gommessage.com/messager/pkg"
	database "gommessage.com/messager/pkg/Database"
)

var addres = []string{
	"localhost:9092",
	"localhost:9093",
}

func main() {
	if err := pkg.InitConfig(); err != nil {
		panic("Error init config main.go")
	}
	session, err := database.SetupDBConnection()
	if err != nil {
		logrus.Errorf("Error conn db: %w", err)
	}
	chatRepo := cassandra.NewChatRepository(session)
	c, err := kafka.NewConsumer(addres, viper.GetString("kafka.topic"), viper.GetString("kafka.consumer-group"), chatRepo)
	if err != nil {
		logrus.Fatalf("Error while creating consumer: %w", err)
	}

	wsSrv := wsserver.NewWsServer(viper.GetString("ws.addr"), chatRepo)

	logrus.Info("Server started at: ", viper.GetString("ws.addr"))
	go wsSrv.Start()
	logrus.Info("Consumer started")
	c.Consuming()

}
