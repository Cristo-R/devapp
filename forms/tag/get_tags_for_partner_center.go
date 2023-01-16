package tag

import (
	"github.com/gin-gonic/gin"
	formutil "gitlab.shoplazza.site/xiabing/goat.git/form"

	"gitlab.shoplazza.site/shoplaza/cobra/config"
	"gitlab.shoplazza.site/shoplaza/cobra/models"
)

type GetTagsForPartnerCenterForm struct{}

func (form *GetTagsForPartnerCenterForm) Do(c *gin.Context) (interface{}, error) {

	tags, err := models.GetTags(config.DB)
	if err != nil {
		return nil, formutil.NewInternalStateError("failed to get tags")
	}
	if tags == nil {
		return nil, formutil.NewNotFoundError("tags is not existed")
	}

	tagsView, _ := models.ConvertTagsToTagViews(tags)

	return map[string]interface{}{"tags": tagsView}, nil
}
