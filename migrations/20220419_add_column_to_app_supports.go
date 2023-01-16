package migrations

import (
	"github.com/jinzhu/gorm"
	"github.com/kikyousky/gormigrate"
)

func init() {
	Migrations = append(Migrations, &gormigrate.Migration{
		ID: "20220419_add_column_to_app_supports",
		Migrate: func(tx *gorm.DB) error {
			type AppSupport struct {
				SubmitSource string `gorm:"column:submit_source;type:varchar(255);" json:"submit_source"`
			}
			return tx.Debug().Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").
				AutoMigrate(&AppSupport{}).
				Error

		},
		Rollback: func(tx *gorm.DB) error {
			return nil
		},
	})
}
