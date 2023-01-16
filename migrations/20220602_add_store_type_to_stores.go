package migrations

import (
	"github.com/jinzhu/gorm"
	"github.com/kikyousky/gormigrate"
)

func init() {
	Migrations = append(Migrations, &gormigrate.Migration{
		ID: "20220602_add_store_type_to_stores",
		Migrate: func(tx *gorm.DB) error {
			type Store struct {
				StoreType int `gorm:"column:store_type;type:int(1);default:0;NOT NULL" json:"store_type"`
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
