package utils

import (
	"testing"

	"github.com/Shopify/sarama"
	. "github.com/smartystreets/goconvey/convey"
)

func TestClass(t *testing.T) {
	Convey("test class", t, func() {
		kafkaMessage := KafkaMessage{}
		So(kafkaMessage.Class(), ShouldEqual, "kafka")
	})
}

func TestNewKafkaMessage(t *testing.T) {
	Convey("test new kafka message", t, func() {
		var msg = new(sarama.ConsumerMessage)
		res := NewKafkaMessage(msg)
		So(res, ShouldNotBeNil)
	})
}

func TestCapturePanic(t *testing.T) {
	Convey("test capture panic", t, func() {
		CapturePanic("SENTRY_DSN")
		SkipSo("no return")
	})
}
