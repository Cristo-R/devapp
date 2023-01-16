package cmd

import (
	"github.com/getsentry/raven-go"
	"github.com/gin-contrib/sentry"
	"github.com/gin-gonic/gin"
	"gitlab.shoplazza.site/shoplaza/cobra/actions"
	"gitlab.shoplazza.site/shoplaza/cobra/config"
	"gitlab.shoplazza.site/shoplaza/cobra/middlewares"
)

func router() *gin.Engine {

	r := gin.New()
	//r.LoadHTMLGlob(config.Cfg.BasePath + "/templates/*")
	if config.Cfg.Env == "local" || config.Cfg.Env == "test" {
		gin.SetMode("debug")
	} else {
		gin.SetMode("release")
	}
	r.GET("/ping", actions.Ping)

	r.Use(middlewares.Prometheus())

	if config.Cfg.SentryDSN != "" {
		raven.SetDSN(config.Cfg.SentryDSN)
		r.Use(sentry.Recovery(raven.DefaultClient, false))
	} else {
		r.Use(gin.Recovery())
	}

	r.Use(middlewares.AppStoreLocaleContext())

	appsAPI := r.Group("/api/platform/apps")
	{
		appsAPI.GET("", actions.GetApps)
		appsAPI.GET("/:id", actions.GetApp)
	}

	r.GET("/api/app_store/app/:id/install", actions.InstallApp)

	return r
}
