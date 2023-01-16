package migrations

import (
	"github.com/jinzhu/gorm"
	"github.com/kikyousky/gormigrate"
)

func init() {
	Migrations = append(Migrations, &gormigrate.Migration{
		ID: "20211123_add_column_to_app_locales",
		Migrate: func(tx *gorm.DB) error {
			type ApplicationLocale struct {
				SupportWebsiteUrl       string `gorm:"type:varchar(255)" json:"support_website_url"`
				SupportPrivacyPolicyUrl string `gorm:"type:varchar(255)" json:"support_privacy_policy_url"`
				SupportFaqUrl           string `gorm:"type:varchar(255)" json:"support_faq_url"`
				VideoUrl                string `gorm:"type:varchar(255)" json:"video_url"`
				SupportEmail            string `gorm:"type:varchar(255)" json:"support_email"`
			}
			return tx.Debug().Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").
				AutoMigrate(&ApplicationLocale{}).
				Error
		},
		Rollback: func(tx *gorm.DB) error {
			return nil
		},
	})
}
