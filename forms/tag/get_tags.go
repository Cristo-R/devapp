package tag

import (
	"github.com/gin-gonic/gin"
	"gitlab.shoplazza.site/shoplaza/cobra/config"
	"gitlab.shoplazza.site/shoplaza/cobra/forms"
	"gitlab.shoplazza.site/shoplaza/cobra/models"
	formutil "gitlab.shoplazza.site/xiabing/goat.git/form"
	"gitlab.shoplazza.site/xiabing/goat.git/i18n"
	"gitlab.shoplazza.site/xiabing/goat.git/middlewares"
)

type GetTagsForm struct{}

func (form *GetTagsForm) Do(c *gin.Context) (interface{}, error) {
	appctx, ok := c.MustGet("context").(*middlewares.Context)
	if !ok {
		return nil, formutil.NewInternalStateError(i18n.T("no_context", appctx.Locale))
	}

	switch {
	case appctx.Origin == "platform":
		if err := forms.PlatformAuth(c, "read"); err != nil {
			return nil, err
		}
		break
	default:
		return nil, formutil.NewUnprocessableError("not support the value of origin field in httpHeader")
	}

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
