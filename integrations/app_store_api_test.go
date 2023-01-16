package integrations

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"

	"bou.ke/monkey"

	"github.com/gin-gonic/gin"
	. "github.com/smartystreets/goconvey/convey"
	"gitlab.shoplazza.site/shoplaza/cobra/models"
	"gitlab.shoplazza.site/shoplaza/cobra/service"
	"gitlab.shoplazza.site/shoplaza/cobra/utils"
	"gitlab.shoplazza.site/xiabing/goat.git/testhelper"
)

func (suite *BaseSuite) TestGetTags() {
	CreateTag(map[string]interface{}{
		"Id":       uint64(115),
		"NameZh":   "测试数据",
		"NameEn":   "Test Data",
		"ParentId": uint64(0),
	})

	req := testhelper.NewJSONRequest(
		"GET",
		fmt.Sprintf("%s%s", suite.server.URL(), "/api/app_store/tags"),
		nil,
	)

	resp := testhelper.RunRequest(req)
	defer resp.Body.Close()

	Convey("get tags", suite.T(), func() {
		Convey("should return 200", func() {
			So(resp.StatusCode, ShouldEqual, 200)
		})
	})
}

func (suite *BaseSuite) TestGetCollections() {
	CreateCollectionLocale(map[string]interface{}{
		"ID":           uint64(100001),
		"Name":         "测试数据",
		"Description":  "测试数据",
		"Locale":       "zh-CN",
		"CollectionId": uint64(10000),
	})
	CreatePage(map[string]interface{}{
		"ID":     uint64(100001),
		"Title":  "test",
		"Handle": "test",
	})
	CreatePageCollect(map[string]interface{}{
		"ID":                 uint64(100001),
		"CollectionId":       uint64(10000),
		"PageId":             uint64(100001),
		"CollectionPosition": 1,
	})

	req := testhelper.NewJSONRequest(
		"GET",
		fmt.Sprintf("%s%s", suite.server.URL(), "/api/app_store/collections?handle=test"),
		nil,
	)

	resp := testhelper.RunRequest(req)
	defer resp.Body.Close()

	Convey("get collections", suite.T(), func() {
		Convey("should return 200", func() {
			So(resp.StatusCode, ShouldEqual, 200)
		})
	})
}

func (suite *BaseSuite) TestGetAppsForAppStore() {
	uid, _ := models.GenerateAppKey(32)
	secretKey, _ := models.GenerateSecretKey(32)
	app := CreateApp(map[string]interface{}{
		"ID":                uint64(1),
		"Name":              "test",
		"StoreId":           sql.NullString{Valid: false},
		"UID":               uid,
		"Secret":            secretKey,
		"Category":          "partner",
		"Scopes":            "read_product write_product read_order write_order",
		"Link":              "/app_store/plugins/test",
		"Confidential":      true,
		"Icon":              "\"oss/operation/ee5d8f1c0b641fa075539bd140593ba7.svg\"",
		"Embbed":            false,
		"PrivateApp":        false,
		"OauthDancable":     true,
		"Subscribable":      false,
		"StoreCharge":       false,
		"NotifyUri":         "",
		"ChargeMinPlan":     "",
		"WebhookApiVersion": "2020-07",
		"Status":            "published",
		"Listing":           true,
		"Regions":           []string{"1", "2"},
	})
	CreateLocal(map[string]interface{}{
		"Id":            utils.NewUUIDBinary(),
		"Locale":        "zh-CN",
		"ApplicationId": app.ID,
	})
	tag := CreateTag(map[string]interface{}{
		"Id":       uint64(115),
		"NameZh":   "测试数据",
		"NameEn":   "Test Data",
		"ParentId": uint64(0),
	})
	CreateAppTag(map[string]interface{}{
		"Id":            uint64(115),
		"ApplicationId": app.ID,
		"TagId":         tag.Id,
	})

	req := testhelper.NewJSONRequest(
		"GET",
		fmt.Sprintf("%s%s", suite.server.URL(), "/api/app_store/apps"),
		nil,
	)

	resp := testhelper.RunRequest(req)
	defer resp.Body.Close()

	Convey("get get apps for app store", suite.T(), func() {
		Convey("should return 200", func() {
			So(resp.StatusCode, ShouldEqual, 200)
		})
	})
}

func (suite *BaseSuite) TestGetAppsByCollectionId() {
	uid, _ := models.GenerateAppKey(32)
	secretKey, _ := models.GenerateSecretKey(32)
	app := CreateApp(map[string]interface{}{
		"ID":                uint64(1),
		"Name":              "test",
		"StoreId":           sql.NullString{Valid: false},
		"UID":               uid,
		"Secret":            secretKey,
		"Category":          "partner",
		"Scopes":            "read_product write_product read_order write_order",
		"Link":              "/app_store/plugins/test",
		"Confidential":      true,
		"Icon":              "\"oss/operation/ee5d8f1c0b641fa075539bd140593ba7.svg\"",
		"Embbed":            false,
		"PrivateApp":        false,
		"OauthDancable":     true,
		"Subscribable":      false,
		"StoreCharge":       false,
		"NotifyUri":         "",
		"ChargeMinPlan":     "",
		"WebhookApiVersion": "2020-07",
		"Status":            "published",
		"Listing":           true,
		"Regions":           []string{"1", "2"},
	})
	CreateLocal(map[string]interface{}{
		"Id":            utils.NewUUIDBinary(),
		"Locale":        "en-US",
		"ApplicationId": app.ID,
	})
	CreateCollection(map[string]interface{}{
		"ID":           uint64(10000),
		"IsShowAtNgv":  false,
		"IsShowAtHome": true,
		"Position":     1,
	})
	CreateCollectionLocale(map[string]interface{}{
		"ID":           uint64(100001),
		"Name":         "测试数据",
		"Description":  "测试数据",
		"Locale":       "en-US",
		"CollectionId": uint64(10000),
	})
	CreateApplicationCollect(map[string]interface{}{
		"ID":            uint64(1),
		"CollectionId":  uint64(10000),
		"ApplicationId": app.ID,
	})
	tag := CreateTag(map[string]interface{}{
		"Id":       uint64(115),
		"NameZh":   "测试数据",
		"NameEn":   "Test Data",
		"ParentId": uint64(0),
	})
	CreateAppTag(map[string]interface{}{
		"Id":            uint64(115),
		"ApplicationId": app.ID,
		"TagId":         tag.Id,
	})

	req := testhelper.NewJSONRequest(
		"GET",
		fmt.Sprintf("%s%s", suite.server.URL(), "/api/app_store/collections/10000"),
		nil,
	)

	resp := testhelper.RunRequest(req)
	defer resp.Body.Close()
	fmt.Println(resp.Body)
	Convey("get get apps by collection id", suite.T(), func() {
		Convey("should return 200", func() {
			So(resp.StatusCode, ShouldEqual, 200)
		})
	})
}

func (suite *BaseSuite) TestGetAppCountGroupByTag() {
	uid, _ := models.GenerateAppKey(32)
	secretKey, _ := models.GenerateSecretKey(32)
	app := CreateApp(map[string]interface{}{
		"ID":                uint64(1),
		"Name":              "test",
		"StoreId":           sql.NullString{Valid: false},
		"UID":               uid,
		"Secret":            secretKey,
		"Category":          "partner",
		"Scopes":            "read_product write_product read_order write_order",
		"Link":              "/app_store/plugins/test",
		"Confidential":      true,
		"Icon":              "\"oss/operation/ee5d8f1c0b641fa075539bd140593ba7.svg\"",
		"Embbed":            false,
		"PrivateApp":        false,
		"OauthDancable":     true,
		"Subscribable":      false,
		"StoreCharge":       false,
		"NotifyUri":         "",
		"ChargeMinPlan":     "",
		"WebhookApiVersion": "2020-07",
		"Status":            "published",
		"Listing":           true,
		"Regions":           []string{"1", "2"},
	})
	CreateLocal(map[string]interface{}{
		"Id":            utils.NewUUIDBinary(),
		"Locale":        "zh-CN",
		"ApplicationId": app.ID,
	})
	tag := CreateTag(map[string]interface{}{
		"Id":       uint64(115),
		"NameZh":   "测试数据",
		"NameEn":   "Test Data",
		"ParentId": uint64(0),
	})
	CreateAppTag(map[string]interface{}{
		"Id":            uint64(115),
		"ApplicationId": app.ID,
		"TagId":         tag.Id,
	})
	monkey.Patch(service.GetShoplazzaAccount, func(storeId string) (*service.User, error) {
		user := &service.User{
			Region: "cn",
		}
		return user, nil
	})
	monkey.Patch(service.GetStoreOriginUserID, func(storeId uint64) (string, error) {

		return "", nil
	})
	req := testhelper.NewJSONRequest(
		"GET",
		fmt.Sprintf("%s%s", suite.server.URL(), "/api/app_store/tags/apps/count"),
		nil,
	)

	req.Header.Add("Login-Store-ID", "0")
	resp := testhelper.RunRequest(req)
	defer resp.Body.Close()

	Convey("get get app count group by tag", suite.T(), func() {
		Convey("should return 200", func() {
			So(resp.StatusCode, ShouldEqual, 200)
		})
	})
}

func (suite *BaseSuite) TestGetAppByIdForAppStore() {
	uid, _ := models.GenerateAppKey(32)
	secretKey, _ := models.GenerateSecretKey(32)
	app := CreateApp(map[string]interface{}{
		"ID":                uint64(1),
		"Name":              "test",
		"StoreId":           sql.NullString{Valid: false},
		"UID":               uid,
		"Secret":            secretKey,
		"Category":          "partner",
		"Scopes":            "read_product write_product read_order write_order",
		"Link":              "/app_store/plugins/test",
		"Confidential":      true,
		"Icon":              "\"oss/operation/ee5d8f1c0b641fa075539bd140593ba7.svg\"",
		"Embbed":            false,
		"PrivateApp":        false,
		"OauthDancable":     true,
		"Subscribable":      false,
		"StoreCharge":       false,
		"NotifyUri":         "",
		"ChargeMinPlan":     "",
		"WebhookApiVersion": "2020-07",
		"Status":            "published",
		"Listing":           true,
		"Regions":           []string{"1", "2"},
	})
	CreateLocal(map[string]interface{}{
		"Id":            utils.NewUUIDBinary(),
		"Locale":        "en-US",
		"ApplicationId": app.ID,
	})
	tag := CreateTag(map[string]interface{}{
		"Id":       uint64(115),
		"NameZh":   "测试数据",
		"NameEn":   "Test Data",
		"ParentId": uint64(0),
	})
	CreateAppTag(map[string]interface{}{
		"Id":            uint64(115),
		"ApplicationId": app.ID,
		"TagId":         tag.Id,
	})

	req := testhelper.NewJSONRequest(
		"GET",
		fmt.Sprintf("%s%s", suite.server.URL(), "/api/app_store/apps/1"),
		nil,
	)

	resp := testhelper.RunRequest(req)
	defer resp.Body.Close()

	Convey("get app by id", suite.T(), func() {
		Convey("should return 200", func() {
			So(resp.StatusCode, ShouldEqual, 200)
		})
	})
}

func (suite *BaseSuite) TestGetAccountStores() {
	var guard *monkey.PatchGuard
	guard = monkey.PatchInstanceMethod(reflect.TypeOf(utils.NewHttpPool()), "HttpRequest", func(httpPool *utils.HttpPool, method string, reqUrl string, body interface{}) ([]byte, error) {
		guard.Unpatch()
		defer guard.Restore()

		type Store struct {
			ID            uint64 `json:"id"`
			Name          string `json:"name"`
			Slug          string `json:"slug"`
			PrimaryDomain string `json:"primary_domain"`
		}

		stores := []Store{
			{
				ID:            1,
				Name:          "test",
				Slug:          "test",
				PrimaryDomain: "test.myshoplaza.com",
			},
		}

		res, _ := json.Marshal(stores)
		return res, nil
	})

	req := testhelper.NewJSONRequest(
		"GET",
		fmt.Sprintf("%s%s", suite.server.URL(), "/api/account/stores"),
		nil,
	)
	req.Header.Set("Login-User-Id", "dafd283d-1274-4412-b86d-21a68ab1172f")

	resp := testhelper.RunRequest(req)
	defer resp.Body.Close()

	defer Convey("get account stores", suite.T(), func() {
		Convey("should return 200", func() {
			So(resp.StatusCode, ShouldEqual, 200)
		})
	})
}

func (suite *BaseSuite) TestCreateAppSupportFeedback() {
	CreateAppSupport(map[string]interface{}{
		"ID":                     uint64(1),
		"StoreId":                uint64(7295),
		"ApplicationId":          uint64(7295),
		"ApplicationListingName": "test",
		"ApplicationListingIcon": "test",
		"IssueDescription":       "test",
		"ContactEmail":           "pengjunjie@shoplazza.com",
		"StoreName":              "test",
		"StoreDomain":            "test",
		"StoreLocale":            "zh-CN",
		"HasSendFeedbackEmail":   true,
		"FeedbackAccessToken":    "test",
	})

	req := testhelper.NewJSONRequest(
		"POST",
		fmt.Sprintf("%s%s", suite.server.URL(), "/api/app_store/apps/support_demand/feedback"),
		gin.H{
			"feedback_access_token": "test",
			"reply_by_developer":    true,
			"rating":                2,
		},
	)

	resp := testhelper.RunRequest(req)
	defer resp.Body.Close()

	Convey("create app support feedback success", suite.T(), func() {
		Convey("should return 200", func() {
			So(resp.StatusCode, ShouldEqual, 200)
		})
	})
}

func (suite *BaseSuite) TestGetAppSupportFeedbackStatus() {
	CreateAppSupport(map[string]interface{}{
		"ID":                     uint64(1),
		"StoreId":                uint64(7295),
		"ApplicationId":          uint64(7295),
		"ApplicationListingName": "test",
		"ApplicationListingIcon": "test",
		"IssueDescription":       "test",
		"ContactEmail":           "pengjunjie@shoplazza.com",
		"StoreName":              "test",
		"StoreDomain":            "test",
		"StoreLocale":            "zh-CN",
		"HasSendFeedbackEmail":   true,
		"FeedbackAccessToken":    "test",
	})

	req := testhelper.NewJSONRequest(
		"GET",
		fmt.Sprintf("%s%s", suite.server.URL(), "/api/app_store/apps/support_demand/feedback_guard_check?feedback_access_token=test"),
		nil,
	)

	resp := testhelper.RunRequest(req)
	defer resp.Body.Close()

	Convey("get app support feedback status", suite.T(), func() {
		Convey("should return 200", func() {
			So(resp.StatusCode, ShouldEqual, 200)
		})
	})
}

func (suite *BaseSuite) TestCreateAppSupport() {
	monkey.PatchInstanceMethod(reflect.TypeOf(utils.NewHttpPool()), "HttpRequest", func(httpPool *utils.HttpPool, method string, reqUrl string, body interface{}) ([]byte, error) {
		type Store struct {
			ID            uint64 `json:"id"`
			Name          string `json:"name"`
			Slug          string `json:"slug"`
			PrimaryDomain string `json:"primary_domain"`
		}

		stores := []Store{
			{
				ID:            123,
				Name:          "test",
				Slug:          "test",
				PrimaryDomain: "test.myshoplaza.com",
			},
		}

		res, _ := json.Marshal(stores)
		return res, nil
	})

	monkey.Patch(service.GetStoreFromTotoro, func(storeId string) (*service.StoreInfo, error) {
		store := &service.StoreInfo{
			Id:         123,
			Name:       "asd",
			UserLocale: "zh-CN",
		}
		return store, nil
	})

	monkey.Patch(service.GetShoplazzaAccount, func(userId string) (*service.User, error) {
		user := &service.User{
			Contact: "pengjunjie@shoplazza.com",
			Locale:  "zh-CN",
		}
		return user, nil
	})

	uid, _ := models.GenerateAppKey(32)
	secretKey, _ := models.GenerateSecretKey(32)
	CreateApp(map[string]interface{}{
		"ID":                uint64(1),
		"Name":              "test",
		"StoreId":           sql.NullString{Valid: false},
		"UID":               uid,
		"Secret":            secretKey,
		"Category":          "partner",
		"Scopes":            "read_product write_product read_order write_order",
		"Link":              "/app_store/plugins/test",
		"Confidential":      true,
		"Icon":              "\"oss/operation/ee5d8f1c0b641fa075539bd140593ba7.svg\"",
		"Embbed":            false,
		"PrivateApp":        false,
		"OauthDancable":     true,
		"Subscribable":      false,
		"StoreCharge":       false,
		"NotifyUri":         "",
		"ChargeMinPlan":     "",
		"WebhookApiVersion": "2020-07",
		"Status":            "published",
		"Listing":           true,
		"Regions":           []string{"1", "2"},
	})

	CreateAppSupport(map[string]interface{}{
		"ID":                     uint64(1),
		"StoreId":                uint64(7295),
		"ApplicationId":          uint64(7295),
		"ApplicationListingName": "test",
		"ApplicationListingIcon": "test",
		"IssueDescription":       "test",
		"ContactEmail":           "pengjunjie@shoplazza.com",
		"StoreName":              "test",
		"StoreDomain":            "test",
		"StoreLocale":            "zh-CN",
		"HasSendFeedbackEmail":   true,
		"FeedbackAccessToken":    "test",
	})

	req := testhelper.NewJSONRequest(
		"POST",
		fmt.Sprintf("%s%s", suite.server.URL(), "/api/app_store/apps/support_demand"),
		gin.H{
			"app_id":            uint64(1),
			"issue_description": "test",
			"contact_email":     "pengjunjie@shoplazza.com",
			"submit_source":     "app_store_app_listing",
		},
	)

	req.Header.Set("Login-Store-ID", "123")
	req.Header.Set("Login-Store-Domain", "pengjunjie233.preview.shoplazza.com")

	resp := testhelper.RunRequest(req)
	defer resp.Body.Close()

	Convey("create app support success", suite.T(), func() {
		Convey("should return 200", func() {
			So(resp.StatusCode, ShouldEqual, 200)
		})
	})
}

func (suite *BaseSuite) TestGetAppInstalledStatus() {
	monkey.PatchInstanceMethod(reflect.TypeOf(utils.NewHttpPool()), "HttpRequest", func(httpPool *utils.HttpPool, method string, reqUrl string, body interface{}) ([]byte, error) {
		type Store struct {
			ID            uint64 `json:"id"`
			Name          string `json:"name"`
			Slug          string `json:"slug"`
			PrimaryDomain string `json:"primary_domain"`
		}

		stores := []Store{
			{
				ID:            123,
				Name:          "test",
				Slug:          "test",
				PrimaryDomain: "test.myshoplaza.com",
			},
		}

		res, _ := json.Marshal(stores)
		return res, nil
	})

	uid, _ := models.GenerateAppKey(32)
	secretKey, _ := models.GenerateSecretKey(32)
	CreateApp(map[string]interface{}{
		"ID":                uint64(1),
		"Name":              "test",
		"StoreId":           sql.NullString{Valid: false},
		"UID":               uid,
		"Secret":            secretKey,
		"Category":          "partner",
		"Scopes":            "read_product write_product read_order write_order",
		"Link":              "/app_store/plugins/test",
		"Confidential":      true,
		"Icon":              "\"oss/operation/ee5d8f1c0b641fa075539bd140593ba7.svg\"",
		"Embbed":            false,
		"PrivateApp":        false,
		"OauthDancable":     true,
		"Subscribable":      false,
		"StoreCharge":       false,
		"NotifyUri":         "",
		"ChargeMinPlan":     "",
		"WebhookApiVersion": "2020-07",
		"Status":            "published",
		"Listing":           true,
		"Regions":           []string{"1", "2"},
	})

	req := testhelper.NewJSONRequest(
		"GET",
		fmt.Sprintf("%s%s", suite.server.URL(), "/api/app_store/current_user/apps/1/installed_status"),
		nil,
	)

	req.Header.Set("Login-Store-ID", "123")
	req.Header.Set("Login-User-Id", "test")

	resp := testhelper.RunRequest(req)
	defer resp.Body.Close()

	Convey("get app installed status", suite.T(), func() {
		Convey("should return 200", func() {
			So(resp.StatusCode, ShouldEqual, 200)
		})
	})
}

func (suite *BaseSuite) TestGetStoreInfo() {
	monkey.PatchInstanceMethod(reflect.TypeOf(utils.NewHttpPool()), "HttpRequest", func(httpPool *utils.HttpPool, method string, reqUrl string, body interface{}) ([]byte, error) {
		type Store struct {
			ID            uint64 `json:"id"`
			Name          string `json:"name"`
			Slug          string `json:"slug"`
			PrimaryDomain string `json:"primary_domain"`
		}

		stores := []Store{
			{
				ID:            123,
				Name:          "test",
				Slug:          "test",
				PrimaryDomain: "test.myshoplaza.com",
			},
		}

		res, _ := json.Marshal(stores)
		return res, nil
	})

	monkey.Patch(service.GetStore, func(storeId string) (*service.Store, error) {
		store := &service.Store{
			Name:   "asd",
			Email:  "pengjunjie@shoplazza.com",
			Locale: "zh-CN",
		}
		return store, nil
	})

	req := testhelper.NewJSONRequest(
		"GET",
		fmt.Sprintf("%s%s", suite.server.URL(), "/api/account/stores/123"),
		nil,
	)

	resp := testhelper.RunRequest(req)
	defer resp.Body.Close()

	Convey("get app account store success", suite.T(), func() {
		Convey("should return 200", func() {
			So(resp.StatusCode, ShouldEqual, 200)
		})
	})
}
