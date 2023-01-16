package integrations

import (
	"database/sql"
	"fmt"

	. "github.com/smartystreets/goconvey/convey"
	"gitlab.shoplazza.site/shoplaza/cobra/models"
	"gitlab.shoplazza.site/xiabing/goat.git/testhelper"
)

func (suite *BaseSuite) TestGetApps() {
	uid, _ := models.GenerateAppKey(32)
	secretKey, _ := models.GenerateSecretKey(32)
	CreateApp(map[string]interface{}{
		"ID":            uint64(115),
		"Name":          "tesssst",
		"StoreId":       sql.NullString{Valid: false},
		"UID":           uid,
		"Secret":        secretKey,
		"Category":      "partner",
		"Scopes":        "read_product write_product read_order write_order",
		"Link":          "/app_store/plugins/test",
		"Confidential":  true,
		"Icon":          "\"oss/operation/ee5d8f1c0b641fa075539bd140593ba7.svg\"",
		"Embbed":        false,
		"PrivateApp":    false,
		"OauthDancable": true,
		"Subscribable":  false,
		"StoreCharge":   false,
		"NotifyUri":     "",
		"ChargeMinPlan": "",
		//"FreeStores": "",
		"WebhookApiVersion": "2020-07",
		// "status":            "published",
		"Regions": []string{"1", "2"},
	})
	uid, _ = models.GenerateAppKey(32)

	CreateApp(map[string]interface{}{
		"ID":            uint64(116),
		"Name":          "tesssst",
		"StoreId":       sql.NullString{Valid: false},
		"UID":           uid,
		"Secret":        secretKey,
		"Category":      "partner",
		"Scopes":        "read_product write_product read_order write_order",
		"Link":          "/app_store/plugins/test",
		"Confidential":  true,
		"Icon":          "\"oss/operation/ee5d8f1c0b641fa075539bd140593ba7.svg\"",
		"Embbed":        false,
		"PrivateApp":    false,
		"OauthDancable": true,
		"Subscribable":  false,
		"StoreCharge":   false,
		"NotifyUri":     "",
		"ChargeMinPlan": "",
		//"FreeStores": "",
		"WebhookApiVersion": "2020-07",
		// "status":            "published",
		"Regions": []string{"1", "2"},
	})

	req := testhelper.NewJSONRequest(
		"GET",
		fmt.Sprintf("%s%s", suite.internalServer.URL(), "/api/internal/platforms/apps?page=2&limit=20&status=published&name=111&created_at_end=2021-10-03T00%3A00%3A00%2B08%3A00&created_at_start=2021-10-01T00%3A00%3A00%2B08%3A00"),
		nil,
	)

	resp := testhelper.RunRequest(req)
	defer resp.Body.Close()

	Convey("get dev apps success", suite.T(), func() {
		Convey("should return 200", func() {
			So(resp.StatusCode, ShouldEqual, 200)
		})
	})
}
