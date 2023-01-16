package migrations

import (
	"github.com/jinzhu/gorm"
	"github.com/kikyousky/gormigrate"
)

func init() {
	Migrations = append(Migrations, &gormigrate.Migration{
		ID: "20220608_add_column_collections",
		Migrate: func(tx *gorm.DB) error {
			type Collection struct {
				SortOrder string `gorm:"column:sort_order;type:varchar(255);" json:"sort_order"`
			}
			return tx.Debug().Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").
				AutoMigrate(&Collection{}).
				Error

		},
		Rollback: func(tx *gorm.DB) error {
			return nil
		},
	})
}
