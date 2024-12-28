package kafka

import (
	"bytes"
	"encoding/binary"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/sirupsen/logrus"
	"gommessage.com/messager/internal/models"
	"gommessage.com/messager/internal/storage"
)

const (
	sesstionTimeout = 7000 //ms
	noTimeout       = -1
)

type Consumer struct {
	consumer *kafka.Consumer
	chatRepo storage.ChatRepository
}

func NewConsumer(addr []string, topic, consumerGroup string) (*Consumer, error) {
	cfg := &kafka.ConfigMap{
		"bootstrap.servers":        strings.Join(addr, ","),
		"group.id":                 consumerGroup,
		"session.timeout.ms":       sesstionTimeout,
		"enable.auto.offset.store": false,
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
	return &Consumer{consumer: c}, nil
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

		var chatsmod models.ChatsModel
		err = binary.Read(bytes.NewReader(kafkaMsg.Value), binary.LittleEndian, &chatsmod)
		if err != nil {
			logrus.Error("Error decoding message: ", err)
			continue 
		}

		if err := c.processChat(chatsmod); err != nil {
			logrus.Error("Error processing chat: ", err)
		}
	}
}

func (c *Consumer) processChat(chatsmod models.ChatsModel) error {
	switch chatsmod.Action {
	case "CREATE":
		return c.chatRepo.CreateChat(&chatsmod)
	case "DELETE":
		return c.chatRepo.DeleteChat(chatsmod.Chat_id)
	default:
		logrus.Warn("Unknown action: ", chatsmod.Action)
		return nil
	}
}
