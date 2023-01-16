package admin

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.shoplazza.site/shoplaza/cobra/config"
	"gitlab.shoplazza.site/shoplaza/cobra/models"
	formutil "gitlab.shoplazza.site/xiabing/goat.git/form"
	"gitlab.shoplazza.site/xiabing/goat.git/i18n"
	"gitlab.shoplazza.site/xiabing/goat.git/middlewares"
)

type (
	InstallAppFromAppStoreForm struct {
		ClientID string `form:"client_id" json:"client_id" binding:"required"`
	}
	InstallAppFromPartnerCenterForm struct {
		ClientID string `form:"client_id" json:"client_id" binding:"required"`
	}
)

func (form *InstallAppFromAppStoreForm) Do(c *gin.Context) (interface{}, error) {
	appctx, ok := c.MustGet("context").(*middlewares.Context)
	if !ok {
		return nil, formutil.NewInternalStateError(i18n.T("no_context", appctx.Locale))
	}

	slug := c.GetHeader("Slug")
	storeDomain := c.GetHeader("Store-Domain")
	domain := slug + "." + storeDomain
	storeId := strconv.Itoa(int(appctx.StoreID))

	db := config.DB
	app, err := models.GetApplicationByUid(db, form.ClientID)
	if err != nil {
		return nil, err
	}
	if app == nil {
		return nil, formutil.NewNotFoundError("app_not_exits")
	}

	// 按key值的字典排序
	querystring := "install_from=app_store&shop=" + domain + "&store_id=" + storeId

	key := []byte(app.Secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(querystring))
	sha := hex.EncodeToString(h.Sum(nil))

	querystring = "?hmac=" + sha + "&install_from=app_store&shop=" + domain + "&store_id=" + storeId

	var redirectUrl string
	if app.AppUri != "" {
		redirectUrl = app.AppUri + querystring
	} else {
		redirectUrl = "https://" + domain + "/admin/oauth/authorize?client_id=" + form.ClientID + "&scope=" + app.Scopes + "&redirect_uri=" + url.QueryEscape(app.RedirectUri) + "&response_type=code&install_from=app_store"
	}

	c.Request.Header.Add("X-Shoplazza-Shop-Id", storeId)
	c.Redirect(http.StatusFound, redirectUrl)

	return nil, nil
}

func (form *InstallAppFromPartnerCenterForm) Do(c *gin.Context) (interface{}, error) {
	appctx, ok := c.MustGet("context").(*middlewares.Context)
	if !ok {
		return nil, formutil.NewInternalStateError(i18n.T("no_context", appctx.Locale))
	}

	slug := c.GetHeader("Slug")
	storeDomain := c.GetHeader("Store-Domain")
	domain := slug + "." + storeDomain
	storeId := strconv.Itoa(int(appctx.StoreID))

	db := config.DB
	app, err := models.GetApplicationByUid(db, form.ClientID)
	if err != nil {
		return nil, err
	}
	if app == nil {
		return nil, formutil.NewNotFoundError(i18n.T("app_not_exits", appctx.Locale))
	}
	query := url.Values{
		"install_from": {"partner_center"},
		"shop":         {domain},
		"store_id":     {storeId},
	}

	key := []byte(app.Secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(query.Encode()))
	sha := hex.EncodeToString(h.Sum(nil))
	query.Set("hmac", sha)

	var redirectUrl string
	if app.AppUri != "" {
		redirectUrl = app.AppUri + "?" + query.Encode()
	} else {
		query = url.Values{}
		query.Set("client_id", form.ClientID)
		query.Set("scope", app.Scopes)
		query.Set("redirect_uri", app.RedirectUri)
		query.Set("response_type", "code")
		query.Set("install_from", "partner_center")
		redirectUrl = "https://" + domain + "/admin/oauth/authorize?" + query.Encode()
	}

	c.Request.Header.Add("X-Shoplazza-Shop-Id", storeId)
	c.Redirect(http.StatusFound, redirectUrl)

	return nil, nil
}
