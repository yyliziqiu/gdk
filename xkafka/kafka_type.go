package xkafka

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

const DefaultId = "default"

type Config struct {
	Server    ServerConfig
	Producers []ProducerConfig
	Consumers []ConsumerConfig
}

type ServerConfig struct {
	BootstrapServers string // must
	SecurityProtocol string // optional
	SASLUsername     string // optional
	SASLPassword     string // optional
	SASLMechanism    string // optional
	SSLCaLocation    string // optional
}

func (c ServerConfig) Default() ServerConfig {
	return c
}

func (c ServerConfig) Join(m kafka.ConfigMap) *kafka.ConfigMap {
	// 服务地址
	m["bootstrap.servers"] = c.BootstrapServers

	// 认证相关
	if c.SecurityProtocol != "" {
		m["security.protocol"] = c.SecurityProtocol
		switch c.SecurityProtocol {
		case "sasl_plaintext":
			m["sasl.username"] = c.SASLUsername
			m["sasl.password"] = c.SASLPassword
			m["sasl.mechanism"] = c.SASLMechanism
		case "sasl_ssl":
			m["sasl.username"] = c.SASLUsername
			m["sasl.password"] = c.SASLPassword
			m["sasl.mechanism"] = c.SASLMechanism
			m["ssl.ca.location"] = c.SSLCaLocation
		}
	}

	return &m
}

type ProducerConfig struct {
	server                ServerConfig
	Id                    string                     // optional
	Topics                []string                   // optional
	RequestRequiredAcks   int                        // optional
	DeliveredCallback     func(kafka.TopicPartition) `json:"-"` // optional
	DeliverFailedCallback func(kafka.TopicPartition) `json:"-"` // optional
}

func (c ProducerConfig) Default() ProducerConfig {
	c.server.Default()
	if c.Id == "" {
		c.Id = DefaultId
	}
	return c
}

func (c ProducerConfig) Map() *kafka.ConfigMap {
	m := kafka.ConfigMap{
		"request.required.acks": c.RequestRequiredAcks,
	}
	return c.server.Join(m)
}

type ConsumerConfig struct {
	server                 ServerConfig
	Id                     string   // optional
	Group                  string   // must
	Topics                 []string // must
	AutoOffsetReset        string   // optional
	MaxPollIntervalMS      int      // optional
	SessionTimeoutMS       int      // optional
	HeartbeatIntervalMS    int      // optional
	FetchMaxBytes          int      // optional
	MaxPartitionFetchBytes int      // optional
}

func (c ConsumerConfig) Default() ConsumerConfig {
	c.server.Default()
	if c.Id == "" {
		c.Id = DefaultId
	}
	if c.AutoOffsetReset == "" {
		c.AutoOffsetReset = "latest"
	}
	if c.MaxPollIntervalMS == 0 {
		c.MaxPollIntervalMS = 10000 // 10s
	}
	if c.SessionTimeoutMS == 0 {
		c.SessionTimeoutMS = 10000 // 10s
	}
	if c.HeartbeatIntervalMS == 0 {
		c.HeartbeatIntervalMS = 3000 // 3s
	}
	if c.FetchMaxBytes == 0 {
		c.FetchMaxBytes = 1024000 // 1M
	}
	if c.MaxPartitionFetchBytes == 0 {
		c.MaxPartitionFetchBytes = 512000 // 500K
	}
	return c
}

func (c ConsumerConfig) Map() *kafka.ConfigMap {
	m := kafka.ConfigMap{
		"group.id":                  c.Group,
		"auto.offset.reset":         c.AutoOffsetReset,
		"max.poll.interval.ms":      c.MaxPollIntervalMS,
		"session.timeout.ms":        c.SessionTimeoutMS,
		"heartbeat.interval.ms":     c.HeartbeatIntervalMS,
		"fetch.max.bytes":           c.FetchMaxBytes,
		"max.partition.fetch.bytes": c.MaxPartitionFetchBytes,
	}
	return c.server.Join(m)
}
