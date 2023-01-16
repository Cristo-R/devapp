package integrations

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	. "github.com/smartystreets/goconvey/convey"
	"gitlab.shoplazza.site/shoplaza/cobra/models"
	"gitlab.shoplazza.site/shoplaza/cobra/utils"
	"gitlab.shoplazza.site/xiabing/goat.git/testhelper"
)

func (suite *BaseSuite) TestCreateAccessToken() {
	uid, _ := models.GenerateAppKey(32)
	secretKey, _ := models.GenerateSecretKey(32)
	app := CreateApp(map[string]interface{}{
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

	guides := models.Guides{}
	guides = append(guides, models.Guide{Name: "test", Status: "init"})
	store := CreateStore(map[string]interface{}{
		"Id":       utils.NewUUIDBinary(),
		"OriginId": "7295",
		"Name":     "pjs",
		"Guides":   &guides,
	})

	CreateOauthAccessToken(map[string]interface{}{
		"Id":              uint64(115),
		"ResourceOwnerId": store.Id,
		"ApplicationId":   uint64(115),
		"Token":           "test",
		"RefreshToken":    "test",
		"ExpiresIn":       123,
		"CreatedAt":       time.Now().AddDate(0, 0, 1),
	})

	req := testhelper.NewJSONRequest(
		"POST",
		fmt.Sprintf("%s%s", suite.internalServer.URL(), "/api/internal/apps/access_token"),
		gin.H{
			"client_id":  app.UID,
			"expires_in": 123,
			"scopes":     []string{"read_product", "write_product"},
			"store_id":   7295,
			"slug":       "phj",
		},
	)

	resp := testhelper.RunRequest(req)
	defer resp.Body.Close()

	Convey("create access token success", suite.T(), func() {
		Convey("should return 200", func() {
			So(resp.StatusCode, ShouldEqual, 200)
		})
	})
}

func (suite *BaseSuite) TestRevokeAccessToken() {
	guides := models.Guides{}
	guides = append(guides, models.Guide{Name: "test", Status: "init"})
	store := CreateStore(map[string]interface{}{
		"Id":       utils.NewUUIDBinary(),
		"OriginId": "7295",
		"Name":     "pjs",
		"Guides":   &guides,
	})

	CreateOauthAccessToken(map[string]interface{}{
		"Id":              uint64(115),
		"ResourceOwnerId": store.Id,
		"ApplicationId":   uint64(115),
		"Token":           "test",
		"RefreshToken":    "test",
		"ExpiresIn":       123,
		"CreatedAt":       time.Now().AddDate(0, 0, 1),
	})

	req := testhelper.NewJSONRequest(
		"DELETE",
		fmt.Sprintf("%s%s", suite.internalServer.URL(), "/api/internal/apps/revoke_access_token"),
		gin.H{
			"access_token": "test",
			"store_id":     7295,
		},
	)

	resp := testhelper.RunRequest(req)
	defer resp.Body.Close()

	Convey("revoke access token success", suite.T(), func() {
		Convey("should return 200", func() {
			So(resp.StatusCode, ShouldEqual, 200)
		})
	})
}

func (suite *BaseSuite) TestUpdateAccessToken() {
	guides := models.Guides{}
	guides = append(guides, models.Guide{Name: "test", Status: "init"})
	store := CreateStore(map[string]interface{}{
		"Id":       utils.NewUUIDBinary(),
		"OriginId": "7295",
		"Name":     "pjs",
		"Guides":   &guides,
	})

	CreateOauthAccessToken(map[string]interface{}{
		"Id":              uint64(115),
		"ResourceOwnerId": store.Id,
		"ApplicationId":   uint64(115),
		"Token":           "test",
		"RefreshToken":    "test",
		"ExpiresIn":       123,
		"CreatedAt":       time.Now().AddDate(0, 0, 1),
	})

	req := testhelper.NewJSONRequest(
		"PUT",
		fmt.Sprintf("%s%s", suite.internalServer.URL(), "/api/internal/apps/access_token/test"),
		gin.H{
			"scopes":   []string{"read_product", "write_product"},
			"store_id": 7295,
		},
	)

	resp := testhelper.RunRequest(req)
	defer resp.Body.Close()

	Convey("update access token success", suite.T(), func() {
		Convey("should return 200", func() {
			So(resp.StatusCode, ShouldEqual, 200)
		})
	})
}
