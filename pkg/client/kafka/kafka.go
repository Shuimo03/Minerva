package kafka

import (
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"net/url"
	"time"
)

const (
	eventTopics = ""
	ClientID    = ""
	DialTimeout = 40 * time.Second
)

type kafka struct {
	syncProducer sarama.SyncProducer
	topic        string
}

// TODO 需要实现消费者和消费组
type ProducerMessage interface {
	ProducerMessage(msg interface{}) (bool, error)
}

func (k *kafka) ProducerMessage(msg interface{}) (bool, error) {
	msgJson, err := json.Marshal(msg)
	if err != nil {
		return false, fmt.Errorf("failed to transform the event to json : %s", err)
	}
	_, _, err = k.syncProducer.SendMessage(&sarama.ProducerMessage{
		Topic: k.topic,
		Key:   nil,
		Value: sarama.ByteEncoder(msgJson),
	})
	if err != nil {
		return false, fmt.Errorf("failed syncProduer to kafka: %s", err)
	}
	return true, nil
}

func NewKafkaClient(broker *url.URL) (*kafka, error) {
	opts, err := url.ParseQuery(broker.RawQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to parse url %s", err)
	}
	brokers := make([]string, 3)

	config := sarama.NewConfig()
	config.ClientID = ClientID
	config.Net.DialTimeout = DialTimeout
	config.Producer.Return.Errors = true
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(broker, config)
	if err != nil {

	}
	return &kafka{
		producer,
		"",
	}, nil
}

func getBrokers(rawQuery url.URL) []string {
	opts, err := url.ParseQuery(rawQuery.RawQuery)
	if err != nil {
		fmt.Println(err)
	}
	brokerList := make([]string, 3)
	brokerList = append(brokerList, opts["brokers"]...)
	return brokerList
}

func getTopic(rawQuery url.URL) string {
	opts, err := url.ParseQuery(rawQuery.RawQuery)
	if err != nil {
		fmt.Println(err)
	}
	topic := opts["topic"][0]
	return topic
}

//TODO 添加认证方式
