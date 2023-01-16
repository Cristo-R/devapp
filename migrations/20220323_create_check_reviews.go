package migrations

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/kikyousky/gormigrate"
)

func init() {
	Migrations = append(Migrations, &gormigrate.Migration{
		ID: "20220310_create_check_reviews",
		Migrate: func(tx *gorm.DB) error {
			type CheckReviews struct {
				Id                int64     `gorm:"column:id;type:bigint(20);primary_key;AUTO_INCREMENT" json:"id"`                  // 主键id
				ReviewId          int64     `gorm:"column:review_id;type:bigint(20);default:0;NOT NULL" json:"review_id"`            // 评论id
				StoreId           int64     `gorm:"column:store_id;type:bigint(20);NOT NULL DEFAULT '0'" json:"store_id"`            // 店铺id
				AppId             int64     `gorm:"column:app_id;type:bigint(20);NOT NULL" json:"app_id"`                            // Appid
				CheckResultStatus string    `gorm:"column:check_result_status;type:varchar(50);NOT NULL" json:"check_result_status"` // 状态:通过accepted驳回rejected
				RejectedReason    string    `gorm:"column:rejected_reason;type:varchar(255);NOT NULL" json:"rejected_reason"`        // 违规类型
				ReviewTitle       string    `gorm:"column:review_title;type:varchar(60);NOT NULL" json:"review_title"`               // 评论标题
				ReviewContent     string    `gorm:"column:review_content;type:varchar(255);NOT NULL" json:"review_content"`          // 评论内容
				Rating            string    `gorm:"column:rating;type:varchar(255);NOT NULL" json:"rating"`                          // 评论评分
				Note              string    `gorm:"column:note;type:varchar(1000);NOT NULL" json:"note"`                             // 备注
				Checker           string    `gorm:"column:checker;type:varchar(255);NOT NULL" json:"checker"`                        // 审核人
				SubmittedAt       time.Time `gorm:"column:submitted_at;type:datetime;default:CURRENT_TIMESTAMP;NOT NULL" json:"submitted_at"`
				CreatedAt         time.Time `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP;NOT NULL" json:"created_at"`
				UpdatedAt         time.Time `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP;NOT NULL" json:"updated_at"`
			}
			return tx.Debug().Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").
				AutoMigrate(&CheckReviews{}).
				Error

		},
		Rollback: func(tx *gorm.DB) error {
			return nil
		},
	})
}
