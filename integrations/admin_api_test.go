package integrations

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"bou.ke/monkey"
	"gopkg.in/h2non/gock.v1"

	"github.com/gin-gonic/gin"
	. "github.com/smartystreets/goconvey/convey"
	"gitlab.shoplazza.site/shoplaza/cobra/config"
	"gitlab.shoplazza.site/shoplaza/cobra/models"
	"gitlab.shoplazza.site/shoplaza/cobra/service"
	"gitlab.shoplazza.site/shoplaza/cobra/utils"
	"gitlab.shoplazza.site/xiabing/goat.git/testhelper"
)

func (suite *BaseSuite) TestGetInstalledApps() {
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
		"Regions":           []string{"1", "2"},
	})

	req := testhelper.NewJSONRequest(
		"GET",
		fmt.Sprintf("%s%s", suite.server.URL(), "/api/admin/installed_apps"),
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

	Convey("get installed apps", suite.T(), func() {
		Convey("should return 200", func() {
			So(resp.StatusCode, ShouldEqual, 200)
		})
	})
}

func (suite *BaseSuite) TestGetSessionToken() {
	defer gock.Off()

	gock.EnableNetworking()
	defer gock.DisableNetworking()
	u, _ := url.Parse(config.Cfg.StoreService)
	gock.NetworkingFilter(func(req *http.Request) bool {
		return !strings.Contains(req.URL.Host, u.Host)
	})

	mockUserInfoApiReturnOk()

	type UserInfoResp struct {
		Contact string `json:"contact"`
		Locale  string `json:"locale"`
	}

	tests := []struct {
		Params           gin.H
		ExpectStatusCode int
		ExpectReturn     *UserInfoResp
		Desc             string
	}{
		{
			Params:           nil,
			ExpectStatusCode: http.StatusOK,
			ExpectReturn: &UserInfoResp{
				Contact: "test.preview.shoplazza.com",
				Locale:  "zh-CN",
			},
			Desc: "call totoro service return user info",
		},
	}

	for _, tt := range tests {
		Convey(tt.Desc, suite.T(), func() {
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
				"Regions":           []string{"1", "2"},
			})

			req := testhelper.NewJSONRequest(
				"GET",
				fmt.Sprintf("%s%s", suite.server.URL(), fmt.Sprintf("/api/admin/session_token?%s", uid)),
				nil,
			)

			req.Header.Set("System-Domain", "test.preview.shoplazza.com")
			req.Header.Set("Login-User-Id", "test")
			req.Header.Set("Store-ID", strconv.FormatUint(1, 10))
			cookie := http.Cookie{
				Name:     "awesomev2",
				Value:    "MTY0MDIyMzE5MHxRaHMzanN1OF9leGdWQTNYZmdqS2tvcnQ0UXpmVlhrZVlhZlJSSG1URTBnOUY4WFNVdl9BVWVmNHozbkVnYU5yc3NwRG9MZFptSGs9fPCmLb7qbttCuZl79rEcRKho9lRqTLZsvs_OESW0um8I",
				HttpOnly: true,
				Path:     "/",
			}
			req.AddCookie(&cookie)
			resp := testhelper.RunRequest(req)
			So(resp.StatusCode, ShouldEqual, 200)
		})
	}
}

func mockUserInfoApiReturnOk() {
	gock.New(config.Cfg.StoreService).
		Reply(http.StatusOK).
		JSON(map[string]interface{}{
			"locale":  "zh-CN",
			"contact": "test.preview.shoplazza.com",
		})
}

func (suite *BaseSuite) TestGetAppInfo() {
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
		"IsSmartApp":        true,
		"Regions":           []string{"1", "2"},
	})

	req := testhelper.NewJSONRequest(
		"GET",
		fmt.Sprintf("%s%s", suite.server.URL(), "/api/admin/app/115"),
		nil,
	)

	resp := testhelper.RunRequest(req)
	defer resp.Body.Close()

	Convey("get app info", suite.T(), func() {
		Convey("should return 200", func() {
			So(resp.StatusCode, ShouldEqual, 200)
		})
	})
}

func (suite *BaseSuite) TestCheckOauth() {
	uid, _ := models.GenerateAppKey(32)
	secretKey, _ := models.GenerateSecretKey(32)
	CreateApp(map[string]interface{}{
		"ID":                uint64(116),
		"Name":              "tesssst",
		"StoreId":           sql.NullString{Valid: false},
		"UID":               uid,
		"Secret":            secretKey,
		"Category":          "partner",
		"Scopes":            "read_product write_product read_order write_order",
		"Link":              "/app_store/plugins/test",
		"Confidential":      true,
		"Icon":              "\"oss/operation/ee5d8f1c0b641fa075539bd140593ba7.svg\"",
		"Embbed":            true,
		"PrivateApp":        false,
		"OauthDancable":     true,
		"Subscribable":      false,
		"StoreCharge":       false,
		"NotifyUri":         "",
		"ChargeMinPlan":     "",
		"AppUri":            "https://app.com",
		"WebhookApiVersion": "2020-07",
		"IsSmartApp":        true,
		"Regions":           []string{"1", "2"},
	})

	guides := models.Guides{}
	guides = append(guides, models.Guide{Name: "test", Status: "init"})
	store := CreateStore(map[string]interface{}{
		"Id":       utils.NewUUIDBinary(),
		"OriginId": "7295",
		"Name":     "pjs",
		"Guides":   &guides,
	})

	CreateInstallTrack(map[string]interface{}{
		"Id":            utils.NewUUIDBinary(),
		"StoreId":       store.Id,
		"ApplicationId": uint64(116),
		"InstalledAt":   time.Now(),
	})

	CreateOauthAccessToken(map[string]interface{}{
		"Id":              uint64(115),
		"ResourceOwnerId": utils.NewUUIDBinary(),
		"ApplicationId":   uint64(116),
		"Token":           "test",
		"RefreshToken":    "test",
		"ExpiresIn":       123,
		"CreatedAt":       time.Now().AddDate(0, 0, 1),
	})

	req := testhelper.NewJSONRequest(
		"GET",
		fmt.Sprintf("%s%s", suite.server.URL(), fmt.Sprintf("/api/admin/app/116/check_oauth")),
		nil,
	)

	req.Header.Set("Store-Domain", "test.preview.shoplazza.com")
	req.Header.Set("Slug", "test")
	req.Header.Set("Store-ID", "7295")
	cookie := http.Cookie{
		Name:     "awesomev2",
		Value:    "MTY0MDIyMzE5MHxRaHMzanN1OF9leGdWQTNYZmdqS2tvcnQ0UXpmVlhrZVlhZlJSSG1URTBnOUY4WFNVdl9BVWVmNHozbkVnYU5yc3NwRG9MZFptSGs9fPCmLb7qbttCuZl79rEcRKho9lRqTLZsvs_OESW0um8I",
		HttpOnly: true,
		Path:     "/",
	}
	req.AddCookie(&cookie)
	resp := testhelper.RunRequest(req)
	defer resp.Body.Close()

	Convey("check oauth success", suite.T(), func() {
		Convey("should return 200", func() {
			So(resp.StatusCode, ShouldEqual, 200)
		})
	})

}

func (suite *BaseSuite) TestAuthrizedEmbbedApp() {
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
		"Embbed":        true,
		"PrivateApp":    false,
		"OauthDancable": true,
		"Subscribable":  false,
		"StoreCharge":   false,
		"AppUri":        "https://app.com",
		"NotifyUri":     "",
		"ChargeMinPlan": "",
		//"FreeStores": "",
		"WebhookApiVersion": "2020-07",
		"Regions":           []string{"1", "2"},
	})

	req := testhelper.NewJSONRequest(
		"GET",
		fmt.Sprintf("%s%s", suite.server.URL(), fmt.Sprintf("/authorized_embbed_apps/%s", uid)),
		nil,
	)

	req.Header.Set("Store-Domain", "test.preview.shoplazza.com")
	req.Header.Set("Slug", "test")
	req.Header.Set("Store-ID", "7295")

	resp := testhelper.RunRequest(req)
	defer resp.Body.Close()

	Convey("authorize embbed app", suite.T(), func() {
		Convey("should return 302", func() {
			So(resp.StatusCode, ShouldEqual, 302)
		})
	})

}

func (suite *BaseSuite) TestGetGuides() {
	monkey.Patch(service.GetStoreFromTotoro, func(storeId string) (*service.StoreInfo, error) {
		store := &service.StoreInfo{
			Id:         7295,
			Name:       "pjs",
			UserLocale: "zh-CN",
			CreatedAt:  time.Now(),
		}
		return store, nil
	})

	guides := models.Guides{}
	guides = append(guides, models.Guide{Name: "test", Status: "init"})
	CreateStore(map[string]interface{}{
		"Id":       utils.NewUUIDBinary(),
		"OriginId": "7295",
		"Name":     "pjs",
		"Guides":   &guides,
	})

	req := testhelper.NewJSONRequest(
		"GET",
		fmt.Sprintf("%s%s", suite.server.URL(), fmt.Sprintf("/api/admin/app/guides")),
		nil,
	)

	req.Header.Set("Store-Domain", "test.preview.shoplazza.com")
	req.Header.Set("Slug", "test")
	req.Header.Set("Store-ID", "7295")
	cookie := http.Cookie{
		Name:     "awesomev2",
		Value:    "MTY0MDIyMzE5MHxRaHMzanN1OF9leGdWQTNYZmdqS2tvcnQ0UXpmVlhrZVlhZlJSSG1URTBnOUY4WFNVdl9BVWVmNHozbkVnYU5yc3NwRG9MZFptSGs9fPCmLb7qbttCuZl79rEcRKho9lRqTLZsvs_OESW0um8I",
		HttpOnly: true,
		Path:     "/",
	}
	req.AddCookie(&cookie)
	resp := testhelper.RunRequest(req)
	defer resp.Body.Close()

	Convey("get guides success", suite.T(), func() {
		Convey("should return 200", func() {
			So(resp.StatusCode, ShouldEqual, 200)
		})
	})

}

func (suite *BaseSuite) TestUpdateGuides() {
	guides := models.Guides{}
	guides = append(guides, models.Guide{Name: "test", Status: "init"})
	CreateStore(map[string]interface{}{
		"Id":       utils.NewUUIDBinary(),
		"OriginId": "7295",
		"Name":     "pjs",
		"Guides":   &guides,
	})

	req := testhelper.NewJSONRequest(
		"PUT",
		fmt.Sprintf("%s%s", suite.server.URL(), fmt.Sprintf("/api/admin/app/guides/close")),
		gin.H{
			"key": "test",
		},
	)

	req.Header.Set("Store-Domain", "test.preview.shoplazza.com")
	req.Header.Set("Slug", "test")
	req.Header.Set("Store-ID", "7295")
	cookie := http.Cookie{
		Name:     "awesomev2",
		Value:    "MTY0MDIyMzE5MHxRaHMzanN1OF9leGdWQTNYZmdqS2tvcnQ0UXpmVlhrZVlhZlJSSG1URTBnOUY4WFNVdl9BVWVmNHozbkVnYU5yc3NwRG9MZFptSGs9fPCmLb7qbttCuZl79rEcRKho9lRqTLZsvs_OESW0um8I",
		HttpOnly: true,
		Path:     "/",
	}
	req.AddCookie(&cookie)
	resp := testhelper.RunRequest(req)
	defer resp.Body.Close()

	Convey("update guides success", suite.T(), func() {
		Convey("should return 200", func() {
			So(resp.StatusCode, ShouldEqual, 200)
		})
	})

}
