package config

import (
	"github.com/Shopify/sarama"
)

var Producer sarama.SyncProducer

func InitProducer() {
	if Cfg.Env != "test" {
		var err error
		Producer, err = sarama.NewSyncProducer(Cfg.KafkaHost, nil)
		if err != nil {
			panic(err)
		}
	}
}
