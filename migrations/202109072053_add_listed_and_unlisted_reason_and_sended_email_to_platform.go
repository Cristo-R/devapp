package migrations

import (
	"github.com/jinzhu/gorm"
	"github.com/kikyousky/gormigrate"
)

func init() {
	Migrations = append(Migrations, &gormigrate.Migration{
		ID: "202109072053_add_listed_to_platform",
		Migrate: func(tx *gorm.DB) error {
			type ApplicationPlatform struct {
				Listed         bool   `gorm:"type:tinyint(1);default 0"`
				UnlistedReason string `gorm:"default ''"`
				SendedEmail    bool   `gorm:"type:tinyint(1);default 0"`
			}

			return tx.Debug().Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").
				AutoMigrate(&ApplicationPlatform{}).
				Error
		},
		Rollback: func(tx *gorm.DB) error {
			return nil
		},
	})
}
