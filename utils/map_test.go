package utils

import (
	"encoding/json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestStruct2Map(t *testing.T) {
	Convey("test struct to map", t, func() {
		type A struct {
			a int
			b string
			c float64
			d []int
		}

		type B struct {
			a A
			b []string
		}

		a := A{
			a: 1,
			b: "b",
			c: 12.12,
			d: []int{1, 2, 3},
		}

		b := B{
			a: a,
			b: []string{"asd", "qwe"},
		}

		_, err1 := Struct2Map(a)
		_, err2 := Struct2Map(b)
		So(err1, ShouldBeNil)
		So(err2, ShouldBeNil)
	})
}

func TestBytes2Map(t *testing.T) {
	Convey("test []bytes to map", t, func() {
		type A struct {
			a int
			b string
			c float64
			d []int
		}

		type B struct {
			a A
			b []string
		}

		a := A{
			a: 1,
			b: "b",
			c: 12.12,
			d: []int{1, 2, 3},
		}

		b := B{
			a: a,
			b: []string{"asd", "qwe"},
		}

		byte1, err1 := json.Marshal(a)
		byte2, err2 := json.Marshal(b)

		_, err3 := Bytes2Map(byte1)
		_, err4 := Bytes2Map(byte2)

		So(err1, ShouldBeNil)
		So(err2, ShouldBeNil)
		So(err3, ShouldBeNil)
		So(err4, ShouldBeNil)
	})
}
