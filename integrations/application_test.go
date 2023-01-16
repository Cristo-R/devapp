package integrations

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.shoplazza.site/shoplaza/cobra/config"
	"gitlab.shoplazza.site/shoplaza/cobra/forms/platform"

	//"gitlab.shoplazza.site/shoplaza/cobra/config"
	//"net/http"
	//"strconv"
	//"strings"
	//"time"
	"gitlab.shoplazza.site/shoplaza/cobra/models"
	"gitlab.shoplazza.site/xiabing/goat.git/testhelper"

	. "github.com/smartystreets/goconvey/convey"

	"encoding/json"
)

func (suite *BaseSuite) TestCreateAppSuccess() {
	appLocal := make([]platform.AppLocaleForm, 3)
	testhelper.HaveRecord(config.DB, &appLocal[0], map[string]interface{}{
		"Locale": "en",
	})
	testhelper.HaveRecord(config.DB, &appLocal[1], map[string]interface{}{
		"Locale": "en-US",
	})
	testhelper.HaveRecord(config.DB, &appLocal[2], map[string]interface{}{
		"Locale": "zh-CN",
	})
	scopes := [...]string{"read_product", "write_product", "read_order", "write_order", "read_script_tags", "write_script_tags"}
	req := testhelper.NewJSONRequest(
		"POST",
		fmt.Sprintf("%s%s", suite.server.URL(), "/api/platform/apps"),
		gin.H{
			"name":                "size_chart",
			"redirect_uri":        "https://size-chart1024.apps.shoplazza.com/callback/shoplazza/oauth",
			"scopes":              scopes,
			"icon":                "oss/operation/d57d9bfb5eae56c09fdd1dca7d1736c3.svg",
			"email":               "test@shoplazza.com",
			"app_uri":             "http://baidu.com",
			"webhook_api_version": "2020-07",
			"app_locales":         appLocal,
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

func (suite *BaseSuite) TestCreateAppWithoutPermissions() {
	appLocal := make([]platform.AppLocaleForm, 3)
	testhelper.HaveRecord(config.DB, &appLocal[0], map[string]interface{}{
		"Locale": "en",
	})
	testhelper.HaveRecord(config.DB, &appLocal[1], map[string]interface{}{
		"Locale": "en-US",
	})
	testhelper.HaveRecord(config.DB, &appLocal[2], map[string]interface{}{
		"Locale": "zh-CN",
	})
	scopes := [...]string{"read_product", "write_product", "read_order", "write_order", "read_script_tags", "write_script_tags"}
	req := testhelper.NewJSONRequest(
		"POST",
		fmt.Sprintf("%s%s", suite.server.URL(), "/api/platform/apps"),
		gin.H{
			"name":                "size_chart",
			"redirect_uri":        "https://size-chart1024.apps.shoplazza.com/callback/shoplazza/oauth",
			"scopes":              scopes,
			"icon":                "oss/operation/d57d9bfb5eae56c09fdd1dca7d1736c3.svg",
			"email":               "test@shoplazza.com",
			"app_uri":             "http://baidu.com",
			"webhook_api_version": "2020-07",
			"app_locales":         appLocal,
		},
	)
	req.Header.Set("Origin", "platform")
	req.Header.Set("Platform-User-Id", "test")
	userPermissions, _ := json.Marshal(map[string]interface{}{
		"apps": map[string]interface{}{
			"read":  true,
			"write": false,
		},
	})
	req.Header.Set("Platform-User-Permissions", string(userPermissions))

	resp := testhelper.RunRequest(req)
	defer resp.Body.Close()

	Convey("create app fail because this user is no permission ", suite.T(), func() {
		Convey("should return 403", func() {
			So(resp.StatusCode, ShouldEqual, 403)
		})
	})
}

func (suite *BaseSuite) TestCreateAppWithoutRequiredField() {
	appLocal := make([]platform.AppLocaleForm, 3)
	testhelper.HaveRecord(config.DB, &appLocal[0], map[string]interface{}{
		"Locale": "en",
	})
	testhelper.HaveRecord(config.DB, &appLocal[1], map[string]interface{}{
		"Locale": "en-US",
	})
	testhelper.HaveRecord(config.DB, &appLocal[2], map[string]interface{}{
		"Locale": "zh-CN",
	})
	scopes := [...]string{"read_product", "write_product", "read_order", "write_order", "read_script_tags", "write_script_tags"}
	req := testhelper.NewJSONRequest(
		"POST",
		fmt.Sprintf("%s%s", suite.server.URL(), "/api/platform/apps"),
		gin.H{
			//"name": "size_chart",
			"redirect_uri":        "https://size-chart1024.apps.shoplazza.com/callback/shoplazza/oauth",
			"scopes":              scopes,
			"icon":                "oss/operation/d57d9bfb5eae56c09fdd1dca7d1736c3.svg",
			"email":               "test@shoplazza.com",
			"app_uri":             "",
			"webhook_api_version": "2020-07",
			"app_locales":         appLocal,
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

	Convey("create app fail because some param of required is not included", suite.T(), func() {
		Convey("should return 422", func() {
			So(resp.StatusCode, ShouldEqual, 422)
		})
	})
}

func (suite *BaseSuite) TestGetAppById() {
	uid, _ := models.GenerateAppKey(32)
	secretKey, _ := models.GenerateSecretKey(32)
	app := CreateApp(map[string]interface{}{
		"Name":           "test",
		"StoreId":        sql.NullString{Valid: false},
		"UID":            uid,
		"Secret":         secretKey,
		"Category":       "partner",
		"Scopes":         "read_product write_product read_order write_order",
		"Link":           "/app_store/plugins/test",
		"Confidential":   true,
		"Icon":           "\"oss/operation/ee5d8f1c0b641fa075539bd140593ba7.svg\"",
		"Embbed":         false,
		"PrivateApp":     false,
		"OauthDancable":  true,
		"Subscribable":   false,
		"StoreCharge":    false,
		"NotifyUri":      "",
		"ChargeMinPlan":  "",
		"AdditionalInfo": "",
		//"FreeStores": "",
		"WebhookApiVersion":  "2020-07",
		"DevelopByShoplazza": false,
		"Regions":            []string{"1", "2"},
	})

	req := testhelper.NewJSONRequest(
		"GET",
		fmt.Sprintf("%s%s%d", suite.server.URL(), "/api/platform/apps/", app.ID),
		nil,
	)
	req.Header.Set("Origin", "platform")
	req.Header.Set("Platform-User-Id", "test")
	userPermissions, _ := json.Marshal(map[string]interface{}{
		"apps": map[string]interface{}{
			"read":  true,
			"write": false,
		},
	})
	req.Header.Set("Platform-User-Permissions", string(userPermissions))

	resp := testhelper.RunRequest(req)
	defer resp.Body.Close()

	Convey("get app info by id success", suite.T(), func() {
		Convey("should return 200", func() {
			So(resp.StatusCode, ShouldEqual, 200)
		})
		Convey("should return correct app info", func() {
			js := testhelper.JSONResp(resp)
			So(js.Get("app").Get("id").MustUint64(), ShouldEqual, app.ID)
			So(js.Get("app").Get("uid").MustString(), ShouldEqual, app.UID)
		})
	})
}

func (suite *BaseSuite) TestGetAppByIdFail() {
	uid, _ := models.GenerateAppKey(32)
	secretKey, _ := models.GenerateSecretKey(32)
	app := CreateApp(map[string]interface{}{
		"Name":           "test",
		"StoreId":        sql.NullString{Valid: false},
		"UID":            uid,
		"Secret":         secretKey,
		"Category":       "partner",
		"Scopes":         "read_product write_product read_order write_order",
		"Link":           "/app_store/plugins/test",
		"Confidential":   true,
		"Icon":           "\"oss/operation/ee5d8f1c0b641fa075539bd140593ba7.svg\"",
		"Embbed":         false,
		"PrivateApp":     false,
		"OauthDancable":  true,
		"Subscribable":   false,
		"StoreCharge":    false,
		"NotifyUri":      "",
		"ChargeMinPlan":  "",
		"AdditionalInfo": "",
		//"FreeStores": "",
		"WebhookApiVersion": "2020-07",
		"Regions":           []string{"1", "2"},
	})

	req := testhelper.NewJSONRequest(
		"GET",
		fmt.Sprintf("%s%s%d", suite.server.URL(), "/api/platform/apps/", app.ID+1),
		nil,
	)
	req.Header.Set("Origin", "platform")
	req.Header.Set("Platform-User-Id", "test")
	userPermissions, _ := json.Marshal(map[string]interface{}{
		"apps": map[string]interface{}{
			"read":  true,
			"write": false,
		},
	})
	req.Header.Set("Platform-User-Permissions", string(userPermissions))

	resp := testhelper.RunRequest(req)
	defer resp.Body.Close()

	Convey("get app info by id success", suite.T(), func() {
		Convey("should return 404", func() {
			So(resp.StatusCode, ShouldEqual, 404)
		})
	})
}

func (suite *BaseSuite) TestGetAppListSuccess() {
	uid, _ := models.GenerateAppKey(32)
	secretKey, _ := models.GenerateSecretKey(32)
	CreateApp(map[string]interface{}{
		"Name":           "test",
		"StoreId":        sql.NullString{Valid: false},
		"UID":            uid,
		"Secret":         secretKey,
		"Category":       "partner",
		"Scopes":         "read_product write_product read_order write_order",
		"Link":           "/app_store/plugins/test",
		"Confidential":   true,
		"Icon":           "\"oss/operation/ee5d8f1c0b641fa075539bd140593ba7.svg\"",
		"Embbed":         false,
		"PrivateApp":     false,
		"OauthDancable":  true,
		"Subscribable":   false,
		"StoreCharge":    false,
		"NotifyUri":      "",
		"ChargeMinPlan":  "",
		"AdditionalInfo": "",
		//"FreeStores": "",
		"WebhookApiVersion": "2020-07",
		"Regions":           []string{"1", "2"},
	})

	req := testhelper.NewJSONRequest(
		"GET",
		fmt.Sprintf("%s%s", suite.server.URL(), "/api/platform/apps"),
		nil,
	)
	req.Header.Set("Origin", "platform")
	req.Header.Set("Platform-User-Id", "test")
	userPermissions, _ := json.Marshal(map[string]interface{}{
		"apps": map[string]interface{}{
			"read":  true,
			"write": false,
		},
	})
	req.Header.Set("Platform-User-Permissions", string(userPermissions))

	resp := testhelper.RunRequest(req)
	defer resp.Body.Close()

	Convey("get app info by id success", suite.T(), func() {
		Convey("should return 200", func() {
			So(resp.StatusCode, ShouldEqual, 200)
		})
		Convey("should return correct app info", func() {
			js := testhelper.JSONResp(resp)
			So(js.Get("count").MustUint64(), ShouldEqual, 1)
		})
	})
}

func (suite *BaseSuite) TestUpdateAppByIdSuccess() {
	uid, _ := models.GenerateAppKey(32)
	secretKey, _ := models.GenerateSecretKey(32)
	app := CreateApp(map[string]interface{}{
		"Name":           "test",
		"StoreId":        sql.NullString{Valid: false},
		"UID":            uid,
		"Secret":         secretKey,
		"Category":       "partner",
		"Scopes":         "read_product write_product read_order write_order",
		"Link":           "/app_store/plugins/test",
		"Confidential":   true,
		"Icon":           "\"oss/operation/ee5d8f1c0b641fa075539bd140593ba7.svg\"",
		"Embbed":         false,
		"PrivateApp":     false,
		"OauthDancable":  true,
		"Subscribable":   false,
		"StoreCharge":    false,
		"NotifyUri":      "",
		"ChargeMinPlan":  "",
		"AdditionalInfo": "",
		//"FreeStores": "",
		"WebhookApiVersion": "2020-07",
		"Regions":           []string{"1", "2"},
	})

	scopes := [...]string{"read_product", "write_product", "read_order", "write_order", "read_script_tags", "write_script_tags"}
	req := testhelper.NewJSONRequest(
		"PATCH",
		fmt.Sprintf("%s%s%d", suite.server.URL(), "/api/platform/apps/", app.ID),
		gin.H{
			"name":                "size_chart1",
			"redirect_uri":        "https://size-chart1024.apps.shoplazza.com/callback/shoplazza/oauth1",
			"scopes":              scopes,
			"icon":                "oss/operation/d57d9bfb5eae56c09fdd1dca7d1736c3.svg",
			"email":               "test1@shoplazza.com",
			"category":            "operation",
			"app_uri":             "https://size-chart1024.apps.shoplazza.com/callback/shoplazza/oauth1",
			"webhook_api_version": "2020-01",
			"link":                "/app_store/plugins/test1",
			"embbed":              true,
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

	Convey("update app success by id", suite.T(), func() {
		Convey("should return 200", func() {
			So(resp.StatusCode, ShouldEqual, 200)
		})
		Convey("should return correct app info", func() {
			js := testhelper.JSONResp(resp)
			So(js.Get("id").MustUint64(), ShouldEqual, app.ID)
			So(js.Get("name").MustString(), ShouldEqual, "size_chart1")
			So(js.Get("redirect_uri").MustString(), ShouldEqual, "https://size-chart1024.apps.shoplazza.com/callback/shoplazza/oauth1")
			// So(js.Get("icon").MustString(), ShouldEqual, "oss/operation/d57d9bfb5eae56c09fdd1dca7d1736c3.svg")
			So(js.Get("email").MustString(), ShouldEqual, "test1@shoplazza.com")
			So(js.Get("category").MustString(), ShouldEqual, "operation")
			So(js.Get("app_uri").MustString(), ShouldEqual, "https://size-chart1024.apps.shoplazza.com/callback/shoplazza/oauth1")
			So(js.Get("webhook_api_version").MustString(), ShouldEqual, "2020-01")
			So(js.Get("link").MustString(), ShouldEqual, "/app_store/plugins/test1")
			So(js.Get("embbed").MustBool(), ShouldEqual, true)
		})
	})
}

func (suite *BaseSuite) TestPublishApp() {
	uid, _ := models.GenerateAppKey(32)
	secretKey, _ := models.GenerateSecretKey(32)
	app := CreateApp(map[string]interface{}{
		"ID":             uint64(115),
		"Name":           "tesssst",
		"StoreId":        sql.NullString{Valid: false},
		"UID":            uid,
		"Secret":         secretKey,
		"Category":       "partner",
		"Scopes":         "read_product write_product read_order write_order",
		"Link":           "/app_store/plugins/test",
		"Confidential":   true,
		"Icon":           "\"oss/operation/ee5d8f1c0b641fa075539bd140593ba7.svg\"",
		"Embbed":         false,
		"PrivateApp":     false,
		"OauthDancable":  true,
		"Subscribable":   false,
		"StoreCharge":    false,
		"NotifyUri":      "",
		"ChargeMinPlan":  "",
		"AdditionalInfo": "",
		//"FreeStores": "",
		"WebhookApiVersion": "2020-07",
		"Status":            "in-review",
		"Regions":           []string{"1", "2"},
	})

	req := testhelper.NewJSONRequest(
		"PUT",
		fmt.Sprintf("%s%s%d%s", suite.server.URL(), "/api/platform/apps/", app.ID, "/publish"),
		gin.H{
			"partner_locale": "zh-CN",
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

	Convey("publish app success", suite.T(), func() {
		Convey("should return 200", func() {
			So(resp.StatusCode, ShouldEqual, 200)
		})
	})
}

func (suite *BaseSuite) TestUnpublishApp() {
	uid, _ := models.GenerateAppKey(32)
	secretKey, _ := models.GenerateSecretKey(32)
	app := CreateApp(map[string]interface{}{
		"ID":             uint64(115),
		"Name":           "tesssst",
		"StoreId":        sql.NullString{Valid: false},
		"UID":            uid,
		"Secret":         secretKey,
		"Category":       "partner",
		"Scopes":         "read_product write_product read_order write_order",
		"Link":           "/app_store/plugins/test",
		"Confidential":   true,
		"Icon":           "\"oss/operation/ee5d8f1c0b641fa075539bd140593ba7.svg\"",
		"Embbed":         false,
		"PrivateApp":     false,
		"OauthDancable":  true,
		"Subscribable":   false,
		"StoreCharge":    false,
		"NotifyUri":      "",
		"ChargeMinPlan":  "",
		"AdditionalInfo": "",
		//"FreeStores": "",
		"WebhookApiVersion": "2020-07",
		"Status":            models.APPSTATUS_PUBLISHED,
		"Regions":           []string{"1", "2"},
	})
	CreateDevApp(map[string]interface{}{
		"ApplicationID":     uint64(115),
		"Name":              "tesssaasst",
		"RedirectUri":       "test",
		"Email":             "test",
		"AppUri":            "test",
		"Status":            models.DEVAPPSTATUS_DRAFT,
		"Icon":              "\"oss/operation/ee5d8f1c0b641fa075539bd140593ba7.svg\"",
		"WebhookApiVersion": "2020-07",
	})

	req := testhelper.NewJSONRequest(
		"PUT",
		fmt.Sprintf("%s%s%d%s", suite.server.URL(), "/api/platform/apps/", app.ID, "/unpublish"),
		gin.H{
			"locale":             "zh-CN",
			"level":              1,
			"send_email":         false,
			"unpublished_reason": "test",
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

	Convey("publish app success", suite.T(), func() {
		Convey("should return 200", func() {
			So(resp.StatusCode, ShouldEqual, 200)
		})
	})
}

func (suite *BaseSuite) TestRejectApp() {
	uid, _ := models.GenerateAppKey(32)
	secretKey, _ := models.GenerateSecretKey(32)
	app := CreateApp(map[string]interface{}{
		"ID":             uint64(115),
		"Name":           "tesssst",
		"StoreId":        sql.NullString{Valid: false},
		"UID":            uid,
		"Secret":         secretKey,
		"Category":       "partner",
		"Scopes":         "read_product write_product read_order write_order",
		"Link":           "/app_store/plugins/test",
		"Confidential":   true,
		"Icon":           "\"oss/operation/ee5d8f1c0b641fa075539bd140593ba7.svg\"",
		"Embbed":         false,
		"PrivateApp":     false,
		"OauthDancable":  true,
		"Subscribable":   false,
		"StoreCharge":    false,
		"NotifyUri":      "",
		"ChargeMinPlan":  "",
		"AdditionalInfo": "",
		//"FreeStores": "",
		"WebhookApiVersion": "2020-07",
		"Status":            "in-review",
		"Regions":           []string{"1", "2"},
	})

	req := testhelper.NewJSONRequest(
		"PUT",
		fmt.Sprintf("%s%s%d%s", suite.server.URL(), "/api/platform/apps/", app.ID, "/reject"),
		gin.H{
			"partner_locale":  "zh-CN",
			"rejected_reason": "test",
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

	Convey("publish app success", suite.T(), func() {
		Convey("should return 200", func() {
			So(resp.StatusCode, ShouldEqual, 200)
		})
	})
}

func (suite *BaseSuite) TestAcceptApp() {
	uid, _ := models.GenerateAppKey(32)
	secretKey, _ := models.GenerateSecretKey(32)
	app := CreateApp(map[string]interface{}{
		"ID":             uint64(118),
		"Name":           "tesssst",
		"StoreId":        sql.NullString{Valid: false},
		"UID":            uid,
		"Secret":         secretKey,
		"Category":       "partner",
		"Scopes":         "read_product write_product read_order write_order",
		"Link":           "/app_store/plugins/test",
		"Confidential":   true,
		"Icon":           "\"oss/operation/ee5d8f1c0b641fa075539bd140593ba7.svg\"",
		"Embbed":         false,
		"PrivateApp":     false,
		"OauthDancable":  true,
		"Subscribable":   false,
		"StoreCharge":    false,
		"NotifyUri":      "",
		"ChargeMinPlan":  "",
		"AdditionalInfo": "",
		//"FreeStores": "",
		"WebhookApiVersion":  "2020-07",
		"Status":             "submitted",
		"DevelopByShoplazza": false,
		"Regions":            []string{"1", "2"},
	})

	req := testhelper.NewJSONRequest(
		"PUT",
		fmt.Sprintf("%s%s%d%s", suite.server.URL(), "/api/platform/apps/", app.ID, "/accept"),
		gin.H{
			"partner_locale": "zh-CN",
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

	Convey("publish app success", suite.T(), func() {
		Convey("should return 200", func() {
			So(resp.StatusCode, ShouldEqual, 200)
		})
	})
}

func (suite *BaseSuite) TestReSubmitApp() {
	uid, _ := models.GenerateAppKey(32)
	secretKey, _ := models.GenerateSecretKey(32)
	app := CreateApp(map[string]interface{}{
		"ID":             uint64(118),
		"Name":           "tesssst",
		"StoreId":        sql.NullString{Valid: false},
		"UID":            uid,
		"Secret":         secretKey,
		"Category":       "partner",
		"Scopes":         "read_product write_product read_order write_order",
		"Link":           "/app_store/plugins/test",
		"Confidential":   true,
		"Icon":           "\"oss/operation/ee5d8f1c0b641fa075539bd140593ba7.svg\"",
		"Embbed":         false,
		"PrivateApp":     false,
		"OauthDancable":  true,
		"Subscribable":   false,
		"StoreCharge":    false,
		"NotifyUri":      "",
		"ChargeMinPlan":  "",
		"AdditionalInfo": "",
		//"FreeStores": "",
		"WebhookApiVersion": "2020-07",
		"Status":            "unpublished",
		"Regions":           []string{"1", "2"},
	})

	req := testhelper.NewJSONRequest(
		"PUT",
		fmt.Sprintf("%s%s%d%s", suite.server.URL(), "/api/platform/apps/", app.ID, "/resubmit"),
		gin.H{
			"partner_locale": "zh-CN",
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

	Convey("publish app success", suite.T(), func() {
		Convey("should return 200", func() {
			So(resp.StatusCode, ShouldEqual, 200)
		})
	})
}
