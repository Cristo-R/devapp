package testing

import (
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"gitlab.shoplazza.site/shoplaza/cobra/config"
	"gitlab.shoplazza.site/shoplaza/cobra/models"
	mw "gitlab.shoplazza.site/xiabing/goat.git/middlewares"
)

func ClearDB() {
	config.DB.BlockGlobalUpdate(false)
	config.DB.Delete(&models.Application{})
	config.DB.Delete(&models.AppLocal{})
	config.DB.Delete(&models.AppPlatform{})
	config.DB.Delete(&models.PaymentExtension{})
	config.DB.Delete(&models.Store{})
	config.DB.Delete(&models.Tag{})
	config.DB.Delete(&models.AppTag{})
	config.DB.Delete(&models.ApplicationCollect{})
	config.DB.Delete(&models.Collection{})
	config.DB.Delete(&models.CollectionLocale{})
	config.DB.Delete(&models.AccountAlias{})
	config.DB.Delete(&models.DevApplicationLocale{})
	config.DB.Delete(&models.DevOauthApplication{})
	config.DB.Delete(&models.AppSupport{})
	config.DB.Delete(&models.AppSupportFeedback{})
	config.DB.Delete(&models.OauthAccessToken{})
	config.DB.Delete(&models.InstallTracks{})
	config.DB.Delete(&models.Page{})
	config.DB.Delete(&models.PageCollect{})
	config.DB.Delete(&models.ApplicationPricing{})
	config.DB.Delete(&models.ApplicationPricingPlan{})
	config.DB.Delete(&models.DevApplicationPricing{})
	config.DB.Delete(&models.DevApplicationPricingPlan{})
	config.DB.Delete(&models.ApplicationEvent{})
	config.DB.Delete(&models.ApplicationFeature{})
	config.DB.Delete(&models.DevApplicationFeature{})
	config.DB.Delete(&models.ApplicationExtend{})
	config.DB.BlockGlobalUpdate(true)
}

func GetTestDefaultContext(storeID uint64) (c *gin.Context) {
	ctx := &mw.Context{
		StoreDomain:   "contextStoreDomain",
		AccountDomain: "contextAccountDomain",
		ImageDomain:   "contextImageDomain",
		Origin:        "contextOrigin",
		MerchantID:    "contextOrigin",
		StoreID:       storeID,
		Locale:        "zh-CN",
	}

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()

	c, _ = gin.CreateTestContext(w)
	c.Set("context", ctx)
	return
}
