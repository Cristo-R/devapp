package admin

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cast"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gitlab.shoplazza.site/shoplaza/cobra/config"
	"gitlab.shoplazza.site/shoplaza/cobra/models"
	formutil "gitlab.shoplazza.site/xiabing/goat.git/form"
	"gitlab.shoplazza.site/xiabing/goat.git/i18n"
	"gitlab.shoplazza.site/xiabing/goat.git/middlewares"
)

type AppView struct {
	Id            uint64     `json:"id"`
	AppUri        string     `json:"app_uri"`
	Embbed        bool       `json:"embbed"`
	Icon          string     `json:"icon"`
	Link          string     `json:"link"`
	OauthDancable bool       `json:"oauth_dancable"`
	RedirectUri   string     `json:"redirect_uri"`
	Scopes        []string   `json:"scopes"`
	Status        string     `json:"status"`
	Uid           string     `json:"uid"`
	Subtitle      string     `json:"subtitle"`
	Title         string     `json:"title"`
	AppOwner      string     `json:"app_owner"`
	Locale        string     `json:"locale"`
	Listing       bool       `json:"listing"`
	Name          string     `json:"name"`
	IsSmartApp    bool       `json:"is_smart_app"`
	InstalledAt   *time.Time `json:"installed_at"`
	Email         string     `json:"email"`
	HasReview     bool       `json:"has_review"`
}

type GetAppForm struct {
	AppId string `form:"app_id" json:"app_id" binding:"required"`
}

func (form *GetAppForm) Do(c *gin.Context) (interface{}, error) {
	appctx, ok := c.MustGet("context").(*middlewares.Context)
	if !ok {
		return nil, formutil.NewInternalStateError(i18n.T("no_context", appctx.Locale))
	}

	db := config.DB
	appView := AppView{}
	var app *models.Application

	searchCondition, err := strconv.ParseUint(form.AppId, 10, 64)
	// 如果转换为整数出错，说明传入的是app_name，需要根据app_name来查询app信息
	// 之所以加上这段根据app_name查询的逻辑，是因为插件组在其他地方仍然通过app_name来调用接口
	// 后续回让他们都改为根据app_id调用，等他们都改成app_id调用后，这段根据app_name查询的逻辑可以去除
	if err != nil {
		app, err = models.GetApplicationByName(db, form.AppId)
	} else {
		app, err = models.GetApplication(db, searchCondition)
	}

	if err != nil {
		log.WithError(err).Errorln("failed to get application by id")
		return nil, formutil.NewInternalStateError(i18n.T("app_error", appctx.Locale))
	}
	if app == nil {
		return nil, formutil.NewNotFoundError(i18n.T("app_not_exits", appctx.Locale))
	}

	appLocales, err := models.GetAppLocalesByAppIds(db, []uint64{app.ID})
	if err != nil {
		log.WithError(err).Errorln("failed to get application locale by id")
		return formutil.NewInternalStateError(i18n.T("app_locale_error", appctx.Locale)), nil
	}

	appView.Id = app.ID
	appView.AppUri = app.AppUri
	appView.Embbed = app.Embbed
	appView.Link = app.Link
	appView.OauthDancable = app.OauthDancable
	appView.RedirectUri = app.RedirectUri
	appView.Scopes = strings.Split(app.Scopes, " ")
	appView.Status = app.Status
	appView.Uid = app.UID
	appView.Listing = app.Listing
	appView.Name = app.Name
	appView.IsSmartApp = app.IsSmartApp
	appView.Email = app.Email

	storeId := strconv.Itoa(int(appctx.StoreID))
	installTrack, err := models.GetInstallTrack(db, storeId, int64(app.ID))
	if err != nil {
		log.WithError(err).Errorln("failed to get install track")
		return formutil.NewInternalStateError(i18n.T("install_track_error", appctx.Locale)), nil
	}
	if installTrack != nil {
		appView.InstalledAt = &installTrack.InstalledAt
	}

	// 如果app状态为已上架，显示listing信息；否则，显示app信息
	if app.Status == models.APPSTATUS_PUBLISHED {
		appLocalsMap := models.ResolveApplocale(appLocales, appctx.Locale)
		if appLocale, ok := appLocalsMap[app.ID]; ok {
			appView.Title = appLocale.Name
			appView.Subtitle = appLocale.Subtitle
			appView.Icon = appLocale.Icon
			appView.Locale = appLocale.Locale
		}
	} else {
		appView.Title = app.Name
		appView.Subtitle = ""
		json.Unmarshal([]byte(app.Icon), &appView.Icon)
		appView.Locale = appctx.Locale
	}

	owner := app.Owner
	if app.PartnerId == 0 {
		owner = ""
	}
	appView.AppOwner = owner
	review := &models.Reviews{
		StoreId: cast.ToInt64(appctx.StoreID),
		AppId:   cast.ToUint64(form.AppId),
	}

	review, err = review.GetValidReviewsByAppIdAndStoreId(config.DB)
	if err != nil {
		return nil, err
	}
	if review == nil {
		appView.HasReview = false
	} else {
		appView.HasReview = true
	}

	return appView, nil
}
