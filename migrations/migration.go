package migrations

import (
	"fmt"
	"sort"

	"github.com/jinzhu/gorm"
	"github.com/kikyousky/gormigrate"
)

var Migrations []*gormigrate.Migration

func Migrate(db *gorm.DB) error {
	if len(Migrations) == 0 {
		return nil
	}
	sort.Slice(Migrations, func(i, j int) bool {
		n := Migrations
		return n[i].ID < n[j].ID
	})
	gormigrate.DefaultOptions.IDColumnSize = 250
	m := gormigrate.New(
		db, /*.LogMode(true)*/
		gormigrate.DefaultOptions,
		Migrations,
	)
	return m.Migrate()
}

func IndexExist(db *gorm.DB, table, index string) (error, bool) {
	var Indexs []struct {
		KeyName string `gorm:"column:Key_name"`
	}

	err := db.Raw(fmt.Sprintf("show index from %s", table)).Find(&Indexs).Error
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return err, false
	}

	for _, item := range Indexs {
		if item.KeyName == index {
			return nil, true
		}
	}
	return nil, false
}
