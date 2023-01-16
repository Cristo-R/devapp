package admin

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.shoplazza.site/shoplaza/cobra/config"
	"gitlab.shoplazza.site/shoplaza/cobra/models"
	"gitlab.shoplazza.site/shoplaza/cobra/service"
	"gitlab.shoplazza.site/shoplaza/cobra/utils"
	formutil "gitlab.shoplazza.site/xiabing/goat.git/form"
	"gitlab.shoplazza.site/xiabing/goat.git/i18n"
	"gitlab.shoplazza.site/xiabing/goat.git/middlewares"
)

type GetGuidesForm struct{}

func (form *GetGuidesForm) Do(c *gin.Context) (interface{}, error) {
	appctx, ok := c.MustGet("context").(*middlewares.Context)
	if !ok {
		return nil, formutil.NewInternalStateError(i18n.T("no_context", appctx.Locale))
	}

	// 根据 store id 查询 stores表
	storeId := strconv.Itoa(int(appctx.StoreID))
	store, err := models.GetStore(config.DB, storeId)
	if err != nil {
		return nil, err
	}
	if store != nil && store.Guides != nil {
		return map[string]interface{}{"guides": store.Guides}, nil
	}

	// 调取店铺组接口查询店铺信息
	storeInfo, err := service.GetStoreFromTotoro(storeId)
	if err != nil {
		return nil, err
	}
	if storeInfo == nil {
		return nil, formutil.NewNotFoundError(i18n.T("get_account_store_error", appctx.Locale))
	}

	// 如果店铺创建时间在指定时间之前，说明是旧店铺
	guides := make(models.Guides, 0)
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", config.Cfg.StoreCreatedAt, time.Local)
	if storeInfo.CreatedAt.Before(t) {
		guides = append(guides, models.Guide{Name: models.GUIDE_TYPE_NEW_STORE_ENTRY, Status: models.GUIDE_STATUS_INIT})
		guides = append(guides, models.Guide{Name: models.GUIDE_TYPE_RECOMMEND_APPS, Status: models.GUIDE_STATUS_INIT})
	} else {
		guides = append(guides, models.Guide{Name: models.GUIDE_TYPE_NEW_STORE_ENTRY, Status: models.GUIDE_STATUS_FINISHED})
		guides = append(guides, models.Guide{Name: models.GUIDE_TYPE_RECOMMEND_APPS, Status: models.GUIDE_STATUS_FINISHED})
	}

	if store == nil {
		s := models.Store{
			Id:        utils.NewUUIDBinary(),
			OriginId:  storeId,
			Name:      storeInfo.Name,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Guides:    &guides,
		}
		if err := models.CreateStore(config.DB, s); err != nil {
			return nil, err
		}
	}
	// 更新店铺
	updateParam := make(map[string]interface{})
	updateParam["guides"] = guides
	if err := models.UpdateStoreById(config.DB, storeId, updateParam); err != nil {
		return nil, err
	}

	return map[string]interface{}{"guides": guides}, nil
}
