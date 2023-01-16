package integrations

import (
	"encoding/json"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"gitlab.shoplazza.site/xiabing/goat.git/testhelper"
)

func (suite *BaseSuite) TestGetAppSupports() {
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
		fmt.Sprintf("%s%s", suite.server.URL(), "/api/platform/apps/1/support_demands"),
		nil,
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

	Convey("get app supports", suite.T(), func() {
		Convey("should return 200", func() {
			So(resp.StatusCode, ShouldEqual, 200)
		})
	})
}
