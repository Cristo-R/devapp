package admin

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"text/template"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"gitlab.shoplazza.site/common/shoplazza-common/xid"
	"gitlab.shoplazza.site/shoplaza/cobra/config"
	"gitlab.shoplazza.site/shoplaza/cobra/models"
	"gitlab.shoplazza.site/shoplaza/cobra/service"
	formutil "gitlab.shoplazza.site/xiabing/goat.git/form"
	"gitlab.shoplazza.site/xiabing/goat.git/i18n"
	"gitlab.shoplazza.site/xiabing/goat.git/middlewares"
)

type CreateAppSupportForm struct {
	AppId            uint64 `form:"app_id" json:"app_id" binding:"required"`
	IssueDescription string `form:"issue_description" json:"issue_description" binding:"required"`
	ContactEmail     string `form:"contact_email" json:"contact_email" binding:"required"`
	SubmitSource     string `form:"submit_source" json:"submit_source"`
}

func (form *CreateAppSupportForm) Do(c *gin.Context) (interface{}, error) {
	appctx, ok := c.MustGet("context").(*middlewares.Context)
	if !ok {
		return nil, formutil.NewInternalStateError(i18n.T("no_context", appctx.Locale))
	}

	db := config.DB
	appSupport := &models.AppSupport{}
	appSupport.ID = xid.Get()
	appSupport.ApplicationId = form.AppId
	appSupport.ContactEmail = form.ContactEmail
	appSupport.IssueDescription = template.HTMLEscapeString(form.IssueDescription)
	appSupport.HasSendFeedbackEmail = false
	appSupport.StoreId = appctx.StoreID
	appSupport.SubmitSource = form.SubmitSource
	appSupport.StoreDomain = fmt.Sprintf("%s.%s", c.GetHeader("Slug"), appctx.StoreDomain)

	h := hmac.New(sha256.New, []byte(config.Cfg.FeebackLinkSecret))
	h.Write([]byte(strconv.Itoa(int(appSupport.ID))))
	feedbackAccesstoken := hex.EncodeToString(h.Sum(nil))
	appSupport.FeedbackAccessToken = feedbackAccesstoken

	app, err := models.GetApplication(db, form.AppId)
	if err != nil {
		log.WithError(err).Errorln("failed to get application by id")
		return nil, formutil.NewInternalStateError(i18n.T("app_info_error", appctx.Locale))
	}
	if app == nil {
		return nil, formutil.NewNotFoundError(i18n.T("app_not_exists", appctx.Locale))
	}

	appLocales, err := models.GetAppLocalesByAppIds(db, []uint64{form.AppId})
	if err != nil {
		log.WithError(err).Errorln("failed to get application locale by id")
		return nil, formutil.NewInternalStateError(i18n.T("app_locale_error", appctx.Locale))
	}

	appLocalesMap := models.ResolveApplocale(appLocales, appctx.Locale)

	if appLocale, ok := appLocalesMap[form.AppId]; ok {
		appSupport.ApplicationListingName = appLocale.Name
		appSupport.ApplicationListingIcon = appLocale.Icon
	}

	// 调用totoro内部接口查询store信息
	store, err := service.GetStoreFromTotoro(cast.ToString(appctx.StoreID))
	if err != nil {
		log.WithError(err).Errorln("Failed to obtain the store")
		return nil, formutil.NewInternalStateError(i18n.T("get_account_store_error", appctx.Locale))
	}
	if store == nil {
		return nil, formutil.NewNotFoundError(i18n.T("get_account_store_error", appctx.Locale))
	}

	appSupport.StoreName = store.Name
	appSupport.StoreLocale = store.UserLocale

	// 查询app开发者的账号信息，取出locale，没有取英文
	appDeveloperLocale := models.LOCALE_EN_US
	user, err := service.GetShoplazzaAccount(app.UserId)
	if err == nil && user != nil {
		appDeveloperLocale = user.Locale
	}

	if err := models.CreateAppSupport(db, appSupport); err != nil {
		return nil, formutil.NewInternalStateError(i18n.T("app_support_create_error", appctx.Locale))
	}

	go func(appSupport *models.AppSupport, appDeveloperEmail, locale string) {
		if err := models.SendAppSupportCreatedEmail(appSupport, appDeveloperEmail, locale); err != nil {
			log.WithError(err).Warn("send create application support message fail")
		}
	}(appSupport, app.Email, appDeveloperLocale)

	return appSupport, nil
}
