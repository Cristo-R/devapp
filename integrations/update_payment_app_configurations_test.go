package integrations

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	. "github.com/smartystreets/goconvey/convey"
	"gitlab.shoplazza.site/shoplaza/cobra/config"
	"gitlab.shoplazza.site/shoplaza/cobra/models"
	"gitlab.shoplazza.site/xiabing/goat.git/testhelper"
)

func (suite *BaseSuite) TestUpdatePaymentAppConfiguration() {
	uid, _ := models.GenerateAppKey(32)
	secretKey, _ := models.GenerateSecretKey(32)
	var app1 models.Application
	testhelper.MakeRecord(config.DB, &app1, map[string]interface{}{
		"ID":       uint64(1),
		"UID":      uid,
		"Secret":   secretKey,
		"Name":     "app1",
		"Category": "payments",
		"Icon":     "\"/app_icon1\"",
		"Regions":  []string{"1", "2"},
	})

	var pe1 models.PaymentExtension
	testhelper.MakeRecord(config.DB, &pe1, map[string]interface{}{
		"AppId":             app1.ID,
		"PaymentSessionUrl": "https://apppay.com",
		"AppIcon":           "/icon",
	})

	Convey("update payment app configuration", suite.T(), func() {
		ready := true
		req := testhelper.NewJSONRequest(
			http.MethodPut,
			fmt.Sprintf(suite.internalServer.URL()+fmt.Sprintf("/api/internal/payments_apps/%s/configure", uid)),
			gin.H{
				"external_handle": "external_handle_test",
				"ready":           &ready,
			},
		)
		req.Header.Add("Store-Id", "7295")

		resp := testhelper.RunRequest(req)
		So(resp.StatusCode, ShouldEqual, http.StatusOK)
		defer resp.Body.Close()

		type Configuration struct {
			ExternalHandle string `json:"external_handle"`
			Ready          *bool  `json:"ready"`
		}

		type Result struct {
			Configuration Configuration `json:"configuration"`
		}
		data := Result{}

		err := json.Unmarshal([]byte(testhelper.StringResp(resp)), &data)
		So(err, ShouldBeNil)
		So(data.Configuration.ExternalHandle, ShouldEqual, "external_handle_test")
		So(*data.Configuration.Ready, ShouldEqual, true)
	})
}
