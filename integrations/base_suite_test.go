package integrations

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"gitlab.shoplazza.site/shoplaza/cobra/cmd"
	"gitlab.shoplazza.site/shoplaza/cobra/config"
	"gitlab.shoplazza.site/shoplaza/cobra/migrations"
	"gitlab.shoplazza.site/shoplaza/cobra/models"
	. "gitlab.shoplazza.site/shoplaza/cobra/testing"
	"gitlab.shoplazza.site/xiabing/goat.git/testhelper"
	"gitlab.shoplazza.site/xiabing/goat.git/webserver"

	"log"
)

type BaseSuite struct {
	suite.Suite
	server         *webserver.Server
	internalServer *webserver.Server
}

func TestBaseSuite(t *testing.T) {
	suite.Run(t, new(BaseSuite))
}

func (suite *BaseSuite) SetupSuite() {
	config.Cfg.Env = "test"
	config.Cfg.Port = 0

	if err := migrations.Migrate(config.DB); err != nil {
		log.Fatalf("Could not migrate: %v", err)
	}

	s, err := cmd.StartServer(context.Background())
	if err != nil {
		panic(err)
	}
	internalS, err := cmd.StartInternalServer(context.Background())
	if err != nil {
		panic(err)
	}
	suite.server = s
	suite.internalServer = internalS
}

func (suite *BaseSuite) SetupTest() {
	ClearDB()
}

func (suite *BaseSuite) TearDownSuite() {
	ClearDB()
	suite.server.Close()
	suite.internalServer.Close()
}

func CreateApp(appInfo map[string]interface{}) *models.Application {
	app := &models.Application{}
	testhelper.MakeRecord(config.DB, app, appInfo)
	return app
}

func CreateDevApp(appInfo map[string]interface{}) *models.DevOauthApplication {
	app := &models.DevOauthApplication{}
	testhelper.MakeRecord(config.DB, app, appInfo)
	return app
}

func CreateTag(tagInfo map[string]interface{}) *models.Tag {
	tag := &models.Tag{}
	testhelper.MakeRecord(config.DB, tag, tagInfo)
	return tag
}

func CreateCollection(collectionInfo map[string]interface{}) *models.Collection {
	collection := &models.Collection{}
	testhelper.MakeRecord(config.DB, collection, collectionInfo)
	return collection
}

func CreateCollectionLocale(collectionLocaleInfo map[string]interface{}) *models.CollectionLocale {
	collectionLocale := &models.CollectionLocale{}
	testhelper.MakeRecord(config.DB, collectionLocale, collectionLocaleInfo)
	return collectionLocale
}

func CreateApplicationCollect(appCollectionInfo map[string]interface{}) *models.ApplicationCollect {
	appCollection := &models.ApplicationCollect{}
	testhelper.MakeRecord(config.DB, appCollection, appCollectionInfo)
	return appCollection
}

func CreateAppTag(appTagInfo map[string]interface{}) *models.AppTag {
	appTag := &models.AppTag{}
	testhelper.MakeRecord(config.DB, appTag, appTagInfo)
	return appTag
}

func CreateLocal(appLocal map[string]interface{}) *models.AppLocal {
	local := &models.AppLocal{}
	testhelper.MakeRecord(config.DB, local, appLocal)
	return local
}

func CreateDevLocal(appLocal map[string]interface{}) *models.DevApplicationLocale {
	local := &models.DevApplicationLocale{}
	testhelper.MakeRecord(config.DB, local, appLocal)
	return local
}

func CreatePlatform(appPlatform map[string]interface{}) *models.AppPlatform {
	platform := &models.AppPlatform{}
	testhelper.MakeRecord(config.DB, platform, appPlatform)
	return platform
}

func CreateAppSupportFeedback(info map[string]interface{}) *models.AppSupportFeedback {
	res := &models.AppSupportFeedback{}
	testhelper.MakeRecord(config.DB, res, info)
	return res
}

func CreateAppSupport(info map[string]interface{}) *models.AppSupport {
	res := &models.AppSupport{}
	testhelper.MakeRecord(config.DB, res, info)
	return res
}

func CreateOauthAccessToken(info map[string]interface{}) *models.OauthAccessToken {
	res := &models.OauthAccessToken{}
	testhelper.MakeRecord(config.DB, res, info)
	return res
}

func CreateInstallTrack(info map[string]interface{}) *models.InstallTracks {
	res := &models.InstallTracks{}
	testhelper.MakeRecord(config.DB, res, info)
	return res
}

func CreateStore(info map[string]interface{}) *models.Store {
	res := &models.Store{}
	testhelper.MakeRecord(config.DB, res, info)
	return res
}

func CreateKeyBenefits(keyBenefits map[string]interface{}) *models.ApplicationKeyBenefit {
	keyBenefit := &models.ApplicationKeyBenefit{}
	testhelper.MakeRecord(config.DB, keyBenefit, keyBenefits)
	return keyBenefit
}

func CreateDevKeyBenefits(devKeyBenefits map[string]interface{}) *models.DevApplicationKeyBenefit {
	keyBenefit := &models.DevApplicationKeyBenefit{}
	testhelper.MakeRecord(config.DB, keyBenefit, devKeyBenefits)
	return keyBenefit
}

func CreatePage(info map[string]interface{}) *models.Page {
	page := &models.Page{}
	testhelper.MakeRecord(config.DB, page, info)
	return page
}

func CreatePageCollect(info map[string]interface{}) *models.PageCollect {
	pageCollect := &models.PageCollect{}
	testhelper.MakeRecord(config.DB, pageCollect, info)
	return pageCollect
}
