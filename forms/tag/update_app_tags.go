package tag

import (
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.shoplazza.site/shoplaza/cobra/config"
	"gitlab.shoplazza.site/shoplaza/cobra/forms"
	"gitlab.shoplazza.site/shoplaza/cobra/models"
	formutil "gitlab.shoplazza.site/xiabing/goat.git/form"
	"gitlab.shoplazza.site/xiabing/goat.git/i18n"
	"gitlab.shoplazza.site/xiabing/goat.git/middlewares"
)

type UpdateTagsForm struct {
	ID     uint64   `json:"id" binding:"required"`
	TagIds []uint64 `json:"tag_ids" binding:"required"`
}

func (form *UpdateTagsForm) Do(c *gin.Context) (interface{}, error) {
	appctx, ok := c.MustGet("context").(*middlewares.Context)
	if !ok {
		return nil, formutil.NewInternalStateError(i18n.T("no_context", appctx.Locale))
	}

	switch {
	case appctx.Origin == "platform":
		if err := forms.PlatformAuth(c, "write"); err != nil {
			return nil, err
		}
		break
	default:
		return nil, formutil.NewUnprocessableError("not support the value of origin field in httpHeader")
	}

	tx := config.DB.Begin()
	defer tx.Rollback()

	err := models.DeleteAppTags(tx, form.ID)
	if err != nil {
		return nil, formutil.NewInternalStateError("failed to update tags")
	}

	for _, v := range form.TagIds {
		appTag := models.AppTag{
			ApplicationId: form.ID,
			TagId:         v,
			CreatedAt:     time.Now().UTC(),
			UpdatedAt:     time.Now().UTC(),
		}

		if err := models.CreateAppTag(tx, appTag); err != nil {
			return nil, formutil.NewInternalStateError("failed to update tags")
		}
	}
	tx.Commit()

	return map[string]interface{}{"status": "success"}, nil
}
