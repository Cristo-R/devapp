package cmd

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gitlab.shoplazza.site/xiabing/goat.git/i18n"
	"gitlab.shoplazza.site/xiabing/goat.git/validators"
	"gitlab.shoplazza.site/xiabing/goat.git/webserver"

	"gitlab.shoplazza.site/shoplaza/cobra/config"
)

func StartServerAndBlock(ctx context.Context) error {
	if server, err := StartServer(ctx); err != nil {
		return err
	} else {
		server.Block()
	}

	return nil
}

func StartServer(ctx context.Context) (*webserver.Server, error) {
	if config.Cfg.Env != "dev" {
		gin.SetMode(gin.ReleaseMode)
		gin.DisableConsoleColor()
	}

	binding.Validator = &validators.DefaultValidator{}

	if err := i18n.Bundle(config.Cfg.BasePath+"/en.json", config.Cfg.BasePath+"/zh.json"); err != nil {
		return nil, err
	}

	return webserver.StartHTTPServer(ctx, router(), webserver.Port(config.Cfg.Port))
}
