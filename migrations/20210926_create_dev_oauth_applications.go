package migrations

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/kikyousky/gormigrate"
)

func init() {
	Migrations = append(Migrations, &gormigrate.Migration{
		ID: "20210926_create_dev_oauth_application",
		Migrate: func(tx *gorm.DB) error {
			type DevOauthApplication struct {
				ApplicationID     uint64    `gorm:"AUTO_INCREMENT" json:"application_id" faker:"boundary_start=1000000, boundary_end=9999999"`
				Name              string    `gorm:"not null;default ''" json:"name" faker:"word"`
				RedirectUri       string    `gorm:"type:text;not null;default ''" json:"redirect_uri" faker:"url"`
				Icon              string    `gorm:"type:text" json:"icon"`
				Email             string    `gorm:"not null;default ''" json:"email" faker:"email"`
				AppUri            string    `gorm:"type:text" json:"app_uri" faker:"url"` //mysql text
				WebhookApiVersion string    `json:"webhook_api_version"`                  //mysql text
				AdditionalInfo    string    `gorm:"type:varchar(4096)" json:"additional_info"`
				Status            string    `json:"status"`
				RejectedReason    string    `gorm:"type:text" json:"rejected_reason"`
				CreatedAt         time.Time `gorm:"type:datetime" json:"created_at"`
				UpdatedAt         time.Time `gorm:"type:datetime" json:"-"`
			}
			return tx.Debug().Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").
				AutoMigrate(&DevOauthApplication{}).
				Error
		},
		Rollback: func(tx *gorm.DB) error {
			return nil
		},
	})
}
