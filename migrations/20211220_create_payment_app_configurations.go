package migrations

import (
	"github.com/jinzhu/gorm"
	"github.com/kikyousky/gormigrate"
	"gitlab.shoplazza.site/shoplaza/cobra/models"
)

func init() {
	Migrations = append(Migrations, &gormigrate.Migration{
		ID: "20211220_create_payment_app_configurations",
		Migrate: func(tx *gorm.DB) error {
			return tx.Debug().Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").
				AutoMigrate(&models.PaymentAppConfiguration{}).
				Error
		},
		Rollback: func(tx *gorm.DB) error {
			return nil
		},
	})
}
