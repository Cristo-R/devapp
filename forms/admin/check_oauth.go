package admin

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gitlab.shoplazza.site/shoplaza/cobra/config"
	"gitlab.shoplazza.site/shoplaza/cobra/models"
	formutil "gitlab.shoplazza.site/xiabing/goat.git/form"
	"gitlab.shoplazza.site/xiabing/goat.git/i18n"
	"gitlab.shoplazza.site/xiabing/goat.git/middlewares"
)

type CheckOauthForm struct {
	ApplicationId int64 `form:"application_id" json:"application_id" binding:"required"`
}

func (form *CheckOauthForm) Do(c *gin.Context) (interface{}, error) {
	appctx, ok := c.MustGet("context").(*middlewares.Context)
	if !ok {
		return nil, formutil.NewInternalStateError(i18n.T("no_context", appctx.Locale))
	}

	app, err := models.GetApplication(config.DB, uint64(form.ApplicationId))
	if err != nil {
		log.WithError(err)
		return nil, formutil.NewInternalStateError(i18n.T("app_error", appctx.Locale))
	}
	if app == nil {
		return nil, formutil.NewNotFoundError(i18n.T("app_not_exits", appctx.Locale))
	}

	// 查询 app 是否被安装到指定店铺
	storeId := strconv.Itoa(int(appctx.StoreID))
	installTrack, err := models.GetInstallTrack(config.DB, storeId, form.ApplicationId)
	if err != nil {
		log.WithError(err).Errorln("failed to get install track")
		return nil, formutil.NewInternalStateError(i18n.T("install_track_error", appctx.Locale))
	}
	if installTrack == nil {
		return nil, formutil.NewNotFoundError(i18n.T("app_not_install", appctx.Locale))
	}

	token, err := models.GetOauthAccessToken(config.DB, uint64(form.ApplicationId))
	if err != nil {
		log.WithError(err)
		return nil, formutil.NewInternalStateError(i18n.T("oauth_access_token_error", appctx.Locale))
	}
	if token == nil {
		return nil, formutil.NewNotFoundError(i18n.T("oauth_access_token_no_exist", appctx.Locale))
	}

	// 如果app的授权时间在安装时间之前，等同于未授权
	if token.CreatedAt.Before(installTrack.InstalledAt) {
		return nil, formutil.NewNotFoundError(i18n.T("app_unauthorized", appctx.Locale))
	}

	// 判断 app oauthdancable是否为 true，embbed是否为true
	if app.OauthDancable == false || app.Embbed == false || app.AppUri == "" {
		return nil, formutil.NewNotFoundError(i18n.T("check_embbed_app_error", appctx.Locale))
	}

	session, err := c.Cookie("awesomev2")
	if err != nil {
		return nil, formutil.NewInternalStateError(i18n.T("session_id_not_exist", appctx.Locale))
	}

	slug := c.GetHeader("Slug")
	storeDomain := c.GetHeader("Store-Domain")
	domain := slug + "." + storeDomain

	query := url.Values{
		"shop":    {domain},
		"session": {session},
	}

	key := []byte(app.Secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(query.Encode()))
	sha := hex.EncodeToString(h.Sum(nil))
	query.Set("hmac", sha)

	redirectUrl := app.AppUri + "?" + query.Encode()

	return map[string]interface{}{"authorized": true, "redirectUrl": redirectUrl}, nil
}
