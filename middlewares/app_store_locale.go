package middlewares

import (
	"github.com/gin-gonic/gin"
)

func AppStoreLocaleContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, _ := c.Cookie("app_store_locale")
		c.Set("app_store_locale", cookie)
	}
}
