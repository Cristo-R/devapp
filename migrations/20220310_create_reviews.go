package migrations

import (
	"database/sql"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/kikyousky/gormigrate"
)

func init() {
	Migrations = append(Migrations, &gormigrate.Migration{
		ID: "20220310_create_reviews",
		Migrate: func(tx *gorm.DB) error {
			type Reviews struct {
				Id                int64        `gorm:"column:id;type:bigint(20);primary_key" json:"id"`
				Title             string       `gorm:"column:title;type:varchar(60);NOT NULL" json:"title"`                                   // 标题
				Content           string       `gorm:"column:content;type:varchar(1000);NOT NULL" json:"content"`                             // 评论内容
				Rating            int          `gorm:"column:rating;type:int(1);default:0;NOT NULL" json:"rating"`                            // 评分
				StoreId           int64        `gorm:"column:store_id;type:bigint(20);default:0;NOT NULL DEFAULT '0'" json:"store_id"`        // 店铺id
				ReviewsStatus     string       `gorm:"column:reviews_status;type:varchar(50);NOT NULL" json:"reviews_status"`                 // 评论状态
				Creator           string       `gorm:"column:creator;type:varchar(100);NOT NULL" json:"creator"`                              // 店铺账号，创建者
				Modificator       string       `gorm:"column:modificator;type:varchar(100);NOT NULL" json:"modificator"`                      // 修改者账号
				ModificatorLocale string       `gorm:"column:modificator_locale;type:varchar(100);NOT NULL" json:"modificator_locale"`        // 修改者语言
				AppId             int64        `gorm:"column:app_id;type:bigint(20);NOT NULL" json:"app_id"`                                  // appid
				PartnerReplies    string       `gorm:"column:partner_replies;type:varchar(1000);NOT NULL" json:"partner_replies"`             // 开发者回复
				CreatedAt         time.Time    `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP;NOT NULL" json:"created_at"` // 创建时间
				UpdatedAt         time.Time    `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP;NOT NULL" json:"updated_at"` // 更新时间
				DeletedAt         sql.NullTime //`gorm:"column:deleted_at;type:timestamp;default:NULL" json:"deleted_at"`                       // 软删除时间
			}
			return tx.Debug().Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").
				AutoMigrate(&Reviews{}).
				Error

		},
		Rollback: func(tx *gorm.DB) error {
			return nil
		},
	})
}
