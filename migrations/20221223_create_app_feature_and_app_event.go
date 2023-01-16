package migrations

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/kikyousky/gormigrate"
	"gitlab.shoplazza.site/common/plugin-common/xtypes"

	"gitlab.shoplazza.site/shoplaza/cobra/utils"
)

func init() {
	Migrations = append(Migrations, &gormigrate.Migration{
		ID: "20221223_create_app_feature_and_app_event",
		Migrate: func(tx *gorm.DB) error {
			return tx.Debug().Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").
				AutoMigrate(&ApplicationFeature{}).
				AutoMigrate(&DevApplicationFeature{}).
				AutoMigrate(&ApplicationEvent{}).
				AutoMigrate(&ApplicationExtend{}).
				AutoMigrate(&DevApplicationLocale{}).
				AutoMigrate(&ApplicationLocale{}).
				Error

		},
		Rollback: func(tx *gorm.DB) error {
			return nil
		},
	})
}

const (
	devApplicationFeatureTableName = "dev_application_features"
	applicationFeatureTableName    = "application_features"
	applicationEventTableName      = "app_events"
	applicationExtendTableName     = "oauth_application_extends"
	devApplicationLocaleTableName  = "dev_application_locales"
	applicationLocaleTableName     = "application_locales"
)

type (
	ApplicationFeature struct {
		Id                  uint64    `gorm:"column:id;type:bigint(20) unsigned;primary_key" json:"id"`                                           // 主键ID
		ImageUrl            string    `gorm:"column:image_url;type:varchar(255);NOT NULL" json:"image_url"`                                       // app 收费类型
		CreatedAt           time.Time `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP;NOT NULL" json:"created_at"`               // 创建时间
		UpdatedAt           time.Time `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP;NOT NULL" json:"updated_at"`               // 更新时间
		ApplicationLocaleId string    `gorm:"column:application_locale_id;type:varbinary(16);default:0x30;NOT NULL" json:"application_locale_id"` // listing locale id
	}

	DevApplicationFeature struct {
		ID                     uint64    `gorm:"column:id;type:bigint(20) unsigned;primary_key" json:"id"`                                                      // 主键ID
		ImageUrl               string    `gorm:"column:image_url;type:varchar(255);NOT NULL" json:"image_url"`                                                  // app 收费类型
		CreatedAt              time.Time `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP;NOT NULL" json:"created_at"`                          // 创建时间
		UpdatedAt              time.Time `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP;NOT NULL" json:"updated_at"`                          // 更新时间
		DevApplicationLocaleId uint64    `gorm:"column:dev_application_locale_id;type:bigint(20) unsigned;default:0;NOT NULL" json:"dev_application_locale_id"` // locale id
	}

	ApplicationEvent struct {
		Id            xtypes.XId `gorm:"column:id;type:bigint(20) unsigned;primary_key" json:"id"`
		ApplicationId xtypes.XId `gorm:"column:application_id;type:bigint(20) unsigned; :0;NOT NULL" json:"application_id"`
		EventType     string     `gorm:"column:event_type;type:varchar(30);NOT NULL" json:"event_type"`
		Reason        string     `gorm:"column:reason;type:varchar(4096);NOT NULL" json:"reason"`
		Extra         string     `gorm:"column:extra" json:"extra"`
		EventTime     time.Time  `gorm:"column:event_time;type:datetime;default:CURRENT_TIMESTAMP;NOT NULL" json:"event_time"`
		CreatedAt     time.Time  `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP;NOT NULL" json:"created_at"`
		UpdatedAt     time.Time  `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP;NOT NULL" json:"updated_at"`
		Readed        bool       `gorm:"column:readed;type:bool;default:false"  json:"readed"`
	}

	ApplicationExtend struct {
		Id                   xtypes.XId     `gorm:"column:id;type:bigint(20) unsigned;primary_key" json:"id"`
		ApplicationId        xtypes.XId     `gorm:"column:application_id;type:bigint(20) unsigned; :0;NOT NULL" json:"application_id"`
		Handle               string         `gorm:"column:handle;type:varchar(255);NOT NULL" json:"handle"`
		IsAllowTest          uint8          `gorm:"column:is_allow_test;type:tinyint(4);default 0;NOT NULL" json:"is_allow_test"`
		IsAllowSubmit        uint8          `gorm:"column:is_allow_submit;type:tinyint(4);default 0;NOT NULL" json:"is_allow_submit"`
		NotAllowTestReason   utils.JSON     `gorm:"column:not_allow_test_reason;type:json;NOT NULL" json:"not_allow_test_reason"`
		NotAllowSubmitReason utils.JSON     `gorm:"column:not_allow_submit_reason;type:json;NOT NULL" json:"not_allow_submit_reason"`
		CreatedAt            xtypes.UTCTime `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP;NOT NULL" json:"created_at"`
		UpdatedAt            xtypes.UTCTime `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP;NOT NULL" json:"updated_at"`
	}

	DevApplicationLocale struct {
		SupportTutorialUrl string `gorm:"column:support_tutorial_url;type:varchar(255);NOT NULL" json:"support_tutorial_url"`
	}
	ApplicationLocale struct {
		SupportTutorialUrl string `gorm:"column:support_tutorial_url;type:varchar(255);NOT NULL" json:"support_tutorial_url"`
	}
)

func (s ApplicationFeature) TableName() string {
	return applicationFeatureTableName
}

func (DevApplicationFeature) TableName() string {
	return devApplicationFeatureTableName
}

func (ApplicationEvent) TableName() string {
	return applicationEventTableName
}

func (ApplicationExtend) TableName() string {
	return applicationExtendTableName
}

func (DevApplicationLocale) TableName() string {
	return devApplicationLocaleTableName
}

func (ApplicationLocale) TableName() string {
	return applicationLocaleTableName
}
