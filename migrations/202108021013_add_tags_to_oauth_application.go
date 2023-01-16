package migrations

import (
	"github.com/jinzhu/gorm"
	"github.com/kikyousky/gormigrate"
)

func init() {
	Migrations = append(Migrations, &gormigrate.Migration{
		ID: "202108021013_add_tags_to_oauth_application",
		Migrate: func(tx *gorm.DB) error {
			type OauthApplication struct {
				Tags string `gorm:"type:varchar(1024)"`
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
