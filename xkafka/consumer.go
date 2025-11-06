package xkafka

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func NewConsumer(config ConsumerConfig) (*kafka.Consumer, error) {
	consumer, err := kafka.NewConsumer(config.Map())
	if err != nil {
		return nil, fmt.Errorf("create consumer failed [%v]", err)
	}

	err = consumer.SubscribeTopics(config.Topics, nil)
	if err != nil {
		return nil, fmt.Errorf("subscribe topic failed [%v]", err)
	}

	return consumer, nil
}

func NewConsumer2(sc ServerConfig, cc ConsumerConfig) (*kafka.Consumer, error) {
	cc.server = sc
	return NewConsumer(cc.Default())
}

// func consume(consumer *kafka.Consumer) {
// 	for {
// 		select {
// 		case <-quit:
// 			log.Info("[runConsumeKafkaMessage] Quit.")
// 			return
// 		default:
// 			msg, err := consumer.ReadMessage(-1)
// 			if err != nil {
// 				continue
// 			}
// 			// to do something
// 		}
// 	}
// }
