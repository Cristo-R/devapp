package migrations

import (
	"github.com/jinzhu/gorm"
	"github.com/kikyousky/gormigrate"
	"time"
)

func init() {
	Migrations = append(Migrations, &gormigrate.Migration{
		ID: "202108271955_create_account_alias_table",
		Migrate: func(tx *gorm.DB) error {
			type AccountAlias struct {
				ID        uint64    `gorm:"column:id;primary_key"`
				AppId     uint64    `gorm:"column:app_id;type:bigint unsigned;not null;default:0"`
				AccountId string    `gorm:"column:account_id;type:varchar(255);not null;default:''"`
				SignKey   string    `gorm:"column:sign_key;type:varchar(255);not null;default:''"`
				KeyA      string    `gorm:"column:key_a;type:varchar(255);not null;default:''"`
				KeyB      string    `gorm:"column:key_b;type:varchar(255);not null;default:''"`
				Locale    string    `gorm:"column:locale;type:varchar(50);not null;default:''"`
				CreatedAt time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP"`
				UpdatedAt time.Time `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
			}

			return tx.Table("account_alias").Debug().Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").
				AutoMigrate(&AccountAlias{}).
				AddUniqueIndex("idx_account_alias_app_id_locale", "app_id", "locale").
				Error
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Debug().DropTableIfExists("account_alias").Error
		},
	})
}
