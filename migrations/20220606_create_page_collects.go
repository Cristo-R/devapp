package migrations

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/kikyousky/gormigrate"
)

func init() {
	Migrations = append(Migrations, &gormigrate.Migration{
		ID: "20220606_create_page_collects",
		Migrate: func(tx *gorm.DB) error {
			type PageCollect struct {
				Id                 uint64    `gorm:"PRIMARY_KEY" json:"id"`
				CollectionId       uint64    `gorm:"type:bigint(20);not null" json:"collection_id"`
				PageId             uint64    `gorm:"type:bigint(20);not null" json:"page_id"`
				CollectionPosition int       `gorm:"type:int(10);not null" json:"collection_position"`
				UpdatedAt          time.Time `gorm:"type:datetime" json:"updated_at"`
				CreatedAt          time.Time `gorm:"type:datetime" json:"created_at"`
			}
			return tx.Debug().Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").
				AutoMigrate(&PageCollect{}).
				Error

		},
		Rollback: func(tx *gorm.DB) error {
			return nil
		},
	})
}
