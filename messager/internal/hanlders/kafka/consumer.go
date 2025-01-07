package kafka

import (
	"bytes"
	"encoding/gob"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/sirupsen/logrus"
	"gommessage.com/messager/internal/models"
	"gommessage.com/messager/internal/storage"
)

const (
	sesstionTimeout = 7000 //ms
	noTimeout       = -1
	nilID           = ""
)

type Consumer struct {
	consumer *kafka.Consumer
	chatRepo storage.ChatRepository
}

func NewConsumer(addr []string, topic, consumerGroup string, repo storage.ChatRepository) (*Consumer, error) {
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

		var chatsmod models.ChatsModel
		decoder := gob.NewDecoder(bytes.NewReader(kafkaMsg.Value))
		err = decoder.Decode(&chatsmod)
		if err != nil {
			logrus.Error("Error decoding message: ", err)
			continue
		}
		if err := c.processChat(chatsmod); err != nil {
			logrus.Error("Error processing chat: ", err)
		}
		logrus.Infof("Message consumed: %+v", chatsmod)
	}
}

func (c *Consumer) processChat(chatsmod models.ChatsModel) error {

	switch chatsmod.Action {
	case "CREATE":
		uid := chatsmod.User_id1
		uid2 := chatsmod.User_id2
		logrus.Infof("Id %s, id %s", uid, uid2)
		if err := c.chatRepo.CreateChat(uid, uid2); err != nil {
			logrus.Error("Error creating chat: ", err)
			return err
		}

		logrus.Infof("Chat created successfully: %+v", chatsmod)
		return nil
	case "DELETE":
		uid := chatsmod.User_id1
		uid2 := chatsmod.User_id2
		logrus.Infof("Id %s, id %s", uid, uid2)
		if err := c.chatRepo.DeleteChat(uid, uid2); err != nil {
			logrus.Error("Error deleting chat: ", err)
			return err
		}
		logrus.Infof("Chat deleted successfully: %+v", chatsmod)
		return nil
	default:
		logrus.Warn("Unknown action: ", chatsmod.Action)
		return nil
	}
}
