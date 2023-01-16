package migrations

import (
	"github.com/jinzhu/gorm"
	"github.com/kikyousky/gormigrate"
)

func init() {
	Migrations = append(Migrations, &gormigrate.Migration{
		ID: "20220614_add_regions_to_oauth_applications",
		Migrate: func(tx *gorm.DB) error {
			type OauthApplication struct {
				Regions string `gorm:"type:varchar(1024);not null;" json:"regions"`
			}
			return tx.Debug().Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").
				AutoMigrate(&OauthApplication{}).
				Error
		},
		Rollback: func(tx *gorm.DB) error {
			return nil
		},
	})
}
