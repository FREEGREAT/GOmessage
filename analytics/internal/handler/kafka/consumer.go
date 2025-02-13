package consumer

import (
	"bytes"
	"context"
	"encoding/gob"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/sirupsen/logrus"
	"gomessage.com/analytics/internal/models"
	"gomessage.com/analytics/internal/storage" // Імпортуємо пакет з репозиторієм
)

const (
	sesstionTimeout = 7000 //ms
	noTimeout       = -1
	nilID           = ""
)

type Consumer struct {
	consumer *kafka.Consumer
	chatRepo storage.AnalyticsRepository
}

func NewConsumer(addr []string, topic, consumerGroup string, repo storage.AnalyticsRepository) (*Consumer, error) {
	cfg := &kafka.ConfigMap{
		"bootstrap.servers":        strings.Join(addr, ","),
		"group.id":                 consumerGroup,
		"session.timeout.ms":       sesstionTimeout,
		"enable.auto.offset.store": true,
		"enable.auto.commit":       true,
		"auto.commit.interval.ms":  6000,
		"auto.offset.reset":        "earliest",
	}
	c, err := kafka.NewConsumer(cfg)
	if err != nil {
		return nil, err
	}
	if err = c.Subscribe(topic, nil); err != nil {
		return nil, err
	}
	return &Consumer{consumer: c, chatRepo: repo}, nil
}

func (c *Consumer) Consuming() error {

	for {
		logrus.Info("Waiting for messages...")
		kafkaMsg, err := c.consumer.ReadMessage(noTimeout)
		if err != nil {
			logrus.Error("Error reading message: ", err)
			continue
		}

		if kafkaMsg == nil {
			continue
		}

		var analyticsmod models.Analytics
		decoder := gob.NewDecoder(bytes.NewReader(kafkaMsg.Value))
		err = decoder.Decode(&analyticsmod)
		if err != nil {
			logrus.Error("Error decoding : ", err)
			continue
		}
		if err := c.processAnalytics(analyticsmod); err != nil {
			logrus.Error("Error processing chat: ", err)
		}
		logrus.Infof("Message consumed: %+v", analyticsmod)
	}
}

func (c *Consumer) processAnalytics(aModels models.Analytics) error {
	_, err := c.chatRepo.AddData(context.Background(), &aModels)
	if err != nil {
		logrus.Errorf("Error while writing data in DB")
		return err
	}
	return nil
}
