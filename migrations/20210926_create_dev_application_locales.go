package migrations

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/kikyousky/gormigrate"
)

func init() {
	Migrations = append(Migrations, &gormigrate.Migration{
		ID: "20210926_create_dev_application_locales",
		Migrate: func(tx *gorm.DB) error {
			type DevApplicationLocale struct {
				ID            uint64 `gorm:"AUTO_INCREMENT" json:"id" faker:"boundary_start=1000000, boundary_end=9999999"`
				Name          string `gorm:"not null;" json:"name"`
				Subtitle      string `gorm:"type:text;" json:"subtitle"`
				ApplicationId uint64 `gorm:"type:bigint(20);not null;UNIQUE_INDEX:idx_dev_application_locales_on_application_id_and_locale;index:idx_dev_application_locales_on_application_id" json:"application_id"`
				Locale        string `gorm:"not null;UNIQUE_INDEX:idx_dev_application_locales_on_application_id_and_locale" json:"locale"`
				Desc          string `gorm:"type:mediumtext" json:"desc"`
				IsPrimary     bool
				Icon          string
				Status        string
				CreatedAt     time.Time `gorm:"type:datetime" json:"created_at"`
				UpdatedAt     time.Time `gorm:"type:datetime" json:"-"`
			}
			return tx.Debug().Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").
				AutoMigrate(&DevApplicationLocale{}).
				Error
		},
		Rollback: func(tx *gorm.DB) error {
			return nil
		},
	})
}
