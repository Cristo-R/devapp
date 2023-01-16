package migrations

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/kikyousky/gormigrate"
)

func init() {
	Migrations = append(Migrations, &gormigrate.Migration{
		ID: "20220606_create_pages",
		Migrate: func(tx *gorm.DB) error {
			type Page struct {
				Id        uint64    `gorm:"PRIMARY_KEY" json:"id"`
				Title     string    `gorm:"type:varchar(255);not null" json:"title"`
				Handle    string    `gorm:"type:varchar(255);not null" json:"handle"`
				UpdatedAt time.Time `gorm:"type:datetime" json:"updated_at"`
				CreatedAt time.Time `gorm:"type:datetime" json:"created_at"`
			}
			return tx.Debug().Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").
				AutoMigrate(&Page{}).
				Error

		},
		Rollback: func(tx *gorm.DB) error {
			return nil
		},
	})
}
