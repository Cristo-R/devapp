package actions

import (
	"strconv"

	"github.com/gin-gonic/gin"
	tagForm "gitlab.shoplazza.site/shoplaza/cobra/forms/tag"
	formutil "gitlab.shoplazza.site/xiabing/goat.git/form"
)

func GetTags(c *gin.Context) {
	form := &tagForm.GetTagsForm{}
	ok, ret := formutil.DoThisForm(form, c)
	if ok {
		c.JSON(200, ret)
	}
}

func UpdateAppTags(c *gin.Context) {
	form := &tagForm.UpdateTagsForm{}
	form.ID, _ = strconv.ParseUint(c.Param("id"), 10, 64)
	ok, ret := formutil.DoThisForm(form, c)
	if ok {
		c.JSON(200, ret)
	}
}

func GetTagsForInternal(c *gin.Context) {
	form := &tagForm.GetTagsForPartnerCenterForm{}
	ok, ret := formutil.DoThisForm(form, c)
	if ok {
		c.JSON(200, ret)
	}
}
