package migrations

import (
	"github.com/jinzhu/gorm"
	"github.com/kikyousky/gormigrate"
	gormutil "gitlab.shoplazza.site/xiabing/goat.git/gorm"
)

func init() {
	Migrations = append(Migrations, &gormigrate.Migration{
		ID: "20210907_create_tags",
		Migrate: func(tx *gorm.DB) error {
			type Tag struct {
				gormutil.BaseModelAutoIncrement
				NameZh   string `gorm:"type:varchar(255)" json:"name_zh"`
				NameEn   string `gorm:"type:varchar(255)" json:"name_en"`
				TagLevel string `gorm:"type:varchar(255);not null" json:"tag_level"`
				ParentId uint64 `gorm:"type:bigint(20);not null" json:"parent_id"`
			}
			return tx.Debug().Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").
				AutoMigrate(&Tag{}).
				Error
		},
		Rollback: func(tx *gorm.DB) error {
			return nil
		},
	})
}
