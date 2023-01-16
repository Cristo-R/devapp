package actions

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.shoplazza.site/shoplaza/cobra/config"
	adminForm "gitlab.shoplazza.site/shoplaza/cobra/forms/admin"
	"gitlab.shoplazza.site/shoplaza/cobra/models"

	formutil "gitlab.shoplazza.site/xiabing/goat.git/form"
)

func GetAppInfo(c *gin.Context) {
	form := &adminForm.GetAppForm{}
	form.AppId = c.Param("id")

	ok, ret := formutil.DoThisForm(form, c)
	if ok {
		c.JSON(200, ret)
	}
}

func InstallAppFromAppStore(c *gin.Context) {
	form := &adminForm.InstallAppFromAppStoreForm{}
	clientId := c.Param("client_id")

	form.ClientID = clientId
	ok, ret := formutil.DoThisForm(form, c)
	if ok {
		c.JSON(302, ret)
	}
}

func GetInstalledAppsOfStore(c *gin.Context) {
	form := &adminForm.GetInstalledAppsForm{}

	ok, ret := formutil.DoThisForm(form, c)
	if ok {
		c.JSON(200, ret)
	}
}

func InstallAppFromPartnerCenter(c *gin.Context) {
	form := &adminForm.InstallAppFromPartnerCenterForm{}
	clientId := c.Param("client_id")

	form.ClientID = clientId
	ok, ret := formutil.DoThisForm(form, c)
	if ok {
		c.JSON(302, ret)
	}
}

func GetSessionToken(c *gin.Context) {
	form := &adminForm.SessionTokenForm{}

	ok, ret := formutil.DoThisForm(form, c)
	if ok {
		c.JSON(200, ret)
	}
}

func CreateAppSupport(c *gin.Context) {
	form := &adminForm.CreateAppSupportForm{}

	ok, ret := formutil.DoThisForm(form, c)
	if ok {
		c.JSON(200, ret)
	}
}

func AuthorizeEmbbedApp(c *gin.Context) {
	form := &adminForm.AuthorizedEmbbedAppForm{}
	form.ClientID = c.Param("client_id")

	ok, ret := formutil.DoThisForm(form, c)
	if ok {
		c.JSON(302, ret)
	}
}

func CheckOauth(c *gin.Context) {
	form := &adminForm.CheckOauthForm{}
	form.ApplicationId, _ = strconv.ParseInt(c.Param("id"), 10, 64)

	ok, ret := formutil.DoThisForm(form, c)
	if ok {
		c.JSON(200, ret)
	}
}

func GetCollectionsApps(c *gin.Context) {
	form := adminForm.NewGetCollectionsAppsForm(models.NewApplicationRepo(config.DB))

	ok, ret := formutil.DoThisForm(form, c)
	if ok {
		c.JSON(200, ret)
	}
}

func GetGuides(c *gin.Context) {
	form := &adminForm.GetGuidesForm{}

	ok, ret := formutil.DoThisForm(form, c)
	if ok {
		c.JSON(200, ret)
	}
}

func UpdateGuides(c *gin.Context) {
	form := &adminForm.UpdateGuidesForm{}

	ok, ret := formutil.DoThisForm(form, c)
	if ok {
		c.JSON(200, ret)
	}
}
