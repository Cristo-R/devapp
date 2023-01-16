package config

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"gitlab.shoplazza.site/common/ssoauth"
)

var SSOAuth *ssoauth.GinAuth

func InitSSO() {
	authConfig := ssoauth.AuthConfig{
		CookieName:      Cfg.CookieName,
		Host:            Cfg.AppStoreHost,
		CallbackPath:    Cfg.SSOCallbackPath,
		SSOClientId:     Cfg.SSOClientId,
		SSOClientSecret: Cfg.SSOClientSecret,
		SSOHost:         Cfg.SSOHost,
		StateCookieKey:  Cfg.StateCookieKey,
		SkipStateCheck:  true,
	}

	fmt.Printf("auct config: %+v\n", authConfig)

	v := reflect.ValueOf(authConfig)
	t := reflect.TypeOf(authConfig)
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).String() == "" {
			panic(fmt.Sprintf("auth config %s field can't be empty", t.Field(i).Name))
		}
	}

	SSOAuth = ssoauth.NewGinAuth(&authConfig, SSOLoginSuccess, SSOLoginFailure)
}

// SSO登录认证成功
func SSOLoginSuccess(c *gin.Context) {
	user := c.MustGet(ssoauth.UserKey).(*ssoauth.User)
	c.SetCookie(
		Cfg.CookieName,
		user.EncodeSession,
		60*60*24,
		"/",
		"",
		true,
		true,
	)

	path := "/"
	if c.Query("continue") != "" {
		path = c.Query("continue")
	}
	c.Redirect(http.StatusFound, path)
}

// SSO登录认证失败
func SSOLoginFailure(c *gin.Context) {
	v, ok := c.Get(ssoauth.ErrorKey)
	if !ok {
		c.JSON(500, gin.H{"errors": []string{"get oauth error failed"}})
		return
	}

	err, ok := v.(error)
	if !ok {
		c.JSON(500, gin.H{"errors": []string{"format oauth error failed"}})
		return
	}

	c.JSON(500, gin.H{"errors": []string{err.Error()}})
}

// SSO参数检查（防止恶意攻击）
func SSOParamCheckHandler(c *gin.Context) {
	// state 参数校验
	state := c.Query("state")
	bytes, err := base64.RawURLEncoding.DecodeString(state)
	if err != nil || len(bytes) != 32 {
		c.JSON(400, gin.H{"errors": []string{"invalid state"}})
		c.Abort()
		return
	}
}
