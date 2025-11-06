package xkafka

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var (
	_pConfigs  map[string]ProducerConfig
	_cConfigs  map[string]ConsumerConfig
	_producers map[string]*kafka.Producer
	_consumers map[string]*kafka.Consumer
)

func Init(configs ...Config) error {
	_pConfigs = make(map[string]ProducerConfig, 16)
	_cConfigs = make(map[string]ConsumerConfig, 16)
	for _, kc := range configs {
		for _, pc := range kc.Producers {
			pc.server = kc.Server
			dc := pc.Default()
			_pConfigs[dc.Id] = dc
		}
		for _, cc := range kc.Consumers {
			cc.server = kc.Server
			dc := cc.Default()
			_cConfigs[dc.Id] = dc
		}
	}

	_producers = make(map[string]*kafka.Producer, 8)
	for _, pc := range _pConfigs {
		producer, err := NewProducer(pc)
		if err != nil {
			Finally()
			return err
		}
		_producers[pc.Id] = producer
	}

	_consumers = make(map[string]*kafka.Consumer, 8)
	for _, cc := range _cConfigs {
		consumer, err := NewConsumer(cc)
		if err != nil {
			Finally()
			return err
		}
		_consumers[cc.Id] = consumer
	}

	return nil
}

func Finally() {
	for _, consumer := range _consumers {
		_ = consumer.Close()
	}
	for _, producer := range _producers {
		producer.Close()
	}
}

func GetProducerConfig(id string) ProducerConfig {
	return _pConfigs[id]
}

func GetProducerConfigDefault() ProducerConfig {
	return GetProducerConfig(DefaultId)
}

func GetConsumerConfig(id string) ConsumerConfig {
	return _cConfigs[id]
}

func GetConsumerConfigDefault() ConsumerConfig {
	return GetConsumerConfig(DefaultId)
}

func GetProducer(id string) *kafka.Producer {
	return _producers[id]
}

func GetProducerDefault() *kafka.Producer {
	return GetProducer(DefaultId)
}

func GetConsumer(id string) *kafka.Consumer {
	return _consumers[id]
}

func GetConsumerDefault() *kafka.Consumer {
	return GetConsumer(DefaultId)
}
