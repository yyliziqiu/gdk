package xkafka

import (
	"encoding/json"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func NewProducer(config ProducerConfig) (*kafka.Producer, error) {
	producer, err := kafka.NewProducer(config.Map())
	if err != nil {
		return nil, fmt.Errorf("create producer failed [%v]", err)
	}

	go func(deliveredCallback func(partition kafka.TopicPartition), deliverFailedCallback func(partition kafka.TopicPartition)) {
		for e := range producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error == nil {
					if deliveredCallback != nil {
						deliveredCallback(ev.TopicPartition)
					}
				} else {
					if deliverFailedCallback != nil {
						deliverFailedCallback(ev.TopicPartition)
					}
				}
			}
		}
	}(config.DeliveredCallback, config.DeliverFailedCallback)

	return producer, nil
}

func NewProducer2(sc ServerConfig, pc ProducerConfig) (*kafka.Producer, error) {
	pc.server = sc
	return NewProducer(pc.Default())
}

func Produce(producer *kafka.Producer, topic string, message []byte) error {
	return producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          message,
	}, nil)
}

func ProduceModel(producer *kafka.Producer, topic string, model interface{}) error {
	message, err := json.Marshal(model)
	if err != nil {
		return err
	}
	return Produce(producer, topic, message)
}

func Push(topic string, message []byte) error {
	return Produce(GetProducerDefault(), topic, message)
}

func PushModel(topic string, model interface{}) error {
	return ProduceModel(GetProducerDefault(), topic, model)
}
