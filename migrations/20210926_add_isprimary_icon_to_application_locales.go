package migrations

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/kikyousky/gormigrate"
)

func init() {
	Migrations = append(Migrations, &gormigrate.Migration{
		ID: "20210926_add_columns_application_locales",
		Migrate: func(tx *gorm.DB) error {
			type ApplicationLocale struct {
				IsPrimary bool
				Icon      string
				CreatedAt time.Time `gorm:"type:datetime" json:"created_at"`
				UpdatedAt time.Time `gorm:"type:datetime" json:"-"`
			}
			return tx.Debug().Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").
				AutoMigrate(&ApplicationLocale{}).
				Error
		},
		Rollback: func(tx *gorm.DB) error {
			return nil
		},
	})
}
