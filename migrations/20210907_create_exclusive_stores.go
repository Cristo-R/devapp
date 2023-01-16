package migrations

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/kikyousky/gormigrate"
	"gitlab.shoplazza.site/shoplaza/cobra/utils"
)

func init() {
	Migrations = append(Migrations, &gormigrate.Migration{
		ID: "20210907_create_exclusive_stores",
		Migrate: func(tx *gorm.DB) error {
			type ExclusiveStore struct {
				Id            utils.UUIDBinary `gorm:"PRIMARY_KEY;type:varbinary(16)" json:"id" faker:"uuid_digit,len=16"`
				ApplicationId uint64           `gorm:"type:bigint(20);not null;UNIQUE_INDEX:index_exclusive_stores_on_store_id_and_application_id" json:"application_id"`
				StoreId       utils.UUIDBinary `gorm:"type:varbinary(16);index:index_exclusive_stores_on_store_id;UNIQUE_INDEX:index_exclusive_stores_on_store_id_and_application_id" json:"store_id" faker:"uuid_digit,len=16"`
				CreatedAt     time.Time        `json:"created_at"`
				UpdatedAt     time.Time        `json:"updated_at"`
			}
			return tx.Debug().Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").
				AutoMigrate(&ExclusiveStore{}).
				Error
		},
		Rollback: func(tx *gorm.DB) error {
			return nil
		},
	})
}
