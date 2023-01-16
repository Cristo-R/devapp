package script

import (
	"context"
	"fmt"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"gitlab.shoplazza.site/xiabing/goat.git/testhelper"

	"gitlab.shoplazza.site/shoplaza/cobra/config"
	"gitlab.shoplazza.site/shoplaza/cobra/migrations"
	"gitlab.shoplazza.site/shoplaza/cobra/models"
	"gitlab.shoplazza.site/shoplaza/cobra/utils"
)

func test_flushAppAndAppListing(db *gorm.DB, appId uint64) error {
	app, err := models.GetApplication(db, appId)
	if err != nil {
		return err
	}
	if app == nil {
		return fmt.Errorf("app not existed. appID: %d", appId)
	}
	devAppListings, err := models.GetDevAppLocales(db, appId)
	if err != nil {
		return err
	}

	if err := flushAppAndAppListing(context.Background(), db, app, devAppListings); err != nil {
		return err
	}
	return nil
}

func Test_flushAppAndAppListing(t *testing.T) {
	if err := migrations.Migrate(config.DB); err != nil {
		panic(err)
	}

	deleteAppFunc := func(app *models.Application) {
		models.DeleteApplication(config.DB, app)
		models.DeleteDevApplocal(config.DB, app.ID, models.LOCALE_EN_US)
		models.DeleteDevApplocal(config.DB, app.ID, models.LOCALE_ZH_CN)
	}
	t.Run("only english listing", func(t *testing.T) {
		assert := assert.New(t)
		db := config.DB.LogMode(false)
		app := &models.Application{}
		appOverrides := map[string]interface{}{
			"Regions": utils.StringArray{models.RegionCN},
		}
		testhelper.MakeRecord(db, app, appOverrides)
		enListing := &models.DevApplicationLocale{}
		listingOverrides := map[string]interface{}{
			"ApplicationId": app.ID,
			"Locale":        models.LOCALE_EN_US,
			"IsPrimary":     false,
		}
		testhelper.MakeRecord(db, enListing, listingOverrides)
		err := test_flushAppAndAppListing(db, app.ID)
		assert.Nil(err)
		retApp, err := models.GetApplication(db, app.ID)
		assert.Nil(err)
		assert.NotNil(retApp)
		assert.Equal(utils.StringArray{models.RegionWW}, retApp.Regions)
		deleteAppFunc(retApp)
	})

	t.Run("only chinese listing", func(t *testing.T) {
		assert := assert.New(t)
		db := config.DB.LogMode(false)
		app := &models.Application{}
		appOverrides := map[string]interface{}{
			"Regions": utils.StringArray{models.RegionWW},
		}
		testhelper.MakeRecord(db, app, appOverrides)
		zhListing := &models.DevApplicationLocale{}
		listingOverrides := map[string]interface{}{
			"ApplicationId": app.ID,
			"Locale":        models.LOCALE_ZH_CN,
			"IsPrimary":     true,
		}
		testhelper.MakeRecord(db, zhListing, listingOverrides)
		err := test_flushAppAndAppListing(db, app.ID)
		assert.Nil(err)
		retApp, err := models.GetApplication(db, app.ID)
		assert.Nil(err)
		assert.NotNil(retApp)
		assert.Equal(utils.StringArray{models.RegionCN}, retApp.Regions)
		retEnListing, err := models.GetDevAppLocale(db, app.ID, models.LOCALE_EN_US)
		assert.Nil(err)
		assert.NotNil(retEnListing)
		assert.True(retEnListing.IsPrimary)
		deleteAppFunc(retApp)
	})

	t.Run("english and chinese listing all exited, english listing is primary", func(t *testing.T) {
		assert := assert.New(t)
		db := config.DB //.LogMode(false)
		app := &models.Application{}
		appOverrides := map[string]interface{}{
			"Regions": utils.StringArray{models.RegionGlobal},
		}
		testhelper.MakeRecord(db, app, appOverrides)
		enListing := &models.DevApplicationLocale{}
		enListingOverrides := map[string]interface{}{
			"ApplicationId": app.ID,
			"Locale":        models.LOCALE_EN_US,
			"IsPrimary":     true,
		}
		testhelper.MakeRecord(db, enListing, enListingOverrides)
		zhListing := &models.DevApplicationLocale{}
		zhListingOverrides := map[string]interface{}{
			"ApplicationId": app.ID,
			"Locale":        models.LOCALE_ZH_CN,
			"IsPrimary":     false,
		}
		testhelper.MakeRecord(db, zhListing, zhListingOverrides)
		err := test_flushAppAndAppListing(db, app.ID)
		assert.Nil(err)
		retApp, err := models.GetApplication(db, app.ID)
		assert.Nil(err)
		assert.NotNil(retApp)
		assert.Equal(utils.StringArray{models.RegionWW}, retApp.Regions)
		retEnListing, err := models.GetDevAppLocale(db, app.ID, models.LOCALE_EN_US)
		assert.Nil(err)
		assert.NotNil(retEnListing)
		assert.True(retEnListing.IsPrimary)
		retZhListing, err := models.GetDevAppLocale(db, app.ID, models.LOCALE_ZH_CN)
		assert.Nil(err)
		assert.NotNil(retZhListing)
		assert.False(retZhListing.IsPrimary)
		deleteAppFunc(retApp)
	})

	t.Run("english and chinese listing all exited, chinese listing is primary", func(t *testing.T) {
		assert := assert.New(t)
		db := config.DB.LogMode(false)
		app := &models.Application{}
		appOverrides := map[string]interface{}{
			"Regions": utils.StringArray{models.RegionCN},
		}
		testhelper.MakeRecord(db, app, appOverrides)
		enListing := &models.DevApplicationLocale{}
		enListingOverrides := map[string]interface{}{
			"ApplicationId": app.ID,
			"Locale":        models.LOCALE_EN_US,
			"IsPrimary":     false,
		}
		testhelper.MakeRecord(db, enListing, enListingOverrides)
		zhListing := &models.DevApplicationLocale{}
		zhListingOverrides := map[string]interface{}{
			"ApplicationId": app.ID,
			"Locale":        models.LOCALE_ZH_CN,
			"IsPrimary":     true,
		}
		testhelper.MakeRecord(db, zhListing, zhListingOverrides)
		err := test_flushAppAndAppListing(db, app.ID)
		assert.Nil(err)
		retApp, err := models.GetApplication(db, app.ID)
		assert.Nil(err)
		assert.NotNil(retApp)
		assert.Equal(utils.StringArray{models.RegionCN}, retApp.Regions)
		retEnListing, err := models.GetDevAppLocale(db, app.ID, models.LOCALE_EN_US)
		assert.Nil(err)
		assert.NotNil(retEnListing)
		assert.True(retEnListing.IsPrimary)
		retZhListing, err := models.GetDevAppLocale(db, app.ID, models.LOCALE_ZH_CN)
		assert.Nil(err)
		assert.NotNil(retZhListing)
		assert.False(retZhListing.IsPrimary)
		deleteAppFunc(retApp)
	})
}
