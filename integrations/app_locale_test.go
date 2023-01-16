package integrations

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.shoplazza.site/shoplaza/cobra/models"
	"gitlab.shoplazza.site/shoplaza/cobra/utils"
	"gitlab.shoplazza.site/xiabing/goat.git/testhelper"

	. "github.com/smartystreets/goconvey/convey"
)

func (suite *BaseSuite) TestCreateAppLocale() {
	uid, _ := models.GenerateAppKey(32)
	secretKey, _ := models.GenerateSecretKey(32)
	app := CreateApp(map[string]interface{}{
		"Name":               "test",
		"StoreId":            sql.NullString{Valid: false},
		"UID":                uid,
		"Secret":             secretKey,
		"Category":           "partner",
		"Scopes":             "read_product write_product read_order write_order",
		"Link":               "/app_store/plugins/test",
		"Confidential":       true,
		"Icon":               "\"oss/operation/ee5d8f1c0b641fa075539bd140593ba7.svg\"",
		"Embbed":             false,
		"PrivateApp":         false,
		"OauthDancable":      true,
		"Subscribable":       false,
		"StoreCharge":        false,
		"NotifyUri":          "",
		"ChargeMinPlan":      "",
		"FreeStores":         sql.NullString{Valid: false},
		"WebhookApiVersion":  "2020-07",
		"DevelopByShoplazza": false,
		"Regions":            []string{"1", "2"},
	})

	req := testhelper.NewJSONRequest(
		"POST",
		fmt.Sprintf("%s%s%d%s", suite.server.URL(), "/api/platform/apps/", app.ID, "/locales"),
		gin.H{
			"name":     "size_chart",
			"subtitle": "test",
			"locale":   "en",
			"Desc":     "test test test",
		},
	)
	req.Header.Set("Origin", "platform")
	req.Header.Set("Platform-User-Id", "test")
	userPermissions, _ := json.Marshal(map[string]interface{}{
		"apps": map[string]interface{}{
			"read":  true,
			"write": true,
		},
	})
	req.Header.Set("Platform-User-Permissions", string(userPermissions))

	resp := testhelper.RunRequest(req)
	defer resp.Body.Close()

	Convey("create app success", suite.T(), func() {
		Convey("should return 201", func() {
			So(resp.StatusCode, ShouldEqual, 201)
		})
	})
}

func (suite *BaseSuite) TestGetAppLocaleSuccess() {
	uid, _ := models.GenerateAppKey(32)
	secretKey, _ := models.GenerateSecretKey(32)
	app := CreateApp(map[string]interface{}{
		"ID":            uint64(1),
		"Name":          "test",
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
		"Regions":           []string{"1", "2"},
	})
	CreateDevLocal(map[string]interface{}{
		"Id":            uint64(1),
		"Locale":        "en",
		"ApplicationId": app.ID,
	})

	req := testhelper.NewJSONRequest(
		"GET",
		fmt.Sprintf("%s%s", suite.server.URL(), "/api/platform/apps/1/locales"),
		nil,
	)
	req.Header.Set("Origin", "platform")
	req.Header.Set("Platform-User-Id", "test")
	req.Header.Set("Store-ID", strconv.FormatUint(1, 10))
	userPermissions, _ := json.Marshal(map[string]interface{}{
		"apps": map[string]interface{}{
			"read":  true,
			"write": false,
		},
	})
	req.Header.Set("Platform-User-Permissions", string(userPermissions))

	resp := testhelper.RunRequest(req)
	defer resp.Body.Close()

	Convey("get app local success by app id", suite.T(), func() {
		Convey("should return 200", func() {
			So(resp.StatusCode, ShouldEqual, 200)
		})
		Convey("should return correct app info", func() {
			js := testhelper.JSONResp(resp)
			So(js.Get("count").MustUint64(), ShouldEqual, 1)
		})
	})
}

func (suite *BaseSuite) TestUpdateAppLocaleSuccess() {
	uid, _ := models.GenerateAppKey(32)
	secretKey, _ := models.GenerateSecretKey(32)
	app := CreateApp(map[string]interface{}{
		"ID":            uint64(1),
		"Name":          "test",
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
		"Regions":           []string{"1", "2"},
	})
	CreateDevLocal(map[string]interface{}{
		"Id":            uint64(1),
		"Locale":        "en-US",
		"ApplicationId": app.ID,
	})

	req := testhelper.NewJSONRequest(
		"PATCH",
		fmt.Sprintf("%s%s", suite.server.URL(), "/api/platform/apps/1/locales/en-US"),
		gin.H{
			"name": "size_chart",
		},
	)
	req.Header.Set("Origin", "platform")
	req.Header.Set("Platform-User-Id", "test")
	req.Header.Set("Store-ID", strconv.FormatUint(1, 10))
	userPermissions, _ := json.Marshal(map[string]interface{}{
		"apps": map[string]interface{}{
			"read":  true,
			"write": true,
		},
	})
	req.Header.Set("Platform-User-Permissions", string(userPermissions))

	resp := testhelper.RunRequest(req)
	defer resp.Body.Close()

	Convey("update app local success by app id", suite.T(), func() {
		Convey("should return 200", func() {
			So(resp.StatusCode, ShouldEqual, 200)
		})
		Convey("should return correct app info", func() {
			js := testhelper.JSONResp(resp)
			So(js.Get("name").MustString(), ShouldEqual, "size_chart")
			So(js.Get("locale").MustString(), ShouldEqual, "en-US")
		})
	})
}

func (suite *BaseSuite) TestDeleteAppLocaleSuccess() {
	uid, _ := models.GenerateAppKey(32)
	secretKey, _ := models.GenerateSecretKey(32)
	app := CreateApp(map[string]interface{}{
		"ID":            uint64(1),
		"Name":          "test",
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
		"Regions":           []string{"1", "2"},
	})
	CreateLocal(map[string]interface{}{
		"Id":            utils.NewUUIDBinary(),
		"Locale":        "en",
		"ApplicationId": app.ID,
	})

	req := testhelper.NewJSONRequest(
		"DELETE",
		fmt.Sprintf("%s%s", suite.server.URL(), "/api/platform/apps/1/locales/en"),
		nil,
	)
	req.Header.Set("Origin", "platform")
	req.Header.Set("Platform-User-Id", "test")
	req.Header.Set("Store-ID", strconv.FormatUint(1, 10))
	userPermissions, _ := json.Marshal(map[string]interface{}{
		"apps": map[string]interface{}{
			"read":  true,
			"write": true,
		},
	})
	req.Header.Set("Platform-User-Permissions", string(userPermissions))

	resp := testhelper.RunRequest(req)
	defer resp.Body.Close()

	Convey("delete app local success by app id and local", suite.T(), func() {
		Convey("should return 200", func() {
			So(resp.StatusCode, ShouldEqual, 200)
		})
	})
}
