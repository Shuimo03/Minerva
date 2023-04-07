package kafka

import (
	"backend/pkg/util"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"k8s.io/klog/v2"
	"net/url"
	"time"
)

const (
	ClientID               = "clientKafka"
	brokerDialRetryLimit   = 1
	brokerDialRetryWait    = 0
	brokerLeaderRetryLimit = 1
	brokerLeaderRetryWait  = 0
	DialTimeout            = 40 * time.Second
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

func NewKafkaClient(uri *url.URL) (*kafka, error) {
	opts, err := url.ParseQuery(uri.RawQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to parse url %s", err)
	}
	brokers, brokerError := getBrokers(opts)
	if brokerError != nil {
		return nil, fmt.Errorf("failed to parse Broker %s", brokerError)
	}
	compression, compressionError := getCompression(opts)
	if compressionError != nil {
		//fmt.Errorf(compressionError)
		klog.Infoln(compression)
	}
	topic := getTopic(opts)

	config := sarama.NewConfig()
	config.ClientID = ClientID
	config.Net.DialTimeout = DialTimeout
	config.Metadata.Retry.Max = brokerDialRetryLimit
	config.Metadata.Retry.Backoff = brokerDialRetryWait
	config.Producer.Retry.Max = brokerLeaderRetryLimit
	config.Producer.Retry.Backoff = brokerLeaderRetryWait
	config.Producer.Compression = compression
	config.Producer.Partitioner = sarama.NewRoundRobinPartitioner
	config.Producer.Return.Errors = true
	config.Producer.Return.Successes = true
	config.Net.SASL.User, config.Net.SASL.Password, config.Net.SASL.Enable, err = getSASLConfiguration(opts)
	if err != nil {
		return nil, err
	}

	config.Net.SASL.User, config.Net.SASL.Password, config.Net.SASL.Enable, err = getSASLConfiguration(opts)
	if err != nil {
		return nil, err
	}

	config.Net.SASL.SCRAMClientGeneratorFunc, config.Net.SASL.Mechanism, err = getSCRAMClientAlgorithm(opts)
	if err != nil {
		return nil, err
	}

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("Failed to setup Producer: - %v", err)
	}
	return &kafka{
		producer,
		topic,
	}, nil
}

func getBrokers(opts url.Values) ([]string, error) {
	if len(opts["brokers"]) < 1 {
		return nil, fmt.Errorf("broker is NULL!!!!!")
	}

	brokerList := make([]string, 0)
	brokerList = append(brokerList, opts["brokers"]...)
	return brokerList, nil
}

func getSASLConfiguration(opts url.Values) (string, string, bool, error) {
	if len(opts["user"]) == 0 {
		return "", "", false, nil
	}
	user := opts["user"][0]
	if len(opts["password"]) == 0 {
		return "", "", false, nil
	}
	password := opts["password"][0]
	return user, password, true, nil
}

func getCompression(opts url.Values) (sarama.CompressionCodec, error) {
	compression := opts["compression"]
	if len(compression) == 0 {
		return sarama.CompressionNone, nil
	}
	toLowerCompersion := []byte(compression[0])
	toLowerCompersion = bytes.ToLower(toLowerCompersion)
	alg := string(toLowerCompersion)
	switch alg {
	case "none":
		return sarama.CompressionNone, nil
	case "gzip":
		return sarama.CompressionGZIP, nil
	case "snappy":
		return sarama.CompressionSnappy, nil
	case "lz4":
		return sarama.CompressionLZ4, nil
	case "zstd":
		return sarama.CompressionZSTD, nil
	default:
		return sarama.CompressionNone, fmt.Errorf("compression %s is NULL", alg)
	}
}

func getTopic(opts url.Values) string {
	topic := opts["topic"][0]
	if topic == "" {
		return fmt.Sprintf("Topic is NULL")
	}
	return topic
}

func getSCRAMClientAlgorithm(opts url.Values) (func() sarama.SCRAMClient, sarama.SASLMechanism, error) {
	algorithm := opts["mechanism"][0]
	if algorithm == "sha512" {
		return func() sarama.SCRAMClient { return &util.XDGSCRAMClient{HashGeneratorFcn: util.SHA512} }, sarama.SASLTypeSCRAMSHA512, nil
	} else if algorithm == "sha256" {
		return func() sarama.SCRAMClient { return &util.XDGSCRAMClient{HashGeneratorFcn: util.SHA256} }, sarama.SASLTypeSCRAMSHA256, nil

	} else {
		return nil, "", fmt.Errorf("invalid SHA algorithm \"%s\": can be either \"sha256\" or \"sha512\"", algorithm)
	}
}
