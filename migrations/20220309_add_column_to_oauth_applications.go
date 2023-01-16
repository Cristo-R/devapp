package migrations

import (
	"github.com/jinzhu/gorm"
	"github.com/kikyousky/gormigrate"
)

func init() {
	Migrations = append(Migrations, &gormigrate.Migration{
		ID: "20220309_add_column_to_oauth_applications",
		Migrate: func(tx *gorm.DB) error {
			type OauthApplication struct {
				UserId    string `gorm:"type:VARCHAR(100);not null;" json:"user_id"`
				PartnerId uint64 `gorm:"type:bigint(20);not null;" json:"partner_id"`
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
