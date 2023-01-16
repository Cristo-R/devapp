package migrations

import (
	"github.com/jinzhu/gorm"
	"github.com/kikyousky/gormigrate"
)

func init() {
	Migrations = append(Migrations, &gormigrate.Migration{
		ID: "20211026_add_column_to_applications",
		Migrate: func(tx *gorm.DB) error {
			type OauthApplication struct {
				AdditionalInfo string `gorm:"type:varchar(4096)" json:"additional_info"`
				RejectedReason string `gorm:"type:varchar(4096)" json:"rejected_reason"`
			}

			return tx.Debug().Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").
				AutoMigrate(&OauthApplication{}).
				Error
		},
		Rollback: func(tx *gorm.DB) error {
			return nil
		},
	})
}
