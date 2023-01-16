package utils

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestValue(t *testing.T) {
	Convey("test uuid binary value", t, func() {
		var id UUIDBinary = "e92e55e1-aa69-4b18-8d64-b51a992eb610"
		res, err := id.Value()
		So(res, ShouldNotBeNil)
		So(err, ShouldBeNil)
	})
}

func TestScan(t *testing.T) {
	Convey("test uuid binary scan", t, func() {
		var id UUIDBinary = "e92e55e1-aa69-4b18-8d64-b51a992eb610"
		var v1 interface{}
		var v2 string = "abc"
		b2 := []byte(v2)

		err1 := id.Scan(v1)
		err2 := id.Scan(b2)

		So(err1, ShouldBeNil)
		So(err2, ShouldNotBeNil)
	})
}

func TestString(t *testing.T) {
	Convey("test uuid binary string", t, func() {
		var id UUIDBinary = "e92e55e1-aa69-4b18-8d64-b51a992eb610"
		res := id.String()
		So(res, ShouldNotBeNil)
	})
}

func TestBinary(t *testing.T) {
	Convey("test uuid binary", t, func() {
		var id UUIDBinary = "e92e55e1-aa69-4b18-8d64-b51a992eb610"
		res, err := id.Binary()
		So(res, ShouldNotBeNil)
		So(err, ShouldBeNil)
	})
}

func TestMustBinary(t *testing.T) {
	Convey("test uuid must binary", t, func() {
		var id UUIDBinary = "e92e55e1-aa69-4b18-8d64-b51a992eb610"
		res := id.MustBinary()
		So(res, ShouldNotBeNil)
	})
}

func TestNewUUIDBinary(t *testing.T) {
	Convey("test new uuid binary", t, func() {
		res := NewUUIDBinary()
		So(res, ShouldNotBeNil)
	})
}
