package migrations

import (
	"github.com/jinzhu/gorm"
	"github.com/kikyousky/gormigrate"
	"time"
)

func init() {
	Migrations = append(Migrations, &gormigrate.Migration{
		ID: "202108181640_create_payment_extensions_table",
		Migrate: func(tx *gorm.DB) error {
			type PaymentExtension struct {
				AppId             uint64    `gorm:"column:app_id;primary_key"`
				PaymentSessionUrl string    `gorm:"column:payment_session_url;type:varchar(255);not null;default:''"`
				RefundSessionUrl  string    `gorm:"column:refund_session_url;type:varchar(255);not null;default:''"`
				CaptureSessionUrl string    `gorm:"column:capture_session_url;type:varchar(255);not null;default:''"`
				VoidSessionUrl    string    `gorm:"column:void_session_url;type:varchar(255);not null;default:''"`
				Methods           string    `gorm:"column:methods;type:varchar(255);not null;default:''"`
				Countries         string    `gorm:"column:countries;type:varchar(255);not null;default:''"`
				MerchantAdminName string    `gorm:"column:merchant_admin_name;type:varchar(50);not null;default:''"`
				CheckoutName      string    `gorm:"column:checkout_name;type:varchar(255);not null;default:''"`
				TestMode          bool      `gorm:"column:test_mode;default:0"`
				AppIcon           string    `gorm:"column:app_icon;type:varchar(255);not null;default:''"`
				AppIconWhite      string    `gorm:"column:app_icon_white;type:varchar(255);not null;default:''"`
				OfficialWebsite   string    `gorm:"column:official_website;type:varchar(255);not null;default:''"`
				BackendWebsite    string    `gorm:"column:backend_website;type:varchar(255);not null;default:''"`
				CreatedAt         time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP"`
				UpdatedAt         time.Time `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
			}

			return tx.Table("payment_extensions").Debug().Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").
				AutoMigrate(&PaymentExtension{}).
				Error
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Debug().DropTableIfExists("payment_extensions").Error
		},
	})
}
