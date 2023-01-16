package forms

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	formutil "gitlab.shoplazza.site/xiabing/goat.git/form"
)

func PlatformAuth(c *gin.Context, needPermission string) error {
	userId := c.Request.Header.Get("Platform-User-Id")
	if userId == "" {
		return formutil.NewUnauthorizedError("user is not login")
	}

	permissions := c.Request.Header.Get("Platform-User-Permissions")
	permissionsMap := make(map[string]interface{})

	if err := json.Unmarshal([]byte(permissions), &permissionsMap); err != nil {
		return formutil.NewWithCodeError(403, "parse platform permissions failed")
	}

	appsPermission, ok := permissionsMap["apps"]
	if !ok {
		return formutil.NewWithCodeError(403, "parse platform permissions failed")
	}

	hasPermission, ok := appsPermission.(map[string]interface{})
	if !ok {
		return formutil.NewUnprocessableError("permission denied")
	}

	if ok, _ := hasPermission[needPermission].(bool); !ok {
		return formutil.NewWithCodeError(403, "parse platform permissions failed")
	}

	return nil
}
