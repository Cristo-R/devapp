package script

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"gitlab.shoplazza.site/shoplaza/cobra/config"
	"gitlab.shoplazza.site/shoplaza/cobra/models"
	"gitlab.shoplazza.site/shoplaza/cobra/service"
	"gitlab.shoplazza.site/shoplaza/cobra/utils"
)

func FlushAppListingScript(ctx context.Context) cli.Command {
	return cli.Command{
		Name:  "flush-app-and-app-listing",
		Usage: "flush app and app listing",
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
			flushAppListingScript(ctx, appIds)
			return nil
		},
		Flags: []cli.Flag{
			cli.StringFlag{Name: "app_ids", Usage: "app ids", Required: false},
		},
	}
}

func flushAppListingScript(ctx context.Context, appIds []uint64) {
	beginTime := time.Now()
	appCh := make(chan *models.Application, 20)
	db := config.DB.LogMode(false)

	go getAllAppInformation(ctx, db, appIds, true, appCh)

	totalCount := 0
	errorAppIds := make([]uint64, 0)
	for app := range appCh {
		totalCount++

		devAppListings, err := models.GetDevAppLocales(db, app.ID)
		if err != nil {
			errorAppIds = append(errorAppIds, app.ID)
			log.Errorf("failed to get dev app locales. app: %+v", app)
			continue
		}

		if err := flushAppAndAppListing(ctx, db, app, devAppListings); err != nil {
			errorAppIds = append(errorAppIds, app.ID)
			log.Errorf("failed to flush app and app listing. app: %+v", app)
			continue
		}
	}
	log.Infof("success to flush app and app listing. total time: %s total count: %d error count: %d", time.Since(beginTime).String(), totalCount, len(errorAppIds))

	if len(errorAppIds) != 0 {
		promptStr := ""
		for _, ids := range errorAppIds {
			promptStr += strconv.FormatUint(ids, 10)
			promptStr += ","
		}
		log.Infof("some data fails to be flush. please run command: %s %s", "flush-app-and-app-listing-script --app_ids", promptStr[:len(promptStr)-1])
	}
}

// getAllAppInformation get all the app information, filter private app
func getAllAppInformation(ctx context.Context, db *gorm.DB, appIDs []uint64, isFilterPrivateApp bool, appCh chan<- *models.Application) {
	defer close(appCh)
	var (
		lastAppId uint64
		pageSize  = 100
	)

	for {
		apps := make([]*models.Application, 0)
		// find all public app
		sql := db.Table("oauth_applications").Order("id ASC")
		if isFilterPrivateApp {
			// filter private app
			sql = sql.Where("private_app = 0")
		}
		if lastAppId > 0 {
			sql = sql.Where("id > ?", lastAppId)
		}

		if len(appIDs) != 0 {
			sql = sql.Where("id IN (?)", appIDs)
		}

		if len(appIDs) == 0 {
			sql = sql.Limit(pageSize)
		}

		err := sql.Find(&apps).Error
		if err != nil && !gorm.IsRecordNotFoundError(err) {
			log.Errorf("failed to get app information. err: %s", err.Error())
			return
		}

		if gorm.IsRecordNotFoundError(err) || len(apps) == 0 {
			return
		}

		for _, app := range apps {
			select {
			case <-ctx.Done():
				log.Infof("context cannel")
				return
			case appCh <- app:
			}
		}

		if len(appIDs) != 0 {
			return
		}

		lastAppId = apps[len(apps)-1].ID
	}
}

func flushAppAndAppListing(ctx context.Context, db *gorm.DB, app *models.Application, appListings []models.DevApplicationLocale) error {
	localeMap := service.DevAppLocaleConvertToLocaleMap(appListings)

	if len(localeMap) == 0 {
		return nil
	}

	enListing, isExistedEnListing := localeMap[models.LOCALE_EN_US]
	zhListing, isExistedZhListing := localeMap[models.LOCALE_ZH_CN]
	tx := db.Begin()
	defer tx.Callback()
	// only english listing
	// set app regions -> new version (Worldwide) -> old version (All)
	// update english listing `IsPrimary` to true
	if isExistedEnListing && !isExistedZhListing {
		app.Regions = utils.StringArray{models.RegionWW}
		if err := models.UpdateApplication(tx, app); err != nil {
			return err
		}

		enListing.IsPrimary = true
		if err := models.UpdateDevAppLocal(tx, enListing); err != nil {
			return err
		}
	}

	// only chinese listing
	// set app regions -> new version (CN) -> old version (CN)
	// create english listing and set `IsPrimary` to true
	// update chinese listing `IsPrimary` to false
	if !isExistedEnListing && isExistedZhListing {
		app.Regions = utils.StringArray{models.RegionCN}
		if err := models.UpdateApplication(tx, app); err != nil {
			return err
		}

		// create default enlisting listing
		devAppLocale := &models.DevApplicationLocale{
			Name:          app.Name,
			ApplicationId: app.ID,
			Locale:        models.LOCALE_EN_US,
			IsPrimary:     true,
			CreatedAt:     time.Now().UTC(),
			UpdatedAt:     time.Now().UTC(),
			SupportEmail:  app.Email,
			Status:        models.APPSTATUS_DRAFT,
		}
		if err := models.CreateDevAppLocal(tx, devAppLocale); err != nil {
			return err
		}

		if err := models.UpdateDevAppLocaleIsPrimaryField(tx, app.ID, zhListing.Locale, false); err != nil {
			return err
		}
	}

	// english and chinese all existed
	if isExistedEnListing && isExistedZhListing {
		if enListing.IsPrimary && zhListing.IsPrimary {
			return fmt.Errorf("unknown app listings. enListing: %+v zhListing: %+v", *enListing, *zhListing)
		}
		// if english listing is primary, update app region to Worldwide
		if enListing.IsPrimary {
			app.Regions = utils.StringArray{models.RegionWW}
			if err := models.UpdateApplication(tx, app); err != nil {
				return err
			}
		}

		// if chinese listing is primary, update app region to CN and update english listing to primary listing.
		if zhListing.IsPrimary {
			app.Regions = utils.StringArray{models.RegionCN}
			if err := models.UpdateApplication(tx, app); err != nil {
				return err
			}
		}

		if err := models.UpdateDevAppLocaleIsPrimaryField(tx, app.ID, enListing.Locale, true); err != nil {
			return err
		}
		if err := models.UpdateDevAppLocaleIsPrimaryField(tx, app.ID, zhListing.Locale, false); err != nil {
			return err
		}
	}

	tx.Commit()
	return nil
}
