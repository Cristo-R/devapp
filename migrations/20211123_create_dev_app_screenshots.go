package migrations

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/kikyousky/gormigrate"
)

func init() {
	Migrations = append(Migrations, &gormigrate.Migration{
		ID: "20211123_create_dev_app_screenshots",
		Migrate: func(tx *gorm.DB) error {
			type DevApplicationScreenshot struct {
				ID                     uint64    `gorm:"AUTO_INCREMENT" json:"id" faker:"boundary_start=1000000, boundary_end=9999999"`
				CreatedAt              time.Time `gorm:"type:datetime" json:"created_at"`
				UpdatedAt              time.Time `gorm:"type:datetime" json:"-"`
				ImageUrl               string    `gorm:"type:varchar(255)" json:"image_url"`
				Description            string    `gorm:"type:varchar(100)" json:"description"`
				DevApplicationLocaleId string    `gorm:"type:bigint(20);not null;index:idx_dev_application_screenshots_on_dev_application_locale_id" json:"dev_application_locale_id"`
			}
			return tx.Debug().Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").
				AutoMigrate(&DevApplicationScreenshot{}).
				Error
		},
		Rollback: func(tx *gorm.DB) error {
			return nil
		},
	})
}
