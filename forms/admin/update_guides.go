package admin

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.shoplazza.site/shoplaza/cobra/config"
	"gitlab.shoplazza.site/shoplaza/cobra/models"
	formutil "gitlab.shoplazza.site/xiabing/goat.git/form"
	"gitlab.shoplazza.site/xiabing/goat.git/i18n"
	"gitlab.shoplazza.site/xiabing/goat.git/middlewares"
)

type UpdateGuidesForm struct {
	Key string `form:"key" json:"key" binding:"required"`
}

func (form *UpdateGuidesForm) Do(c *gin.Context) (interface{}, error) {
	appctx, ok := c.MustGet("context").(*middlewares.Context)
	if !ok {
		return nil, formutil.NewInternalStateError(i18n.T("no_context", appctx.Locale))
	}

	storeId := strconv.Itoa(int(appctx.StoreID))
	store, err := models.GetStore(config.DB, storeId)
	if err != nil {
		return nil, err
	}
	if store == nil {
		return nil, formutil.NewNotFoundError(i18n.T("store_no_exist", appctx.Locale))
	}
	if store.Guides == nil {
		return nil, formutil.NewNotFoundError(i18n.T("guides_no_existed", appctx.Locale))
	}

	guides := *store.Guides
	for k, v := range guides {
		if v.Name == form.Key {
			guides[k].Status = models.GUIDE_STATUS_FINISHED
		}
	}

	// 更新店铺
	modifyField := make(map[string]interface{})
	modifyField["guides"] = guides
	if err := models.UpdateStoreById(config.DB, storeId, modifyField); err != nil {
		return nil, err
	}

	return map[string]interface{}{"guides": guides}, nil
}
