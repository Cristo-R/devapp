package tag

import (
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
)

func TestGetTagsForm_Do(t *testing.T) {
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

	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name    string
		form    *GetTagsForm
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "ok",
			form: &GetTagsForm{},
			args: args{
				c: ctx,
			},
			want:    map[string]interface{}{"tags": []models.TagView{}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := &GetTagsForm{}
			got, err := form.Do(tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTagsForm.Do() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTagsForm.Do() = %v, want %v", got, tt.want)
			}
		})
	}
}
