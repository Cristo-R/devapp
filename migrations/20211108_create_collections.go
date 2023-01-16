package migrations

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/kikyousky/gormigrate"
)

func init() {
	Migrations = append(Migrations, &gormigrate.Migration{
		ID: "20211108_create_collections",
		Migrate: func(tx *gorm.DB) error {
			type Collection struct {
				ID           uint64    `gorm:"AUTO_INCREMENT" json:"id" faker:"boundary_start=1000000, boundary_end=9999999"`
				IsShowAtNgv  bool      `gorm:"type:tinyint(1);not null;default 0" json:"is_show_at_ngv"`
				IsShowAtHome bool      `gorm:"type:tinyint(1);not null;default 0" json:"is_show_at_home"`
				Position     uint64    `gorm:"type:int(10);not null" json:"position"`
				CreatedAt    time.Time `gorm:"type:datetime" json:"created_at"`
				UpdatedAt    time.Time `gorm:"type:datetime" json:"-"`
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
