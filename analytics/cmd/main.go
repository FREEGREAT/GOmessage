package main

import (
	"log"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	consumer "gomessage.com/analytics/internal/handler/kafka"
	clickhouse "gomessage.com/analytics/internal/storage/clickHouse"
	"gomessage.com/analytics/pkg"
)

var addres = []string{
	"localhost:9092",
	"localhost:9093",
}

func main() {
	logrus.Info("Launching analytics microservise")
	pkg.InitConfig()
	conn := pkg.ConnectDB()
	repo := clickhouse.NewAnalyticsRepository(conn)
	logrus.Info("Launching consumers")
	c, err := consumer.NewConsumer(addres, viper.GetString("kafka.topic"), viper.GetString("kafka.consumer-group"), repo)
	if err != nil {
		log.Fatalf("Error starting Kafka consumer: %v", err)
	}
	c.Consuming()
	logrus.Info("Done")
}
