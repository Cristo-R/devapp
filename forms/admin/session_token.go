package admin

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.shoplazza.site/shoplaza/cobra/config"
	"gitlab.shoplazza.site/shoplaza/cobra/models"
	"gitlab.shoplazza.site/shoplaza/cobra/utils"
	formutil "gitlab.shoplazza.site/xiabing/goat.git/form"
	"gitlab.shoplazza.site/xiabing/goat.git/i18n"
	"gitlab.shoplazza.site/xiabing/goat.git/middlewares"
)

type SessionTokenForm struct {
	ClientId string `form:"client_id" json:"client_id" binding:"required"`
}

type User struct {
	Contact string `json:"contact" faker:"email"`
	Locale  string `json:"locale"`
}

type SessionTokenClaims struct {
	Locale      string `json:"locale"`
	Account     string `json:"account"`
	Destination string `json:"dest"`
	SessionId   string `json:"sid"`
	jwt.StandardClaims
}

func (form *SessionTokenForm) Do(c *gin.Context) (interface{}, error) {
	appctx, ok := c.MustGet("context").(*middlewares.Context)
	if !ok {
		return nil, formutil.NewInternalStateError(i18n.T("no_context", appctx.Locale))
	}

	sessionId, err := c.Cookie("awesomev2")
	if err != nil {
		return nil, formutil.NewInternalStateError(i18n.T("session_id_not_exist", appctx.Locale))
	}

	app, err := models.GetApplicationByUid(config.DB, form.ClientId)
	if err != nil {
		return nil, err
	}
	if app == nil {
		return nil, formutil.NewNotFoundError(i18n.T("app_not_exits", appctx.Locale))
	}

	domain := c.GetHeader("System-Domain")

	user := User{}
	url := fmt.Sprintf("%s/api/user/%s", config.Cfg.StoreService, appctx.LoginUserID)
	content, status := utils.HttpRequest("GET", url, nil)
	if status/100 == 2 {
		err := json.Unmarshal(content, &user)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, formutil.NewInternalStateError("failed to get user info")
	}

	jwtSecret := []byte(app.Secret)

	nowTime := time.Now()
	claims := SessionTokenClaims{
		Account:     user.Contact,
		Locale:      user.Locale,
		Destination: domain,
		SessionId:   sessionId,
		StandardClaims: jwt.StandardClaims{
			Audience:  app.UID,
			Subject:   appctx.LoginUserID,
			Issuer:    fmt.Sprintf("https://%s/admin", domain),
			ExpiresAt: nowTime.Add(1 * time.Minute).Unix(),
			IssuedAt:  nowTime.Unix(),
			NotBefore: nowTime.Unix(),
			Id:        uuid.New().String(),
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{"session_token": token}, nil
}
