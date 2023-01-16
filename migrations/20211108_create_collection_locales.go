package migrations

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/kikyousky/gormigrate"
)

func init() {
	Migrations = append(Migrations, &gormigrate.Migration{
		ID: "20211108_create_collection_locales",
		Migrate: func(tx *gorm.DB) error {
			type CollectionLocale struct {
				ID           uint64    `gorm:"AUTO_INCREMENT" json:"id" faker:"boundary_start=1000000, boundary_end=9999999"`
				Name         string    `gorm:"type:varchar(255);not null;default ''" json:"name"`
				Description  string    `gorm:"type:varchar(1024);not null;default ''" json:"description"`
				Locale       string    `gorm:"type:varchar(20);not null;" json:"locale"`
				CollectionId uint64    `gorm:"type:int(10);not null;index:idx_collection_locales_on_collection_id" json:"collection_id"`
				CreatedAt    time.Time `gorm:"type:datetime" json:"created_at"`
				UpdatedAt    time.Time `gorm:"type:datetime" json:"-"`
			}
			return tx.Debug().Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").
				AutoMigrate(&CollectionLocale{}).
				Error
		},
		Rollback: func(tx *gorm.DB) error {
			return nil
		},
	})
}
