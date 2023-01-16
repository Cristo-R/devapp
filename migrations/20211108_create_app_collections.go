package migrations

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/kikyousky/gormigrate"
)

func init() {
	Migrations = append(Migrations, &gormigrate.Migration{
		ID: "20211108_create_app_collections",
		Migrate: func(tx *gorm.DB) error {
			type ApplicationCollection struct {
				ID            uint64    `gorm:"AUTO_INCREMENT" json:"id" faker:"boundary_start=1000000, boundary_end=9999999"`
				ApplicationId uint64    `gorm:"type:bigint(20);not null" json:"application_id"`
				CollectionId  uint64    `gorm:"type:int(10);not null;index:idx_application_collections_on_collection_id_and_position" json:"collection_id"`
				Position      uint64    `gorm:"type:int(10);not null;index:idx_application_collections_on_collection_id_and_position" json:"position"`
				CreatedAt     time.Time `gorm:"type:datetime" json:"created_at"`
				UpdatedAt     time.Time `gorm:"type:datetime" json:"-"`
			}
			return tx.Debug().Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").
				AutoMigrate(&ApplicationCollection{}).
				Error
		},
		Rollback: func(tx *gorm.DB) error {
			return nil
		},
	})
}
