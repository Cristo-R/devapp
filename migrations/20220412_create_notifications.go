package migrations

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/kikyousky/gormigrate"
)

func init() {
	Migrations = append(Migrations, &gormigrate.Migration{
		ID: "20220412_create_notifications",
		Migrate: func(tx *gorm.DB) error {
			type Notifications struct {
				Id                 int64     `gorm:"column:id;type:bigint(20);primary_key;default:0" json:"id"`                            // 主键id
				StoreId            int64     `gorm:"column:store_id;type:bigint(20);default:0;NOT NULL" json:"store_id"`                   // 店铺id
				AppId              int64     `gorm:"column:app_id;type:bigint(20);default:0;NOT NULL" json:"app_id"`                       // 应用id
				NotificationsEvent string    `gorm:"column:notifications_event;type:varchar(100);NOT NULL" json:"notifications_event"`     // 消息事件
				CreatedAt          time.Time `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP;NOT NULL" json:"created_at"` // 创建时间
				UpdatedAt          time.Time `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP;NOT NULL" json:"updated_at"` // 更新时间
			}
			return tx.Debug().Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").
				AutoMigrate(&Notifications{}).
				Error

		},
		Rollback: func(tx *gorm.DB) error {
			return nil
		},
	})
}
