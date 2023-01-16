package integrations

import (
	"encoding/json"
	"fmt"
	"net/http"

	. "github.com/smartystreets/goconvey/convey"
	"gitlab.shoplazza.site/shoplaza/cobra/config"
	"gitlab.shoplazza.site/shoplaza/cobra/forms/payments"
	"gitlab.shoplazza.site/shoplaza/cobra/models"
	"gitlab.shoplazza.site/shoplaza/cobra/utils"
	"gitlab.shoplazza.site/xiabing/goat.git/testhelper"
)

const getPaymentsAppsApi = "/api/internal/payments_apps"

func (suite *BaseSuite) TestGetPaymentApps() {
	var app1 models.Application
	testhelper.MakeRecord(config.DB, &app1, map[string]interface{}{
		"ID":       uint64(1),
		"Name":     "app1",
		"Category": "xxx",
		"Icon":     "\"/app_icon1\"",
		"Regions":  []string{"1", "2"},
	})

	var app2 models.Application
	testhelper.MakeRecord(config.DB, &app2, map[string]interface{}{
		"ID":       uint64(2),
		"Name":     "app2",
		"Category": models.PaymentCategory,
		"Icon":     "\"/app_icon2\"",
		"Regions":  []string{"1", "2"},
	})

	var app3 models.Application
	testhelper.MakeRecord(config.DB, &app3, map[string]interface{}{
		"ID":       uint64(3),
		"Name":     "app3",
		"Category": models.PaymentCategory,
		"Icon":     "\"/app_icon3\"",
		"Regions":  []string{"1", "2"},
	})

	var store models.Store
	guides := models.Guides{}
	guides = append(guides, models.Guide{Name: "test", Status: "init"})
	testhelper.MakeRecord(config.DB, &store, map[string]interface{}{
		"Id":       utils.NewUUIDBinary(),
		"OriginId": "1234",
		"Guides":   &guides,
	})

	var it models.InstallTracks
	testhelper.MakeRecord(config.DB, &it, map[string]interface{}{
		"Id":            utils.NewUUIDBinary(),
		"ApplicationId": app2.ID,
		"StoreId":       store.Id,
	})

	var pe1 models.PaymentExtension
	testhelper.MakeRecord(config.DB, &pe1, map[string]interface{}{
		"AppId":             app2.ID,
		"PaymentSessionUrl": "https://apppay.com",
		"AppIcon":           "/icon",
	})

	var pe2 models.PaymentExtension
	testhelper.MakeRecord(config.DB, &pe2, map[string]interface{}{
		"AppId":             app3.ID,
		"PaymentSessionUrl": "https://apppay2.com",
		"AppIcon":           "/icon2",
	})

	var aa1 models.AccountAlias
	testhelper.MakeRecord(config.DB, &aa1, map[string]interface{}{
		"AppId":     app2.ID,
		"AccountId": "商户号",
		"SignKey":   "密钥",
		"KeyA":      "网管接入号",
		"Locale":    "zh-CN",
	})

	var aa2 models.AccountAlias
	testhelper.MakeRecord(config.DB, &aa2, map[string]interface{}{
		"AppId":     app3.ID,
		"AccountId": "商户号2",
		"SignKey":   "密钥2",
		"KeyA":      "网管接入号2",
		"Locale":    "zh-CN",
	})

	Convey("get payments apps", suite.T(), func() {
		req := testhelper.NewJSONRequest(
			http.MethodGet,
			fmt.Sprintf(suite.internalServer.URL()+getPaymentsAppsApi),
			nil,
		)

		req.Header.Add("store-id", "1234")

		resp := testhelper.RunRequest(req)
		So(resp.StatusCode, ShouldEqual, http.StatusOK)
		defer resp.Body.Close()

		var data payments.GetPaymentsAppsResp
		err := json.Unmarshal([]byte(testhelper.StringResp(resp)), &data)
		So(err, ShouldBeNil)

		So(len(data.Apps), ShouldEqual, 2)
		So(data.Apps[0].ID, ShouldEqual, app2.ID)
		So(data.Apps[0].Installed, ShouldBeTrue)
		So(data.Apps[0].Icon, ShouldEqual, "/app_icon2")
		So(data.Apps[0].Extension.AppId, ShouldEqual, app2.ID)
		So(data.Apps[0].Extension.PaymentSessionUrl, ShouldEqual, "https://apppay.com")
		So(data.Apps[0].Extension.AppIcon, ShouldEqual, "/icon")
		So(data.Apps[0].Extension.AccountAlias.AccountId, ShouldEqual, "商户号")

		So(data.Apps[1].ID, ShouldEqual, app3.ID)
		So(data.Apps[1].Installed, ShouldBeFalse)
		So(data.Apps[1].Extension.AppId, ShouldEqual, app3.ID)
		So(data.Apps[1].Icon, ShouldEqual, "/app_icon3")
		So(data.Apps[1].Extension.PaymentSessionUrl, ShouldEqual, "https://apppay2.com")
		So(data.Apps[1].Extension.AppIcon, ShouldEqual, "/icon2")
		So(data.Apps[1].Extension.AccountAlias.AccountId, ShouldEqual, "商户号2")
	})
}
