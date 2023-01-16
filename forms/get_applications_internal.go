package forms

import (
	//"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"gitlab.shoplazza.site/shoplaza/cobra/config"

	formutil "gitlab.shoplazza.site/xiabing/goat.git/form"
	//"gitlab.shoplazza.site/xiabing/goat.git/middlewares"
	//"strings"
	//"time"
)

type GetAppFormForInternal struct {
	UID string `form:"uid" json:"uid"`
}

func (form *GetAppFormForInternal) Do(c *gin.Context) (interface{}, error) {
	result := &struct {
		Uid    string `form:"uid" json:"uid"`
		Secret string `form:"secret" json:"secret"`
	}{}
	if err := config.DB.Table("oauth_applications").Select("uid,secret").Where("uid = ?", form.UID).Take(&result).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, formutil.NewNotFoundError("app_not_exits")
		}
		return nil, err
	}

	return result, nil
}
