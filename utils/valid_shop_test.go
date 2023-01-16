package utils

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestValidShop(t *testing.T) {
	Convey("test valid shop failed", t, func() {
		var shop string = "baidu.com"
		res := ValidShop(shop)
		So(res, ShouldBeFalse)
	})

	Convey("test valid shop success in dev env", t, func() {
		var shop string = "xxx.preview.shoplazza.com"
		res := ValidShop(shop)
		So(res, ShouldBeTrue)
	})

	Convey("test valid shop success in prod env", t, func() {
		var shop string = "xxx.myshoplaza.com"
		res := ValidShop(shop)
		So(res, ShouldBeTrue)
	})
}
