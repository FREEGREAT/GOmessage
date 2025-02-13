package kafka

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/sirupsen/logrus"
	"gomessage.com/users/internal/models"
)

var errUnknownType = errors.New("unknown event type")

const (
	flushTimeout = 5000
)

type Producer struct {
	producer *kafka.Producer
}

func NewProducer(addr []string) (*Producer, error) {
	conf := &kafka.ConfigMap{
		"bootstrap.servers": strings.Join(addr, ","),
	}

	p, err := kafka.NewProducer(conf)
	if err != nil {
		return nil, fmt.Errorf("error with creating producer: %w", err)
	}
	return &Producer{producer: p}, nil
}

func (p *Producer) Produce(data_message models.ChatsModel, topic string) error {
	var message bytes.Buffer

	enc := gob.NewEncoder(&message)

	err := enc.Encode(data_message)
	if err != nil {
		logrus.Fatalf("encode error:", err)
	}

	kafkaMsg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: message.Bytes(),
		Key:   nil,
	}
	kafkaChan := make(chan kafka.Event)
	if err := p.producer.Produce(kafkaMsg, kafkaChan); err != nil {
		return err
	}

	e := <-kafkaChan
	switch ev := e.(type) {
	case *kafka.Message:
		return nil
	case kafka.Error:
		return ev
	default:
		return errUnknownType
	}
}

func (p *Producer) ProduceAnalytic(data_message models.Analytics, topic string) error {
	var message bytes.Buffer

	enc := gob.NewEncoder(&message)

	err := enc.Encode(data_message)
	if err != nil {
		logrus.Fatalf("encode error:", err)
	}

	kafkaMsg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: message.Bytes(),
		Key:   nil,
	}
	kafkaChan := make(chan kafka.Event)
	if err := p.producer.Produce(kafkaMsg, kafkaChan); err != nil {
		return err
	}

	e := <-kafkaChan
	switch ev := e.(type) {
	case *kafka.Message:
		return nil
	case kafka.Error:
		return ev
	default:
		return errUnknownType
	}
}

func (p *Producer) Close() {
	p.producer.Flush(flushTimeout)
	p.producer.Close()
}
