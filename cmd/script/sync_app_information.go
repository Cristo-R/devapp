package script

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"gitlab.shoplazza.site/common/plugin-common/xtypes"

	"gitlab.shoplazza.site/shoplaza/cobra/config"
	"gitlab.shoplazza.site/shoplaza/cobra/models"
	"gitlab.shoplazza.site/shoplaza/cobra/pkg/errs"
	"gitlab.shoplazza.site/shoplaza/cobra/service"
)

// SyncAppInformationScript sync app information script, generate app handel to app extend table
func SyncAppInformationScript(ctx context.Context) cli.Command {
	return cli.Command{
		Name:  "sync-app-information",
		Usage: "app information sync to app extend table",
		Action: func(c *cli.Context) (err error) {
			appIds := make([]uint64, 0)
			flags := c.String("app_ids")
			if flags != "" {
				for _, str := range strings.Split(flags, ",") {
					id, err := strconv.ParseUint(str, 10, 64)
					if err != nil {
						log.Errorf("unknown app id: %s err: %s", str, err.Error())
					}
					appIds = append(appIds, id)
				}
			}
			syncAppInformationScript(ctx, appIds)
			return nil
		},
		Flags: []cli.Flag{
			cli.StringFlag{Name: "app_ids", Usage: "app ids", Required: false},
		},
	}
}

// syncAppHandleToAppExtendTable sync app handle to app extend table
func syncAppInformationScript(ctx context.Context, appIds []uint64) {
	beginTime := time.Now()
	appCh := make(chan *models.Application, 20)
	db := config.DB.LogMode(false)

	go getAllAppInformation(ctx, db, appIds, false, appCh)

	appRepo := models.NewApplicationRepo(db)
	totalCount := 0
	errorAppIds := make([]uint64, 0)
	for app := range appCh {
		totalCount++
		_, err := appRepo.GetApplicationExtendByAppId(ctx, xtypes.XId(app.ID))
		if err == nil {
			continue
		}
		if !errors.Is(err, errs.NotFoundErr) {
			errorAppIds = append(errorAppIds, app.ID)
			log.Errorf("failed query app extend is exited. appID: %d", app.ID)
			continue
		}

		handle, err := service.GetUniqueApplicationHandle(ctx, appRepo, app.Name)
		if err != nil {
			errorAppIds = append(errorAppIds, app.ID)
			log.Errorf("%s. appid: %d appName: %s err count: %d", err.Error(), app.ID, app.Name, len(errorAppIds))
			continue
		}

		now := time.Now().UTC()
		appExtend := &models.ApplicationExtend{
			ApplicationId: xtypes.XId(app.ID),
			Handle:        handle,
			CreatedAt:     now,
			UpdatedAt:     now,
		}
		if err := appRepo.UpsertApplicationExtend(ctx, appExtend); err != nil {
			errorAppIds = append(errorAppIds, app.ID)
			log.Errorf("failed to upsert app extend. err: %s appExtend: %+v", err.Error(), appExtend)
			continue
		}
	}

	log.Infof("success to sync app information to app extend. total time: %s total count: %d error count: %d", time.Since(beginTime).String(), totalCount, len(errorAppIds))

	// output error app ids
	if len(errorAppIds) != 0 {
		promptStr := ""
		for _, ids := range errorAppIds {
			promptStr = promptStr + strconv.FormatUint(ids, 10) + ","
		}
		log.Infof("some data fails to be sync. please run command: %s %s", "sync-app-information-script --app_ids", promptStr[:len(promptStr)-1])
	}
}
