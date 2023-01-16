package migrations

import (
	"github.com/jinzhu/gorm"
	"github.com/kikyousky/gormigrate"
	gormutil "gitlab.shoplazza.site/xiabing/goat.git/gorm"
)

func init() {
	Migrations = append(Migrations, &gormigrate.Migration{
		ID: "20210907_create_app_tags",
		Migrate: func(tx *gorm.DB) error {
			type AppTag struct {
				gormutil.BaseModelAutoIncrement
				ApplicationId uint64 `gorm:"type:bigint(20);not null;UNIQUE_INDEX:uk_appid_and_tagid" json:"application_id"`
				TagId         uint64 `gorm:"type:bigint(20);not null;UNIQUE_INDEX:uk_appid_and_tagid" json:"tag_id"`
			}
			return tx.Debug().Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").
				AutoMigrate(&AppTag{}).
				Error
		},
		Rollback: func(tx *gorm.DB) error {
			return nil
		},
	})
}
