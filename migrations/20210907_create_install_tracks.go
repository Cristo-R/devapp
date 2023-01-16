package migrations

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/kikyousky/gormigrate"
	"gitlab.shoplazza.site/shoplaza/cobra/utils"
)

func init() {
	Migrations = append(Migrations, &gormigrate.Migration{
		ID: "20210929_create_install_tracks",
		Migrate: func(tx *gorm.DB) error {
			type InstallTrack struct {
				Id            utils.UUIDBinary `gorm:"PRIMARY_KEY;type:varbinary(16)" json:"id" faker:"uuid_digit,len=16"`
				ApplicationId uint64           `gorm:"type:bigint(20);null;UNIQUE_INDEX:index_install_tracks_on_store_id_and_application_id" json:"application_id"`
				StoreId       utils.UUIDBinary `gorm:"type:varbinary(16);index:index_install_tracks_on_store_id;UNIQUE_INDEX:index_install_tracks_on_store_id_and_application_id" json:"store_id" faker:"uuid_digit,len=16"`
				InstalledAt   time.Time        `json:"installed_at"`
				CreatedAt     time.Time        `json:"created_at"`
				UpdatedAt     time.Time        `json:"updated_at"`
			}
			return tx.Debug().Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").
				AutoMigrate(&InstallTrack{}).
				Error
		},
		Rollback: func(tx *gorm.DB) error {
			return nil
		},
	})
}
