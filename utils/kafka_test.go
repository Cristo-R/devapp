package utils

import (
	"testing"

	"github.com/Shopify/sarama"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSetup(t *testing.T) {
	Convey("test group consumer set up", t, func() {
		groupConsumer := new(GroupConsumer)
		consumerGroupSession := new(sarama.ConsumerGroupSession)
		err := groupConsumer.Setup(*consumerGroupSession)
		So(err, ShouldBeNil)
	})
}

func TestCleanup(t *testing.T) {
	Convey("test group consumer clean up", t, func() {
		groupConsumer := new(GroupConsumer)
		consumerGroupSession := new(sarama.ConsumerGroupSession)
		err := groupConsumer.Cleanup(*consumerGroupSession)
		So(err, ShouldBeNil)
	})
}
