package migrations

import (
	"github.com/jinzhu/gorm"
	"github.com/kikyousky/gormigrate"
	"gitlab.shoplazza.site/shoplaza/cobra/utils"
)

func init() {
	Migrations = append(Migrations, &gormigrate.Migration{
		ID: "20220602_add_new_store_id_to_check_reviews",
		Migrate: func(tx *gorm.DB) error {
			type CheckReviews struct {
				NewStoreId utils.UUIDBinary `gorm:"column:new_store_id;type:varbinary(16);NOT NULL" json:"new_store_id"` //新的店铺ID
			}
			return tx.Debug().Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").
				AutoMigrate(&CheckReviews{}).
				Error
		},
		Rollback: func(tx *gorm.DB) error {
			return nil
		},
	})
}
