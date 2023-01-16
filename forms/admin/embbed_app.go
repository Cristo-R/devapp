package admin

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.shoplazza.site/shoplaza/cobra/config"
	"gitlab.shoplazza.site/shoplaza/cobra/models"
	formutil "gitlab.shoplazza.site/xiabing/goat.git/form"
	"gitlab.shoplazza.site/xiabing/goat.git/i18n"
	"gitlab.shoplazza.site/xiabing/goat.git/middlewares"
)

type AuthorizedEmbbedAppForm struct {
	ClientID string `form:"client_id" json:"client_id" binding:"required"`
}

func (form *AuthorizedEmbbedAppForm) Do(c *gin.Context) (interface{}, error) {
	appctx, ok := c.MustGet("context").(*middlewares.Context)
	if !ok {
		return nil, formutil.NewInternalStateError(i18n.T("no_context", appctx.Locale))
	}

	query := strings.Split(c.Request.RequestURI, "?")
	if len(query) > 2 {
		return nil, formutil.NewUnprocessableError("invalid url")
	}

	if len(query) == 2 && c.Query("hmac") == "" {
		return nil, formutil.NewUnprocessableError("if url wants to contains other params, must use app secret validate other params and generate hmac param")
	}

	app, err := models.GetApplicationByUid(config.DB, form.ClientID)
	if err != nil {
		return nil, err
	}
	if app == nil {
		return nil, formutil.NewNotFoundError("app is not exsited")
	}

	params := url.Values{}
	hm := hmac.New(sha256.New, []byte(app.Secret))

	// uri = https://pjs.myshoplaza.com/admin/authorized_embbed_apps/123?a=aa&b=bb&hmac=ajhsdkjahskdjashdk
	if len(query) == 2 {
		params, err = url.ParseQuery(query[1])
		if err != nil {
			return nil, formutil.NewUnprocessableError("invalid url")
		}

		// 校验参数签名
		v := params.Get("hmac")
		params.Del("hmac")
		hm.Write([]byte(params.Encode()))
		signature := hex.EncodeToString(hm.Sum(nil))
		if !hmac.Equal([]byte(signature), []byte(v)) {
			return nil, formutil.NewUnprocessableError("signature does not match, it may have been tampered with")
		}
	}
	params.Set("authorized", "true")

	slug := c.GetHeader("Slug")
	storeDomain := c.GetHeader("Store-Domain")
	domain := slug + "." + storeDomain
	redirectUrl := fmt.Sprintf("https://%s/admin/smart_apps/angora/app_store/plugins/%d?%s", domain, app.ID, params.Encode())

	c.Redirect(http.StatusFound, redirectUrl)

	return nil, nil
}
