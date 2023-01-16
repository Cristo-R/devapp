package utils

import (
	"context"
	"fmt"

	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"

	"gitlab.shoplazza.site/shoplaza/cobra/config"
)

func SendMessage(topic string, message []byte) error {
	if config.Cfg.Env == "test" {
		return nil
	}

	msg := &sarama.ProducerMessage{Topic: topic, Value: sarama.ByteEncoder(message)}
	partition, offset, err := config.Producer.SendMessage(msg)
	if err != nil {
		log.Printf("FAILED to send message: %s\n", err)
	} else {
		log.Printf("> message sent to partition %d at offset %d\n", partition, offset)
	}
	return err
}

type GroupConsumer struct {
	ProcessFunc func(*sarama.ConsumerMessage) error
	KafkaServer []string
	ListenTopic []string
	GroupId     string
}

func (GroupConsumer) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (GroupConsumer) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h *GroupConsumer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	var msg *sarama.ConsumerMessage
	if config.Cfg.SentryDSN != "" {
		defer func() {
			if msg != nil {
				CapturePanic(config.Cfg.SentryDSN, NewKafkaMessage(msg))
			} else {
				CapturePanic(config.Cfg.SentryDSN)
			}
		}()
	}

	for msg = range claim.Messages() {
		log.Printf("Message topic:%q partition:%d offset:%d timestamp:%v\n", msg.Topic, msg.Partition, msg.Offset, msg.Timestamp)
		err := h.ProcessFunc(msg)
		if err != nil {
			log.WithError(err).Warn("consumer err")
		}
		sess.MarkMessage(msg, "")
	}
	return nil
}

func (h *GroupConsumer) Run(ctx context.Context) error {
	if h.GroupId == "" || len(h.ListenTopic) < 1 || len(h.KafkaServer) < 1 {
		panic(fmt.Sprintf("kafka config invalid, GroupId:%s ListenTopic:%s KafkaServer:%v\n", h.GroupId, h.ListenTopic, h.KafkaServer))
	}

	saramaConfig := sarama.NewConfig()
	saramaConfig.Consumer.Return.Errors = true
	saramaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest
	saramaConfig.Version = sarama.V1_1_0_0

	client, err := sarama.NewClient(h.KafkaServer, saramaConfig)
	if err != nil {
		return err
	}
	defer func() {
		_ = client.Close()
	}()

	group, err := sarama.NewConsumerGroupFromClient(h.GroupId, client)
	if err != nil {
		return err
	}

	for {
		err = group.Consume(ctx, h.ListenTopic, h)
		if err != nil {
			return err
		}
	}

	return nil
}
