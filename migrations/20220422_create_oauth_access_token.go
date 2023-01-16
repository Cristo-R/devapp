package migrations

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/kikyousky/gormigrate"
	"gitlab.shoplazza.site/shoplaza/cobra/utils"
)

func init() {
	Migrations = append(Migrations, &gormigrate.Migration{
		ID: "20220422_create_oauth_access_token",
		Migrate: func(tx *gorm.DB) error {
			type OauthAccessToken struct {
				Id                   uint64           `gorm:"PRIMARY_KEY" json:"id"`
				ResourceOwnerId      utils.UUIDBinary `gorm:"type:varbinary(16)" json:"resource_owner_id"`
				ApplicationId        uint64           `gorm:"type:bigint(20);not null" json:"application_id"`
				Token                string           `gorm:"type:varchar(255)" json:"token"`
				RefreshToken         string           `gorm:"type:varchar(255)" json:"refresh_token"`
				ExpiresIn            int              `gorm:"type:int(11)" json:"expires_in"`
				RevokedAt            *time.Time       `gorm:"type:datetime" json:"revoked_at"`
				CreatedAt            time.Time        `gorm:"type:datetime" json:"created_at"`
				Scopes               string           `gorm:"type:varchar(1024)" json:"scopes"`
				PreviousRefreshToken string           `gorm:"type:varchar(255)" json:"previous_refresh_token"`
			}
			return tx.Debug().Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").
				AutoMigrate(&OauthAccessToken{}).
				Error

		},
		Rollback: func(tx *gorm.DB) error {
			return nil
		},
	})
}
