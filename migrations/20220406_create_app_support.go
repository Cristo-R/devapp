package migrations

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/kikyousky/gormigrate"
)

func init() {
	Migrations = append(Migrations, &gormigrate.Migration{
		ID: "20220406_create_app_support",
		Migrate: func(tx *gorm.DB) error {
			type AppSupport struct {
				Id                     uint64     `gorm:"column:id;type:bigint(20);primary_key" json:"id"`                      // 主键id
				StoreId                uint64     `gorm:"column:store_id;type:bigint(20);NOT NULL" json:"store_id"`             // 店铺id
				ApplicationId          int64      `gorm:"column:application_id;type:bigint(20);NOT NULL" json:"application_id"` // Appid
				ApplicationListingName string     `gorm:"column:application_listing_name;type:varchar(255);NOT NULL" json:"application_listing_name"`
				ApplicationListingIcon string     `gorm:"column:application_listing_icon;type:varchar(255);NOT NULL" json:"application_listing_icon"`
				IssueDescription       string     `gorm:"column:issue_description;type:varchar(4096);NOT NULL" json:"issue_description"`
				ContactEmail           string     `gorm:"column:contact_email;type:varchar(255);NOT NULL" json:"contact_email"`
				StoreName              string     `gorm:"column:store_name;type:varchar(255)" json:"store_name"`
				StoreDomain            string     `gorm:"column:store_domain;type:varchar(255)" json:"store_domain"`
				StoreLocale            string     `gorm:"column:store_locale;type:varchar(20)" json:"store_locale"`
				HasSendFeedbackEmail   bool       `gorm:"column:has_send_feedback_email;type:tinyint(1)" json:"has_send_feedback_email"`
				FeedbackAccessToken    string     `gorm:"column:feedback_access_token;type:varchar(255)" json:"feedback_access_token"`
				TokenExpiredAt         *time.Time `gorm:"column:token_expired_at;type:datetime;default:CURRENT_TIMESTAMP" json:"token_expired_at"`
				CreatedAt              time.Time  `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP;NOT NULL" json:"created_at"`
				UpdatedAt              time.Time  `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP;NOT NULL" json:"updated_at"`
			}
			return tx.Debug().Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").
				AutoMigrate(&AppSupport{}).
				Error

		},
		Rollback: func(tx *gorm.DB) error {
			return nil
		},
	})
}
