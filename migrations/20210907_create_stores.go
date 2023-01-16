package migrations

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/kikyousky/gormigrate"
	"gitlab.shoplazza.site/shoplaza/cobra/utils"
)

func init() {
	Migrations = append(Migrations, &gormigrate.Migration{
		ID: "20210907_create_stores",
		Migrate: func(tx *gorm.DB) error {
			type Store struct {
				Id        utils.UUIDBinary `gorm:"PRIMARY_KEY;type:varbinary(16)" json:"id" faker:"uuid_digit,len=16"`
				OriginId  string           `gorm:"type:varchar(255);UNIQUE_INDEX:index_stores_on_origin_id" json:"origin_id"`
				Name      string           `gorm:"type:varchar(255)" json:"name"`
				CreatedAt time.Time        `json:"created_at"`
				UpdatedAt time.Time        `json:"updated_at"`
				Platform  string           `gorm:"type:varchar(255)" json:"platform"`
			}
			return tx.Debug().Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").
				AutoMigrate(&Store{}).
				Error
		},
		Rollback: func(tx *gorm.DB) error {
			return nil
		},
	})
}
