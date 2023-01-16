package migrations

import (
	"github.com/jinzhu/gorm"
	"github.com/kikyousky/gormigrate"
)

func init() {
	Migrations = append(Migrations, &gormigrate.Migration{
		ID: "20220308_add_column_to_oauth_applications",
		Migrate: func(tx *gorm.DB) error {
			type OauthApplication struct {
				DevelopByShoplazza bool `gorm:"type:tinyint(1);default:'0'" json:"develop_by_shoplazza"`
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
