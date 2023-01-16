package migrations

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/kikyousky/gormigrate"
)

func init() {
	Migrations = append(Migrations, &gormigrate.Migration{
		ID: "20211123_create_app_key_benefits",
		Migrate: func(tx *gorm.DB) error {
			type ApplicationKeyBenefit struct {
				ID                  uint64    `gorm:"AUTO_INCREMENT" json:"id" faker:"boundary_start=1000000, boundary_end=9999999"`
				CreatedAt           time.Time `gorm:"type:datetime" json:"created_at"`
				UpdatedAt           time.Time `gorm:"type:datetime" json:"-"`
				ImageUrl            string    `gorm:"type:varchar(255)" json:"image_url"`
				Title               string    `gorm:"type:varchar(64)" json:"title"`
				Description         string    `gorm:"type:varchar(150)" json:"description"`
				ApplicationLocaleId string    `gorm:"type:varbinary(16);not null;index:idx_application_key_benefits_on_application_locale_id" json:"application_locale_id"`
			}
			return tx.Debug().Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").
				AutoMigrate(&ApplicationKeyBenefit{}).
				Error
		},
		Rollback: func(tx *gorm.DB) error {
			return nil
		},
	})
}
