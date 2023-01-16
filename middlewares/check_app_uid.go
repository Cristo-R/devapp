package middlewares

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gitlab.shoplazza.site/xiabing/goat.git/i18n"

	"gitlab.shoplazza.site/shoplaza/cobra/config"
	"gitlab.shoplazza.site/shoplaza/cobra/models"
)

func CheckAppUidMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		appUID := c.Param("uid")
		if appUID == "" {
			log.Errorf("no app uid. request params: %+v", c.Params)
			c.AbortWithStatus(422)
			return
		}

		db := config.DB
		app, err := models.GetApplicationByUid(db, appUID)
		if err != nil {
			log.Errorf("failed to get app by uid, err %s", err.Error())
			c.AbortWithStatus(500)
			return
		}
		if app == nil {
			log.Infof("app is not exits, uid=%s", appUID)
			c.Abort()
			c.String(403, i18n.T("app_not_exits", models.LOCALE_EN_US))
			return
		}

		// add app to gin context
		// reduce the number of queries to the app
		c.Set(CtxAppKey, app)
	}
}
