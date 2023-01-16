package migrations

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/kikyousky/gormigrate"
)

func init() {
	Migrations = append(Migrations, &gormigrate.Migration{
		ID: "20220310_create_application_pricing",
		Migrate: func(tx *gorm.DB) error {
			return tx.Debug().Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").
				AutoMigrate(&DevApplicationPricing{}).
				AutoMigrate(&DevApplicationPricingPlan{}).
				AutoMigrate(&ApplicationPricing{}).
				AutoMigrate(&ApplicationPricingPlan{}).
				Error

		},
		Rollback: func(tx *gorm.DB) error {
			return nil
		},
	})
}

// 应用收费审核主表
type DevApplicationPricing struct {
	ID                     uint64    `gorm:"column:id;type:bigint(20) unsigned;primary_key" json:"id"`                                                      // 主键ID
	ChargeType             string    `gorm:"column:charge_type;type:varchar(30);NOT NULL" json:"charge_type"`                                               // app 收费类型
	CreatedAt              time.Time `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP;NOT NULL" json:"created_at"`                          // 创建时间
	UpdatedAt              time.Time `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP;NOT NULL" json:"updated_at"`                          // 更新时间
	DevApplicationLocaleID uint64    `gorm:"column:dev_application_locale_id;type:bigint(20) unsigned;default:0;NOT NULL" json:"dev_application_locale_id"` // listing locale id
	ApplicationID          uint64    `gorm:"column:application_id;type:bigint(20) unsigned;default:0;NOT NULL" json:"application_id"`                       // 应用ID
}

func (s DevApplicationPricing) TableName() string {
	return "dev_application_pricing"
}

// 收费套餐审核表
type DevApplicationPricingPlan struct {
	ID                     uint64    `gorm:"column:id;type:bigint(20) unsigned;primary_key" json:"id"`                                                      // 主键ID
	PlanName               string    `gorm:"column:plan_name;type:varchar(100);NOT NULL" json:"plan_name"`                                                  // 套餐名字
	Features               string    `gorm:"column:features;type:varchar(255);NOT NULL" json:"features"`                                                    // 套餐特点
	Price                  string    `gorm:"column:price;type:decimal(19,2);default:0.00;NOT NULL" json:"price"`                                            // 月度套餐价格
	PriceType              string    `gorm:"column:price_type;type:varchar(20);default:'monthly';NOT NULL" json:"price_type"`                               // 价格类型： free  monthly  yearly  one-time
	TrailDays              int       `gorm:"column:trail_days;type:int(11);default:0;NOT NULL" json:"trail_days"`                                           // 试用期天数
	ParentPlanID           uint64    `gorm:"column:parent_plan_id;type:bigint(20) unsigned;default:0;NOT NULL" json:"parent_plan_id"`                       // 非primay listing时,关联的primary listing plan id
	CreatedAt              time.Time `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP;NOT NULL" json:"created_at"`                          // 创建时间
	UpdatedAt              time.Time `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP;NOT NULL" json:"updated_at"`                          // 更新时间
	DevApplicationLocaleID uint64    `gorm:"column:dev_application_locale_id;type:bigint(20) unsigned;default:0;NOT NULL" json:"dev_application_locale_id"` // locale id
	ApplicationID          uint64    `gorm:"column:application_id;type:bigint(20) unsigned;default:0;NOT NULL" json:"application_id"`                       // 应用ID
}

func (s DevApplicationPricingPlan) TableName() string {
	return "dev_application_pricing_plan"
}

// 应用收费主表
type ApplicationPricing struct {
	ID                  uint64    `gorm:"column:id;type:bigint(20) unsigned;primary_key" json:"id"`                                           // 主键ID
	ChargeType          string    `gorm:"column:charge_type;type:varchar(30);NOT NULL" json:"charge_type"`                                    // app 收费类型
	CreatedAt           time.Time `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP;NOT NULL" json:"created_at"`               // 创建时间
	UpdatedAt           time.Time `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP;NOT NULL" json:"updated_at"`               // 更新时间
	ApplicationLocaleID string    `gorm:"column:application_locale_id;type:varbinary(16);default:0x30;NOT NULL" json:"application_locale_id"` // listing locale id
	ApplicationID       uint64    `gorm:"column:application_id;type:bigint(20) unsigned;default:0;NOT NULL" json:"application_id"`            // 应用ID
}

func (s ApplicationPricing) TableName() string {
	return "application_pricing"
}

// 收费套餐表
type ApplicationPricingPlan struct {
	ID                  uint64    `gorm:"column:id;type:bigint(20) unsigned;primary_key" json:"id"`                                           // 主键ID
	PlanName            string    `gorm:"column:plan_name;type:varchar(100);NOT NULL" json:"plan_name"`                                       // 套餐名字
	Features            string    `gorm:"column:features;type:varchar(255);NOT NULL" json:"features"`                                         // 套餐特点
	Price               string    `gorm:"column:price;type:decimal(19,2);default:0.00;NOT NULL" json:"price"`                                 // 月度套餐价格
	PriceType           string    `gorm:"column:price_type;type:varchar(20);default:'monthly';NOT NULL" json:"price_type"`                    // 价格类型： free  monthly  yearly  one-time
	TrailDays           int       `gorm:"column:trail_days;type:int(11);default:0;NOT NULL" json:"trail_days"`                                // 试用期天数
	CreatedAt           time.Time `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP;NOT NULL" json:"created_at"`               // 创建时间
	UpdatedAt           time.Time `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP;NOT NULL" json:"updated_at"`               // 更新时间
	ApplicationLocaleID string    `gorm:"column:application_locale_id;type:varbinary(16);default:0x30;NOT NULL" json:"application_locale_id"` // locale id
	ApplicationID       uint64    `gorm:"column:application_id;type:bigint(20) unsigned;default:0;NOT NULL" json:"application_id"`            // 应用ID
	ParentPlanID        uint64    `gorm:"column:parent_plan_id;type:bigint(20) unsigned;default:0;NOT NULL" json:"parent_plan_id"`            // 关联parent plan id
}

func (s ApplicationPricingPlan) TableName() string {
	return "application_pricing_plan"
}
