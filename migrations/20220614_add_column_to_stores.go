package migrations

import (
	"github.com/jinzhu/gorm"
	"github.com/kikyousky/gormigrate"
)

func init() {
	Migrations = append(Migrations, &gormigrate.Migration{
		ID: "20220614_add_column_to_stores",
		Migrate: func(tx *gorm.DB) error {
			type Store struct {
				Guides string `gorm:"column:guides;type:varchar(255);" json:"guides"`
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
