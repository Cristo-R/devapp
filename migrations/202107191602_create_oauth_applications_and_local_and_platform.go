package migrations

import (
	"github.com/jinzhu/gorm"
	"github.com/kikyousky/gormigrate"
	"gitlab.shoplazza.site/shoplaza/cobra/utils"
	gormutil "gitlab.shoplazza.site/xiabing/goat.git/gorm"
)

func init() {
	Migrations = append(Migrations, &gormigrate.Migration{
		ID: "202107131602_create_oauth_applications_and_local_a",
		Migrate: func(tx *gorm.DB) error {
			type OauthApplication struct {
				gormutil.BaseModelAutoIncrement
				Name              string `gorm:"not null;default ''" json:"name" faker:"word"`
				UID               string `gorm:"unique;not null;index:index_oauth_applications_on_uid" json:"uid"`
				Secret            string `gorm:"not null;" json:"secret"`
				RedirectUri       string `gorm:"type:text;not null;default ''" json:"redirect_uri" faker:"url"`
				Scopes            string `gorm:"type:varchar(1024);not null;default ''" json:"scopes"`
				Confidential      bool   `gorm:"not null;" json:"confidential"`
				Icon              string `gorm:"type:text" json:"icon"`
				Embbed            bool   `gorm:"not null;" json:"embbed"`
				PrivateApp        bool   `gorm:"not null;" json:"private_app"`
				StoreId           string `gorm:"type:varbinary(16)" json:"store_id"` //mysql binary limit(16)
				Email             string `gorm:"not null;default ''" json:"email" faker:"email"`
				Category          string `gorm:"not null;default 'partner'" json:"category"`
				OauthDancable     bool   `gorm:"not null;default true" json:"oauth_dancable"`
				Subscribable      bool   `gorm:"not null;default false" json:"subscribable"`
				StoreCharge       bool   `json:"store_charge"`
				NotifyUri         string `gorm:"type:text" json:"notify_uri" faker:"url"` //mysql text
				Link              string `gorm:"type:text" json:"link"`                   //mysql text
				AppUri            string `gorm:"type:text" json:"app_uri" faker:"url"`    //mysql text
				ChargeMinPlan     string `json:"charge_min_plan"`
				FreeStores        string `gorm:"type:text" json:"free_stores"` //mysql text
				WebhookApiVersion string `json:"webhook_api_version"`          //mysql text
				Listing           bool   `json:"listing"`
				UnpublishedReason string `json:"unpublished_reason"`
				Status            string `json:"status"`
			}
			type ApplicationLocale struct {
				Id            utils.UUIDBinary `gorm:"PRIMARY_KEY;type:varbinary(16)" json:"id" faker:"uuid_digit,len=16"`
				Name          string           `gorm:"not null;" json:"name"`
				Subtitle      string           `gorm:"type:text;" json:"subtitle"`
				ApplicationId uint64           `gorm:"type:bigint(20);not null;UNIQUE_INDEX:index_application_locales_on_application_id_and_locale;index:index_application_locales_on_application_id" json:"application_id"`
				Locale        string           `gorm:"not null;UNIQUE_INDEX:index_application_locales_on_application_id_and_locale" json:"locale"`
				Desc          string           `gorm:"type:mediumtext" json:"desc"`
			}
			type ApplicationPlatform struct {
				Id            utils.UUIDBinary `gorm:"PRIMARY_KEY;type:varbinary(16)" json:"id"`
				ApplicationId uint64           `gorm:"not null;type:bigint(20);UNIQUE_INDEX:index_application_platforms_on_application_id_and_platform;index:index_application_platforms_on_application_id" json:"application_id"`
				Platform      string           `gorm:"not null;UNIQUE_INDEX:index_application_platforms_on_application_id_and_platform" json:"platform"`
			}
			return tx.Debug().Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").
				AutoMigrate(&OauthApplication{}).
				AutoMigrate(&ApplicationLocale{}).
				AutoMigrate(&ApplicationPlatform{}).
				Error
		},
		Rollback: func(tx *gorm.DB) error {
			//return tx.DropTable("oauth_applications").DropTable("application_locales").DropTable("application_platforms").Error
			return nil
		},
	})
}
