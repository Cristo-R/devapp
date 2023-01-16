package tag

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"gitlab.shoplazza.site/shoplaza/cobra/config"
	"gitlab.shoplazza.site/shoplaza/cobra/migrations"
	"gitlab.shoplazza.site/shoplaza/cobra/models"
	"gitlab.shoplazza.site/xiabing/goat.git/middlewares"
	"gitlab.shoplazza.site/xiabing/goat.git/testhelper"
)

func TestUpdateTagsForm_Do(t *testing.T) {
	if err := migrations.Migrate(config.DB); err != nil {
		panic(err)
	}
	gin.SetMode(gin.TestMode)

	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = new(http.Request)
	ctx.Request.Header = make(map[string][]string)
	ctx.Request.Header.Set("Origin", "platform")
	ctx.Request.Header.Set("Store-ID", strconv.FormatUint(1, 10))
	userPermissions, _ := json.Marshal(map[string]interface{}{
		"apps": map[string]interface{}{
			"read":  true,
			"write": true,
		},
	})
	ctx.Request.Header.Set("Platform-User-Permissions", string(userPermissions))
	middlewaresContext := new(middlewares.Context)
	middlewaresContext.Origin = "platform"
	ctx.Set("context", middlewaresContext)
	ctx.Request.Header.Set("Platform-User-Id", "789416f3-c73b-4fb6-bbd2-259da3309760")
	ctx.Request.Header.Set("User-Agent", "test")
	ctx.Request.Method = "PUT"
	ctx.Request.URL = &url.URL{}
	ctx.Request.URL.Path = "/test"

	app := &models.Application{}
	uid, _ := models.GenerateAppKey(32)
	secretKey, _ := models.GenerateSecretKey(32)
	testhelper.MakeRecord(config.DB, app, map[string]interface{}{
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
		"Status":            "published",
		"Regions":           []string{"1", "2"},
	})

	devApplicationLocale := &models.DevApplicationLocale{}
	testhelper.MakeRecord(config.DB, devApplicationLocale, map[string]interface{}{
		"Id":            uint64(879876),
		"Name":          "ok",
		"Subtitle":      "1",
		"ApplicationId": app.ID,
		"Locale":        "zh-CN",
		"Desc":          "desc",
		"IsPrimary":     true,
		"Icon":          "icon",
		"Status":        models.LOCALE_STATUS_SUBMITTED,
	})
	type fields struct {
		ID     uint64
		TagIds []uint64
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				ID:     app.ID,
				TagIds: []uint64{},
			},
			args: args{
				c: ctx,
			},
			want:    map[string]interface{}{"status": "success"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := &UpdateTagsForm{
				ID:     tt.fields.ID,
				TagIds: tt.fields.TagIds,
			}
			got, err := form.Do(tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateTagsForm.Do() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateTagsForm.Do() = %v, want %v", got, tt.want)
			}
		})
	}
}
