package utils

import (
	"errors"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/Shopify/sarama"
	"github.com/getsentry/raven-go"
)

type KafkaMessage struct {
	Value     string    `json:"value"`
	Key       string    `json:"key"`
	Topic     string    `json:"topic"`
	Partition int32     `json:"partition"`
	Offset    int64     `json:"offset"`
	Timestamp time.Time `json:"timestamp"`
}

func (*KafkaMessage) Class() string {
	return "kafka"
}

func NewKafkaMessage(m *sarama.ConsumerMessage) *KafkaMessage {
	return &KafkaMessage{
		Key:       string(m.Key),
		Value:     string(m.Value),
		Topic:     m.Topic,
		Partition: m.Partition,
		Offset:    m.Offset,
		Timestamp: m.Timestamp,
	}
}

func CapturePanic(sentryDSN string, interfaces ...raven.Interface) {
	if rval := recover(); rval != nil {
		raven.SetDSN(sentryDSN)
		client := raven.DefaultClient

		debug.PrintStack()
		rvalStr := fmt.Sprint(rval)
		client.CaptureMessage(rvalStr,
			make(map[string]string),
			append(interfaces, raven.NewException(errors.New(rvalStr), raven.NewStacktrace(2, 3, nil)))...,
		)
	}

}
