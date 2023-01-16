package middlewares

import (
	"encoding/json"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gitlab.shoplazza.site/shoplaza/cobra/config"
	"gitlab.shoplazza.site/shoplaza/cobra/models"
	"gitlab.shoplazza.site/shoplaza/cobra/utils"
	"gitlab.shoplazza.site/xiabing/goat.git/i18n"
)

type Store struct {
	ID            uint64 `json:"id"`
	Name          string `json:"name"`
	Slug          string `json:"slug"`
	PrimaryDomain string `json:"primary_domain"`
}

func CheckLoginStoreId() gin.HandlerFunc {
	return func(c *gin.Context) {
		appStoreLocale, ok := c.MustGet("app_store_locale").(string)
		if !ok {
			appStoreLocale = models.LOCALE_EN_US
		}

		loginUserId := c.Request.Header.Get("Login-User-Id")

		storeId := ""
		if c.Param("id") != "" {
			storeId = c.Param("id")
		}
		if c.Request.Header.Get("Login-Store-ID") != "" {
			storeId = c.Request.Header.Get("Login-Store-ID")
		}

		loginStoreId, err := strconv.ParseUint(storeId, 10, 64)
		if err != nil {
			log.WithError(err)
			c.String(422, i18n.T("store_id_check_type", appStoreLocale))
			c.Abort()
			return
		}

		var stores []Store
		httpPool := utils.NewHttpPool()

		body, err := httpPool.HttpRequest("GET", config.Cfg.StoreService+"/api/user/"+loginUserId+"/stores", nil)
		if err != nil {
			log.WithError(err)
			c.String(500, i18n.T("user_obtain_error", appStoreLocale))
			c.Abort()
			return
		}

		err = json.Unmarshal(body, &stores)
		if err != nil {
			log.WithError(err)
			c.String(500, i18n.T("user_obtain_error", appStoreLocale))
			c.Abort()
			return
		}

		var i int
		for i = 0; i < len(stores); i++ {
			if stores[i].ID == loginStoreId {
				c.Set("Login-Store-ID", stores[i].ID)
				c.Set("Login-Store-Domain", stores[i].PrimaryDomain)
				c.Set("Login-Store-Name", stores[i].Name)
				break
			}
		}
		if i == len(stores) {
			c.String(422, i18n.T("check_user_store", appStoreLocale))
			c.Abort()
			return
		}
	}
}
