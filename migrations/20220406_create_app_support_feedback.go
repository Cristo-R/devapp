package migrations

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/kikyousky/gormigrate"
)

func init() {
	Migrations = append(Migrations, &gormigrate.Migration{
		ID: "20220406_create_app_support_feedback",
		Migrate: func(tx *gorm.DB) error {
			type AppSupportFeedback struct {
				Id               uint64    `gorm:"column:id;type:bigint(20);primary_key" json:"id"`                      // 主键id
				AppSupportId     uint64    `gorm:"column:app_support_id;type:bigint(20);NOT NULL" json:"app_support_id"` // 店铺id
				AdditionFeedback string    `gorm:"column:addition_feedback;type:varchar(4096)" json:"addition_feedback"`
				Rating           uint      `gorm:"column:rating;type:int(1)" json:"rating"`
				ReplyByDeveloper bool      `gorm:"column:reply_by_developer;type:tinyint(1)" json:"reply_by_developer"`
				CreatedAt        time.Time `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP;NOT NULL" json:"created_at"`
				UpdatedAt        time.Time `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP;NOT NULL" json:"updated_at"`
			}
			return tx.Debug().Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").
				AutoMigrate(&AppSupportFeedback{}).
				Error

		},
		Rollback: func(tx *gorm.DB) error {
			return nil
		},
	})
}
