package integrations

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gitlab.shoplazza.site/xiabing/goat.git/i18n"
	mw "gitlab.shoplazza.site/xiabing/goat.git/middlewares"
	"gitlab.shoplazza.site/xiabing/goat.git/testhelper"

	"gitlab.shoplazza.site/shoplaza/cobra/config"
	"gitlab.shoplazza.site/shoplaza/cobra/forms"
	app "gitlab.shoplazza.site/shoplaza/cobra/forms/internal_app"
	app_locale "gitlab.shoplazza.site/shoplaza/cobra/forms/internal_app_locale"
	app_locale_c2 "gitlab.shoplazza.site/shoplaza/cobra/forms/internal_app_locale_v2"
	app_v2 "gitlab.shoplazza.site/shoplaza/cobra/forms/internal_app_v2"
	"gitlab.shoplazza.site/shoplaza/cobra/middlewares"
	"gitlab.shoplazza.site/shoplaza/cobra/migrations"
	"gitlab.shoplazza.site/shoplaza/cobra/models"
	"gitlab.shoplazza.site/shoplaza/cobra/service"
	"gitlab.shoplazza.site/shoplaza/cobra/utils"
)

func resetAppTime(app *models.Application) *models.Application {
	nowTime := time.Time{}
	app.UpdatedAt = nowTime
	app.CreatedAt = nowTime
	return app
}

func TestCreateAndSubmitApp(t *testing.T) {
	if err := migrations.Migrate(config.DB); err != nil {
		panic(err)
	}
	assert := assert.New(t)
	err := i18n.Bundle(config.Cfg.BasePath+"/en.json", config.Cfg.BasePath+"/zh.json")
	assert.Nil(err)
	db := config.DB.LogMode(false)
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Set("context", &mw.Context{
		Locale: "en-US",
	})

	var defaultApp = &models.Application{}

	defer func() {
		deleteApp, err := service.GetAppFromGinContext(ctx)
		if err != nil {
			panic(err.Error())
		}

		db.Delete(&models.Application{}).Where("id = ?", deleteApp.ID)
		db.Delete(&models.AppLocal{}).Where("application_id = ?", deleteApp.ID)
		db.Delete(&models.DevApplicationLocale{}).Where("application_id = ?", deleteApp.ID)
		db.Delete(&models.ApplicationExtend{}).Where("application_id = ?", deleteApp.ID)
	}()

	t.Run("success create app", func(t *testing.T) {
		form := app_v2.NewCreateAppInternalForm(models.NewApplicationRepo(db))
		form.Name = utils.RandString(10)
		form.WebhookApiVersion = models.WebhookAPIVersion202007
		form.Email = "23432"

		resp, err := form.Do(ctx)
		assert.Nil(err)
		assert.NotNil(resp)
		createDevAppResult, ok := resp.(*app_v2.CreateDevAppResult)
		assert.True(ok)
		ctx.Set(middlewares.CtxAppKey, createDevAppResult.Application)
		defaultApp = createDevAppResult.Application
	})

	t.Run("check test", func(t *testing.T) {
		form := app_v2.NewCheckTestAppInternalForm(models.NewApplicationRepo(db))
		resp, err := form.Do(ctx)
		assert.Nil(err)

		wantResp := &forms.AppCheckTestResult{
			Results: []*forms.AppCheckResult{
				{
					Title:         "redirect_url",
					InvalidFields: []string{"redirect_url"},
					IsValid:       false,
				},
			},
		}

		j1, err1 := json.Marshal(resp)
		j2, err2 := json.Marshal(wantResp)
		assert.Nil(err1)
		assert.Nil(err2)
		assert.Equal(j1, j2)
	})

	t.Run("update app setting", func(t *testing.T) {
		const redirectUri = "http://baidu.com"
		form := app.NewUpdateAppForm(models.NewApplicationRepo(db))
		form.UID = defaultApp.UID
		form.RedirectUri = redirectUri
		resp, err := form.Do(ctx)
		assert.Nil(err)
		respApp, ok := resp.(*models.Application)
		assert.True(ok)
		defaultApp.RedirectUri = redirectUri
		resetAppTime(defaultApp)
		resetAppTime(respApp)
		assert.Equal(*defaultApp, *respApp)
	})

	t.Run("get app test status", func(t *testing.T) {
		form := app.NewGetAppFormForInternal(models.NewApplicationRepo(db))
		form.UID = defaultApp.UID
		form.WithCheckResult = true

		resp, err := form.Do(ctx)
		assert.Nil(err)
		result, ok := resp.(app.GetAppFormForInternalResult)
		assert.True(ok)
		assert.NotNil(result.IsAllowTest)
		assert.Equal(*result.IsAllowTest, models.AppNotAllowStatus)
		wantTestCheckResult := &forms.AppCheckTestResult{
			Results: []*forms.AppCheckResult{
				{
					Title:         "redirect_url",
					InvalidFields: []string{},
					IsValid:       true,
				},
			},
		}

		j1, err1 := json.Marshal(wantTestCheckResult)
		j2, err2 := json.Marshal(result.CheckTestResult)
		assert.Nil(err1)
		assert.Nil(err2)
		assert.Equal(j1, j2)
	})

	t.Run("allow app test", func(t *testing.T) {
		form := app_v2.NewCheckTestAppInternalForm(models.NewApplicationRepo(db))
		resp, err := form.Do(ctx)
		assert.Nil(err)

		wantResp := &forms.AppCheckTestResult{Results: []*forms.AppCheckResult{}}

		j1, err1 := json.Marshal(resp)
		j2, err2 := json.Marshal(wantResp)
		assert.Nil(err1)
		assert.Nil(err2)
		assert.Equal(j1, j2)
	})

	t.Run("check submit", func(t *testing.T) {
		form := app_v2.NewCheckSubmitAppInternalForm(models.NewApplicationRepo(db))
		resp, err := form.Do(ctx)
		assert.Nil(err)
		checkSubmitResult, ok := resp.(*forms.AppCheckSubmitResult)
		assert.True(ok)
		assert.NotNil(checkSubmitResult)

		// app listing not existed
		appCheckResults := []*forms.AppCheckResult{
			{
				Title:         "app_information",
				InvalidFields: []string{"icon"},
				IsValid:       false,
			},
			{
				Title:         "urls",
				InvalidFields: []string{"app_url"},
				IsValid:       false,
			},
			{

				Title:         app_locale.AppEnglishListingNotExisted,
				InvalidFields: []string{"create_listing"},
				IsValid:       false,
			},
		}

		j1, err1 := json.Marshal(checkSubmitResult.Results)
		j2, err2 := json.Marshal(appCheckResults)
		assert.Nil(err1)
		assert.Nil(err2)
		assert.Contains(string(j1), string(j2))
	})

	t.Run("create app locale", func(t *testing.T) {
		form := app_locale_c2.NewCreateAppLocaleForm(models.NewApplicationRepo(db))
		form.Locale = models.LOCALE_EN_US

		_, err := form.Do(ctx)
		assert.Nil(err)
	})

	t.Run("update app setting", func(t *testing.T) {
		const (
			icon   = "fsdfsf"
			appURL = "http://baidu.com"
		)
		form := app.NewUpdateAppForm(models.NewApplicationRepo(db))
		form.UID = defaultApp.UID
		form.Icon = icon
		form.AppUri = appURL
		resp, err := form.Do(ctx)
		assert.Nil(err)
		respApp, ok := resp.(*models.Application)
		assert.True(ok)
		defaultApp.Icon = icon
		defaultApp.AppUri = appURL
		resetAppTime(defaultApp)
		resetAppTime(respApp)
		assert.Equal(*defaultApp, *respApp)
	})

	t.Run("update app listing", func(t *testing.T) {
		tag := &models.Tag{}
		testhelper.MakeRecord(db, tag, nil)
		defer db.Delete(tag)
		const featureURL = "http://baidu.com"
		form := app_locale.NewUpdateApplocaleForm(models.NewApplicationRepo(db), models.NewApplicationEventRepo(db))
		form.UID = defaultApp.UID
		form.SubTitle = "SubTitle"
		form.Desc = "desc"
		form.Locale = models.LOCALE_EN_US
		form.TagIds = []uint64{tag.Id}
		form.Features = []*app_locale.Feature{
			{
				ImageUrl: featureURL,
			},
			{
				ImageUrl: featureURL,
			},
			{
				ImageUrl: featureURL,
			},
		}
		form.Pricing = &app_locale.Pricing{
			ChargeType: models.FreeChargeType,
		}

		resp, err := form.Do(ctx)
		assert.Nil(err)
		assert.NotNil(resp)
	})

	t.Run("allow app submit", func(t *testing.T) {
		form := app_v2.NewCheckSubmitAppInternalForm(models.NewApplicationRepo(db))
		resp, err := form.Do(ctx)
		assert.Nil(err)
		checkSubmitResult, ok := resp.(*forms.AppCheckSubmitResult)
		assert.True(ok)
		assert.NotNil(checkSubmitResult)

		appCheckResults := []*forms.AppCheckResult{}

		j1, err1 := json.Marshal(checkSubmitResult.Results)
		j2, err2 := json.Marshal(appCheckResults)
		assert.Nil(err1)
		assert.Nil(err2)
		assert.Equal(string(j1), string(j2))
	})

	t.Run("submit app", func(t *testing.T) {
		form := app_v2.NewUpdateAppStatusInternalForm(models.NewApplicationEventRepo(db), models.NewApplicationRepo(db))
		_, err := form.Do(ctx)
		assert.Nil(err)
	})
}
